package cli

import (
	"crypto/rand"
	"math/big"

	"github.com/pkg/errors"
)

var errInvalidParams = errors.New("invalid params")
var errRPC = errors.New("RPC request failed")
var errInvalidResponse = errors.New("unexpected response format")

var (
	rpcaddr  = "localhost:12345"
	endpoint string
)

func reqID() *string {
	v, _ := rand.Int(rand.Reader, new(big.Int).SetInt64(1<<62))
	str := v.String()
	return &str
}
