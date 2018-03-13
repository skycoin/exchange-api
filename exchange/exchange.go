package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

// Order types
const (
	Buy  = "buy"
	Sell = "sell"
)

// Order is a normalized order
type Order struct {
	Type      string
	Market    string
	Amount    decimal.Decimal
	Price     decimal.Decimal
	Submitted time.Time

	//Mutable fields
	OrderID         int
	Fee             decimal.Decimal
	CompletedAmount decimal.Decimal
	Status          string
	Accepted        time.Time
	Completed       time.Time
}

// Client provides functionality for placing orders,
// gets orders info and statuses and gets balance from exchange
type Client interface {
	// Cancel cancels one order by order id
	Cancel(int) (Order, error)
	// CancelMarket cancels all orders in given market
	CancelMarket(string) ([]Order, error)
	// CancelAll cancels all orders that executed in exchange
	CancelAll() ([]Order, error)
	// GetBalance returns the balance for a given currency
	GetBalance(string) (decimal.Decimal, error)
	// Buy places buy order
	Buy(string, decimal.Decimal, decimal.Decimal) (int, error)
	// Sell places sell order
	Sell(string, decimal.Decimal, decimal.Decimal) (int, error)
	// Completed gets completed orders
	Completed() []int
	// Executed gets opened orders
	Executed() []int
	// OrderStatus gets a string representation of order status
	// possible statuses defined below
	OrderStatus(int) (string, error)
	// OrderDetails gets detailed information about order with given order id
	OrderDetails(int) (Order, error)
	// GetMarketRecord gets the orderbook for a given tradepair symbol
	GetMarketRecord(string) (*MarketRecord, error)
}

// Statuses
const (
	Submitted = "submitted"
	Opened    = "opened"
	Partial   = "partial"
	Completed = "completed"
	Cancelled = "cancelled"
)
