package redis

import "fmt"

// RCC:[PREFIX]:[UID]:[TIMESTAMP_MINUTE]
func requestConsumptionCounter(prefix, uid, timestampMinute string) string {
	return fmt.Sprintf("RCC:%s:%s:%s", prefix, uid, timestampMinute)
}
