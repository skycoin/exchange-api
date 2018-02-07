package c2cx

import (
	"testing"
	"time"

	"github.com/skycoin/exchange-api/db"

	exchange "github.com/skycoin/exchange-api/exchange"
)

var (
	orderMarket = "CNY_SHL"
	orderPrice  = 0.01
	orderAmount = 10.0
	orderID     int
)

func newClient() (*Client, error) {
	orderBookDatabase, err := db.NewOrderbookTracker("memory", "c2cx", "")

	if err != nil {
		return nil, err
	}

	cl := &Client{
		Key:                      key,
		Secret:                   secret,
		OrdersRefreshInterval:    time.Second * 5,
		OrderbookRefreshInterval: time.Second * 5,
		Orders:     exchange.NewTracker(),
		Orderbooks: orderBookDatabase,
	}

	return cl, nil
}

func TestClientCreateOrder(t *testing.T) {
	var err error
	cl, err := newClient()

	if err != nil {
		t.Fatal(err)
	}

	orderID, err = cl.Buy(orderMarket, orderPrice, orderAmount)
	if err != nil {
		t.Error(err)
	}
}

func TestClientUpdateOrders(t *testing.T) {
	cl, err := newClient()

	if err != nil {
		t.Fatal(err)
	}

	cl.updateOrders()
	cl.updateOrderbook()
}

func TestClientGetExecuted(t *testing.T) {
	cl, err := newClient()

	if err != nil {
		t.Fatal(err)
	}

	orders := cl.Executed()
	if len(orders) != 1 {
		t.Fatal("placed order not found in tracker", len(orders), orders)
	}
	if orders[0] != orderID {
		t.Fatalf("want %d orderID, expected %d", orderID, orders[0])
	}
}

func TestClientGetCompleted(t *testing.T) {
	cl, err := newClient()

	if err != nil {
		t.Fatal(err)
	}

	orders := cl.Completed()
	if len(orders) > 0 {
		t.Error("it should not contains completed orders")
	}
}

func TestClientCancelMarket(t *testing.T) {
	cl, err := newClient()

	if err != nil {
		t.Fatal(err)
	}

	_, err = cl.Cancel(orderID)
	if err != nil {
		t.Error(err)
	}
	if len(cl.Orders.GetCompleted()) != 1 {
		t.Error("it should have one completed order")
	}
}
