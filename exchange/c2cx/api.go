package c2cx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	// The following is nolinted because it's part of c2cx's authentication scheme
	"crypto/md5" // nolint: gas

	"github.com/shopspring/decimal"
)

const (
	getOrderbookEndpoint     = "getorderbook"
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

// Error represents an error in the C2CX API wrapper
type Error interface {
	error
	APIError() bool // Is the error an API error?
}

// OtherError is returned when an error other than an API error occurs, such as a net.Error or a JSON parsing error
type OtherError struct {
	error
}

// NewOtherError creates an Error
func NewOtherError(err error) OtherError {
	return OtherError{err}
}

// APIError returns false
func (e OtherError) APIError() bool {
	return false
}

// APIError is returned when an API response has an error code
type APIError struct {
	Code     int
	Message  string
	Endpoint string
}

// NewAPIError creates an APIError
func NewAPIError(endpoint string, code int, message string) APIError {
	return APIError{
		Code:     code,
		Message:  message,
		Endpoint: endpoint,
	}
}

// APIError returns true
func (e APIError) APIError() bool {
	return true
}

func (e APIError) Error() string {
	return fmt.Sprintf("C2CX request failed: endpoint=%s code=%d message=%s", e.Endpoint, e.Code, e.Message)
}

// Client implements a wrapper around the C2CX API interface
type Client struct {
	Key    string
	Secret string
	Debug  bool
}

// CancelMultiError is returned when an error was encountered while cancelling multiple orders
type CancelMultiError struct {
	OrderIDs []OrderID
	Errors   []error
}

func (e CancelMultiError) Error() string {
	return fmt.Sprintf("these orders failed to cancel: %v", e.OrderIDs)
}

type status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type pagination struct {
	PageIndex   *int `json:"pageindex"`
	PageSize    *int `json:"pagesize"`
	RecordCount int  `json:"recordcount"`
	PageCount   int  `json:"pagecount"`
}

type getOrderbookResponse struct {
	status
	Orderbook Orderbook `json:"data"`
}

// GetOrderbook gets all open orders by given symbol
// This method does not required API key and signing
func (c *Client) GetOrderbook(symbol TradePair) (*Orderbook, error) {
	params := url.Values{}
	params.Set("symbol", string(symbol))

	data, err := c.get(getOrderbookEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp getOrderbookResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return nil, NewAPIError(getOrderbookEndpoint, resp.status.Code, resp.status.Message)
	}

	return &resp.Orderbook, nil
}

type getBalanceResponse struct {
	status
	BalanceSummary BalanceSummary `json:"data"`
}

// GetBalanceSummary returns user balance for all available currencies
func (c *Client) GetBalanceSummary() (*BalanceSummary, error) {
	data, err := c.post(getBalanceEndpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp getBalanceResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return nil, NewAPIError(getBalanceEndpoint, resp.status.Code, resp.status.Message)
	}

	return &resp.BalanceSummary, nil
}

type createOrderResponse struct {
	status
	Order struct {
		OrderID OrderID `json:"orderId"`
	} `json:"data"`
}

// CreateOrder creates order with given orderType and parameters
// advanced is a advanced options for order creation
// if advanced is nil, isAdvancedOrder sets to zero, else advanced will be used as advanced options
func (c *Client) CreateOrder(symbol TradePair, price, quantity decimal.Decimal, orderType OrderType, priceType PriceType, customerID *string, advanced *AdvancedOrderParams) (OrderID, error) {
	params := url.Values{}
	params.Set("symbol", string(symbol))
	params.Set("price", price.String())
	params.Set("quantity", quantity.String())
	params.Set("orderType", string(orderType))
	params.Set("priceTypeId", string(priceType))

	if customerID != nil {
		params.Set("cid", *customerID)
	}

	isAdvanced := false
	if advanced != nil {
		if advanced.StopLoss != nil {
			params.Set("stopLoss", advanced.StopLoss.String())
			isAdvanced = true
		}
		if advanced.TakeProfit != nil {
			params.Set("takeProfit", advanced.TakeProfit.String())
			isAdvanced = true
		}
		if advanced.TriggerPrice != nil {
			params.Set("triggerPrice", advanced.TriggerPrice.String())
			isAdvanced = true
		}
	}

	if isAdvanced {
		params.Set("isAdvancedOrder", "1")
	} else {
		params.Set("isAdvancedOrder", "0")
	}

	data, err := c.post(createOrderEndpoint, params)
	if err != nil {
		return 0, err
	}

	var resp createOrderResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return 0, NewAPIError(createOrderEndpoint, resp.status.Code, resp.status.Message)
	}

	return resp.Order.OrderID, nil
}

type getOrderInfoResponse struct {
	status
	Order Order `json:"data"`
}

// GetOrderInfo returns extended information about given order
func (c *Client) GetOrderInfo(symbol TradePair, orderID OrderID) (*Order, error) {
	params := url.Values{}
	params.Set("orderId", fmt.Sprint(orderID))
	params.Set("symbol", string(symbol))

	data, err := c.post(getOrderInfoEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp getOrderInfoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return nil, NewAPIError(getOrderInfoEndpoint, resp.status.Code, resp.status.Message)
	}

	return &resp.Order, nil
}

type getOrderInfoAllResponse struct {
	status
	Orders []Order `json:"data"`
}

// GetOrderInfoAll returns extended information about all orders
// Returns a 400 if it decides there are no orders (there may be orders but it can disagree).
func (c *Client) GetOrderInfoAll(symbol TradePair) ([]Order, error) {
	params := url.Values{}
	params.Set("orderId", allOrders)
	params.Set("symbol", string(symbol))

	data, err := c.post(getOrderInfoEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp getOrderInfoAllResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return nil, NewAPIError(getOrderInfoEndpoint, resp.status.Code, resp.status.Message)
	}

	return resp.Orders, nil
}

type cancelOrderResponse struct {
	status
	// Data is an empty dict
}

// CancelOrder cancel order with given orderID
func (c *Client) CancelOrder(orderID OrderID) error {
	params := url.Values{}
	params.Set("orderId", fmt.Sprint(orderID))

	data, err := c.post(cancelOrderEndpoint, params)
	if err != nil {
		return err
	}

	var resp cancelOrderResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return NewAPIError(cancelOrderEndpoint, resp.status.Code, resp.status.Message)
	}

	return nil
}

// NOTE: different from c2cx API docs, see note in api_notes.go
type getOrderByStatusResponse struct {
	status
	Data struct {
		pagination
		Rows []Order `json:"rows"`
	} `json:"data"`
}

// GetOrderByStatusPaged get all orders with given status for a given pagination page.
// NOTE: GetOrderByStatusPaged may returns orders with a different status than specified
func (c *Client) GetOrderByStatusPaged(symbol TradePair, status OrderStatus, page int) ([]Order, int, error) {
	params := url.Values{}
	params.Set("symbol", string(symbol))
	params.Set("status", fmt.Sprint(status))
	params.Set("pageindex", fmt.Sprint(page))
	params.Set("pagesize", fmt.Sprint(maxPageSize))

	data, err := c.post(getOrderByStatusEndpoint, params)
	if err != nil {
		return nil, 0, err
	}

	var resp getOrderByStatusResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, 0, NewOtherError(err)
	}

	if resp.status.Code != http.StatusOK {
		return nil, 0, NewAPIError(getOrderByStatusEndpoint, resp.status.Code, resp.status.Message)
	}

	return resp.Data.Rows, resp.Data.pagination.PageCount, nil
}

// GetOrderByStatus get all orders with given status. Makes multiple calls in the event of pagination.
// NOTE: GetOrderByStatus may returns orders with a different status than specified
func (c *Client) GetOrderByStatus(symbol TradePair, status OrderStatus) ([]Order, error) {
	page := 1

	pageOrders, nPages, err := c.GetOrderByStatusPaged(symbol, status, page)
	if err != nil {
		return nil, err
	}

	orders := pageOrders
	page++

	for page <= nPages {
		pageOrders, nPages, err = c.GetOrderByStatusPaged(symbol, status, page)
		if err != nil {
			return nil, err
		}

		orders = append(orders, pageOrders...)
		page++
	}

	return orders, nil
}

// CancelAll cancels all executed orders for an orderbook.
// If it encounters an error, it aborts and returns the order IDs that had been cancelled to that point.
func (c *Client) CancelAll(symbol TradePair) ([]OrderID, error) {
	orders, err := c.GetOrderInfoAll(symbol)
	if err != nil {
		return nil, err
	}

	var orderIDs []OrderID
	for _, o := range orders {
		if o.Status != StatusCancelled {
			orderIDs = append(orderIDs, o.OrderID)
		}
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
func (c *Client) LimitBuy(symbol TradePair, price, amount decimal.Decimal, customerID *string) (OrderID, error) {
	orderID, err := c.CreateOrder(symbol, price, amount, OrderTypeBuy, PriceTypeLimit, customerID, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// LimitSell place limit sell order
func (c *Client) LimitSell(symbol TradePair, price, amount decimal.Decimal, customerID *string) (OrderID, error) {
	orderID, err := c.CreateOrder(symbol, price, amount, OrderTypeSell, PriceTypeLimit, customerID, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// MarketBuy place market buy order. A market buy order will sell the entire amount of
// the trade pair's first coin in exchange for the second coin.
// e.g. for BTC_SKY, the amount is the amount of BTC you want to spend on SKY.
func (c *Client) MarketBuy(symbol TradePair, amount decimal.Decimal, customerID *string) (OrderID, error) {
	orderID, err := c.CreateOrder(symbol, amount, decimal.Zero, OrderTypeBuy, PriceTypeMarket, customerID, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// MarketSell place market sell order. A market sell order will sell the entire amount
// of the trade pair's second coin in exchange for the first coin.
// e.g. for BTC_SKY, the amount is the amount of SKY you want to sell for BTC.
func (c *Client) MarketSell(symbol TradePair, amount decimal.Decimal, customerID *string) (OrderID, error) {
	orderID, err := c.CreateOrder(symbol, amount, decimal.Zero, OrderTypeSell, PriceTypeMarket, customerID, nil)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (c *Client) get(method string, params url.Values) ([]byte, error) { // nolint: unparam
	reqURL := apiroot
	reqURL.Path += method
	reqURL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Get(reqURL.String())
	if err != nil {
		return nil, NewOtherError(err)
	}

	defer resp.Body.Close() // nolint: errcheck

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewOtherError(err)
	}

	// NOTE:
	// c2cx's API always returns 200 OK except for 500 errors
	// Instead, it places the status code inside of a JSON object per request
	// The caller must handle the status code, since the structure of the JSON
	// is different across requests
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, NewAPIError(method, resp.StatusCode, "Internal Server Error")
	}

	if c.Debug {
		fmt.Printf("GET endpoint=%s response=%s\n", reqURL.String(), string(b))
	}

	return b, nil
}

func (c *Client) post(method string, params url.Values) ([]byte, error) {
	reqURL := apiroot
	reqURL.Path += method

	if params == nil {
		params = url.Values{}
	}
	params.Set("apiKey", c.Key)

	body := encodeParamsSorted(params) + "&sign=" + sign(c.Secret, params)

	req, _ := http.NewRequest("POST", reqURL.String(), strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, NewOtherError(err)
	}

	defer resp.Body.Close() // nolint: errcheck

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewOtherError(err)
	}

	// NOTE:
	// c2cx's API always returns 200 OK except for 500 errors
	// Instead, it places the status code inside of a JSON object per request
	// The caller must handle the status code, since the structure of the JSON
	// is different across requests
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, NewAPIError(method, resp.StatusCode, "Internal Server Error")
	}

	if c.Debug {
		fmt.Printf("POST endpoint=%s body=%s response=%s\n", reqURL.String(), body, string(b))
	}

	return b, nil
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
