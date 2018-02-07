package db

import (
	"fmt"
	"github.com/skycoin/exchange-api/exchange"
	"sync"
	"time"
)

type MemoryDb struct {
	storage sync.Map
}

func (db *MemoryDb) Get(key string) (*exchange.MarketRecord, error) {
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

func (db *MemoryDb) Update(sym string, Bids []exchange.MarketOrder, Asks []exchange.MarketOrder) {
	book := exchange.MarketRecord{
		Symbol:    normalize(sym),
		Timestamp: time.Now(),
		Asks:      Asks,
		Bids:      Bids,
	}

	db.storage.Store(normalize(sym), book)
}
