package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/skycoin/exchange-api/db"
	"github.com/skycoin/exchange-api/exchange"
	"github.com/skycoin/exchange-api/exchange/c2cx.com"

	"github.com/skycoin/exchange-api/api/rpc"
	"github.com/skycoin/exchange-api/exchange/cryptopia.co.nz"
)

/*
var (
	c2cxKey         = "2A4C851A-1B86-4E08-863B-14582094CE0F"         // = "censored"
	c2cxSecret      = "83262169-B473-4BF2-9288-5E5AC898F4B0"         // = "this too"
	cryptopiaKey    = "23a69c51c746446e819b213ef3841920"             // = "and this"
	cryptopiaSecret = "poPwm3OQGOb85L0Zf3DL4TtgLPc2OpxZg9n8G7Sv2po=" // = ":)"
)
*/
var keys struct {
	cryptopia struct {
		key    string
		secret string
	}
	c2cx struct {
		key    string
		secret string
	}
}
var (
	c2cxClient      *c2cx.Client
	cryptopiaClient *cryptopia.Client
)

func main() {
	var dbaddr = flag.String("db", "localhost:6379", "Redis address")
	var srvaddr = flag.String("srv", "localhost:12345", "RPC listener address")
	var cryptopiaKey = flag.String("cryptopia", "", "key and secret joined \":\"")
	var c2cxKey = flag.String("c2cx", "", "c2cx key and secret joined \":\"")
	flag.Parse()

	vals := strings.SplitN(*cryptopiaKey, ":", 2)
	if len(vals) != 2 {
		log.Println("-cryptopia flag does not set")
		return
	}
	keys.cryptopia.key = vals[0]
	keys.cryptopia.secret = vals[1]
	vals = strings.SplitN(*c2cxKey, ":", 2)
	if len(vals) != 2 {
		log.Println("-c2cx flag does not set")
		return
	}
	keys.c2cx.key = vals[0]
	keys.c2cx.secret = vals[1]

	cryptopiaClient = &cryptopia.Client{
		Key:    keys.cryptopia.key,
		Secret: keys.cryptopia.secret,
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: *dbaddr,
		}, "cryptopia"),
		OrderbookRefreshInterval: time.Second * 5,
		OrdersRefreshInterval:    time.Second * 5,
	}
	c2cxClient = &c2cx.Client{
		Key:    keys.c2cx.key,
		Secret: keys.c2cx.secret,
		Orders: exchange.NewTracker(),
		Orderbooks: db.NewOrderbookTracker(&redis.Options{
			Addr: *dbaddr,
		}, "c2cx"),
		OrderbookRefreshInterval: time.Second * 5,
		OrdersRefreshInterval:    time.Second * 5,
	}

	go c2cxClient.Update()
	go cryptopiaClient.Update()

	var server = rpc.Server{
		Handlers: map[string]rpc.Wrapper{
			"cryptopia": rpc.Wrapper{
				Client: cryptopiaClient,
				Env: map[string]string{
					"key":    keys.cryptopia.key,
					"secret": keys.cryptopia.secret,
				},
				Handlers: cryptopiaHandlers,
			},
			"c2cx": rpc.Wrapper{
				Client: c2cxClient,
				Env: map[string]string{
					"key":    keys.c2cx.key,
					"secret": keys.c2cx.secret,
				},
				Handlers: c2cxHandlers,
			},
		},
	}
	var stop = make(chan struct{})
	go server.Start(*srvaddr, stop)
	<-stop
	// Send anything for exit
}

// exchange-specific functions, that not handles by Client interface
var cryptopiaHandlers = map[string]rpc.HandlerFunc{
	"deposit": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidRequest, err)
		}
		currency, err := rpc.GetStringParam(params, "currency")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		addr, err := cryptopia.GetDepositAddress(env["key"], env["secret"], currency)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, addr)
	},
	"withdraw": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidRequest, err)
		}
		currency, err := rpc.GetStringParam(params, "currency")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		addr, err := rpc.GetStringParam(params, "address")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		amount, err := rpc.GetFloatParam(params, "amount")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		var (
			paymentid *string
		)
		if pid, err := rpc.GetStringParam(params, "payment_id"); err == nil {
			paymentid = &pid
		}
		result, err := cryptopia.SubmitWithdraw(env["key"], env["secret"], currency, addr, paymentid, amount)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, result)
	},
	"transactions": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidRequest, err)
		}
		txType, err := rpc.GetStringParam(params, "type")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		if txType != cryptopia.Deposit && txType != cryptopia.Withdraw {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, errors.New("invalid type"))
		}
		txs, err := cryptopia.GetTransactions(env["key"], env["secret"], txType, nil)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, txs)
	},
	"tracking_add": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidRequest, err)
		}
		market, err := rpc.GetStringParam(params, "market")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		cryptopiaClient.AddOrderbookTracking(market)
		return rpc.MakeSuccessResponse(r, nil)

	},
	"tracking_rm": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidRequest, err)
		}
		market, err := rpc.GetStringParam(params, "market")
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		err = cryptopiaClient.RemoveOrderbookTracking(market)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, nil)
	},
}

var c2cxHandlers = map[string]rpc.HandlerFunc{
	"submit_trade": func(r rpc.Request, env map[string]string) rpc.Response {
		params, err := rpc.DecodeParams(r)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidRequest, err)
		}
		var (
			advanced        *c2cx.AdvancedOrderParams
			priceTypeID     int
			orderType       string
			price, quantity float64
			market          string
		)
		if priceTypeID, err = rpc.GetIntParam(params, "price_type_id"); err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		if orderType, err = rpc.GetStringParam(params, "order_type"); err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		if market, err = rpc.GetStringParam(params, "market"); err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		if price, err = rpc.GetFloatParam(params, "price"); err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		if quantity, err = rpc.GetFloatParam(params, "quantity"); err != nil {
			return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
		}
		if adv, ok := params["advanced"]; ok {
			if data, ok := adv.(json.RawMessage); ok {
				advanced = new(c2cx.AdvancedOrderParams)
				err = json.Unmarshal(data, advanced)
				if err != nil {
					return rpc.MakeErrorResponse(r, rpc.InvalidParams, err)
				}
			}
		}
		orderid, err := c2cx.CreateOrder(env["key"], env["secret"], market, price, quantity, orderType, priceTypeID, advanced)
		if err != nil {
			return rpc.MakeErrorResponse(r, rpc.InternalError, err)
		}
		return rpc.MakeSuccessResponse(r, orderid)
	},
}
