# Common exchange client interface
  * Each exchange should implement client interface
  * RPC handler works with this interface
 ## Client requirements
  * Client should track orders and orderbook
  * Client should allow to set refresh interval
  * Client should convert order from exchange format to exchange.Order
 ## Orders
  Orders is interface for order tracking, it allows to:
  * Get open orders
  * Get competed orders
  * Get orderinfo by orderID
  * Push new orders to store
  * Update order
  #### Updating order: 
  Order has mutable and immutable fields

  Immutable fields is used for tracking orders, cause OrderID may be different in opened and completed order.

  Before calling `Orders.UpdateOrder(Order)` Order should has proprely filled all mutable and immutable fields
  
  Immutable fields are: 
   - Type
   - Price
   - Amount
   - Timestamp
 ## Orderbooks
  Orderbooks is interface for orderbooks tracking, it wraps underlying DB or other store. Client manipulates orderbook through this interface.
  
  MarketRecord is a orderbook record for one tradepair symbol. MarketRecord also contains timestamp, truncated to milliseconds and tradepair label.
