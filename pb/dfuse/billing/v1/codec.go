package pbbilling

import (
	"github.com/linkedin/goavro/v2"
)

var Codec *goavro.Codec

func init() {
	var err error
	Codec, err = goavro.NewCodec(`{
		"namespace": "io.dfuse",
		"type": "record",
		"name": "BillableEvent",
		"fields": [
			{ "name": "timestamp", "type" : {"type": "long", "logicalType" : "timestamp-millis"} },
			{ "name": "user_id", "type": "string" },
			{ "name": "api_key_id", "type": "string" },
			{ "name": "source", "type": "string" },
			{ "name": "kind", "type": "string" },
			{ "name": "usage", "type": "string" },
			{ "name": "network", "type": "string" },
			{ "name": "requests_count", "type": "long" },
			{ "name": "responses_count", "type": "long" },
			{ "name": "ratelimit_hit_count", "type": "long" },
			{ "name": "ingress_bytes", "type": "long" },
			{ "name": "egress_bytes", "type": "long" },
			{ "name": "idle_time_ms", "type": "long" },
			{ "name": "ip_address", "type": "string", "default": "" },
			{ "name": "method", "type": "string", "default": "" }
		]
	}`)
	if err != nil {
		panic("Unable to parse AVRO schema")
	}
}
