package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/uberfurrer/tradebot/exchange"
	"github.com/urfave/cli"
)

func cancelCmd() cli.Command {
	name := "cancel"
	return cli.Command{
		Name:  name,
		Usage: "type [args]",
		Subcommands: []cli.Command{
			cancelTradeCmd(),
			cancelMarketCmd(),
			cancelAllCmd(),
		},
	}
}

// endpoint cancel_trade
// params {"orderid":1}
func cancelTradeCmd() cli.Command {
	name := "trade"
	return cli.Command{
		Name:      name,
		ArgsUsage: "[orderids...]",
		Usage:     "Cancel orders with specified orderids",
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) == 0 {
				return errInvalidParams
			}
			var orderids = make([]int, len(args))
			for i, v := range args {
				id, err := strconv.Atoi(v)
				if err != nil {
					fmt.Printf("Error: orderid must be integer, given %v, position %d", v, i+1)
					return errInvalidParams
				}
				orderids[i] = id
			}

			for _, v := range orderids {
				params, _ := json.Marshal(map[string]int{"orderid": v})
				var req = rpc.Request{
					ID:      reqID(),
					JSONRPC: rpc.JSONRPC,
					Method:  "cancel_trade",
					Params:  params,
				}
				resp, err := rpc.Do(rpcaddr, endpoint, req)
				if err != nil {
					fmt.Printf("Error processing request %s", err.Error())
					return errRPC
				}
				var orderinfo exchange.Order
				if err := json.Unmarshal(resp.Result, &orderinfo); err != nil {
					fmt.Printf("Error: invalid response format, %s", err.Error())
					return errInvalidResponse
				}
				b, _ := json.MarshalIndent(orderinfo, "", "    ")
				fmt.Println(string(b))
			}
			return nil
		},
	}
}

// endpoint cancel_market
// params {"symbol":"CNY_BTC"}
// symbol isnt case-sensetive
func cancelMarketCmd() cli.Command {
	name := "market"
	return cli.Command{
		Name:      name,
		ArgsUsage: "[tradepairs...]",
		Usage:     "Cancel orders with specified tradepairs",
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) == 0 {
				return errInvalidParams
			}
			for _, v := range args {
				params, _ := json.Marshal(map[string]string{"symbol": v})
				var req = rpc.Request{
					ID:      reqID(),
					JSONRPC: rpc.JSONRPC,
					Method:  "cancel_market",
					Params:  params,
				}
				resp, err := rpc.Do(rpcaddr, endpoint, req)
				if err != nil {
					fmt.Printf("Error processing request, %s", err)
					return errRPC
				}
				var orders []exchange.Order
				if err := json.Unmarshal(resp.Result, &orders); err != nil {
					fmt.Printf("Error: invalid response format %s", err)
					return errInvalidResponse
				}
				for _, info := range orders {
					b, _ := json.MarshalIndent(info, "", "    ")
					fmt.Println(string(b))
				}
			}
			return nil
		},
	}
}

// endpoint cancel_all
// params: null
func cancelAllCmd() cli.Command {
	name := "all"
	return cli.Command{
		Name:  name,
		Usage: "Cancel all executed orders",
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) != 0 {
				return errInvalidParams
			}
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "cancel_all",
				Params:  nil,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request %s", err)
				return errRPC
			}
			var orders []exchange.Order
			if err := json.Unmarshal(resp.Result, &orders); err != nil {
				fmt.Printf("Error: invalid response %s", err)
				return errInvalidResponse
			}
			for _, v := range orders {
				b, _ := json.MarshalIndent(v, "", "    ")
				fmt.Println(string(b))
			}
			return nil
		},
	}
}
