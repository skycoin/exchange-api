package c2cx

import (
	"testing"

	exchange "github.com/uberfurrer/tradebot/exchange"
)

func TestGetOrderByStatus(t *testing.T) {
	var key, secret = "2A4C851A-1B86-4E08-863B-14582094CE0F", "83262169-B473-4BF2-9288-5E5AC898F4B0"
	orders, err := GetOrderByStatus(key, secret, "cny/btc", exchange.StatusCompleted, -1)
	if err != nil {
		t.Error(err)
	}
	t.Log(orders)
}
