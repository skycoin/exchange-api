package cryptopia

import (
	"reflect"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	//fill cache
	_, _ = getCurrencyID("")
	_, _ = getMarketID("")
}

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
	crs, err := GetCurrencies()
	if err != nil {
		t.Fatal(err)
	}
	if len(crs) != 446 {
		t.Fatal("Incorrect count of currencies")
	}
}
func TestGetTradePairs(t *testing.T) {
	tps, err := GetTradePairs()
	if err != nil {
		t.Fatal(err)
	}
	if len(tps) < 1 {
		t.Fatal("Incorrect count of tradepairs")
	}
}
func TestGetMarkets(t *testing.T) {
	mkts, err := GetMarkets("ALL", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(mkts) < 1 {
		t.Fatal("empty")
	}
}
func TestGetMarket(t *testing.T) {
	mkt, err := GetMarket("LTC/BTC", -1)
	if err != nil {
		t.Fatal(err)
	}
	if mkt.TradePairID != 101 {
		t.Fatal("API error", "want 101 TradePairID")
	}
}
func TestGetMarketHistory(t *testing.T) {
	hst, err := GetMarketHistory("LTC/BTC", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(hst) < 1 {
		t.Fatal("empty")
	}
}
func TestGetMarketOrders(t *testing.T) {
	orders, err := GetMarketOrders("LTC/BTC", -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(orders.Buy) < 1 || len(orders.Sell) < 1 {
		t.Fatal("empty")
	}
}
func TestGetMarketOrderGroups(t *testing.T) {
	groups, err := GetMarketOrderGroups(-1, "LTC/BTC", "SKY/BTC")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatal("count of groups should be 2")
	}
	t.Logf("%#+v", groups)
}

// mock private API tests

func TestGetBalance(t *testing.T) {
	type result struct {
		Data string
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: "Total: 10300.00000000 Available: 6700.00000000 Unconfirmed: 2.00000000 Held: 3400.00000000 Pending: 200.00000000",
			Err:  nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":[{\"CurrencyId\":1,\"Symbol\":\"BTC\",\"Total\":10300,\"Available\":6700,\"Unconfirmed\":2,\"HeldForTrades\":3400,\"PendingWithdraw\":200,\"Address\":\"4HMjBARzTNdUpXCYkZDTHq8vmJQkdxXyFg\",\"BaseAddress\":\"ZDTHq8vmJQkdxXyFgZDTHq8vmJQkdxXyFgZDTHq8vmJQkdxXyFg\",\"Status\":\"OK\",\"StatusMessage\":\"\"}]}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"getbalance", httpmock.NewBytesResponder(200, v))
		data, err := GetBalance("key", "secret", "nonce", "BTC")
		if err != nil {
			t.Fatal(err)
		}
		httpmock.DeactivateAndReset()
		if data != k.Data {
			t.Fatalf("want %s\nexpected %s\n", k.Data, data)
		}
	}
}

func TestGetDepositAddress(t *testing.T) {
	type result struct {
		Data DepositAddress
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: DepositAddress{
				Address:     "4HMjBARzTNdUpXCYkZDTHq8vmJQkdxXyFg",
				BaseAddress: "ZDTHq8vmJQkdxXyFgZDTHq8vmJQkdxXyFgZDTHq8vmJQkdxXyFg",
				Currency:    "DOT",
			},
			Err: nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":{\"Currency\":\"DOT\",\"Address\":\"4HMjBARzTNdUpXCYkZDTHq8vmJQkdxXyFg\",\"BaseAddress\":\"ZDTHq8vmJQkdxXyFgZDTHq8vmJQkdxXyFgZDTHq8vmJQkdxXyFg\"}}"),
	}

	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"getdepositaddress", httpmock.NewBytesResponder(200, v))
		data, err := GetDepositAddress("key", "secret", "nonce", "DOT")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(k.Data, *data) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, *data)
		}
	}
}

func TestGetOpenOrders(t *testing.T) {
	type result struct {
		Data OpenedOrder
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: OpenedOrder{
				OrderID:     23467,
				TradePairID: 100,
				Market:      "DOT/BTC",
				Type:        "Buy",
				Rate:        0.00000034,
				Amount:      145.98,
				Total:       "0.00004963",
				Remaining:   "23.98760000",
				Timestamp:   "2014-12-07T20:04:05.3947572",
			},
			Err: nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":[{\"OrderId\":23467,\"TradePairId\":100,\"Market\":\"DOT/BTC\",\"Type\":\"Buy\",\"Rate\":3.4e-7,\"Amount\":145.98,\"Total\":\"0.00004963\",\"Remaining\":\"23.98760000\",\"TimeStamp\":\"2014-12-07T20:04:05.3947572\"}]}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"getopenorders", httpmock.NewBytesResponder(200, v))
		data, err := GetOpenOrders("key", "secret", "nonce", "DOT/BTC", 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) != 1 {
			t.Fatal("response entries count incorrect")
		}
		if !reflect.DeepEqual(k.Data, data[0]) {
			t.Fatalf("want %v\n expected%v\n", k.Data, data[0])
		}
		httpmock.DeactivateAndReset()
	}
}
func TestGetTradeHistory(t *testing.T) {
	type result struct {
		Data ClosedOrder
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: ClosedOrder{
				OrderID:     23467,
				TradePairID: 100,
				Market:      "DOT/BTC",
				Type:        "Buy",
				Rate:        0.00000034,
				Amount:      145.98,
				Total:       "0.00004963",
				Fee:         "0.98760000",
				Timestamp:   "2014-12-07T20:04:05.3947572",
			},
			Err: nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":[{\"TradeId\":23467,\"TradePairId\":100,\"Market\":\"DOT/BTC\",\"Type\":\"Buy\",\"Rate\":3.4e-7,\"Amount\":145.98,\"Total\":\"0.00004963\",\"Fee\":\"0.98760000\",\"TimeStamp\":\"2014-12-07T20:04:05.3947572\"}]}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"gettradehistory", httpmock.NewBytesResponder(200, v))
		data, err := GetTradeHistory("key", " secret", "nonce", "DOT/BTC", 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) != 1 {
			t.Fatal("response entries count incorrect")
		}
		if !reflect.DeepEqual(k.Data, data[0]) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, data[0])
		}
		httpmock.DeactivateAndReset()
	}
}
func TestGetTransactions(t *testing.T) {
	type result struct {
		Data Transaction
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: Transaction{
				ID:            23467,
				Currency:      "DOT",
				TxID:          "6ddbaca454c97ba4e8a87a1cb49fa5ceace80b89eaced84b46a8f52c2b8c8ca3",
				Type:          "Deposit",
				Amount:        145.98000000,
				Fee:           "0.00000000",
				Status:        "Confirmed",
				Confirmations: "20",
				Timestamp:     "2014-12-07T20:04:05.3947572",
				Address:       "",
			},
			Err: nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":[{\"Id\":23467,\"Currency\":\"DOT\",\"TxId\":\"6ddbaca454c97ba4e8a87a1cb49fa5ceace80b89eaced84b46a8f52c2b8c8ca3\",\"Type\":\"Deposit\",\"Amount\":145.98,\"Fee\":\"0.00000000\",\"Status\":\"Confirmed\",\"Confirmations\":\"20\",\"TimeStamp\":\"2014-12-07T20:04:05.3947572\",\"Address\":\"\"}]}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"gettransactions", httpmock.NewBytesResponder(200, v))
		data, err := GetTransactions("key", " secret", "nonce", TxTypeDeposit, 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) != 1 {
			t.Fatal("response entries incorrect count")
		}
		if !reflect.DeepEqual(k.Data, data[0]) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, data[0])
		}
		httpmock.DeactivateAndReset()
	}
}
func TestSubmitTrade(t *testing.T) {
	type result struct {
		Data NewTradeInfo
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: NewTradeInfo{
				OrderID:      23467,
				FilledOrders: []int{44310, 44311},
			},
			Err: nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":{\"OrderId\":23467,\"FilledOrders\":[44310,44311]}}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"submittrade", httpmock.NewBytesResponder(200, v))
		data, err := SubmitTrade("key", "secret", "nonce", "LTC/BTC", OfTypeBuy, 0.1, 0.1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(k.Data, *data) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, *data)
		}
		httpmock.DeactivateAndReset()
	}
}
func TestCancelTrade(t *testing.T) {
	type result struct {
		Data []int
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: []int{44310, 44311},
			Err:  nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":[44310,44311]}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"canceltrade", httpmock.NewBytesResponder(200, v))
		data, err := CancelTrade("key", "secret", "nonce", CancelAll, 0, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(k.Data) != len(data) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, data)
		}
		httpmock.DeactivateAndReset()
	}
}
func TestSubmitTip(t *testing.T) {
	type result struct {
		Data TipMessage
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: TipMessage("You tipped 45 users 0.00034500 DOT each."),
			Err:  nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":\"You tipped 45 users 0.00034500 DOT each.\"}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"submittip", httpmock.NewBytesResponder(200, v))
		data, err := SubmitTip("key", "", "nonce", "DOT", 45, 0.00034500)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(k.Data, data) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, data)
		}
		httpmock.DeactivateAndReset()
	}
}
func TestSubmitWithdraw(t *testing.T) {
	type result struct {
		Data WithdrawID
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: WithdrawID(405667),
			Err:  nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":405667}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"submitwithdraw", httpmock.NewBytesResponder(200, v))
		data, err := SubmitWithdraw("key", "secret", "nonce", "BTC", "", "", 1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(k.Data, data) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, data)
		}
		httpmock.DeactivateAndReset()
	}
}
func TestSubmitTrasnfer(t *testing.T) {
	type result struct {
		Data TransferMessage
		Err  error
	}
	var tests = map[*result][]byte{
		&result{
			Data: "Successfully transfered 200 DOT to Hex.",
			Err:  nil,
		}: []byte("{\"Success\":true,\"Error\":null,\"Data\":\"Successfully transfered 200 DOT to Hex.\"}"),
	}
	for k, v := range tests {
		httpmock.Activate()
		httpmock.RegisterResponder("POST", apiroot.String()+"submittransfer", httpmock.NewBytesResponder(200, v))
		data, err := SubmitTransfer("key", "secret", "nonce", "DOT", "Hex", 200)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(k.Data, data) {
			t.Fatalf("want %v\nexpected %v\n", k.Data, data)
		}
		httpmock.DeactivateAndReset()
	}
}
