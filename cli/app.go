package cli

import "github.com/urfave/cli"

var rpcaddr = new(string)
var endpoint string

// App is a cli app
var App = cli.App{
	Commands: []cli.Command{
		{
			Name: "c2cx",
			Subcommands: []cli.Command{
				orderCMDs(),
				completedCMDs(),
				executedCMDs(),
				cancelCMDs(),
				buyCMD(),
				sellCMD(),
				orderbookCMD(),
				sumbitTradeCMD(),
				balanceCMD(),
			},
			Before: func(c *cli.Context) error {
				endpoint = "c2cx"
				return nil
			},
		},
		{
			Name: "cryptopia",
			Subcommands: []cli.Command{
				orderCMDs(),
				completedCMDs(),
				executedCMDs(),
				cancelCMDs(),
				buyCMD(),
				sellCMD(),
				orderbookCMD(),
				depositCMD(),
				submitWithdrawCMD(),
				transactionsCMD(),
				trackingCMDs(),
				balanceCMD(),
			},
			Before: func(c *cli.Context) error {
				endpoint = "cryptopia"
				return nil
			},
		},
	},
	Flags: []cli.Flag{
		cli.StringFlag{Name: "rpc", Destination: rpcaddr, Value: "localhost:12345"},
	},
}
