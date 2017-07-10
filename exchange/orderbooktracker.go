package exchange

import (
	"encoding/json"
	"time"
)

// OrderBookTracker provides functionality for managing OrderBook
type OrderBookTracker interface {
	// UpdateSym updates buffer for the given sym
	// After updating all symbols you need to call OrderBookTracker.Flush()
	UpdateSym(string, []OrderbookEntry, []OrderbookEntry)
	//Get gets orderbook for given tradepair symbol
	// It is case-sensetive
	// First returned value - bids, second - asks
	GetRecord(string) (Orderbook, error)
}

// OrderbookEntry entry represnets a single order in orderbook
type OrderbookEntry struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

// Orderbook represents orderbook for one market
type Orderbook struct {
	Timestamp time.Time        `json:"timestamp"`
	Symbol    string           `json:"symbol"`
	Bids      []OrderbookEntry `json:"bids"`
	Asks      []OrderbookEntry `json:"asks"`
}

// MarshalJSON implements json.Marshaler interface
func (r Orderbook) MarshalJSON() ([]byte, error) {
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
func (r *Orderbook) UnmarshalJSON(b []byte) error {

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
		bids []OrderbookEntry
		asks []OrderbookEntry
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
