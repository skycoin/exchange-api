package rpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-redis/redis"
	"github.com/uberfurrer/tradebot/db"
	exchange "github.com/uberfurrer/tradebot/exchange"
)

var client = new(mockExchange)

func Test_defaultHandler_GetBalance(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"currency\":\"BTC\"}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["GetBalance"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_Cancel(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"orderid\":1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["Cancel"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_CancelAll(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage(nil),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["CancelAll"](req, client)
	fmt.Printf("%s\n", resp.Result)
}

func Test_defaulthandler_CancelMarket(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\"}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["CancelMarket"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_Buy(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\",\"rate\":1.0,\"amount\":1.1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["Buy"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_Sell(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\",\"rate\":1.0,\"amount\":1.1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["Sell"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_OrderDetails(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"orderid\":1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["OrderDetails"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_OrderStatus(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage("{\"orderid\":1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["OrderStatus"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_Completed(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage(nil),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["Completed"](req, client)
	fmt.Printf("%s\n", resp.Result)
}
func Test_defaulthandler_Executed(t *testing.T) {
	var req = Request{
		Params:  json.RawMessage(nil),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["Executed"](req, client)
	fmt.Printf("%s\n", resp.Result)
}

type mockExchange int

func (d *mockExchange) Cancel(orderID int) (*exchange.OrderInfo, error) {
	return &exchange.OrderInfo{
		Type:      "buy",
		Status:    exchange.StatusCompleted,
		TradePair: "BTC/LTC",

		Volume:  1.0,
		Price:   1.0,
		OrderID: orderID,

		Submitted: int64(0),
		Accepted:  int64(0),
		Completed: int64(0),
	}, nil
}
func (d *mockExchange) CancelMarket(sym string) ([]*exchange.OrderInfo, error) {
	return []*exchange.OrderInfo{
		&exchange.OrderInfo{
			Type:      "buy",
			Status:    exchange.StatusCompleted,
			TradePair: sym,

			Volume:  1.0,
			Price:   1.0,
			OrderID: 1,

			Submitted: int64(0),
			Accepted:  int64(0),
			Completed: int64(0),
		},
	}, nil
}
func (d *mockExchange) CancelAll() ([]*exchange.OrderInfo, error) {
	return []*exchange.OrderInfo{
		&exchange.OrderInfo{
			Type:      "buy",
			Status:    exchange.StatusCompleted,
			TradePair: "BTC/LTC",

			Volume:  1.0,
			Price:   1.0,
			OrderID: 1,

			Submitted: int64(0),
			Accepted:  int64(0),
			Completed: int64(0),
		},
	}, nil
}
func (d *mockExchange) Buy(sym string, price, amount float64) (int, error) {
	return 1, nil
}
func (d *mockExchange) Sell(sym string, price, amount float64) (int, error) {
	return 2, nil
}
func (d *mockExchange) Completed() []*exchange.OrderInfo {
	return []*exchange.OrderInfo{
		&exchange.OrderInfo{
			Type:      "buy",
			Status:    exchange.StatusCompleted,
			TradePair: "BTC/LTC",

			Volume:  1.0,
			Price:   1.0,
			OrderID: 1,

			Submitted: int64(0),
			Accepted:  int64(0),
			Completed: int64(0),
		},
	}
}
func (d *mockExchange) Executed() []*exchange.OrderInfo {
	return []*exchange.OrderInfo{
		&exchange.OrderInfo{
			Type:      "sell",
			Status:    exchange.StatusOpened,
			TradePair: "BTC/LTC",

			Volume:  1.0,
			Price:   1.0,
			OrderID: 2,

			Submitted: int64(0),
			Accepted:  int64(0),
			Completed: int64(0),
		},
	}
}
func (d *mockExchange) OrderDetails(orderID int) (exchange.OrderInfo, error) {
	return exchange.OrderInfo{
		Type:      "buy",
		Status:    exchange.StatusCompleted,
		TradePair: "BTC/LTC",

		Volume:  1.0,
		Price:   1.0,
		OrderID: orderID,

		Submitted: int64(0),
		Accepted:  int64(0),
		Completed: int64(0),
	}, nil
}
func (d *mockExchange) OrderStatus(orderID int) (string, error) {
	return exchange.StatusCompleted, nil
}
func (d *mockExchange) GetBalance(sym string) (string, error) {
	return "You has 21 * 10e9 BTC", nil
}
func (d *mockExchange) OrderBook() exchange.OrderBookTracker {
	return db.NewOrderbookTracker(&redis.Options{Addr: "localhost:6379"}, "dummy")
}
