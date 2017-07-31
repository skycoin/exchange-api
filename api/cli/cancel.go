package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/uberfurrer/tradebot/exchange"
	"github.com/urfave/cli"
)

func cancelTradeCMD() cli.Command {
	var name = "trade"
	return cli.Command{
		Name:      name,
		Usage:     "Cancel order",
		ArgsUsage: "<orderid>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			orderid, err := strconv.Atoi(c.Args().First())
			if err != nil {
				return err
			}
			var params = map[string]interface{}{
				"orderid": orderid,
			}
			resp, err := rpcRequest("cancel_trade", params)
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
func cancelMarketCMD() cli.Command {
	var name = "market"
	return cli.Command{
		Name:      name,
		Usage:     "Cancel all orders in market",
		ArgsUsage: "<market>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			var params = map[string]interface{}{
				"symbol": c.Args().First(),
			}
			resp, err := rpcRequest("cancel_market", params)
			if err != nil {
				return err
			}
			var orders []exchange.Order
			err = json.Unmarshal(resp, &orders)
			if err != nil {
				return err
			}
			for _, v := range orders {
				fmt.Println(orderShort(v))
			}
			fmt.Printf("Cancelled %d orders", len(orders))
			return nil
		},
	}
}
func cancelAllCMD() cli.Command {
	var name = "all"
	return cli.Command{
		Name:  name,
		Usage: "Cancel all orders",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return errInvalidInput
			}
			resp, err := rpcRequest("cancel_all", nil)
			if err != nil {
				return err
			}
			var orders []exchange.Order
			err = json.Unmarshal(resp, &orders)
			if err != nil {
				return err
			}
			for _, v := range orders {
				fmt.Println(orderShort(v))
			}
			fmt.Printf("Cancelled %d orders", len(orders))
			return nil
		},
	}
}

func cancelCMDs() cli.Command {
	return cli.Command{
		Name:  "cancel",
		Usage: "Cancel order(s)",
		Subcommands: cli.Commands{
			cancelTradeCMD(),
			cancelMarketCMD(),
			cancelAllCMD(),
		},
	}
}
