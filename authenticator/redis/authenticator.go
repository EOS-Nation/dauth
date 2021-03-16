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
	ipQuotaHandler         *dredd.IpLimitHandler
	kmsVerificationKeyFunc jwt.Keyfunc
	enforceQuota           bool
	enforceAuth            bool
}

func init() {
	// redis://redis1,redis2,redis3?quotaEnforce=true&jwtKey=abc123&quotaBlacklistUpdateInterval=5s&ipQuotaFile=/etc/quota.yml&defaultIpQuota=10
	authenticator.Register("redis", func(configURL string) (authenticator.Authenticator, error) {

		redisNodes, db, enforceQuota, jwtKey, quotaBlacklistUpdateInterval, ipQuotaHandler, err := parseURL(configURL)

		if err != nil {
			return nil, fmt.Errorf("redis auth factory: %w", err)
		}
		return newAuthenticator(redisNodes, db, enforceQuota, jwtKey, quotaBlacklistUpdateInterval, ipQuotaHandler), nil
	})
}

func parseURL(configURL string) (redisNodes []string, db int, enforceQuota bool, jwtKey string, quotaBlacklistUpdateInterval time.Duration, ipQuotaHandler *dredd.IpLimitHandler, err error) {
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

	dbString := values.Get("redisDB")
	if dbString == "" {
		db = 0
	} else {
		db, err = strconv.Atoi(dbString)

		if err != nil {
			err = fmt.Errorf("failed to parse redisDB parameter, not an integer: %s", dbString)
			return
		}
	}

	quotaBlacklistUpdateInterval, err = time.ParseDuration(values.Get("quotaBlacklistUpdateInterval"))
	if err != nil {
		return
	}

	// don't parse and error handle quota settings if we don't enforce it anyways
	if !enforceQuota {
		return
	}

	ipLimitsFile := values.Get("ipLimitsFile")
	defaultIpQuotaString := values.Get("defaultIpQuota")
	defaultIpRateString := values.Get("defaultIpRate")

	if defaultIpQuotaString != "" || defaultIpRateString != "" {
		var defaultIpQuota int
		defaultIpQuota, err = strconv.Atoi(defaultIpQuotaString)

		if err != nil {
			err = fmt.Errorf("failed to parse default ip quota, expected integer: %s", defaultIpQuotaString)
			return
		}

		var defaultIpRate int
		defaultIpRate, err = strconv.Atoi(defaultIpRateString)

		if err != nil {
			err = fmt.Errorf("failed to parse default ip rate, expected integer: %s", defaultIpQuotaString)
			return
		}

		if ipLimitsFile == "" {
			ipQuotaHandler = dredd.NewIpLimitsHandler(defaultIpQuota, defaultIpRate)
		} else {
			ipQuotaHandler, err = dredd.NewIpLimitsHandlerFromFile(ipLimitsFile, defaultIpQuota, defaultIpRate)

			if err != nil {
				err = fmt.Errorf("failed to parse ip limits file: %e", err)
				return
			}
		}
	} else {
		if ipLimitsFile != "" {
			err = fmt.Errorf("ip limits file given, but defaultIpQuota or defaultIpRate is not set")
			return
		}
		if jwtKey == "" {
			err = fmt.Errorf("enforceQuota is set but neither a jwt key or ip based quota handling is configured")
			return
		}
	}

	return
}

func newAuthenticator(redisNodes []string, db int, enforceQuota bool, jwtKey string, quotaBlacklistUpdateInterval time.Duration, ipQuotaHandler *dredd.IpLimitHandler) *authenticatorPlugin {
	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: redisNodes,
		DB:            db,
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

	validToken := false

	// if we have a token, try to get the credentials from it. A given token must always be valid
	if token != "" {
		parsedToken, err := jwt.ParseWithClaims(token, credentials, a.kmsVerificationKeyFunc)

		if err != nil { // todo ||parsedToken == nil
			zlog.Warn("failed to decode token", zap.Error(err))
			// todo return ctx, err
		}
		if parsedToken != nil && !parsedToken.Valid {
			zlog.Warn("failed to verify token", zap.Any("token", parsedToken))
			// todo return ctx, errors.New("unable to verify token")
		} else {
			validToken = true
			zlog.Info("created token based credentials", zap.Any("credentials", credentials))
		}
	}

	// we don't have a valid token, try ip based quota limiting
	if !validToken {
		credentials.Subject = "ip:" + ipAddress

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
