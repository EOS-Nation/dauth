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

package redis

import (
	"context"
	"fmt"
	"github.com/dfuse-io/dauth/dredd"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"

	"github.com/dfuse-io/dauth/authenticator"
	"github.com/go-redis/redis/v8"
)

func init() {
	// redis://redis1,redis2,redis3?quotaEnforce=true&jwtKey=abc123&quotaBlacklistUpdateInterval=5s
	authenticator.Register("redis", func(configURL string) (authenticator.Authenticator, error) {

		redisNodes, enforceQuota, jwtKey, quotaBlacklistUpdateInterval, whitelistedIps, err := parseURL(configURL)

		if err != nil {
			return nil, fmt.Errorf("redis auth factory: %w", err)
		}
		return newAuthenticator(redisNodes, enforceQuota, jwtKey, quotaBlacklistUpdateInterval, whitelistedIps), nil
	})
}

func parseURL(configURL string) (redisNodes []string, enforceQuota bool, jwtKey string, quotaBlacklistUpdateInterval time.Duration, whitelistedIps map[string]bool, err error) {
	urlObject, err := url.Parse(configURL)
	if err != nil {
		return
	}

	redisNodes = strings.Split(urlObject.Host, ",")
	if len(redisNodes) == 1 && redisNodes[0] == "" {
		err = fmt.Errorf("missing redis nodes")
		return
	} else {
		for _, redisNode := range redisNodes {
			if !strings.Contains(redisNode, ":") {
				err = fmt.Errorf("invalid host [%s], needs to be specified as host:port", redisNode)
				return
			}
		}
	}

	values := urlObject.Query()
	enforceQuota = values.Get("quotaEnforce") == "true"
	jwtKey = values.Get("jwtKey")

	quotaBlacklistUpdateInterval, err = time.ParseDuration(values.Get("quotaBlacklistUpdateInterval"))
	if err != nil {
		return
	}

	whitelistedIpsString := values.Get("whitelist")
	if whitelistedIpsString == "" {
		// whitelist is optional
		whitelistedIps = map[string]bool{}
	} else {
		whitelistEntries := strings.Split(whitelistedIpsString, ",")
		whitelistedIps = make(map[string]bool)

		for _, entry := range whitelistEntries {
			// todo check if valid ip?
			whitelistedIps[entry] = true
		}
	}

	return
}

type authenticatorPlugin struct {
	// kmsVerificationKeyFunc jwt.Keyfunc
	enforceAuth  bool
	enforceQuota bool
}

func newAuthenticator(redisNodes []string, enforceQuota bool, jwtKey string, quotaBlacklistUpdateInterval time.Duration, whitelistedIps map[string]bool) *authenticatorPlugin {

	// todo add jwt Keyfunc if jwtKey is set
	enforceAuth := jwtKey != ""

	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: redisNodes,
	})

	dreddDB := dredd.NewDB(redisClient)
	Setup(dreddDB, quotaBlacklistUpdateInterval)

	return &authenticatorPlugin{
		// kmsVerificationKeyFunc: kmsVerificationKeyFunc,
		enforceAuth:  enforceAuth,
		enforceQuota: enforceQuota,
	}
}

func (a *authenticatorPlugin) IsAuthenticationTokenRequired() bool {
	return true
}

func (a *authenticatorPlugin) Check(ctx context.Context, token, ipAddress string) (context.Context, error) {
	credentials := &Credentials{}
	credentials.IP = ipAddress
	// todo only if jwt is disabled
	credentials.Subject = "uid:" + ipAddress

	zlog.Info("access token", zap.String("token", token))

	if a.enforceAuth {
		// todo implement
		/*
			parsedToken, err := jwt.ParseWithClaims(token, credentials, a.kmsVerificationKeyFunc)

			if err != nil {
				return ctx, err
			}
			expectedSigningAlgorithm := gcpjwt.SigningMethodKMSES256.Alg()
			actualSigningAlgorithm := parsedToken.Header["alg"]

			if expectedSigningAlgorithm != actualSigningAlgorithm {
				return ctx, fmt.Errorf("expected %s signing method but token specified %s", expectedSigningAlgorithm, actualSigningAlgorithm)
			}

			if !parsedToken.Valid {
				return ctx, errors.New("unable to verify token")
			}
		*/
	}

	authContext := authenticator.WithCredentials(ctx, credentials)

	// todo remove hard coded eosq token
	if a.enforceQuota || token == "csjBpe8I3UoJP6oqk5iYCCKF" {
		//zlog.Debug("adding cutoff to context", zap.String("user_id", credentials.Subject))
		withCutOffCtx, setCredentials := ContextWithCutOff(authContext)
		err := setCredentials(credentials)
		if err != nil {
			return withCutOffCtx, err
		}
		authContext = withCutOffCtx
	}
	return authContext, nil
}
