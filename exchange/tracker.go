package exchange

import (
	"errors"
	"time"
)

// NewOrder creates new Order with given params
func (tracker *OrderTracker) NewOrder(sym, action, status string, orderID int, vol, price float64) {
	switch status {
	case StatusCancelled, StatusCompleted:
		tracker.completed = append(tracker.completed, &OrderInfo{
			TradePair: sym,
			Type:      action,
			Status:    status,

			OrderID: orderID,
			Price:   price,
			Volume:  vol,

			Submitted: int64(time.Now().UnixNano() / 10e5),
			Accepted:  int64(time.Now().UnixNano() / 10e5),
			Completed: int64(time.Now().UnixNano() / 10e5),
		})
	default:
		tracker.executed[orderID] = &OrderInfo{
			TradePair: sym,
			Type:      action,
			Status:    status,

			OrderID: orderID,
			Price:   price,
			Volume:  vol,

			Submitted: int64(time.Now().UnixNano() / 10e5),
		}

	}
}

// UpdateOrderDetails updates info of order with given orderID, if acceptedAt == nil, only status will be updated
func (tracker *OrderTracker) UpdateOrderDetails(orderID int, status string, acceptedAt *time.Time) {
	if v, ok := tracker.executed[orderID]; ok {
		if acceptedAt != nil {
			v.Accepted = int64(acceptedAt.UnixNano() / 10e5)
		}
		v.Status = status
	}
}

// Complete sets time and move Order from executed orders
func (tracker *OrderTracker) Complete(orderID int, completedAt time.Time) {
	if v, ok := tracker.executed[orderID]; ok {
		v.Completed = int64(completedAt.UnixNano() / 10e5)
		v.Status = StatusCompleted
		tracker.completed = append(tracker.completed, v)

		delete(tracker.executed, orderID)
	}
}

// Cancel set comlpleted time as now and move Order form executed orders
func (tracker *OrderTracker) Cancel(orderIDs ...int) {
	for _, id := range orderIDs {
		if v, ok := tracker.executed[id]; ok {
			v.Completed = int64(time.Now().UnixNano() / 10e5)
			v.Status = StatusCancelled
			tracker.completed = append(tracker.completed, v)

			delete(tracker.executed, id)
		}
	}
}

// Completed returns currently completed or cancelled orders
func (tracker *OrderTracker) Completed() []*OrderInfo {
	return tracker.completed[:]
}

// Clear clear history of completed orders
func (tracker *OrderTracker) Clear() {
	tracker.completed = tracker.completed[:0]
}

// Status returns a string status of order with given orderID,
// if this order not found( after call flush, order was created from another client), Status returns non-nil error
func (tracker *OrderTracker) Status(orderID int) (string, error) {
	if v, ok := tracker.executed[orderID]; ok {
		return v.Status, nil
	}
	for _, v := range tracker.completed {
		if v.OrderID == orderID {
			return v.Status, nil
		}
	}
	return "", errors.New("Order not found")
}

// Get return detalied info of Order with given orderID
// return error as Status()
func (tracker *OrderTracker) Get(orderID int) (OrderInfo, error) {
	if v, ok := tracker.executed[orderID]; ok {
		return OrderInfo{
			TradePair: v.TradePair,
			Type:      v.Type,
			Status:    v.Status,
			OrderID:   v.OrderID,
			Price:     v.Price,
			Volume:    v.Volume,

			Submitted: v.Submitted,
			Accepted:  v.Accepted,
			Completed: v.Completed,
		}, nil
	}
	for _, v := range tracker.completed {
		if v.OrderID == orderID {
			return OrderInfo{
				TradePair: v.TradePair,
				Type:      v.Type,
				Status:    v.Status,
				OrderID:   v.OrderID,
				Price:     v.Price,
				Volume:    v.Volume,

				Submitted: v.Submitted,
				Accepted:  v.Accepted,
				Completed: v.Completed,
			}, nil
		}
	}

	return OrderInfo{}, errors.New("Order not found")
}

// Executed returns incompleted orders
func (tracker *OrderTracker) Executed() []*OrderInfo {
	var result = make([]*OrderInfo, 0, len(tracker.executed))
	for _, v := range tracker.executed {
		result = append(result, &OrderInfo{
			TradePair: v.TradePair,
			Status:    v.Status,
			Type:      v.Type,
			OrderID:   v.OrderID,
			Price:     v.Price,
			Volume:    v.Volume,

			Submitted: v.Submitted,
			Accepted:  v.Accepted,
			Completed: v.Completed,
		})
	}
	return result[:]
}

// NewTracker returns new OrderTracker instanse
func NewTracker() *OrderTracker {
	return &OrderTracker{
		executed:  make(map[int]*OrderInfo),
		completed: make([]*OrderInfo, 0),
	}
}
