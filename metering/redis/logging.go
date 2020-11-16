package redis

import (
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func init() {
	logging.Register("github/dfuse-io/dauth/metering/redis", &zlog)
}
