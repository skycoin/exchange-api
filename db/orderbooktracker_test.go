package db

import (
	"os"
	"testing"
	"time"

	"reflect"

	exchange "github.com/skycoin/exchange-api/exchange"
)

var redisAddr, testingRedis = os.LookupEnv("REDIS_TEST_ADDR")

func TestRecord_MarshalJSON_UnmarshalJSON(t *testing.T) {
	var r = exchange.MarketRecord{
		Timestamp: time.Unix(1499202345, 0),
		Symbol:    "BTC/LTC",
		Asks:      []exchange.MarketOrder{{Price: 1, Volume: 1}},
		Bids:      []exchange.MarketOrder{{Price: 1, Volume: 1}},
	}
	data, err := r.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var result exchange.MarketRecord
	err = result.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(r, result) {
		t.Fatal("results isnt equal")
	}
}

func Test_orderbooktracker_UpdateSym(t *testing.T) {
	testCases := []struct {
		dbType string
		dbURL  string
		hash   string
	}{
		{
			RedisDatabase,
			redisAddr,
			"test",
		},
		{
			MemoryDatabase,
			"",
			"test",
		},
	}

	for _, test := range testCases {
		if test.dbType != RedisDatabase || testingRedis {
			orderBookTracker, err := NewOrderbookTracker(test.dbType,
				test.dbURL,
				test.hash)

			if err != nil {
				t.Fatal(err)
			}

			var r = exchange.MarketRecord{
				Timestamp: time.Now(),
				Symbol:    "BTC/LTC",
				Bids:      []exchange.MarketOrder{{Price: 1, Volume: 1}},
				Asks:      []exchange.MarketOrder{{Price: 1, Volume: 1}},
			}
			orderBookTracker.Update(
				"BTC/LTC",
				[]exchange.MarketOrder{{Price: 1, Volume: 1}},
				[]exchange.MarketOrder{{Price: 1, Volume: 1}},
			)
			rec, err := orderBookTracker.Get("BTC/LTC")
			if err != nil {
				t.Fatal(err)
			}
			if rec.Symbol != r.Symbol || !reflect.DeepEqual(r.Asks, rec.Asks) || !reflect.DeepEqual(r.Bids, rec.Bids) {
				t.Fatal("records isnt equals")
			}
		}
	}
}
