package c2cx

import (
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/uberfurrer/tradebot/db"
)

func TestClientUpdateOrderbook(t *testing.T) {
	var c = Client{
		Key:     "",
		Secret:  "",
		Tracker: nil,
		OrderBookTracker: db.NewOrderbookTracker(&redis.Options{
			Addr: "localhost:6379",
		}, "c2cx"),
		sem: make(chan struct{}, 1),
	}
	c.checkUpdate()
	time.Sleep(time.Second * 20)
}
func TestAPI_requestPost(t *testing.T) {
	var (
		Key    = "2A4C851A-1B86-4E08-863B-14582094CE0F"
		Secret = "83262169-B473-4BF2-9288-5E5AC898F4B0"
	)
	var params = url.Values{}
	params.Add("symbol", "CNY_BTC")
	params.Add("interval", "100")
	log.Println(requestPost("getorderbystatus", Key, Secret, params))
}
