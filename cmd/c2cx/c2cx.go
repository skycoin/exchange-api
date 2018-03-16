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
	tradePair := c2cx.DrgDash

	fmt.Println("TradePair", tradePair)

	// fmt.Println()
	// fmt.Println("GetOrderbook")
	// orderbook, err := c.GetOrderbook(tradePair)
	// exitOnError(err)
	// fmt.Printf("Orderbook:\n%+v\n", orderbook)

	fmt.Println()
	fmt.Println("GetBalanceSummary")
	balances, err := c.GetBalanceSummary()
	exitOnError(err)

	fmt.Printf("Balances:\n%+v\n", balances)

	price, err := decimal.NewFromString("0.00102")
	exitOnError(err)

	amount, err := decimal.NewFromString("2")
	exitOnError(err)

	customerID := "foo"

	fmt.Println()
	fmt.Println("LimitBuy customerID", customerID)
	orderID, err := c.LimitBuy(tradePair, price, amount, &customerID)
	exitOnError(err)

	fmt.Println("Order ID:", orderID)

	fmt.Println()
	fmt.Println("LimitBuy customerID", customerID)
	orderID2, err := c.LimitBuy(tradePair, price, amount, &customerID)
	exitOnError(err)

	fmt.Println("Order ID2:", orderID2)

	fmt.Println()
	fmt.Println("GetOrderByStatus")
	orders, err := c.GetOrderByStatus(tradePair, c2cx.StatusAll)
	exitOnError(err)

	fmt.Println("Orders found:", len(orders))
	for _, o := range orders {
		fmt.Printf("%+v\n", o)
	}

	fmt.Println()
	fmt.Println("GetOrderInfoAll")
	orders, err = c.GetOrderInfoAll(tradePair)
	exitOnError(err)

	fmt.Println("Orders found:", len(orders))
	for _, o := range orders {
		fmt.Printf("%+v\n", o)
	}

	fmt.Println()
	fmt.Println("GetOrderInfo (one)")
	order, err := c.GetOrderInfo(tradePair, orderID)
	exitOnError(err)

	fmt.Printf("Order:\n%+v\n", order)

	err = c.CancelOrder(orderID)
	exitOnError(err)

	fmt.Println("Cancelled order", orderID)

	orderIDs, err := c.CancelAll(tradePair)
	exitOnError(err)

	fmt.Println("Cancelled all orders:", orderIDs)
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
	runExamples(c)
}
