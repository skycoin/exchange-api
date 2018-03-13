// +build c2cx_integration_test

package c2cx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

var key, secret = func() (key string, secret string) {
	var found bool
	if key, found = os.LookupEnv("C2CX_TEST_KEY"); found {
		if secret, found = os.LookupEnv("C2CX_TEST_SECRET"); found {
			return
		}
		panic("C2CX secret not provided")
	}
	panic("C2CX key not provided")
}()

var (
	// declaring these as globals so client_test can test with the same params
	orderMarket = "BTC_SKY"
	orderPrice  = decimal.NewFromFloat(0.5) // 0.5 btc/sky? I like that price!
	orderAmount = decimal.NewFromFloat(1.2)
	orderType   = "Sell"
	priceTypeId = PriceTypeLimit
)

func ExampleMarketBuy(t *testing.T) {
	c := Client{
		Key:    "your-key-here",
		Secret: "your-secret-here",
	}

	amount := decimal.NewFromString("2.12345")
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

func TestGetUserInfo(t *testing.T) {
	b, err := GetBalance(key, secret)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 5 {
		t.Fatal("invalid balance response")
	}
}

type balanceJSONEntry struct {
	Btc decimal.Decimal `json:"btc"`
	Etc decimal.Decimal `json:"etc"`
	Eth decimal.Decimal `json:"eth"`
	Cny decimal.Decimal `json:"cny"`
	Sky decimal.Decimal `json:"sky"`
}

type balanceJSON struct {
	Balance balanceJSONEntry `json:balance`
	Frozen  balanceJSONEntry `json:balance`
}

func availableSKY() (decimal.Decimal, error) {
	var zero = decimal.NewFromFloat(0.0)
	var endpoint = "getbalance"
	resp, err := requestPost(endpoint, key, secret, nil)
	if err != nil {
		return zero, err
	}
	if resp.Code != 200 {
		return zero, apiError(endpoint, resp.Message)
	}
	var balance balanceJSON
	err = json.Unmarshal(resp.Data, &balance)
	return balance.Balance.Sky, err
}

func TestOrderManipulation(t *testing.T) {
	// confirming we have enough SKY to experiment with
	availSky, err := availableSKY()
	if err != nil {
		t.Fatal(err)
	}
	if availSky.LessThan(orderAmount) {
		t.Fatal(errors.New("Test wallet doesn't have enough SKY"))
	}

	// placing test order
	orderId, err := CreateOrder(key, secret, orderMarket, orderPrice, orderAmount, orderType, priceTypeId, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetOrderInfo", func(t *testing.T) {
		orders, err := GetOrderInfo(key, secret, orderMarket, -1, nil)
		if err != nil {
			t.Fatal(err)
		}

		var order *Order
		for _, v := range orders {
			if orderId == v.OrderID {
				order = &v
			}
		}

		if order == nil {
			t.Errorf("order %d not found", orderId)
		}
	})

	// cleaning up test order
	err = CancelOrder(key, secret, orderId)
	if err != nil {
		t.Fatal(err)
	}
}
