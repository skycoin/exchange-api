package exchange

// Client describes interface that we could work with exchange for
// placing orders, check it status, etc...
type Client interface {
	Cancel(int) (*OrderInfo, error)
	CancelMarket(string) ([]*OrderInfo, error)
	CancelAll() ([]*OrderInfo, error)

	GetBalance(string) (string, error)
	Buy(string, float64, float64) (int, error)
	Sell(string, float64, float64) (int, error)
	Completed() []*OrderInfo
	Executed() []*OrderInfo
	OrderStatus(int) (string, error)
	OrderDetails(int) (OrderInfo, error)

	OrderBook() OrderBookTracker
}

// Balance is a total amount of funds on account
type Balance struct {
}

// OrderInfo type represents Order with tracking time and any additional info
// Each exchange must returns Orders in this format
type OrderInfo struct {
	TradePair string `json:"TradePair"`
	Type      string `json:"Type"`
	Status    string `json:"Status"`

	OrderID int     `json:"OrderID"`
	Price   float64 `json:"Price"`
	Volume  float64 `json:"Volume"`

	Submitted int64 `json:"Submitted"`
	Accepted  int64 `json:"Accepted"`
	Completed int64 `json:"Completed"`
}

// OrderTracker manages order statusees and track time
type OrderTracker struct {
	executed  map[int]*OrderInfo
	completed []*OrderInfo
}

// possible statusees of order
const (
	StatusOpened    = "Opened"
	StatusPartial   = "Partial"
	StatusCompleted = "Completed"
	StatusCancelled = "Cancelled"
	StatusSubmitted = "Submitted"
)

// possible OrderInfo actions
const (
	ActionBuy  = "buy"
	ActionSell = "sell"
)
