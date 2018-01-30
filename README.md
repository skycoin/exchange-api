
# Node
 Node is main executable
 Node needs exchange API key and secret in format: "key:secret" for initilaize client,
 client is interface that handles major methods of each exchange  

 -srv flag define rpc addres for bind  
 -db flag define redis address  

 Client functionality: 
  - Place buy and sell orders
  - Getting order information and status
  - Getting all completed and executed orders, that was created using this client 
  - Getting orderbook for given market
  - Cancel orders by orderid and market or cancel all orders
  - Getting balance

 Node using jsonrpc 2.0 protocol for rpc calls:

## RPC methods:

More details on the RPC interface are available in [the RPC README file](https://github.com/skycoin/exchange-api/tree/master/rpc#readme).

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

More details on the CLI are available in [the CLI README file](https://github.com/skycoin/exchange-api/tree/master/cli#readme).

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
