// Package c2cx provides api methods methods for communication with c2cx exchange
package c2cx

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/skycoin/exchange-api/exchange"
)

// OrderStatus is an order's status
type OrderStatus string

// OrderType is an order's type
type OrderType string

// PriceType is an order's price type
type PriceType string

// OrderID is an order's ID
type OrderID int

// TradePair is a market trade pair
type TradePair string

// TradePairRules defines variable configuration per trading pair
type TradePairRules struct {
	// PricePrecision is the maximum number of decimals for price
	PricePrecision int
	// VolumePrecision is the maximum number of decimals for volume
	VolumePrecision int
	// VolumeMinimum is the minimum volume value
	VolumeMinimum decimal.Decimal
}

const (
	// StatusAll all orders
	StatusAll OrderStatus = "all"
	// StatusSubmitted submitted orders
	StatusSubmitted OrderStatus = "submitted"
	// StatusOpened open orders
	StatusOpened OrderStatus = "opened"
	// StatusPartial partial orders
	StatusPartial OrderStatus = "partial"
	// StatusCompleted completed orders
	StatusCompleted OrderStatus = "completed"
	// StatusCancelled cancelled orders
	StatusCancelled OrderStatus = "cancelled"
	// StatusSuspended suspended orders
	StatusSuspended OrderStatus = "suspended"
	// StatusErrored errored orders
	StatusErrored OrderStatus = "errored"
	// StatusTriggerPending trigger pending orders
	StatusTriggerPending OrderStatus = "trigger_pending"
	// StatusStopLossPending stop loss pending orders
	StatusStopLossPending OrderStatus = "stop_loss_pending"
	// StatusExpired expired orders
	StatusExpired OrderStatus = "expired"
	// StatusCancelling cancelling orders
	StatusCancelling OrderStatus = "cancelling"

	// OrderTypeBuy is a buy order
	OrderTypeBuy OrderType = "buy"
	// OrderTypeSell is a sell order
	OrderTypeSell OrderType = "sell"

	// PriceTypeLimit a limit order
	PriceTypeLimit PriceType = "limit"
	// PriceTypeMarket a market order
	PriceTypeMarket PriceType = "market"

	// AllOrders is used to include all orders when orderID is required
	AllOrders OrderID = -1

	// CnyBtc trade pair
	CnyBtc TradePair = "CNY_BTC"
	// CnyEth trade pair
	CnyEth TradePair = "CNY_ETH"
	// CnyEtc trade pair
	CnyEtc TradePair = "CNY_ETC"
	// CnySky trade pair
	CnySky TradePair = "CNY_SKY"
	// EthSky trade pair
	EthSky TradePair = "ETH_SKY"
	// BtcSky trade pair
	BtcSky TradePair = "BTC_SKY"
	// CnyShl trade pair
	CnyShl TradePair = "CNY_SHL"
	// BtcBcc trade pair
	BtcBcc TradePair = "BTC_BCC"
)

var (
	// OrderStatuses maps OrderStatus strings to status codes.
	// From the C2CX API Docs:
	// 0. All
	// 1. Pending
	// 2. Active
	// 3. Partially Completed
	// 4. Completed
	// 5. Cancelled
	// 6. Error
	// 7. Suspended
	// 8. TriggerPending
	// 9. StopLossPending
	// 11. Expired
	// 12. Cancelling
	OrderStatuses = map[OrderStatus]int{
		StatusAll:             0,
		StatusSubmitted:       1,
		StatusOpened:          2,
		StatusPartial:         3,
		StatusCompleted:       4,
		StatusCancelled:       5,
		StatusErrored:         6,
		StatusSuspended:       7,
		StatusTriggerPending:  8,
		StatusStopLossPending: 9,
		StatusExpired:         11,
		StatusCancelling:      12,
	}

	// ReverseOrderStatuses maps status ID to name
	ReverseOrderStatuses map[int]OrderStatus

	// ErrUnknownStatus is returned for an unknown status
	ErrUnknownStatus = errors.New("unknown status")

	// TradePairRulesTable maps TradePairs to their TradePairRules
	TradePairRulesTable = map[TradePair]TradePairRules{
		BtcSky: {
			PricePrecision:  5,
			VolumePrecision: 2,
			VolumeMinimum:   decimal.New(1, 0),
		},
	}
)

func init() {
	ReverseOrderStatuses = make(map[int]OrderStatus)
	for k, v := range OrderStatuses {
		ReverseOrderStatuses[v] = k
	}
}

// AdvancedOrderParams is extended parameters, that can be used for set stoploss, takeprofit and trigger price
type AdvancedOrderParams struct {
	TakeProfit   decimal.Decimal `json:"take_profit"`
	StopLoss     decimal.Decimal `json:"stop_loss"`
	TriggerPrice decimal.Decimal `json:"trigger_price"`
}

// Order represents all information about order
type Order struct {
	Amount          decimal.Decimal
	AvgPrice        decimal.Decimal
	CompletedAmount decimal.Decimal
	Fee             decimal.Decimal
	CreateDate      int64
	CompleteDate    int64
	OrderID         OrderID
	Price           decimal.Decimal
	Status          OrderStatus
	Type            OrderType
}

type orderJSON struct {
	Amount          decimal.Decimal `json:"amount"`
	AvgPrice        decimal.Decimal `json:"avgPrice"`
	CompletedAmount decimal.Decimal `json:"completedAmount"`
	Fee             decimal.Decimal `json:"fee"`
	CreateDate      int64           `json:"createDate"`
	CompleteDate    int64           `json:"completeDate"`
	OrderID         OrderID         `json:"orderId"`
	Price           decimal.Decimal `json:"price"`
	Status          int             `json:"status"`
	Type            OrderType       `json:"type"`
}

// UnmarshalJSON implements json.Unmarshaler
func (order *Order) UnmarshalJSON(b []byte) error {
	var orderinfo orderJSON
	err := json.Unmarshal(b, &orderinfo)
	if err != nil {
		return err
	}

	status, ok := ReverseOrderStatuses[orderinfo.Status]
	if !ok {
		return ErrUnknownStatus
	}

	*order = Order{
		OrderID:         orderinfo.OrderID,
		Status:          status,
		Amount:          orderinfo.Amount,
		Price:           orderinfo.Price,
		AvgPrice:        orderinfo.AvgPrice,
		Type:            orderinfo.Type,
		CompletedAmount: orderinfo.CompletedAmount,
		Fee:             orderinfo.Fee,
		CreateDate:      orderinfo.CreateDate,
		CompleteDate:    orderinfo.CompleteDate,
	}
	return nil
}

// Orderbook with timestamp
type Orderbook struct {
	Timestamp int                   `json:"timestamp"`
	Bids      exchange.MarketOrders `json:"bids"`
	Asks      exchange.MarketOrders `json:"asks"`
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
	r.Bids = make(exchange.MarketOrders, len(vals))
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
	r.Asks = make(exchange.MarketOrders, len(vals))
	for i := 0; i < len(vals); i++ {
		r.Asks[i] = exchange.MarketOrder{
			Price:  vals[i][0],
			Volume: vals[i][1],
		}
	}
	return nil
}

// newOrder represents an response from CreateOrder function
type newOrder struct {
	OrderID OrderID `json:"orderId"`
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
