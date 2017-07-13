package c2cx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"time"

	"github.com/pkg/errors"
	"github.com/uberfurrer/tradebot/exchange"
	"github.com/uberfurrer/tradebot/logger"
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
		Path:   "/v1/",
	}
)

func requestGet(method string, params url.Values) (*response, error) {
	reqURL := apiroot
	reqURL.Path += method
	reqURL.RawQuery = params.Encode()
	resp, err := httpclient.Get(reqURL.String())
	if err != nil {
		logger.Error("c2cx: requestGet http error, ", err)
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
	params.Add("sign", sign(secret, params))
	req, _ := http.NewRequest("POST", reqURL.String(), bytes.NewReader([]byte(params.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpclient.Do(req)
	if err != nil {
		logger.Error("c2cx: requestPost http error, ", err)
		return nil, err
	}
	return readResponse(resp.Body)
}

// GetOrderBook gets all open orders by given symbol
// This method does not required API key and signing
func GetOrderBook(symbol string) (*Orderbook, error) {
	var params = url.Values{}
	var err error
	if symbol, err = normalize(symbol); err != nil {
		logger.Error("c2cx: invalid market,", err)
		return nil, err
	}
	params.Add("symbol", symbol)
	resp, err := requestGet("getorderbook", params)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, errors.Errorf("GetOrderBook failed: %d %s", resp.Code, resp.Message)
	}
	var result Orderbook
	return &result, json.Unmarshal(resp.Data, &result)

}

// GetBalance returns user balance for all avalible currencies
// return value is a map[string]string
// all keys should be a lowercase
func GetBalance(key, secret string) (userInfo Balance, err error) {
	resp, err := requestPost("getbalance", key, secret, nil)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, errors.Errorf("getbalance failed: %d %s", resp.Code, resp.Message)
	}

	err = json.Unmarshal(resp.Data, &userInfo)
	return

}

// Allowed priceTypeID for CreateOrder
const (
	PriceTypeLimit  = 1
	PriceTypeMarket = 2
)

// CreateOrder creates order with given orderType and parameters
func CreateOrder(key, secret, symbol, orderType string, priceTypeID int, triggerPrice, quantity, price float64,
	takeprofit, stoploss *float64, exptime *time.Time) (orderID int, err error) {
	var (
		params = url.Values{}
	)
	if symbol, err = normalize(symbol); err != nil {
		logger.Error("c2cx: invalid market", err)
		return 0, err
	}
	params.Add("symbol", symbol)
	params.Add("orderType", orderType)
	if priceTypeID != 1 && priceTypeID != 2 {
		logger.Error("priceTypeId must be 1 - limit or 2 - market")
		return 0, errors.New("Wrong priceTypeID")
	}
	params.Add("priceTypeId", strconv.Itoa(priceTypeID))
	params.Add("triggerPrice", strconv.FormatFloat(triggerPrice, 'f', -1, 64))
	params.Add("quantity", strconv.FormatFloat(quantity, 'f', -1, 64))
	params.Add("price", strconv.FormatFloat(price, 'f', -1, 64))

	if takeprofit != nil {
		params.Add("takeProfit", strconv.FormatFloat(*takeprofit, 'f', -1, 64))
	}
	if stoploss != nil {
		params.Add("stopLoss", strconv.FormatFloat(*stoploss, 'f', -1, 64))
	}
	if exptime != nil {
		params.Add("expirationDate", exptime.Format("2005-1-2 15:04:05"))
	}

	resp, err := requestPost("createorder", key, secret, params)
	if err != nil {
		return 0, err
	}
	if resp.Code != 200 {
		return 0, errors.Errorf("CreateOrder failed: %d, %s", resp.Code, resp.Message)
	}
	err = json.Unmarshal(resp.Data, &orderID)
	if err != nil {
		return 0, errors.Wrap(err, "error while parsing response")
	}
	return
}

// GetOrderInfo returns extended information about given order
// if orderID == -1, then GetOrderInfo returns array of all unfilled orders
func GetOrderInfo(key, secret, symbol string, orderID int) (orders []OrderInfo, err error) {
	var params = url.Values{}
	if symbol, err = normalize(symbol); err != nil {
		return nil, err
	}
	params.Add("symbol", symbol)
	params.Add("orderId", strconv.Itoa(orderID))
	resp, err := requestPost("getorderinfo", key, secret, params)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, errors.Errorf("GetOrderInfo failed: %d, %s", resp.Code, resp.Message)
	}
	err = json.Unmarshal(resp.Data, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// CancelOrder cancel order with given orderID and returns error
// error == nil if CancelOrder was finished successfully
func CancelOrder(key, secret string, orderID int) (err error) {
	var params = url.Values{}
	params.Add("orderId", strconv.Itoa(orderID))

	resp, err := requestPost("cancelorder", key, secret, params)
	if err != nil {
		return err
	}
	if resp.Code != 200 {
		return errors.Errorf("CancelOrder failed: %d, %s", resp.Code, resp.Message)
	}
	var result bool
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("CancelOrder failed, result is false")
	}
	return nil
}

// Statusees is a possible statusees of order
var Statusees = map[string]int{
	exchange.StatusOpened:    2,
	exchange.StatusPartial:   3,
	exchange.StatusCompleted: 4,
	exchange.StatusCancelled: 5,
}

// GetOrderByStatus get all orders with given status
// interval is time in seconds between now and start time, if interval == -1, then returns all orders
// statusees defined below
func GetOrderByStatus(key, secret, symbol, status string, interval int) (orders []OrderInfo, err error) {
	var params = url.Values{}
	if symbol, err = normalize(symbol); err != nil {
		return nil, err
	}
	params.Add("symbol", symbol)
	params.Add("interval", strconv.Itoa(interval))
	params.Add("status", strconv.Itoa(Statusees[status]))
	resp, err := requestPost("getorderbystatus", key, secret, params)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, errors.Errorf("c2cx: GetOrderByStatus failed, %d %s", resp.Code, resp.Message)
	}
	if err = json.Unmarshal(resp.Data, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
