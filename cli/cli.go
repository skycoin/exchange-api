package cli

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/skycoin/exchange-api/rpc"
)

var errInvalidInput = errors.New("invalid input params")

func reqID() *string {
	v, err := rand.Int(rand.Reader, new(big.Int).SetInt64(1<<62))
	if err != nil {
		panic(err)
	}
	str := v.String()
	return &str
}

func rpcRequest(method string, params map[string]interface{}) (json.RawMessage, error) {
	p, err := json.Marshal(params)
	req := rpc.Request{
		ID:      reqID(),
		JSONRPC: rpc.JSONRPC,
		Method:  method,
		Params:  p,
	}
	if err != nil {
		return nil, err
	}
	resp, err := rpc.Do(*rpcaddr, endpoint, req)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil

}
