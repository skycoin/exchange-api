package cli

import (
	"encoding/json"
	"fmt"

	"github.com/uberfurrer/tradebot/exchange"
	"github.com/urfave/cli"
)

func balanceCMD() cli.Command {
	var name = "balance"
	return cli.Command{
		Name:      name,
		Usage:     "Print balance",
		ArgsUsage: "<currency>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			var params = map[string]interface{}{
				"currency": c.Args().First(),
			}
			resp, err := rpcRequest("balance", params)
			if err != nil {
				return err
			}
			var order exchange.Order
			err = json.Unmarshal(resp, &order)
			if err != nil {
				return err
			}
			fmt.Println(orderFull(order))
			return nil
		},
	}
}
