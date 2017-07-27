package c2cx

import (
	"github.com/pkg/errors"
	exchange "github.com/uberfurrer/tradebot/exchange"
)

// add new markets here
var markets = []string{"CNY_BTC", "CNY_ETH", "CNY_ETC", "CNY_SKY", "ETH_SKY", "BTC_SKY", "CNY_SHL"}

// GetOrderInfo gets extended information about orders with given tradepair ando orderID
// if orderID is -1, then all orders in given market will be returned
// if orders count greater than 100, use page parameter
func GetOrderInfo(key, secret string, symbol string, orderID int, page *int) ([]Order, error) {
	return getOrderinfo(key, secret, symbol, orderID, page)
}

// GetOrderByStatus gets extended information about orders in market with given status
// interval is time range between now and time after that you need get orders, in seconds
func GetOrderByStatus(key, secret string, symbol, status string, interval *int) ([]Order, error) {
	var i = -1
	if interval != nil {
		i = *interval
	}
	return getOrderByStatus(key, secret, symbol, status, i)

}

// GetBalance gets balance of all currencies
func GetBalance(key, secret string) (Balance, error) {
	return getBalance(key, secret)
}

// CancelOrder cancels order with given orderID, it return nil error if cancellation was success
func CancelOrder(key, secret string, orderID int) error {
	return cancelOrder(key, secret, orderID)
}

// AdvancedOrderParams is extended parameters, that may be used for setting stop loss,take profit and trigger price
type AdvancedOrderParams struct {
	TakeProfit   float64
	StopLoss     float64
	TriggerPrice float64
}

// CreateOrder creates new order with given parameters
// if adv == nil, then isAdvancedOrder will set to zero
// availible priceTypeIDs defined below
func CreateOrder(key, secret string, market string, price, quantity float64, orderType string, priceTypeID int, adv *AdvancedOrderParams) (int, error) {
	return createOrder(key, secret, market, price, quantity, orderType, priceTypeID, adv)
}

// GetOrderbook returns Orderbook with timestamp
// if symbol is not found,GetOrderbook also returns non-nli error
func GetOrderbook(symbol string) (Orderbook, error) {
	book, err := getOrderbook(symbol)
	if err != nil {
		return Orderbook{}, err
	}
	return *book, nil
}

// Order represents all information about order
type Order struct {
	Amount          float64
	AvgPrice        float64
	CompletedAmount float64
	Fee             float64
	CreateDate      int64
	CompleteDate    int64
	OrderID         int
	Price           float64
	Status          int
	Type            string
}

// Orderbook with timestamp
type Orderbook struct {
	Timestamp int                    `json:"timestamp"`
	Bids      []exchange.MarketOrder `json:"bids"`
	Asks      []exchange.MarketOrder `json:"asks"`
}

func apiError(endpoint, message string) error {
	return errors.Errorf("c2cx: %s falied, %s", endpoint, message)
}
