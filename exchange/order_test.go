package exchange

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrder_jsonMarshaler(t *testing.T) {
	testtime := time.Now().Truncate(time.Millisecond)
	var order = Order{
		OrderID:         1,
		Status:          Completed,
		Type:            Buy,
		Market:          "BTC/LTC",
		Price:           decimal.NewFromFloat(2250.01),
		Amount:          decimal.NewFromFloat(10.0),
		CompletedAmount: decimal.NewFromFloat(2250.0),
		Fee:             decimal.NewFromFloat(0.1),
		Accepted:        testtime,
		Submitted:       testtime,
		Completed:       testtime,
	}
	b, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}
	var result Order
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(order, result) {
		t.Errorf("want %v, expected %v", order, result)
	}
}
