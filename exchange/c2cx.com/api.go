package c2cx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	// the following is nolinted because it's part of c2cx's authentication scheme
	// nolint: gas
	"crypto/md5"

	"github.com/shopspring/decimal"
)

const (
	getOrderBookEndpoint     = "getorderbook"
	getBalanceEndpoint       = "getbalance"
	createOrderEndpoint      = "createorder"
	getOrderInfoEndpoint     = "getorderinfo"
	cancelOrderEndpoint      = "cancelorder"
	getOrderByStatusEndpoint = "getorderbystatus"
)

var (
	apiroot = url.URL{
		Scheme: "https",
		Host:   "api.c2cx.com",
		Path:   "/v1/",
	}
)

type response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Client implements a wrapper around the C2CX API interface
type Client struct {
	Key    string
	Secret string
}

// CancelMultiError is returned when an error was encountered while cancelling multiple orders
type CancelMultiError struct {
	OrderIDs []OrderID
	Errors   []error
}

func (e CancelMultiError) Error() string {
	return fmt.Sprintf("these orders failed to cancel: %v", e.OrderIDs)
}

// GetOrderbook gets all open orders by given symbol
// This method does not required API key and signing
func (c *Client) GetOrderbook(symbol TradePair) (*Orderbook, error) {
	params := url.Values{}
	params.Add("symbol", string(symbol))

	resp, err := c.get(getOrderBookEndpoint, params)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, apiError(getOrderBookEndpoint, resp.Message)
	}

	var result Orderbook
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBalance returns user balance for all avalible currencies
// return value is a map[string]string
// all keys should be a lowercase
func (c *Client) GetBalance() (*Balance, error) {
	resp, err := c.post(getBalanceEndpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, apiError(getBalanceEndpoint, resp.Message)
	}

	var balance Balance
	if err := json.Unmarshal(resp.Data, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// CreateOrder creates order with given orderType and parameters
// advanced is a advanced options for order creation
// if advanced is nil, isAdvancedOrder sets to zero, else advanced will be used as advanced options
func (c *Client) CreateOrder(symbol TradePair, price, quantity decimal.Decimal, orderType OrderType, priceType PriceType, advanced *AdvancedOrderParams) (OrderID, error) {
	params := url.Values{
		"symbol":      []string{string(symbol)},
		"price":       []string{price.String()},
		"quantity":    []string{quantity.String()},
		"orderType":   []string{string(orderType)},
		"priceTypeId": []string{string(priceType)},
	}

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

	resp, err := c.post(createOrderEndpoint, params)
	if err != nil {
		return 0, err
	}

	if resp.Code != http.StatusOK {
		return 0, apiError(createOrderEndpoint, resp.Message)
	}

	var orderid newOrder
	if err := json.Unmarshal(resp.Data, &orderid); err != nil {
		return 0, err
	}

	return orderid.OrderID, nil
}

// GetOrderInfo returns extended information about given order
// if orderID is -1, then GetOrderInfo returns array of all unfilled orders
func (c *Client) GetOrderInfo(symbol TradePair, orderID OrderID, page *int) ([]Order, error) {
	params := url.Values{
		"orderId": []string{fmt.Sprint(orderID)},
		"symbol":  []string{string(symbol)},
	}

	if page != nil {
		params.Add("page", strconv.Itoa(*page))
	}

	resp, err := c.post(getOrderInfoEndpoint, params)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, apiError(getOrderInfoEndpoint, resp.Message)
	}

	// if we're requesting a specific order, c2cx returns a single object, not an array
	if orderID != AllOrders {
		var order Order
		if err := json.Unmarshal(resp.Data, &order); err != nil {
			return nil, err
		}

		return []Order{order}, nil
	}

	var orders []Order
	if err := json.Unmarshal(resp.Data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// CancelOrder cancel order with given orderID
func (c *Client) CancelOrder(orderID OrderID) error {
	params := url.Values{
		"orderId": []string{fmt.Sprint(orderID)},
	}

	resp, err := c.post(cancelOrderEndpoint, params)
	if err != nil {
		return err
	}

	if resp.Code != http.StatusOK {
		return apiError(cancelOrderEndpoint, resp.Message)
	}

	return nil
}

// GetOrderByStatus get all orders with given status
func (c *Client) GetOrderByStatus(symbol TradePair, status OrderStatus) ([]Order, error) {
	// TODO -- this endpoint is paginated -- automatically follow pages to extract all orders
	params := url.Values{
		"symbol": []string{string(symbol)},
		"status": []string{strconv.Itoa(OrderStatuses[status])},
	}

	resp, err := c.post(getOrderByStatusEndpoint, params)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, apiError(getOrderByStatusEndpoint, resp.Message)
	}

	var orders []Order
	if err := json.Unmarshal(resp.Data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// CancelAll cancels all executed orders for an orderbook.
// If it encounters an error, it aborts and returns the order IDs that had been cancelled to that point.
func (c *Client) CancelAll(symbol TradePair) ([]OrderID, error) {
	orders, err := c.GetOrderInfo(symbol, AllOrders, nil)
	if err != nil {
		return nil, err
	}

	var orderIDs []OrderID
	for _, o := range orders {
		orderIDs = append(orderIDs, o.OrderID)
	}

	return c.CancelMultiple(orderIDs)
}

// CancelMultiple cancels multiple orders.  It will try to cancel all of them, not
// stopping for any individual error. If any orders failed to cancel, a CancelMultiError is returned
// along with the array of order IDs which were successfully cancelled.
func (c *Client) CancelMultiple(orderIDs []OrderID) ([]OrderID, error) {
	var cancelledOrderIDs []OrderID
	var cancelErr CancelMultiError

	for _, v := range orderIDs {
		if err := c.CancelOrder(v); err != nil {
			cancelErr.OrderIDs = append(cancelErr.OrderIDs, v)
			cancelErr.Errors = append(cancelErr.Errors, err)
			continue
		}
		cancelledOrderIDs = append(cancelledOrderIDs, v)
	}

	if len(cancelErr.OrderIDs) != 0 {
		return cancelledOrderIDs, &cancelErr
	}

	return cancelledOrderIDs, nil
}

// LimitBuy place limit buy order
func (c *Client) LimitBuy(symbol TradePair, price, amount decimal.Decimal) (OrderID, error) {
	orderID, err := c.CreateOrder(symbol, price, amount, OrderTypeBuy, PriceTypeLimit, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// LimitSell place limit sell order
func (c *Client) LimitSell(symbol TradePair, price, amount decimal.Decimal) (OrderID, error) {
	orderID, err := c.CreateOrder(symbol, price, amount, OrderTypeSell, PriceTypeLimit, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// MarketBuy place market buy order. A market buy order will spend the entire amount to buy
// the symbol's second coin
func (c *Client) MarketBuy(symbol TradePair, amount decimal.Decimal) (OrderID, error) {
	// For "market" orders, the amount to spend is placed in the "price" field
	orderID, err := c.CreateOrder(symbol, amount, decimal.Zero, OrderTypeBuy, PriceTypeMarket, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// MarketSell place market sell order. A market sell order will sell the entire amount
// of the symbol's first coin
func (c *Client) MarketSell(symbol TradePair, amount decimal.Decimal) (OrderID, error) {
	// For "market" orders, the amount to sell is placed in the "price" field
	orderID, err := c.CreateOrder(symbol, amount, decimal.Zero, OrderTypeSell, PriceTypeMarket, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (c *Client) get(method string, params url.Values) (*response, error) { // nolint: unparam
	reqURL := apiroot
	reqURL.Path += method
	reqURL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Get(reqURL.String())
	if err != nil {
		return nil, err
	}

	return readResponse(resp.Body)
}

func (c *Client) post(method string, params url.Values) (*response, error) {
	reqURL := apiroot
	reqURL.Path += method

	if params == nil {
		params = url.Values{}
	}
	params.Add("apiKey", c.Key)

	req, _ := http.NewRequest("POST", reqURL.String(), strings.NewReader(params.Encode()+"&"+"sign="+sign(c.Secret, params)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return readResponse(resp.Body)
}

func apiError(endpoint, message string) error {
	return fmt.Errorf("c2cx: %s falied, %s", endpoint, message)
}

func sign(secret string, params url.Values) string {
	var paramString = encodeParamsSorted(params)
	if len(paramString) > 0 {
		paramString += "&secretKey=" + secret
	} else {
		paramString += "secretKey=" + secret
	}

	sum := md5.Sum([]byte(paramString)) // nolint: gas
	return strings.ToUpper(fmt.Sprintf("%x", sum))
}

// returns sorted string for signing
func encodeParamsSorted(params url.Values) string {
	if params == nil {
		return ""
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	result := bytes.Buffer{}
	for i, k := range keys {
		result.WriteString(url.QueryEscape(k))
		result.WriteString("=")
		result.WriteString(url.QueryEscape(params.Get(k)))

		if i != len(keys)-1 {
			result.WriteString("&")
		}
	}

	return result.String()
}

func readResponse(r io.ReadCloser) (*response, error) {
	var tmp struct {
		Fail    []json.RawMessage `json:"fail,omitempty"`
		Success json.RawMessage   `json:"success,omitempty"`
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := json.Unmarshal(b, &tmp); err != nil {
		return nil, err
	}

	if len(tmp.Fail) != 0 {
		return nil, errors.New(string(tmp.Fail[0]))
	}

	var resp response
	if err := json.Unmarshal(tmp.Success, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
