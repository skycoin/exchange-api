// Package cryptopia provides api methods for communicating with cryptopia exchange
package cryptopia

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	// OfferTypeBuy a buy order
	OfferTypeBuy = "Buy"
	// OfferTypeSell a sell order
	OfferTypeSell = "Sell"
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

// MarketOrdersWithLabel is a response that was received from GetMarketOrderGroups function
type MarketOrdersWithLabel struct {
	TradePairID int           `json:"TradePairId"`
	Label       string        `json:"Market"`
	Buy         []MarketOrder `json:"Buy"`
	Sell        []MarketOrder `json:"Sell"`
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

// DepositAddress is a representation of deposit address for single currency
type DepositAddress struct {
	Currency    string `json:"Currency"`
	Address     string `json:"Address"`
	BaseAddress string `json:"BaseAddress"`
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

// Types of cancellation
const (
	All       = "All"
	ByMarket  = "TradePair"
	ByOrderID = "Trade"
)
