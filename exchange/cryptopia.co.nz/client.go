package cryptopia

import (
	"time"

	"errors"

	"github.com/shopspring/decimal"

	"github.com/skycoin/exchange-api/exchange"
)

// Client impletments an exchange.Client interface
type Client struct {
	Key, Secret              string
	OrdersRefreshInterval    time.Duration
	OrderbookRefreshInterval time.Duration

	Orders     exchange.Orders
	Orderbooks exchange.Orderbooks

	// Cause cryptopia.co.nz supported greater than 1000 curencies, you need to select markets, thats will be tracked
	TrackedBooks []string

	Stop chan struct{}
	// Add concurrecny for updating
	instantOrdersCounter int
}

// Cancel cancels one order by given orderID
func (c *Client) Cancel(orderID int) (exchange.Order, error) {
	cancelled, err := cancelTrade(c.Key, c.Secret, nonce(), ByOrderID, nil, &orderID)
	if err != nil {
		return exchange.Order{}, err
	}
	if len(cancelled) != 1 {
		return exchange.Order{}, errors.New("no orders cancelled")
	}

	v, err := lookupOrder(c.Key, c.Secret, orderID)
	if err != nil {
		return exchange.Order{}, err
	}
	order := convert(v)
	order.Status = exchange.Cancelled
	order.Completed = time.Now()
	if err = c.Orders.UpdateOrder(order); err != nil {
		return order, err
	}
	return c.Orders.GetOrderInfo(order.OrderID)
}

// CancelAll cancels all executed orders on account
func (c *Client) CancelAll() ([]exchange.Order, error) {
	orders, err := cancelTrade(c.Key, c.Secret, nonce(), All, nil, nil)
	if err != nil {
		return nil, err
	}
	var cancelled []int
	for _, v := range orders {
		j, err := lookupOrder(c.Key, c.Secret, v)
		if err != nil {
			continue
		}
		order := convert(j)
		order.Status = exchange.Cancelled
		order.Completed = time.Now()
		if err = c.Orders.UpdateOrder(order); err != nil {
			return nil, err
		}
		cancelled = append(cancelled, v)
	}
	result := make([]exchange.Order, len(cancelled))

	for i, v := range cancelled {
		result[i], _ = c.Orders.GetOrderInfo(v)
	}

	return result, nil
}

// CancelMarket cancel all orders opened in given market
func (c *Client) CancelMarket(symbol string) ([]exchange.Order, error) {
	orders, err := cancelTrade(c.Key, c.Secret, nonce(), ByMarket, &symbol, nil)
	if err != nil {
		return nil, err
	}
	var cancelled []int
	for _, v := range orders {
		j, err := lookupOrder(c.Key, c.Secret, v)
		if err != nil {
			continue
		}
		order := convert(j)
		order.Status = exchange.Cancelled
		order.Completed = time.Now()
		if err = c.Orders.UpdateOrder(order); err != nil {
			return nil, err
		}
		cancelled = append(cancelled, v)
	}
	var result = make([]exchange.Order, len(cancelled))
	for i, v := range cancelled {
		result[i], _ = c.Orders.GetOrderInfo(v)
	}
	return result, nil
}

// Buy places buy order
func (c *Client) Buy(symbol string, rate, amount decimal.Decimal) (int, error) {
	var order = exchange.Order{
		Submitted: time.Now(),
		Type:      exchange.Buy,
		Market:    symbol,
		Price:     rate,
		Amount:    amount,
		Status:    exchange.Completed,
	}
	orderID, err := submitTrade(c.Key, c.Secret, nonce(), symbol, Buy, rate, amount)
	if err != nil && err != ErrInstant {
		return 0, err
	}
	// Order placing successfully, order instant completed
	if err == ErrInstant {
		c.instantOrdersCounter--
		order.OrderID = c.instantOrdersCounter
	}
	// Order placing successfully, order status - opened or partial
	order.Status = exchange.Opened
	order.OrderID = orderID
	err = c.Orders.Push(order)
	return order.OrderID, err

}

// Sell places sell order
func (c *Client) Sell(symbol string, rate, amount decimal.Decimal) (int, error) {
	var order = exchange.Order{
		Submitted: time.Now(),
		Type:      exchange.Sell,
		Market:    symbol,
		Price:     rate,
		Amount:    amount,
		Status:    exchange.Completed,
	}
	orderID, err := submitTrade(c.Key, c.Secret, nonce(), symbol, Sell, rate, amount)
	if err != nil && err != ErrInstant {
		return 0, err
	}
	// Order placing successfully, order instant completed
	if err == ErrInstant {
		c.instantOrdersCounter--
		order.OrderID = c.instantOrdersCounter
	}
	// Order placeing successfully, order status - opened
	order.Status = exchange.Opened
	order.OrderID = orderID
	err = c.Orders.Push(order)
	return order.OrderID, err
}

// OrderStatus returns a string representation of order status
func (c *Client) OrderStatus(orderID int) (string, error) {
	return c.Orders.GetOrderStatus(orderID)
}

// OrderDetails returns detailed informaiton of order with given orderID
func (c *Client) OrderDetails(orderID int) (exchange.Order, error) {
	order, err := lookupOrder(c.Key, c.Secret, orderID)
	if err != nil {
		return exchange.Order{}, err
	}
	if err = c.Orders.UpdateOrder(convert(order)); err != nil {
		return exchange.Order{}, err
	}
	return c.Orders.GetOrderInfo(orderID)
}

// GetBalance returns string representation of balance informaiton for given currency
func (c *Client) GetBalance(symbol string) (string, error) {
	return getBalance(c.Key, c.Secret, nonce(), symbol)
}

// Completed wraps Tracker.Completed()
func (c *Client) Completed() []int {
	return c.Orders.GetCompleted()
}

// Executed wraps Tracker.Executed()
func (c *Client) Executed() []int {
	return c.Orders.GetOpened()
}

// Orderbook returns interface for managing Orderbook
func (c *Client) Orderbook() exchange.Orderbooks {
	return c.Orderbooks
}

func (c *Client) updateOrderbook() {
	ordergroups, err := getMarketOrderGroups(100, c.TrackedBooks...)
	if err != nil {
		return
	}
	for _, v := range ordergroups {
		var (
			bids = make([]exchange.MarketOrder, len(v.Buy))
			asks = make([]exchange.MarketOrder, len(v.Sell))
		)
		for i, k := range v.Buy {
			bids[i] = exchange.MarketOrder{
				Price:  k.Price,
				Volume: k.Volume,
			}
		}
		for i, k := range v.Sell {
			asks[i] = exchange.MarketOrder{
				Price:  k.Price,
				Volume: k.Volume,
			}
		}
		c.Orderbooks.Update(v.Label, bids, asks)
	}
}
func (c *Client) updateOrders() {
	var count = len(c.Orders.GetOpened())
	orders, err := getOpenOrders(c.Key, c.Secret, nonce(), nil, &count)
	if err != nil {
		return
	}
	for _, v := range orders {
		j := convert(v)
		if j.CompletedAmount > 0 {
			j.Status = exchange.Partial
		}
		if err = c.Orders.UpdateOrder(j); err != nil {
			panic(err)
		}
	}
	orders, err = getTradeHistory(c.Key, c.Secret, nonce(), nil, &count)
	if err != nil {
		return
	}
	// Update only orders, that was completed it last refresh interval
	for _, orderID := range c.Orders.GetOpened() {
		for _, v := range orders {
			if v.OrderID == orderID {
				j := convert(v)
				j.Completed = time.Now()
				if err = c.Orders.UpdateOrder(j); err != nil {
					panic(err)
				}
			}
		}
	}

}

// Update starts update cycle
func (c *Client) Update() {
	var t = time.NewTicker(c.OrdersRefreshInterval)
	var bookt = time.NewTicker(c.OrderbookRefreshInterval)
	for {
		select {
		case <-t.C:
			c.updateOrders()
		case <-bookt.C:
			c.updateOrderbook()
		case <-c.Stop:
			t.Stop()
			bookt.Stop()
			return
		}
	}
}

// AddOrderbookTracking normalize and adds market to Client.TrackedBooks
func (c *Client) AddOrderbookTracking(market string) {
	if c.TrackedBooks == nil {
		c.TrackedBooks = make([]string, 0, 1)
	}
	market = normalize(market)
	c.TrackedBooks = append(c.TrackedBooks, market)
}

// RemoveOrderbookTracking delete markets from Client.TrackedBooks
func (c *Client) RemoveOrderbookTracking(market string) error {
	if c.TrackedBooks == nil || len(c.TrackedBooks) == 0 {
		return ErrNoOrderbooks
	}
	market = normalize(market)
	var tmp = make([]string, 0, len(c.TrackedBooks))
	for _, v := range c.TrackedBooks {
		if market != v {
			tmp = append(tmp, v)
		}
	}
	if len(c.TrackedBooks) == len(tmp) {
		return ErrOrderbookNotFound
	}
	c.TrackedBooks = tmp
	return nil
}

var errNotFound = errors.New("order with given orderid does not found")

func lookupOrder(key, secret string, orderID int) (Order, error) {
	orders, err := GetOpenOrders(key, secret, nil, nil)
	if err != nil {
		return Order{}, err
	}
	for _, v := range orders {
		if v.OrderID == orderID {
			return v, nil
		}
	}
	return Order{}, errNotFound
}

// ErrNoOrderbooks returns by Client.RemoveOrderbookTracking() if Client does not track any orderbooks
var ErrNoOrderbooks = errors.New("orderbooks not tracked")

// ErrOrderbookNotFound returns by Client.RemoveOrderbookTracking() if this market isnt tracked by Client
var ErrOrderbookNotFound = errors.New("this orderbook isnt tracked")
