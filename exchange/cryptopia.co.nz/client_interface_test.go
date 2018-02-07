package cryptopia_test

import (
	"testing"
	"time"

	"github.com/skycoin/exchange-api/db"
	"github.com/skycoin/exchange-api/exchange"
	cryptopia "github.com/skycoin/exchange-api/exchange/cryptopia.co.nz"
)

func TestClientInit(t *testing.T) {
	var c exchange.Client

	orderBook, err := db.NewOrderbookTracker(db.MemoryDatabase,
		"",
		"cryptopia")

	if err != nil {
		t.Fatal(err)
	}

	var client = cryptopia.Client{
		Key:                      "23a69c51c746446e819b213ef3841920",
		Secret:                   "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po=",
		OrdersRefreshInterval:    time.Second * 5,
		OrderbookRefreshInterval: time.Second * 1,
		Stop:         make(chan struct{}),
		TrackedBooks: []string{"LTC/BTC", "SKY/DOGE"},
		Orders:       exchange.NewTracker(),
		Orderbooks:   orderBook,
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
