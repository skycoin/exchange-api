package cli

import "github.com/urfave/cli"

// App is a cli app
var App = cli.App{
	Commands: []cli.Command{
		cli.Command{
			Name: "c2cx",
			Subcommands: []cli.Command{
				balanceCmd(),
				buyCmd(),
				sellCmd(),
				cancelCmd(),
				orderbookCmd(),
				orderCmd(),
				completedCmd(),
				executedCmd(),
			},
			Before: func(c *cli.Context) error {
				endpoint = "c2cx"
				return nil
			},
		},
		cli.Command{
			Name: "cryptopia",
			Subcommands: []cli.Command{
				balanceCmd(),
				buyCmd(),
				sellCmd(),
				cancelCmd(),
				orderbookCmd(),
				orderCmd(),
				completedCmd(),
				executedCmd(),
			},
			Before: func(c *cli.Context) error {
				endpoint = "cryptopia"
				return nil
			},
		},
	},
}
