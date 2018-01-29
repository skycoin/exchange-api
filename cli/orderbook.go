package cli

import (
	"encoding/json"
	"fmt"

	"github.com/skycoin/exchange-api/exchange"
	"github.com/urfave/cli"
)

func orderbookCMD() cli.Command {
	var name = "orderbook"
	var short bool
	return cli.Command{
		Name:      name,
		Usage:     "Print orderbook",
		ArgsUsage: "<market>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			var params = map[string]interface{}{
				"symbol": c.Args().First(),
			}
			resp, err := rpcRequest("orderbook", params)
			if err != nil {
				return err
			}
			var orderbook exchange.MarketRecord
			err = json.Unmarshal(resp, &orderbook)
			if err != nil {
				return err
			}
			if short {
				fmt.Println(orderbookShort(orderbook))
			} else {
				fmt.Println(orderbookFull(orderbook))
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "short", Destination: &short, Usage: "Short output"},
		},
	}
}
