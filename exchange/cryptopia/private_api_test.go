// +build cryptopia_integration_test

package cryptopia

import (
	"os"
	"testing"

	"github.com/skycoin/exchange-api/exchange"
)

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

// Works
func TestSubmitTrade(t *testing.T) {
	c := Client{
		Key:    key,
		Secret: secret,
	}

	order, err := c.SubmitTrade(key, secret, nonce(), "SKY/BTC", exchange.Buy, 0.0005, 10)
	t.Log(order, err)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.CancelTrade(key, secret, nonce(), ByOrderID, nil, &order)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTradeHistory(t *testing.T) {
	c := Client{
		Key:    key,
		Secret: secret,
	}

	if orders, err := c.GetTradeHistory(nil, nil); err != nil || len(orders) != 1 {
		t.Fatalf("invalid order history info %v %v", orders, err)
	}
}

func TestGetOpenOrders(t *testing.T) {
	c := Client{
		Key:    key,
		Secret: secret,
	}

	if orders, err := c.GetOpenOrders(nil, nil); err != nil || len(orders) != 0 {
		t.Fatalf("invalid open orders info %v %v", orders, err)
	}
}

func TestWithdrawHistory(t *testing.T) {
	c := Client{
		Key:    key,
		Secret: secret,
	}

	if txs, err := c.GetTransactions(TxTypeWithdraw, 1); err != nil || len(txs) != 1 {
		t.Fatalf("invalid transactions %v %v", txs, err)
	}
}
