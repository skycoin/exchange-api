package rpc

import (
	"encoding/json"
	"testing"
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
	if string(resp.Result) != "\"You has 21 * 10e9 BTC\"" {
		t.Fatal("want \"You has 21 * 10e9 BTC\", expected", string(resp.Result))
	}
	close(stop)
}
