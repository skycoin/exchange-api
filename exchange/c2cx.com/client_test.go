package c2cx

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/uberfurrer/tradebot/db"
	"github.com/uberfurrer/tradebot/exchange"
)

func TestClientUpdateOrderbook(t *testing.T) {
	var c = Client{
		Key:    "",
		Secret: "",
		Orders: nil,
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx"),
	}
	c.updateOrderbook()
	for _, v := range markets {
		book, err := c.Orderbook().Get(v)
		if err != nil {
			t.Error(err)
		}
		if sym, err := normalize(book.Symbol); err != nil || sym != v {
			t.Error("corrupted orderbook record")
		}
	}

}

func TestClientUpdate(t *testing.T) {
	var c = &Client{
		Key:                      "",
		Secret:                   "",
		OrderRefreshInterval:     time.Second * 5,
		OrderbookRefreshInterval: time.Second * 5,
		Stop:   make(chan struct{}),
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx_test"),
	}
	go c.Update()
	defer func() { c.Stop <- struct{}{} }()
	for {
		_, err := c.Orderbook().Get("BTC/SKY")
		if err != nil {
			continue
		}
		break
	}
}

func TestClientGetBalance(t *testing.T) {
	var c = &Client{
		Secret:                   "83262169-B473-4BF2-9288-5E5AC898F4B0",
		Key:                      "2A4C851A-1B86-4E08-863B-14582094CE0F",
		OrderRefreshInterval:     time.Second * 5,
		OrderbookRefreshInterval: time.Second * 5,
		Stop:   make(chan struct{}),
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx_test"),
	}
	if b, err := c.GetBalance("BTC"); err != nil || len(b) == 0 {
		t.Error("GetBalance failed")
	}
}

func TestClientCreateOrder(t *testing.T) {
	var c = &Client{
		Secret:                   "83262169-B473-4BF2-9288-5E5AC898F4B0",
		Key:                      "2A4C851A-1B86-4E08-863B-14582094CE0F",
		OrderRefreshInterval:     time.Second * 5,
		OrderbookRefreshInterval: time.Second * 5,
		Stop:   make(chan struct{}),
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx_test"),
	}
	go c.Update()
	defer func() { c.Stop <- struct{}{} }()
	orderid, err := c.Buy("BTC/SKY", 0.00001, 0.00001)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("created", orderid)

	t.Log("tracker", c.Executed())

	order, err := c.Cancel(orderid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("cancelled", order)
	t.Log("tracker", c.Completed())
}
