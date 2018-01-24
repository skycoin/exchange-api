package c2cx

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	exchange "github.com/skycoin/exchange-api/exchange"
)

func Test_sign(t *testing.T) {
	var params = url.Values{}
	params.Add("apiKey", "C821DB84-6FBD-11E4-A9E3-C86000D26D7C")
	want := "BC0DE7EBA50C730BDFC575FE2CD54082"
	expected := sign("12D857DE-7A92-F555-10AC-7566A0D84D1B", params)
	if want != expected {
		t.Fatalf("Incorrect sign!\nwant %s, expected %s", want, expected)
	}
}

func Test_convert(t *testing.T) {
	var createDate = int64(1500757124000)
	var accepted = time.Unix(int64(createDate/1000),
		int64(createDate%1000*int64(time.Millisecond)))
	var tests = []struct {
		in  Order
		out exchange.Order
	}{
		{
			in: Order{
				OrderID:         1,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 0,

				Status:     statuses[exchange.Opened],
				Type:       exchange.Buy,
				CreateDate: createDate,
			},
			out: exchange.Order{
				OrderID:         1,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 0,

				Status:   exchange.Opened,
				Accepted: accepted,
				Type:     exchange.Buy,
			},
		},
		{
			in: Order{
				OrderID:         2,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 0.5,

				Status:     statuses[exchange.Partial],
				Type:       exchange.Buy,
				CreateDate: createDate,
			},
			out: exchange.Order{
				OrderID:         2,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 0.5,

				Status:   exchange.Partial,
				Accepted: accepted,
				Type:     exchange.Buy,
			},
		},
		{
			in: Order{
				OrderID:         3,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 1.0,

				Status:     statuses[exchange.Completed],
				Type:       exchange.Buy,
				CreateDate: createDate,
			},
			out: exchange.Order{
				OrderID:         3,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 1.0,

				Status:   exchange.Completed,
				Accepted: accepted,
				Type:     exchange.Buy,
			},
		},
		{
			in: Order{
				OrderID:         4,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 0.7,

				Status:     statuses[exchange.Cancelled],
				Type:       exchange.Buy,
				CreateDate: createDate,
			},
			out: exchange.Order{
				OrderID:         4,
				Amount:          1.0,
				Price:           1.0,
				Fee:             0.01,
				CompletedAmount: 0.7,

				Status:   exchange.Cancelled,
				Type:     exchange.Buy,
				Accepted: accepted,
			},
		},
	}
	for i, v := range tests {
		if !reflect.DeepEqual(convert(v.in), v.out) {
			t.Fatalf("test %d falied, in %#v, out %#v", i, convert(v.in), v.out)
		}
	}
}
