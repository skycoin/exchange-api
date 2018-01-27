package cryptopia

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"errors"
)

// balance represents balance of all avalible currencies
type balance map[string]string

// UnmarshalJSON implements json.Unmarshaler interface
func (r *balance) UnmarshalJSON(b []byte) error {
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
	var result = make(balance)
	for _, v := range tmp {
		result[strings.ToUpper(v.Symbol)] = fmt.Sprintf("Total: %.8f Available: %.8f Unconfirmed: %.8f Held: %.8f Pending: %.8f",
			v.Total, v.Available, v.Unconfirmed, v.HeldForTrades, v.PendingWithdraw)
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

// UnmarshalJSON implements an json.Unmarshaler interface
func (order *Order) UnmarshalJSON(b []byte) error {
	var (
		tmp = struct {
			OrderID     *int    `json:"OrderId,omitempty"`
			TradeID     *int    `json:"TradeId,omitempty"`
			TradePairID int     `json:"TradePairId"`
			Market      string  `json:"Market"`
			Type        string  `json:"Type"`
			Rate        float64 `json:"Rate"`
			Amount      float64 `json:"Amount"`
			Total       float64 `json:"Total"`
			Fee         float64 `json:"Fee,omitempty"`
			Remaining   float64 `json:"Remaining,omitempty"`
			Timestamp   string  `json:"TimeStamp"`
		}{}
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
