package pbbilling

import (
	"github.com/golang/protobuf/ptypes"
)

func ToAvro(e *Event) map[string]interface{} {
	timestamp, err := ptypes.Timestamp(e.Timestamp)
	if err != nil {
		panic(err)
	}
	return map[string]interface{}{
		"timestamp":           timestamp,
		"user_id":             e.UserId,
		"api_key_id":          e.ApiKeyId,
		"source":              e.Source,
		"kind":                e.Kind,
		"usage":               e.Usage,
		"network":             e.Network,
		"requests_count":      e.RequestsCount,
		"responses_count":     e.ResponsesCount,
		"ratelimit_hit_count": e.RateLimitHitCount,
		"ingress_bytes":       e.IngressBytes,
		"egress_bytes":        e.EgressBytes,
		"idle_time_ms":        e.IdleTime,
		"ip_address":          e.IpAddress,
		"method":              e.Method,
	}
}
