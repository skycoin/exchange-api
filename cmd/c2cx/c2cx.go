package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shopspring/decimal"
	"github.com/skycoin/exchange-api/exchange/c2cx"
)

func exitOnError(err error) { // nolint: megacheck
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func runExamples(c *c2cx.Client) { // nolint: deadcode,megacheck
	fmt.Println("GetOrderbook")
	orderbook, err := c.GetOrderbook(c2cx.BtcSky)
	exitOnError(err)
	fmt.Printf("Orderbook:\n%+v\n", orderbook)

	fmt.Println()
	fmt.Println("GetBalanceSummary")
	balances, err := c.GetBalanceSummary()
	exitOnError(err)

	fmt.Printf("Balances:\n%+v\n", balances)

	price, err := decimal.NewFromString("0.00102")
	exitOnError(err)

	amount, err := decimal.NewFromString("2")
	exitOnError(err)

	fmt.Println()
	fmt.Println("LimitBuy")
	orderID, err := c.LimitBuy(c2cx.BtcSky, price, amount)
	exitOnError(err)

	fmt.Println("Order ID:", orderID)

	fmt.Println()
	fmt.Println("GetOrderByStatus")
	orders, err := c.GetOrderByStatus(c2cx.BtcSky, c2cx.StatusAll)
	exitOnError(err)

	fmt.Println("Orders found:", len(orders))
	for _, o := range orders {
		fmt.Printf("%+v\n", o)
	}

	fmt.Println()
	fmt.Println("GetOrderInfoAll")
	orders, err = c.GetOrderInfoAll(c2cx.BtcSky)
	exitOnError(err)

	fmt.Println("Orders found:", len(orders))
	for _, o := range orders {
		fmt.Printf("%+v\n", o)
	}

	fmt.Println()
	fmt.Println("GetOrderInfo (one)")
	order, err := c.GetOrderInfo(c2cx.BtcSky, orderID)
	exitOnError(err)

	fmt.Printf("Order:\n%+v\n", order)

	err = c.CancelOrder(orderID)
	exitOnError(err)

	fmt.Println("Cancelled order", orderID)

	orderIDs, err := c.CancelAll(c2cx.BtcSky)
	exitOnError(err)

	fmt.Println("Cancelled all BTC_SKY orders:", orderIDs)
}

func doNothing(c *c2cx.Client) {}

func main() {
	key := flag.String("key", "", "API key")
	secret := flag.String("secret", "", "API secret key")
	flag.Parse()

	if *key == "" {
		fmt.Println("-key is required")
		os.Exit(1)
	}

	if *secret == "" {
		fmt.Println("-secret is required")
		os.Exit(1)
	}

	c := &c2cx.Client{
		Key:    *key,
		Secret: *secret,
	}

	doNothing(c)
	// runExamples(c)
}
