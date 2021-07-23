package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eosnationftw/dauth/ratelimiter"
	"github.com/go-redis/redis/v8"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

func init() {
	// cloud-gcp://redis.default.svc.cluster.local?rates=search:60,block:60,blockmeta:60,token:60"
	ratelimiter.Register("cloud-gcp", func(configURL string) (ratelimiter.RateLimiter, error) {
		zlog.Info("parsing rate limiter settings", zap.String("url", configURL))
		redisAddr, userRateLimits, err := parseURL(configURL)
		if err != nil {
			return nil, fmt.Errorf("cloud-gcp factory: %w", err)
		}
		serviceNames := ratelimiter.GetServices()
		err = validateServices(userRateLimits, serviceNames)
		if err != nil {
			return nil, fmt.Errorf("cloud-gcp factory: %w", err)
		}

		zlog.Info("setting up rate limiter",
			zap.String("redis_addr", redisAddr),
			zap.Reflect("rate_limits", userRateLimits),
		)

		redisClient := redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})

		requestRateLimiter := NewRequestRateLimiter("USR", redisClient, userRateLimits)
		return requestRateLimiter, nil
	})
}

func NewRequestRateLimiter(prefix string, redisClient *redis.Client, limits map[string]int64) *RequestRateLimiter {
	return &RequestRateLimiter{
		prefix:      prefix,
		redisClient: redisClient,
		limits:      limits,
		counters:    make(map[string]*requestCounter),
	}
}

type RequestRateLimiter struct {
	prefix      string
	redisClient *redis.Client
	counters    map[string]*requestCounter
	mutex       sync.RWMutex
	limits      map[string]int64
}

type requestCounter struct {
	uid         string
	prefix      string
	redisClient *redis.Client
	mutex       sync.Mutex

	remoteCountCurrent    int64
	remoteCountPrevious   int64
	remotePreviousSeconds int // number of seconds out of 60 to consider from previous minute

	localCount  *atomic.Int64
	cleanupFunc func()
}

func (c *requestCounter) updateRemoteCounter(at time.Time) (empty bool, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	currMinute, prevMinute, secs := minuteTimestamps(at)
	currKey := requestConsumptionCounter(c.prefix, c.uid, currMinute)
	prevKey := requestConsumptionCounter(c.prefix, c.uid, prevMinute)

	localCount := c.localCount.Swap(0)

	if localCount > 0 {
		ret := c.redisClient.IncrBy(context.Background(), currKey, localCount)
		res, err := ret.Result()
		if err != nil {
			return false, err
		}
		c.redisClient.Expire(context.Background(), currKey, 3*time.Minute)
		if err != nil {
			return false, err
		}

		c.remoteCountCurrent = res
	} else {
		curr, err := c.redisClient.Get(context.Background(), currKey).Int64()
		if err != nil && err != redis.Nil {
			return false, err
		}
		c.remoteCountCurrent = curr
	}

	prev, err := c.redisClient.Get(context.Background(), prevKey).Int64()
	if err != nil && err != redis.Nil {
		return false, err
	}
	c.remoteCountPrevious = prev
	c.remotePreviousSeconds = secs

	empty = c.remoteCountPrevious == 0 && c.remoteCountCurrent == 0
	return empty, nil

}

func (r *requestCounter) Launch(cleanup func()) {
	go func() {
		defer cleanup()
		for {
			time.Sleep(3 * time.Second)
			empty, err := r.updateRemoteCounter(time.Now().UTC())
			if err != nil {
				zlog.Warn("error updating counter", zap.Error(err))
			}
			if empty {
				return
			}
		}
	}()
}

const minuteTimeFmt = "200601021504"

func minuteTimestamps(now time.Time) (cur string, prev string, prevSeconds int) {
	remaining := now.Sub(now.Truncate(time.Minute)).Seconds()
	return now.Format(minuteTimeFmt), now.Add(-time.Minute).Format(minuteTimeFmt), 60 - int(remaining)
}

func estimateLastMinuteCount(currMinuteCount, prevMinuteCount int64, prevMinuteSignificantSeconds int, localCount *atomic.Int64) int64 {
	return currMinuteCount + (prevMinuteCount * int64(prevMinuteSignificantSeconds) / 60) + localCount.Load()
}

func (r *RequestRateLimiter) getCounter(uid string, at time.Time) *requestCounter {
	r.mutex.RLock()
	c, ok := r.counters[uid]
	r.mutex.RUnlock()

	if ok {
		return c
	}

	currMinute, prevMinute, secs := minuteTimestamps(at)
	currKey := requestConsumptionCounter(r.prefix, uid, currMinute)
	prevKey := requestConsumptionCounter(r.prefix, uid, prevMinute)

	var cur int64
	var prev int64
	var err error

	cur, err = r.redisClient.Get(context.Background(), currKey).Int64()
	if err != nil {
		zlog.Debug("redis error getting current", zap.Error(err), zap.String("uid", uid))
	}

	prev, err = r.redisClient.Get(context.Background(), prevKey).Int64()
	if err != nil {
		zlog.Debug("redis error getting previous", zap.Error(err), zap.String("uid", uid))
	}

	c = &requestCounter{
		remoteCountCurrent:    cur,
		remoteCountPrevious:   prev,
		remotePreviousSeconds: secs,
		localCount:            atomic.NewInt64(0),
		prefix:                r.prefix,
		redisClient:           r.redisClient,
		uid:                   uid,
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	c.Launch(func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		delete(r.counters, uid)
	})

	r.counters[uid] = c
	return c
}

func (r *RequestRateLimiter) Gate(
	id string,
	method string) (allow bool) {
	/*
		1. Create the key (remote value key, &)
		2. retrieve value from local cache
		3. if value greater then limit 429 request
		4. let through and increment local key
	*/

	limit, ok := r.limits[method]
	if !ok {
		// allow request if the service is not rate limited
		return true
	}

	uid := fmt.Sprintf("%s:%s", id, method)

	c := r.getCounter(uid, time.Now().UTC())
	lastMinCount := estimateLastMinuteCount(c.remoteCountCurrent, c.remoteCountPrevious, c.remotePreviousSeconds, c.localCount)
	if lastMinCount < limit {
		c.localCount.Inc()
		allow = true
	}

	return
}
