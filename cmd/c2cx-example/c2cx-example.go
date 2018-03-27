package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/skycoin/exchange-api/exchange/c2cx"
)

// nolint
func exitOnError(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

// nolint
func runExamples(c *c2cx.Client) {
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

// nolint
func highMarketSell(c *c2cx.Client) {
	tradePair := c2cx.BtcSky

	fmt.Println("TradePair", tradePair)

	fmt.Println()
	fmt.Println("GetBalanceSummary")

	balances, err := c.GetBalanceSummary()
	exitOnError(err)

	fmt.Printf("Balances:\n%+v\n", balances)

	amount, err := decimal.NewFromString("5")
	exitOnError(err)

	customerID := "foo-2"

	fmt.Println()
	fmt.Println("MarketSell customerID", customerID)
	orderID, err := c.MarketSell(tradePair, amount, &customerID)
	exitOnError(err)

	fmt.Println("Order ID:", orderID)

	fmt.Println()
	fmt.Println("GetOrderInfo (one)")
	order, err := c.GetOrderInfo(tradePair, orderID)
	exitOnError(err)

	fmt.Printf("Order:\n%+v\n", order)
}

// nolint
func highMarketBuy(c *c2cx.Client) {
	tradePair := c2cx.BtcSky

	fmt.Println("TradePair", tradePair)

	fmt.Println()
	fmt.Println("GetBalanceSummary")

	balances, err := c.GetBalanceSummary()
	exitOnError(err)

	fmt.Printf("Balances:\n%+v\n", balances)

	amount, err := decimal.NewFromString("1")
	exitOnError(err)

	customerID := "foo-3"

	fmt.Println()
	fmt.Println("MarketBuy amount customerID", amount.String(), customerID)
	orderID, err := c.MarketBuy(tradePair, amount, &customerID)
	exitOnError(err)

	fmt.Println("Order ID:", orderID)

	fmt.Println()
	fmt.Println("GetOrderInfo (one)")
	order, err := c.GetOrderInfo(tradePair, orderID)
	exitOnError(err)

	fmt.Printf("Order:\n%+v\n", order)
}

// nolint
func triggerRatelimit(c *c2cx.Client) {
	// Docs state limit of 60 requests per minute per endpoint
	// This method tries to trigger the ratelimit to determine what error is returned,
	// since the error response is undocumented

	var wg sync.WaitGroup

	type ReqError struct {
		err error
		n   int
	}

	errs := make(chan ReqError)

	for i := 0; ; i++ {
		j := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(j)
			_, err := c.GetBalanceSummary()
			if err != nil {
				errs <- ReqError{
					err: err,
					n:   j,
				}
			}
		}()
		time.Sleep(time.Millisecond * 100)
	}

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		k := 0
		defer wg2.Done()
		for err := range errs {
			fmt.Println()
			fmt.Println(err.n)
			printError(err.err)
			k++
		}
		fmt.Println("number of errors:", k)
	}()

	wg.Wait()
	close(errs)
	wg2.Wait()
}

// nolint
func lowMarketBuy(c *c2cx.Client) {
	// Make a very low market buy order to discover the error message

	// amount=0 tradepair=DrgBtc c2cx request failed: endpoint=createorder code=400 message=limit value: 13
	// amount=0.0000000001 tradepair=BtcSky c2cx request failed: endpoint=createorder code=400 message=limit value: 0.00159
	// amount=0.0000000001 tradepair=BtcSky c2cx request failed: endpoint=createorder code=400 message=limit value: 0.00158

	c.Debug = true

	tradePair := c2cx.BtcSky

	amount, err := decimal.NewFromString("0.00001")
	exitOnError(err)

	cid := fmt.Sprintf("lmb-%d", rand.Uint32())

	fmt.Println("making a market buy from", tradePair, amount, cid)

	_, err = c.MarketBuy(tradePair, amount, &cid)
	printError(err)
}

// nolint
func printError(err error) {
	if err == nil {
		fmt.Println("no error")
	} else {
		fmt.Println("ERROR:", err)
		apiErr, ok := err.(c2cx.APIError)
		if ok {
			fmt.Printf("%+v\n", apiErr)
		} else {
			fmt.Println("not an api error")
		}
	}
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

	c := c2cx.NewAPIClient(*key, *secret)


	doNothing(c)
	// runExamples(c)
	// triggerRatelimit(c)
	// lowMarketBuy(c)
	// highMarketSell(c)
	// highMarketBuy(c)
}
