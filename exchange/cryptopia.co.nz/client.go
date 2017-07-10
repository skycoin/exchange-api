package cryptopia

import (
	"log"
	"strconv"
	"time"

	"sync"

	"github.com/pkg/errors"
	exchange "github.com/uberfurrer/tradebot/exchange"
)

// Client impletments an exchange.Client interface
type Client struct {
	Key, Secret     string
	RefreshInterval time.Duration

	Tracker   *exchange.OrderTracker
	Orderbook exchange.OrderBookTracker

	Stop chan struct{}
	// Add concurrecny for updating
	sem chan struct{}
}

// Cancel cancels one order by given orderID
func (c *Client) Cancel(orderID int) (*exchange.OrderInfo, error) {
	orders, err := CancelTrade(c.Key, c.Secret, c.nonce(), CancelOne, orderID, "")
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("no orders cancelled")
	}
	c.Tracker.Cancel(orders...)
	info, err := c.Tracker.Get(orders[0])
	return &info, err
}

// CancelAll cancels all executed orders on account
func (c *Client) CancelAll() ([]*exchange.OrderInfo, error) {
	var result = make([]*exchange.OrderInfo, 0, len(c.Tracker.Executed()))
	orders, err := CancelTrade(c.Key, c.Secret, c.nonce(), CancelAll, 0, "")
	if err != nil {
		return nil, err
	}
	// maybe if you have several managing sources
	// in normal usecase, this does not executed
	if len(orders) > cap(result) {
		err = errors.New("additional orders cancelled")
	}
	c.Tracker.Cancel(orders...)
	for _, v := range orders {
		// if order does not found and err != nil
		// not-tracked concelled orders appear in error message
		info, cerr := c.Tracker.Get(v)
		if cerr != nil {
			err = errors.Wrap(err, strconv.Itoa(v))
			continue
		}
		result = append(result, &info)
	}
	return result, err

}

// CancelMarket cancel all orders opened in given market
func (c *Client) CancelMarket(symbol string) ([]*exchange.OrderInfo, error) {
	var result = make([]*exchange.OrderInfo, 0, len(c.Tracker.Executed()))
	orders, err := CancelTrade(c.Key, c.Secret, c.nonce(), CancelTradePair, 0, symbol)
	if err != nil {
		return nil, err
	}
	// maybe if you have several managing sources
	// in normal usecase, this does not executed
	if len(orders) > cap(result) {
		err = errors.New("additional orders cancelled")
	}
	c.Tracker.Cancel(orders...)
	for _, v := range orders {
		// if order does not found and err != nil
		// not-tracked concelled orders appear in error message
		info, cerr := c.Tracker.Get(v)
		if cerr != nil {
			err = errors.Wrap(err, strconv.Itoa(v))
			continue
		}
		result = append(result, &info)
	}
	return result, err
}

// Buy places buy order
func (c *Client) Buy(symbol string, rate, amount float64) (int, error) {
	orderData, err := SubmitTrade(c.Key, c.Secret, c.nonce(), symbol, OfTypeBuy, rate, amount)
	if err != nil {
		return 0, err
	}
	symbol = normalize(symbol)
	if orderData.OrderID == 0 {
		c.Tracker.NewOrder(symbol, OfTypeBuy, exchange.StatusCompleted, 0, amount, rate)
		return 0, nil
	}
	if orderData.FilledOrders == nil || len(orderData.FilledOrders) == 0 {
		c.Tracker.NewOrder(symbol, OfTypeBuy, exchange.StatusOpened, orderData.OrderID, amount, rate)
		return orderData.OrderID, nil
	}
	c.Tracker.NewOrder(symbol, OfTypeBuy, exchange.StatusPartial, orderData.OrderID, amount, rate)
	return orderData.OrderID, nil
}

// Sell places sell order
func (c *Client) Sell(symbol string, rate, amount float64) (int, error) {
	orderData, err := SubmitTrade(c.Key, c.Secret, c.nonce(), symbol, OfTypeSell, rate, amount)
	if err != nil {
		return 0, err
	}
	symbol = normalize(symbol)
	if orderData.OrderID == 0 {
		c.Tracker.NewOrder(symbol, OfTypeSell, exchange.StatusCompleted, 0, amount, rate)
		return 0, nil
	}
	if orderData.FilledOrders == nil || len(orderData.FilledOrders) == 0 {
		c.Tracker.NewOrder(symbol, OfTypeSell, exchange.StatusOpened, orderData.OrderID, amount, rate)
		return orderData.OrderID, nil
	}
	c.Tracker.NewOrder(symbol, OfTypeSell, exchange.StatusPartial, orderData.OrderID, amount, rate)
	return orderData.OrderID, nil
}

// OrderStatus returns a string representation of order status
func (c *Client) OrderStatus(orderID int) (string, error) {
	return c.Tracker.Status(orderID)
}

// OrderDetails returns an exchange.OrderInfo struct with detailed informaiton of order with given orderID
func (c *Client) OrderDetails(orderID int) (exchange.OrderInfo, error) {
	info, err := c.Tracker.Get(orderID)
	return info, err
}

// GetBalance returns string representation of balance informaiton for given currency
func (c *Client) GetBalance(symbol string) (string, error) {
	return GetBalance(c.Key, c.Secret, c.nonce(), symbol)
}

// Completed wraps Tracker.Completed()
func (c *Client) Completed() []*exchange.OrderInfo { return c.Tracker.Completed() }

// Executed wraps Tracker.Executed()
func (c *Client) Executed() []*exchange.OrderInfo { return c.Tracker.Executed() }

func (c *Client) nonce() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

// OrderBook returns interface for managing Orderbook
func (c *Client) OrderBook() exchange.OrderBookTracker { return c.Orderbook }

func (c *Client) checkUpdate() {
	// Update orderbook
	if c.Orderbook != nil {
		c.sem <- struct{}{}
		var wg sync.WaitGroup
		wg.Add(len(marketCache))
		for k := range marketCache {
			go func(w *sync.WaitGroup, market string) {
				defer w.Done()
				book, err := GetMarketOrders(market, 100)
				if err != nil {
					log.Println("cryptopia: update error:", err)
					return
				}
				var (
					bids = make([]exchange.OrderbookEntry, len(book.Buy))
					asks = make([]exchange.OrderbookEntry, len(book.Sell))
				)
				for i, v := range book.Buy {
					bids[i] = exchange.OrderbookEntry{
						Price:  v.Price,
						Volume: v.Volume,
					}
				}
				for i, v := range book.Sell {
					asks[i] = exchange.OrderbookEntry{
						Price:  v.Price,
						Volume: v.Volume,
					}
				}
				c.Orderbook.UpdateSym(market, bids, asks)
			}(&wg, k)
		}
		wg.Wait()
		<-c.sem
	}
	// Update placed orders
	if c.Tracker != nil {
		var openedCount = len(c.Tracker.Executed())
		for {
			orders, err := GetOpenOrders(c.Key, c.Key, c.nonce(), AllMarkets, openedCount)
			if err != nil {
				log.Println("cryptopia: update error", err)
				break
			}
			for _, v := range orders {
				var (
					t, _         = time.Parse(time.RFC3339, v.Timestamp)
					status       string
					remaining, _ = strconv.ParseFloat(v.Remaining, 64)
					total, _     = strconv.ParseFloat(v.Total, 64)
				)

				//TODO: check float equaling,
				// maybe need compare thru epsilon value
				if total-remaining < total {
					status = exchange.StatusPartial
				} else {
					status = exchange.StatusOpened
				}
				c.Tracker.UpdateOrderDetails(v.OrderID, status, &t)
			}
			break
		}
		for {

			//TODO: check here, what time returns from exchange
			// means that GetTradeHistory return time of completion
			orders, err := GetTradeHistory(c.Key, c.Key, c.nonce(), AllMarkets, openedCount)
			if err != nil {
				log.Println("cryptopia: update error", err)
				break
			}
			for _, v := range orders {
				var t, _ = time.Parse(time.RFC3339, v.Timestamp)
				c.Tracker.Complete(v.OrderID, t)
			}
			break
		}
	}
}

// Update runs update cycle for Client
// It also updates orderbook
func (c *Client) Update() {
	c.Stop = make(chan struct{})
	c.sem = make(chan struct{}, 1)
	ticker := time.NewTicker(c.RefreshInterval)
	for {
		select {
		case <-c.Stop:
			ticker.Stop()
			return
		case <-ticker.C:
			c.checkUpdate()
		}
	}
}
