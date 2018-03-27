package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"os/user"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/skycoin/exchange-api/exchange/c2cx"
)

var (
	client  c2cx.Client
	rootCmd *cobra.Command
)

const null = "null"

func getCommands() map[string]*cobra.Command {
	return map[string]*cobra.Command{
		"getOrderBook": {
			Use:   "get_orderbook",
			Short: "gets all open orders by given symbol",
			Long: `
GetOrderbook gets all open orders by given symbol This method does not required API key and signing.
	Params:
	trade_pair -  market trade pair`,
			Example: "c2cx get_orderbook <trade_pair>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				res, err := client.GetOrderbook(c2cx.TradePair(args[0]))
				handleResult(res, err)
			},
		},
		"getOrderInfo": {
			Use:   "get_orderinfo",
			Short: "GetOrderInfo returns extended information about given order",
			Long: `
GetOrderInfo returns extended information about given order
	Params:
		trade_pair -  market trade pair,
		orderID - Id of wanted order info`,
			Example: "c2cx get_orderinfo <trade_pair> <orderID>",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := args[0]
				orderID, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				res, err := client.GetOrderInfo(c2cx.TradePair(tradePair), c2cx.OrderID(orderID))
				handleResult(res, err)
			},
		},
		"cancelAll": {
			Use:   "cancel_all",
			Short: "CancelAll cancels all executed orders for an orderbook.",
			Long: `
CancelAll cancels all executed orders for an orderbook. If it encounters an error,
it aborts and returns the order IDs that had been cancelled to that point.
	Params:
		trade_pair -  market trade pair`,
			Example: "c2cx cancel_all <trade_pair>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := args[0]
				res, err := client.CancelAll(c2cx.TradePair(tradePair))
				handleResult(res, err)
			},
		},
		"cancelMultiple": {
			Use:   "cancel_multiple",
			Short: "CancelMultiple cancels multiple orders",
			Long: `
CancelMultiple cancels multiple orders. It will try to cancel all of them, not stopping for any individual error. 
If any orders failed to cancel, a CancelMultiError is returned along with the array of order IDs 
which were successfully cancelled.
	Params:
		orderID[] - list Ids of orders for cancellation`,
			Example: "c2cx cancel_multiple <orderID> <orderID> <orderID>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				orderIDs := make([]c2cx.OrderID, 0)
				for _, arg := range args {
					orderID, err := strconv.Atoi(arg)
					if err != nil {
						printErrorWithExit(err)
					}
					orderIDs = append(orderIDs, c2cx.OrderID(orderID))
				}
				res, err := client.CancelMultiple(orderIDs)
				handleResult(res, err)
			},
		},
		"cancelOrder": {
			Use:   "cancel_order",
			Short: "CancelOrder cancels order with given orderID",
			Long: `
CancelOrder cancel order with given orderID
	Params:
		orderID - ID of order for cancellation`,
			Example: "c2cx cancel_order <orderID>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				orderID, err := strconv.Atoi(args[0])
				if err != nil {
					printErrorWithExit(err)
				}
				err = client.CancelOrder(c2cx.OrderID(orderID))
				handleResult(map[string]string{"result": "OK"}, err)
			},
		},
		"createOrder": {
			Use: "create_order",
			Short: `CreateOrder creates order with given orderType and parameters advanced is a advanced options for order
					creation if advanced is nil, isAdvancedOrder sets to zero, else advanced will be used as advanced options`,
			Long: `
CreateOrder creates order with given orderType and parameters advanced is a advanced options 
for order creation if advanced is nil, isAdvancedOrder sets to zero, else advanced will be used as advanced options
	Params:
		trade_pair -  market trade pair,
		price - order Price,
		quantity - quantity of order
		order_type - “buy” or “sell”
		price_type - “limit” for limit orders, “market” for market orders
		customerID - user submitted id
		
		advanced params:
			take_profit - take profit price, ”null” if empty
			stop_loss - stop loss price, ”null” if empty
			trigger_price - trigger price, ”null” if empty`,
			Example: "c2cx create_order <symbol> <price> <quantity> <order_type> <price_type> <customerID> <take_profit> <stop_loss> <trigger_price>",
			Args:    cobra.MinimumNArgs(9),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := c2cx.TradePair(args[0])
				price, err := decimal.NewFromString(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				quantity, err := decimal.NewFromString(args[2])
				if err != nil {
					printErrorWithExit(err)
				}
				orderType := c2cx.OrderType(args[3])
				priceType := c2cx.PriceType(args[4])
				customerID := args[5]

				advancedParams := &c2cx.AdvancedOrderParams{}

				if args[6] != null {
					takeProfit, err := decimal.NewFromString(args[6])
					if err != nil {
						printErrorWithExit(err)
					}
					advancedParams.TakeProfit = &takeProfit
				}

				if args[7] != null {
					stopLoss, err := decimal.NewFromString(args[7])
					if err != nil {
						printErrorWithExit(err)
					}
					advancedParams.StopLoss = &stopLoss
				}

				if args[8] != null {
					triggerPrice, err := decimal.NewFromString(args[8])
					if err != nil {
						printErrorWithExit(err)
					}
					advancedParams.TriggerPrice = &triggerPrice
				}

				orderID, err := client.CreateOrder(tradePair, price, quantity, orderType, priceType, &customerID, advancedParams)
				handleResult(map[string]c2cx.OrderID{"orderID": orderID}, err)
			},
		},
		"getBalanceSummary": {
			Use:   "get_balance_summary",
			Short: "GetBalanceSummary returns user balance for all available currencies",
			Long: `
GetBalanceSummary returns user balance for all available currencies
	Params:
		----`,
			Example: "c2cx get_balance_summary",
			Args:    cobra.MinimumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				res, err := client.GetBalanceSummary()
				handleResult(res, err)
			},
		},
		"getOrderByStatus": {
			Use:   "get_order_by_status",
			Short: "GetOrderByStatus get all orders with given status.",
			Long: `
GetOrderByStatus get all orders with given status. Makes multiple calls in the event of pagination. 
NOTE: GetOrderByStatus may returns orders with a different status than specified
	Params:
		trade_pair - market trade pair
		order_status - requested orders status
			0 = all,
			2=Active,
			3=Partially Completed,
			4=completed,
			5=cancelled,
			6=Suspended`,
			Example: "c2cx get_order_by_status <trade_pair> <order_status>",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := args[0]
				orderStatus, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				orders, err := client.GetOrderByStatus(c2cx.TradePair(tradePair), c2cx.OrderStatus(orderStatus))
				handleResult(orders, err)
			},
		},
		"getOrderByStatusPaged": {
			Use:   "get_order_by_status_paged",
			Short: "GetOrderByStatusPaged get all orders with given status for a given pagination page.",
			Long: `
GetOrderByStatusPaged get all orders with given status for a given pagination page.
	Params:
		trade_pair - market trade pair
		order_status - requested orders status
		page - page number for request with pagination`,
			Example: "c2cx get_order_by_status_paged <trade_pair> <order_status> <page>",
			Args:    cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := args[0]
				orderStatus, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				page, err := strconv.Atoi(args[2])
				if err != nil {
					printErrorWithExit(err)
				}
				orders, pageCount, err := client.GetOrderByStatusPaged(c2cx.TradePair(tradePair), c2cx.OrderStatus(orderStatus), page)
				handleResult(c2cx.Orders{Orders: orders, Page: pageCount}, err)
			},
		},
		"getOrderInfoAll": {
			Use: "get_order_info_all",
			Short: `
GetOrderInfoAll returns extended information about all orders Returns a 400 if it decides there are
no orders (there may be orders but it can disagree).`,
			Long: `
GetOrderInfoAll returns extended information about all orders Returns a 400
if it decides there are no orders (there may be orders but it can disagree).
	Params:
		trade_pair - market trade pair`,
			Example: "c2cx get_order_info_all <trade_pair>",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := args[0]
				orders, err := client.GetOrderInfoAll(c2cx.TradePair(tradePair))
				handleResult(c2cx.Orders{Orders: orders}, err)
			},
		},
		"limitBuy": {
			Use:   "limit_buy",
			Short: "LimitBuy place limit buy order",
			Long: `
LimitBuy place limit buy order.
	Params:
		trade_pair - market trade pair
		price - order price
		amount - amount of buying currency
		customerID - user submitted id`,
			Example: "c2cx limit_buy <trade_pair> <price> <amount> <customerID>",
			Args:    cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				symbol := args[0]
				price, err := decimal.NewFromString(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				amount, err := decimal.NewFromString(args[2])
				if err != nil {
					printErrorWithExit(err)
				}
				customerID := args[3]
				orderID, err := client.LimitBuy(c2cx.TradePair(symbol), price, amount, &customerID)
				handleResult(map[string]c2cx.OrderID{"orderID": orderID}, err)
			},
		},
		"limitSell": {
			Use:   "limit_sell",
			Short: "LimitSell place limit sell order",
			Long: `
LimitSell place limit sell order.
	Params:
		trade_pair - market trade pair
		price - order price
		amount - amount of buying currency
		customerID - user submitted id`,
			Example: "c2cx limit_sell <trade_pair> <price> <amount> <customerID>",
			Args:    cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				tradePair := args[0]
				price, err := decimal.NewFromString(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				amount, err := decimal.NewFromString(args[2])
				if err != nil {
					printErrorWithExit(err)
				}
				customerID := args[3]
				orderID, err := client.LimitSell(c2cx.TradePair(tradePair), price, amount, &customerID)
				handleResult(map[string]c2cx.OrderID{"orderID": orderID}, err)
			},
		},
		"marketBuy": {
			Use:   "market_buy",
			Short: "MarketBuy place market buy order",
			Long: `
MarketBuy place market buy order. 
A market buy order will sell the entire amount of the trade pair's first coin in exchange for the second coin. 
e.g. for BTC_SKY, the amount is the amount of BTC you want to spend on SKY.
	Params:
		trade_pair - market trade pair
		amount - amount of buying currency
		customerID - user submitted id`,
			Example: "c2cx market_buy <trade_pair> <amount> <customerID>",
			Args:    cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				symbol := args[0]
				amount, err := decimal.NewFromString(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				customerID := args[2]
				orderID, err := client.MarketBuy(c2cx.TradePair(symbol), amount, &customerID)
				handleResult(map[string]c2cx.OrderID{"orderID": orderID}, err)
			},
		},
		"marketSell": {
			Use:   "market_sell",
			Short: "MarketSell place market sell order",
			Long: `
MarketSell place market sell order. A market sell order will sell the entire amount 
of the trade pair's second coin in exchange for the first coin. e.g. for BTC_SKY, 
the amount is the amount of SKY you want to sell for BTC.
	Params:
		trade_pair - market trade pair
		amount - amount of selling currency
		customerID - user submitted id`,
			Example: "c2cx market_sell <trade_pair> <amount> <customerID>",
			Args:    cobra.MinimumNArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				symbol := args[0]
				amount, err := decimal.NewFromString(args[1])
				if err != nil {
					printErrorWithExit(err)
				}
				customerID := args[2]
				orderID, err := client.MarketSell(c2cx.TradePair(symbol), amount, &customerID)
				handleResult(map[string]c2cx.OrderID{"orderID": orderID}, err)
			},
		},
	}
}

func init() {
	user, err := user.Current()
	if err != nil {
		log.Panicf("failed to get the current user. err: %v", err)
	}

	var config = filepath.Join(user.HomeDir, ".exchangectl/config.toml")
	viper.SetConfigFile(config)
	err = viper.ReadInConfig()
	if err != nil {
		log.Panicf("failed to read the config file %s, err: %v", config, err)
	}
	key := viper.GetString("c2cx.key")
	if key == "" {
		panic("key param is empty")
	}
	secret := viper.GetString("c2cx.secret")
	if secret == "" {
		panic("secret param is empty")
	}
	client = c2cx.Client{
		Key:    key,
		Secret: secret,
	}
	rootCmd = &cobra.Command{Use: "c2cx"}
	for _, v := range getCommands() {
		rootCmd.AddCommand(v)
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

func main() {
	rootCmd.Execute()
}
