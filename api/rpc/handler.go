package rpc

import (
	"encoding/json"
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
	// GetBalance params: {"currency": string}
	"GetBalance": func(r Request, c exchange.Client) Response {
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
	// Cancel params: {"orderid": int}
	"Cancel": func(r Request, c exchange.Client) Response {
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
	// CancelAll params should be omitted, empty or null
	"CancelAll": func(r Request, c exchange.Client) Response {
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
	// CancelMarket params: {"symbol": string}
	"CancelMarket": func(r Request, c exchange.Client) Response {
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
	// Buy params: {"symbol": string; "rate": float64, "amount": float64}
	"Buy": func(r Request, c exchange.Client) Response {
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
			rate, err := GetFloatParam(params, "rate")
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
	// Sell params: {"symbol": string, "rate": float64, "amount": float64}
	"Sell": func(r Request, c exchange.Client) Response {
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
			rate, err := GetFloatParam(params, "rate")
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
	// OrderDetails params: {"orderid": int}
	"OrderDetails": func(r Request, c exchange.Client) Response {
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
	// OrderStatus params: {"orderid": int}
	"OrderStatus": func(r Request, c exchange.Client) Response {
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
	// Completed params should be omitted, empty or null
	"Completed": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		_, err = DecodeParams(r)
		if err != nil && err != errEmptyParams {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
			return resp
		}
		orders := c.Completed()
		resp.setBody(orders)
		return resp
	},
	// Executed params should be omitted, empty or null
	"Executed": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		_, err = DecodeParams(r)
		if err != nil && err != errEmptyParams {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
			return resp
		}
		orders := c.Executed()
		resp.setBody(orders)
		return resp
	},
	//OrderBook params: {"market": string}
	"OrderBook": func(r Request, c exchange.Client) Response {
		resp, err := validateRequest(r)
		if err != nil {
			return resp
		}
		params, err := DecodeParams(r)
		if err != nil {
			resp.Error = makeError(ParseError, parseErrorMsg, err)
			return resp
		}
		market, err := GetStringParam(params, "market")
		if err != nil {
			resp.Error = makeError(InvalidParams, invalidParamsMsg, err)
		}
		orderbook := c.OrderBook()
		book, err := orderbook.GetRecord(market)
		if err != nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
			return resp
		}
		data, err := json.Marshal(book)
		log.Println("Orderbook data:", data)
		if err != nil {
			resp.Error = makeError(InternalError, internalErrorMsg, err)
			return resp
		}
		resp.setBody(data)
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

// Process execute request and return response
func (h *PackageHandler) Process(r Request) Response {
	log.Println(r, string(r.Params))
	if f, ok := defaultHandlers[r.Method]; ok {
		return f(r, h.Client)
	}
	if f, ok := h.Handlers[r.Method]; ok {
		return f(r, h.Env)
	}
	if r.ID != nil {
		return Response{
			JSONRPC: JSONRPC,
			ID:      *r.ID,
			Error:   makeError(MethodNotFound, methodNotFoundMsg, nil),
			Result:  nil,
		}
	}
	return Response{}

}
