// +build c2cx_integration_test

package c2cx

import (
	"encoding/json"
	"errors"
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

const (
	// declaring these as globals so client_test can test with the same params
	orderMarket = "BTC_SKY"
	orderPrice  = decimal.NewFromFloat(0.5) // 0.5 btc/sky? I like that price!
	orderAmount = decimal.NewFromFloat(1.2)
	orderType   = "Sell"
	priceTypeId = PriceTypeLimit
)

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
	Btc float64 `json:"btc"`
	Etc float64 `json:"etc"`
	Eth float64 `json:"eth"`
	Cny float64 `json:"cny"`
	Sky float64 `json:"sky"`
}

type balanceJSON struct {
	Balance balanceJSONEntry `json:balance`
	Frozen  balanceJSONEntry `json:balance`
}

func availableSKY() (float64, error) {
	var endpoint = "getbalance"
	resp, err := requestPost(endpoint, key, secret, nil)
	if err != nil {
		return 0.0, err
	}
	if resp.Code != 200 {
		return 0.0, apiError(endpoint, resp.Message)
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
	if availSky < orderAmount {
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
