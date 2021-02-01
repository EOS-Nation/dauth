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
	"sync"
	"sync/atomic"
	"time"

	"github.com/dfuse-io/dauth/authenticator"
	"go.uber.org/zap"
)

var limitValidator func(credentials authenticator.Credentials) (exceeded bool, reason string)
var dreddDB *dredd.DB
var blackListLock sync.Mutex
var blacklistedUserIDs = map[string]bool{}
var blackListVersion = -1

func Setup(ddb *dredd.DB, blacklistUpdateInterval time.Duration) {
	dreddDB = ddb
	if dreddDB == nil {
		panic("dreddDB can not be nil")
	}

	go func() {
		sleepDuration := time.Second * 0
		for {
			time.Sleep(sleepDuration)
			sleepDuration = blacklistUpdateInterval

			v, err := dreddDB.UserBlackListVersion()
			if err != nil {
				zlog.Error("failed to retrieve black listed version", zap.Error(err))
				continue
			}
			/*
				if v == blackListVersion {
					continue
				}
			*/
			blackListVersion = v
			zlog.Info("updating black, new version available", zap.Int("user_black_list", blackListVersion))

			ids, err := dreddDB.BlackListedUsers()
			blackListLock.Lock()
			blacklistedUserIDs = map[string]bool{}

			zlog.Info("blacklisted user ids", zap.Any("user_ids", blacklistedUserIDs))

			if err != nil {
				zlog.Error("failed to retrieve black listed user ids", zap.Error(err))
				blackListLock.Unlock()
				continue
			}
			zlog.Debug("updating black list:", zap.Int("id_count", len(ids)))
			for _, id := range ids {
				blacklistedUserIDs[id] = true
			}
			blackListLock.Unlock()
		}
	}()

	limitValidator = func(creds authenticator.Credentials) (exceeded bool, reason string) {
		blackListLock.Lock()
		defer blackListLock.Unlock()

		zlog.Info("checking limits", creds.GetLogFields()...)

		// userID := strings.TrimPrefix(creds.(*Credentials).Subject, "uid:")
		userID := creds.(*Credentials).IP
		if blackListed, _ := blacklistedUserIDs[userID]; blackListed {
			zlog.Debug("canceling context", zap.String("user_id", userID))
			return true, "document quota exceeded"
		}
		return false, ""
	}
}

// WithCutOff will cancel context when rate limit events occur for the
// concerned key. NOTE: Make sure you pass a Context that will be
// *canceled* in the future, to shutdown this goroutine.
func ContextWithCutOff(ctx context.Context) (*Context, SetCredentialsFunc) {
	ctx, cancel := context.WithCancel(ctx)
	wrapped := &Context{
		Context:    ctx,
		cancelFunc: cancel,
	}
	return wrapped, wrapped.SetCredentials
}

type Context struct {
	context.Context
	err             atomic.Value
	cancelFunc      context.CancelFunc
	credentials     authenticator.Credentials
	credentialsLock sync.Mutex
}

type SetCredentialsFunc func(authenticator.Credentials) error

func (ctx *Context) SetCredentials(credentials authenticator.Credentials) error {
	ctx.credentialsLock.Lock()
	defer ctx.credentialsLock.Unlock()
	if ctx.credentials != nil {
		return fmt.Errorf("calling SetCredentials twice on the same Context object")
	}
	ContextWithCutoffCounter.Inc()
	ctx.credentials = credentials
	if exceeded, reason := limitValidator(credentials); exceeded {
		err := fmt.Errorf("black listed: %s", reason)
		ctx.err.Store(err)
		ctx.cancelFunc()
		ContextWithCutoffCanceledCounter.Inc()
		return err
	}

	go ctx.listen()
	return nil
}

type ctxKeyType int

const ctxKey ctxKeyType = iota

func (ctx *Context) Value(key interface{}) interface{} {
	if key == ctxKey {
		return ctx.credentials
	}
	return ctx.Context.Value(key)
}

func (ctx *Context) listen() {
	for {
		if ctx.Context.Err() != nil {
			return
		}
		ctx.credentialsLock.Lock()
		credentials := ctx.credentials
		ctx.credentialsLock.Unlock()

		if exceeded, reason := limitValidator(credentials); exceeded {
			ctx.err.Store(fmt.Errorf("black listed: %s", reason))
			ctx.cancelFunc()
			ContextWithCutoffCanceledCounter.Inc()
			return
		}
		time.Sleep(5 * time.Second)
	}
}
func (ctx *Context) Err() error {
	if err := ctx.err.Load(); err != nil {
		return err.(error)
	}
	return ctx.Context.Err()
}
