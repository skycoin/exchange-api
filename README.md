# Exchange-API

[![GoDoc](https://godoc.org/github.com/skycoin/exchange-api?status.svg)](https://godoc.org/github.com/skycoin/exchange-api)
[![Build Status](https://travis-ci.org/skycoin/exchange-api.svg?branch=master)](https://travis-ci.org/skycoin/exchange-api)

Exchange-API implements an interface to various cryptocurrency exchanges APIs in Go.

It can used as a library, or as a standalone JSON-RPC 2.0 server.

These two primary operations are abstracted across all 3rd party exchange interfaces:

* Placing a bid/ask order
* Tracking the status of an existing order

<!-- MarkdownTOC autolink="true" bracket="round" depth="5" -->

- [Server](#server)
    - [RPC API](#rpc-api)
        - [Types](#types)
            - [`Order`](#order)
            - [`OrderbookRecord`](#orderbookrecord)
        - [Methods](#methods)
            - [`buy`](#buy)
            - [`sell`](#sell)
            - [`cancel_trade`](#canceltrade)
            - [`cancel_market`](#cancelmarket)
            - [`cancel_all`](#cancelall)
            - [`balance`](#balance)
            - [`order_info`](#orderinfo)
            - [`order_status`](#orderstatus)
            - [`completed`](#completed)
            - [`executed`](#executed)
            - [`orderbook`](#orderbook)
    - [RPC Development](#rpc-development)
        - [Adding new RPC method](#adding-new-rpc-method)
- [CLI](#cli)
- [3 Additional C2CX commands](#3-additional-c2cx-commands)
    - [Additional Cryptopia commands](#additional-cryptopia-commands)
    - [Usage examples](#usage-examples)
        - [Getting order information by orderid](#getting-order-information-by-orderid)
        - [Getting order status by orderid](#getting-order-status-by-orderid)
        - [Create new order](#create-new-order)
        - [Get open orders](#get-open-orders)
        - [Get open orders by market](#get-open-orders-by-market)
        - [Cancel orders by market](#cancel-orders-by-market)
        - [Get balance](#get-balance)
- [Integration Tests](#integration-tests)

<!-- /MarkdownTOC -->


## Server

`node` is the main executable.
`node` needs exchange API key and secret in format: "key:secret" for initilaizing the client.
The client is interface that handles the major methods of each exchange.

Flags:

* `-srv`: RPC bind address

Client functionality:

- Place buy and sell orders
- Getting order information and status
- Getting all completed and executed orders, that was created using this client
- Getting orderbook for given market
- Cancel orders by orderid and market or cancel all orders
- Getting balance

`node` uses jsonrpc 2.0 protocol for RPC calls

### RPC API

#### Types

The RPC API accepts/returns simple JSON types, like strings and numbers, as well as complex structures.

Two reoccurring structures the API uses are:

##### `Order`

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

##### `OrderbookRecord`

```json
{
  "timestamp": 1516255342,
  "symbol": "LTC/BTC",
  "asks": [{"price": 12345.67, "volume": 12.345},{"price": 23456.78, "volume": 23.456}],
  "bids": [{"price": 34567.89, "volume": 34.567},{"price": 14567.89, "volume": 543.21}]
}
```

#### Methods

##### `buy`

Parameters: `{"symbol": string, "price": number, "amount", number}`

Returns: `number` (order ID)

##### `sell`

Parameters: `{"symbol": string, "price": number, "amount", number}`

Returns: `number` (order ID)

##### `cancel_trade`

Parameters: `{"orderid": number}` (order ID)

Returns: `Order`

##### `cancel_market`

Parameters: `{"symbol": string}`

Returns: `array of Order`

##### `cancel_all`

Parameters: `{}` or `null`

Returns: `array of Order`

##### `balance`

Parameters: `{"currency": string}`

Returns: `string`

##### `order_info`

Parameters: `{"orderid": number}` (order ID, obviously enough)

Returns: `Order`

##### `order_status`

Parameters: `{"orderid": number}` (order ID again)

Returns: `string`

##### `completed`

Parameters: `{}` or `null`

Returns: `array of number` (array of order IDs)

##### `executed`

Parameters: `{}` or `null`

Returns: `array of number` (array of order IDs)

##### `orderbook`

Parameters: `{"symbol": string}`

Returns: `Orderbook_record`

### RPC Development

Package `rpc` uses JSON-RPC 2.0.

`PackageHandler` handles the `exchange.Client` interface and allows one to add additional,
exchange-specific handlers to each method.
You can add any variables that you need using `Env` and the `PackageHandler.Setenv()` function.

#### Adding new RPC method

```go
PackageHandler.Handlers["method_name"] = func(r Request, env map[string]string) (Response) {
    var result interface{}
    var err error
    //some actions
    if err != nil {
        return MakeErrorResponse(r, InternalError, err)
    }
    return MakeSuccessResponse(r, result)
}
```

Each function should check errors and return an empty body and non-empty error field on error.
Use `MakeErrorResponse()` and `MakeSuccessResponse()` for this.

## CLI

A command line interface to the RPC API is provided.

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

A specific exchange may have additional functions:

##3 Additional C2CX commands

- `submittrade <market> <price> <amount>`
 This command sets additional order params.
 Require: `pricetype` with a value of `limit` or `market` and `ordertype` with a value of `buy` or `sell`
 Optional: `triggerprice`, `takeprofit` and `stoploss`

### Additional Cryptopia commands

- `deposit <currency>`
 Gets the deposit address for given currency
- `withdraw <address> <currency> <amount> [paymentid]`
 Creates withdrawal request, address should be added to list of allowed addresses
 `paymentid` is needed only for CryptoNote-based cryptocurrencies
- `transactions <type>`
 `type` must be either `deposit` or `withdraw`
- `tracking add <symbol>`
 Add symbol to orderbook tracking list
- `tracking rm <symbol>`
 Remove symbol from orderbook tracking list

### Usage examples

#### Getting order information by orderid

```sh
> ./cli c2cx order info 2832693
{
    "amount": 100,
    "market": "CNY/SHL",
    "orderid": 2832693,
    "price": 0.01,
    "status": "submitted",
    "submitted_at": "2017-08-05T06:54:12.469+07:00",
    "type": "buy"
}
```

#### Getting order status by orderid

```sh
> ./cli c2cx order status 2832693
Order 2832693 status: "submitted"
```

#### Create new order

first parameter is marketpair, second - price, third - amount

```sh
./cli c2cx buy cny/shl 0.01 100
Order 2832694 created
```

#### Get open orders

```sh
./cli c2cx executed all
{
    "amount": 100,
    "market": "BTC/BCC",
    "orderid": 2832693,
    "price": 0.000011,
    "status": "submitted",
    "submitted_at": "2017-08-05T06:54:12.469+07:00",
    "type": "buy"
}
{
    "amount": 100,
    "market": "CNY/SHL",
    "orderid": 2832694,
    "price": 0.01,
    "status": "submitted",
    "submitted_at": "2017-08-05T07:30:14.063+07:00",
    "type": "buy"
}
```

#### Get open orders by market

```sh
./cli c2cx executed market cny/shl
{
    "amount": 100,
    "market": "CNY/SHL",
    "orderid": 2832694,
    "price": 0.01,
    "status": "submitted",
    "submitted_at": "2017-08-05T07:30:14.063+07:00",
    "type": "buy"
}
```

#### Cancel orders by market

```sh
./cli c2cx cancel market cny/shl
{
"amount": 10,
"market": "CNY/SHL",
"orderid": 2832707,
"price": 0.01
}
{
"amount": 100,
"market": "CNY/SHL",
"orderid": 2832708,
"price": 0.01
}
Cancelled 2 orders
```

#### Get balance

```sh
./cli c2cx balance btc
"Availible 0.00000000, frozen 0.00000000"
```

## Integration Tests

To run the integration tests for the C2CX API:
1. Obtain a [C2CX account](https://www.c2cx.com) and deposit at least 1.2 SKY.
2. Create an [API key](https://www.c2cx.com/in/myaccount/api)
3. Set the environment variables `C2CX_TEST_KEY` and `C2CX_TEST_SECRET` to the key/secret from step 2.
4. Run `make test`.

Example:

```sh
C2CX_TEST_KEY=ABABABAB-ABAB-ABAB-ABAB-ABABABABABAB C2CX_TEST_SECRET=CDCDCDCD-CDCD-CDCD-CDCD-CDCDCDCDCDCD make test
```

The tests test (among other things) placing orders, retrieving orders and canceling orders.

The tests will place an order to sell 1.2 SKY at a rate of 0.5 BTC/SKY. Then it will test retrieving the order information. Then it will cancel the order.

Since we probably won't hit 0.5 BTC/SKY anytime soon, there's not much concern about losing any funds.
However, if the tests fail halfway, you may need to [manually cancel orders in C2CX](https://www.c2cx.com/in/orders) to run the tests again.


