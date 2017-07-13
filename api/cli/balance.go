package cli

import (
	"encoding/json"
	"fmt"

	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/urfave/cli"
)

// endpoint balance
// params: {"currency":"btc"}
func balanceCmd() cli.Command {
	var name = "balance"
	return cli.Command{
		Name:      name,
		Usage:     "Gets balance of given currency",
		ArgsUsage: "[currency]",
		Action: func(c *cli.Context) error {
			var args = c.Args()
			if len(args) != 1 {
				return errInvalidParams
			}
			params, _ := json.Marshal(map[string]string{"currency": args[0]})
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "balance",
				Params:  params,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request %s\n", err)
				return errRPC
			}
			var result string
			if err = json.Unmarshal(resp.Result, &result); err != nil {
				fmt.Printf("Error: invalid response, %s\n", err)
				return errInvalidResponse
			}
			fmt.Println(result)
			return nil
		},
	}
}
