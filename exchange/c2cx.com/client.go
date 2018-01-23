package c2cx

import (
	"log"
	"time"

	"strings"

	"github.com/pkg/errors"
	"github.com/skycoin/exchange-api/exchange"
)

// Client implements exchange.Client interface
// Client track all orders that was created using it
type Client struct {
	// Key and Secret needs for creating and accessing orders, update them
	// You may use Client without it for tracking OrderBook
	Key, Secret              string
	OrdersRefreshInterval    time.Duration
	OrderbookRefreshInterval time.Duration

	// Tracker provides provides functionality for tracking orders
	// if Tracker == nil then orders does not tracked and Client will be update only Orderbook
	// For properly work it requires that GetOrderByStatus with given key and secret executes without error
	Orders exchange.Orders

	// OrderBookTracker provides functionality for tracking OrderBook
	// It use RefreshRate in milliseconds for updating
	// OrderBookTracker should be free for concurrent use
	Orderbooks exchange.Orderbooks

	// Stop stops updating
	// After sending to this, you need to restart Client.Update()
	Stop chan struct{}

	prevUpdate time.Time
}

// Cancel cancels order with given orderID
func (c *Client) Cancel(orderID int) (exchange.Order, error) {
	err := cancelOrder(c.Key, c.Secret, orderID)
	if err != nil {
		return exchange.Order{}, err
	}
	order, err := c.Orders.GetOrderInfo(orderID)
	if err != nil {
		return exchange.Order{}, err
	}
	orders, err := getOrderinfo(c.Key, c.Secret, order.Market, order.OrderID, nil)
	if err != nil {
		return exchange.Order{}, err
	}
	if len(orders) != 1 {
		return exchange.Order{}, errors.New("order does not found")
	}
	var completedTime time.Time
	log.Println(orders[0].CompleteDate)
	if orders[0].CompleteDate != 0 {
		completedTime = unix(orders[0].CompleteDate)
	} else {
		completedTime = time.Now()
	}
	c.Orders.UpdateOrder(
		exchange.Order{
			OrderID:         orderID,
			Price:           orders[0].Price,
			Amount:          orders[0].Amount,
			Status:          exchange.Cancelled,
			Completed:       completedTime,
			Accepted:        unix(orders[0].CreateDate),
			Fee:             orders[0].Fee,
			CompletedAmount: orders[0].CompletedAmount,
		})
	return c.Orders.GetOrderInfo(orderID)

}

// CancelAll cancels all executed orders, that was created using this cilent
func (c *Client) CancelAll() ([]exchange.Order, error) {
	var (
		orderids = c.Orders.GetOpened()
		orders   = make([]exchange.Order, 0, len(orderids))
	)
	for _, v := range orderids {
		order, err := c.Cancel(v)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}

// CancelMarket cancels all order with given symbol that was created using this client
func (c *Client) CancelMarket(symbol string) ([]exchange.Order, error) {
	var (
		orderids []int
		orders   []exchange.Order
	)
	symbol = strings.ToUpper(strings.Replace(symbol, "_", "/", -1))
	for _, v := range c.Orders.GetOpened() {
		order, err := c.Orders.GetOrderInfo(v)
		if err != nil {
			continue
		}
		if order.Market == symbol {
			orderids = append(orderids, order.OrderID)
		}
	}
	orders = make([]exchange.Order, 0, len(orderids))
	var rejected []int
	for _, v := range orderids {
		order, err := c.Cancel(v)
		if err != nil {
			rejected = append(rejected, v)
			continue
		}
		orders = append(orders, order)
	}
	if rejected != nil {
		return orders, errors.Errorf("this orders does not cancelled: %v", rejected)
	}
	return orders, nil

}

// Buy place buy order
func (c *Client) Buy(symbol string, price, amount float64) (orderID int, err error) {
	var order = exchange.Order{
		Submitted: time.Now(),
		Market:    symbol,
		Price:     price,
		Amount:    amount,
		Type:      exchange.Buy,
		Status:    exchange.Submitted,
	}
	orderID, err = c.createOrder(symbol, price, amount, "buy")
	if err != nil {
		return
	}

	orders, err := getOrderinfo(c.Key, c.Secret, symbol, orderID, nil)
	if err != nil {
		return
	}
	if len(orders) != 1 {
		return /// error update info
	}
	order.Accepted = convert(orders[0]).Accepted
	order.OrderID = orderID
	err = c.Orders.Push(order)
	return
}

// Sell place sell order
func (c *Client) Sell(symbol string, price, amount float64) (orderID int, err error) {
	var order = exchange.Order{
		Submitted: time.Now(),
		Market:    symbol,
		Price:     price,
		Amount:    amount,
		Type:      exchange.Sell,
		Status:    exchange.Submitted,
	}
	orderID, err = c.createOrder(symbol, price, amount, "sell")
	if err != nil {
		return
	}
	order.OrderID = orderID
	err = c.Orders.Push(order)

	return
}

func (c *Client) createOrder(symbol string, price, quantity float64, Type string) (int, error) {
	var err error
	if symbol, err = normalize(symbol); err != nil {
		return 0, err
	}
	return createOrder(c.Key, c.Secret, symbol, price, quantity, Type, PriceTypeLimit, nil)
}

// OrderStatus returns string status of order with given orderID
// Handles only orders that was created using this client
func (c *Client) OrderStatus(orderID int) (string, error) {
	return c.Orders.GetOrderStatus(orderID)

}

// OrderDetails returns all avalible info about order
// Handles only orders that was created using this client
func (c *Client) OrderDetails(orderID int) (exchange.Order, error) {
	return c.Orders.GetOrderInfo(orderID)

}

// Executed wraps Tracker.Executed()
func (c *Client) Executed() []int {
	return c.Orders.GetOpened()
}

// Completed wraps Tracker.Completed()
func (c *Client) Completed() []int {
	return c.Orders.GetCompleted()
}

// Orderbook returns interface for managing Orderbook
func (c *Client) Orderbook() exchange.Orderbooks {
	return c.Orderbooks
}

// GetBalance gets balance information about given currency
func (c *Client) GetBalance(currency string) (string, error) {
	info, err := getBalance(c.Key, c.Secret)
	if err != nil {
		return "", err
	}
	if v, ok := info[strings.ToLower(currency)]; ok {
		return v, nil
	}
	return "", errors.Errorf("currency %s does not found", currency)
}
func (c *Client) updateOrderbook() {
	for _, v := range Markets {
		orderbook, err := getOrderbook(v)
		if err != nil {
			continue
		}
		c.Orderbooks.Update(v, orderbook.Bids, orderbook.Asks)
	}
}
func (c *Client) updateOrders() {
	for _, v := range Markets {
		orders, err := getOrderinfo(c.Key, c.Secret, v, -1, nil)
		if err != nil {
			continue
		}
		for _, v := range orders {
			t := convert(v)
			if err := c.Orders.UpdateOrder(t); err != nil {
				continue
			}
		}
	}
}

// Update starts update cycle
func (c *Client) Update() {
	bookt := time.NewTicker(c.OrderbookRefreshInterval)
	t := time.NewTicker(c.OrdersRefreshInterval)
	for {
		select {
		case <-bookt.C:
			c.updateOrderbook()
		case <-t.C:
			c.updateOrders()
		case <-c.Stop:
			t.Stop()
			bookt.Stop()
			return
		}
	}
}
