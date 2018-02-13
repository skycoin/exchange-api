package cli

import (
	"encoding/json"

	"github.com/shopspring/decimal"
	"github.com/skycoin/exchange-api/exchange"
)

func orderShort(order exchange.Order) string {
	r := map[string]interface{}{
		"orderid": order.OrderID,
		"market":  order.Market,
		"price":   order.Price,
		"amount":  order.Amount,
	}
	str, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func orderFull(order exchange.Order) string {
	r := map[string]interface{}{
		"orderid": order.OrderID,
		"market":  order.Market,
		"type":    order.Type,
		"price":   order.Price,
		"amount":  order.Amount,
		"status":  order.Status,
	}
	switch order.Status {
	case exchange.Submitted:
		r["submitted_at"] = order.Submitted
	case exchange.Opened, exchange.Partial:
		r["submtted_at"] = order.Submitted
		r["accepted_at"] = order.Accepted
		r["completed_amount"] = order.CompletedAmount
	case exchange.Cancelled:
		r["submtted_at"] = order.Submitted
		r["accepted_at"] = order.Accepted
		r["completed_at"] = order.Completed
		r["completed_amount"] = order.CompletedAmount
		r["fee"] = order.Fee
	case exchange.Completed:
		r["submtted_at"] = order.Submitted
		r["accepted_at"] = order.Accepted
		r["completed_at"] = order.Completed
		r["fee"] = order.Fee
	}
	str, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func orderbookFull(orderbook exchange.MarketRecord) string {
	str, err := json.MarshalIndent(orderbook, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func orderbookShort(orderbook exchange.MarketRecord) string {
	var (
		averageBuyPrice  decimal.Decimal
		averageSellPrice decimal.Decimal
		totalBuyVolume   decimal.Decimal
		totalSellVolume  decimal.Decimal
	)
	for _, v := range orderbook.Bids {
		totalBuyVolume = totalBuyVolume.Add(v.Volume)
	}
	for _, v := range orderbook.Bids {
		averageBuyPrice = averageBuyPrice.Add(v.Price.Mul(v.Volume.Div(totalBuyVolume)))
	}
	for _, v := range orderbook.Asks {
		totalSellVolume = totalSellVolume.Add(v.Volume)
	}
	for _, v := range orderbook.Asks {
		averageSellPrice = averageSellPrice.Add(v.Price.Mul(v.Volume.Div(totalSellVolume)))
	}
	representation := map[string]interface{}{
		"timestamp":          orderbook.Timestamp,
		"symbol":             orderbook.Symbol,
		"average_sell_price": averageSellPrice,
		"average_buy_price":  averageBuyPrice,
		"total_sell_volume":  totalSellVolume,
		"total_buy_volume":   totalBuyVolume,
	}
	str, err := json.MarshalIndent(representation, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(str)
}
