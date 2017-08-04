package c2cx

import (
	"encoding/json"
	"fmt"
)

// newOrder represents an response from CreateOrder function
type newOrder struct {
	OrderID int `json:"orderId"`
}

// Balance is a map with strings of balances
// all keys must be lowercase
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
			Sky float64 `json:"sky"`
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
