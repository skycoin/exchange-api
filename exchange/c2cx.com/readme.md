 # c2cx.com API wrapper
  This package implements c2cx.com api and exchange.Client interface
  ## Public API functions:
   * GetOrderbook
  ## Private API functions:
   * CreateOrder
   * CancelOrder
   * GetBalance
   * GetOrderInfo
   * GetOrderByStatus 
  ## Client 
   Client track orderbooks and orders, that was created through it.
   Close Client.Stop for cancel updating.
 # Bugs:
  * CreateOrder does not working. I sent report to c2cx.com support