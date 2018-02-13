package db

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/skycoin/exchange-api/exchange"
)

type redisDb struct {
	storage *redis.Client
	hash    string
}

func (db *redisDb) Get(key string) (*exchange.MarketRecord, error) {
	result := db.storage.HGet(db.hash, normalize(key))
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

func (db *redisDb) Update(sym string, Bids []exchange.MarketOrder, Asks []exchange.MarketOrder) {
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
	db.storage.HSet(db.hash, normalize(sym), data)
}
