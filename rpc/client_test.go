package rpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestClientDo(t *testing.T) {
	server := Server{
		Handlers: map[string]Wrapper{
			"test": {
				Client:   new(ex),
				Handlers: nil,
				Env:      nil,
			},
		},
	}
	addr := "localhost:12345"
	stop := make(chan struct{})
	go server.Start(addr, stop)
	time.Sleep(1 * time.Second)
	params, _ := json.Marshal(map[string]interface{}{"currency": "BTC"})
	resp, err := Do(addr, "test", Request{
		ID:      new(string),
		JSONRPC: JSONRPC,
		Params:  params,
		Method:  "balance",
	})
	if err != nil {
		t.Fatal(err)
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
	close(stop)
}
