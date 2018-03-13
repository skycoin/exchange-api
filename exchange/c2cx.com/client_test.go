// +build c2cx_integration_test

package c2cx

import (
	"errors"
	"testing"
	"time"

	exchange "github.com/skycoin/exchange-api/exchange"
)

func TestClientOperations(t *testing.T) {
	cl := Client{
		Key:                      key,
		Secret:                   secret,
		OrdersRefreshInterval:    time.Second * 5,
		Orders:     exchange.NewTracker(),
	}

	// verifying we've got enough SKY to play with
	availSky, err := availableSKY()
	if err != nil {
		t.Fatal(err)
	}

	if availSky.LessThan(orderAmount) {
		t.Fatal(errors.New("Test wallet doesn't have enough SKY"))
	}

	// creating an order
	orderId, err := cl.Sell(orderMarket, orderPrice, orderAmount)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("updateOrdersOrderbook", func(t *testing.T) {
		cl.updateOrders()
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
