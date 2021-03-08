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
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/dfuse-io/dauth/authenticator"
	"github.com/dfuse-io/dauth/dredd"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithRateLimit(t *testing.T) {
	ctx := context.Background()
	ctx, f := ContextWithCutOff(ctx)
	Setup(dredd.NewDB(newTestClient()), 3*time.Second)

	assert.NotNil(t, f)
	assert.NotNil(t, ctx)
}

func TestContext_SetCredentials(t *testing.T) {
	ctx := context.Background()
	ctx, setCredentials := ContextWithCutOff(ctx)
	Setup(dredd.NewDB(newTestClient()), 3*time.Second)
	done := make(chan bool)
	var once sync.Once
	limitValidator = func(credentials authenticator.Credentials) (exceeded bool, reason string) {
		once.Do(func() {
			close(done)
		})

		assert.Equal(t, "api.key.1", credentials.(*Credentials).APIKeyID)
		return false, ""
	}

	credentials := &Credentials{
		APIKeyID: "api.key.1",
	}
	setCredentials(credentials)
	<-done
}

func TestContext_ExceededLimit(t *testing.T) {
	ctx := context.Background()
	ctx, setCredentials := ContextWithCutOff(ctx)

	Setup(dredd.NewDB(newTestClient()), 3*time.Second)
	blacklistedUserIDs["user.id.1"] = true
	credentials := &Credentials{
		StandardClaims: jwt.StandardClaims{
			Subject: "uid:user.id.1",
		},
		APIKeyID: "api.key.1",
	}
	setCredentials(credentials)
	<-ctx.Done()
	assert.Error(t, ctx.Err())
	assert.Equal(t, "blocked: document quota exceeded", ctx.Err().Error())
}

func TestContext_ListenOnCloseContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx, setCredentials := ContextWithCutOff(ctx)
	Setup(dredd.NewDB(newTestClient()), 3*time.Second)
	credentials := &Credentials{
		APIKeyID: "api.key.1",
	}
	cancel()
	setCredentials(credentials)
	<-ctx.Done()
	assert.Error(t, ctx.Err())
	assert.Equal(t, "context canceled", ctx.Err().Error())
}

func TestContext_Err(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx, _ = ContextWithCutOff(ctx)

	require.NoError(t, ctx.Err())
	cancel()
	assert.Error(t, ctx.Err())
	assert.Equal(t, "context canceled", ctx.Err().Error())
}

func TestContext_Value(t *testing.T) {
	ctx := context.Background()
	ctx, setCredentials := ContextWithCutOff(ctx)

	expectedCredentials := &Credentials{
		APIKeyID: "api.key.1",
	}
	setCredentials(expectedCredentials)

	require.Nil(t, ctx.Value("key.1"))
	require.Equal(t, expectedCredentials, ctx.Value(ctxKey))
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
