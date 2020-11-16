package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_parseURL(t *testing.T) {
	cfgURL := "cloud-gcp://projects/eoscanada-public/locations/global/keyRings/eosws-api-auth/cryptoKeys/default/cryptoKeyVersions/1?quotaEnforce=true&quotaRedisAddr=redis.default.svc.cluster.local:6379&quotaBlacklistUpdateInterval=3s"

	kmsKeyPath, enforceQuota, redisAddr, quotaBlacklistUpdateInterval, err := parseURL(cfgURL)

	require.NoError(t, err)
	require.Equal(t, "projects/eoscanada-public/locations/global/keyRings/eosws-api-auth/cryptoKeys/default/cryptoKeyVersions/1", kmsKeyPath)
	require.True(t, enforceQuota)
	require.Equal(t, "redis.default.svc.cluster.local:6379", redisAddr)
	require.Equal(t, 3*time.Second, quotaBlacklistUpdateInterval)
}
