package cli

import (
	"fmt"

	c2cx "github.com/skycoin/exchange-api/exchange/c2cx.com"
	"github.com/urfave/cli"
)

func sumbitTradeCMD() cli.Command {
	name := "submittrade"
	pricetype    := new(string)
	ordertype    := new(string)
	takeprofit   := new(float64)
	stoploss     := new(float64)
	triggerprice := new(float64)
	return cli.Command{
		Name:      name,
		Usage:     "Create new order with advanced parameters",
		ArgsUsage: "<market> <price> <amount>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 3 {
				return errInvalidInput
			}
			symbol := ""
			price  := 0.0
			amount := 0.0
			params := map[string]interface{}{
				"price_type_id": pricetype,
				"order_type":    ordertype,
				"symbol":        symbol,
				"price":         price,
				"amount":        amount,
				"advanced": c2cx.AdvancedOrderParams{
					StopLoss:     *stoploss,
					TakeProfit:   *takeprofit,
					TriggerPrice: *triggerprice,
				},
			}
			resp, err := rpcRequest("submit_trade", params)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return nil
			}
			fmt.Printf("Order %s created", resp)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "pricetype",
				Destination: pricetype,
			},
			cli.StringFlag{
				Name:        "type",
				Destination: ordertype,
			},
			cli.Float64Flag{
				Name:        "takeprofit",
				Destination: takeprofit,
			},
			cli.Float64Flag{
				Name:        "stoploss",
				Destination: stoploss,
			},
			cli.Float64Flag{
				Name:        "triggerprice",
				Destination: triggerprice,
			},
		},
	}
}
