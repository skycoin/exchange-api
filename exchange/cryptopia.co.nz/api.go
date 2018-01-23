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
)

type response struct {
	Success bool            `json:"Success"`
	Message string          `json:"Error"`
	Data    json.RawMessage `json:"Data"`
}

var (
	httpclient = http.Client{}
	apiroot    = url.URL{
		Scheme: "https",
		Host:   "www.cryptopia.co.nz",
		Path:   "api/",
	}
)

func requestGet(endpoint string, params string) (*response, error) {
	reqURL := apiroot
	reqURL.Path += endpoint
	if len(params) > 0 {
		reqURL.Path += "/" + params
	}
	resp, err := httpclient.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	return readResponse(resp.Body)
}
func requestPost(endpoint, key, secret, nonce string, params map[string]interface{}) (*response, error) {
	reqURL := apiroot
	reqURL.Path += endpoint
	reqData := encodeValues(params)
	req, _ := http.NewRequest("POST", reqURL.String(), bytes.NewReader(reqData))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", header(key, secret, nonce, reqURL, reqData))
	resp, err := httpclient.Do(req)

	if err != nil {
		return nil, err
	}
	return readResponse(resp.Body)
}

//Public API functions

// getCurrencies gets all currencies
func getCurrencies() ([]CurrencyInfo, error) {
	resp, err := requestGet("getcurrencies", "")
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetCurrencies failed: %s",
			resp.Message)
	}
	var result []CurrencyInfo
	return result, json.Unmarshal(resp.Data, &result)
}

// getTradePairs gets all TradePairs on exchange
func getTradePairs() ([]TradepairInfo, error) {
	resp, err := requestGet("gettradepairs", "")
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetTradePairs failed: %s",
			resp.Message)
	}
	var result []TradepairInfo
	return result, json.Unmarshal(resp.Data, &result)
}

// getMarkets return all Market info by given baseMarket
// if baseMarket is empty or "all" getMarkets return all markets
// if hours < 1 it will be omitted, default value is 24
func getMarkets(baseMarket string, hours int) ([]MarketInfo, error) {
	var (
		requestParams string
		err           error
	)
	if len(baseMarket) > 0 && strings.ToUpper(baseMarket) != "ALL" {
		if _, err = getCurrencyID(baseMarket); err != nil {
			return nil, err
		}
		requestParams += normalize(baseMarket)
	}
	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}
	resp, err := requestGet("getmarkets", requestParams)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetMarkets failed: %s",
			resp.Message)

	}
	var result []MarketInfo
	return result, json.Unmarshal(resp.Data, &result)
}

// getMarket return market with given label
// if hours < 1, it will be omitted, default value is 24
func getMarket(market string, hours int) (MarketInfo, error) {
	var (
		requestParams string
		marketID      int
		err           error
		result        MarketInfo
	)

	if marketID, err = getMarketID(market); err != nil {
		return result, err
	}
	requestParams += strconv.Itoa(marketID)
	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}
	resp, err := requestGet("getmarket", requestParams)
	if err != nil {
		return result, err
	}
	if !resp.Success {
		return result, fmt.Errorf("GetMarket failed: %s, Market: %s",
			resp.Message, market)
	}
	return result, json.Unmarshal(resp.Data, &result)
}

// getMarketHistory return market history with given label
// if hours < 1, it will be omitted, default value is 24
func getMarketHistory(market string, hours int) ([]MarketHistory, error) {
	var (
		requestParams string
		err           error
		marketID      int
	)
	if marketID, err = getMarketID(market); err != nil {
		return nil, err
	}
	requestParams += strconv.Itoa(marketID)
	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}
	resp, err := requestGet("getmarkethistory", requestParams)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetMarketHistory failed: %s, Market: %s",
			resp.Message, market)
	}
	var result []MarketHistory
	return result, json.Unmarshal(resp.Data, &result)
}

// getMarketOrders returns count orders from market with given label
// if count < 1, its will be omitted, default value is 100
func getMarketOrders(market string, count int) (MarketOrders, error) {
	var (
		requestParams string
		err           error
		marketID      int
		result        MarketOrders
	)
	if marketID, err = getMarketID(market); err != nil {
		return result, err
	}
	requestParams += strconv.Itoa(marketID)
	if count > 0 {
		requestParams += "/" + strconv.Itoa(count)
	}
	resp, err := requestGet("getmarketorders", requestParams)
	if err != nil {
		return result, err
	}
	if !resp.Success {
		return result, fmt.Errorf("GetMarketOrders failed: %s, Market: %s",
			resp.Message, market)
	}

	return result, json.Unmarshal(resp.Data, &result)
}

// getMarketOrderGroups returns count Orders to each market
// If count < 1, it will be omitted
func getMarketOrderGroups(count int, markets ...string) ([]MarketOrdersWithLabel, error) {
	var (
		requestParams string
		err           error
		marketID      int
	)
	if len(markets) == 0 {
		return nil, errNoOrders
	}
	for _, v := range markets {
		if marketID, err = getMarketID(v); err != nil {
			return nil, err
		}
		requestParams += strconv.Itoa(marketID) + "-"
	}
	requestParams = requestParams[:len(requestParams)-1]
	if count > 0 {
		requestParams += "/" + strconv.Itoa(count)
	}
	resp, err := requestGet("getmarketordergroups", requestParams)
	if err != nil {
		return nil, fmt.Errorf("GetMarketOrderGroups failed, markets: %s", strings.Join(markets, " "))
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetMarketOrderGroups failed: %s, Market: %s",
			resp.Message, strings.Join(markets, " "))
	}
	var result []MarketOrdersWithLabel
	return result, json.Unmarshal(resp.Data, &result)
}

var errNoOrders = errors.New("no orders for updating")

// Private API functions

//getBalance return a string representation of balance by given currency
func getBalance(key, secret, nonce, currency string) (string, error) {
	resp, err := requestPost("getbalance", key, secret, nonce, nil)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("GetBalance failed: %s, Currency %s Rawdata %s",
			resp.Message, currency, string(resp.Data))
	}
	var result balance
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return "", err
	}
	if v, ok := result[normalize(currency)]; ok {
		return v, nil
	}
	return "", errors.New("currency does not found")
}

// getDepositAddress returns a deposit address of given currency
func getDepositAddress(key, secret, nonce, currency string) (DepositAddress, error) {
	var result DepositAddress
	var params = make(map[string]interface{})
	cID, err := getCurrencyID(currency)
	if err != nil {
		return result, fmt.Errorf("Currency %s does not found", currency)
	}
	params["CurrencyId"] = cID
	resp, err := requestPost("getdepositaddress", key, secret, nonce, params)
	if err != nil {
		return result, err
	}
	if !resp.Success {
		return result, fmt.Errorf("GetDepositAddress failed: %s, Currency %s",
			resp.Message, currency)
	}

	return result, json.Unmarshal(resp.Data, &result)
}

// getOpenOrders return a list of opened orders by specific market or all markets
func getOpenOrders(key, secret, nonce string, market *string, count *int) ([]Order, error) {
	var (
		params = make(map[string]interface{})
		mID    int
		err    error
	)
	if market != nil {
		if mID, err = getMarketID(*market); err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	}
	if count != nil {
		params["Count"] = *count
	}
	resp, err := requestPost("getopenorders", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetOpenOrders failed: %s Market %#v Count %#v",
			resp.Message, market, count)
	}
	var result []Order
	return result, json.Unmarshal(resp.Data, &result)
}

// getTradeHistory return a list of all executed orders by specific market or all markets
func getTradeHistory(key, secret, nonce string, market *string, count *int) ([]Order, error) {
	var (
		params = make(map[string]interface{})
		err    error
		mID    int
	)
	if market != nil {
		if mID, err = getMarketID(*market); err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	}
	if count != nil {
		params["Count"] = count
	}
	resp, err := requestPost("gettradehistory", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetTradeHistory failed: %s Market %s Count %d",
			resp.Message, *market, count)
	}
	var result []Order
	return result, json.Unmarshal(resp.Data, &result)
}

// getTransactions returns a list of transactions by given type
// if count < 1, it will be omitted
func getTransactions(key, secret, nonce, Type string, count int) ([]Transaction, error) {
	var params = make(map[string]interface{})
	if Type = strings.Title(Type); Type != TxTypeDeposit && Type != TxTypeWithdraw {
		return nil, fmt.Errorf("Icorrect trasnaction type %s; avalible types: %s %s",
			Type, TxTypeDeposit, TxTypeWithdraw)
	}
	params["Type"] = Type
	if count > 0 {
		params["Count"] = count
	}
	resp, err := requestPost("gettransactions", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("GetTransactions failed: %s Type %s Count %d",
			resp.Message, Type, count)
	}
	var result []Transaction
	return result, json.Unmarshal(resp.Data, &result)

}

// Possible offer types
const (
	OfTypeBuy  = "Buy"
	OfTypeSell = "Sell"
)

// submitTrade submits a new trade offer
func submitTrade(key, secret, nonce, market, Type string, rate, amount float64) (int, error) {
	var (
		params = make(map[string]interface{})
		err    error
		mID    int
	)
	if Type = strings.Title(Type); Type != OfTypeBuy && Type != OfTypeSell {
		return 0, fmt.Errorf("Incorrect offer type %s; avalible types: %s %s",
			Type, OfTypeBuy, OfTypeSell)
	}
	if mID, err = getMarketID(market); err != nil {
		return 0, err
	}
	params["TradePairId"] = mID
	params["Type"] = Type
	params["Rate"] = rate
	params["Amount"] = amount
	resp, err := requestPost("submittrade", key, secret, nonce, params)
	if err != nil {
		return 0, err
	}
	if !resp.Success {
		return 0, fmt.Errorf("SubmitTrade failed: %s, Type %s Market %s Rate %f Amount %f",
			resp.Message, Type, market, rate, amount)
	}
	var result newOrder
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return 0, err
	}
	if result.OrderID != nil {
		return *result.OrderID, nil
	}
	return 0, ErrInstant
}

// ErrInstant willbe returned if order instantly executed and hasn't OrderID
var ErrInstant = errors.New("Order instantly executed")

// CancelTrade cancel trades by given orderid, market or add active
// depends of type argument
func cancelTrade(key, secret, nonce string, Type string, TradePair *string, orderID *int) ([]int, error) {
	var params = map[string]interface{}{
		"Type": Type,
	}
	switch Type {
	case ByOrderID:
		if orderID == nil {
			return nil, errors.New("for this type orderID should be valid")
		}
		params["OrderId"] = *orderID
	case ByMarket:
		if TradePair == nil {
			return nil, errors.New("for this type TradePair should be valid")
		}
		if tradepairID, err := getMarketID(*TradePair); err == nil {
			params["TradePairId"] = tradepairID
		} else {
			return nil, errors.New("invalid tradepair")
		}
	case All:
	// all ok
	default:
		return nil, errors.New("invalid cancel type")
	}
	resp, err := requestPost("CancelTrade", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}
	var orders []int
	return orders, json.Unmarshal(resp.Data, &orders)
}

// SubmitTip submits a tip to Trollbox
func submitTip(key, secret, nonce, currency string, activeUsers int, amount float64) (string, error) {
	var (
		params = make(map[string]interface{})
		cID    int
		err    error
	)
	if activeUsers < 2 || activeUsers > 100 {
		return "", errors.New("activeUsers range 2-100")
	}
	if cID, err = getCurrencyID(currency); err != nil {
		return "", err
	}
	params["ActiveUsers"] = activeUsers
	params["CurrencyId"] = cID
	params["Amount"] = amount
	resp, err := requestPost("submittip", key, secret, nonce, params)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("SubmitTip failed: %s",
			resp.Message)
	}
	var result string
	return result, json.Unmarshal(resp.Data, &result)
}

// SubmitWithdraw submits a withdrawal request. If address does not exists in you AddressBook, it will fail
// paymentid will be used only for currencies, based of CryptoNote algorhitm
func submitWithdraw(key, secret, nonce, currency, address, paymentid string, amount float64) (int, error) {
	var (
		params = make(map[string]interface{})
		err    error
		cID    int
	)
	if cID, err = getCurrencyID(currency); err != nil {
		return 0, err
	}
	if v, ok := currencyCache[normalize(currency)]; ok {
		if v.Algorithm == "CryptoNote" {
			params["PaymentId"] = paymentid
		}
	}
	params["CurrencyId"] = cID
	params["Address"] = address
	params["Amount"] = amount
	resp, err := requestPost("submitwithdraw", key, secret, nonce, params)
	if err != nil {
		return 0, err
	}
	if !resp.Success {
		return 0, fmt.Errorf("SubmitWithdraw failed: %s, %s %f to %s ",
			resp.Message, currency, amount, address)
	}
	var result int
	return result, json.Unmarshal(resp.Data, &result)
}

// submitTransfer submit a transfer funds to another user
func submitTransfer(key, secret, nonce, currency, username string, amount float64) (string, error) {
	var (
		params = make(map[string]interface{})
		err    error
		cID    int
	)
	if cID, err = getCurrencyID(currency); err != nil {
		return "", err
	}
	params["CurrencyId"] = cID
	params["Username"] = username
	params["Amount"] = amount
	resp, err := requestPost("submittransfer", key, secret, nonce, params)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("SubmitTransfer failed: %s",
			resp.Message)
	}
	var result string
	return result, json.Unmarshal(resp.Data, &result)
}
