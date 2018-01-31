package cli

import "github.com/skycoin/exchange-api/exchange"
import "encoding/json"

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
	str, _ := json.MarshalIndent(r, "", "    ")
	return string(str)
}

func orderbookFull(orderbook exchange.MarketRecord) string {
	str, _ := json.MarshalIndent(orderbook, "", "    ")
	return string(str)
}

func orderbookShort(orderbook exchange.MarketRecord) string {
	averageBuyPrice := 0.0
	averageSellPrice := 0.0
	totalBuyVolume := 0.0
	totalSellVolume := 0.0
	for _, v := range orderbook.Bids {
		if v.Price == 0 {
			continue
		}
		totalBuyVolume += v.Volume

	}
	for _, v := range orderbook.Bids {
		if v.Price == 0 {
			continue
		}
		averageBuyPrice += v.Price * (v.Volume / totalBuyVolume)
	}
	for _, v := range orderbook.Asks {
		if v.Price == 0 {
			continue
		}
		totalSellVolume += v.Volume

	}
	for _, v := range orderbook.Asks {

		averageSellPrice += v.Price * (v.Volume / totalSellVolume)
	}
	representation := map[string]interface{}{
		"timestamp":          orderbook.Timestamp,
		"symbol":             orderbook.Symbol,
		"average_sell_price": averageSellPrice,
		"average_buy_price":  averageBuyPrice,
		"total_sell_volume":  totalSellVolume,
		"total_buy_volume":   totalBuyVolume,
	}
	str, _ := json.MarshalIndent(representation, "", "    ")
	return string(str)
}
