package c2cx

import (
	"net/url"
	"testing"

	"github.com/uberfurrer/tradebot/exchange"
)

func TestGetOrderByStatus(t *testing.T) {
	var key, secret = "2A4C851A-1B86-4E08-863B-14582094CE0F", "83262169-B473-4BF2-9288-5E5AC898F4B0"
	orders, err := getOrderByStatus(key, secret, "ETH_SKY", exchange.Cancelled, -1)
	if err != nil && len(orders) > 0 {
		t.Error(err)
	}
	t.Log(orders)
}

func TestGetEthSkyOrderbook(t *testing.T) {
	orderbook, err := getOrderbook("ETH_SKY")
	if err != nil {
		t.Error(err)
	}
	if orderbook.Timestamp == 0 {
		t.Error("corrupted record", orderbook)
	}
}

func TestCreateOrder(t *testing.T) {
	var (
		key    = "2A4C851A-1B86-4E08-863B-14582094CE0F"
		secret = "83262169-B473-4BF2-9288-5E5AC898F4B0"
	)
	orderid, err := CreateOrder(key, secret, "CNY_SHL", 0.10, 10, "Buy", PriceTypeLimit, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orderid, err)
}

func TestCreateOrder2(t *testing.T) {
	var (
		key    = "2A4C851A-1B86-4E08-863B-14582094CE0F"
		secret = "83262169-B473-4BF2-9288-5E5AC898F4B0"
		params = url.Values{
			"symbol":         []string{"BTC_SKY"},
			"priceTypeId":    []string{"1"},
			"isAdvanceOrder": []string{"0"},
			"orderType":      []string{"Buy"},
			"price":          []string{"0.001"},
			"quantity":       []string{"5"},
		}
	)
	resp, err := requestPost("createOrder", key, secret, params)
	t.Logf("%#v, %#v", resp, err)
}
