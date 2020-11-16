package pbbilling

import (
	"time"
)

func TestAvroMap() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":           time.Unix(0, 0),
		"ip_addr":             "127.0.0.1",
		"user_id":             "uid:mdfuse2f4c5791b9a7315",
		"api_key_id":          "69548590991bccfabc68fed790cbfe57a2bed388cab28702c48f417f3c040ab7",
		"source":              "grapheos",
		"kind":                "GET /v0/search/transactions",
		"usage":               "usage.1",
		"network":             "network.1",
		"requests_count":      1,
		"responses_count":     2,
		"ratelimit_hit_count": 3,
		"ingress_bytes":       4,
		"egress_bytes":        5,
		"idle_time_ms":        int64(10),
	}
}

func TestAvroData() []byte {
	data, err := Codec.BinaryFromNative(nil, TestAvroMap())
	if err != nil {
		panic(err)
	}
	return data
}
