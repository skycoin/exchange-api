package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/skycoin/exchange-api/exchange"
)

type memoryDb struct {
	storage sync.Map
}

func (db *memoryDb) Get(key string) (*exchange.MarketRecord, error) {
	val, ok := db.storage.Load(normalize(key))

	if !ok {
		return nil, fmt.Errorf("key not found %s", key)
	}

	record, ok := val.(*exchange.MarketRecord)

	if !ok {
		return nil, fmt.Errorf("error converting value %v to *exchange.MarketRecord", val)
	}

	return record, nil
}

func (db *memoryDb) Update(sym string, Bids []exchange.MarketOrder, Asks []exchange.MarketOrder) {
	book := &exchange.MarketRecord{
		Symbol:    normalize(sym),
		Timestamp: time.Now(),
		Asks:      Asks,
		Bids:      Bids,
	}

	db.storage.Store(normalize(sym), book)
}
