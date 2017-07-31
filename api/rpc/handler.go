package rpc

import (
	"log"

	exchange "github.com/uberfurrer/tradebot/exchange"
)

// PackageHandler handles one exchange, resolve methods and returns json formatted responses
type PackageHandler struct {
	Client exchange.Client
	// Handlers contains all Handlers for all functions, provided by exchange
	Handlers map[string]PackageFunc
	// env needs if you want call additional, package-specific functions, that does not included in exchange.Client interface
	// env typically contains api keys
	Env map[string]string
}

// PackageFunc wraps function, provides by exchange package
// PackageFunc should validating params and checks correct for request scheme - matching jsonrpc version and has id
type PackageFunc func(r Request, env map[string]string) Response

// defaultHandlers contain functions that will executed, if called functionality, handled by exchange.Client interface
var defaultHandlers = map[string]func(Request, exchange.Client) Response{
	// balance params: {"currency": string}
	"balance": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, nil)
				break
			}
			currency, err := GetStringParam(params, "currency")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			balance, err := c.GetBalance(currency)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(balance)
			break
		}
		return resp
	},
	// cancel_trade params: {"orderid": int}
	"cancel_trade": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			orderID, err := GetIntParam(params, "orderid")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			order, err := c.Cancel(orderID)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(order)
			break
		}
		return resp
	},
	// cancel_all params should be omitted, empty or null
	"cancel_all": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			_, err = DecodeParams(r)
			if err != nil && err != errEmptyParams {
				resp.Error = makeError(ParseError, parseErrorMsg, err)
				break
			}
			orders, err := c.CancelAll()
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(orders)
			break
		}
		return resp
	},
	// cancel_market params: {"symbol": string}
	"cancel_market": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(ParseError, parseErrorMsg, err)
				break
			}
			symbol, err := GetStringParam(params, "symbol")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			orders, err := c.CancelMarket(symbol)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(orders)
			break
		}
		return resp
	},
	// buy params: {"symbol": string; "rate": float64, "amount": float64}
	"buy": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(ParseError, parseErrorMsg, err)
				break
			}
			symbol, err := GetStringParam(params, "symbol")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			rate, err := GetFloatParam(params, "price")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			amount, err := GetFloatParam(params, "amount")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			orderID, err := c.Buy(symbol, rate, amount)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(orderID)
			break
		}

		return resp
	},
	// sell params: {"symbol": string, "rate": float64, "amount": float64}
	"sell": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(ParseError, parseErrorMsg, err)
				break
			}
			symbol, err := GetStringParam(params, "symbol")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			rate, err := GetFloatParam(params, "price")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			amount, err := GetFloatParam(params, "amount")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			orderID, err := c.Sell(symbol, rate, amount)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(orderID)
			break
		}

		return resp
	},
	// order_info params: {"orderid": int}
	"order_info": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			orderID, err := GetIntParam(params, "orderid")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
			}
			order, err := c.OrderDetails(orderID)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(order)
			break
		}
		return resp
	},
	// order_status params: {"orderid": int}
	"order_status": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		for {
			params, err := DecodeParams(r)
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
				break
			}
			orderID, err := GetIntParam(params, "orderid")
			if err != nil {
				resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
			}
			order, err := c.OrderStatus(orderID)
			if err != nil {
				resp.Error = makeError(InternalError, internalErrorMsg, err)
				break
			}
			resp.setBody(order)
			break
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
	//orderbook params: {"symbol": string}
	"orderbook": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
			return resp
		}
		market, err := GetStringParam(params, "symbol")
		if err != nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		book, err := c.Orderbook().Get(market)
		if err != nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
			return resp
		}
		resp.setBody(book)
		return resp
	},
}

// PackageFunc adds a PackageFunc for specified method
func (h *PackageHandler) PackageFunc(method string, f PackageFunc) {
	if h.Handlers == nil {
		h.Handlers = make(map[string]PackageFunc)
	}
	h.Handlers[method] = f
}

// Setenv sets a environment variable
func (h *PackageHandler) Setenv(key, value string) {
	if h.Env == nil {
		h.Env = make(map[string]string)
	}
	h.Env[key] = value
}

// Process lookups given method and calls it
func (h *PackageHandler) Process(r Request) *Response {
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
