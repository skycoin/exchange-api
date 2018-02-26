[![GoDoc](https://godoc.org/github.com/skycoin/exchange-api?status.svg)](https://godoc.org/github.com/skycoin/exchange-api)
[![Build Status](https://travis-ci.org/skycoin/exchange-api.svg?branch=master)](https://travis-ci.org/skycoin/exchange-api)

# Node
 Node is main executable
 Node needs exchange API key and secret in format: "key:secret" for initilaize client,
 client is interface that handles major methods of each exchange  

 -srv flag define rpc addres for bind

 Client functionality: 
  - Place buy and sell orders
  - Getting order information and status
  - Getting all completed and executed orders, that was created using this client 
  - Getting orderbook for given market
  - Cancel orders by orderid and market or cancel all orders
  - Getting balance

 Node using jsonrpc 2.0 protocol for rpc calls:

## RPC methods:

More details on the RPC interface are available in [the RPC README file](rpc/readme.md).

  - `buy`
  - `sell`
  - `order_info`
  - `order_status`
  - `cancel_all`
  - `cancel_market`
  - `cancel_trade`
  - `completed`
  - `executed`
  - `orderbook`

  and you can add specific handlers for each exchange
  rpc package contains more documentation - methods, params structures for each call, return values and description for each method
  rpc endpoints for each exchange are different -  
  http://rpcaddr:rpcport/c2cx and http://rpcaddr:rpcport/cryptopia, respectively 
# Cli

 Cli command structure:  
 `cli <exchange> <command> [subcommand] [params]`  
 Commands tree: 
 - `order`  
   - `info <orderid>`
   - `status <orderid>`
 - `cancel`  
   - `all`
   - `market <symbol>`
   - `trade <orderid>`
 - `buy <symbol> <price> <amount>`
 - `sell <symbol> <price> <amount>`
 - `orderbook <symbol>`
 - `executed`  
   - `market <symbol>`
   - `all`
 - `completed`  
   - `market <symbol>`
   - `all`
 - `balance <currency>`

 Exchange may contains additional functions, that does not handles client interface  
  c2cx commands:  
   - `submittrade <market> <price> <amount>`  
     this command allow to set additional order params  
     required flags - pricetype `limit/market` and ordertype `buy/sell`  
     additional flags - triggerprice, takeprofit and stoploss

  cryptopia commands:
   - `deposit <currency>`  
   gets deposit address for given currency
   - `withdraw <address> <currency> <amount> [paymentid]`  
   creates withdrawal request, address should be added to list of allowed addresses   
   paymentid needs only for CryptoNote based cryptocurrencies
   - `transactions <type>`  
   type `deposit/withdraw`
   - `tracking`
     - `add <symbol>`  
     add symbol to orderbook tracking list
     - `rm <symbol>`  
     remove symbol from ordebook tracking list

 ### Usage examples: 
  - Getting order information by orderid  
    ```
     ./cli c2cx order info 2832693
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
  - Getting order status by orderid  
    ```
     ./cli c2cx order status 2832693
     Order 2832693 status: "submitted"
    ```
  - Create new order  
    first parameter is marketpair, second - price, third - amount  
    ```
     ./cli c2cx buy cny/shl 0.01 100
     Order 2832694 created
    ```
  - Get open orders  
    ```
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
  - Get open orders by market 
    ```
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
  - Cancel orders by market  
    ```
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
  - Getting balance  
    ```
     ./cli c2cx balance btc
     "Availible 0.00000000, frozen 0.00000000"
    ```
  #### TODO:
   * Components interaction tests

## Integration Tests

To run the integration tests for the C2CX API:
1. Obtain a [C2CX account](https://www.c2cx.com) and deposit at least 1.2 SKY.
2. Create an [API key](https://www.c2cx.com/in/myaccount/api)
3. Set the environment variables `C2CX_TEST_KEY` and `C2CX_TEST_SECRET` to the key/secret from step 2.
4. Run `make test`.

    Example: `C2CX_TEST_KEY=ABABABAB-ABAB-ABAB-ABAB-ABABABABABAB C2CX_TEST_SECRET=CDCDCDCD-CDCD-CDCD-CDCD-CDCDCDCDCDCD make test`

The tests test (among other things) placing orders, retrieving orders and canceling orders.

The tests will place an order to sell 1.2 SKY at a rate of 0.5 BTC/SKY. Then it will test retrieving the order information. Then it will cancel the order.

Since we probably won't hit 0.5 BTC/SKY anytime soon, there's not much concern about losing any funds. However, if the tests fail halfway, you may need to [manually cancel orders in C2CX](https://www.c2cx.com/in/orders) to run the tests again.


