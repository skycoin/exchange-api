package cryptopia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"errors"

	"github.com/shopspring/decimal"
	"time"
	"net"
)

const (
	// InstantOrderID is returned if an order executed instantly and was not assigned an OrderID
	InstantOrderID = -1
)

const (
	dialTimeout         = 60 * time.Second
	httpClientTimeout   = 120 * time.Second
	tlsHandshakeTimeout = 60 * time.Second
)

var (
	apiroot = url.URL{
		Scheme: "https",
		Host:   "www.cryptopia.co.nz",
		Path:   "api/",
	}

	// ErrCurrencyNotFound is returned if a currency is not found in the currencies list
	ErrCurrencyNotFound = errors.New("Currency not found")

	// ErrTradePairNotFound is returned is a trade pair is not found in the markets
	ErrTradePairNotFound = errors.New("Trade pair not found")
)

type response struct {
	Success bool            `json:"Success"`
	Message string          `json:"Error"`
	Data    json.RawMessage `json:"Data"`
}

// Client implements a wrapper around the Cryptopia API interface
type Client struct {
	Key           string
	Secret        string
	httpClient    *http.Client
	currencyCache map[string]CurrencyInfo
	marketCache   map[string]int
}

func NewAPIClient(key string, secret string) *Client {
	var netTransport = http.Transport{
		Dial: (&net.Dialer{
			Timeout: dialTimeout,
		}).Dial,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
	var client = &http.Client{
		Transport: &netTransport,
		Timeout:   httpClientTimeout,
	}
	return &Client{
		Key:        key,
		Secret:     secret,
		httpClient: client,
	}
}

//Public API functions

// GetCurrencies gets all currencies
func (c *Client) GetCurrencies() ([]CurrencyInfo, error) {
	resp, err := c.get("getcurrencies", "")
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetCurrencies failed: %s", resp.Message)
	}

	var result []CurrencyInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTradePairs gets all TradePairs on exchange
func (c *Client) GetTradePairs() ([]TradepairInfo, error) {
	resp, err := c.get("gettradepairs", "")
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetTradePairs failed: %s", resp.Message)
	}

	var result []TradepairInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMarkets return all Market info by given baseMarket
// if baseMarket is empty or "all" getMarkets return all markets
// if hours < 1 it will be omitted, default value is 24
func (c *Client) GetMarkets(baseMarket string, hours int) ([]MarketInfo, error) {
	var requestParams string

	if len(baseMarket) > 0 && strings.ToUpper(baseMarket) != "ALL" {
		if _, err := c.GetCurrencyID(baseMarket); err != nil {
			return nil, err
		}
		requestParams += normalize(baseMarket)
	}

	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}

	resp, err := c.get("getmarkets", requestParams)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetMarkets failed: %s", resp.Message)
	}

	var result []MarketInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMarket return market with given label
// if hours < 1, it will be omitted, default value is 24
func (c *Client) GetMarket(market string, hours int) (*MarketInfo, error) {
	marketID, err := c.GetMarketID(market)
	if err != nil {
		return nil, err
	}

	requestParams := strconv.Itoa(marketID)

	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}

	resp, err := c.get("getmarket", requestParams)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetMarket failed: %s, Market: %s", resp.Message, market)
	}

	var result MarketInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMarketHistory return market history with given label
// if hours < 1, it will be omitted, default value is 24
func (c *Client) GetMarketHistory(market string, hours int) ([]MarketHistory, error) {
	marketID, err := c.GetMarketID(market)
	if err != nil {
		return nil, err
	}

	requestParams := strconv.Itoa(marketID)

	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}

	resp, err := c.get("getmarkethistory", requestParams)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetMarketHistory failed: %s, Market: %s", resp.Message, market)
	}

	var result []MarketHistory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMarketOrders returns count orders from market with given label
// if count < 1, its will be omitted, default value is 100
func (c *Client) GetMarketOrders(market string, count int) (*MarketOrders, error) {
	marketID, err := c.GetMarketID(market)
	if err != nil {
		return nil, err
	}

	requestParams := strconv.Itoa(marketID)

	if count > 0 {
		requestParams += "/" + strconv.Itoa(count)
	}

	resp, err := c.get("getmarketorders", requestParams)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetMarketOrders failed: %s, Market: %s", resp.Message, market)
	}

	var result MarketOrders
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMarketOrderGroups returns count Orders to each market
// If count < 1, it will be omitted
func (c *Client) GetMarketOrderGroups(count int, markets []string) ([]MarketOrdersWithLabel, error) {
	if len(markets) == 0 {
		return nil, errors.New("markets must not be empty")
	}

	var requestParams string

	for _, v := range markets {
		marketID, err := c.GetMarketID(v)
		if err != nil {
			return nil, err
		}
		requestParams += strconv.Itoa(marketID) + "-"
	}

	requestParams += requestParams[:len(requestParams)-1]

	if count > 0 {
		requestParams += "/" + strconv.Itoa(count)
	}

	resp, err := c.get("getmarketordergroups", requestParams)
	if err != nil {
		return nil, fmt.Errorf("GetMarketOrderGroups failed, markets: %s; original error: %v", strings.Join(markets, " "), err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetMarketOrderGroups failed: %s, Market: %s", resp.Message, strings.Join(markets, " "))
	}

	var result []MarketOrdersWithLabel
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Private API functions

// GetBalance return a string representation of balance by given currency
func (c *Client) GetBalance(currency string) (decimal.Decimal, error) {
	cID, err := c.GetCurrencyID(currency)
	if err != nil {
		return decimal.Zero, fmt.Errorf("Currency %s does not found", currency)
	}
	params := make(map[string]interface{})
	params["CurrencyId"] = cID
	resp, err := c.post("getbalance", params)
	if err != nil {
		return decimal.Zero, err
	}

	if !resp.Success {
		return decimal.Zero, fmt.Errorf("GetBalance failed: %s, Currency %s Rawdata %s", resp.Message, currency, string(resp.Data))
	}

	var result balance
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return decimal.Zero, err
	}

	if v, ok := result[normalize(currency)]; ok {
		return v, nil
	}

	return decimal.Zero, errors.New("currency was not found")
}

// GetDepositAddress returns a deposit address of given currency
func (c *Client) GetDepositAddress(currency string) (*DepositAddress, error) {
	cID, err := c.GetCurrencyID(currency)
	if err != nil {
		return nil, fmt.Errorf("Currency %s does not found", currency)
	}

	params := make(map[string]interface{})
	params["CurrencyId"] = cID

	resp, err := c.post("getdepositaddress", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetDepositAddress failed: %s, Currency %s", resp.Message, currency)
	}
	var result DepositAddress
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOpenOrders return a list of opened orders by specific market or all markets
func (c *Client) GetOpenOrders(market *string, count *int) ([]Order, error) {
	params := make(map[string]interface{})

	if market != nil {
		mID, err := c.GetMarketID(*market)
		if err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	}

	if count != nil {
		params["Count"] = *count
	}

	resp, err := c.post("getopenorders", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetOpenOrders failed: %s Market %#v Count %#v", resp.Message, market, count)
	}

	var result []Order
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTradeHistory return a list of all executed orders by specific market or all markets
func (c *Client) GetTradeHistory(market *string, count *int) ([]Order, error) {
	params := make(map[string]interface{})

	if market != nil {
		mID, err := c.GetMarketID(*market)
		if err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	}

	if count != nil {
		params["Count"] = count
	}

	resp, err := c.post("gettradehistory", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetTradeHistory failed: %s Market %s Count %d", resp.Message, *market, count)
	}

	var result []Order
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTransactions returns a list of transactions by given type
// if count < 1, it will be omitted
func (c *Client) GetTransactions(txType string, count int) ([]Transaction, error) {
	if txType = strings.Title(txType); txType != TxTypeDeposit && txType != TxTypeWithdraw {
		return nil, fmt.Errorf("Icorrect trasnaction type %s; avalible types: %s %s", txType, TxTypeDeposit, TxTypeWithdraw)
	}

	params := make(map[string]interface{})
	params["Type"] = txType
	if count > 0 {
		params["Count"] = count
	}

	resp, err := c.post("gettransactions", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("GetTransactions failed: %s Type %s Count %d", resp.Message, txType, count)
	}

	var result []Transaction
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SubmitTrade submits a new trade offer
func (c *Client) SubmitTrade(market, offerType string, rate, amount decimal.Decimal) (int, error) {
	if offerType = strings.Title(offerType); offerType != OfferTypeBuy && offerType != OfferTypeSell {
		return 0, fmt.Errorf("Incorrect offer type %s; avalible types: %s %s", offerType, OfferTypeBuy, OfferTypeSell)
	}

	mID, err := c.GetMarketID(market)
	if err != nil {
		return 0, err
	}

	params := make(map[string]interface{})
	params["TradePairId"] = mID
	params["Type"] = offerType
	params["Rate"] = rate
	params["Amount"] = amount

	resp, err := c.post("submittrade", params)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("SubmitTrade failed: %s, Type %s Market %s Rate %s Amount %s", resp.Message, offerType, market, rate.String(), amount.String())
	}

	var result newOrder
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return 0, err
	}

	if result.OrderID != nil {
		return *result.OrderID, nil
	}

	return InstantOrderID, nil
}

// CancelTrade cancel trades by given orderid, market or add active
// depends of type argument
func (c *Client) CancelTrade(tradeType string, TradePair *string, orderID *int) ([]int, error) {
	params := map[string]interface{}{
		"Type": tradeType,
	}

	switch tradeType {
	case ByOrderID:
		if orderID == nil {
			return nil, errors.New("for this type orderID should be valid")
		}
		params["OrderId"] = *orderID
	case ByMarket:
		if TradePair == nil {
			return nil, errors.New("for this type TradePair should be valid")
		}
		if tradepairID, err := c.GetMarketID(*TradePair); err == nil {
			params["TradePairId"] = tradepairID
		} else {
			return nil, errors.New("invalid tradepair")
		}
	case All:
		// all ok
	default:
		return nil, errors.New("invalid cancel type")
	}

	resp, err := c.post("CancelTrade", params)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New(resp.Message)
	}

	var orders []int
	if err := json.Unmarshal(resp.Data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// SubmitTip submits a tip to Trollbox
func (c *Client) SubmitTip(currency string, activeUsers int, amount decimal.Decimal) (string, error) {
	if activeUsers < 2 || activeUsers > 100 {
		return "", errors.New("activeUsers range 2-100")
	}

	cID, err := c.GetCurrencyID(currency)
	if err != nil {
		return "", err
	}

	params := make(map[string]interface{})
	params["ActiveUsers"] = activeUsers
	params["CurrencyId"] = cID
	params["Amount"] = amount

	resp, err := c.post("submittip", params)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", fmt.Errorf("SubmitTip failed: %s", resp.Message)
	}

	var result string
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", err
	}

	return result, nil
}

// SubmitWithdraw submits a withdrawal request. If address does not exists in you AddressBook, it will fail
// paymentid will be used only for currencies, based of CryptoNote algorhitm
func (c *Client) SubmitWithdraw(currency, address, paymentid string, amount decimal.Decimal) (int, error) {
	cID, err := c.GetCurrencyID(currency)
	if err != nil {
		return 0, err
	}

	params := make(map[string]interface{})
	if v, ok := c.currencyCache[normalize(currency)]; ok {
		if v.Algorithm == "CryptoNote" {
			params["PaymentId"] = paymentid
		}
	}

	params["CurrencyId"] = cID
	params["Address"] = address
	params["Amount"] = amount

	resp, err := c.post("submitwithdraw", params)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("SubmitWithdraw failed: %s, %s %s to %s ", resp.Message, currency, amount.String(), address)
	}

	var result int
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// SubmitTransfer submit a transfer funds to another user
func (c *Client) SubmitTransfer(currency, username string, amount decimal.Decimal) (string, error) {
	cID, err := c.GetCurrencyID(currency)
	if err != nil {
		return "", err
	}

	params := make(map[string]interface{})
	params["CurrencyId"] = cID
	params["Username"] = username
	params["Amount"] = amount

	resp, err := c.post("submittransfer", params)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", fmt.Errorf("SubmitTransfer failed: %s", resp.Message)
	}

	var result string
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", err
	}

	return result, nil
}

func (c *Client) get(endpoint string, params string) (*response, error) {
	reqURL := apiroot
	reqURL.Path += endpoint
	if len(params) > 0 {
		reqURL.Path += "/" + params
	}

	resp, err := http.DefaultClient.Get(reqURL.String())
	if err != nil {
		return nil, err
	}

	return readResponse(resp.Body)
}

func (c *Client) post(endpoint string, params map[string]interface{}) (*response, error) {
	reqURL := apiroot
	reqURL.Path += endpoint
	reqData, err := encodeValues(params)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", reqURL.String(), bytes.NewReader(reqData))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", header(c.Key, c.Secret, nonce(), reqURL, reqData))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	return readResponse(resp.Body)
}

// GetCurrencyID returns the ID of a currency
func (c *Client) GetCurrencyID(currency string) (int, error) {
	if v, ok := c.currencyCache[normalize(currency)]; ok {
		return v.ID, nil
	}

	// If not found, try update first
	if err := c.updateCurrencyCache(); err != nil {
		return 0, err
	}

	if v, ok := c.currencyCache[normalize(currency)]; ok {
		return v.ID, nil
	}

	return 0, ErrCurrencyNotFound
}

func (c *Client) updateCurrencyCache() error {
	crs, err := c.GetCurrencies()
	if err != nil {
		return err
	}

	c.currencyCache = make(map[string]CurrencyInfo)
	for _, v := range crs {
		c.currencyCache[v.Symbol] = v
	}

	return nil
}

// GetMarketID returns the ID of a trade pair
func (c *Client) GetMarketID(market string) (int, error) {
	if v, ok := c.marketCache[normalize(market)]; ok {
		return v, nil
	}

	// If not found, try update first
	if err := c.updateMarketCache(); err != nil {
		return 0, err
	}

	if v, ok := c.marketCache[normalize(market)]; ok {
		return v, nil
	}

	return 0, ErrTradePairNotFound
}

func (c *Client) updateMarketCache() error {
	mrkts, err := c.GetTradePairs()
	if err != nil {
		return err
	}

	c.marketCache = make(map[string]int)
	for _, v := range mrkts {
		c.marketCache[v.Label] = v.ID
	}

	return nil
}

// CancelAll cancels all executed orders on account
func (c *Client) CancelAll() ([]int, error) {
	orderIDs, err := c.CancelTrade(All, nil, nil)
	if err != nil {
		return nil, err
	}
	return orderIDs, nil
}

// CancelMarket cancel all orders opened in given market
func (c *Client) CancelMarket(symbol string) ([]int, error) {
	orderIDs, err := c.CancelTrade(ByMarket, &symbol, nil)
	if err != nil {
		return nil, err
	}
	return orderIDs, nil
}

// Buy places buy order
func (c *Client) Buy(symbol string, rate, amount decimal.Decimal) (int, error) {
	orderID, err := c.SubmitTrade(symbol, Buy, rate, amount)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// Sell places sell order
func (c *Client) Sell(symbol string, rate, amount decimal.Decimal) (int, error) {
	orderID, err := c.SubmitTrade(symbol, Sell, rate, amount)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}
