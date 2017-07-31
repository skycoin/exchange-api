package cli

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli"
)

func buyCMD() cli.Command {
	var name = "buy"
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
				price, amount float64
			)
			symbol = c.Args().Get(0)
			if price, err = strconv.ParseFloat(c.Args().Get(1), 64); err != nil {
				return err
			}
			if amount, err = strconv.ParseFloat(c.Args().Get(2), 64); err != nil {
				return err
			}
			var params = map[string]interface{}{
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
	var name = "sell"
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
				price, amount float64
			)
			symbol = c.Args().Get(0)
			if price, err = strconv.ParseFloat(c.Args().Get(1), 64); err != nil {
				return err
			}
			if amount, err = strconv.ParseFloat(c.Args().Get(2), 64); err != nil {
				return err
			}
			var params = map[string]interface{}{
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
