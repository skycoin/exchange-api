// +build cryptopia_integration_test

package cryptopia

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
	binaryName = "cryptopia"
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
	args := []string{"build", "-o", binaryPath, "./cli/cryptopia.go"}
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

// MarketOrders is a orderbook for market
type MarketOrdersTest struct {
	Buy  []MarketOrderTest `json:"Buy"`
	Sell []MarketOrderTest `json:"Sell"`
}

// MarketOrder represents a single order info
type MarketOrderTest struct {
	TradePairID int             `json:"TradePairId"`
	Label       string          `json:"Label"`
	Price       decimal.Decimal `json:"Price"`
	Volume      decimal.Decimal `json:"Volume"`
	Total       decimal.Decimal `json:"Total"`
}

func (m *MarketOrderTest) UnmarshalJSON(b []byte) error {
	var tmp struct {
		TradePairID int    `json:"TradePairId"`
		Label       string `json:"Label"`
		Price       string `json:"Price"`
		Volume      string `json:"Volume"`
		Total       string `json:"Total"`
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	price, err := decimal.NewFromString(tmp.Price)
	if err != nil {
		return err
	}

	volume, err := decimal.NewFromString(tmp.Volume)
	if err != nil {
		return err
	}

	total, err := decimal.NewFromString(tmp.Total)
	if err != nil {
		return err
	}
	m.TradePairID = tmp.TradePairID
	m.Label = tmp.Label
	m.Price = price
	m.Volume = volume
	m.Total = total
	return nil
}

func TestGetMarketOrdersIntegration(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		label      string
		errMessage string
	}{
		{
			name:       "get_market_orders - trade pair not found",
			label:      "LTC/BTC",
			errMessage: "Trade pair not found\n",
			args:       []string{"get_market_orders", "LTC/DUMMY", "-1"},
		},
		{
			name:  "get_market_orders - OK",
			label: "LTC/BTC",
			args:  []string{"get_market_orders", "LTC/BTC", "-1"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res MarketOrdersTest
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))
			for _, order := range res.Buy {
				require.Equal(t, tc.label, order.Label)
				require.True(t, order.Price.GreaterThan(decimal.New(0, 0)))
				require.True(t, order.Volume.GreaterThan(decimal.New(0, 0)))
				require.True(t, order.Total.GreaterThan(decimal.New(0, 0)))
			}

			for _, order := range res.Sell {
				require.Equal(t, tc.label, order.Label)
				require.True(t, order.Price.GreaterThan(decimal.New(0, 0)))
				require.True(t, order.Volume.GreaterThan(decimal.New(0, 0)))
				require.True(t, order.Total.GreaterThan(decimal.New(0, 0)))
			}
		})
	}
}

type MarketTest struct {
	TradePairId int               `json:"tradePairId"`
	Market      string            `json:"market"`
	Buy         []MarketOrderTest `json:"Buy"`
	Sell        []MarketOrderTest `json:"Sell"`
}

func TestGetMarketOrderGroupsIntegration(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		labels     map[string]bool
		errMessage string
	}{
		{
			name: "get_market_order_groups - trade pair not found",
			labels: map[string]bool{
				"LTC_BTC": true,
				"SKY_BTC": true,
			},
			errMessage: "Trade pair not found\n",
			args:       []string{"get_market_order_groups", "1", "LTC/DUMMY", "SKY/BTC"},
		},
		{
			name: "get_market_order_groups - OK",
			labels: map[string]bool{
				"LTC_BTC": true,
				"SKY_BTC": true,
			},
			args: []string{"get_market_order_groups", "1", "LTC/BTC", "SKY/BTC"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res []MarketTest
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))

			for _, market := range res {
				require.True(t, tc.labels[market.Market])
				for _, order := range market.Buy {
					require.True(t, order.Price.GreaterThan(decimal.New(0, 0)))
					require.True(t, order.Volume.GreaterThan(decimal.New(0, 0)))
					require.True(t, order.Total.GreaterThan(decimal.New(0, 0)))
				}

				for _, order := range market.Sell {
					require.True(t, order.Price.GreaterThan(decimal.New(0, 0)))
					require.True(t, order.Volume.GreaterThan(decimal.New(0, 0)))
					require.True(t, order.Total.GreaterThan(decimal.New(0, 0)))
				}
			}

		})
	}
}

func TestGetBalanceIntegration(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "get_balance - OK",
			args:       []string{"get_balance", "BTC"},
			errMessage: "GetBalance failed: No balance found, Currency BTC Rawdata null\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res []MarketTest
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))

		})
	}
}

func TestGetDepositAddressIntegration(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name: "get_deposit_address - OK",
			args: []string{"get_deposit_address", "DOT"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if err != nil {
				fmt.Printf("stderr: %v\n", err.Error())
			}

			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res DepositAddress
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))
			require.NotEmpty(t, res.Address)
		})
	}
}

func TestGetOpenOrdersIntegration(t *testing.T) {
	tt := []struct {
		name        string
		args        []string
		ordersCount int
		errMessage  string
	}{
		{
			name:        "get_open_orders - OK",
			args:        []string{"get_open_orders", "SKY/BTC", "10"},
			ordersCount: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if err != nil {
				fmt.Printf("stderr: %v\n", err.Error())
			}

			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res []Order
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))
			require.Equal(t, len(res), tc.ordersCount)
		})
	}
}

func TestGetTradeHistoryIntegration(t *testing.T) {
	tt := []struct {
		name        string
		args        []string
		ordersCount int
		errMessage  string
	}{
		{
			name:        "get_trade_history - OK",
			args:        []string{"get_trade_history", "SKY/BTC", "10"},
			ordersCount: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if err != nil {
				fmt.Printf("stderr: %v\n", err.Error())
			}

			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res []Order
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))
			fmt.Println(string(output))
			require.Equal(t, len(res), tc.ordersCount)
		})
	}
}

func TestSubmitTradeIntegration(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "submit_trade - insufficient Funds",
			args:       []string{"submit_trade", "SKY/BTC", OfferTypeBuy, "0.00050000", "123.00000000"},
			errMessage: "SubmitTrade failed: Insufficient Funds., Type Buy Market SKY/BTC Rate 0.0005 Amount 123\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if err != nil {
				fmt.Printf("stderr: %v\n", err.Error())
			}

			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res []Order
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))
		})
	}
}

func TestCancelTradeIntegration(t *testing.T) {
	tt := []struct {
		name       string
		args       []string
		errMessage string
	}{
		{
			name:       "cancel_trade - no trades found to cancel for tradepair",
			args:       []string{"cancel_trade", "TradePair", "SKY/BTC", "null"},
			errMessage: "No trades found to cancel for tradepair.\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output, err := exec.Command(binaryPath, tc.args...).CombinedOutput()
			if err != nil {
				fmt.Printf("stderr: %v\n", err.Error())
			}

			if tc.errMessage != "" && err != nil {
				require.Equal(t, tc.errMessage, string(output))
				return
			}
			require.NoError(t, err, fmt.Sprintf("stdout: %v", string(output)))
			var res []Order
			err = json.Unmarshal(output, &res)
			require.NoError(t, err, "stdout: %v", string(output))
		})
	}
}
