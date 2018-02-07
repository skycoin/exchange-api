package db

import (
	"fmt"

	"github.com/go-redis/redis"

	"github.com/skycoin/exchange-api/exchange"
)

// Order book database types
const (
	RedisDatabase  = "redis"
	MemoryDatabase = "memory"
)

// OrderDatabase interface for manipulating with orders
type OrderDatabase interface {
	Get(string) (*exchange.MarketRecord, error)
	Update(string, []exchange.MarketOrder, []exchange.MarketOrder)
}

// NewDatabase factory method for creating order database
func NewDatabase(dbType, dbURL, hash string) (OrderDatabase, error) {
	switch dbType {
	case RedisDatabase:
		storage := redis.NewClient(&redis.Options{
			Addr: dbURL,
		})

		cmd := storage.Ping()

		if cmd.Err() != nil {
			return nil, cmd.Err()
		}

		return &redisDb{
			storage,
			hash,
		}, nil
	case MemoryDatabase:
		return &memoryDb{}, nil
	}

	return nil, fmt.Errorf("unknown db type %s", dbType)
}
