package cryptopia_test

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/uberfurrer/tradebot/db"
	"github.com/uberfurrer/tradebot/exchange"
	cryptopia "github.com/uberfurrer/tradebot/exchange/cryptopia.co.nz"
)

func TestClientInit(t *testing.T) {
	var c exchange.Client

	var client = cryptopia.Client{
		Key:                      "23a69c51c746446e819b213ef3841920",
		Secret:                   "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po=",
		OrdersRefreshInterval:    time.Second * 5,
		OrderbookRefreshInterval: time.Second * 1,
		Stop:         make(chan struct{}),
		TrackedBooks: []string{"LTC/BTC", "SKY/DOGE"},
		Orders:       exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "cryptopia_test"),
	}
	go client.Update()
	c = &client
	balance, err := c.GetBalance("BTC")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
	for {
		orderbook, err := c.Orderbook().Get("ltc_btc")
		if err != nil {
			continue
		}
		if orderbook.Symbol != "LTC/BTC" {
			t.Fatal("invalid symbol in db")
		}
		break
	}
	client.Stop <- struct{}{}
}
