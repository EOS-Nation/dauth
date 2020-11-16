package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name                 string
		url                  string
		expectRedisAddr      string
		expectUserRateLimits map[string]int64
		expectErr            bool
	}{
		{
			name:            "Golden path",
			url:             "cloud-gcp://redis.default.svc.cluster.local?rates=search:60,block:60,blockmeta:40,token:20",
			expectRedisAddr: "redis.default.svc.cluster.local",
			expectUserRateLimits: map[string]int64{
				"search":    60,
				"block":     60,
				"blockmeta": 40,
				"token":     20,
			},
		},
		{
			name:      "invalid url",
			url:       "redis.default.svc.cluster.local?rates=search:60,block:60,blockmeta:60,token:60",
			expectErr: true,
		},
		{
			name:                 "rate limits are optionals",
			url:                  "cloud-gcp://redis.default.svc.cluster.local",
			expectRedisAddr:      "redis.default.svc.cluster.local",
			expectUserRateLimits: map[string]int64{},
		},
		{
			name:                 "rate limits are optionals",
			url:                  "cloud-gcp://redis.default.svc.cluster.local?rates=",
			expectRedisAddr:      "redis.default.svc.cluster.local",
			expectUserRateLimits: map[string]int64{},
		},
		{
			name:      "invalid rates param",
			url:       "cloud-gcp://redis.default.svc.cluster.local?rates=search",
			expectErr: true,
		},
		{
			name:            "redis path with port",
			url:             "cloud-gcp://redis.default.svc.cluster.local:6379?rates=search:60,block:60,blockmeta:60,token:60",
			expectRedisAddr: "redis.default.svc.cluster.local:6379",
			expectUserRateLimits: map[string]int64{
				"search":    60,
				"block":     60,
				"blockmeta": 60,
				"token":     60,
			},
		},
		{
			name:            "redis path with ip and port",
			url:             "cloud-gcp://192.168.1.1:6379?rates=search:60,block:60,blockmeta:60,token:60",
			expectRedisAddr: "192.168.1.1:6379",
			expectUserRateLimits: map[string]int64{
				"search":    60,
				"block":     60,
				"blockmeta": 60,
				"token":     60,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			redisAddr, userRateLimits, err := parseURL(test.url)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectRedisAddr, redisAddr)
				assert.Equal(t, test.expectUserRateLimits, userRateLimits)
			}
		})
	}
}

func TestValidateServices(t *testing.T) {

	tests := []struct {
		name           string
		userRateLimits map[string]int64
		serviceNames   []string
		expectErr      bool
	}{
		{
			name: "golden path",
			userRateLimits: map[string]int64{
				"search":    60,
				"block":     60,
				"blockmeta": 60,
				"tokenmeta": 60,
			},
			serviceNames: []string{
				"search", "block", "blockmeta", "tokenmeta",
			},
			expectErr: false,
		},
		{
			name: "invalid service name",
			userRateLimits: map[string]int64{
				"broken":    60,
				"block":     60,
				"blockmeta": 60,
				"tokenmeta": 60,
			},
			serviceNames: []string{
				"search", "block", "blockmeta", "tokenmeta",
			},
			expectErr: true,
		},
		{
			name: "not all services rate limited",
			userRateLimits: map[string]int64{
				"block":     60,
				"blockmeta": 60,
				"tokenmeta": 60,
			},
			serviceNames: []string{
				"search", "block", "blockmeta", "tokenmeta",
			},
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateServices(test.userRateLimits, test.serviceNames)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
