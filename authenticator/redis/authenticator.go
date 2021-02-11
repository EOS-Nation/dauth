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
	"errors"
	"fmt"
	"github.com/dfuse-io/dauth/dredd"
	"github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dfuse-io/dauth/authenticator"
	"github.com/go-redis/redis/v8"
)

type authenticatorPlugin struct {
	ipQuotaHandler         *dredd.IpQuotaHandler
	kmsVerificationKeyFunc jwt.Keyfunc
	enforceQuota           bool
	enforceAuth            bool
}

func init() {
	// redis://redis1,redis2,redis3?quotaEnforce=true&jwtKey=abc123&quotaBlacklistUpdateInterval=5s&ipQuotaFile=/etc/quota.yml&defaultIpQuota=10
	authenticator.Register("redis", func(configURL string) (authenticator.Authenticator, error) {

		redisNodes, enforceQuota, jwtKey, quotaBlacklistUpdateInterval, ipQuotaHandler, err := parseURL(configURL)

		if err != nil {
			return nil, fmt.Errorf("redis auth factory: %w", err)
		}
		return newAuthenticator(redisNodes, enforceQuota, jwtKey, quotaBlacklistUpdateInterval, ipQuotaHandler), nil
	})
}

func parseURL(configURL string) (redisNodes []string, enforceQuota bool, jwtKey string, quotaBlacklistUpdateInterval time.Duration, ipQuotaHandler *dredd.IpQuotaHandler, err error) {
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

	// if we didn't get a jwt key here, try the env variables
	if jwtKey == "" {
		jwtKey = os.Getenv("JWT_SIGNING_KEY")
	}

	quotaBlacklistUpdateInterval, err = time.ParseDuration(values.Get("quotaBlacklistUpdateInterval"))
	if err != nil {
		return
	}

	// don't parse and error handle quota settings if we don't enforce it anyways
	if !enforceQuota {
		return
	}

	ipQuotaFile := values.Get("ipQuotaFile")
	defaultIpQuotaString := values.Get("defaultIpQuota")

	if defaultIpQuotaString != "" {
		var defaultIpQuota int
		defaultIpQuota, err = strconv.Atoi(defaultIpQuotaString)

		if err != nil {
			err = fmt.Errorf("failed to parse default ip quota, expected integer: %s", defaultIpQuotaString)
			return
		}

		if ipQuotaFile == "" {
			ipQuotaHandler = dredd.NewIpQuotaHandler(defaultIpQuota)
		} else {
			ipQuotaHandler, err = dredd.NewIpQuotaHandlerFromFile(ipQuotaFile, defaultIpQuota)

			if err != nil {
				err = fmt.Errorf("failed to parse ip quota file: %e", err)
				return
			}
		}
	} else {
		if ipQuotaFile != "" {
			err = fmt.Errorf("ip quota file given, but defaultIpQuota is not set")
			return
		}
		if jwtKey == "" {
			err = fmt.Errorf("enforceQuota is set but neither a jwt key or ip based quota handling is configured")
			return
		}
	}

	/*whitelistedIpsString := values.Get("whitelist")
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
	}*/

	return
}

func newAuthenticator(redisNodes []string, enforceQuota bool, jwtKey string, quotaBlacklistUpdateInterval time.Duration, ipQuotaHandler *dredd.IpQuotaHandler) *authenticatorPlugin {
	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: redisNodes,
	})

	dreddDB := dredd.NewDB(redisClient)
	Setup(dreddDB, quotaBlacklistUpdateInterval)

	return &authenticatorPlugin{
		kmsVerificationKeyFunc: func(token *jwt.Token) (interface{}, error) {
			if jwtKey == "" {
				return nil, fmt.Errorf("no jwt key set")
			}
			return []byte(jwtKey), nil
		},
		ipQuotaHandler: ipQuotaHandler,
		enforceQuota:   enforceQuota,
		// token is required if we don't have a ip quota handler but enforce doc quota
		enforceAuth: ipQuotaHandler == nil && enforceQuota,
	}
}

func (a *authenticatorPlugin) IsAuthenticationTokenRequired() bool {
	return a.enforceAuth
}

func (a *authenticatorPlugin) Check(ctx context.Context, token, ipAddress string) (context.Context, error) {
	credentials := &Credentials{}
	credentials.IP = ipAddress

	// if we have a token, try to get the credentials from it. A given token must always be valid
	// we exclude dfuse.io tokens from validation checks here which are in the format of (web|server|mobile)_abcdef
	if token != "" && !strings.Contains(token, "_") {
		parsedToken, err := jwt.ParseWithClaims(token, credentials, a.kmsVerificationKeyFunc)

		zlog.Info("access token", zap.String("token", token))
		zlog.Info("decoding issue", zap.Error(err))
		zlog.Info("parsed token", zap.Any("parsed_token", parsedToken))

		if err != nil {
			return ctx, err
		}
		if !parsedToken.Valid {
			return ctx, errors.New("unable to verify token")
		}

		zlog.Info("created token based credentials", zap.Any("credentials", credentials))
	} else {
		credentials.Subject = "uid:" + ipAddress

		// if we don't have a token, see if ip based quota handling is enabled and retrieve credentials from there
		if a.ipQuotaHandler != nil {
			quota, err := a.ipQuotaHandler.GetQuota(ipAddress)
			credentials.Quota = quota

			if err != nil {
				return ctx, err
			}

			zlog.Info("created ip quota based credentials", zap.Any("credentials", credentials))
		} else if a.enforceQuota {
			zlog.Info("didn't get a token but required one")
			return ctx, errors.New("no token given")
		}
	}

	authContext := authenticator.WithCredentials(ctx, credentials)

	if a.enforceQuota {
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
