// +build c2cx_integration_test
// +build redis_integration_test

package c2cx

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"

	"github.com/skycoin/exchange-api/db"

	exchange "github.com/skycoin/exchange-api/exchange"
)

var redisAddr = func() string {
	res, found := os.LookupEnv("REDIS_TEST_ADDR")
	if !found {
		panic("redis test address not provided")
	}
	return res
}()

func TestClientOperations(t *testing.T) {
	cl := Client{
		Key:                      key,
		Secret:                   secret,
		OrdersRefreshInterval:    time.Second * 5,
		OrderbookRefreshInterval: time.Second * 5,
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: redisAddr,
		}, "c2cx"),
	}

	// verifying we've got enough SKY to play with
	availSky, err := availableSKY()
	if err != nil {
		t.Fatal(err)
	}
	if availSky < orderAmount {
		t.Fatal(errors.New("Test wallet doesn't have enough SKY"))
	}

	// creating an order
	orderId, err := cl.Sell(orderMarket, orderPrice, orderAmount)
	if err != nil {
		t.Error(err)
	}

	t.Run("updateOrdersOrderbook", func(t *testing.T) {
		cl.updateOrders()
		cl.updateOrderbook()
	})

	t.Run("GetExecuted", func(t *testing.T) {
		orderIds := cl.Executed()
		for _, v := range orderIds {
			if v == orderId {
				return
			}
		}
		t.Errorf("couldn't locate order #%d", orderId)
	})

	t.Run("GetCompletedFirstPass", func(t *testing.T) {
		orderIds := cl.Completed()
		for _, v := range orderIds {
			if v == orderId {
				t.Errorf("order #%d shouldn't have completed", orderId)
			}
		}
		return
	})

	// finally cleanup our order
	_, err = cl.Cancel(orderId)
	if err != nil {
		t.Error(err)
	}

	// and confirm it shows up in completed
	t.Run("GetCompletedSecondPass", func(t *testing.T) {
		orderIds := cl.Completed()
		for _, v := range orderIds {
			if v == orderId {
				return
			}
		}
		t.Errorf("couldn't locate order #%d", orderId)
	})
}
