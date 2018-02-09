package exchange

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Orderbooks provides functionality for managing OrderBook
type Orderbooks interface {
	// Update updates orderbook for given market
	Update(string, []MarketOrder, []MarketOrder)
	//Get gets orderbook for given tradepair symbol
	Get(string) (*MarketRecord, error)
}

// MarketOrder is a one order in stock
type MarketOrder struct {
	Price  decimal.Decimal `json:"price"`
	Volume decimal.Decimal `json:"volume"`
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
		Time:   r.Timestamp.UnixNano() / 10e5,
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
	r.Timestamp = time.Unix(tmp.Time/10e2, (tmp.Time%10e2)*10e5)
	r.Symbol = tmp.Symbol
	return nil
}
