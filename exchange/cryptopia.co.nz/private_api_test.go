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

func TestRequestSignature(t *testing.T) {
	var (
		key    = "abababababababababababababababab"
		secret = "YWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWJhYmFiYWI="
		nonce  = "3"
		requrl = apiroot
	)
	requrl.Path += "getbalance"
	var want = "amx abababababababababababababababab:QRB4yf+QkSxxzPg6JLDeNFdAsTu24wpiDozHNQZ3Jkc=:3"
	if expected := header(key, secret, nonce, requrl, []byte("{}")); want != expected {
		t.Fatal("invalid request signature")
	}
}

// Works
func TestSubmitTrade(t *testing.T) {
	order, err := submitTrade(key, secret, nonce(), "SKY/BTC", exchange.Buy, 0.0005, 10)
	t.Log(order, err)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cancelTrade(key, secret, nonce(), ByOrderID, nil, &order)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTradeHistory(t *testing.T) {
	var (
		key    = "23a69c51c746446e819b213ef3841920"
		secret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po="
	)
	if orders, err := getTradeHistory(key, secret, nonce(), nil, nil); err != nil || len(orders) != 1 {
		t.Fatalf("invalid order history info %v %v", orders, err)
	}
}

func TestGetOpenOrders(t *testing.T) {
	var (
		key    = "23a69c51c746446e819b213ef3841920"
		secret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po="
	)
	if orders, err := getOpenOrders(key, secret, nonce(), nil, nil); err != nil || len(orders) != 0 {
		t.Fatalf("invalid open orders info %v %v", orders, err)
	}
}

func TestWithdrawHistory(t *testing.T) {
	var (
		key    = "23a69c51c746446e819b213ef3841920"
		secret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po="
	)
	if txs, err := getTransactions(key, secret, nonce(), TxTypeWithdraw, 1); err != nil || len(txs) != 1 {
		t.Fatalf("invalid transactions %v %v", txs, err)
	}
}
