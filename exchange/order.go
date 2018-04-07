package exchange

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// MarketOrder is a one order in stock
type MarketOrder struct {
	Price  decimal.Decimal `json:"price"`
	Volume decimal.Decimal `json:"volume"`
}

// MarketOrder is a one order in stock
type MarketOrderString struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

// TotalCost returns the total cost of executing the given order
func (marketOrder MarketOrder) TotalCost() decimal.Decimal {
	return marketOrder.Price.Mul(marketOrder.Volume)
}

// MarketRecord represents orderbook for one market
type MarketRecord struct {
	Timestamp time.Time     `json:"timestamp"`
	Symbol    string        `json:"symbol"`
	Bids      []MarketOrder `json:"bids"`
	Asks      []MarketOrder `json:"asks"`
}

// MarshalJSON implements json.Marshaler interface
func (r MarketRecord) MarshalJSON() ([]byte, error) {
	type rec struct {
		Time   int64           `json:"timestamp"`
		Symbol string          `json:"symbol"`
		Bids   json.RawMessage `json:"bids"`
		Asks   json.RawMessage `json:"asks"`
	}
	var (
		bids, asks json.RawMessage
		err        error
	)
	if bids, err = json.Marshal(r.Bids); err != nil {
		return nil, err
	}
	if asks, err = json.Marshal(r.Asks); err != nil {
		return nil, err
	}
	var tmp = rec{
		Time:   r.Timestamp.Unix(),
		Symbol: r.Symbol,
		Bids:   bids,
		Asks:   asks,
	}
	return json.Marshal(tmp)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (r *MarketRecord) UnmarshalJSON(b []byte) error {

	type rec struct {
		Time   int64           `json:"timestamp"`
		Symbol string          `json:"symbol"`
		Bids   json.RawMessage `json:"bids"`
		Asks   json.RawMessage `json:"asks"`
	}

	var tmp rec
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	var (
		bids []MarketOrder
		asks []MarketOrder
	)
	if err = json.Unmarshal(tmp.Bids, &bids); err != nil {
		return err
	}
	if err = json.Unmarshal(tmp.Asks, &asks); err != nil {
		return err
	}
	r.Asks = asks
	r.Bids = bids
	r.Timestamp = time.Unix(tmp.Time, 0)
	r.Symbol = tmp.Symbol
	return nil
}

// MarketOrders alias for []MarketOrder
type MarketOrders []MarketOrder

// MarketOrders alias for []MarketOrderString
type MarketOrdersString []MarketOrderString

// Volume returns the sum of a set of MarketOrders' volumes
func (marketOrders MarketOrders) Volume() decimal.Decimal {
	var sum = decimal.Zero

	for _, order := range marketOrders {
		sum = sum.Add(order.Volume)
	}

	return sum
}

var (
	// ErrNegativeAmount SpendItAll error for when called with <0 currency
	ErrNegativeAmount = errors.New("can't spend negative quantities of currency")
	// ErrOrdersRanOut SpendItAll error for when the caller tries to purchase more coins than are available in the orderbook
	ErrOrdersRanOut = errors.New("ran out of orders before we ran out of currency")
)

// SpendItAll determines the cheapest series of purchases necessary to spend the specified quantity of coins. It can fail if there aren't enough standing orders available to cover the purchase or if the user specifies a negative quantity of coins.
func (r *MarketRecord) SpendItAll(amount decimal.Decimal) (MarketOrders, error) {
	if amount.LessThan(decimal.Zero) {
		return nil, ErrNegativeAmount
	}

	if amount.Equal(decimal.Zero) {
		return nil, nil
	}

	sort.Slice(r.Asks, func(first, second int) bool {
		return r.Asks[first].Price.LessThan(r.Asks[second].Price)
	})

	var orders []MarketOrder

	for _, order := range r.Asks {
		maxSpend := order.Price.Mul(order.Volume)
		actualSpend := decimal.Min(maxSpend, amount)
		volume := actualSpend.Div(order.Price)

		newOrder := MarketOrder{
			Price:  order.Price,
			Volume: volume,
		}

		orders = append(orders, newOrder)

		amount = amount.Sub(newOrder.Price.Mul(newOrder.Volume))

		if amount.Equal(decimal.Zero) {
			break
		}
	}

	if amount.GreaterThan(decimal.Zero) {
		return orders, ErrOrdersRanOut
	}

	return orders, nil
}

// CheapestAsk returns the cheapest Ask order. If there are two ask orders with the same price, it returns the one with the larger volume.
func (r *MarketRecord) CheapestAsk() *MarketOrder {
	if len(r.Asks) == 0 {
		return nil
	}

	result := r.Asks[0]
	for _, order := range r.Asks[1:] {
		switch {
		case order.Price.LessThan(result.Price):
			result = order
		case order.Price.Equal(result.Price) && order.Volume.GreaterThan(result.Volume):
			result = order
		default:
		}
	}

	return &result
}
