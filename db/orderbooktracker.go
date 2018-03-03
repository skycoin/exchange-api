package db

import (
	"strings"

	"github.com/skycoin/exchange-api/exchange"
)

type orderbooktracker struct {
	db OrderDatabase
}

func (t *orderbooktracker) Update(sym string, Bids []exchange.MarketOrder, Asks []exchange.MarketOrder) {
	t.db.Update(sym, Bids, Asks)
}

// Get gets information about stock
func (t *orderbooktracker) Get(sym string) (*exchange.MarketRecord, error) {
	return t.db.Get(sym)
}

// NewOrderbookTracker returns exchange.OrderbookTracker
// that wraps either redis connection or sync.Map
// For in-memory tracker dbURL and hash are optional parameters
func NewOrderbookTracker() (exchange.Orderbooks, error) {
	db, err := NewDatabase()

	if err != nil {
		return nil, err
	}

	return &orderbooktracker{
		db: db,
	}, nil
}

func normalize(symbol string) string {
	return strings.ToUpper(strings.Replace(symbol, "_", "/", -1))
}
