package c2cx

import (
	"os"
	"path/filepath"
	"fmt"
	"testing"
	"os/exec"
	"github.com/stretchr/testify/require"
	"time"
	"github.com/skycoin/exchange-api/exchange"
	"encoding/json"
	"github.com/shopspring/decimal"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"encoding/hex"
	"strconv"
	"errors"
	"strings"
)

const (
	binaryName = "c2cx"
)

var (
	binaryPath string
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

type orderbookJSONToolResponse struct {
	Timestamp *string `json:"Timestamp,omitempty"`
	Bids      exchange.MarketOrdersString
	Asks      exchange.MarketOrdersString
}

// Orderbook with timestamp
type OrderbookTool struct {
	TradePair TradePair
	Timestamp time.Time
	Bids      exchange.MarketOrders
	Asks      exchange.MarketOrders
}

// UnmarshalJSON implements json.Unmarshaler
func (r *OrderbookTool) UnmarshalJSON(b []byte) error {
	var v orderbookJSONToolResponse
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
			var resp OrderbookTool
			err = resp.UnmarshalJSON(output)
			require.NoError(t, err)
			for _, bid := range resp.Bids {
				require.True(t, bid.Volume.Cmp(decimal.New(0, 10)) > 0)
				require.True(t, bid.Price.Cmp(decimal.New(0, 10)) > 0)
			}

			for _, ask := range resp.Asks {
				require.True(t, ask.Volume.Cmp(decimal.New(0, 10)) > 0)
				require.True(t, ask.Price.Cmp(decimal.New(0, 10)) > 0)
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
			args:       []string{"create_order", "USDT_BTG", "14", "0", "buy", "market", cid, "null", "null", "null"},
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

	tt_get_order_info := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name: "get_orderinfo - OK",
			args: []string{"get_orderinfo", "USDT_BTG", strconv.FormatUint(orderID, 10)},
		},
	}

	for _, tc := range tt_get_order_info {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			require.NoError(t, err)
			var resp Order
			err = json.Unmarshal(output, &resp)
			require.NoError(t, err)

			require.Equal(t, OrderID(orderID), resp.OrderID, "can't find created order")
		})

	}

	tt_get_orderinfo_all := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name: "get_order_info_all - OK",
			args: []string{"get_order_info_all", "USDT_BTG"},
		},
	}

	for _, tc := range tt_get_orderinfo_all {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			fmt.Println(tc.name)
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
			require.NoError(t, err)
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
			require.NoError(t, err)
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

func TestGetOrderByStatusPaged(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "get_order_by_status_paged - symbol not exists",
			args:       []string{"get_order_by_status_paged", "USDT_BT"},
			errMessage: "C2CX request failed: endpoint=getorderbook code=400 message=symbol not exists\n",
		},
		{
			name: "get_order_by_status_paged - OK",
			args: []string{"get_order_by_status_paged", "USDT_BTG", "cancelled", "1"},
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
			var resp OrderbookTool
			err = resp.UnmarshalJSON(output)
			require.NoError(t, err)
			for _, bid := range resp.Bids {
				require.True(t, bid.Volume.Cmp(decimal.New(0, 10)) > 0)
				require.True(t, bid.Price.Cmp(decimal.New(0, 10)) > 0)
			}

			for _, ask := range resp.Asks {
				require.True(t, ask.Volume.Cmp(decimal.New(0, 10)) > 0)
				require.True(t, ask.Price.Cmp(decimal.New(0, 10)) > 0)
			}

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
			require.NoError(t, err)
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
			require.NoError(t, err)
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

			require.NoError(t, err)
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
			require.NoError(t, err)
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
			require.NoError(t, err)
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
			args: []string{"limit_sell", "DRG_ETC", "14", "1", cid},
		},
		{
			name:       "limit_sell - cid already exists",
			args:       []string{"limit_sell", "DRG_ETC", "14", "1", cid},
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
			require.NoError(t, err)
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
