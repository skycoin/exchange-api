// +build redis_integration_test
// +build cryptopia_integration_test

package cryptopia_test

import (
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"

	"github.com/skycoin/exchange-api/db"
	"github.com/skycoin/exchange-api/exchange"
	cryptopia "github.com/skycoin/exchange-api/exchange/cryptopia.co.nz"
)

var redisAddr = func() string {
	res, found := os.LookupEnv("REDIS_TEST_ADDR")
	if !found {
		panic("redis test address not provided")
	}
	return res
}()

var key, secret = func() (key string, secret string) {
	var found bool
	if key, found = os.LookupEnv("CRYPTOPIA_TEST_KEY"); found {
		if secret, found = os.LookupEnv("CRYPTOPIA_TEST_SECRET"); found {
			return
		}
		panic("Cryptopia secret not provided")
	}
	panic("Cryptopia key not provided")
}()

func TestClientInit(t *testing.T) {
	var c exchange.Client

	var client = cryptopia.Client{
		Key:                      key,
		Secret:                   secret,
		OrdersRefreshInterval:    time.Second * 5,
		OrderbookRefreshInterval: time.Second * 1,
		Stop:         make(chan struct{}),
		TrackedBooks: []string{"LTC/BTC", "SKY/DOGE"},
		Orders:       exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: redisAddr,
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
