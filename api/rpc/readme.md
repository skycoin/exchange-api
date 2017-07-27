# RPC interface
 ## Overview
  package rpc use JSONRPC 2.0 scheme
  PackageHandler handles exchange.Client interface and allow to add additional, exchange specific handlers to each functions, that does not supported for this interface. You can add any variables that you need using Env map and PackageHandler.Setenv() function.  
  Adding new function:  
    PackageHandler.Handlers["method_name"] = func(r Request, env map[string]string) (resp Response) {  
        var result interface{}  
        //some function actions  
        resp.SetBody(result)  
        return resp  
    }  
 ## exchange.Client interface:
  #### Datatypes:
   * order   
   `{`  
     `"orderid": integer,`  
     ` "type": string,`  
     ` "market": string,`  
     ` "amount": float64,`  
     ` "price": float64,`  
     ` "submitted_at": integer,`  
     ` "fee": float64,`  
     ` "completed_amount": float64,`  
     ` "status": string,`  
     ` "completed_at": integer,`  
     ` "accepted_at": integer`  
   `}`
   * orderbook_record  
   `{`  
       `"timestamp": integer,`  
       `"symbol": string,`  
       `"asks": [{"price":float64, "volume": float64}...],`  
       `"bids": [{"price":float64, "volume": float64}...]`  
    `}`
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