package exchange

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type orderInternal struct {
	OrderID   int    `json:"orderid"`
	Type      string `json:"type"`
	Market    string `json:"market"`
	Submitted int64  `json:"submitted_at"`

	Amount decimal.Decimal `json:"amount"`
	Price  decimal.Decimal `json:"price"`

	Fee             decimal.Decimal `json:"fee"`
	CompletedAmount decimal.Decimal `json:"completed_amount"`

	Status    string `json:"status"`
	Accepted  int64  `json:"accepted_at"`
	Completed int64  `json:"completed_at"`
}

// MarshalJSON implements json.Marshaler interface
func (order Order) MarshalJSON() ([]byte, error) {
	var internal = orderInternal{
		OrderID:   order.OrderID,
		Type:      order.Type,
		Market:    order.Market,
		Amount:    order.Amount,
		Price:     order.Price,
		Submitted: order.Submitted.Truncate(time.Millisecond).UnixNano() / int64(time.Millisecond),

		Fee:             order.Fee,
		CompletedAmount: order.CompletedAmount,
		Status:          order.Status,
		Accepted:        order.Accepted.Truncate(time.Millisecond).UnixNano() / int64(time.Millisecond),
		Completed:       order.Completed.Truncate(time.Millisecond).UnixNano() / int64(time.Millisecond),
	}
	return json.Marshal(internal)
}

//UnmarshalJSON implements json.Unmarshaler interface
func (order *Order) UnmarshalJSON(b []byte) error {
	var internal orderInternal
	err := json.Unmarshal(b, &internal)
	if err != nil {
		return err
	}
	order.OrderID = internal.OrderID
	order.Type = internal.Type
	order.Market = internal.Market
	order.Amount = internal.Amount
	order.Price = internal.Price
	order.Submitted = time.Unix(internal.Submitted/int64(time.Second/time.Millisecond),
		(internal.Submitted%int64(time.Second/time.Millisecond))*int64(time.Millisecond))

	order.Fee = internal.Fee
	order.CompletedAmount = internal.CompletedAmount
	order.Status = internal.Status
	order.Accepted = time.Unix(internal.Accepted/int64(time.Second/time.Millisecond),
		(internal.Accepted%int64(time.Second/time.Millisecond))*int64(time.Millisecond))
	order.Completed = time.Unix(internal.Completed/int64(time.Second/time.Millisecond),
		(internal.Completed%int64(time.Second/time.Millisecond))*int64(time.Millisecond))
	return nil
}

func truncate(order Order) Order {
	order.Submitted = order.Submitted.Truncate(time.Millisecond)
	order.Accepted = order.Accepted.Truncate(time.Millisecond)
	order.Completed = order.Completed.Truncate(time.Millisecond)
	return order
}
