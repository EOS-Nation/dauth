package keyer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentConsumptionLast30Days(t *testing.T) {
	dd, err := time.Parse("2006-01-02", "2020-01-15")
	require.NoError(t, err)

	keys := DocumentConsumptionLast30Days("uid:user.id", dd)
	require.Len(t, keys, 30)
	assert.Equal(t, "DCD:uid:user.id:20200115", keys[0])
	for i := 0; i < 30; i++ {
		h := -24 * time.Hour * time.Duration(i)
		fmt.Println("h:", h)
		d := dd.Add(h)
		expectedKey := "DCD:uid:user.id:" + d.Format("20060102")
		assert.Equal(t, expectedKey, keys[i])
		fmt.Println(keys[i])
	}

}
