package main

import (
	"os"

	"github.com/skycoin/exchange-api/cli"
)

func main() {
	err := cli.App.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
