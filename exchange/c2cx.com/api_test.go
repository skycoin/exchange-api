// +build c2cx_integration_test

package c2cx

import (
	"os"
	"testing"
)

var key, secret = func() (key string, secret string) {
	var found bool
	if key, found = os.LookupEnv("C2CX_TEST_KEY"); found {
		if secret, found = os.LookupEnv("C2CX_TEST_SECRET"); found {
			return
		}
		panic("C2CX secret not provided")
	}
	panic("C2CX key not provided")
}()

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
