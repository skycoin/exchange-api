package c2cx

import (
	"log"
	"time"

	"github.com/shopspring/decimal"

	exchange "github.com/skycoin/exchange-api/exchange"
)

// ExampleClient shows how to initialize a C2CX exchange client and then use it to place orders, check order status and cancel orders.
func ExampleClient() {
	client := &Client{
		Key:        "ABABABAB-ABAB-ABAB-ABAB-ABABABABABAB",
		Secret:     "CDCDCDCD-CDCD-CDCD-CDCD-CDCDCDCDCDCD",
		Orders:     exchange.NewTracker(),

		OrdersRefreshInterval:    time.Second * 5,
	}

	// exchange clients must periodically update their orderbooks and other data, so we'll run that in a separate goroutine
	go client.Update()

	// now for the main event, manipulating orders

	// let's place a buy order for SKY/BTC
	orderID, err := client.Buy(
		"SKY_BTC",
		decimal.NewFromFloat(0.00150), // price
		decimal.NewFromFloat(20.0))    // quantity

	if err != nil {
		log.Fatal("Failed to place order: ", err)
	}

	// now let's check on the order status
	status, err := client.OrderStatus(orderID)

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
	case exchange.Cancelled:
		log.Println("Order was submitted")
	default:
		log.Fatal("Unrecognized order status: ", status)
	}

	// and finally, let's cancel our order
	orderInfo, err := client.Cancel(orderID)

	if err != nil {
		log.Fatal("Failed to cancel order: ", err)
	}

	log.Println("Successfully canceled order", orderInfo)
}
