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
	Btc float64 `json:"btc"`
	Etc float64 `json:"etc"`
	Eth float64 `json:"eth"`
	Cny float64 `json:"cny"`
	Sky float64 `json:"sky"`
}

type balanceResponse struct {
	Balance balanceResponseEntry `json:"balance"`
	Frozen  balanceResponseEntry `json:"frozen"`
}

func subFloatsToDecimal(a, b float64) decimal.Decimal {
	return decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b))
}

func (br balanceResponse) Balances() map[string]decimal.Decimal {
	res := make(map[string]decimal.Decimal)
	res["btc"] = subFloatsToDecimal(br.Balance.Btc, br.Frozen.Btc)
	res["etc"] = subFloatsToDecimal(br.Balance.Etc, br.Frozen.Etc)
	res["eth"] = subFloatsToDecimal(br.Balance.Eth, br.Frozen.Eth)
	res["cny"] = subFloatsToDecimal(br.Balance.Cny, br.Frozen.Cny)
	res["sky"] = subFloatsToDecimal(br.Balance.Sky, br.Frozen.Sky)
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
	(*r)["btc"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Btc, v.Frozen.Btc)
	(*r)["etc"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Etc, v.Frozen.Etc)
	(*r)["eth"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Eth, v.Frozen.Eth)
	(*r)["sky"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Sky, v.Frozen.Sky)
	(*r)["cny"] = fmt.Sprintf("Availible %.8f, frozen %.8f", v.Balance.Cny, v.Frozen.Cny)
	return nil
}
