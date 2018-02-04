package cli

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/urfave/cli"
)

func buyCMD() cli.Command {
	name := "buy"
	return cli.Command{
		Name:      name,
		Usage:     "Place buy order",
		ArgsUsage: "<symbol> <price> <amount>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 3 {
				return errInvalidInput
			}
			var (
				err           error
				symbol        string
				price, amount decimal.Decimal
			)
			symbol = c.Args().Get(0)
			if price, err = decimal.NewFromString(c.Args().Get(1)); err != nil {
				return err
			}
			if amount, err = decimal.NewFromString(c.Args().Get(2)); err != nil {
				return err
			}
			params := map[string]interface{}{
				"symbol": symbol,
				"price":  price,
				"amount": amount,
			}

			resp, err := rpcRequest("buy", params)
			if err != nil {
				return err
			}
			fmt.Printf("Order %s created\n", resp)
			return nil
		},
	}
}

func sellCMD() cli.Command {
	name := "sell"
	return cli.Command{
		Name:      name,
		Usage:     "Place sell order",
		ArgsUsage: "<symbol> <price> <amount>",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 3 {
				return errInvalidInput
			}
			var (
				err           error
				symbol        string
				price, amount decimal.Decimal
			)
			symbol = c.Args().Get(0)
			if price, err = decimal.NewFromString(c.Args().Get(1)); err != nil {
				return err
			}
			if amount, err = decimal.NewFromString(c.Args().Get(2)); err != nil {
				return err
			}
			params := map[string]interface{}{
				"symbol": symbol,
				"price":  price,
				"amount": amount,
			}

			resp, err := rpcRequest("sell", params)
			if err != nil {
				return err
			}
			fmt.Printf("Order %s created\n", resp)
			return nil
		},
	}
}
