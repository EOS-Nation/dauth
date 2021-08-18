package dredd

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streamingfast/dauth/dredd/keyer"
)

type DB struct {
	redisClient *redis.Client
}

func NewDB(redisClient *redis.Client) *DB {
	return &DB{
		redisClient: redisClient,
	}
}

/*
func (d *DB) CacheUserQuotaWithDurations(userID string, apiKeyID string, documentQuota int64, ipQuota int64, duration time.Duration) error {
	key := keyer.UserQuotaCacheKey(userID, apiKeyID)

	status := d.redisClient.ZAdd(context.Background(), key,
		&redis.Z{
			Score:  float64(documentQuota),
			Member: "docs",
		},
		&redis.Z{
			Score:  float64(ipQuota),
			Member: "ips",
		},
	)

	if status.Err() != nil {
		return fmt.Errorf("failed to cache user quota: %w", status.Err())
	}

	expStatus := d.redisClient.Expire(context.Background(), key, duration)
	if expStatus.Err() != nil {
		return fmt.Errorf("failed to set expiration on key %s: %w", key, expStatus.Err())
	}

	return nil
}

func (d *DB) CacheUserQuota(userID string, apiKeyID string, documentQuota int64, ipQuota int64) error {
	return d.CacheUserQuotaWithDurations(userID, apiKeyID, documentQuota, ipQuota, 24*time.Hour)
}

var QuotaNotFoundErr = errors.New("quota not found")

func (d *DB) UserQuota(userID string, apiKeyID string) (docs int64, ips int64, err error) {
	key := keyer.UserQuotaCacheKey(userID, apiKeyID)

	result := d.redisClient.ZRangeWithScores(context.Background(), key, 0, -1)
	if result.Err() != nil {
		return 0, 0, fmt.Errorf("failed to retreive user quota for key : %s: %w", key, result.Err())
	}

	if len(result.Val()) == 0 {
		return 0, 0, QuotaNotFoundErr
	}

	docs = 0
	ips = 0

	for _, z := range result.Val() {
		if z.Member == "docs" {
			docs = int64(z.Score)

		} else if z.Member == "ips" {
			ips = int64(z.Score)
		}
	}
	return
}

func (d *DB) EvictUserFromQuotaCache(userID string) error {
	prefix := keyer.UserQuotaCachePrefix(userID)

	s := `
	local prefix = ARGV[1] .. "*"
	local keys = redis.call('keys', prefix)
	if table.getn(keys) == 0 then return 0 end
	return redis.call('del', unpack(keys))
`
	status := d.redisClient.Eval(context.Background(), s, []string{}, prefix)
	if status.Err() != nil {
		return fmt.Errorf("failed to evict userID %s from user quota cache : %w", userID, status.Err())
	}

	return nil
}
*/

func (d *DB) BlackListUser(userID string, reason string, duration time.Duration) error {
	key := keyer.UserIDBlackListKey(userID)
	status := d.redisClient.Set(context.Background(), key, reason, duration)
	if status.Err() != nil {
		return fmt.Errorf("failed to black list user: %w", status.Err())
	}
	return nil
}

func (d *DB) UnBlackListUser(userID string) error {
	key := keyer.UserIDBlackListKey(userID)
	status := d.redisClient.Del(context.Background(), key)
	if status.Err() != nil {
		return fmt.Errorf("failed to black list user: %w", status.Err())
	}

	_, err := d.BlackListVersionIncr()
	if status.Err() != nil {
		return err
	}

	return nil
}

func (d *DB) IsUserBlackListed(userID string) (bool, error) {
	key := keyer.UserIDBlackListKey(userID)
	result := d.redisClient.Exists(context.Background(), key)
	if result.Err() != nil {
		return false, fmt.Errorf("failed to check if user black listed: %w", result.Err())
	}

	return result.Val() == 1, nil
}

func (d *DB) UserBlackListedDetails(userID string) (string, time.Duration, error) {
	key := keyer.UserIDBlackListKey(userID)
	result := d.redisClient.Get(context.Background(), key)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return "", 0, nil
		}
		return "", 0, fmt.Errorf("failed to check if user black listed: %w", result.Err())
	}

	return result.Val(), d.redisClient.TTL(context.Background(), key).Val(), nil
}

func (d *DB) BlackListedUsers() (userIDs []string, err error) {
	result := d.redisClient.Keys(context.Background(), keyer.BLACK_LIST_USER_KEY_PREFIX+"*")
	if result.Err() != nil {
		return nil, fmt.Errorf("failed retreive black listed users: %w", result.Err())
	}

	for _, v := range result.Val() {
		parts := strings.Split(v, keyer.DELIMITER)

		if len(parts) > 1 && parts[1] != "VERSION" {
			userIDs = append(userIDs, parts[1])
		}
	}

	return
}

func (d *DB) UserBlackListVersion() (version int, err error) {
	key := keyer.UserBlackListVersionKey()
	result := d.redisClient.Get(context.Background(), key)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return 0, nil
		}
		return -1, fmt.Errorf("failed retreive black list version: %w", result.Err())
	}
	return result.Int()
}

func (d *DB) BlackListVersionIncr() (version int64, err error) {
	key := keyer.UserBlackListVersionKey()
	result := d.redisClient.Incr(context.Background(), key)
	if result.Err() != nil {
		return -1, fmt.Errorf("failed incremewnt black list version: %w", result.Err())
	}
	return result.Val(), nil
}

type DocumentCount struct {
	Key             string
	Count           int64
	CumulativeCount int64
}

func (d *DB) UserDocumentCounts(userID string) ([]*DocumentCount, int64, error) {
	var out []*DocumentCount
	keys := keyer.DocumentConsumptionLast30Days(userID, time.Now())
	keys = reverse(keys)
	result := d.redisClient.MGet(context.Background(), keys...)
	if result.Err() != nil {
		return nil, 0, result.Err()
	}
	cumulativeCount := int64(0)
	for i, v := range result.Val() {
		dc := &DocumentCount{Key: keys[i]}
		if v == nil {
			dc.Count = 0
		} else {
			count, err := strconv.Atoi(v.(string))
			if err != nil {
				return nil, 0, err
			}
			dc.Count = int64(count)
		}
		cumulativeCount += dc.Count
		dc.CumulativeCount = cumulativeCount
		out = append(out, dc)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Key > out[j].Key
	})

	return out, cumulativeCount, nil
}

func (d *DB) TTL(key string) (time.Duration, error) {
	result := d.redisClient.TTL(context.Background(), key)
	if result.Err() != nil {
		return 0, fmt.Errorf("failed to retreive ttl for key %s : %w", key, result.Err())
	}

	return result.Val(), nil
}

type MStat struct {
	TotalDocumentCount       int64
	DocumentCountPerApiKeyID map[string]int64
	Key                      string
}

func (d *DB) MStats() ([]*MStat, error) {
	script := `
local data = {}
local cursor = 0
repeat

    local t = redis.call("SCAN", cursor, "MATCH", "DCCP:*", "COUNT", 20000);
    local keys = t[2];
	local counts = redis.call("MGET", unpack(keys));
	for i = 1, #counts do
		table.insert(data, keys[i])
		table.insert(data, counts[i])
	end
    cursor = t[1];
until cursor == "0";

return data
`
	result := d.redisClient.Eval(context.Background(), script, []string{})
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to retreive mstats: %w", result.Err())
	}

	var out []*MStat
	stats := result.Val().([]interface{})

	for i := 0; i < len(stats); i += 2 {
		k := stats[i].(string)

		v, err := strconv.Atoi(stats[i+1].(string))
		if err != nil {
			return nil, err
		}
		mstat := &MStat{
			Key:                      k,
			TotalDocumentCount:       int64(v),
			DocumentCountPerApiKeyID: map[string]int64{},
		}
		out = append(out, mstat)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TotalDocumentCount > out[j].TotalDocumentCount
	})

	return out, nil
}

func reverse(a []string) []string {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}
