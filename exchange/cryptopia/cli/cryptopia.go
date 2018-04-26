package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/skycoin/exchange-api/exchange/cryptopia"
)

var (
	rootCmd *cobra.Command
	client  *cryptopia.Client
)

const null = "null"

func getCommands() map[string]*cobra.Command {
	return map[string]*cobra.Command{
		"get_currencies": {
			Use:   "get_currencies",
			Short: "get_currencies returns all currency data",
			Long: `
get_currencies returns all currency data.
	Params:
		-
`,
			Example: "cryptopia get_currencies",
			Run: func(cmd *cobra.Command, args []string) {
				currencies, err := client.GetCurrencies()
				handleResult(currencies, err)
			},
		},
		"get_trade_pairs": {
			Use:   "get_trade_pairs",
			Short: "get_trade_pairs returns all trade pair data on exchange",
			Long: `
get_trade_pairs returns all trade pair data on exchange.
	Params:
		-
`,
			Example: "cryptopia get_trade_pairs",
			Run: func(cmd *cobra.Command, args []string) {
				tradePairs, err := client.GetTradePairs()
				handleResult(tradePairs, err)
			},
		},
		"get_markets": {
			Use:   "get_markets",
			Short: "get_markets returns all Market info by given baseMarket",
			Long: `
get_market returns all Market info by given baseMarket
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		hours - period in hours
`,
			Example: "cryptopia get_markets <markets> <hours>",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				markets := args[0]
				hours, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					printErrorWithExit(err)
					return
				}
				marketsResult, err := client.GetMarkets(markets, int(hours))
				handleResult(marketsResult, err)
			},
		},
		"get_market": {
			Use:   "get_market",
			Short: "get_market returns market with given label",
			Long: `
get_market returns market with given label.
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		hours - period in hours
`,
			Example: "cryptopia get_market <market> <hours>",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				market := args[0]
				hours, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					printErrorWithExit(err)
					return
				}
				marketsResult, err := client.GetMarket(market, int(hours))
				handleResult(marketsResult, err)
			},
		},
		"get_market_history": {
			Use:   "get_market_history",
			Short: "get_market_history returns market history with given label",
			Long: `
get_market_history returns market history with given label.
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		hours - period in hours
`,
			Example: "cryptopia get_market_history <market> <hours>",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				market := args[0]
				hours, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					printErrorWithExit(err)
					return
				}
				marketsResult, err := client.GetMarketHistory(market, int(hours))
				handleResult(marketsResult, err)
			},
		},
		"get_market_orders": {
			Use:   "get_market_orders",
			Short: "get_market_orders returns count orders from market with given label",
			Long: `
get_market_orders returns count orders from market with given label.
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		count - orders count, if count < 1, its will be omitted, default value is 100
`,
			Example:            "cryptopia get_market_orders <market> <count>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				market := args[0]
				count, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				marketsResult, err := client.GetMarketOrders(market, int(count))
				handleResult(marketsResult, err)
			},
		},
		"get_market_order_groups": {
			Use:   "get_market_order_groups",
			Short: "get_market_order_groups returns count orders from market with given label",
			Long: `
get_market_order_groups returns count orders from market with given label.
	Params:
		count - orders count, if count < 1, it will be omitted
		market - the market symbol of the trade e.g. 'SKY/BTC'
`,
			Example:            "cryptopia get_market_order_groups <count> <market> <market> <market>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				count, err := strconv.Atoi(args[0])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				var markets = make([]string, 0)
				markets = append(markets, args[1:]...)
				marketsResult, err := client.GetMarketOrderGroups(count, markets)
				handleResult(marketsResult, err)
			},
		},
		"get_balance": {
			Use:   "get_balance",
			Short: "get_balance returns a string representation of balance by given currency",
			Long: `
get_balance returns a string representation of balance by given currency.
	Params:
		currency - The currency symbol of the coins e.g. 'SKY'
`,
			Example:            "cryptopia get_balance <currency>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				currency := args[0]

				balance, err := client.GetBalance(currency)
				handleResult(balance, err)
			},
		},
		"get_deposit_address": {
			Use:   "get_deposit_address",
			Short: "get_deposit_address returns a deposit address of given currency",
			Long: `
get_deposit_address returns a deposit address of given currency.
	Params:
		currency - the currency symbol of the coins e.g. 'SKY'
`,
			Example:            "cryptopia get_deposit_address <currency>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				currency := args[0]
				depositAddress, err := client.GetDepositAddress(currency)
				handleResult(depositAddress, err)
			},
		},
		"get_open_orders": {
			Use:   "get_open_orders",
			Short: "get_open_orders returns a list of opened orders by specific market or all markets",
			Long: `
get_open_orders returns a list of opened orders by specific market or all markets.
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		count - orders count
`,
			Example:            "cryptopia get_open_orders <market> <count>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				var market *string
				var count *int
				var err error
				market = &args[0]
				if len(args) == 2 {
					tmp, err := strconv.Atoi(args[1])
					if err != nil {
						printErrorWithExit(err)
						return
					}
					count = &tmp
				}
				openOrders, err := client.GetOpenOrders(market, count)
				handleResult(openOrders, err)
			},
		},
		"get_trade_history": {
			Use:   "get_trade_history",
			Short: "get_trade_history return a list of all executed orders by specific market or all markets",
			Long: `
get_trade_history return a list of all executed orders by specific market or all markets.
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		count - orders count
`,
			Example:            "cryptopia get_trade_history <market> <count>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				var market *string
				var count *int
				market = &args[0]
				if len(args) == 2 {
					tmp, err := strconv.Atoi(args[1])
					if err != nil {
						printErrorWithExit(err)
						return
					}
					count = &tmp
				}
				tradeHistory, err := client.GetTradeHistory(market, count)
				handleResult(tradeHistory, err)
			},
		},
		"get_transactions": {
			Use:   "get_transactions",
			Short: "get_transactions a list of transactions",
			Long: `
get_transactions a list of transactions.
	Params:
		txType - The type of transactions to return e.g. 'Deposit' or 'Withdraw'
		count - The maximum amount of transactions to return e.g. '10'. If count < 1, it will be omitted
`,
			Example:            "cryptopia get_trade_history <txType> <count>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				txType := args[0]
				count, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				transactions, err := client.GetTransactions(txType, count)
				handleResult(transactions, err)
			},
		},

		"submit_trade": {
			Use:   "submit_trade",
			Short: "submit_trade submits a new trade offer",
			Long: `
submit_trade submits a new trade offer.
	Params:
		market - the market symbol of the trade e.g. 'SKY/BTC'
		offer_type - the type of trade e.g. 'buy' or 'sell'
		rate - the rate or price to pay for the coins e.g. 0.00000034
		amount - the amount of coins to buy e.g. 123.00000000
`,
			Example:            "cryptopia submit_trade <market> <offer_type> <rate> <amount>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(4),
			Run: func(cmd *cobra.Command, args []string) {
				market := args[0]
				offerType := args[1]
				rate, err := decimal.NewFromString(args[2])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				amount, err := decimal.NewFromString(args[3])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				orderID, err := client.SubmitTrade(market, offerType, rate, amount)
				handleResult(struct {
					OrderID int
				}{OrderID: orderID}, err)
			},
		},
		"cancel_trade": {
			Use:   "cancel_trade",
			Short: "cancel_trade cancel trades by given orderID, market or add active depends of type argument",
			Long: `
cancel_trade cancel trades by given orderID, market or add active depends of type argument.
	Params:
		trade_type - the type of cancellation, Valid Types: 'All',  'Trade', 'TradePair'
		trade_pair - the Cryptopia tradepair symbol of trades to cancel e.g. 'SKY/BTC', ”null” if empty
		orderID - the order identifier of trade to cancel, ”null” if empty
`,
			Example:            "cryptopia cancel_trade <trade_type> <trade_pair> <orderID>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				tradeType := args[0]
				var tradePair *string
				var orderID *int
				if args[1] != null {
					tradePair = &args[1]
				}
				if args[2] != null {
					tmp, err := strconv.Atoi(args[2])
					if err != nil {
						printErrorWithExit(err)
						return
					}
					orderID = &tmp
				}
				orders, err := client.CancelTrade(tradeType, tradePair, orderID)
				handleResult(struct {
					Orders []int
				}{Orders: orders}, err)
			},
		},
		"submit_tip": {
			Use:   "submit_tip",
			Short: "submit_tip submits a tip to Trollbox",
			Long: `
submit_tip submits a tip to Trollbox
	Params:
		currency - the currency symbol of the coins e.g. 'SKY'
		activeUsers - the amount of last active users to tip (Min: 2 Max: 100)
		amount - the amount of coins to buy e.g. 123.00000000
`,
			Example:            "cryptopia submit_tip <currency> <active_users> <amount>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				currency := args[0]
				activeUsers, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				amount, err := decimal.NewFromString(args[2])
				if err != nil {
					printErrorWithExit(err)
					return
				}

				result, err := client.SubmitTip(currency, activeUsers, amount)
				handleResult(struct {
					Result string
				}{Result: result}, err)
			},
		},
		"submit_withdraw": {
			Use:   "submit_withdraw",
			Short: "submit_withdraw submits a withdrawal request.",
			Long: `
submit_withdraw submits a withdrawal request. If address does not exists in you AddressBook, it will fail
paymentID will be used only for currencies, based of CryptoNote algorithm
	Params:
		currency - the currency symbol of the coins e.g. 'SKY'
		address - the address to send the currency to
		paymentID - the unique paimentID to identify the payment
		amount - the amount of coins to withdraw e.g. 123.00000000
`,
			Example:            "cryptopia submit_withdraw <currency> <address> <paymentID> <amount>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(4),
			Run: func(cmd *cobra.Command, args []string) {
				currency := args[0]
				address := args[1]
				paymentID := args[2]
				amount, err := decimal.NewFromString(args[3])
				if err != nil {
					printErrorWithExit(err)
					return
				}

				result, err := client.SubmitWithdraw(currency, address, paymentID, amount)
				handleResult(struct {
					Result int
				}{Result: result}, err)
			},
		},
		"submit_transfer": {
			Use:   "submit_transfer",
			Short: "submit_transfer submits a transfer funds to another user",
			Long: `
submit a transfer funds to another user
	Params:
		currency - the currency symbol of the coins e.g. 'SKY'
		userName - the Cryptopia username of the person to transfer the funds to
		amount - the amount of coins to withdraw e.g. 123.00000000
`,
			Example:            "cryptopia submit_transfer <currency> <userName> <amount>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				currency := args[0]
				userName := args[1]
				amount, err := decimal.NewFromString(args[2])
				if err != nil {
					printErrorWithExit(err)
					return
				}

				result, err := client.SubmitTransfer(currency, userName, amount)
				handleResult(struct {
					Result string
				}{Result: result}, err)
			},
		},
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Panicf("failed to execute root cobra command. err: %v", err)
	}
}

func init() {
	var key, secret string
	if os.Getenv("CRYPTOPIA_API_KEY") != "" && os.Getenv("CRYPTOPIA_API_SECRET") != "" {
		key = os.Getenv("CRYPTOPIA_API_KEY")
		secret = os.Getenv("CRYPTOPIA_API_SECRET")
	} else {
		currentUser, err := user.Current()
		if err != nil {
			log.Panicf("failed to get the new user. err: %v", err)
		}
		var config = filepath.Join(currentUser.HomeDir, ".exchangectl/config.toml")
		viper.SetConfigFile(config)
		err = viper.ReadInConfig()
		if err != nil {
			log.Panicf("failed to read config from %v. err: %v", config, err)
		}
		key = viper.GetString("cryptopia.key")
		if key == "" {
			log.Panic("cryptopia key is empty")
		}
		secret = viper.GetString("cryptopia.secret")
		if secret == "" {
			log.Panic("cryptopia secret is empty")
		}
	}

	client = cryptopia.NewAPIClient(key, secret)
	rootCmd = &cobra.Command{Use: "cryptopia"}

	for _, command := range getCommands() {
		rootCmd.AddCommand(command)
	}
}

func printResultWithExit(res interface{}) {
	output, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		fmt.Println("Error formating result to JSON. Error:", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", output)
	os.Exit(0)
}

func printErrorWithExit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func handleResult(res interface{}, err error) {
	if err != nil {
		printErrorWithExit(err)
	} else {
		printResultWithExit(res)
	}
}
