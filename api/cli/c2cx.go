package cli

import (
	"fmt"

	c2cx "github.com/skycoin/exchange-api/exchange/c2cx.com"
	"github.com/urfave/cli"
)

func sumbitTradeCMD() cli.Command {
	var name = "submittrade"
	var (
		pricetype string
		ordertype string
		takeprofit string
		stoploss string
		triggerprice string
	)
	return cli.Command{
		Name:      name,
		Usage:     "Create new order with advanced parameters",
		ArgsUsage: "<market> <price> <amount>",
		Action: func(c *cli.Context) error {
			if c.NArg() != 3 {
				return errInvalidInput
			}
			var (
				symbol        string
				price, amount decimal.Decimal
				stopLoss, takeProfit, triggerPrice decimal.Decimal
			)
			if stopLoss, err := decimal.NewFromString(stoploss); err != nil {
				return err
			}
			if takeProfit, err := decimal.NewFromString(takeprofit); err != nil {
				return err
			}
			if triggerPrice, err := decimal.NewFromString(triggerprice); err != nil {
				return err
			}
			var params = map[string]interface{}{
				"price_type_id": pricetype,
				"order_type":    ordertype,
				"symbol":        symbol,
				"price":         price,
				"amount":        amount,
				"advanced": c2cx.AdvancedOrderParams{
					StopLoss:     stopLoss,
					TakeProfit:   takeProfit,
					TriggerPrice: triggerPrice,
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
				Destination: &pricetype,
			},
			cli.StringFlag{
				Name:        "type",
				Destination: &ordertype,
			},
			cli.StringFlag{
				Name:        "takeprofit",
				Destination: &takeprofit,
			},
			cli.StringFlag{
				Name:        "stoploss",
				Destination: &stoploss,
			},
			cli.StringFlag{
				Name:        "triggerprice",
				Destination: &triggerprice,
			},
		},
	}
}
