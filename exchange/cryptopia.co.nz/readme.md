 #### cryptopia.co.nz API wrapper
 This package wraps cryptopia.co.nz API and implements exchange.Client interface 
 ## Public API functions: 
  * GetCurrencies
  * GetTradepairs
  * GetMarket
  * GetMarkets
  * GetMarketHistory
  * GetMarketOrders
  * GetMarketOrderGroups
 ## Private API functions:
  * GetBalance
  * SubmitTrade
  * CancelTrade
  * GetDepositAddress
  * GetOpenOrders
  * GetTradeHistory
  * GetTransactions
  * SubmitWithdraw
  * SubmitTransfer
  * SubmitTip
 ## Client
  Client implements order tracking for orders, that was created using it.
  Using Client.TrackedBooks you can set markets, where orderbooks will tracked.
  Close Client.Stop channel for stop updating