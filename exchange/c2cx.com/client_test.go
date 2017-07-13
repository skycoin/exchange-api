package c2cx

import (
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
