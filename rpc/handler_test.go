package rpc

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"

	exchange "github.com/skycoin/exchange-api/exchange"

	"reflect"
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
	var (
		result decimal.Decimal
		target = decimal.NewFromFloat(1.234)
	)
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatal(err)
	}
	if !target.Equal(result) {
		t.Fatalf("expected: %s; received: %s", target.String(), result.String())
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
	want := map[string]interface{}{
		"orderid":          float64(0),
		"type":             "",
		"market":           "",
		"amount":           "0",
		"price":            "0",
		"submitted_at":     float64(-6795364578871),
		"fee":              "0",
		"completed_amount": "0",
		"status":           "",
		"accepted_at":      float64(-6795364578871),
		"completed_at":     float64(-6795364578871),
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(resp.Result, &parsed); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, parsed) {
		t.Fatalf("expected: %v, received: %v", want, parsed)
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
		t.Fatalf("expected: %s, received: %s", want, resp.Result)
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
		t.Fatalf("expected: %s, received: %s", want, resp.Result)
	}
}
func Test_defaulthandler_Buy(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\",\"price\":\"1.0\",\"amount\":\"1.1\"}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["buy"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != "1" {
		t.Fatalf("expected: 1, received: %s", resp.Result)
	}
}
func Test_defaulthandler_Sell(t *testing.T) {
	req := Request{
		Params:  json.RawMessage("{\"symbol\":\"BTC/LTC\",\"price\":\"1.0\",\"amount\":\"1.1\"}"),
		Method:  "/",
		ID:      new(string),
		JSONRPC: JSONRPC,
	}
	resp := defaultHandlers["sell"](req, client)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if string(resp.Result) != "2" {
		t.Fatalf("expected: 2, received: %s", resp.Result)
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
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	want := map[string]interface{}{
		"orderid":          float64(0),
		"type":             "",
		"market":           "",
		"amount":           "0",
		"price":            "0",
		"submitted_at":     float64(-6795364578871),
		"fee":              "0",
		"completed_amount": "0",
		"status":           "",
		"accepted_at":      float64(-6795364578871),
		"completed_at":     float64(-6795364578871),
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(resp.Result, &parsed); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, parsed) {
		t.Fatalf("expected: %v, received: %v", want, parsed)
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
		t.Fatalf("expected: \"Completed\", received: %s", resp.Result)
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
		t.Fatalf("expected %s, received: %s", want, resp.Result)
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
		t.Fatalf("expected: %s, received: %s", want, resp.Result)
	}
}

type ex struct{}

func (ex) Buy(string, decimal.Decimal, decimal.Decimal) (int, error)  { return 1, nil }
func (ex) Sell(string, decimal.Decimal, decimal.Decimal) (int, error) { return 2, nil }
func (ex) Cancel(int) (exchange.Order, error)                         { return exchange.Order{}, nil }
func (ex) CancelMarket(string) ([]exchange.Order, error)              { return []exchange.Order{}, nil }
func (ex) CancelAll() ([]exchange.Order, error)                       { return []exchange.Order{}, nil }
func (ex) Executed() []int                                            { return []int{} }
func (ex) Completed() []int                                           { return []int{} }
func (ex) GetBalance(string) (decimal.Decimal, error)                 { return decimal.NewFromFloat(1.234), nil }
func (ex) OrderDetails(int) (exchange.Order, error)                   { return exchange.Order{}, nil }
func (ex) OrderStatus(int) (string, error)                            { return exchange.Completed, nil }
