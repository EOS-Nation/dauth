package dredd

import (
	"context"
	"fmt"
	"github.com/dfuse-io/dauth/dredd/lua"
	pbbilling "github.com/dfuse-io/dauth/pb/dfuse/billing/v1"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/dfuse-io/dauth/dredd/keyer"
	"github.com/dfuse-io/dmetrics"
)

var Metrics = dmetrics.NewSet()
var LuaExecutionTimeSumMetric = Metrics.NewCounter("rate_limit_lua_execution_time_sum_micro", "rate limit lua script execution time")
var LuaExecutionCountMetric = Metrics.NewCounter("rate_limit_lua_execution_count", "number of rate limit lua script execution")

type LuaEventHandler struct {
	redisClient *redis.Client
	scriptSHA1  string
}

func NewLuaEventHandler(redisClient *redis.Client) (*LuaEventHandler, error) {
	data, err := lua.Asset("cutoff.lua")
	if err != nil {
		return nil, fmt.Errorf("bin data err: %w", err)
	}

	loadResult := redisClient.ScriptLoad(context.Background(), string(data))
	if loadResult.Err() != nil {
		return nil, fmt.Errorf("failed to upload lua script: %w", loadResult.Err())
	}
	scriptSHA1 := loadResult.Val()
	//fmt.Println("scriptSHA1:", scriptSHA1)
	return &LuaEventHandler{
		redisClient: redisClient,
		scriptSHA1:  scriptSHA1,
	}, nil
}

func (l *LuaEventHandler) Test() {

}

func (l *LuaEventHandler) HandleEvent(ev *pbbilling.Event, docQuota int) (bool, error) {

	keys := []string{
		keyer.CurrentPeriodDocumentConsumption(ev.UserId),
		keyer.UserIDBlackListKey(ev.UserId),
		keyer.UserBlackListVersionKey(),
		keyer.DocumentConsumptionDaily(ev.UserId, time.Now()),
		keyer.DocumentConsumptionMinutely(ev.UserId, time.Now()),
	}
	// keys = append(keys, keyer.DocumentConsumptionLast30Days(ev.UserId, time.Now())...)
	keys = append(keys, keyer.DocumentConsumptionLastWindow(ev.UserId, 10, time.Now())...)

	now := time.Now()
	// end of window is the end of the current minute, afterwards we want to calculate a new 10 min moving average
	endOfWindow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 59, 999999, time.UTC)

	start := time.Now()

	blacklisted := false

	result := l.redisClient.EvalSha(context.Background(), l.scriptSHA1, keys, docQuota, ev.ResponsesCount, endOfWindow.Unix(), 600, 10)
	luaRespStr, err := result.Result()
	if err == nil {
		blacklisted = (luaRespStr == "bl")
	}

	if result.Err() != nil {
		return blacklisted, fmt.Errorf("failed to eval rate limit script: %w", result.Err())
	}
	LuaExecutionTimeSumMetric.AddInt64(time.Since(start).Microseconds())
	LuaExecutionCountMetric.Inc()

	return blacklisted, nil
}
