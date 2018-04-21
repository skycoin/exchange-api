package main

import (
	"github.com/spf13/cobra"
	"github.com/skycoin/exchange-api/exchange/cryptopia"
	"os/user"
	"log"
	"path/filepath"
	"github.com/spf13/viper"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	rootCmd *cobra.Command
	client  *cryptopia.Client
)

func getCommands() map[string]*cobra.Command {
	return map[string]*cobra.Command{
		"get_currencies": {
			Use:     "get_currencies",
			Short:   "get_currencies gets all currencies",
			Example: "cryptopia get_currencies",
			Run: func(cmd *cobra.Command, args []string) {
				currencies, err := client.GetCurrencies()
				handleResult(currencies, err)
			},
		},
		"get_trade_pairs": {
			Use:     "get_trade_pairs",
			Short:   "get_trade_pairs gets all TradePairs on exchange",
			Example: "cryptopia get_trade_pairs",
			Run: func(cmd *cobra.Command, args []string) {
				tradePairs, err := client.GetTradePairs()
				handleResult(tradePairs, err)
			},
		},
		"get_markets": {
			Use:     "get_markets",
			Short:   "GetMarkets return all Market info by given baseMarket",
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
			Use:     "get_market",
			Short:   "get_market returns market with given label",
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
			Use:     "get_market_history",
			Short:   "get_market_history returns market history with given label",
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
			Use:                "get_market_orders",
			Short:              "get_market_orders returns count orders from market with given label",
			Example:            "cryptopia get_market_orders <market> <count>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				market := args[0]
				count, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					printErrorWithExit(err)
					return
				}
				marketsResult, err := client.GetMarketOrders(market, int(count))
				handleResult(marketsResult, err)
			},
		},
		"get_market_order_groups": {
			Use:                "get_market_order_groups",
			Short:              "get_market_order_groups returns count orders from market with given label",
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
				marketsResult, err := client.GetMarketOrderGroups(int(count), markets)
				handleResult(marketsResult, err)
			},
		},
		"get_balance": {
			Use:                "get_balance",
			Short:              "get_balance returns a string representation of balance by given currency",
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
			Use:                "get_deposit_address",
			Short:              "get_deposit_address returns a deposit address of given currency",
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
			Use:                "get_open_orders",
			Short:              "get_open_orders return a list of opened orders by specific market or all markets",
			Example:            "cryptopia get_open_orders <market> <count>",
			DisableFlagParsing: true,
			Args:               cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				market := args[0]
				count, err := strconv.Atoi(args[1])
				if err != nil {
					printErrorWithExit(err)
					return
				}
				depositAddress, err := client.GetOpenOrders(&market, &count)
				handleResult(depositAddress, err)
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
	key := viper.GetString("cryptopia.key")
	if key == "" {
		log.Panic("cryptopia key is empty")
	}
	secret := viper.GetString("cryptopia.secret")
	if secret == "" {
		log.Panic("cryptopia secret is empty")
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
