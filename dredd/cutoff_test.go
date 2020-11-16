package dredd

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
)

func TestLuaScript(t *testing.T) {
	t.Skip()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	handler, err := NewLuaEventHandler(client)
	require.NoError(t, err)
	start := time.Now()
	n := time.Now()
	users := 1
	days := 31
	for c := 0; c < users; c++ {
		for d := 0; d < days; d++ {
			date := n.Add((-24 * time.Hour) * time.Duration(d))
			//fmt.Println("Day: ", d, " date: ", date)
			_, err = handler.HandleEvent(fmt.Sprintf("user.id.%d", c+1), "api.key.1", "0.0.0.0", 10, date, 200, 0)
			require.NoError(t, err)
		}
	}

	fmt.Println("result in :", time.Since(start))

}
