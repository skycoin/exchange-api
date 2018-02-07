package db

import (
	"github.com/skycoin/exchange-api/exchange"
	"strings"
)

type orderbooktracker struct {
	db   OrderDatabase
	hash string //name of hash where values will be stored
}

func (t *orderbooktracker) Update(sym string, Bids []exchange.MarketOrder, Asks []exchange.MarketOrder) {
	t.db.Update(sym, Bids, Asks)
}

// Get gets information about stock
func (t *orderbooktracker) Get(sym string) (*exchange.MarketRecord, error) {
	return t.db.Get(sym)
}

// NewOrderbookTracker returns exchange.OrderbookTracker that wraps redis connection
func NewOrderbookTracker(dbType, dbUrl, hash string) (exchange.Orderbooks, error) {
	db, err := NewDatabase(dbType, dbUrl, hash)

	if err != nil {
		return nil, err
	}

	return &orderbooktracker{
		db:   db,
		hash: hash,
	}, nil
}

func normalize(symbol string) string {
	return strings.ToUpper(strings.Replace(symbol, "_", "/", -1))
}
