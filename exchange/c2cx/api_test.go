package c2cx

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
)

func ExampleMarketBuy() { // nolint: vet
	c := &Client{
		Key:    "your-key-here",
		Secret: "your-secret-here",
	}

	amount, err := decimal.NewFromString("2.12345")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	orderID, err := c.MarketBuy(BtcSky, amount)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("Placed market buy order:", orderID)
}

func Test_sign(t *testing.T) {
	var params = url.Values{}
	params.Add("apiKey", "C821DB84-6FBD-11E4-A9E3-C86000D26D7C")
	want := "BC0DE7EBA50C730BDFC575FE2CD54082"
	expected := sign("12D857DE-7A92-F555-10AC-7566A0D84D1B", params)
	if want != expected {
		t.Fatalf("Incorrect sign!\nwant %s, expected %s", want, expected)
	}
}
