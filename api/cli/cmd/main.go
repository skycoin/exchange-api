package main

import (
	"os"

	"github.com/skycoin/exchange-api/api/cli"
)

func main() {
	cli.App.Run(os.Args)
}
