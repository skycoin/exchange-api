# exchange-api

Skycoin's exchange-api is a Go abstraction of cryptocoin trading exchanges. It abstract these four operations:
* placing a bid/ask order
* tracking the status of an existing order
* withdrawing bitcoin from the exchange
* depositing bitcoin to the exchange

It provides an internal Go API, a REST API, and a command line interface.

Currently it supports the `cryptopia` and `c2cx` exchanges. Additional exchanges can be implemented in Go.

You can view [the source code](https://github.com/skycoin/exchange-api) via github.

## Terminology

briancaine note:

  so, I added this section in case we might want to clarify any terminology.

  I'm guessing we're using standard terminology but

  maybe we might need to (for technical reasons) clarify what we mean

  (ie, like, does a transaction include fees or not? do we want to clarify that somewhere?)

  if not, then we can just delete this section

## Go

### Library API

briancaine note:

  here's where we document how you'd link it into another program as a library. straightforward enough

### Exchange API

The [`exchange` package](https://github.com/skycoin/exchange-api/tree/master/exchange) contains the necessary interfaces/types to implement support for a new exchange.

The [c2cx](https://github.com/skycoin/exchange-api/tree/master/exchange/c2cx.com) and [Cryptopia](https://github.com/skycoin/exchange-api/tree/master/exchange/cryptopia.co.nz) exchanges are already implemented and are useful as examples for future exchanges.

Implementing a new exchange only requires implementing the [Client interface](#client). The other datatypes are documented for reference purposes.

#### Client

The [`exchange.Client` interface](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/exchange/exchange.go#L30) is the primary interface to manipulating orders on the exchange.

```golang
type Client interface {
	// Cancel cancels one order by order id
	Cancel(int) (Order, error)
	// CancelMarket cancels all orders in given market
	CancelMarket(string) ([]Order, error)
	// CancelAll cancels all orders that executed in exchange
	CancelAll() ([]Order, error)
	// GetBalance gets a information about balance in a string format, depends of exchange representation format
	GetBalance(string) (string, error)
	// Buy places buy order
	Buy(string, float64, float64) (int, error)
	// Sell places sell order
	Sell(string, float64, float64) (int, error)
	// Completed gets completed orders
	Completed() []int
	// Executed gets opened orders
	Executed() []int
	// OrderStatus gets a string representation of order status
	// possible statuses defined below
	OrderStatus(int) (string, error)
	// OrderDetails gets detalied inforamtion about order with given order id
	OrderDetails(int) (Order, error)
	// Orderbook return Orderbooks interface
	Orderbook() Orderbooks
}
```

#### Orderbooks

The [`exchange.Orderbooks` interface](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/exchange/orderbooks.go#L9) provides access to an exchange's orderbooks, logically enough.

A [redis-backed implementation of Orderbooks](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/db/orderbooktracker.go) is provided.

```golang
type Orderbooks interface {
	// Update updates orderbook for given market
	Update(string, []MarketOrder, []MarketOrder)
	//Get gets orderbook for given tradepair symbol
	Get(string) (MarketRecord, error)
}
```

#### Order

[Source](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/exchange/exchange.go#L11)

```golang
type Order struct {
	Type      string
	Market    string
	Amount    float64
	Price     float64
	Submitted time.Time

	//Mutable fields
	OrderID         int
	Fee             float64
	CompletedAmount float64
	Status          string
	Accepted        time.Time
	Completed       time.Time
}
```

#### Order types

[Source](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/exchange/exchange.go#L5)

```golang
const (
	Buy  = "buy"
	Sell = "sell"
)
```

#### MarketOrder

A [MarketOrder](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/exchange/orderbooks.go#L16) is one order in stock.

```golang
type MarketOrder struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}
```

#### MarketRecord

A [MarketRecord](https://github.com/skycoin/exchange-api/blob/0b17f1aaf8967d3423495918ab350e290eaeafa8/exchange/orderbooks.go#L22) represents the order book for one market.

```golang
type MarketRecord struct {
	Timestamp time.Time     `json:"timestamp"`
	Symbol    string        `json:"symbol"`
	Bids      []MarketOrder `json:"bids"`
	Asks      []MarketOrder `json:"asks"`
}
```

## REST API

The REST API is based on [JSON-RPC 2.0](http://www.jsonrpc.org/specification).

### Types

The REST API accepts/returns simple JSON types, like strings and numbers, as well as complex structures.

Two reoccurring structures the API uses are:

* Order:
```json
{
  "orderid": 1,
  "type": "Buy",
  "market": "LTC/BTC",
  "amount": 100.123,
  "price": 0.0001,
  "submitted_at": 1501200244000,
  "fee": 0.001,
  "completed_amount": 99.9,
  "status": "Cancelled",
  "completed_at": 1501200244000,
  "accepted_at": 1501200244000
}
```
* Orderbook_record:
```json
{
  "timestamp": 1516255342,
  "symbol": "LTC/BTC",
  "asks": [{"price": 12345.67, "volume": 12.345},{"price": 23456.78, "volume": 23.456}],
  "bids": [{"price": 34567.89, "volume": 34.567},{"price": 14567.89, "volume": 543.21}]
}
```

### Methods

#### `buy`

Parameters: `{"symbol": string, "price": number, "amount", number}`

Returns: `number` (order ID)

#### `sell`

Parameters: `{"symbol": string, "price": number, "amount", number}`

Returns: `number` (order ID)

#### `cancel_trade`

Parameters: `{"orderid": number}` (order ID)

Returns: `Order`

#### `cancel_market`

Parameters: `{"symbol": string}`

Returns: `array of Order`

#### `cancel_all`

Parameters: `{}` or `null`

Returns: `array of Order`

#### `balance`

Parameters: `{"currency": string}`

Returns: `string`

#### `order_info`

Parameters: `{"orderid": number}` (order ID, obviously enough)

Returns: `Order`

#### `order_status`

Parameters: `{"orderid": number}` (order ID again)

Returns: `string`

#### `completed`

Parameters: `{}` or `null`

Returns: `array of number` (array of order IDs)

#### `executed`

Parameters: `{}` or `null`

Returns: `array of number` (array of order IDs)

#### `orderbook`

Parameters: `{"symbol": string}`

Returns: `Orderbook_record`

## CLI

A command line interface to the REST API is provided.

A CLI call takes the form: `cli <exchange> <command> [subcommand] [params...]`

Available CLI commands:
* `order`
  * `info <orderid>`
  * `status <orderid>`
* `cancel`
  * `all`
  * `market <symbol>`
  * `trade <orderid>`
* `buy <symbol> <price> <amount>`
* `sell <symbol> <price> <amount>`
* `orderbook <symbol>`
* `executed`
  * `market <symbol>`
  * `all`
* `completed`
  * `market <symbol>`
  * `all`
* `balance <currency>`
