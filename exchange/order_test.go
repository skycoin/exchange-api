package exchange

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestOrder_jsonMarshaler(t *testing.T) {
	testtime := time.Now().Truncate(time.Millisecond)
	var order = Order{
		OrderID:         1,
		Status:          Completed,
		Type:            Buy,
		Market:          "BTC/LTC",
		Price:           2250.01,
		Amount:          10,
		CompletedAmount: 2250,
		Fee:             0.1,
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
