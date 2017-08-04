package c2cx

import (
	"testing"
)

var (
	key    = "2A4C851A-1B86-4E08-863B-14582094CE0F"
	secret = "83262169-B473-4BF2-9288-5E5AC898F4B0"
)
var (
	order  int
	market = "CNY_SHL"
)

func TestCreateOrder(t *testing.T) {
	orderid, err := CreateOrder(key, secret, "CNY/SHL", 0.01, 10, "Buy", PriceTypeLimit, nil)
	if err != nil {
		t.Fatal(err)
	}
	//t.Logf("Order %d successfully created", orderid)
	order = orderid
}

func TestGetOrderInfo(t *testing.T) {
	orders, err := GetOrderInfo(key, secret, market, -1, nil)
	if err != nil {
		t.Fatal(err)
	}
	var found = false
	for _, v := range orders {
		if order == v.OrderID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("order %d not found", order)
	}

}
func TestCancelOrder(t *testing.T) {
	err := CancelOrder(key, secret, order)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserInfo(t *testing.T) {
	b, err := GetBalance(key, secret)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 5 {
		t.Fatal("invalid balance response")
	}
}
