package c2cx

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/skycoin/exchange-api/exchange"
)

type response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

var (
	httpclient = &http.Client{}
	apiroot    = url.URL{
		Scheme: "https",
		Host:   "api.c2cx.com",
		Path:   "/rest/",
	}
)

func requestGet(method string, params url.Values) (*response, error) {
	reqURL := apiroot
	reqURL.Path += method
	reqURL.RawQuery = params.Encode()
	resp, err := httpclient.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	return readResponse(resp.Body)

}
func requestPost(method, key, secret string, params url.Values) (*response, error) {
	reqURL := apiroot
	reqURL.Path += method
	if params == nil {
		params = url.Values{}
	}
	params.Add("apiKey", key)
	req, _ := http.NewRequest("POST", reqURL.String(), strings.NewReader(params.Encode()+"&"+"sign="+sign(secret, params)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	return readResponse(resp.Body)
}

// getOrderbook gets all open orders by given symbol
// This method does not required API key and signing
func getOrderbook(symbol string) (*Orderbook, error) {
	var (
		params   = url.Values{}
		endpoint = "getorderbook"
		err      error
	)
	if symbol, err = normalize(symbol); err != nil {
		return nil, err
	}
	params.Add("symbol", symbol)
	resp, err := requestGet(endpoint, params)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, apiError(endpoint, resp.Message)
	}
	var result Orderbook
	return &result, json.Unmarshal(resp.Data, &result)

}

// getBalance returns user balance for all avalible currencies
// return value is a map[string]string
// all keys should be a lowercase
func getBalance(key, secret string) (userInfo Balance, err error) {
	var endpoint = "getuserinfo" // for new api it should be getbalance
	resp, err := requestPost(endpoint, key, secret, nil)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, apiError(endpoint, resp.Message)
	}
	err = json.Unmarshal(resp.Data, &userInfo)
	return
}

// Allowed priceTypeID for CreateOrder
const (
	PriceTypeLimit  = 1
	PriceTypeMarket = 2
)

// createOrder creates order with given orderType and parameters
// advanced is a advanced options for order creation
// if advanced is nil, isAdvancedOrder sets to zero, else advanced will be used as advanced options
func createOrder(key, secret string, market string, price, quantity decimal.Decimal, orderType string, priceTypeID int, advanced *AdvancedOrderParams) (int, error) {
	var (
		params = url.Values{
			"symbol":      []string{market},
			"price":       []string{price.String()},
			"quantity":    []string{quantity.String()},
			"orderType":   []string{orderType},
			"priceTypeId": []string{strconv.Itoa(priceTypeID)},
		}
		endpoint = "createorder"
	)
	if advanced != nil {
		params.Add("isAdvancedOrder", "1")
		if !advanced.StopLoss.Equal(decimal.Zero) {
			params.Add("stopLoss", advanced.StopLoss.String())
		}
		if !advanced.TakeProfit.Equal(decimal.Zero) {
			params.Add("takeProfit", advanced.TakeProfit.String())
		}
		if !advanced.TriggerPrice.Equal(decimal.Zero) {
			params.Add("triggerPrice", advanced.TriggerPrice.String())
		}
	} else {
		params.Add("isAdvancedOrder", "0")
	}
	resp, err := requestPost(endpoint, key, secret, params)
	if err != nil {
		return 0, err
	}
	if resp.Code != 200 {
		return 0, apiError(endpoint, resp.Message)
	}
	var orderid newOrder
	err = json.Unmarshal(resp.Data, &orderid)
	return orderid.OrderID, err
}

// getOrderinfo returns extended information about given order
// if orderID is -1, then GetOrderInfo returns array of all unfilled orders
func getOrderinfo(key, secret string, symbol string, orderID int, page *int) (orders []Order, err error) {
	if symbol, err = normalize(symbol); err != nil {
		return nil, err
	}
	var (
		params = url.Values{
			"orderId": []string{strconv.Itoa(orderID)},
			"symbol":  []string{symbol},
		}
		endpoint = "getorderinfo"
	)
	if page != nil {
		params.Add("page", strconv.Itoa(*page))
	}
	resp, err := requestPost(endpoint, key, secret, params)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, apiError(endpoint, resp.Message)
	}
	return orders, json.Unmarshal(resp.Data, &orders)
}

// cancelOrder cancel order with given orderID and returns error
// error == nil if cancelOrder was finished successfully
func cancelOrder(key, secret string, orderID int) (err error) {
	var (
		params = url.Values{
			"orderId": []string{strconv.Itoa(orderID)},
		}
		endpoint = "cancelorder"
	)
	resp, err := requestPost(endpoint, key, secret, params)
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return apiError(endpoint, resp.Message)
	}
	return nil
}

// statuses is a possible statusees of order
var statuses = map[string]int{
	exchange.Opened:    2,
	exchange.Partial:   3,
	exchange.Completed: 4,
	exchange.Cancelled: 5,
	exchange.Submitted: 7,
}

// getOrderByStatus get all orders with given status
// interval is time in seconds between now and start time, if interval == -1, then returns all orders
// statuses defined below
func getOrderByStatus(key, secret, symbol, status string, interval int) (orders []Order, err error) {
	if symbol, err = normalize(symbol); err != nil {
		return nil, err
	}
	var (
		params = url.Values{
			"symbol": []string{symbol},
			"status": []string{strconv.Itoa(statuses[status])},
		}
		endpoint = "getorderbystatus"
	)
	if interval > 0 {
		params.Add("interval", strconv.Itoa(interval))
	}
	resp, err := requestPost(endpoint, key, secret, params)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, apiError(endpoint, resp.Message)
	}
	return orders, json.Unmarshal(resp.Data, &orders)
}
