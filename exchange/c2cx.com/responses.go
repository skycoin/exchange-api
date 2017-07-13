package c2cx

import (
	"encoding/json"
	"fmt"
	"strconv"

	exchange "github.com/uberfurrer/tradebot/exchange"
)

// CreateOrderResponse represents an response from CreateOrder function
type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

// OrderInfo is a single OrderInfo that returns by GetOrderInfo function
type OrderInfo struct {
	Amount          float64 `json:"amount"`
	AvgPrice        float64 `json:"avgPrice"`
	CompletedAmount string  `json:"completedAmount"`
	CreateDate      int64   `json:"createDate"`
	OrderID         int     `json:"orderId"`
	Price           float64 `json:"price"`
	Status          int     `json:"status"`
	Type            string  `json:"type"`
}

// Orderbook represents a response from GetOrderBook function
type Orderbook struct {
	Timestamp int                       `json:"timestamp"`
	Bids      []exchange.OrderbookEntry `json:"bids"`
	Asks      []exchange.OrderbookEntry `json:"asks"`
}

// UnmarshalJSON implements json.Unmarshaler
func (r *Orderbook) UnmarshalJSON(b []byte) error {
	var v = struct {
		Timestamp string          `json:"timestamp"`
		Bids      json.RawMessage `json:"bids"`
		Asks      json.RawMessage `json:"asks"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	r.Timestamp, err = strconv.Atoi(v.Timestamp)
	var vals = make([][2]float64, 0)
	err = json.Unmarshal(v.Bids, &vals)
	if err != nil {
		return err
	}
	r.Bids = make([]exchange.OrderbookEntry, len(vals))
	for i := 0; i < len(vals); i++ {
		r.Bids[i] = exchange.OrderbookEntry{
			Price:  vals[i][0],
			Volume: vals[i][1],
		}
	}
	err = json.Unmarshal(v.Asks, &vals)
	if err != nil {
		return err
	}
	r.Asks = make([]exchange.OrderbookEntry, len(vals))
	for i := 0; i < len(vals); i++ {
		r.Asks[i] = exchange.OrderbookEntry{
			Price:  vals[i][0],
			Volume: vals[i][1],
		}
	}
	return nil
}

// Balance represents a response from GetUserInfo
// Note: all keys must be lowercase
// Keys: "btc", "etc", "eth", "cny", "sky"
type Balance map[string]string

// UnmarshalJSON implements json.Unmarshaler
func (r *Balance) UnmarshalJSON(b []byte) error {
	if *r == nil {
		(*r) = make(map[string]string)
	}
	var v struct {
		Balance struct {
			Btc float64 `json:"btc"`
			Etc float64 `json:"etc"`
			Eth float64 `json:"eth"`
			Cny float64 `json:"cny"`
			Sky float64 `json:""`
		} `json:"balance"`
		Frozen struct {
			Btc float64 `json:"btc"`
			Etc float64 `json:"etc"`
			Eth float64 `json:"eth"`
			Cny float64 `json:"cny"`
			Sky float64 `json:"sky"`
		} `json:"frozen"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	(*r)["btc"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Btc, v.Frozen.Btc)
	(*r)["etc"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Etc, v.Frozen.Etc)
	(*r)["eth"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Eth, v.Frozen.Eth)
	(*r)["sky"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Sky, v.Frozen.Sky)
	(*r)["cny"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Cny, v.Frozen.Cny)
	return nil
}
