package ratelimiter

import (
	"context"
	"fmt"
	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"strings"
	"sync"
	"time"

	"github.com/eosnationftw/dauth/ratelimiter"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

func init() {
	// "olric://local?rates=search:60,block:60,blockmeta:60,token:60&whitelist=1.2.3.4,1.1.1.1"
	ratelimiter.Register("olric", func(configURL string) (ratelimiter.RateLimiter, error) {
		zlog.Info("parsing rate limiter settings", zap.String("url", configURL))
		olricPeers, userRateLimits, whitelistedIps, err := parseURL(configURL)
		if err != nil {
			return nil, fmt.Errorf("olric factory: %w", err)
		}
		serviceNames := ratelimiter.GetServices()
		err = validateServices(userRateLimits, serviceNames)
		if err != nil {
			return nil, fmt.Errorf("olric factory: %w", err)
		}

		zlog.Info("setting up rate limiter",
			zap.String("olric_peers", strings.Join(olricPeers, ",")),
			zap.Reflect("rate_limits", userRateLimits),
		)

		var olricConfig *config.Config

		if len(olricPeers) == 1 && olricPeers[0] == "local" {
			olricConfig = config.New("local")
		} else {
			olricConfig = config.New("lan")
			olricConfig.Peers = olricPeers
		}

		// Callback function. It's called when this node is ready to accept connections.
		ctx, cancel := context.WithCancel(context.Background())
		olricConfig.Started = func() {
			defer cancel()
			zlog.Info("olric is ready to accept connections")
		}

		db, err := olric.New(olricConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create olric instance: %w", err)
		}

		go func() {
			// Call Start at background. It's a blocker call.
			err = db.Start()
			if err != nil {
				zlog.Fatal("olric.Start returned an error", zap.Error(err))
			}
		}()

		<-ctx.Done()

		olricDMap, err := db.NewDMap("rate_limits")
		if err != nil {
			return nil, fmt.Errorf("failed to create dmap: %w", err)
		}

		requestRateLimiter := NewRequestRateLimiter(db, olricDMap, userRateLimits, whitelistedIps)
		return requestRateLimiter, nil
	})
}

func NewRequestRateLimiter(olricClient *olric.Olric, olricDMap *olric.DMap, limits map[string]int64, whitelist map[string]bool) *RequestRateLimiter {
	return &RequestRateLimiter{
		olricClient: olricClient,
		olricDMap:   olricDMap,
		limits:      limits,
		whitelist:   whitelist,
		counters:    make(map[string]*requestCounter),
	}
}

type RequestRateLimiter struct {
	prefix      string
	olricClient *olric.Olric
	olricDMap   *olric.DMap
	counters    map[string]*requestCounter
	mutex       sync.RWMutex
	limits      map[string]int64
	whitelist   map[string]bool
}

type requestCounter struct {
	uid       string
	prefix    string
	olricDMap *olric.DMap
	mutex     sync.Mutex

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
		res, err := c.olricDMap.Incr(currKey, int(localCount))
		// fmt.Printf("increased %s by %d - res: %d", currKey, localCount, res)
		if err != nil {
			return false, err
		}
		err = c.olricDMap.Expire(currKey, 3*time.Minute)
		if err != nil {
			return false, err
		}

		// todo remove
		zlog.Info("incremented counter",
			zap.String("key", currKey),
			zap.Int64("by", localCount),
			zap.Int("to", res),
		)

		c.remoteCountCurrent = int64(res)
	} else {
		curr, err := c.olricDMap.Get(currKey)
		if err != nil && err != olric.ErrKeyNotFound {
			return false, err
		} else if err == olric.ErrKeyNotFound {
			c.remoteCountCurrent = 0
		} else {
			c.remoteCountCurrent = int64(curr.(int))
		}
	}

	prev, err := c.olricDMap.Get(prevKey)
	if err != nil && err != olric.ErrKeyNotFound {
		return false, err
	} else if err == olric.ErrKeyNotFound {
		c.remoteCountPrevious = 0
	} else {
		c.remoteCountPrevious = int64(prev.(int))
	}
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

	res, err := r.olricDMap.Get(currKey)
	if err != nil {
		zlog.Debug("olric error getting current", zap.Error(err), zap.String("uid", uid))
	} else {
		cur = int64(res.(int))
	}

	res, err = r.olricDMap.Get(prevKey)
	if err != nil {
		zlog.Debug("olric error getting previous", zap.Error(err), zap.String("uid", uid))
	} else {
		prev = int64(res.(int))
	}

	c = &requestCounter{
		remoteCountCurrent:    cur,
		remoteCountPrevious:   prev,
		remotePreviousSeconds: secs,
		localCount:            atomic.NewInt64(0),
		prefix:                r.prefix,
		olricDMap:             r.olricDMap,
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

	if lastMinCount < limit || r.whitelist[id] {
		c.localCount.Inc()
		allow = true
	}

	// todo remove
	zlog.Info("gate called",
		zap.Int64("last_min_count", lastMinCount),
		zap.Int64("limit", limit),
		zap.Bool("allowed", allow),
		zap.String("id", id),
		zap.String("uid", uid),
		zap.Bool("whitelisted", r.whitelist[id]),
	)

	return
}
