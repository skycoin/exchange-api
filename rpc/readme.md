# RPC interface
 ## Overview
  package rpc use JSONRPC 2.0 scheme  
  PackageHandler handles exchange.Client interface and allow to add additional, exchange specific handlers to each functions, that does not supported for this interface. You can add any variables that you need using Env map and `PackageHandler.Setenv()` function.  
  Adding new function:  
   ``` 
    PackageHandler.Handlers["method_name"] = func(r Request, env map[string]string) (Response) {  
        var result interface{} 
        var err error 
        //some actions  
        if err != nil {
            return MakeErrorResponse(r, InternalError, err)
        }
        return MakeSuccesResponse(r, result)
    }
   ```
   Each function should check errors and return empty body and non-empty error field on error  
   Use `MakeErrorResponse()` and `MakeSuccessResponse()` for this
  
 ## exchange.Client interface:
  #### Datatypes:
   * order  
   ``` 
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
   * orderbook_record  
   ```
     {  
         "timestamp": integer,
         "symbol": string,
         "asks": [{"price":float64, "volume": float64}...],
         "bids": [{"price":float64, "volume": float64}...]
     }
   ```
  #### Methods:
   * buy   
      request parametes: `{"symbol": string, "price": float64, "amount": float64}`  
     response result: `{integer}`
   * sell  
      request parametes: `{"symbol": string, "price": float64, "amount": float64}`  
      response result: `{integer}`
   * cancel_trade  
      request parameters: `{"orderid": integer}`  
      response result: `{Order}`
   * cancel_market  
      request parameters: `{"symbol": string}`  
      response result: `{array of order}`
   * cancel_all  
      request parameters: `{}` or `null`  
      response result: `{array of order}`
   * balance  
      request parameters: `{"currency": string}`  
      response result: `{string}`
   * order_info  
      request parameters `{"orderid": integer}`  
      response result: `{order}`
   * order_status  
      request parameters `{"orderid": integer}`  
      response result: `{string}`
   * completed  
      request parameters: `{}` or `null`  
      response result: `{[]integer}`
   * executed  
      request parameters: `{}` or `null`  
      response result:`{[]integer}`
   * orderbook  
      request parameters `{"symbol": string}`  
      response result: `{orderbook_record}`