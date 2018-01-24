package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

func balanceCMD() cli.Command {
	name := "balance"
	return cli.Command{
		Name:      name,
		Usage:     "Print balance",
		ArgsUsage: "<currency>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errInvalidInput
			}
			params := map[string]interface{}{
				"currency": c.Args().First(),
			}
			resp, err := rpcRequest("balance", params)
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}
}
