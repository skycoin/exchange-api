package cryptopia

import (
	"testing"

	"github.com/uberfurrer/tradebot/exchange"
)

var (
	key    = "23a69c51c746446e819b213ef3841920"
	secret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po="
)

func TestRequestSignature(t *testing.T) {
	var (
		key    = "23a69c51c746446e819b213ef3841920"
		secret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po="
		nonce  = "3"
		requrl = apiroot
	)
	requrl.Path += "getbalance"
	var want = "amx 23a69c51c746446e819b213ef3841920:VTUkpXJ8Cl2VfoRXH6qaPK887Ejy58UC2mPEwB80w2M=:3"
	if expected := header(key, secret, nonce, requrl, []byte("{}")); want != expected {
		t.Fatal("invalid request signature")
	}
}

// Works
func TestSubmitTrade(t *testing.T) {
	var (
		key    = "23a69c51c746446e819b213ef3841920"
		secret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po="
	)
	order, err := submitTrade(key, secret, nonce(), "LTC/BTC", exchange.Buy, 0.001, 0.01)
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
