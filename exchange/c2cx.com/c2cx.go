// Package c2cx provides api methods methods for communication with c2cx exchange
package c2cx

import (
	"github.com/shopspring/decimal"
)

// Markets is all supported markets
// add new markets here
var Markets = []string{"CNY_BTC", "CNY_ETH", "CNY_ETC", "CNY_SKY", "ETH_SKY", "BTC_SKY", "CNY_SHL", "BTC_BCC"}

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
	OrderID         int
	Price           decimal.Decimal
	Status          int
	Type            string
}
