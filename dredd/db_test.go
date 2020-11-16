package dredd

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/dfuse-io/dauth/dredd/keyer"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB_CacheUserAllocation(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.CacheUserQuota("uid:user.1", "api.key.1", 1, 2)
	require.NoError(t, err)

	d, i, err := db.UserQuota("uid:user.1", "api.key.1")
	require.NoError(t, err)

	assert.Equal(t, int64(1), d)
	assert.Equal(t, int64(2), i)
	//todo : check ttl
	result := cli.TTL(keyer.UserQuotaCacheKey("uid:user.1", "api.key.1"))
	fmt.Println("ttl:", result.Val())
}

func TestDB_EvictUserFromAllocationCache(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.CacheUserQuota("uid:user.1", "api.key.1", 1, 2)
	require.NoError(t, err)
	err = db.CacheUserQuota("uid:user.2", "api.key.2", 1, 2)
	require.NoError(t, err)

	err = db.EvictUserFromQuotaCache("uid:user.1")
	require.NoError(t, err)

	key := keyer.UserQuotaCacheKey("uid:user.1", "api.key.1")
	result := cli.Exists(key)
	require.False(t, result.Val() == 1)

	key = keyer.UserQuotaCacheKey("uid:user.2", "api.key.2")
	result = cli.Exists(key)
	require.True(t, result.Val() == 1)
}

func TestDB_EvictUserNotInCache(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.EvictUserFromQuotaCache("uid:user.1")
	require.NoError(t, err)
}

func TestDB_UserAllocationNotInCache(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	d, i, err := db.UserQuota("uid:user.1", "api.key.1")
	require.Equal(t, QuotaNotFoundErr, err)
	assert.Equal(t, int64(0), d)
	assert.Equal(t, int64(0), i)
}

func TestDB_BlackListUser(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.BlackListUser("uid:user.1", "reason.1", time.Second*1)
	require.NoError(t, err)

	result := cli.Exists(keyer.UserIDBlackListKey("uid:user.1"))
	require.True(t, result.Val() == 1)
}

func TestDB_UnBlackListUser(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.BlackListUser("uid:user.1", "reason.1", time.Second*1)
	require.NoError(t, err)

	result := cli.Exists(keyer.UserIDBlackListKey("uid:user.1"))
	require.True(t, result.Val() == 1)

	err = db.UnBlackListUser("uid:user.1")
	require.NoError(t, err)

	result = cli.Exists(keyer.UserIDBlackListKey("uid:user.1"))
	require.False(t, result.Val() == 1)
}

func TestDB_BlackListedUsers(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.BlackListUser("uid:user.1", "reason.1", time.Second*1)
	require.NoError(t, err)
	err = db.BlackListUser("uid:user.2", "reason.1", time.Second*1)
	require.NoError(t, err)
	err = db.BlackListUser("uid:user.3", "reason.1", time.Second*1)
	require.NoError(t, err)

	ids, err := db.BlackListedUsers()
	require.NoError(t, err)

	require.Equal(t, []string{"user.1", "user.2", "user.3"}, ids)

}

func TestDB_IsUserBlackListed(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	err := db.BlackListUser("uid:user.1", "reason.1", time.Second*1)
	require.NoError(t, err)

	blackListed, err := db.IsUserBlackListed("uid:user.1")
	require.NoError(t, err)
	require.True(t, blackListed)

	blackListed, err = db.IsUserBlackListed("uid:user.2")
	require.NoError(t, err)
	require.False(t, blackListed)
}

func TestDB_UserBlackListVersion(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	version, err := db.UserBlackListVersion()
	require.NoError(t, err)
	assert.Equal(t, 0, version)

	cli.Incr(keyer.UserBlackListVersionKey())
	cli.Incr(keyer.UserBlackListVersionKey())
	version, err = db.UserBlackListVersion()
	require.NoError(t, err)
	assert.Equal(t, 2, version)

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

func TestDB_MStats(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	cli.Incr("DCCP:uid:0beti479656cd07513488")

	stats, err := db.MStats()
	require.NoError(t, err)
	require.Equal(t, int64(1), stats[0].TotalDocumentCount)
}

func TestDB_BlackListVersionIncr(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	v, err := db.BlackListVersionIncr()
	require.NoError(t, err)
	require.Equal(t, int64(1), v)
}

func TestDB_UserDocumentCounts(t *testing.T) {
	cli := newTestClient()
	db := NewDB(cli)

	expectedData := []struct{
		ExpectedKey string
		ExpectedCount int64
		ExpectedCumulativeCount int64
	}{
		{
			ExpectedKey:   keyer.DocumentConsumptionDaily("uid:0butyf99d12f03093f3ca", time.Now()),
			ExpectedCount: 3,
			ExpectedCumulativeCount: 10,
		},
		{
			ExpectedKey:   keyer.DocumentConsumptionDaily("uid:0butyf99d12f03093f3ca", time.Now().Add((-1 * 24 * time.Hour))),
			ExpectedCount: 5,
			ExpectedCumulativeCount: 7,
		},
		{
			ExpectedKey:   keyer.DocumentConsumptionDaily("uid:0butyf99d12f03093f3ca", time.Now().Add((-2 * 24 * time.Hour))),
			ExpectedCount: 2,
			ExpectedCumulativeCount: 2,
		},
	}

	for _, d := range expectedData {
		cli.Set(d.ExpectedKey, d.ExpectedCount, -1)
	}

	counts, total,  err := db.UserDocumentCounts("uid:0butyf99d12f03093f3ca")
	require.NoError(t, err)
	assert.Equal(t, expectedData[0].ExpectedCumulativeCount, total)
	for i, d := range expectedData {
		assert.Equal(t, d.ExpectedCount, counts[i].Count)
		assert.Equal(t, d.ExpectedCumulativeCount, counts[i].CumulativeCount)
		assert.Equal(t, d.ExpectedKey, counts[i].Key)
	}

}
