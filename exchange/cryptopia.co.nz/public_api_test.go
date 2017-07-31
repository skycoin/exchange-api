package cryptopia

import (
	"testing"
)

func Test_getCurrencyID(t *testing.T) {
	btcID, err := getCurrencyID("btc")
	if err != nil {
		t.Fatal(err)
	}
	if btcID != 1 {
		t.Errorf("Incorrect BTC id %d, want %d", btcID, 1)
	}
	ltcID, err := getCurrencyID("ltc")
	if ltcID != 3 {
		t.Errorf("Incorrect BTC id %d, want %d", ltcID, 3)
	}
	if err != nil {
		t.Fatal(err)
	}
	skyID, err := getCurrencyID("sky")
	if err != nil {
		t.Fatal(err)
	}
	if skyID != 504 {
		t.Errorf("Incorrect BTC id %d, want %d", skyID, 504)
	}
}

func Test_getMarketID(t *testing.T) {
	btcltc, err := getMarketID("ltc_btc")
	if err != nil || btcltc != 101 {
		t.Fatal(err, btcltc)
	}
}

func TestGetCurrencies(t *testing.T) {
	_, err := getCurrencies()
	if err != nil {
		t.Fatal(err)
	}

}
func TestGetTradePairs(t *testing.T) {
	_, err := getTradePairs()
	if err != nil {
		t.Fatal(err)
	}
}
func TestGetMarkets(t *testing.T) {
	mkts, err := getMarkets("ALL", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(mkts) < 1 {
		t.Fatal("empty")
	}
}
func TestGetMarket(t *testing.T) {
	mkt, err := getMarket("LTC/BTC", -1)
	if err != nil {
		t.Fatal(err)
	}
	if mkt.TradePairID != 101 {
		t.Fatal("API error", "want 101 TradePairID")
	}
}
func TestGetMarketHistory(t *testing.T) {
	hst, err := getMarketHistory("LTC/BTC", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(hst) < 1 {
		t.Fatal("empty")
	}
}
func TestGetMarketOrders(t *testing.T) {
	orders, err := getMarketOrders("LTC/BTC", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(orders.Buy) < 1 || len(orders.Sell) < 1 {
		t.Fatal("empty")
	}
}
func TestGetMarketOrderGroups(t *testing.T) {
	groups, err := getMarketOrderGroups(-1, "LTC/BTC", "SKY/BTC")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatal("count of groups should be 2")
	}
}
