// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authenticator

import (
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func RealIPFromRequest(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	// todo remove
	zlog.Info("resolving ip address",
		zap.String("xForwardedFor", xForwardedFor),
	)

	return RealIP(xForwardedFor)
}

func RealIP(forwardIPs string) string {
	if forwardIPs != "" {
		addresses := strings.Split(forwardIPs, ",")
		if len(addresses) >= 2 {
			return strings.TrimSpace(addresses[len(addresses)-2])
		}
	}

	return "0.0.0.0"
}
