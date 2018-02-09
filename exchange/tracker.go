package exchange

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"errors"

	"github.com/shopspring/decimal"
)

// Orders is interface that provides functionality for order tracking
type Orders interface {
	GetOpened() []int
	GetCompleted() []int
	GetOrderInfo(int) (Order, error)
	GetOrderStatus(int) (string, error)

	UpdateOrder(Order) error
	Push(Order) error
}

// Errors that might happens in order processing
var (
	ErrInvalidStatus = fmt.Errorf("invalid order status, available statuses: %s, %s, %s, %s, %s",
		Submitted, Opened, Partial, Completed, Cancelled)
	ErrExist       = errors.New("order with given orderid already exist")
	ErrZeroOrderID = errors.New("order with zero orderID does not allowed")
	ErrNotFound    = errors.New("order does not found")
	ErrNotTracked  = errors.New("tracker does not present")
)

type idx struct {
	hash   int
	status int
}
type tracker struct {
	orders    map[int]idx
	opened    map[int]Order
	completed map[int]Order
}

func (t *tracker) GetOpened() (orders []int) {
	orders = make([]int, 0, len(t.opened))
	for _, v := range t.opened {
		orders = append(orders, v.OrderID)
	}
	return orders
}

func (t *tracker) GetCompleted() (orders []int) {
	orders = make([]int, 0, len(t.completed))
	for _, v := range t.completed {
		orders = append(orders, v.OrderID)
	}
	return orders
}

func (t *tracker) GetOrderInfo(orderid int) (Order, error) {
	return t.lookupByOrderID(orderid)
}

func (t *tracker) GetOrderStatus(orderid int) (string, error) {
	if v, err := t.lookupByOrderID(orderid); err == nil {
		return v.Status, nil
	}
	return "", ErrNotFound
}

const (
	statusOpen      = 0
	statusCompleted = 1
)

func (t *tracker) lookupByOrderID(orderid int) (Order, error) {
	if lookupInfo, ok := t.orders[orderid]; ok {
		switch lookupInfo.status {
		case statusOpen:
			return t.opened[lookupInfo.hash], nil
		case statusCompleted:
			return t.completed[lookupInfo.hash], nil
		}
	}
	return Order{}, ErrNotFound
}

func (t *tracker) UpdateOrder(order Order) error {
	if lookupData, ok := t.orders[order.OrderID]; ok {
		switch lookupData.status {
		case statusOpen:
			switch order.Status {
			case Opened, Partial, Submitted:
				t.opened[lookupData.hash] = update(t.opened[lookupData.hash], order)
			case Completed, Cancelled:
				exist := t.opened[lookupData.hash]
				delete(t.opened, lookupData.hash)
				t.orders[order.OrderID] = idx{
					hash:   lookupData.hash,
					status: statusCompleted,
				}
				t.completed[lookupData.hash] = update(exist, order)
			}
		case statusCompleted:
			t.completed[lookupData.hash] = update(t.completed[lookupData.hash], order)
		}
		return nil
	}
	// if not found in t.orders
	// it happens if OrderID of order was changed, but it is same order
	// update lookup data and repeat
	// old OrderID also saved
	var hashValue = hash(order)
	if v, ok := t.opened[hashValue]; ok {
		lookupData := t.orders[v.OrderID]
		t.orders[order.OrderID] = lookupData
		err := t.UpdateOrder(order)
		if err != nil {
			return err
		}
		t.orders[v.OrderID] = t.orders[order.OrderID]
		return nil
	}
	if v, ok := t.completed[hashValue]; ok {
		lookupData := t.orders[v.OrderID]
		t.orders[order.OrderID] = lookupData
		err := t.UpdateOrder(order)
		if err != nil {
			return err
		}
		t.orders[v.OrderID] = t.orders[order.OrderID]
		return nil
	}
	return ErrNotFound
}

func update(exist, upd Order) Order {
	var (
		zerotime   = time.Time{}
		zerostring = ""
	)
	if !upd.CompletedAmount.Equal(decimal.Zero) {
		exist.CompletedAmount = upd.CompletedAmount
	}
	if upd.Status != zerostring {
		exist.Status = upd.Status
	}
	if !upd.Fee.Equal(decimal.Zero) {
		exist.Fee = upd.Fee
	}
	if upd.Accepted != zerotime {
		exist.Accepted = upd.Accepted
	}
	if upd.Completed != zerotime {
		exist.Completed = upd.Completed
	}
	return exist
}

func (t *tracker) Push(order Order) error {
	if order.OrderID == 0 {
		return ErrZeroOrderID
	}
	if _, ok := t.orders[order.OrderID]; ok {
		return ErrExist
	}
	order.Market = normalize(order.Market)
	order.Status = strings.ToLower(order.Status)
	var hashValue = hash(order)
	var posData = idx{
		hash: hashValue,
	}
	switch order.Status {
	case Submitted, Opened, Partial:
		posData.status = statusOpen
		t.opened[hashValue] = order
	case Completed, Cancelled:
		posData.status = statusCompleted
		t.completed[hashValue] = order
	default:
		return ErrInvalidStatus
	}
	t.orders[order.OrderID] = posData
	return nil

}

// NewTracker returns internal tracker instance
func NewTracker() Orders {
	return &tracker{
		opened:    make(map[int]Order),
		completed: make(map[int]Order),
		orders:    make(map[int]idx),
	}
}

func normalize(sym string) string {
	return strings.ToUpper(strings.Replace(sym, "_", "/", -1))
}

func hash(order Order) int {
	order = truncate(order)
	var (
		buf    = make([]byte, 8*2)
		amount = uint64(order.Amount.Mul(decimal.New(10, 8)).IntPart())
		price  = uint64(order.Price.Mul(decimal.New(10, 8)).IntPart())
	)
	binary.BigEndian.PutUint64(buf[0:8], amount)
	binary.BigEndian.PutUint64(buf[8:16], price)
	buf = append(buf, order.Accepted.String()...)
	buf = append(buf, order.Type...)
	hash := sha256.Sum256(buf[:])
	return int(binary.BigEndian.Uint64(hash[0:8]))
}
