package db

import (
	"time"

	"encoding/json"

	"strings"

	"github.com/go-redis/redis"

	"github.com/skycoin/exchange-api/exchange"
)

type orderbooktracker struct {
	db   *redis.Client
	hash string //name of hash where values will be stored
}

func (t *orderbooktracker) Update(sym string, Bids []exchange.MarketOrder, Asks []exchange.MarketOrder) {
	book := exchange.MarketRecord{
		Symbol:    normalize(sym),
		Timestamp: time.Now(),
		Asks:      Asks,
		Bids:      Bids,
	}
	data, err := json.Marshal(book)
	if err != nil {
		return
	}
	t.db.HSet(t.hash, normalize(sym), data)
	return

}

// Get gets information about stock
func (t *orderbooktracker) Get(sym string) (*exchange.MarketRecord, error) {
	result := t.db.HGet(t.hash, normalize(sym))
	if err := result.Err(); err != nil {
		return nil, err
	}

	bb, err := result.Bytes()
	if err != nil {
		return nil, err
	}

	var r exchange.MarketRecord
	if err := json.Unmarshal(bb, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// NewOrderbookTracker returns exchange.OrderbookTracker that wraps redis connection
func NewOrderbookTracker(opts *redis.Options, hash string) exchange.Orderbooks {
	return &orderbooktracker{
		db:   redis.NewClient(opts),
		hash: hash,
	}
}

func normalize(symbol string) string {
	return strings.ToUpper(strings.Replace(symbol, "_", "/", -1))
}
