package rpc

import (
	"encoding/json"
	"fmt"
)

// JSONRPC version
const JSONRPC = "2.0"

// Request represents a JSONRPC 2.0 request message
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      *string         `json:"id"`
}

// Response represents a JSONRPC 2.0 response message
type Response struct {
	ID      string          `json:"id,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	JSONRPC string          `json:"jsonrpc"`
}

func (r *Response) setBody(v interface{}) {
	body, err := json.Marshal(v)
	if err != nil {
		r.Result = nil
		r.Error = makeError(InternalError, internalErrorMsg, err)
		return
	}
	r.Result = body
}

//Predefined messages
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)
const (
	parseErrorMsg     = "Parse Error"
	invalidRequestMsg = "Invalid Request"
	methodNotFoundMsg = "Method Not Found"
	invalidParamsMsg  = "Invalid Params"
	internalErrorMsg  = "Internal Error"
)

type jsonrpcError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    *string `json:"data,omitempty"`
}

// Implements error interface
func (err *jsonrpcError) Error() string {
	return fmt.Sprintf("jsonrpc error: %d %s %s", err.Code, err.Message, *err.Data)
}

func makeError(code int, message string, additional error) *jsonrpcError {
	var datastr *string
	if additional != nil {
		datastr = new(string)
		*datastr = additional.Error()
	}
	return &jsonrpcError{Code: code, Message: message, Data: datastr}
}
