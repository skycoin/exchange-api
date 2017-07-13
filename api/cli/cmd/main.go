package main

import (
	"os"

	"github.com/uberfurrer/tradebot/api/cli"
)

func main() {
	cli.App.Run(os.Args)
}
