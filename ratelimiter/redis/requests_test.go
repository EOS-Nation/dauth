package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func TestRequestRateLimiterSmoketest(t *testing.T) {

	refTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	cli := newTestClient()
	rl := NewRequestRateLimiter("TST", cli, map[string]int64{
		"test_req":  5,
		"other_req": 10,
	})

	for i := 0; i < 5; i++ {
		allowed := rl.Gate("bob", "test_req")
		require.True(t, allowed)
	}
	require.False(t, rl.Gate("bob", "test_req"))
	require.True(t, rl.Gate("bob", "other_req"))
	require.True(t, rl.Gate("alice", "test_req"))

	require.Equal(t, 3, len(rl.counters))
	for _, counter := range rl.counters {
		empty, err := counter.updateRemoteCounter(refTime)
		require.NoError(t, err)
		assert.False(t, empty)
	}

	keys := cli.Keys("RCC:TST*")
	require.NoError(t, keys.Err())

	require.EqualValues(t, []string{
		"RCC:TST:alice:test_req:200911172034",
		"RCC:TST:bob:other_req:200911172034",
		"RCC:TST:bob:test_req:200911172034",
	}, keys.Val())

}

func TestRequestRateLimiterCounter(t *testing.T) {

	currTime := time.Date(2009, 11, 17, 20, 34, 58, 0, time.UTC)

	//	limits := make(map[string]int64)
	//	limits["fetch"] = 5
	//	limits["play"] = 10

	tests := []struct {
		name                  string
		redisKeys             map[string]int64
		uid                   string
		expectedCountCurrent  int64
		expectedCountPrevious int64
	}{
		{
			name:                  "new",
			redisKeys:             nil,
			uid:                   "bob:fetch",
			expectedCountCurrent:  0,
			expectedCountPrevious: 0,
		},
		{
			name: "existingCurrent",
			redisKeys: map[string]int64{
				"RCC:TST:bob:fetch:200911172034": 2,
			},
			uid:                   "bob:fetch",
			expectedCountCurrent:  2,
			expectedCountPrevious: 0,
		},
		{
			name: "almostMatchingKeys",
			redisKeys: map[string]int64{
				"RCC:TST:bob:fetch:200911172034":   2,
				"RCC:NIL:bob:play:200911172034":    3,
				"RCC:TST:alice:fetch:200911172034": 3,
				"NIL:TST:bob:fetch:200911172034":   3,
			},
			uid:                   "bob:play",
			expectedCountCurrent:  0,
			expectedCountPrevious: 0,
		},
		{
			name: "existingPrevious",
			redisKeys: map[string]int64{
				"RCC:TST:bob:fetch:200911172033": 2,
			},
			uid:                   "bob:fetch",
			expectedCountCurrent:  0,
			expectedCountPrevious: 2,
		},
		{
			name: "existingBoth",
			redisKeys: map[string]int64{
				"RCC:TST:bob:fetch:200911172033": 2,
				"RCC:TST:bob:fetch:200911172034": 3,
			},
			uid:                   "bob:fetch",
			expectedCountCurrent:  3,
			expectedCountPrevious: 2,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			cli := newTestClient()
			rl := NewRequestRateLimiter("TST", cli, nil)
			for k, v := range c.redisKeys {
				cli.Set(k, v, 3*time.Minute)
			}
			rq := rl.getCounter(c.uid, currTime)
			assert.Equal(t, c.expectedCountCurrent, rq.remoteCountCurrent)
			assert.Equal(t, c.expectedCountPrevious, rq.remoteCountPrevious)

		})
	}
}

func TestRequestRateLimiterGate(t *testing.T) {

	//currTime := time.Date(2009, 11, 17, 20, 34, 58, 0, time.UTC)

	limits := make(map[string]int64)
	limits["fetch"] = 100
	limits["play"] = 1000

	tests := []struct {
		name         string
		reqCounters  []*requestCounter
		id           string
		method       string
		expectedPass bool
	}{
		{
			name: "pass",
			reqCounters: []*requestCounter{
				{
					uid:                 "bob:fetch",
					localCount:          atomic.NewInt64(1),
					remoteCountCurrent:  2,
					remoteCountPrevious: 3,
				},
			},
			expectedPass: true,
			id:           "bob",
			method:       "fetch",
		},
		{
			name: "blockCurrent",
			reqCounters: []*requestCounter{
				{
					uid:                 "bob:fetch",
					localCount:          atomic.NewInt64(1),
					remoteCountCurrent:  100,
					remoteCountPrevious: 0,
				},
			},
			expectedPass: false,
			id:           "bob",
			method:       "fetch",
		},
		{
			name: "blockLocal",
			reqCounters: []*requestCounter{
				{
					uid:                 "bob:fetch",
					localCount:          atomic.NewInt64(88),
					remoteCountCurrent:  12,
					remoteCountPrevious: 0,
				},
			},
			expectedPass: false,
			id:           "bob",
			method:       "fetch",
		},
		{
			name: "blockPrevious",
			reqCounters: []*requestCounter{
				{
					uid:                   "bob:fetch",
					localCount:            atomic.NewInt64(0),
					remoteCountCurrent:    0,
					remoteCountPrevious:   200,
					remotePreviousSeconds: 40,
				},
			},
			expectedPass: false,
			id:           "bob",
			method:       "fetch",
		},
		{
			name: "passPreviousPercent",
			reqCounters: []*requestCounter{
				{
					uid:                   "bob:fetch",
					localCount:            atomic.NewInt64(1),
					remoteCountCurrent:    1,
					remoteCountPrevious:   200,
					remotePreviousSeconds: 10,
				},
			},
			expectedPass: true,
			id:           "bob",
			method:       "fetch",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			rl := NewRequestRateLimiter("", nil, limits)

			rl.counters = make(map[string]*requestCounter)
			for _, rq := range c.reqCounters {
				rl.counters[rq.uid] = rq
			}
			assert.Equal(t, c.expectedPass, rl.Gate(c.id, c.method))
		})
	}
}

func newTestClient() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
}
