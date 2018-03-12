package c2cx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gz-c/tradebot/exchange"
	"github.com/shopspring/decimal"
)

const (
	// PriceTypeLimit a limit order
	PriceTypeLimit = "limit"
	// PriceTypeMarket a market order
	PriceTypeMarket = "market"

	// OrderTypeBuy a buy order
	OrderTypeBuy = "buy"
	// OrderTypeSell a sell order
	OrderTypeSell = "sell"

	getOrderBookEndpoint     = "getorderbook"
	getBalanceEndpoint       = "getbalance"
	createOrderEndpoint      = "createorder"
	getOrderInfoEndpoint     = "getorderinfo"
	cancelOrderEndpoint      = "cancelorder"
	getOrderByStatusEndpoint = "getorderbystatus"

	// AllStatuses is used to include all statuses when status is required
	AllStatuses = "all"
)

var (
	apiroot = url.URL{
		Scheme: "https",
		Host:   "api.c2cx.com",
		Path:   "/v1/",
	}

	// statuses is a possible statuses of order
	statuses = map[string]int{
		AllStatuses:        0,
		exchange.Opened:    2,
		exchange.Partial:   3,
		exchange.Completed: 4,
		exchange.Cancelled: 5,
		exchange.Submitted: 7,
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

// GetOrderbook gets all open orders by given symbol
// This method does not required API key and signing
func (c *Client) GetOrderbook(symbol string) (*Orderbook, error) {
	symbol, err := normalize(symbol)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("symbol", symbol)

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
	if err := json.Unmarshal(resp.Data, &userInfo); err != nil {
		return nil, err
	}

	return &balance, nil
}

// CreateOrder creates order with given orderType and parameters
// advanced is a advanced options for order creation
// if advanced is nil, isAdvancedOrder sets to zero, else advanced will be used as advanced options
func (c *Client) CreateOrder(symbol string, price, quantity decimal.Decimal, orderType, priceTypeID string, advanced *AdvancedOrderParams) (int, error) {
	symbol, err := normalize(symbol)
	if err != nil {
		return 0, err
	}

	params := url.Values{
		"symbol":      []string{symbol},
		"price":       []string{price.String()},
		"quantity":    []string{quantity.String()},
		"orderType":   []string{orderType},
		"priceTypeId": []string{priceTypeID},
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

// GetOrderinfo returns extended information about given order
// if orderID is -1, then GetOrderInfo returns array of all unfilled orders
func (c *Client) GetOrderinfo(symbol string, orderID int, page *int) ([]Order, error) {
	symbol, err := normalize(symbol)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"orderId": []string{strconv.Itoa(orderID)},
		"symbol":  []string{symbol},
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
	if orderID != -1 {
		orders = make([]Order, 1)
		return orders, json.Unmarshal(resp.Data, &(orders[0]))
	}

	var orders []Order
	if err := json.Unmarshal(resp.Data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// CancelOrder cancel order with given orderID
func (c *Client) CancelOrder(orderID int) error {
	params := url.Values{
		"orderId": []string{strconv.Itoa(orderID)},
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
func (c *Client) GetOrderByStatus(symbol, status string) ([]Order, error) {
	// TODO -- this endpoint is paginated -- automatically follow pages to extract all orders
	symbol, err := normalize(symbol)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"symbol": []string{symbol},
		"status": []string{strconv.Itoa(statuses[status])},
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

// CancelMultiError is returned when an error was encountered while cancelling multiple orders
type CancelMultiError struct {
	OrderIDs []int
	Errors   []error
}

func (e CancelMultiError) Error() string {
	return fmt.Sprintf("these orders failed to cancel: %v", e.OrderIDs)
}

// CancelAll cancels all executed orders for an orderbook.
// If it encounters an error, it aborts and returns the orders that had been cancelled to that point.
func (c *Client) CancelAll(symbol string) ([]exchange.Order, error) {
	symbol, err := normalize(symbol)
	if err != nil {
		return nil, err
	}

	orders, err := c.GetOrderByStatus(symbol, AllStatuses)
	if err != nil {
		return nil, err
	}

	var orderIDs []int
	for _, o := range orders {
		orderIDs = append(orderIDs, o.OrderID)
	}

	return c.CancelMultiple(orderIDs)
}

// CancelMultiple cancels multiple orders.  It will try to cancel all of them, not
// stopping for any individual error. If any orders failed to cancel, a CancelMultiError is returned
// along with the array of orders which successfully cancelled.
func (c *Client) CancelMultiple(orderIDs []int) ([]exchange.Order, error) {
	var orders []exchange.Order
	var cancelErr CancelMultiError

	for _, v := range orderIDs {
		order, err := c.Cancel(v)
		if err != nil {
			cancelErr.OrderIDs = append(cancelErr.OrderIDs, v)
			cancelErr.Errors = append(cancelErr.Errors, err)
			continue
		}
		orders = append(orders, order)
	}

	if len(cancelErr.OrderIDs) != 0 {
		return orders, &cancelErr
	}

	return orders, nil
}

// LimitBuy place limit buy order
func (c *Client) LimitBuy(symbol string, price, amount decimal.Decimal) (int, error) {
	orderID, err := c.CreateOrder(symbol, price, amount, OrderTypeBuy, PriceTypeLimit, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// LimitSell place limit sell order
func (c *Client) LimitSell(symbol string, price, amount decimal.Decimal) (int, error) {
	orderID, err := c.CreateOrder(symbol, price, amount, OrderTypeSell, PriceTypeLimit, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// MarketBuy place market buy order. A market buy order will spend the entire amount to buy
// the symbol's second coin
func (c *Client) MarketBuy(symbol string, amount decimal.Decimal) (int, error) {
	// For "market" orders, the amount to spend is placed in the "price" field
	orderID, err := c.CreateOrder(symbol, amount, decimal.Zero, OrderTypeBuy, PriceTypeMarket, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// MarketSell place market sell order. A market sell order will sell the entire amount
// of the symbol's first coin
func (c *Client) MarketSell(symbol string, amount decimal.Decimal) (int, error) {
	// For "market" orders, the amount to sell is placed in the "price" field
	orderID, err := c.CreateOrder(symbol, amount, decimal.Zero, OrderTypeSell, PriceTypeMarket, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func apiError(endpoint, message string) error {
	return fmt.Errorf("c2cx: %s falied, %s", endpoint, message)
}
