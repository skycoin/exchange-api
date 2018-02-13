package exchange

import (
	"log"
	"time"

	"github.com/shopspring/decimal"

	"github.com/skycoin/exchange-api/db"
	c2cx "github.com/skycoin/exchange-api/exchange/c2cx.com"
)

// This example shows how to initialize a C2CX exchange client and then use it to place orders, check order status and cancel orders.
func Example_c2cx() {
	// first we need to create an orderbook tracker
	orderbook, err := db.NewOrderbookTracker(db.MemoryDatabase, "", "")
	if err != nil {
		log.Fatal("Failed to create orderbook: ", err)
	}

	// now we create the client itself
	client := &c2cx.Client{
		Key:        "ABABABAB-ABAB-ABAB-ABAB-ABABABABABAB",
		Secret:     "CDCDCDCD-CDCD-CDCD-CDCD-CDCDCDCDCDCD",
		Orders:     exchange.NewTracker(),
		Orderbooks: orderbook,

		OrderbookRefreshInterval: time.Second * 5,
		OrdersRefreshInterval:    time.Second * 5,
	}

	// exchange clients must periodically update their orderbooks and other data, so we'll run that in a separate goroutine
	go client.Update()

	// now for the main event, manipulating orders

	// let's place a buy order for SKY/BTC
	orderId, err := client.Buy(
		"SKY_BTC",
		decimal.NewFromFloat(0.00150), // price
		decimal.NewFromFloat(20.0))    // quantity

	if err != nil {
		log.Fatal("Failed to place order: ", err)
	}

	// now let's check on the order status
	status, err := client.OrderStatus(orderId)

	if err != nil {
		log.Fatal("Failed to check order status: ", err)
	}

	switch status {
	case exchange.Submitted:
		log.Println("Order was submitted")
	case exchange.Opened:
		log.Println("Order was opened")
	case exchange.Partial:
		log.Println("Order was partially executed")
	case exchange.Completed:
		log.Println("Order was completed")
	case exchange.Canceled:
		log.Println("Order was submitted")
	default:
		log.Fatal("Unrecognized order status: ", status)
	}

	// and finally, let's cancel our order
	orderInfo, err := client.Cancel(orderId)

	if err != nil {
		log.Fatal("Failed to cancel order: ", err)
	}

	log.Println("Successfully canceled order", orderInfo)
}
