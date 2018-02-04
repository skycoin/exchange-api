package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/skycoin/exchange-api/exchange/cryptopia.co.nz"

	"github.com/urfave/cli"

	"github.com/skycoin/exchange-api/exchange/cryptopia.co.nz"
)

func submitWithdrawCMD() cli.Command {
	name := "withdraw"
	return cli.Command{
		Name:      name,
		Usage:     "Withdraw funds to address",
		ArgsUsage: "<address> <currency> <amount> [paymentid]",
		Action: func(c *cli.Context) error {
			if c.NArg() < 3 || c.NArg() > 4 {
				return errInvalidInput
			}
			var (
				err    error
				amount decimal.Decimal
			)
			if amount, err = decimal.NewFromString(c.Args().Get(2)); err != nil {
				return err
			}
			params := map[string]interface{}{
				"address":  c.Args().Get(0),
				"currency": c.Args().Get(1),
				"amount":   amount,
			}
			if c.NArg() == 4 {
				params["paymentid"] = c.Args().Get(3)
			}
			resp, err := rpcRequest("withdraw", params)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			fmt.Printf("Withdrawal request ID %s\n", resp)
			return nil
		},
	}
}

func depositCMD() cli.Command {
	name := "deposit"
	return cli.Command{
		Name:      name,
		Usage:     "Print address for deposit",
		ArgsUsage: "<currency>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			params := map[string]interface{}{
				"currency": c.Args().First(),
			}
			resp, err := rpcRequest("deposit", params)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			var addr cryptopia.DepositAddress
			if err = json.Unmarshal(resp, &addr); err != nil {
				return err
			}
			str, _ := json.MarshalIndent(addr, "", "    ")
			fmt.Println(string(str))
			return nil
		},
	}
}

func transactionsCMD() cli.Command {
	name := "transactions"
	return cli.Command{
		Name:      name,
		Usage:     "Print list of transactions",
		ArgsUsage: "<type>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			params := map[string]interface{}{
				"type": strings.Title(c.Args().First()),
			}
			resp, err := rpcRequest("transactions", params)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			var txs []cryptopia.Transaction
			err = json.Unmarshal(resp, &txs)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			str, _ := json.MarshalIndent(txs, "", "    ")
			fmt.Println(string(str))
			return nil
		},
	}
}

func trackingAddCMD() cli.Command {
	name := "add"
	return cli.Command{
		Name:      name,
		Usage:     "Add market to orderbook tracking list",
		ArgsUsage: "<market>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			params := map[string]interface{}{
				"market": c.Args().First(),
			}
			_, err := rpcRequest("tracking_add", params)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			fmt.Printf("Market %s added to orderbook tracking list\n", c.Args().First())
			return nil
		},
	}
}

func trackingRemoveCMD() cli.Command {
	name := "remove"
	return cli.Command{
		Name:      name,
		Usage:     "Remove market from orderbook tracking list",
		ArgsUsage: "<market>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			params := map[string]interface{}{
				"market": c.Args().First(),
			}
			_, err := rpcRequest("tracking_rm", params)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			fmt.Printf("Market %s removed from orderbook tracking list\n", c.Args().First())
			return nil
		},
	}
}

func trackingCMDs() cli.Command {
	return cli.Command{
		Name:  "tracking",
		Usage: "Manage tracked markets",
		Subcommands: cli.Commands{
			trackingAddCMD(),
			trackingRemoveCMD(),
		},
	}
}
