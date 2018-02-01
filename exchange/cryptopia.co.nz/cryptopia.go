package cryptopia

import (
	"time"

	"github.com/shopspring/decimal"
)

// Order types, buy or sell
const (
	OrderTypeBuy  = "Buy"
	OrderTypeSell = "Sell"
)

// Transaction types
const (
	TxTypeDeposit  = "Deposit"
	TxTypeWithdraw = "Withdraw"
)

// Cancellation types
const (
	CancelTypeAll    = "All"
	CancelTypeMarket = "Market"
	CancelTypeOrder  = "Trade"
)

var (
	currencyCache map[string]CurrencyInfo
	marketCache   map[string]int
)

//CurrencyInfo represents currency info
type CurrencyInfo struct {
	ID                   int             `json:"Id"`
	Name                 string          `json:"Name"`
	Symbol               string          `json:"Symbol"`
	Algorithm            string          `json:"Algorithm"`
	WithdrawFee          decimal.Decimal `json:"WithdrawFee"`
	MinWithdraw          decimal.Decimal `json:"MinWithdraw"`
	MinBaseTrade         decimal.Decimal `json:"MinBaseTrade"`
	IsTipEnabled         bool            `json:"IsTipEnabled"`
	MinTip               decimal.Decimal `json:"MinTip"`
	DepositConfirmations int             `json:"DepositConfirmations"`
	Status               string          `json:"Status"`
	StatusMessage        string          `json:"StatusMessage"`
	ListingStatus        string          `json:"ListingStatus"`
}

// GetCurrencies gets all availible currencies from exchange
func GetCurrencies() ([]CurrencyInfo, error) {
	return getCurrencies()
}

// TradepairInfo represents tradepair info
type TradepairInfo struct {
	ID               int             `json:"Id"`
	Label            string          `json:"Label"`
	Currency         string          `json:"Currency"`
	Symbol           string          `json:"Symbol"`
	BaseCurrency     string          `json:"BaseCurrency"`
	BaseSymbol       string          `json:"BaseSymbol"`
	Status           string          `json:"Status"`
	StatusMessage    string          `json:"StatusMessage"`
	TradeFee         decimal.Decimal `json:"TradeFee"`
	MinimumTrade     decimal.Decimal `json:"MinimumTrade"`
	MaximumTrade     decimal.Decimal `json:"MaximumTrade"`
	MinimumBaseTrade decimal.Decimal `json:"MinimumBaseTrade"`
	MaximumBaseTrade decimal.Decimal `json:"MaximumBaseTrade"`
	MinimumPrice     decimal.Decimal `json:"MinimumPrice"`
	MaximumPrice     decimal.Decimal `json:"MaximumPrice"`
}

// GetTradepairs gets all availible tradepairs from exchange
func GetTradepairs() ([]TradepairInfo, error) {
	return getTradePairs()
}

// MarketInfo represents market info
type MarketInfo struct {
	TradePairID    int             `json:"TradePairId"`
	Label          string          `json:"Label"`
	AskPrice       decimal.Decimal `json:"AskPrice"`
	BidPrice       decimal.Decimal `json:"BidPrice"`
	Low            decimal.Decimal `json:"Low"`
	High           decimal.Decimal `json:"High"`
	Volume         decimal.Decimal `json:"Volume"`
	LastPrice      decimal.Decimal `json:"LastPrice"`
	BuyVolume      decimal.Decimal `json:"BuyVolume"`
	SellVolume     decimal.Decimal `json:"SellVolume"`
	Change         decimal.Decimal `json:"Change"`
	Open           decimal.Decimal `json:"Open"`
	Close          decimal.Decimal `json:"Close"`
	BaseVolume     decimal.Decimal `json:"BaseVolume"`
	BaseBuyVolume  decimal.Decimal `json:"BaseBuyVolume"`
	BaseSellVolume decimal.Decimal `json:"BaseSellVolume"`
}

// GetMarkets gets all market by specifying baseMarket
// if hours == nil, default value(24) is used
func GetMarkets(baseMarket string, hours *int) ([]MarketInfo, error) {
	if hours != nil {
		return getMarkets(baseMarket, *hours)
	}
	return getMarkets(baseMarket, 24)
}

// GetMarket returns informaiton abot market with given label
// if market not found, it returns error
// if hours = nil, default value(24) is used
func GetMarket(label string, hours *int) (MarketInfo, error) {
	if hours != nil {
		return getMarket(label, *hours)
	}
	return getMarket(label, 24)
}

// MarketOrders is a orderbook for market
type MarketOrders struct {
	Buy  []MarketOrder `json:"Buy"`
	Sell []MarketOrder `json:"Sell"`
}

// MarketOrder represents a single order info
type MarketOrder struct {
	TradePairID int             `json:"TradePairId"`
	Label       string          `json:"Label"`
	Price       decimal.Decimal `json:"Price"`
	Volume      decimal.Decimal `json:"Volume"`
	Total       decimal.Decimal `json:"Total"`
}

// GetMarketOrders returns orderbook for given market
// if count == nil, default value(100) is used
func GetMarketOrders(label string, count *int) (MarketOrders, error) {
	if count != nil {
		return getMarketOrders(label, *count)
	}
	return getMarketOrders(label, 100)
}

// MarketOrdersWithLabel is a response that was received from GetMarketOrderGroups function
type MarketOrdersWithLabel struct {
	TradePairID int           `json:"TradePairId"`
	Label       string        `json:"Market"`
	Buy         []MarketOrder `json:"Buy"`
	Sell        []MarketOrder `json:"Sell"`
}

// GetMarketOrderGroups returns MarketOrders for given markets
func GetMarketOrderGroups(count int, markets ...string) ([]MarketOrdersWithLabel, error) {
	return getMarketOrderGroups(count, markets...)
}

// MarketHistory represents market history
type MarketHistory struct {
	TradePairID int             `json:"TradePairId"`
	Label       string          `json:"Label"`
	Type        string          `json:"Type"`
	Price       decimal.Decimal `json:"Price"`
	Amount      decimal.Decimal `json:"Amount"`
	Total       decimal.Decimal `json:"Total"`
	Timestamp   int             `json:"Timestamp"`
}

// GetMarketHistory returns completed orders in given market for given time
// if hours == nil, default value(24) is used
func GetMarketHistory(label string, hours *int) ([]MarketHistory, error) {
	if hours != nil {
		return getMarketHistory(label, *hours)
	}
	return getMarketHistory(label, 24)
}

// GetBalance returns a string representation of balance by given currency
func GetBalance(key, secret string, currency string) (string, error) {
	return getBalance(key, secret, nonce(), currency)
}

// DepositAddress is a representation of deposit address for single currency
type DepositAddress struct {
	Currency    string `json:"Currency"`
	Address     string `json:"Address"`
	BaseAddress string `json:"BaseAddress"`
}

// GetDepositAddress gets new deposit address for currency
func GetDepositAddress(key, secret string, currency string) (DepositAddress, error) {
	return getDepositAddress(key, secret, nonce(), currency)
}

// GetOpenOrders gets count opened orders from given market
// if market == nil, then gets from all markets
// if count == nil, default value(100) is used
func GetOpenOrders(key, secret string, market *string, count *int) ([]Order, error) {
	return getOpenOrders(key, secret, nonce(), market, count)
}

// GetTradeHistory same as GetOpenOrders
func GetTradeHistory(key, secret string, market *string, count *int) ([]Order, error) {
	return getTradeHistory(key, secret, nonce(), market, count)
}

// Transaction types
const (
	Deposit  = "Deposit"
	Withdraw = "Withdraw"
)

// Transaction represents a single transaction, deposit or withdraw
type Transaction struct {
	ID            int             `json:"Id"`
	Currency      string          `json:"Currency"`
	TxID          string          `json:"TxId"`
	Type          string          `json:"Type"`
	Amount        decimal.Decimal `json:"Amount"`
	Fee           decimal.Decimal `json:"Fee"`
	Status        string          `json:"Status"`
	Confirmations int             `json:"Confirmations"`
	Timestamp     string          `json:"TimeStamp"`
	Address       *string         `json:"Address,omitempty"`
}

// GetTransactions gets count transactions with given type
// if count == nil, default value(?) is used
// Transaction types defined below
func GetTransactions(key, secret string, Type string, count *int) ([]Transaction, error) {
	if count != nil {
		return getTransactions(key, secret, nonce(), Type, *count)
	}
	return getTransactions(key, secret, nonce(), Type, 0)
}

// Order types
const (
	Buy  = "Buy"
	Sell = "Sell"
)

// Order represents single opened or closed order
// If order was closed, fee > 0 && remaining == 0
type Order struct {
	OrderID     int
	TradePairID int
	Market      string
	Type        string

	Rate      decimal.Decimal
	Amount    decimal.Decimal
	Total     decimal.Decimal
	Fee       decimal.Decimal
	Remaining decimal.Decimal

	Timestamp time.Time
}

// SubmitTrade creates new trade offer, if order instantly executed, it returns ErrInstant
// Trade types defined below
func SubmitTrade(key, secret string, market, Type string, rate, amount decimal.Decimal) (int, error) {
	return submitTrade(key, secret, nonce(), market, Type, rate, amount)
}

// Types of cancellation
const (
	All       = "All"
	ByMarket  = "TradePair"
	ByOrderID = "Trade"
)

// CancelTrade cancel all trades, trades in given market, or trade with given orderID
// Cancellation types defined below
// if Type == ByMarket, then tradepair should not be nil
// if Type == ByOrderID, then orderID should be not nil
func CancelTrade(key, secret string, Type string, tradepair *string, orderID *int) ([]int, error) {
	return cancelTrade(key, secret, nonce(), Type, tradepair, orderID)
}

// SubmitTip is useless
func SubmitTip(key, secret string, currency string, activeUsers int, amount decimal.Decimal) (string, error) {
	return submitTip(key, secret, nonce(), currency, activeUsers, amount)
}

// SubmitWithdraw creates withdraw request
// paymentid needs only for currencies, based on CryptoNight
func SubmitWithdraw(key, secret string, currency, address string, paymentid *string, amount decimal.Decimal) (int, error) {
	if paymentid != nil {
		return submitWithdraw(key, secret, nonce(), currency, address, *paymentid, amount)
	}
	return submitWithdraw(key, secret, nonce(), currency, address, "", amount)
}

// SubmitTransfer transfer funds to another user
func SubmitTransfer(key, secret string, currency, username string, amount decimal.Decimal) (string, error) {
	return submitTransfer(key, secret, nonce(), currency, username, amount)
}
