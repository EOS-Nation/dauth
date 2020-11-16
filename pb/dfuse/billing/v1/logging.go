package pbbilling

import (
	"github.com/dfuse-io/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func init() {
	logging.Register("github/dfuse-io/dauth/pb/dfuse/billing/v1", &zlog)
}
