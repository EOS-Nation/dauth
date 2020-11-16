package redis

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/dfuse-io/dauth/dredd"
	pbbilling "github.com/dfuse-io/dauth/pb/dfuse/billing/v1"
	"github.com/dfuse-io/derr"
	"github.com/go-redis/redis/v8"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dfuse-io/dauth/authenticator"
	redis_auth "github.com/dfuse-io/dauth/authenticator/redis"
	"github.com/dfuse-io/dmetering"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

func init() {
	dmetering.Register("redis", func(config string) (dmetering.Metering, error) {
		u, err := url.Parse(config)
		if err != nil {
			return nil, err
		}

		vals := u.Query()
		networkID := vals.Get("networkId")
		if networkID == "" {
			return nil, fmt.Errorf("missing networkId query param to metering config")
		}

		ipQuotaFile := vals.Get("ipQuotaFile")
		defaultQuotaString := vals.Get("defaultQuota")
		var ipQuotaHandler *dredd.IpQuotaHandler

		if defaultQuotaString != "" {
			defaultQuota, err := strconv.Atoi(defaultQuotaString)

			if err != nil {
				return nil, fmt.Errorf("failed to parse default quota, expected integer: %s", defaultQuotaString)
			}

			if ipQuotaFile == "" {
				ipQuotaHandler = dredd.NewIpQuotaHandler(defaultQuota)
			} else {
				ipQuotaHandler, err = dredd.NewIpQuotaHandlerFromFile(ipQuotaFile, defaultQuota)

				if err != nil {
					return nil, fmt.Errorf("failed to parse ip quota file: %e", err)
				}
			}
		} else {
			zlog.Warn("no default quota set, there won't be any rate limiting")
			ipQuotaHandler = dredd.NewIpQuotaHandler(0)
		}

		var emitterDelay = 10 * time.Second
		emitterDelayString := vals.Get("emitterDelay")
		if emitterDelayString != "" {
			if d, err := time.ParseDuration(emitterDelayString); err == nil {
				emitterDelay = d
			}
		}

		project := u.Host
		if project == "" {
			return nil, fmt.Errorf("project not specified (as hostname)")
		}

		hosts := strings.Split(u.Host, ",")

		topic := strings.TrimLeft(u.Path, "/")
		if topic == "" {
			return nil, fmt.Errorf("topic not specified (as path component)")
		}

		warnOnErrors := vals.Get("warnOnErrors") == "true"

		return newMetering(networkID, hosts, topic, warnOnErrors, emitterDelay, nil, ipQuotaHandler), nil
	})
}

type meteringPlugin struct {
	network string

	redisClient        *redis.Client
	warnOnPubSubErrors bool
	pubSubTopic        string

	quotaHandler *dredd.IpQuotaHandler

	messagesCount atomic.Uint64
	errorCount    atomic.Uint64

	luaHandler   *dredd.LuaEventHandler
	accumulator  *Accumulator
}

// type topicProviderFunc func(pubsubProject string, topicName string) *pubsub.Topic
type topicEmitterFunc func(e *pbbilling.Event)

func newMetering(network string, hosts []string, pubSubTopic string, warnOnPubSubErrors bool, emitterDelay time.Duration, /*topicProvider topicProviderFunc,*/ topicEmitter topicEmitterFunc, quotaHandler *dredd.IpQuotaHandler) *meteringPlugin {
	m := &meteringPlugin{
		network:            network,
		warnOnPubSubErrors: warnOnPubSubErrors,
		quotaHandler:       quotaHandler,
	}

	m.redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: hosts,
	})

	m.pubSubTopic = pubSubTopic

	luaHandler, err := dredd.NewLuaEventHandler(m.redisClient)

	if err != nil {
		derr.Check("failed to init dredd handler", err)
	}

	m.luaHandler = luaHandler

	/*	if topicProvider == nil {
			m.topic = defaultTopicProvider(pubSubProject, pubSubTopic)
		} else {
			m.topic = topicProvider(pubSubProject, pubSubTopic)
		}*/

	if topicEmitter == nil {
		m.accumulator = newAccumulator(m.defaultTopicEmitter, emitterDelay)
	} else {
		m.accumulator = newAccumulator(topicEmitter, emitterDelay)
	}

	zlog.Info("metering is ready to emit")
	return m
}

func (m *meteringPlugin) EmitWithContext(ev dmetering.Event, ctx context.Context) {
	credentials := authenticator.GetCredentials(ctx)
	m.EmitWithCredentials(ev, credentials)
}

func (m *meteringPlugin) EmitWithCredentials(ev dmetering.Event, creds authenticator.Credentials) {
	userEvent := &pbbilling.Event{
		Source:            ev.Source,
		Kind:              ev.Kind,
		Network:           m.network,
		RequestsCount:     ev.RequestsCount,
		ResponsesCount:    ev.ResponsesCount,
		RateLimitHitCount: ev.RateLimitHitCount,
		IngressBytes:      ev.IngressBytes,
		EgressBytes:       ev.EgressBytes,
		IdleTime:          ev.IdleTime,
	}

	switch c := creds.(type) {
	case *authenticator.AnonymousCredentials:
		userEvent.UserId = "anonymous"
		userEvent.ApiKeyId = "anonymous"
		userEvent.Usage = "anonymous"
		userEvent.IpAddress = "0.0.0.0"
	case *redis_auth.Credentials:
		// userEvent.UserId = c.Subject
		userEvent.UserId = c.IP
		userEvent.ApiKeyId = c.APIKeyID
		userEvent.Usage = c.Usage
		userEvent.IpAddress = c.IP
	}

	docQuota, err := m.quotaHandler.GetQuota(userEvent.UserId)

	if err != nil {
		zlog.Warn("failed to get doc quota", zap.String("user_id", userEvent.UserId), zap.Error(err))
	}

	// todo add doc quota
	_, err = m.luaHandler.HandleEvent(userEvent, docQuota)

	if err != nil {
		zlog.Warn("failed to execute lua script", zap.Error(err))
	}

	m.emit(userEvent)
}

func (m *meteringPlugin) emit(e *pbbilling.Event) {
	m.messagesCount.Inc()
	if e.Timestamp == nil {
		e.Timestamp = ptypes.TimestampNow()
	}
	m.accumulator.emit(e)
}

func (m *meteringPlugin) GetStatusCounters() (total, errors uint64) {
	return m.messagesCount.Load(), m.errorCount.Load()
}

func (m *meteringPlugin) WaitToFlush() {
	zlog.Info("gracefully shutting down, now flushing pending dbilling events")
	m.accumulator.emitAccumulatedEvents()
	// m.topic.Stop()
	zlog.Info("all billing events have been flushed before shutdown")
}

/*func defaultTopicProvider(pubsubProject string, topicName string) *pubsub.Topic {
	ctx := context.Background()

	client := redis.NewClusterClient(&redis.ClusterOptions{

	})

	topics, err := client.PubSubChannels(ctx, topicName).Result()

	if err != nil || len(topics) == 0 {
		zlog.Panic("unable to setup dbilling PubSub connection", zap.String("project", pubsubProject), zap.String("topic", topicName), zap.Error(err))
	}

	client.Publish(ctx)

	topic := client.Topic(topicName)
	topic.PublishSettings = pubsub.PublishSettings{
		ByteThreshold:  20000,
		CountThreshold: 100,
		DelayThreshold: 1 * time.Second,
	}

	exists, err := topic.Exists(ctx)
	if err != nil || !exists {
		zlog.Panic("unable to setup dbilling PubSub connection", zap.String("project", pubsubProject), zap.String("topic", topicName), zap.Error(err))
	}
	return topic
}*/

func (m *meteringPlugin) defaultTopicEmitter(e *pbbilling.Event) {
	if e.UserId == "" || e.Source == "" || e.Kind == "" {
		zlog.Warn("events SHALL minimally contain UserID, Source and Kind, dropping billing event", zap.Any("event", e))
		return
	}

	cmd := &pbbilling.Command{
		Action: &pbbilling.Command_EventAction{
			EventAction: &pbbilling.EventAction{
				Event: e,
			},
		},
	}

	data, err := proto.Marshal(cmd)
	if err != nil {
		m.errorCount.Inc()
		return
	}

	zlog.Debug("sending message", zap.String("data_hex", hex.EncodeToString(data)))

	newCmd := &pbbilling.Command{}
	err = proto.Unmarshal(data, newCmd)
	if err != nil {
		panic(err)
	}
	zlog.Debug("decoded command", zap.Reflect("cmd", newCmd))

	res := m.redisClient.Publish(context.Background(), m.pubSubTopic, data)

	/*res := m.topic.Publish(context.Background(), &pubsub.Message{
		Data: data,
	})*/

	if m.warnOnPubSubErrors {
		if err := res.Err(); err != nil {
			zlog.Warn("failed to publish", zap.Error(err))
		}
	}
}
