package c2cx

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestUnixMilli(t *testing.T) {
	var x int64 = 1521563904999
	y := fromUnixMilli(x)
	z := toUnixMilli(y)
	require.Equal(t, x, z)

	y = fromUnixMilli(0)
	require.True(t, y.IsZero())

	z = toUnixMilli(y)
	require.Equal(t, int64(0), z)
}

func TestOrderJSON(t *testing.T) {
	var x int64 = 1521563904999
	y := fromUnixMilli(x)

	o := Order{
		Amount:          decimal.New(123, -3),
		AvgPrice:        decimal.New(345, -1),
		CompletedAmount: decimal.New(321, -2),
		Fee:             decimal.New(1, 0),
		CreateDate:      y,
		CompleteDate:    y,
		OrderID:         1234,
		Price:           decimal.New(456, -4),
		Status:          StatusActive,
		Type:            OrderTypeBuy,
		Trigger:         nil,
		CustomerID:      nil,
		Source:          "api",
	}

	p, err := json.Marshal(o)
	require.NoError(t, err)

	t.Log(string(p))

	var q Order
	err = json.Unmarshal(p, &q)
	require.NoError(t, err)
	require.Equal(t, o, q)

	// trigger and customerID non-nil pointer handling
	trigger := decimal.New(789, 0)
	o.Trigger = &trigger
	customerID := "foo-cid"
	o.CustomerID = &customerID

	p, err = json.Marshal(o)
	require.NoError(t, err)

	q = Order{}
	err = json.Unmarshal(p, &q)
	require.NoError(t, err)
	require.Equal(t, o, q)

	// 0 timestamp conversion
	o.CreateDate = time.Time{}
	o.CompleteDate = time.Time{}
	require.True(t, o.CreateDate.IsZero())
	require.True(t, o.CompleteDate.IsZero())
	p, err = json.Marshal(o)
	require.NoError(t, err)

	q = Order{}
	err = json.Unmarshal(p, &q)
	require.NoError(t, err)
	require.Equal(t, o, q)
	require.True(t, o.CreateDate.IsZero())
	require.True(t, o.CompleteDate.IsZero())
}
