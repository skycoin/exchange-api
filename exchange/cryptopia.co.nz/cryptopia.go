package cryptopia

import "time"

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
	ID                   int     `json:"Id"`
	Name                 string  `json:"Name"`
	Symbol               string  `json:"Symbol"`
	Algorithm            string  `json:"Algorithm"`
	WithdrawFee          float64 `json:"WithdrawFee"`
	MinWithdraw          float64 `json:"MinWithdraw"`
	MinBaseTrade         float64 `json:"MinBaseTrade"`
	IsTipEnabled         bool    `json:"IsTipEnabled"`
	MinTip               float64 `json:"MinTip"`
	DepositConfirmations int     `json:"DepositConfirmations"`
	Status               string  `json:"Status"`
	StatusMessage        string  `json:"StatusMessage"`
	ListingStatus        string  `json:"ListingStatus"`
}

// GetCurrencies gets all available currencies from exchange
func GetCurrencies() ([]CurrencyInfo, error) {
	return getCurrencies()
}

// TradepairInfo represents tradepair info
type TradepairInfo struct {
	ID               int     `json:"Id"`
	Label            string  `json:"Label"`
	Currency         string  `json:"Currency"`
	Symbol           string  `json:"Symbol"`
	BaseCurrency     string  `json:"BaseCurrency"`
	BaseSymbol       string  `json:"BaseSymbol"`
	Status           string  `json:"Status"`
	StatusMessage    string  `json:"StatusMessage"`
	TradeFee         float64 `json:"TradeFee"`
	MinimumTrade     float64 `json:"MinimumTrade"`
	MaximumTrade     float64 `json:"MaximumTrade"`
	MinimumBaseTrade float64 `json:"MinimumBaseTrade"`
	MaximumBaseTrade float64 `json:"MaximumBaseTrade"`
	MinimumPrice     float64 `json:"MinimumPrice"`
	MaximumPrice     float64 `json:"MaximumPrice"`
}

// GetTradepairs gets all available tradepairs from exchange
func GetTradepairs() ([]TradepairInfo, error) {
	return getTradePairs()
}

// MarketInfo represents market info
type MarketInfo struct {
	TradePairID    int     `json:"TradePairId"`
	Label          string  `json:"Label"`
	AskPrice       float64 `json:"AskPrice"`
	BidPrice       float64 `json:"BidPrice"`
	Low            float64 `json:"Low"`
	High           float64 `json:"High"`
	Volume         float64 `json:"Volume"`
	LastPrice      float64 `json:"LastPrice"`
	BuyVolume      float64 `json:"BuyVolume"`
	SellVolume     float64 `json:"SellVolume"`
	Change         float64 `json:"Change"`
	Open           float64 `json:"Open"`
	Close          float64 `json:"Close"`
	BaseVolume     float64 `json:"BaseVolume"`
	BaseBuyVolume  float64 `json:"BaseBuyVolume"`
	BaseSellVolume float64 `json:"BaseSellVolume"`
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
	TradePairID int     `json:"TradePairId"`
	Label       string  `json:"Label"`
	Price       float64 `json:"Price"`
	Volume      float64 `json:"Volume"`
	Total       float64 `json:"Total"`
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
	TradePairID int     `json:"TradePairId"`
	Label       string  `json:"Label"`
	Type        string  `json:"Type"`
	Price       float64 `json:"Price"`
	Amount      float64 `json:"Amount"`
	Total       float64 `json:"Total"`
	Timestamp   int     `json:"Timestamp"`
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
	ID            int     `json:"Id"`
	Currency      string  `json:"Currency"`
	TxID          string  `json:"TxId"`
	Type          string  `json:"Type"`
	Amount        float64 `json:"Amount"`
	Fee           float64 `json:"Fee"`
	Status        string  `json:"Status"`
	Confirmations int     `json:"Confirmations"`
	Timestamp     string  `json:"TimeStamp"`
	Address       *string `json:"Address,omitempty"`
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

	Rate      float64
	Amount    float64
	Total     float64
	Fee       float64
	Remaining float64

	Timestamp time.Time
}

// SubmitTrade creates new trade offer, if order instantly executed, it returns ErrInstant
// Trade types defined below
func SubmitTrade(key, secret string, market, Type string, rate, amount float64) (int, error) {
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
func SubmitTip(key, secret string, currency string, activeUsers int, amount float64) (string, error) {
	return submitTip(key, secret, nonce(), currency, activeUsers, amount)
}

// SubmitWithdraw creates withdraw request
// paymentid needs only for currencies, based on CryptoNight
func SubmitWithdraw(key, secret string, currency, address string, paymentid *string, amount float64) (int, error) {
	if paymentid != nil {
		return submitWithdraw(key, secret, nonce(), currency, address, *paymentid, amount)
	}
	return submitWithdraw(key, secret, nonce(), currency, address, "", amount)
}

// SubmitTransfer transfer funds to another user
func SubmitTransfer(key, secret string, currency, username string, amount float64) (string, error) {
	return submitTransfer(key, secret, nonce(), currency, username, amount)
}
