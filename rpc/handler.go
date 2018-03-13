package rpc

import (
	"log"

	exchange "github.com/skycoin/exchange-api/exchange"
)

// Wrapper handles one exchange, resolve methods and returns json formatted responses
type Wrapper struct {
	Client exchange.Client
	// Handlers contains all Handlers for all functions, provided by exchange
	Handlers map[string]HandlerFunc
	// env needs if you want call additional, package-specific functions, that does not included in exchange.Client interface
	// env typically contains api keys
	Env map[string]string
}

// HandlerFunc wraps function, provides by exchange package
// HandlerFunc should validating params and checks correct for request scheme - matching jsonrpc version and has id
type HandlerFunc func(r Request, env map[string]string) Response

// defaultHandlers contain functions that will executed, if called functionality, handled by exchange.Client interface
var defaultHandlers = map[string]func(Request, exchange.Client) Response{
	// balance params: {"currency": string}
	"balance": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, nil)
		}
		currency, err := GetStringParam(params, "currency")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		balance, err := c.GetBalance(currency)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(balance)
		}
		return resp
	},
	// cancel_trade params: {"orderid": int}
	"cancel_trade": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orderID, err := GetIntParam(params, "orderid")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		order, err := c.Cancel(orderID)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(order)
		}
		return resp
	},
	// cancel_all params should be omitted, empty or null
	"cancel_all": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		_, err = DecodeParams(r)
		if err != nil && err != errEmptyParams {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
		}
		orders, err := c.CancelAll()
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		resp.setBody(orders)
		return resp
	},
	// cancel_market params: {"symbol": string}
	"cancel_market": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
		}
		symbol, err := GetStringParam(params, "symbol")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orders, err := c.CancelMarket(symbol)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(orders)
		}
		return resp
	},
	// buy params: {"symbol": string; "rate": decimal.Decimal, "amount": decimal.Decimal}
	"buy": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
		}
		symbol, err := GetStringParam(params, "symbol")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		rate, err := GetDecimalParam(params, "price")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		amount, err := GetDecimalParam(params, "amount")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orderID, err := c.Buy(symbol, rate, amount)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(orderID)
		}
		return resp
	},
	// sell params: {"symbol": string, "rate": decimal.Decimal, "amount": decimal.Decimal}
	"sell": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
		}
		symbol, err := GetStringParam(params, "symbol")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		rate, err := GetDecimalParam(params, "price")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		amount, err := GetDecimalParam(params, "amount")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orderID, err := c.Sell(symbol, rate, amount)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(orderID)
		}
		return resp
	},
	// order_info params: {"orderid": int}
	"order_info": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orderID, err := GetIntParam(params, "orderid")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		order, err := c.OrderDetails(orderID)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(order)
		}
		return resp
	},
	// order_status params: {"orderid": int}
	"order_status": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orderID, err := GetIntParam(params, "orderid")
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		order, err := c.OrderStatus(orderID)
		if err != nil && resp.Error == nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
		}
		if err == nil {
			resp.setBody(order)
		}
		return resp
	},
	// completed params should be omitted, empty or null
	"completed": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		_, err = DecodeParams(r)
		if err != nil && err != errEmptyParams {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
			return resp
		}
		resp.setBody(c.Completed())
		return resp
	},
	// executed params should be omitted, empty or null
	"executed": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		_, err = DecodeParams(r)
		if err != nil && err != errEmptyParams {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
			return resp
		}
		resp.setBody(c.Executed())
		return resp
	},
}

// PackageFunc adds a PackageFunc for specified method
func (h *Wrapper) PackageFunc(method string, f HandlerFunc) {
	if h.Handlers == nil {
		h.Handlers = make(map[string]HandlerFunc)
	}
	h.Handlers[method] = f
}

// Setenv sets a environment variable
func (h *Wrapper) Setenv(key, value string) {
	if h.Env == nil {
		h.Env = make(map[string]string)
	}
	h.Env[key] = value
}

// Process lookups given method and calls it
func (h *Wrapper) Process(r Request) *Response {
	log.Printf("processing request, method %s, params %s\n", r.Method, r.Params)
	if f, ok := defaultHandlers[r.Method]; ok {
		resp := f(r, h.Client)
		return &resp
	}
	if f, ok := h.Handlers[r.Method]; ok {
		resp := f(r, h.Env)
		return &resp
	}
	if r.ID != nil {
		return &Response{
			JSONRPC: JSONRPC,
			ID:      *r.ID,
			Error:   makeError(MethodNotFound, methodNotFoundMsg, nil),
			Result:  nil,
		}
	}
	return nil
}
