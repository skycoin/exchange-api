package cli

import "github.com/uberfurrer/tradebot/exchange"
import "encoding/json"

func orderShort(order exchange.Order) string {
	var representaton = map[string]interface{}{
		"orderid": order.OrderID,
		"market":  order.Market,
		"price":   order.Price,
		"amount":  order.Amount,
	}
	str, _ := json.MarshalIndent(representaton, "", "    ")
	return string(str)
}

func orderFull(order exchange.Order) string {
	str, _ := json.MarshalIndent(order, "", "    ")
	return string(str)
}

func orderbookFull(orderbook exchange.MarketRecord) string {
	str, _ := json.MarshalIndent(orderbook, "", "    ")
	return string(str)
}

func orderbookShort(orderbook exchange.MarketRecord) string {
	var (
		averageBuyPrice  float64
		averageSellPrice float64
		totalBuyVolume   float64
		totalSellVolume  float64
	)
	for _, v := range orderbook.Bids {
		totalBuyVolume += v.Volume
	}
	for _, v := range orderbook.Bids {
		averageBuyPrice += v.Price * (v.Volume / totalBuyVolume)
	}
	for _, v := range orderbook.Asks {
		totalSellVolume += v.Volume
	}
	for _, v := range orderbook.Asks {
		averageSellPrice += v.Price * (v.Volume / totalSellVolume)
	}
	var representation = map[string]interface{}{
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
