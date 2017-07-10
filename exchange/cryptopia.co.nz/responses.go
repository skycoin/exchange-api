package cryptopia

import (
	"encoding/json"
	"fmt"
	"strings"
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

// MarketHistoryResponse represents market history
type MarketHistoryResponse struct {
	TradePairID int     `json:"TradePairId"`
	Label       string  `json:"Label"`
	Type        string  `json:"Type"`
	Price       float64 `json:"Price"`
	Amount      float64 `json:"Amount"`
	Total       float64 `json:"Total"`
	Timestamp   int     `json:"Timestamp"`
}

// MarketOrder represents a single order info
type MarketOrder struct {
	TradePairID int     `json:"TradePairId"`
	Label       string  `json:"Label"`
	Price       float64 `json:"Price"`
	Volume      float64 `json:"Volume"`
	Total       float64 `json:"Total"`
}

// OrderBook is a response that was received from GetMarketOrders function
type OrderBook struct {
	Buy  []MarketOrder `json:"Buy"`
	Sell []MarketOrder `json:"Sell"`
}

// OrderBookLabeled is a response that was received from GetMarketOrderGroups function
type OrderBookLabeled struct {
	TradePairID int           `json:"TradePairId"`
	Label       string        `json:"Market"`
	Buy         []MarketOrder `json:"Buy"`
	Sell        []MarketOrder `json:"Sell"`
}

// Balance represents Balance of all avalible currencies
type Balance map[string]string

// UnmarshalJSON implements json.Unmarshaler interface
func (r *Balance) UnmarshalJSON(b []byte) error {
	if r == nil {
		(*r) = make(map[string]string)
	}
	type currency struct {
		CurrencyID      int     `json:"CurrencyId"`
		Symbol          string  `json:"Symbol"`
		Total           float64 `json:"Total"`
		Available       float64 `json:"Available"`
		Unconfirmed     float64 `json:"Unconfirmed"`
		HeldForTrades   float64 `json:"HeldForTrades"`
		PendingWithdraw float64 `json:"PendingWithdraw"`
		Address         string  `json:"Address"`
		BaseAddress     string  `json:"BaseAddress"`
		Status          string  `json:"Status"`
		StatusMessage   string  `json:"StatusMessage"`
	}

	var tmp = make([]currency, 0)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	var result = make(Balance)
	for _, v := range tmp {
		result[strings.ToUpper(v.Symbol)] = fmt.Sprintf("Total: %.8f Available: %.8f Unconfirmed: %.8f Held: %.8f Pending: %.8f",
			v.Total, v.Available, v.Unconfirmed, v.HeldForTrades, v.PendingWithdraw)
	}
	*r = result
	return nil
}

// DepositAddress is a representation of deposit address for single currency
type DepositAddress struct {
	Currency    string `json:"Currency"`
	Address     string `json:"Address"`
	BaseAddress string `json:"BaseAddress"`
}

// OpenedOrder represpent a single order, that returns from GetOpenOrders
type OpenedOrder struct {
	OrderID     int     `json:"OrderId"`
	TradePairID int     `json:"TradePairId"`
	Market      string  `json:"Market"`
	Type        string  `json:"Type"`
	Rate        float64 `json:"Rate"`
	Amount      float64 `json:"Amount"`
	Total       string  `json:"Total"`
	Remaining   string  `json:"Remaining"`
	Timestamp   string  `json:"TimeStamp"`
}

// ClosedOrder represents a single order, that returns from GetTradeHistory
type ClosedOrder struct {
	OrderID     int     `json:"TradeIdfloat64"`
	TradePairID int     `json:"TradePairId"`
	Market      string  `json:"Market"`
	Type        string  `json:"Type"`
	Rate        float64 `json:"Rate"`
	Amount      float64 `json:"Amount"`
	Total       string  `json:"Total"`
	Fee         string  `json:"Fee"`
	Timestamp   string  `json:"TimeStamp"`
}

// Transaction represents a single transaction, deposit or withdraw
// returns from GetTrasnactions
type Transaction struct {
	ID            int     `json:"Id"`
	Currency      string  `json:"Currency"`
	TxID          string  `json:"TxId"`
	Type          string  `json:"Type"`
	Amount        float64 `json:"Amount"`
	Fee           string  `json:"Fee"`
	Status        string  `json:"Status"`
	Confirmations string  `json:"Confirmations"`
	Timestamp     string  `json:"TimeStamp"`
	Address       string  `json:"Address,omitempty"`
}

// NewTradeInfo represents success created order
// if OrderID == 0, order completed instantly
// if FilledOrders empty - order opened, but does not filled
// else order partitally filled
type NewTradeInfo struct {
	OrderID      int   `json:"OrderId,omitempty"`
	FilledOrders []int `json:"FilledOrders,omitempty"`
}

// CancelledOrders represents all cancelled orders
type CancelledOrders []int

// TipMessage represnets message from SubmiTip
type TipMessage string

// WithdrawID is a withdraw tx ID
type WithdrawID int

// TransferMessage represents message from SubmitTransfer
type TransferMessage string
