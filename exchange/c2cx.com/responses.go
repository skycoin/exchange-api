package c2cx

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gz-c/tradebot/exchange"
	"github.com/shopspring/decimal"
)

// newOrder represents an response from CreateOrder function
type newOrder struct {
	OrderID int `json:"orderId"`
}

// Balance is a map with strings of balances
// all keys must be lowercase
type Balance map[string]decimal.Decimal

// ForCoin returns the balance for a given coin symbol
func (b Balance) ForCoin(coinType string) decimal.Decimal {
	c, ok := b[strings.ToLower(coinType)]
	if !ok {
		return decimal.Zero
	}

	return c
}

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

// Balances returns the available balances of the account on the exchange
func (br balanceResponse) Balances() Balance {
	// Subtract "frozen" amounts from the listed balances to get the true spendable amounts
	res := make(Balance)
	res["btc"] = br.Balance.Btc.Sub(br.Frozen.Btc)
	res["etc"] = br.Balance.Etc.Sub(br.Frozen.Etc)
	res["eth"] = br.Balance.Eth.Sub(br.Frozen.Eth)
	res["cny"] = br.Balance.Cny.Sub(br.Frozen.Cny)
	res["sky"] = br.Balance.Sky.Sub(br.Frozen.Sky)
	return res
}

// UnmarshalJSON implements json.Unmarshaler
func (b *Balance) UnmarshalJSON(d []byte) error {
	var v balanceResponse
	err := json.Unmarshal(d, &v)
	if err != nil {
		return err
	}
	(*b) = v.Balances()
	return nil
}

type orderJSON struct {
	Amount          decimal.Decimal `json:"amount"`
	AvgPrice        decimal.Decimal `json:"avgPrice"`
	CompletedAmount decimal.Decimal `json:"completedAmount"`
	Fee             decimal.Decimal `json:"fee"`
	CreateDate      int64           `json:"createDate"`
	CompleteDate    int64           `json:"completeDate"`
	OrderID         int             `json:"orderId"`
	Price           decimal.Decimal `json:"price"`
	Status          int             `json:"status"`
	Type            string          `json:"type"`
}

// UnmarshalJSON implements json.Unmarshaler
func (order *Order) UnmarshalJSON(b []byte) error {
	var orderinfo orderJSON
	err := json.Unmarshal(b, &orderinfo)
	if err != nil {
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

	if r.Timestamp, err = strconv.Atoi(v.Timestamp); err != nil {
		return err
	}
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
