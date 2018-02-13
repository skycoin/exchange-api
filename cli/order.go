package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/exchange-api/exchange"
	"github.com/urfave/cli"
)

const (
	allName    = "all"
	marketName = "market"
)

func orderInfoCMD() cli.Command {
	name := "info"
	return cli.Command{
		Name:      name,
		Usage:     "Print information about order",
		ArgsUsage: "<orderid>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			orderid, err := strconv.Atoi(c.Args().First())
			if err != nil {
				return err
			}
			params := map[string]interface{}{
				"orderid": orderid,
			}
			resp, err := rpcRequest("order_info", params)
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
func orderStatusCMD() cli.Command {
	name := "status"
	return cli.Command{
		Name:      name,
		Usage:     "Print order status",
		ArgsUsage: "<orderid>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			orderid, err := strconv.Atoi(c.Args().First())
			if err != nil {
				return err
			}
			params := map[string]interface{}{
				"orderid": orderid,
			}
			resp, err := rpcRequest("order_status", params)
			if err != nil {
				return err
			}
			fmt.Printf("Order %d status: %s\n", orderid, resp)
			return nil
		},
	}
}

func completedAllCMD() cli.Command {
	name := allName
	var short bool
	return cli.Command{
		Name:  name,
		Usage: "Print all completed orders",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return errInvalidInput
			}
			resp, err := rpcRequest("completed", nil)
			if err != nil {
				return err
			}
			var orderids []int
			err = json.Unmarshal(resp, &orderids)
			if err != nil {
				return err
			}
			for _, v := range orderids {
				var order exchange.Order
				params := map[string]interface{}{
					"orderid": v,
				}
				resp, err := rpcRequest("order_info", params)
				if err != nil {
					continue
				}
				if err = json.Unmarshal(resp, &order); err != nil {
					panic(err)
				}
				if short {
					fmt.Println(orderShort(order))
				} else {
					fmt.Println(orderFull(order))
				}
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "short", Destination: &short, Usage: "Short output"},
		},
	}
}

func completedMarketCMD() cli.Command {
	name := marketName
	var short bool
	return cli.Command{
		Name:      name,
		Usage:     "Print all completed orders in market",
		ArgsUsage: "<market>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			resp, err := rpcRequest("completed", nil)
			if err != nil {
				return err
			}
			var orderids []int
			err = json.Unmarshal(resp, &orderids)
			if err != nil {
				return err
			}
			market := strings.ToUpper(strings.Replace(c.Args().First(), "_", "/", -1))
			for _, v := range orderids {
				var order exchange.Order
				params := map[string]interface{}{
					"orderid": v,
				}
				resp, err := rpcRequest("order_info", params)
				if err != nil {
					continue
				}
				if err = json.Unmarshal(resp, &order); err != nil {
					panic(err)
				}
				if order.Market != market {
					continue
				}
				if short {
					fmt.Println(orderShort(order))
				} else {
					fmt.Println(orderFull(order))
				}
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "short", Destination: &short, Usage: "Short output"},
		},
	}
}

func executedAllCMD() cli.Command {
	name := allName
	var short bool
	return cli.Command{
		Name:  name,
		Usage: "Print all opened orders",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				return errInvalidInput
			}
			resp, err := rpcRequest("executed", nil)
			if err != nil {
				return err
			}
			var orderids []int
			err = json.Unmarshal(resp, &orderids)
			if err != nil {
				return err
			}
			for _, v := range orderids {
				var order exchange.Order
				params := map[string]interface{}{
					"orderid": v,
				}
				resp, err := rpcRequest("order_info", params)
				if err != nil {
					continue
				}
				if err = json.Unmarshal(resp, &order); err != nil {
					panic(err)
				}
				if short {
					fmt.Println(orderShort(order))
				} else {
					fmt.Println(orderFull(order))
				}
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "short", Destination: &short, Usage: "Short output"},
		},
	}
}
func executedMarketCMD() cli.Command {
	name := marketName
	var short bool
	return cli.Command{
		Name:      name,
		Usage:     "Print all opened orders in market",
		ArgsUsage: "<market>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			resp, err := rpcRequest("executed", nil)
			if err != nil {
				return err
			}
			var orderids []int
			err = json.Unmarshal(resp, &orderids)
			if err != nil {
				return err
			}
			market := strings.ToUpper(strings.Replace(c.Args().First(), "_", "/", -1))
			for _, v := range orderids {
				var order exchange.Order
				params := map[string]interface{}{
					"orderid": v,
				}
				resp, err := rpcRequest("order_info", params)
				if err != nil {
					continue
				}
				if err = json.Unmarshal(resp, &order); err != nil {
					panic(err)
				}
				if order.Market != market {
					continue
				}
				if short {
					fmt.Println(orderShort(order))
				} else {
					fmt.Println(orderFull(order))
				}
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "short", Destination: &short, Usage: "Short output"},
		},
	}
}

func orderCMDs() cli.Command {
	return cli.Command{
		Name:  "order",
		Usage: "Prints information about order",
		Subcommands: cli.Commands{
			orderInfoCMD(),
			orderStatusCMD(),
		},
	}
}

func completedCMDs() cli.Command {
	return cli.Command{
		Name:  "completed",
		Usage: "Print completed orders",
		Subcommands: cli.Commands{
			completedAllCMD(),
			completedMarketCMD(),
		},
	}
}

func executedCMDs() cli.Command {
	return cli.Command{
		Name:  "executed",
		Usage: "Print opened orders",
		Subcommands: cli.Commands{
			executedAllCMD(),
			executedMarketCMD(),
		},
	}
}
