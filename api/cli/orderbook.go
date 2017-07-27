package cli

import (
	"encoding/json"
	"fmt"

	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/uberfurrer/tradebot/exchange"
	"github.com/urfave/cli"
)

func orderbookCmd() cli.Command {
	var name = "orderbook"
	return cli.Command{
		Name:      name,
		Usage:     "Gets orderbook for given market, if this market tracked by node",
		ArgsUsage: "[market]",
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) != 1 {
				return errInvalidParams
			}
			params, _ := json.Marshal(map[string]string{"symbol": args[0]})
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "orderbook",
				Params:  params,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request %s\n", err)
				return errRPC
			}
			var result exchange.MarketRecord
			if err = json.Unmarshal(resp.Result, &result); err != nil {
				fmt.Printf("Error: invalid request, %s\n", err)
				return errInvalidResponse
			}
			b, _ := json.MarshalIndent(result, "", "    ")
			fmt.Println(string(b))
			return nil
		},
	}
}
