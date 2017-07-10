package db

import (
	"testing"
	"time"

	"reflect"

	"github.com/go-redis/redis"
	exchange "github.com/uberfurrer/tradebot/exchange"
)

func TestRecord_MarshalJSON_UnmarshalJSON(t *testing.T) {
	var r = exchange.Orderbook{
		Timestamp: time.Unix(1499202345, 0),
		Symbol:    "BTC/LTC",
		Asks:      []exchange.OrderbookEntry{exchange.OrderbookEntry{Price: 1, Volume: 1}},
		Bids:      []exchange.OrderbookEntry{exchange.OrderbookEntry{Price: 1, Volume: 1}},
	}
	data, err := r.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var result exchange.Orderbook
	err = result.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(r, result) {
		t.Fatal("results isnt equal")
	}
}
func Test_orderbooktracker_UpdateSym(t *testing.T) {
	var tr = orderbooktracker{
		db: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		hash:   "test",
		buffer: make(map[string]exchange.Orderbook),
	}
	var r = exchange.Orderbook{
		Timestamp: time.Now(),
		Symbol:    "BTC/LTC",
		Bids:      []exchange.OrderbookEntry{exchange.OrderbookEntry{Price: 1, Volume: 1}},
		Asks:      []exchange.OrderbookEntry{exchange.OrderbookEntry{Price: 1, Volume: 1}},
	}
	tr.UpdateSym(
		"BTC/LTC",
		[]exchange.OrderbookEntry{exchange.OrderbookEntry{Price: 1, Volume: 1}},
		[]exchange.OrderbookEntry{exchange.OrderbookEntry{Price: 1, Volume: 1}},
	)
	rec, err := tr.GetRecord("BTC/LTC")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Symbol != r.Symbol || !reflect.DeepEqual(r.Asks, rec.Asks) || !reflect.DeepEqual(r.Bids, rec.Bids) {
		t.Fatal("records isnt equals")
	}

}
