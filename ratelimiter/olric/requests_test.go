package ratelimiter

import (
	"context"
	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"go.uber.org/atomic"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestRateLimiterSmoketest(t *testing.T) {

	refTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	olricClient := config.New("local")

	ctx, cancel := context.WithCancel(context.Background())
	olricClient.Started = func() {
		defer cancel()
		log.Println("[INFO] Olric is ready to accept connections")
	}

	db, err := olric.New(olricClient)
	if err != nil {
		log.Fatalf("Failed to create Olric instance: %v", err)
	}

	go func() {
		// Call Start at background. It's a blocker call.
		err = db.Start()
		if err != nil {
			log.Fatalf("olric.Start returned an error: %v", err)
		}
	}()

	<-ctx.Done()

	olricDMap, err := db.NewDMap("rate_limits")
	require.NoError(t, err)

	rl := NewRequestRateLimiter(db, olricDMap, map[string]int64{
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

	keys := []string{
		"RCC::alice:test_req:200911172034",
		"RCC::bob:other_req:200911172034",
		"RCC::bob:test_req:200911172034",
	}

	for _, key := range keys {
		_, err := olricDMap.Get(key)
		require.NoError(t, err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = db.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown Olric: %v", err)
	}
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
				"RCC::bob:fetch:200911172034": 2,
			},
			uid:                   "bob:fetch",
			expectedCountCurrent:  2,
			expectedCountPrevious: 0,
		},
		{
			name: "existingPrevious",
			redisKeys: map[string]int64{
				"RCC::bob:fetch:200911172033": 2,
			},
			uid:                   "bob:fetch",
			expectedCountCurrent:  0,
			expectedCountPrevious: 2,
		},
		{
			name: "existingBoth",
			redisKeys: map[string]int64{
				"RCC::bob:fetch:200911172033": 2,
				"RCC::bob:fetch:200911172034": 3,
			},
			uid:                   "bob:fetch",
			expectedCountCurrent:  3,
			expectedCountPrevious: 2,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			olricClient := config.New("local")

			ctx, cancel := context.WithCancel(context.Background())
			olricClient.Started = func() {
				defer cancel()
				log.Println("[INFO] Olric is ready to accept connections")
			}

			db, err := olric.New(olricClient)
			if err != nil {
				log.Fatalf("Failed to create Olric instance: %v", err)
			}

			go func() {
				// Call Start at background. It's a blocker call.
				err = db.Start()
				if err != nil {
					log.Fatalf("olric.Start returned an error: %v", err)
				}
			}()

			<-ctx.Done()

			olricDMap, err := db.NewDMap(c.name)
			require.NoError(t, err)

			rl := NewRequestRateLimiter(db, olricDMap, nil)
			for k, v := range c.redisKeys {
				err := olricDMap.PutEx(k, v, 3*time.Minute)
				require.NoError(t, err)
			}
			rq := rl.getCounter(c.uid, currTime)
			assert.Equal(t, c.expectedCountCurrent, rq.remoteCountCurrent)
			assert.Equal(t, c.expectedCountPrevious, rq.remoteCountPrevious)

			ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
			err = db.Shutdown(ctx)
			if err != nil {
				log.Printf("Failed to shutdown Olric: %v", err)
			}
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
			rl := NewRequestRateLimiter(nil, nil, limits)
			rl.counters = make(map[string]*requestCounter)
			for _, rq := range c.reqCounters {
				rl.counters[rq.uid] = rq
			}
			assert.Equal(t, c.expectedPass, rl.Gate(c.id, c.method))
		})
	}
}
