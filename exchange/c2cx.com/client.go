package c2cx

import (
	"fmt"
	"time"

	"strings"

	"errors"

	"github.com/shopspring/decimal"

	"github.com/skycoin/exchange-api/exchange"
)

// Client implements exchange.Client interface
// Client track all orders that was created using it
type Client struct {
	// Key and Secret needs for creating and accessing orders, update them
	// You may use Client without it for tracking OrderBook
	Key                      string
	Secret                   string
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
	if len(orders) == 0 {
		return exchange.Order{}, errors.New("order was not found")
	}

	var completedTime time.Time
	if orders[0].CompleteDate != 0 {
		completedTime = unix(orders[0].CompleteDate)
	} else {
		completedTime = time.Now()
	}

	if err := c.Orders.UpdateOrder(exchange.Order{
		OrderID:         orderID,
		Price:           orders[0].Price,
		Amount:          orders[0].Amount,
		Status:          exchange.Cancelled,
		Completed:       completedTime,
		Accepted:        unix(orders[0].CreateDate),
		Fee:             orders[0].Fee,
		CompletedAmount: orders[0].CompletedAmount,
	}); err != nil {
		return exchange.Order{}, err
	}

	return c.Orders.GetOrderInfo(orderID)
}

// CancelMultiError is returned when an error was encountered while cancelling multiple orders
type CancelMultiError struct {
	OrderIDs []int
	Errors   []error
}

func (e CancelMultiError) Error() string {
	return fmt.Sprintf("these orders failed to cancel: %v", e.OrderIDs)
}

// CancelAll cancels all executed orders, that was created using this client.
// If it encounters an error, it aborts and returns the orders that had been cancelled to that point.
func (c *Client) CancelAll() ([]exchange.Order, error) {
	orderIDs := c.Orders.GetOpened()
	return c.CancelMultiple(orderIDs)
}

// CancelMarket cancels all order with given symbol that was created using this client
func (c *Client) CancelMarket(symbol string) ([]exchange.Order, error) {
	var orderIDs []int
	// TODO -- symbol transformation should be a utility method
	symbol = strings.ToUpper(strings.Replace(symbol, "_", "/", -1))
	for _, v := range c.Orders.GetOpened() {
		order, err := c.Orders.GetOrderInfo(v)
		if err != nil {
			continue
		}
		if order.Market == symbol {
			orderIDs = append(orderIDs, order.OrderID)
		}
	}

	return c.CancelMultiple(orderIDs)
}

// CancelMultiple cancels multiple orders.  It will try to cancel all of them, not
// stopping for any individual error. If any orders failed to cancel, a CancelMultiError is returned
// along with the array of orders which successfully cancelled.
func (c *Client) CancelMultiple(orderIDs []int) ([]exchange.Order, error) {
	var orders []exchange.Order
	var cancelErr CancelMultiError

	for _, v := range orderIDs {
		order, err := c.Cancel(v)
		if err != nil {
			cancelErr.OrderIDs = append(cancelErr.OrderIDs, v)
			cancelErr.Errors = append(cancelErr.Errors, err)
			continue
		}
		orders = append(orders, order)
	}

	if len(cancelErr.OrderIDs) != 0 {
		return orders, &cancelErr
	}

	return orders, nil
}

// Buy place buy order
func (c *Client) Buy(symbol string, price, amount decimal.Decimal) (int, error) {
	order := exchange.Order{
		Submitted: time.Now(),
		Market:    symbol,
		Price:     price,
		Amount:    amount,
		Type:      exchange.Buy,
		Status:    exchange.Submitted,
	}
	orderID, err := c.createOrder(symbol, price, amount, "buy")
	if err != nil {
		return 0, err
	}

	orders, err := getOrderinfo(c.Key, c.Secret, symbol, orderID, nil)
	if err != nil {
		return 0, err
	}
	if len(orders) == 0 {
		return 0, errors.New("order not found")
	}

	order.Accepted = convert(orders[0]).Accepted
	order.OrderID = orderID

	if err := c.Orders.Push(order); err != nil {
		return 0, err
	}

	return orderID, nil
}

// Sell place sell order
func (c *Client) Sell(symbol string, price, amount decimal.Decimal) (int, error) {
	order := exchange.Order{
		Submitted: time.Now(),
		Market:    symbol,
		Price:     price,
		Amount:    amount,
		Type:      exchange.Sell,
		Status:    exchange.Submitted,
	}

	orderID, err := c.createOrder(symbol, price, amount, "sell")
	if err != nil {
		return 0, err
	}

	order.OrderID = orderID
	if err := c.Orders.Push(order); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (c *Client) createOrder(symbol string, price, quantity decimal.Decimal, Type string) (int, error) {
	symbol, err := normalize(symbol)
	if err != nil {
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
func (c *Client) GetBalance(currency string) (decimal.Decimal, error) {
	info, err := getBalance(c.Key, c.Secret)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if result, ok := info[strings.ToLower(currency)]; ok {
		return result, nil
	}

	return decimal.Decimal{}, fmt.Errorf("currency %s was not found", currency)
}

func (c *Client) updateOrderbook() {
	for _, v := range Markets {
		orderbook, err := getOrderbook(v)
		if err != nil {
			// TODO -- log or return error?
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
