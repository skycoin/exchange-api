package c2cx

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/skycoin/exchange-api/db"

	exchange "github.com/skycoin/exchange-api/exchange"
)

var (
	cl = Client{
		Key:                      key,
		Secret:                   secret,
		OrdersRefreshInterval:    time.Second * 5,
		OrderbookRefreshInterval: time.Second * 5,
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx"),
	}
	orderMarket = "CNY_SHL"
	orderPrice  = 0.01
	orderAmount = 10.0
	orderID     int
)

func TestClientCreateOrder(t *testing.T) {
	var err error
	orderID, err = cl.Buy(orderMarket, orderPrice, orderAmount)
	if err != nil {
		t.Error(err)
	}
}
func TestClientUpdateOrders(t *testing.T) {
	cl.updateOrders()
	cl.updateOrderbook()
}
func TestClientGetExecuted(t *testing.T) {
	orders := cl.Executed()
	if len(orders) != 1 {
		t.Error("placed order not found in tracker", len(orders), orders)
	}
	if orders[0] != orderID {
		t.Errorf("want %d orderID, expected %d", orderID, orders[0])
	}
}

func TestClientGetCompleted(t *testing.T) {
	orders := cl.Completed()
	if len(orders) > 0 {
		t.Error("it should not contains completed orders")
	}
}

func TestClientCancelMarket(t *testing.T) {
	_, err := cl.Cancel(orderID)
	if err != nil {
		t.Error(err)
	}
	if len(cl.Orders.GetCompleted()) != 1 {
		t.Error("it should have one completed order")
	}
}
