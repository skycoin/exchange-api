package rpc

import (
	"encoding/json"
	"testing"

	"github.com/skycoin/exchange-api/exchange"
)

var client = ex{}

func Test_defaultHandler_GetBalance(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"currency\":\"BTC\"}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["balance"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != "\"You has 21 * 10e9 BTC\"" {
		t.Fatalf("want: %s, expected: %s", "\"You has 21 * 10e9 BTC\"", resp.Result)
	}
}
func Test_defaulthandler_Cancel(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"orderid\":1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["cancel_trade"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	want := "{\"orderid\":0,\"type\":\"\",\"market\":\"\",\"amount\":0,\"price\":0,\"submitted_at\":-6795364578871,\"fee\":0,\"completed_amount\":0,\"status\":\"\",\"accepted_at\":-6795364578871,\"completed_at\":-6795364578871}"
	if string(resp.Result) != want {
		t.Fatalf("want: %s, expected: %s", want, resp.Result)
	}
}
func Test_defaulthandler_CancelAll(t *testing.T) {
	req := Request{
		Params:  json.RawMessage(nil),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	want := "[]"
	resp := defaultHandlers["cancel_all"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != want {
		t.Fatalf("want: %s, expected: %s", want, resp.Result)
	}
}
func Test_defaulthandler_CancelMarket(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\"}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["cancel_market"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	want := "[]"
	if string(resp.Result) != want {
		t.Fatalf("want: %s, expected: %s", want, resp.Result)
	}
}
func Test_defaulthandler_Buy(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\",\"rate\":1.0,\"amount\":1.1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["buy"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != "1" {
		t.Fatalf("want: 1, exepected: %s", resp.Result)
	}
}
func Test_defaulthandler_Sell(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\",\"rate\":1.0,\"amount\":1.1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["sell"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != "2" {
		t.Fatalf("want: 2, exepected: %s", resp.Result)
	}
}
func Test_defaulthandler_OrderDetails(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"orderid\":1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["order_info"](req, client)
	want := "{\"orderid\":0,\"type\":\"\",\"market\":\"\",\"amount\":0,\"price\":0,\"submitted_at\":-6795364578871,\"fee\":0,\"completed_amount\":0,\"status\":\"\",\"accepted_at\":-6795364578871,\"completed_at\":-6795364578871}"
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != want {
		t.Fatalf("want: %s, expected: %s", want, resp.Result)
	}
}
func Test_defaulthandler_OrderStatus(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"orderid\":1}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["order_status"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != "\"completed\"" {
		t.Fatalf("want: \"Completed\", expected: %s", resp.Result)
	}
}
func Test_defaulthandler_Completed(t *testing.T) {
	req := Request{
		Params:  json.RawMessage(nil),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	want := "[]"
	resp := defaultHandlers["completed"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != want {
		t.Fatalf("want %s, expected: %s", want, resp.Result)
	}
}
func Test_defaulthandler_Executed(t *testing.T) {
	req := Request{
		Params:  json.RawMessage(nil),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	want := "[]"
	resp := defaultHandlers["executed"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != want {
		t.Fatalf("want: %s, expected: %s", want, resp.Result)
	}
}

type ex struct{}

func (ex) Buy(string, float64, float64) (int, error)     { return 1, nil }
func (ex) Sell(string, float64, float64) (int, error)    { return 2, nil }
func (ex) Cancel(int) (exchange.Order, error)            { return exchange.Order{}, nil }
func (ex) CancelMarket(string) ([]exchange.Order, error) { return []exchange.Order{}, nil }
func (ex) CancelAll() ([]exchange.Order, error)          { return []exchange.Order{}, nil }
func (ex) Executed() []int                               { return []int{} }
func (ex) Completed() []int                              { return []int{} }
func (ex) GetBalance(string) (string, error)             { return "You has 21 * 10e9 BTC", nil }
func (ex) OrderDetails(int) (exchange.Order, error)      { return exchange.Order{}, nil }
func (ex) OrderStatus(int) (string, error)               { return exchange.Completed, nil }
func (ex) Orderbook() exchange.Orderbooks                { return nil }
