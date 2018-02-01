package db

import (
	"testing"
	"time"

	"reflect"

	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	exchange "github.com/skycoin/exchange-api/exchange"
)

func TestRecord_MarshalJSON_UnmarshalJSON(t *testing.T) {
	decimalOne := decimal.NewFromFloat(1.0)
	var r = exchange.MarketRecord{
		Timestamp: time.Unix(1499202345, 0),
		Symbol:    "BTC/LTC",
		Asks:      []exchange.MarketOrder{exchange.MarketOrder{Price: decimalOne, Volume: decimalOne}},
		Bids:      []exchange.MarketOrder{exchange.MarketOrder{Price: decimalOne, Volume: decimalOne}},
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
	decimalOne := decimal.NewFromFloat(1.0)
	var tr = orderbooktracker{
		db: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		hash: "test",
	}
	var r = exchange.MarketRecord{
		Timestamp: time.Now(),
		Symbol:    "BTC/LTC",
		Bids:      []exchange.MarketOrder{exchange.MarketOrder{Price: decimalOne, Volume: decimalOne}},
		Asks:      []exchange.MarketOrder{exchange.MarketOrder{Price: decimalOne, Volume: decimalOne}},
	}
	tr.Update(
		"BTC/LTC",
		[]exchange.MarketOrder{exchange.MarketOrder{Price: decimalOne, Volume: decimalOne}},
		[]exchange.MarketOrder{exchange.MarketOrder{Price: decimalOne, Volume: decimalOne}},
	)
	rec, err := tr.Get("BTC/LTC")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Symbol != r.Symbol || !reflect.DeepEqual(r.Asks, rec.Asks) || !reflect.DeepEqual(r.Bids, rec.Bids) {
		t.Fatal("records isnt equals")
	}

}
