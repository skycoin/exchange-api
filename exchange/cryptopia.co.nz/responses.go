package cryptopia

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// balance represents balance of all avalible currencies
type balance map[string]string

type currency struct {
	CurrencyID      int             `json:"CurrencyId"`
	Symbol          string          `json:"Symbol"`
	Total           decimal.Decimal `json:"Total"`
	Available       decimal.Decimal `json:"Available"`
	Unconfirmed     decimal.Decimal `json:"Unconfirmed"`
	HeldForTrades   decimal.Decimal `json:"HeldForTrades"`
	PendingWithdraw decimal.Decimal `json:"PendingWithdraw"`
	Address         string          `json:"Address"`
	BaseAddress     string          `json:"BaseAddress"`
	Status          string          `json:"Status"`
	StatusMessage   string          `json:"StatusMessage"`
}

// UnmarshalJSON implements json.Unmarshaler interface
func (r *balance) UnmarshalJSON(b []byte) error {
	if r == nil {
		(*r) = make(map[string]string)
	}

	var tmp = make([]currency, 0)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	var result = make(balance)

	for _, v := range tmp {
		result[strings.ToUpper(v.Symbol)] = fmt.Sprintf("Total: %s Available: %s Unconfirmed: %s Held: %s Pending: %s",
			v.Total.StringFixed(8),
			v.Available.StringFixed(8),
			v.Unconfirmed.StringFixed(8),
			v.HeldForTrades.StringFixed(8),
			v.PendingWithdrawal.StringFixed(8))
	}
	*r = result
	return nil
}

// newOrder represents success created order
// if OrderID == 0, order completed instantly
// if FilledOrders empty - order opened, but does not filled
// else order partitally filled
type newOrder struct {
	OrderID      *int  `json:"OrderId,omitempty"`
	FilledOrders []int `json:"FilledOrders,omitempty"`
}

type orderJSON struct {
	OrderID     *int   `json:"OrderId,omitempty"`
	TradeID     *int   `json:"TradeId,omitempty"`
	TradePairID int    `json:"TradePairId"`
	Market      string `json:"Market"`
	Type        string `json:"Type"`

	Rate      decimal.Decimal `json:"Rate"`
	Amount    decimal.Decimal `json:"Amount"`
	Total     decimal.Decimal `json:"Total"`
	Fee       decimal.Decimal `json:"Fee,omitempty"`
	Remaining decimal.Decimal `json:"Remaining,omitempty"`

	Timestamp string `json:"TimeStamp"`
}

// UnmarshalJSON implements an json.Unmarshaler interface
func (order *Order) UnmarshalJSON(b []byte) error {
	var (
		tmp     = orderJSON{}
		orderID int
		ts      time.Time
	)

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	if tmp.OrderID == nil && tmp.TradeID == nil {
		return errZeroOrderID
	}
	if tmp.OrderID != nil {
		orderID = *tmp.OrderID
	} else {
		orderID = *tmp.TradeID
	}
	ts, err = time.Parse("2006-01-02T15:04:05.0000000", tmp.Timestamp)
	if err != nil {
		return err
	}

	*order = Order{
		OrderID:     orderID,
		TradePairID: tmp.TradePairID,
		Market:      tmp.Market,
		Type:        tmp.Type,
		Rate:        tmp.Rate,
		Amount:      tmp.Amount,
		Total:       tmp.Total,
		Fee:         tmp.Fee,
		Remaining:   tmp.Remaining,
		Timestamp:   ts,
	}
	return nil
}

var errZeroOrderID = errors.New("parse error, OrderID and TradeID is empty")
