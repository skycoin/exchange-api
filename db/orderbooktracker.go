package db

import (
	"time"

	"encoding/json"

	"strings"

	"github.com/go-redis/redis"
	"github.com/uberfurrer/tradebot/exchange"
	"github.com/uberfurrer/tradebot/logger"
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
		logger.Warningf("orderbooktracker: falied to store update: error marshaling symbol %s %s", normalize(sym), err.Error())
		return
	}
	cmd := t.db.HSet(t.hash, normalize(sym), data)
	if cmd.Err() != nil {
		logger.Warningf("orderbooktracker: falied to store update, %s", cmd.Err().Error())
	}
	return

}

// Get gets information about stock
func (t *orderbooktracker) Get(sym string) (exchange.MarketRecord, error) {
	var (
		r   exchange.MarketRecord
		err error
	)
	result := t.db.HGet(t.hash, normalize(sym))
	if err = result.Err(); err != nil {
		logger.Warningf("orderbooktracker: failed to get data from redis for symbol %s, %s", normalize(sym), err)
		return r, err
	}
	if bb, err := result.Bytes(); err == nil {
		err = json.Unmarshal(bb, &r)
	}
	return r, err
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
