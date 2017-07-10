package main

import (
	"fmt"
	"strconv"

	"encoding/json"

	"github.com/urfave/cli"
)

var exchangeCliHandler = []cli.Command{
	cli.Command{
		Name: "cancel",
		Subcommands: []cli.Command{
			cli.Command{
				Name: "all",
				Action: func(c *cli.Context) error {
					if len(c.Args()) != 0 {
						fmt.Println("Error: This command does not have any arguments")
						return ErrInvalidArgs
					}
					result, err := makeRPCCall(c.GlobalString("endpoint"), "CancelAll", nil)
					if err != nil {
						fmt.Printf("Error: %s\n", err)
						return err
					}
					return printOrderInfoArr(result)
				},
				Description: "cancel all orders",
				Usage:       "all",
			},
			cli.Command{
				Name: "market",
				Action: func(c *cli.Context) error {
					if len(c.Args()) == 0 {
						fmt.Println("Error: this command required markets, thus be cancelled")
						return ErrInvalidArgs
					}
					for _, v := range c.Args() {
						// cancel here
						var params = make(map[string]string)
						params["symbol"] = v
						data, err := json.Marshal(params)
						if err != nil {
							///
						}
						result, err := makeRPCCall(c.GlobalString("endpoint"), "CancelMarket", data)
						if err != nil {
							fmt.Printf("Error: %s\n", err)
						}
						printOrderInfoArr(result)
					}
					return nil
				},
				Usage:       "market [TradePair TradePair...]",
				Description: "market cancels all orders in given markets",
			},
			cli.Command{
				Name: "order",
				Action: func(c *cli.Context) error {
					if len(c.Args()) == 0 {
						fmt.Println("Error: this command required orderids, thus be cancelled")
						return ErrInvalidArgs
					}
					for _, v := range c.Args() {
						orderID, err := strconv.Atoi(v)
						if err != nil {
							fmt.Printf("Warning: input orderID %s isn't numeric, ignored\n", v)
							continue
						}
						// cancel here
						var params = make(map[string]int)
						params["orderid"] = orderID
						data, err := json.Marshal(params)
						if err != nil {
							///
						}
						result, err := makeRPCCall(c.GlobalString("endpoint"), "Cancel", data)
						if err != nil {
							fmt.Printf("Error: %s\n", err)
							continue
						}
						printOrderInfo(result)
					}
					return nil
				},
				Usage:       "order [OrderID OrderID....]",
				Description: "cancels order(s) by given orderid",
			},
		},
		Description: "cancel all or seprarated orders, executed on exchange",
		Usage:       "cancel [canceltype args...]",
	},
	cli.Command{
		Name: "buy",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 3 {
				fmt.Println("Error: this command accept three parameters: tradepair, price, amount")
				return ErrInvalidArgs
			}
			pair := c.Args().Get(0)
			price, err := strconv.ParseFloat(c.Args().Get(1), 64)
			if err != nil {
				fmt.Printf("Error: price isn't numeric, strconv.ParseFloat error %s\n", err.Error())
				return ErrInvalidArgs
			}
			amount, err := strconv.ParseFloat(c.Args().Get(2), 64)
			if err != nil {
				fmt.Printf("Error: amount isn't numeric, strconv.ParseFloat error %s\n", err.Error())
				return ErrInvalidArgs
			}
			var params = struct {
				Rate   float64 `json:"rate"`
				Amount float64 `json:"amount"`
				Symbol string  `json:"symbol"`
			}{Rate: price, Amount: amount, Symbol: pair}
			data, err := json.Marshal(params)
			if err != nil {
				///
			}
			result, err := makeRPCCall(c.GlobalString("endpoint"), "Buy", data)
			if err != nil {
				///
			}
			printOrderID(result)
			return nil
		},
		Usage:       "buy [tradepair rate amount]",
		Description: "place buy order on exchange",
	},
	cli.Command{
		Name: "sell",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 3 {
				fmt.Println("Error: this command accept three parameters: tradepair, price, amount")
				return ErrInvalidArgs
			}
			pair := c.Args().Get(0)
			price, err := strconv.ParseFloat(c.Args().Get(1), 64)
			if err != nil {
				fmt.Printf("Error: price isn't numeric, strconv.ParseFloat error %s\n", err.Error())
				return ErrInvalidArgs
			}
			amount, err := strconv.ParseFloat(c.Args().Get(2), 64)
			if err != nil {
				fmt.Printf("Error: amount isn't numeric, strconv.ParseFloat error %s\n", err.Error())
				return ErrInvalidArgs
			}
			var params = struct {
				Rate   float64 `json:"rate"`
				Amount float64 `json:"amount"`
				Symbol string  `json:"symbol"`
			}{Rate: price, Amount: amount, Symbol: pair}
			data, err := json.Marshal(params)
			if err != nil {
				///
			}
			result, err := makeRPCCall(c.GlobalString("endpoint"), "Buy", data)
			if err != nil {
				///
			}
			printOrderID(result)
			return nil
		},
		Usage:       "sell [tradepair rate amount]",
		Description: "place sell order on exchange",
	},
	cli.Command{
		Name: "status",
		Action: func(c *cli.Context) error {
			if len(c.Args()) == 0 {
				fmt.Println("Error: this command accepts list of orderids")
				return ErrInvalidArgs
			}
			for _, v := range c.Args() {
				orderid, err := strconv.Atoi(v)
				if err != nil {
					fmt.Printf("Warning: orderid %s isn't number, ignored\n", v)
					continue
				}
				// handle orderids here
				var params = map[string]int{
					"orderid": orderid,
				}
				data, err := json.Marshal(params)
				if err != nil {
					//
				}
				result, err := makeRPCCall(c.GlobalString("endpoint"), "OrderStatus", data)
				if err != nil {
					///
				}
				printString(result)
			}
			return nil
		},
		Usage:       "status [OrderID OrderID...]",
		Description: "gets statusees of order(s) by given OrderID(s)",
	},
	cli.Command{
		Name: "info",
		Action: func(c *cli.Context) error {
			if len(c.Args()) == 0 {
				fmt.Println("Error: this command accepts list of orderids")
				return ErrInvalidArgs
			}
			for _, v := range c.Args() {
				orderid, err := strconv.Atoi(v)
				if err != nil {
					fmt.Printf("Warning: orderid %s isn't number, ignored\n", v)
					continue
				}
				// handle orderids here
				var params = map[string]int{
					"orderid": orderid,
				}
				data, err := json.Marshal(params)
				if err != nil {
					///
				}
				result, err := makeRPCCall(c.GlobalString("endpoint"), "OrderDetails", data)
				if err != nil {
					///
				}
				printOrderInfo(result)

			}
			return nil
		},
		Usage:       "info [OrderID OrderID...]",
		Description: "print detailed info about orders, such as status, price, volume, time",
	},
	cli.Command{
		Name: "completed",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 0 {
				fmt.Println("Error: this command does not accept parameters")
				return ErrInvalidArgs
			}
			result, err := makeRPCCall(c.GlobalString("endpoint"), "Completed", nil)
			if err != nil {
				///
			}
			printOrderInfoArr(result)
			return nil
		},
		Usage:       "completed",
		Description: "returns all completed orders, created by this bot during current session",
	},
	cli.Command{
		Name: "active",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 0 {
				fmt.Println("Error: this command does not accept parameters")
				return ErrInvalidArgs
			}
			result, err := makeRPCCall(c.GlobalString("endpoint"), "Executed", nil)
			if err != nil {
				//
			}
			printOrderInfoArr(result)
			return nil
		},
		Usage:       "active",
		Description: "prints all executed orders, that was created by this client",
	},
	cli.Command{
		Name: "balance",
		Action: func(c *cli.Context) error {
			if len(c.Args()) == 0 {
				fmt.Println("Error: this command accept list of currencies")
				return ErrInvalidArgs
			}
			for _, v := range c.Args() {
				// get balance here
				var params = map[string]string{
					"currency": v,
				}
				data, err := json.Marshal(params)
				if err != nil {
					///
				}
				result, err := makeRPCCall(c.GlobalString("endpoint"), "GetBalance", data)
				if err != nil {
					///
				}
				printString(result)
			}
			return nil
		},
		Usage:       "balance [currency currency...]",
		Description: "gets balance of given currencies",
	},
	cli.Command{
		Name: "orderbook",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				fmt.Println("Error: this command accepts one tradepair symbol")
				return ErrInvalidArgs
			}
			var params = map[string]string{
				"market": c.Args().First(),
			}
			data, err := json.Marshal(params)
			if err != nil {
				///
			}
			result, err := makeRPCCall(c.GlobalString("endpoint"), "OrderBook", data)
			if err != nil {
				//
			}
			printOrderbook(result)
			return nil
		},
		Usage:       "orderbook tradepair",
		Description: "gets orderbook for given tradepair",
	},
}
