package c2cx

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/exchange-api/exchange"
)

const (
	binaryName = "c2cx"
)

var (
	binaryPath string
)

func TestMain(m *testing.M) {
	var err error
	binaryPath, err = filepath.Abs(binaryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("get binary name absolute path failed: %v\n", err))
	}
	// Build cli binary file.
	args := []string{"build", "-o", binaryPath, "./cli/c2cx.go"}
	fmt.Println(exec.Command("go", args...).Args)
	if err := exec.Command("go", args...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Make %v binary failed: %v\n", binaryName, err))
		os.Exit(1)
	}

	ret := m.Run()

	// Remove the generated cli binary file.
	if err := os.Remove(binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Delete %v failed: %v", binaryName, err))
		os.Exit(1)
	}

	os.Exit(ret)
}

// MarketOrderString is a one order in stock
type MarketOrderString struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

type orderbookJSONTestResponse struct {
	Timestamp *string `json:"Timestamp,omitempty"`
	Bids      []MarketOrderString
	Asks      []MarketOrderString
}

// Orderbook with timestamp
type orderbookTest struct {
	TradePair TradePair
	Timestamp time.Time
	Bids      exchange.MarketOrders
	Asks      exchange.MarketOrders
}

// UnmarshalJSON implements json.Unmarshaler
func (r *orderbookTest) UnmarshalJSON(b []byte) error {
	var v orderbookJSONTestResponse
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, *v.Timestamp)

	if err != nil {
		return err
	}
	r.Timestamp = t

	r.Bids = make(exchange.MarketOrders, len(v.Bids))
	for i := 0; i < len(v.Bids); i++ {
		price, err := decimal.NewFromString(v.Bids[i].Price)
		if err != nil {
			return err
		}
		volume, err := decimal.NewFromString(v.Bids[i].Volume)
		if err != nil {
			return err
		}
		r.Bids[i] = exchange.MarketOrder{
			Price:  price,
			Volume: volume,
		}
	}

	return nil
}

func TestGetOrderbook(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "get_orderbook - symbol not exists",
			args:       []string{"get_orderbook", "USDT_BT"},
			errMessage: "C2CX request failed: endpoint=getorderbook code=400 message=symbol not exists\n",
		},
		{
			name: "get_orderbook - OK",
			args: []string{"get_orderbook", "USDT_BTG"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err)
			var resp orderbookTest
			err = resp.UnmarshalJSON(output)
			require.NoError(t, err)
			for _, bid := range resp.Bids {
				require.True(t, bid.Volume.GreaterThan(decimal.Zero))
				require.True(t, bid.Price.GreaterThan(decimal.Zero))
			}

			for _, ask := range resp.Asks {
				require.True(t, ask.Volume.GreaterThan(decimal.Zero))
				require.True(t, ask.Price.GreaterThan(decimal.Zero))
			}

		})
	}
}

// BalanceSummaryString includes the account balance and its frozen balance
type BalanceSummaryString struct {
	Balance map[string]string `json:"balance"`
	Frozen  map[string]string `json:"frozen"`
}

func TestClient_GetBalanceSummary(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name: "get_balance_summary",
			args: []string{"get_balance_summary"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)
			fmt.Println(string(output))
			var resp BalanceSummaryString
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			require.True(t, len(resp.Balance) > 0)
			require.True(t, len(resp.Frozen) > 0)
			for key := range resp.Balance {
				balance, err := decimal.NewFromString(resp.Balance[key])
				require.NoError(t, err)
				frozen, err := decimal.NewFromString(resp.Frozen[key])
				require.NoError(t, err)
				require.True(t, balance.GreaterThanOrEqual(frozen))
			}
		})

	}
}

func randBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	x, err := rand.Read(b)
	assert.Equal(t, n, x)
	assert.Nil(t, err)
	return b
}

type createOrderResponseTest struct {
	OrderID uint64 `json:"orderID"`
}

// create_order
// get_orderinfo
// get_order_info_all
// cancel_order
// get_order_by_status
// get_order_by_status_paged
func TestCreateOrderCancelOrder(t *testing.T) {
	var orderID uint64
	cid := hex.EncodeToString(randBytes(t, 10))
	tt := []struct {
		name       string
		args       []string
		errMessage error
		message    string
	}{
		{
			name: "create_order - OK",
			args: []string{"create_order", "USDT_BTG", "14", "1", "buy", "limit", cid, "null", "null", "null"},
		},
		{
			name:       "create_order - cid already exists",
			args:       []string{"create_order", "USDT_BTG", "20", "0", "buy", "market", cid, "null", "null", "null"},
			message:    "C2CX request failed: endpoint=createorder code=400 message=cid already exists, please change\n",
			errMessage: errors.New("exit status 1"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.message, string(output))
				require.Equal(t, tc.errMessage.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			var resp createOrderResponseTest
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			orderID = resp.OrderID
		})

	}

	ttgoi := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name: "get_orderinfo - OK",
			args: []string{"get_orderinfo", "USDT_BTG", strconv.FormatUint(orderID, 10)},
		},
	}

	for _, tc := range ttgoi {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)
			var resp Order
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)

			require.Equal(t, OrderID(orderID), resp.OrderID, "can't find created order")
		})

	}

	ttgoa := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name: "get_order_info_all - OK",
			args: []string{"get_order_info_all", "USDT_BTG"},
		},
	}

	for _, tc := range ttgoa {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)
			var resp struct {
				Orders []Order `json:"orders"`
			}
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			var orderExists bool
			for _, order := range resp.Orders {
				if order.OrderID == OrderID(orderID) {
					orderExists = true
					break
				}
			}
			require.True(t, orderExists, "can't find created order")
		})
	}

	ttc := []struct {
		name       string
		args       []string
		message    string
		errMessage error
	}{
		{
			name: "cancel_order - OK",
			args: []string{"cancel_order", strconv.FormatUint(orderID, 10)},
		},
		{
			name:       "cancel_order - order in \"Cancelled\" status",
			args:       []string{"cancel_order", strconv.FormatUint(orderID, 10)},
			errMessage: errors.New("exit status 1"),
			message:    "C2CX request failed: endpoint=cancelorder code=400 message=you can't cancel order in \"Cancelled\" status\n",
		},
		{
			name:       "cancel_order - orderId is wrong",
			args:       []string{"cancel_order", strconv.FormatUint(orderID+1, 10)},
			errMessage: errors.New("exit status 1"),
			message:    "C2CX request failed: endpoint=cancelorder code=400 message=orderId is wrong\n",
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err)
			var resp struct {
				Result string
			}
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			require.Equal(t, "OK", resp.Result)
		})
	}

	ttg := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "get_order_by_status - symbol not exists",
			args:       []string{"get_order_by_status", "USDT_BT", "5"},
			errMessage: "C2CX request failed: endpoint=getorderbystatus code=400 message=symbol not exists\n",
		},
		{
			name: "get_order_by_status - OK",
			args: []string{"get_order_by_status", "USDT_BTG", "5"},
		},
	}

	for _, tc := range ttg {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp []Order
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			var orderExists bool
			for _, order := range resp {
				if order.OrderID == OrderID(orderID) {
					orderExists = true
					break
				}
			}
			require.True(t, orderExists, "can't find created order")
		})
	}

	ttgp := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "get_order_by_status_paged - symbol not exists",
			args:       []string{"get_order_by_status_paged", "USDT_BT", "5", "1"},
			errMessage: "C2CX request failed: endpoint=getorderbystatus code=400 message=symbol not exists\n",
		},
		{
			name: "get_order_by_status_paged - OK",
			args: []string{"get_order_by_status_paged", "USDT_BTG", "5", "1"},
		},
	}

	for _, tc := range ttgp {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp Orders
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			var orderExists bool
			for _, order := range resp.Orders {
				if order.OrderID == OrderID(orderID) {
					orderExists = true
					break
				}
			}
			require.True(t, orderExists, "can't find created order")
		})
	}
}

// create_order
// cancel_multiple
func TestCreateOrderCancelOrdersMultiple(t *testing.T) {
	var orderIDs = make([]string, 0)
	cid := hex.EncodeToString(randBytes(t, 10))
	tt := []struct {
		name       string
		args       []string
		errMessage error
		message    string
	}{
		{
			name: "create_order - OK",
			args: []string{"create_order", "USDT_BTG", "14", "1", "buy", "limit", cid, "null", "null", "null"},
		},
		{
			name: "create_order - OK",
			args: []string{"create_order", "USDT_BTG", "14", "1", "buy", "limit", cid + "_2", "null", "null", "null"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp createOrderResponseTest
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			orderIDs = append(orderIDs, strconv.FormatUint(resp.OrderID, 10))
		})

	}

	ttc := []struct {
		name       string
		args       []string
		message    string
		errMessage error
	}{
		{
			name: "cancel_multiple - OK",
			args: []string{"cancel_multiple"},
		},
		{
			name:       "cancel_multiple - FAIL",
			args:       []string{"cancel_multiple"},
			errMessage: errors.New("exit status 1"),
			message:    "these orders failed to cancel: [" + strings.Join(orderIDs, " ") + "]\n",
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, append(tc.args, orderIDs...)...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp []uint64
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			for idx, val := range resp {
				require.Equal(t, orderIDs[idx], strconv.FormatUint(val, 10))
			}
		})
	}
}

// create_order
// cancel_all
func TestCreateOrderCancelOrdersAll(t *testing.T) {
	var orderIDs = make([]string, 0)
	cid := hex.EncodeToString(randBytes(t, 10))
	tt := []struct {
		name       string
		args       []string
		errMessage error
		message    string
	}{
		{
			name: "create_order - OK",
			args: []string{"create_order", "USDT_BTG", "14", "1", "buy", "limit", cid, "null", "null", "null"},
		},
		{
			name: "create_order - OK",
			args: []string{"create_order", "USDT_BTG", "14", "1", "buy", "limit", cid + "_2", "null", "null", "null"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()

			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp createOrderResponseTest
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			orderIDs = append(orderIDs, strconv.FormatUint(resp.OrderID, 10))
		})

	}

	ttc := []struct {
		name       string
		args       []string
		message    string
		errMessage error
	}{
		{
			name: "cancel_all - OK",
			args: []string{"cancel_all", "USDT_BTG"},
		},
		{
			name:       "cancel_all - there is no items",
			args:       []string{"cancel_all", "USDT_BTG"},
			errMessage: errors.New("exit status 1"),
			message:    "C2CX request failed: endpoint=getorderinfo code=400 message=There is no items\n",
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp []uint64
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			for idx, val := range resp {
				require.Equal(t, orderIDs[idx], strconv.FormatUint(val, 10))
			}
		})
	}
}

// limit_buy
// cancel_order
func TestLimitBuyCancelOrder(t *testing.T) {
	var orderID uint64
	cid := hex.EncodeToString(randBytes(t, 10))
	tt := []struct {
		name       string
		args       []string
		errMessage error
		message    string
	}{
		{
			name: "limit_buy - OK",
			args: []string{"limit_buy", "USDT_BTG", "14", "1", cid},
		},
		{
			name:       "limit_buy - cid already exists",
			args:       []string{"limit_buy", "USDT_BTG", "14", "1", cid},
			errMessage: errors.New("exit status 1"),
			message:    "C2CX request failed: endpoint=createorder code=400 message=cid already exists, please change\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()

			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp createOrderResponseTest
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			orderID = resp.OrderID
		})

	}

	ttc := []struct {
		name       string
		args       []string
		message    string
		errMessage error
	}{
		{
			name: "cancel_order - OK",
			args: []string{"cancel_order", strconv.FormatUint(orderID, 10)},
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp struct {
				Result string
			}
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			require.Equal(t, "OK", resp.Result)
		})
	}
}

/////////////////////////////////////
//FAIL!!!
// Received message: C2CX request failed: endpoint=createorder code=400 message=Your Money must bigger than the min value 13 DRG
//// create_order, get_order_info_all, cancel_order
/////////////////////////////////////

//func TestMarketBuyCancelOrder(t *testing.T) {
//	var orderID uint64
//	cid := hex.EncodeToString(randBytes(t, 10))
//	tt := []struct {
//		name       string
//		args       []string
//		errMessage error
//		message    string
//	}{
//		{
//			name: "market_buy",
//			args: []string{"market_buy", "DRG_ETC", "1", cid},
//		},
//		{
//			name: "market_buy",
//			args: []string{"market_buy", "USDT_BTG", "1", cid},
//			errMessage: errors.New("exit status 1"),
//			message: "C2CX request failed: endpoint=createorder code=400 message=cid already exists, please change\n",
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			fmt.Println("name ==> " + tc.name)
//			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
//			fmt.Println("stdout: " + string(output))
//			if err != nil {
//				fmt.Println("stderr: " + err.Error())
//			}
//			if tc.errMessage != nil && err != nil {
//				require.Equal(t, tc.errMessage.Error(), err.Error())
//				require.Equal(t, tc.message, string(output))
//				return
//			}
//			require.NoError(t, err)
//			var resp createOrderResponseTest
//			err = json.Unmarshal(output, &resp)
//			require.NoError(t, err)
//			orderID = resp.OrderID
//		})
//	}
//
//	ttc := []struct {
//		name       string
//		args       []string
//		message    string
//		errMessage error
//	}{
//		{
//			name: "cancel_order",
//			args: []string{"cancel_order", strconv.FormatUint(orderID, 10)},
//		},
//	}
//
//	for _, tc := range ttc {
//		t.Run(tc.name, func(t *testing.T) {
//			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
//			if tc.errMessage != nil && err != nil {
//				require.Equal(t, tc.errMessage.Error(), err.Error())
//				require.Equal(t, tc.message, string(output))
//				return
//			}
//			require.NoError(t, err)
//			var resp struct {
//				Result string
//			}
//			err = json.Unmarshal(output, &resp)
//			require.NoError(t, err)
//			require.Equal(t, "OK", resp.Result)
//		})
//	}
//}

// limit_sell
// cancel_order
func TestLimitSellCancelOrder(t *testing.T) {
	var orderID uint64
	cid := hex.EncodeToString(randBytes(t, 10))
	tt := []struct {
		name       string
		args       []string
		errMessage error
		message    string
	}{
		{
			name: "limit_sell",
			args: []string{"limit_sell", "DRG_ETC", "14", "100", cid},
		},
		{
			name:       "limit_sell - cid already exists",
			args:       []string{"limit_sell", "DRG_ETC", "14", "100", cid},
			errMessage: errors.New("exit status 1"),
			message:    "C2CX request failed: endpoint=createorder code=400 message=cid already exists, please change\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp createOrderResponseTest
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			orderID = resp.OrderID
		})

	}

	ttc := []struct {
		name       string
		args       []string
		message    string
		errMessage error
	}{
		{
			name: "cancel_order - OK",
			args: []string{"cancel_order", strconv.FormatUint(orderID, 10)},
		},
	}

	for _, tc := range ttc {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err)
			var resp struct {
				Result string
			}
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)
			require.Equal(t, "OK", resp.Result)
		})
	}
}

type tickerDataTest struct {
	Timestamp time.Time        `json:"timestamp,omitempty"`
	High      *decimal.Decimal `json:"high,omitempty"`
	Last      *decimal.Decimal `json:"last,omitempty"`
	Low       *decimal.Decimal `json:"low,omitempty"`
	Buy       *decimal.Decimal `json:"buy,omitempty"`
	Sell      *decimal.Decimal `json:"sell,omitempty"`
	Volume    *decimal.Decimal `json:"volume,omitempty"`
}

func (td *tickerDataTest) UnmarshalJSON(b []byte) error {
	type tickerDataTestString struct {
		Timestamp *string `json:"timestamp"`
		High      *string `json:"high,omitempty"`
		Last      *string `json:"last,omitempty"`
		Low       *string `json:"low,omitempty"`
		Buy       *string `json:"buy,omitempty"`
		Sell      *string `json:"sell,omitempty"`
		Volume    *string `json:"volume,omitempty"`
	}
	var tmp tickerDataTestString
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, *tmp.Timestamp)

	if err != nil {
		return err
	}

	high, err := decimal.NewFromString(*tmp.High)
	if err != nil {
		return err
	}

	last, err := decimal.NewFromString(*tmp.Last)
	if err != nil {
		return err
	}

	low, err := decimal.NewFromString(*tmp.Low)
	if err != nil {
		return err
	}

	buy, err := decimal.NewFromString(*tmp.Buy)
	if err != nil {
		return err
	}

	sell, err := decimal.NewFromString(*tmp.Sell)
	if err != nil {
		return err
	}

	volume, err := decimal.NewFromString(*tmp.Volume)
	if err != nil {
		return err
	}
	td.Timestamp = t
	td.High = &high
	td.Low = &low
	td.Buy = &buy
	td.Sell = &sell
	td.Last = &last
	td.Volume = &volume
	return nil
}

func TestGetTicker(t *testing.T) {
	tc := []struct {
		name       string
		args       []string
		message    string
		errMessage error
	}{
		{
			name: "get_ticker - OK",
			args: []string{"get_ticker", "DRG_ETC"},
		},
		{
			name:       "get_ticker - unknown pair",
			args:       []string{"get_ticker", "DRG_UNKNOWN"},
			message:    "C2CX request failed: endpoint=ticker code=400 message=symbol does not exist\n",
			errMessage: errors.New("exit status 1"),
		},
	}
	now := time.Now()
	time.Sleep(5 * time.Second)
	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != nil && err != nil {
				require.Equal(t, tc.errMessage.Error(), err.Error())
				require.Equal(t, tc.message, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var resp tickerDataTest
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			require.True(t, resp.Timestamp.After(now))
			require.True(t, resp.High.GreaterThan(*resp.Low))
			require.True(t, resp.Sell.GreaterThan(*resp.Buy))
		})
	}
}
