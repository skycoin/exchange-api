package main

import (
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/uberfurrer/tradebot/db"
	c2cx "github.com/uberfurrer/tradebot/exchange/c2cx.com"
	cryptopia "github.com/uberfurrer/tradebot/exchange/cryptopia.co.nz"
	"github.com/uberfurrer/tradebot/logger"
)

var (
	c2cxKey         = "censored"
	c2cxSecret      = "this too"
	cryptopiaKey    = "and this"
	cryptopiaSecret = ":)"
)

func main() {
	logger.SetLogLevel(logger.Errors, os.Stdout)

	var c2cxClient = c2cx.Client{
		Key:    c2cxKey,
		Secret: c2cxSecret,
		OrderBookTracker: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx"),
		RefreshInterval: 1000,
	}
	go c2cxClient.Update()
	var c2cxRPC = rpc.PackageHandler{
		Client:   &c2cxClient,
		Handlers: c2cxHandlers,
		Env: map[string]string{
			"key":    c2cxKey,
			"secret": c2cxSecret,
		},
	}

	var cryptopiaClient = cryptopia.Client{
		Key:    cryptopiaKey,
		Secret: cryptopiaSecret,
		Orderbook: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "cryptopia"),
		RefreshInterval: 1000,
		TrackedBooks:    []string{"LTC/BTC", "SKY/DOGE"},
	}
	go cryptopiaClient.Update()
	var cryptopiaRPC = rpc.PackageHandler{
		Client:   &cryptopiaClient,
		Handlers: cryptopiaHandlers,
		Env: map[string]string{
			"key":    cryptopiaKey,
			"secret": cryptopiaSecret,
		},
	}
	var srv = rpc.Server{
		Handlers: map[string]rpc.PackageHandler{
			"c2cx":      c2cxRPC,
			"cryptopia": cryptopiaRPC,
		},
	}
	var stop = make(chan struct{})
	srv.Start("localhost:12345", stop)

	//close stop for exit
}

// Additional functions, that does not wrapped by exchange.Client interface
var cryptopiaHandlers = map[string]rpc.PackageFunc{
	// GetDepositAddress gets deposit address for given currency
	// params: {"currency":string}
	"GetDepositAddress": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.ParseError, err)
		}
		currency, err := rpc.GetStringParam(params, "currency")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		address, err := cryptopia.GetDepositAddress(env["key"], env["secret"], cryptopia.Nonce(), currency)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, address)
	},
	// SubmitWithdraw creates withdrawal request with given currency and amount
	// params: {"currency": string, "amount": float64, "paymentid":string, optional(needs only for CryptoNote based currencies)}
	"SubmitWithdraw": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.ParseError, err)
		}
		var (
			amount    float64
			paymentID string
			address   string
			currency  string
		)
		amount, err = rpc.GetFloatParam(params, "amount")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		paymentID, _ = rpc.GetStringParam(params, "paymentid")
		address, err = rpc.GetStringParam(params, "address")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		currency, err = rpc.GetStringParam(params, "currency")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		withdraw, err := cryptopia.SubmitWithdraw(env["key"], env["secret"], cryptopia.Nonce(), currency, address, paymentID, amount)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, withdraw)
	},
	// GetTransactions gets list of transactions of given type
	// params: {"type":string("Deposit" or "Withdraw"), "count":int, optional(default value is 100)}
	"GetTransactions": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.ParseError, err)
		}
		var (
			txType string
			count  int
		)
		txType, err = rpc.GetStringParam(params, "type")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		count, _ = rpc.GetIntParam(params, "count")
		transactions, err := cryptopia.GetTransactions(env["key"], env["secret"], cryptopia.Nonce(), txType, count)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, transactions)
	},
}
var c2cxHandlers = map[string]rpc.PackageFunc{
	// This order wont tracked
	// Calling SubmitTrade directly allows to creating limited order
	// params:{ "ordertype":string("buy" or "sell"), "pricetype": string("Market" or "Limit"), "symbol": string,
	//          "triggerprice": float64, "quantity": float64, "price": float64, "takeprofit":float64,
	//          "stoploss":float64, "exptime": string(time in RFC3389 format)}
	// required params: "ordertype", "pricetype","symbol" ,"triggerprice"(if pricetype - market, then this field means amount), "quantity", "price"
	"SubmitTrade": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.ParseError, err)
		}
		var (
			ordertype, exptime, pricetype, symbol string
			triggerprice, price, quantity         float64
			priceTypeID                           int
			takeProfit, stopLoss                  *float64
			expTime                               *time.Time
		)
		ordertype, err = rpc.GetStringParam(params, "ordertype")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		pricetype, err = rpc.GetStringParam(params, "pricetype")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		switch strings.ToLower(pricetype) {
		case "limit":
			priceTypeID = c2cx.PriceTypeLimit
		case "market":
			priceTypeID = c2cx.PriceTypeMarket
		default:
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, errors.Errorf("price type should be \"market\" or \"limit\", given %s", pricetype))
		}
		symbol, err = rpc.GetStringParam(params, "symbol")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		triggerprice, err = rpc.GetFloatParam(params, "triggerprice")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		quantity, err = rpc.GetFloatParam(params, "quantity")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		price, err = rpc.GetFloatParam(params, "price")
		takeprofit, err := rpc.GetFloatParam(params, "takeprofit")
		if err == nil {
			takeProfit = &takeprofit
		}
		stoploss, err := rpc.GetFloatParam(params, "stoploss")
		if err == nil {
			stopLoss = &stoploss
		}
		exptime, err = rpc.GetStringParam(params, "exptime")
		if err != nil {
			t, err := time.Parse(time.RFC3339, exptime)
			if err != nil {
				return rpc.MakeErrorResponse(r, rpc.ParseError, err)
			}
			expTime = &t
		}
		orderID, err := c2cx.CreateOrder(env["key"], env["secret"], symbol, ordertype, priceTypeID, triggerprice, quantity, price,
			takeProfit, stopLoss, expTime)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, orderID)
	},
}
