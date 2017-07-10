package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var app = cli.App{
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Value:  "",
			Hidden: true,
		},
	},
}

func main() {
	app.Run(os.Args)
}

func init() {
	app.Commands = []cli.Command{
		cli.Command{
			Name: "c2cx",
			Subcommands: append([]cli.Command{cli.Command{
				Name: "createorder",
				Action: func(c *cli.Context) error {
					fmt.Println(c.Args())
					return nil
				},
				Description: "allows to create limited orders",
			},
			}, exchangeCliHandler...),
			Before: func(c *cli.Context) error {
				c.GlobalSet("endpoint", "c2cx")
				return nil
			},
		},
		cli.Command{
			Name: "cryptopia",
			Subcommands: append([]cli.Command{
				cli.Command{
					Name: "deposit",
					Action: func(c *cli.Context) error {
						if len(c.Args()) != 1 {
							fmt.Println("Error: this command requires currency for getting address")
						}
						var params = map[string]string{
							"currency": c.Args().First(),
						}
						data, err := json.Marshal(params)
						if err != nil {
							//
						}
						result, err := makeRPCCall(c.GlobalString("endpoint"), "GetDepositAddress", data)
						if err != nil {
							//
						}
						printString(result)
						return nil
					},
					Usage:       "desposit [currency]",
					Description: "gets deposit address for given currency",
				},
				cli.Command{
					Name: "transactions",
					Action: func(c *cli.Context) error {
						if len(c.Args()) < 1 || len(c.Args()) > 2 {
							fmt.Println("Error: this command requires one mandatory and second optional argumnent")
							return ErrInvalidArgs
						}
						var count *int
						if v, err := strconv.Atoi(c.Args().Get(2)); err == nil {
							count = &v
						}
						var params = struct {
							Type  string `json:"type"`
							Count *int   `json:"count,omitempty"`
						}{
							Type:  c.Args().First(),
							Count: count,
						}
						data, err := json.Marshal(params)
						if err != nil {
							///
						}
						result, err := makeRPCCall(c.GlobalString("endpoint"), "GetTransactions", data)
						if err != nil {
							//
						}
						printString(result)
						return nil
					},
					Usage:       "transactions [type count(optional)]",
					Description: "gets list of deposites or withdrawals",
				},
				cli.Command{
					Name: "withdraw",
					Action: func(c *cli.Context) error {
						if len(c.Args()) < 3 || len(c.Args()) > 4 {
							fmt.Println("Error: this command requires 3 or 4 args")
						}
						var pid *string
						amount, err := strconv.ParseFloat(c.Args().Get(2), 64)
						if err != nil {
							//
						}
						if len(c.Args()) == 3 {
							v := c.Args().Get(3)
							pid = &v
						}
						var params = struct {
							Currency  string  `json:"currency"`
							Address   string  `json:"address"`
							PaymentID *string `json:"paymentid,omitempty"`
							Amount    float64 `json:"amount"`
						}{
							Currency:  c.Args().Get(1),
							Address:   c.Args().Get(2),
							PaymentID: pid,
							Amount:    amount,
						}
						data, err := json.Marshal(params)
						if err != nil {
							//
						}
						result, err := makeRPCCall(c.GlobalString("endpoint"), "SubmitWithdraw", data)
						if err != nil {
							///
						}
						printOrderID(result) ///Withdraw ID
						return nil
					},
					Usage:       "withdraw [currency address amount paymentid(optional)]",
					Description: "creates withdrawal request",
				},
			}, exchangeCliHandler...),
			Before: func(c *cli.Context) error {
				c.GlobalSet("endpoint", "cryptopia")
				return nil
			},
		},
	}
}
