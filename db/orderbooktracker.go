package db

import (
	"time"

	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/uberfurrer/tradebot/exchange"
	"github.com/uberfurrer/tradebot/logger"
)

type orderbooktracker struct {
	db   *redis.Client
	hash string //name of hash where values will be stored
}

func (t *orderbooktracker) UpdateSym(sym string, Bids []exchange.OrderbookEntry, Asks []exchange.OrderbookEntry) {
	book := exchange.Orderbook{
		Symbol:    sym,
		Timestamp: time.Now(),
		Asks:      Asks,
		Bids:      Bids,
	}
	data, err := json.Marshal(book)
	if err != nil {
		logger.Warningf("orderbooktracker: falied to store update: error marshaling symbol %s %s", sym, err.Error())
		return
	}
	cmd := t.db.HSet(t.hash, sym, data)
	if cmd.Err() != nil {
		logger.Warningf("orderbooktracker: falied to store update: redis error %s", cmd.Err().Error())
	}
	return

}

// GetRecord gets information about stock
func (t *orderbooktracker) GetRecord(sym string) (exchange.Orderbook, error) {
	var (
		r   exchange.Orderbook
		err error
	)
	result := t.db.HGet(t.hash, sym)
	if err = result.Err(); err != nil {
		logger.Warningf("orderbooktracker: failed to get data from redis for symbol %s, redis: %s", sym, err)
		return r, err
	}
	if bb, err := result.Bytes(); err == nil {
		err = json.Unmarshal(bb, &r)
	}
	return r, err
}

// NewOrderbookTracker returns exchange.OrderbookTracker that wraps redis connection
func NewOrderbookTracker(opts *redis.Options, hash string) exchange.OrderBookTracker {
	return &orderbooktracker{
		db:   redis.NewClient(opts),
		hash: hash,
	}
}
