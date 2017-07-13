package cryptopia

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/uberfurrer/tradebot/logger"
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
		logger.Error("cryptopia: http error:", err)
		return nil, err
	}
	return readResponse(resp.Body)
}
func requestPost(endpoint, key, secret, nonce string, params map[string]interface{}) (*response, error) {
	reqURL := apiroot
	reqURL.Path += endpoint
	reqData := encodeValues(params)
	log.Println(string(reqData), reqURL.String())
	req, _ := http.NewRequest("POST", reqURL.String(), bytes.NewReader(reqData))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", header(key, secret, nonce, reqURL, reqData))
	resp, err := httpclient.Do(req)

	if err != nil {
		logger.Error("cryptopia: http error:", err)
		return nil, err
	}
	return readResponse(resp.Body)
}

//Public API functions

// GetCurrencies gets all currencies
func GetCurrencies() ([]CurrencyInfo, error) {
	resp, err := requestGet("getcurrencies", "")
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetCurrencies failed: %s",
			resp.Message)
	}
	var result []CurrencyInfo
	return result, json.Unmarshal(resp.Data, &result)
}

// GetTradePairs gets all TradePairs on exchange
func GetTradePairs() ([]TradepairInfo, error) {
	resp, err := requestGet("gettradepairs", "")
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetTradePairs failed: %s",
			resp.Message)
	}
	var result []TradepairInfo
	return result, json.Unmarshal(resp.Data, &result)
}

// GetMarkets return all Market info by given baseMarket
// if baseMarket is empty or "all" GetMarkets return all markets
// if hours < 1 it will be omitted, default value is 24
func GetMarkets(baseMarket string, hours int) ([]MarketInfo, error) {
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
		return nil, errors.Errorf("GetMarkets failed: %s",
			resp.Message)
	}
	var result []MarketInfo
	return result, json.Unmarshal(resp.Data, &result)
}

// GetMarket return market with given label
// if hours < 1, it will be omitted, default value is 24
func GetMarket(market string, hours int) (*MarketInfo, error) {
	var (
		requestParams string
		marketID      int
		err           error
	)

	if marketID, err = getMarketID(market); err != nil {
		return nil, err
	}
	requestParams += strconv.Itoa(marketID)
	if hours > 0 {
		requestParams += "/" + strconv.Itoa(hours)
	}
	resp, err := requestGet("getmarket", requestParams)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetMarket failed: %s, Market: %s",
			resp.Message, market)
	}
	var result = &MarketInfo{}
	return result, json.Unmarshal(resp.Data, result)
}

// GetMarketHistory return market history with given label
// if hours < 1, it will be omitted, default value is 24
func GetMarketHistory(market string, hours int) ([]MarketHistoryResponse, error) {
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
		return nil, errors.Errorf("GetMarketHistory failed: %s, Market: %s",
			resp.Message, market)
	}
	var result []MarketHistoryResponse
	return result, json.Unmarshal(resp.Data, &result)
}

// GetMarketOrders returns count orders from market with given label
// if count < 1, its will be omitted, default value is 100
func GetMarketOrders(market string, count int) (*OrderBook, error) {
	var (
		requestParams string
		err           error
		marketID      int
	)
	if marketID, err = getMarketID(market); err != nil {
		return nil, err
	}
	requestParams += strconv.Itoa(marketID)
	if count > 0 {
		requestParams += "/" + strconv.Itoa(count)
	}
	resp, err := requestGet("getmarketorders", requestParams)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetMarketHistory failed: %s, Market: %s",
			resp.Message, market)
	}
	var result = &OrderBook{}
	return result, json.Unmarshal(resp.Data, &result)
}

// GetMarketOrderGroups returns count Orders to each market
// If count < 1, it will be omitted
func GetMarketOrderGroups(count int, markets ...string) ([]OrderBookLabeled, error) {
	var (
		requestParams string
		err           error
		marketID      int
	)
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
		return nil, errors.Wrapf(err, "GetMarketOrderGroups failed: %s, Market: %s")
	}
	if !resp.Success {
		return nil, errors.Errorf("GetMarketOrderGroups failed: %s, Market: %s",
			resp.Message, markets)
	}
	var result []OrderBookLabeled
	return result, json.Unmarshal(resp.Data, &result)
}

// Private API functions

//GetBalance return a string representation of balance by given currency
func GetBalance(key, secret, nonce, currency string) (string, error) {
	resp, err := requestPost("getbalance", key, secret, nonce, nil)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", errors.Errorf("GetBalance failed: %s, Currency %s Rawdata %s",
			resp.Message, currency, string(resp.Data))
	}
	var result Balance
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return "", err
	}
	if v, ok := result[normalize(currency)]; ok {
		return v, nil
	}
	return "", errors.New("currency does not found")
}

// GetDepositAddress returns a deposit address of given currency
func GetDepositAddress(key, secret, nonce, currency string) (*DepositAddress, error) {
	var params = make(map[string]interface{})
	cID, err := getCurrencyID(currency)
	if err != nil {
		return nil, errors.Errorf("Currency %s does not found", currency)
	}
	params["CurrencyId"] = cID
	resp, err := requestPost("getdepositaddress", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetDepositAddress failed: %s, Currency %s",
			resp.Message, currency)
	}
	var result DepositAddress
	return &result, json.Unmarshal(resp.Data, &result)
}

// AllMarkets for GetTradeHistory and GetOpenOrders
const (
	AllMarkets = "ALL"
)

// GetOpenOrders return a list of opened orders by specific market
// if count < 1, it will be omitted
func GetOpenOrders(key, secret, nonce, market string, count int) ([]OpenedOrder, error) {
	var (
		params = make(map[string]interface{})
		mID    int
		err    error
	)
	if market != AllMarkets {
		if mID, err = getMarketID(market); err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	}
	if count > 0 {
		params["Count"] = count
	}
	resp, err := requestPost("getopenorders", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetOpenOrders failed: %s Market %s Count %d",
			resp.Message, market, count)
	}
	var result []OpenedOrder
	return result, json.Unmarshal(resp.Data, &result)
}

// GetTradeHistory return a list of all executed orders by specific market
// if count < 1, it will be omitted
func GetTradeHistory(key, secret, nonce, market string, count int) ([]ClosedOrder, error) {
	var (
		params = make(map[string]interface{})
		err    error
		mID    int
	)
	if market != AllMarkets {
		if mID, err = getMarketID(market); err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	}
	if count > 0 {
		params["Count"] = count
	}
	resp, err := requestPost("gettradehistory", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("GetTradeHistory failed: %s Market %s Count %d",
			resp.Message, market, count)
	}
	var result []ClosedOrder
	return result, json.Unmarshal(resp.Data, &result)
}

//Possible types of transaction
const (
	TxTypeWithdraw = "Withdraw"
	TxTypeDeposit  = "Deposit"
)

// GetTransactions returns a list of transactions by given type
// if count < 1, it will be omitted
func GetTransactions(key, secret, nonce, Type string, count int) ([]Transaction, error) {
	var params = make(map[string]interface{})
	if Type = strings.Title(Type); Type != TxTypeDeposit && Type != TxTypeWithdraw {
		return nil, errors.Errorf("Icorrect trasnaction type %s; avalible types: %s %s",
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
		return nil, errors.Errorf("GetTransactions failed: %s Type %s Count %d",
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

// SubmitTrade submits a new trade offer
func SubmitTrade(key, secret, nonce, market, Type string, rate, amount float64) (*NewTradeInfo, error) {
	var (
		params = make(map[string]interface{})
		err    error
		mID    int
	)
	if Type = strings.Title(Type); Type != OfTypeBuy && Type != OfTypeSell {
		return nil, errors.Errorf("Incorrect offer type %s; avalible types: %s %s",
			Type, OfTypeBuy, OfTypeSell)
	}
	if mID, err = getMarketID(market); err != nil {
		return nil, err
	}
	params["TradePairId"] = mID
	params["Type"] = Type
	params["Rate"] = rate
	params["Amount"] = amount
	resp, err := requestPost("submittrade", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("SubmitTrade failed: %s, Type %s Market %s Rate %f Amount %f",
			resp.Message, Type, market, rate, amount)
	}
	var result NewTradeInfo
	return &result, json.Unmarshal(resp.Data, &result)
}

// Possible types of cancellation
const (
	CancelOne       = "Trade"
	CancelTradePair = "TradePair"
	CancelAll       = "All"
)

// CancelTrade cancel trades by given orderid, market or add active
// depends of type argument
func CancelTrade(key, secret, nonce, Type string, orderID int, market string) (CancelledOrders, error) {
	var params = make(map[string]interface{})
	Type = strings.Title(Type)
	switch Type {
	case CancelAll:
		break
	case CancelOne:
		params["OrderId"] = orderID
	case CancelTradePair:
		var (
			mID int
			err error
		)
		if mID, err = getMarketID(market); err != nil {
			return nil, err
		}
		params["TradePairId"] = mID
	default:
		return nil, errors.Errorf("Incorrect type of cancellation %s; possilble types: %s %s %s",
			Type, CancelAll, CancelOne, CancelTradePair)
	}
	params["Type"] = Type
	resp, err := requestPost("canceltrade", key, secret, nonce, params)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.Errorf("CancelTrade failed: %s params: %v",
			resp.Message, params)
	}
	var result CancelledOrders
	return result, json.Unmarshal(resp.Data, &result)
}

// SubmitTip submits a tip to Trollbox
func SubmitTip(key, secret, nonce, currency string, activeUsers int, amount float64) (TipMessage, error) {
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
		return "", errors.Errorf("SubmitTip failed: %s",
			resp.Message)
	}
	var result TipMessage
	return result, json.Unmarshal(resp.Data, &result)
}

// SubmitWithdraw submits a withdrawal request. If address does not exists in you AddressBook, it will fail
// paymentid will be used only for currencies, based of CryptoNote algorhitm
func SubmitWithdraw(key, secret, nonce, currency, address, paymentid string, amount float64) (WithdrawID, error) {
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
		return 0, errors.Errorf("SubmitWithdraw failed: %s, %s %f to %s ",
			resp.Message, currency, amount, address)
	}
	var result WithdrawID
	return result, json.Unmarshal(resp.Data, &result)
}

// SubmitTransfer submit a transfer funds to another user
func SubmitTransfer(key, secret, nonce, currency, username string, amount float64) (TransferMessage, error) {
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
		return "", errors.Errorf("SubmitTransfer failed: %s",
			resp.Message)
	}
	var result TransferMessage
	return result, json.Unmarshal(resp.Data, &result)
}
