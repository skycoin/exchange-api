package c2cx

import (
	"encoding/json"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/skycoin/exchange-api/exchange"
)

type orderJSON struct {
	Amount          decimal.Decimal `json:"amount"`
	AvgPrice        decimal.Decimal `json:"avgPrice"`
	CompletedAmount string  `json:"completedAmount"`
	Fee             decimal.Decimal `json:"fee"`
	CreateDate      int64   `json:"createDate"`
	CompleteDate    int64   `json:"completeDate,omitempty"`
	OrderID         int     `json:"orderId"`
	Price           decimal.Decimal `json:"price"`
	Status          int     `json:"status"`
	Type            string  `json:"type"`
}

// UnmarshalJSON implements json.Unmarshaler
func (order *Order) UnmarshalJSON(b []byte) error {
	var orderinfo orderJSON
	err := json.Unmarshal(b, &orderinfo)
	if err != nil {
		return err
	}
	var completedAmount decimal.Decimal
	if completedAmount, err = decimal.NewFromString(orderinfo.CompletedAmount); err != nil {
		return err
	}
	*order = Order{
		OrderID:         orderinfo.OrderID,
		Status:          orderinfo.Status,
		Amount:          orderinfo.Amount,
		Price:           orderinfo.Price,
		AvgPrice:        orderinfo.AvgPrice,
		Type:            orderinfo.Type,
		CompletedAmount: completedAmount,
		Fee:             orderinfo.Fee,
		CreateDate:      orderinfo.CreateDate,
		CompleteDate:    orderinfo.CompleteDate,
	}
	return nil
}

type orderbookJSON struct {
	Timestamp string          `json:"timestamp"`
	Bids      json.RawMessage `json:"bids"`
	Asks      json.RawMessage `json:"asks"`
}

// UnmarshalJSON implements json.Unmarshaler
func (r *Orderbook) UnmarshalJSON(b []byte) error {
	var v orderbookJSON
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	r.Timestamp, err = strconv.Atoi(v.Timestamp)
	var vals = make([][2]decimal.Decimal, 0)
	err = json.Unmarshal(v.Bids, &vals)
	if err != nil {
		return err
	}
	r.Bids = make([]exchange.MarketOrder, len(vals))
	for i := 0; i < len(vals); i++ {
		r.Bids[i] = exchange.MarketOrder{
			Price:  vals[i][0],
			Volume: vals[i][1],
		}
	}
	err = json.Unmarshal(v.Asks, &vals)
	if err != nil {
		return err
	}
	r.Asks = make([]exchange.MarketOrder, len(vals))
	for i := 0; i < len(vals); i++ {
		r.Asks[i] = exchange.MarketOrder{
			Price:  vals[i][0],
			Volume: vals[i][1],
		}
	}
	return nil
}
