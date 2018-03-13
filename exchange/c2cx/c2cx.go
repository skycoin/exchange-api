// Package c2cx provides api methods methods for communication with c2cx exchange
package c2cx

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/skycoin/exchange-api/exchange"
)

// OrderStatus is an order's status
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
type OrderStatus int

// Name returns an OrderStatus's human-readable name
func (s OrderStatus) Name() string {
	switch s {
	case StatusAll:
		return "all"
	case StatusSubmitted:
		return "submitted"
	case StatusActive:
		return "opened"
	case StatusPartial:
		return "partially_completed"
	case StatusCompleted:
		return "completed"
	case StatusCancelled:
		return "cancelled"
	case StatusErrored:
		return "errored"
	case StatusSuspended:
		return "suspended"
	case StatusTriggerPending:
		return "trigger_pending"
	case StatusStopLossPending:
		return "stop_loss_pending"
	case StatusExpired:
		return "expired"
	case StatusCancelling:
		return "cancelling"
	default:
		return "unknown"
	}
}

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
	StatusAll OrderStatus = 0
	// StatusSubmitted submitted orders
	StatusSubmitted OrderStatus = 1
	// StatusActive open orders
	StatusActive OrderStatus = 2
	// StatusPartial partial orders
	StatusPartial OrderStatus = 3
	// StatusCompleted completed orders
	StatusCompleted OrderStatus = 4
	// StatusCancelled cancelled orders
	StatusCancelled OrderStatus = 5
	// StatusErrored errored orders
	StatusErrored OrderStatus = 6
	// StatusSuspended suspended orders
	StatusSuspended OrderStatus = 7
	// StatusTriggerPending trigger pending orders
	StatusTriggerPending OrderStatus = 8
	// StatusStopLossPending stop loss pending orders
	StatusStopLossPending OrderStatus = 9
	// StatusExpired expired orders
	StatusExpired OrderStatus = 11
	// StatusCancelling cancelling orders
	StatusCancelling OrderStatus = 12

	// OrderTypeBuy is a buy order
	OrderTypeBuy OrderType = "buy"
	// OrderTypeSell is a sell order
	OrderTypeSell OrderType = "sell"

	// PriceTypeLimit a limit order
	PriceTypeLimit PriceType = "limit"
	// PriceTypeMarket a market order
	PriceTypeMarket PriceType = "market"

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

	// allOrders is used to include all orders when orderID is required
	allOrders string = "-1"

	// pagination parameters
	// minPageSize = 1
	maxPageSize = 100
)

var (
	// TradePairRulesTable maps TradePairs to their TradePairRules
	TradePairRulesTable = map[TradePair]TradePairRules{
		BtcSky: {
			PricePrecision:  5,
			VolumePrecision: 2,
			VolumeMinimum:   decimal.New(1, 0),
		},
	}
)

// AdvancedOrderParams is extended parameters, that can be used for set stoploss, takeprofit and trigger price
type AdvancedOrderParams struct {
	TakeProfit   *decimal.Decimal
	StopLoss     *decimal.Decimal
	TriggerPrice *decimal.Decimal
}

// Order represents all information about order
type Order struct {
	Amount          decimal.Decimal
	AvgPrice        decimal.Decimal
	CompletedAmount decimal.Decimal
	Fee             decimal.Decimal
	CreateDate      time.Time
	CompleteDate    time.Time
	OrderID         OrderID
	Price           decimal.Decimal
	Status          OrderStatus
	Type            OrderType
	Trigger         *decimal.Decimal
	CID             *int
	Source          string
}

type orderJSON struct {
	Amount          decimal.Decimal  `json:"amount"`
	AvgPrice        decimal.Decimal  `json:"avgPrice"`
	CompletedAmount decimal.Decimal  `json:"completedAmount"`
	CreateDate      int64            `json:"createDate"`
	CompleteDate    int64            `json:"completeDate"`
	OrderID         OrderID          `json:"orderId"`
	Price           decimal.Decimal  `json:"price"`
	Status          OrderStatus      `json:"status"`
	Fee             decimal.Decimal  `json:"fee"`
	Type            OrderType        `json:"type"`
	Trigger         *decimal.Decimal `json:"trigger"`
	CID             *int             `json:"cid"`
	Source          string           `json:"source"`
}

func fromUnixMilli(t int64) time.Time {
	base := t / 1e3
	nano := (t % 1e3) * 1e6
	return time.Unix(base, nano)
}

// UnmarshalJSON unmarshals binary data to Order
func (o *Order) UnmarshalJSON(b []byte) error {
	var v orderJSON
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	createDate := fromUnixMilli(v.CreateDate)
	completeDate := fromUnixMilli(v.CompleteDate)

	*o = Order{
		Amount:          v.Amount,
		AvgPrice:        v.AvgPrice,
		CompletedAmount: v.CompletedAmount,
		Fee:             v.Fee,
		CreateDate:      createDate,
		CompleteDate:    completeDate,
		OrderID:         v.OrderID,
		Price:           v.Price,
		Status:          v.Status,
		Type:            v.Type,
		Trigger:         v.Trigger,
		CID:             v.CID,
		Source:          v.Source,
	}

	return nil
}

// Orderbook with timestamp
type Orderbook struct {
	TradePair TradePair
	Timestamp time.Time
	Bids      exchange.MarketOrders
	Asks      exchange.MarketOrders
}

type orderbookJSON struct {
	Timestamp string               `json:"timestamp"`
	Bids      [][2]decimal.Decimal `json:"bids"`
	Asks      [][2]decimal.Decimal `json:"asks"`
}

// UnmarshalJSON implements json.Unmarshaler
func (r *Orderbook) UnmarshalJSON(b []byte) error {
	var v orderbookJSON
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	ts, err := strconv.ParseInt(v.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	r.Timestamp = time.Unix(ts, 0)

	r.Bids = make(exchange.MarketOrders, len(v.Bids))
	for i := 0; i < len(v.Bids); i++ {
		r.Bids[i] = exchange.MarketOrder{
			Price:  v.Bids[i][0],
			Volume: v.Bids[i][1],
		}
	}

	r.Asks = make(exchange.MarketOrders, len(v.Asks))
	for i := 0; i < len(v.Asks); i++ {
		r.Asks[i] = exchange.MarketOrder{
			Price:  v.Asks[i][0],
			Volume: v.Asks[i][1],
		}
	}

	return nil
}

// Balances represents balances held on an exchange
type Balances struct {
	Btc   decimal.Decimal `json:"btc"`
	Etc   decimal.Decimal `json:"etc"`
	Eth   decimal.Decimal `json:"eth"`
	Cny   decimal.Decimal `json:"cny"`
	Sky   decimal.Decimal `json:"sky"`
	Ltc   decimal.Decimal `json:"ltc"`
	Bcc   decimal.Decimal `json:"bcc"`
	Shl   decimal.Decimal `json:"shl"`
	Bch   decimal.Decimal `json:"bch"`
	Zec   decimal.Decimal `json:"zec"`
	Drg   decimal.Decimal `json:"drg"`
	Usdt  decimal.Decimal `json:"usdt"`
	Btg   decimal.Decimal `json:"btg"`
	Fcabs decimal.Decimal `json:"fcabs"`
	Cabs  decimal.Decimal `json:"cabs"`
	Dash  decimal.Decimal `json:"dash"`
	Zrx   decimal.Decimal `json:"zrx"`
	Fun   decimal.Decimal `json:"fun"`
	Tnb   decimal.Decimal `json:"tnb"`
	Etp   decimal.Decimal `json:"etp"`
	Ucash decimal.Decimal `json:"ucash"`
	Total decimal.Decimal `json:"total"`
}

// BalanceSummary includes the account balance and its frozen balance
type BalanceSummary struct {
	Balance Balances `json:"balance"`
	Frozen  Balances `json:"frozen"`
}

// Spendable returns the available balances of the account on the exchange by subtracting
// frozen amounts from the total amounts
func (br BalanceSummary) Spendable() Balances {
	return Balances{
		Btc:   br.Balance.Btc.Sub(br.Frozen.Btc),
		Etc:   br.Balance.Etc.Sub(br.Frozen.Etc),
		Eth:   br.Balance.Eth.Sub(br.Frozen.Eth),
		Cny:   br.Balance.Cny.Sub(br.Frozen.Cny),
		Sky:   br.Balance.Sky.Sub(br.Frozen.Sky),
		Ltc:   br.Balance.Ltc.Sub(br.Frozen.Ltc),
		Bcc:   br.Balance.Bcc.Sub(br.Frozen.Bcc),
		Shl:   br.Balance.Shl.Sub(br.Frozen.Shl),
		Bch:   br.Balance.Bch.Sub(br.Frozen.Bch),
		Zec:   br.Balance.Zec.Sub(br.Frozen.Zec),
		Drg:   br.Balance.Drg.Sub(br.Frozen.Drg),
		Usdt:  br.Balance.Usdt.Sub(br.Frozen.Usdt),
		Btg:   br.Balance.Btg.Sub(br.Frozen.Btg),
		Fcabs: br.Balance.Fcabs.Sub(br.Frozen.Fcabs),
		Cabs:  br.Balance.Cabs.Sub(br.Frozen.Cabs),
		Dash:  br.Balance.Dash.Sub(br.Frozen.Dash),
		Zrx:   br.Balance.Zrx.Sub(br.Frozen.Zrx),
		Fun:   br.Balance.Fun.Sub(br.Frozen.Fun),
		Tnb:   br.Balance.Tnb.Sub(br.Frozen.Tnb),
		Etp:   br.Balance.Etp.Sub(br.Frozen.Etp),
		Ucash: br.Balance.Ucash.Sub(br.Frozen.Ucash),
		Total: br.Balance.Total.Sub(br.Frozen.Total),
	}
}
