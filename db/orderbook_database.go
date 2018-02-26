package db

import (
	"github.com/skycoin/exchange-api/exchange"
)

// OrderDatabase interface for manipulating with orders
type OrderDatabase interface {
	Get(string) (*exchange.MarketRecord, error)
	Update(string, []exchange.MarketOrder, []exchange.MarketOrder)
}

// NewDatabase factory method for creating order database
func NewDatabase() (OrderDatabase, error) {
	return &memoryDb{}, nil
}
