package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/skycoin/exchange-api/exchange"
)

const (
	REDIS_DATABASE = "redis"
	MEMORY_DATABSE = "memory"
)

type OrderDatabase interface {
	Get(string) (*exchange.MarketRecord, error)
	Update(string, []exchange.MarketOrder, []exchange.MarketOrder)
}

func NewDatabase(dbType, dbUrl, hash string) (OrderDatabase, error) {
	switch dbType {
	case REDIS_DATABASE:
		storage := redis.NewClient(&redis.Options{
			Addr: dbUrl,
		})

		cmd := storage.Ping()

		if cmd.Err() != nil {
			return nil, cmd.Err()
		}

		return &RedisDb{
			storage,
			hash,
		}, nil
	case MEMORY_DATABSE:
		return &MemoryDb{}, nil
	}

	return nil, fmt.Errorf("unknown db type %s", dbType)
}
