package c2cx

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// newOrder represents an response from CreateOrder function
type newOrder struct {
	OrderID int `json:"orderId"`
}

// Balance is a map with strings of balances
// all keys must be lowercase
type Balance map[string]string

type balanceResponseEntry struct {
	Btc decimal.Decimal `json:"btc"`
	Etc decimal.Decimal `json:"etc"`
	Eth decimal.Decimal `json:"eth"`
	Cny decimal.Decimal `json:"cny"`
	Sky decimal.Decimal `json:"sky"`
}

type balanceResponse struct {
	Balance balanceResponseEntry `json:"balance"`
	Frozen  balanceResponseEntry `json:"frozen"`
}

func (br balanceResponse) Balances() map[string]decimal.Decimal {
	res := make(map[string]decimal.Decimal)
	res["btc"] = br.Balance.Btc.Sub(br.Frozen.Btc)
	res["etc"] = br.Balance.Etc.Sub(br.Frozen.Etc)
	res["eth"] = br.Balance.Eth.Sub(br.Frozen.Eth)
	res["cny"] = br.Balance.Cny.Sub(br.Frozen.Cny)
	res["sky"] = br.Balance.Sky.Sub(br.Frozen.Sky)
	return res
}

// UnmarshalJSON implements json.Unmarshaler
func (r *Balance) UnmarshalJSON(b []byte) error {
	if *r == nil {
		(*r) = make(map[string]string)
	}
	var v balanceResponse
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	(*r)["btc"] = fmt.Sprintf("Available %.8f, frozen %.8f", v.Balance.Btc, v.Frozen.Btc)
	(*r)["etc"] = fmt.Sprintf("Available %.8f, frozen %.8f", v.Balance.Etc, v.Frozen.Etc)
	(*r)["eth"] = fmt.Sprintf("Available %.8f, frozen %.8f", v.Balance.Eth, v.Frozen.Eth)
	(*r)["sky"] = fmt.Sprintf("Available %.8f, frozen %.8f", v.Balance.Sky, v.Frozen.Sky)
	(*r)["cny"] = fmt.Sprintf("Available %.8f, frozen %.8f", v.Balance.Cny, v.Frozen.Cny)
	return nil
}
