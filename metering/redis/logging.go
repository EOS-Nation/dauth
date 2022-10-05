package redis

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var traceEnabled = logging.IsTraceEnabled("dauth", "github.com/streamingfast/dauth/metering/redis")
var zlog *zap.Logger

func init() {
	logging.Register("github.com/streamingfast/dauth/metering/redis", &zlog)
}
