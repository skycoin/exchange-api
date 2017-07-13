package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/urfave/cli"
)

// endpoint buy
// params; {"symbol":"CNY_BTC", "price":1.0, "amount": 1.0}
func buyCmd() cli.Command {
	var name = "buy"
	return cli.Command{
		Name: name,
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) != 3 {
				return errInvalidParams
			}
			var (
				price, amount float64
				err           error
			)
			if args[0] == "" {
				return errInvalidParams
			}
			if price, err = strconv.ParseFloat(args[1], 64); err != nil {
				return errInvalidParams
			}
			if amount, err = strconv.ParseFloat(args[2], 64); err != nil {
				return errInvalidParams
			}
			params, _ := json.Marshal(map[string]interface{}{
				"symbol": args[0],
				"price":  price,
				"amount": amount,
			})
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "buy",
				Params:  params,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request, %s", err)
				return errRPC
			}
			var orderid int
			if err = json.Unmarshal(resp.Result, &orderid); err != nil {
				fmt.Printf("Error: invalid response, %s", err)
				return errInvalidResponse
			}
			b, _ := json.MarshalIndent(map[string]int{"orderid": orderid}, "", "    ")
			fmt.Println(string(b))

			return nil
		},
		ArgsUsage: "[tradepair price amount]",
		Usage:     "buy places a buy order",
	}
}

// endpoint sell
// params; {"symbol":"CNY_BTC", "price":1.0, "amount": 1.0}}
func sellCmd() cli.Command {
	var name = "sell"
	return cli.Command{
		Name: name,
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) != 3 {
				return errInvalidParams
			}
			var (
				price, amount float64
				err           error
			)
			if args[0] == "" {
				return errInvalidParams
			}
			if price, err = strconv.ParseFloat(args[1], 64); err != nil {
				return errInvalidParams
			}
			if amount, err = strconv.ParseFloat(args[2], 64); err != nil {
				return errInvalidParams
			}
			params, _ := json.Marshal(map[string]interface{}{
				"symbol": args[0],
				"price":  price,
				"amount": amount,
			})
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "sell",
				Params:  params,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request, %s", err)
				return errRPC
			}
			var orderid int
			if err = json.Unmarshal(resp.Result, &orderid); err != nil {
				fmt.Printf("Error: invalid response, %s", err)
			}
			b, _ := json.MarshalIndent(map[string]int{"orderid": orderid}, "", "    ")
			fmt.Println(string(b))
			return nil
		},
		ArgsUsage: "[tradepair price amount]",
		Usage:     "sell places a sell order",
	}
}
