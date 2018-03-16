package c2cx

/*

C2CX documentatation: http://api.c2cx.com

However, the documentation is unclear or wrong in some cases.

Clarifications:

The response status code is always 200 (except in the case of a 500?).
Instead, a non-200 code is set in a field in the JSON response if an error occurred.

Some timestamps are in unix *milliseconds*, others are in regular unix seconds

During pagination, pageindex of 0 and pageindex of 1 are treated the same.

A client-defined ID can be assigned when creating an order. This ID is an arbitrary string.
The parameter name is "cid".  A cid cannot be reused.

For buy market orders you add in the quantity field the amount of BTC you want to spend on SKY.
For sell market orders you add in the quantity field the amount of SKY you want to spend on BTC.

Endpoints with responses that differ from the documentation, along with their correct
responses:

getOrderByStatus
{
  "code": 200,
  "message": "success",
  "data": {
    "rows": [
      {
        "amount": 2,
        "avgPrice": 0,
        "completedAmount": "0",
        "createDate": 1520934562420,
        "updateDate": 1520934562420,
        "orderId": 3266582,
        "price": 0.00102,
        "status": 7,
        "fee": 0,
        "type": "buy",
        "trigger": 0,
        "cid": null,
        "source": "api"
      }
    ],
    "pageindex": null,
    "pagesize": null,
    "recordcount": 2,
    "pagecount": 1
  }
}


getOrderInfo [when orderID == -1, returning all]

{
  "code": 200,
  "message": "succcess",
  "data": [
    {
      "amount": 2,
      "avgPrice": 0,
      "completedAmount": "0",
      "createDate": 1520934423387,
      "orderId": 3266573,
      "price": 0.00102,
      "status": 7,
      "type": "buy",
      "fee": 0,
      "cid": null,
      "source": "api"
    }
  ]
}


getOrderInfo [for a single order]

{
  "code": 200,
  "message": "succcess",
  "data": {
    "amount": 2,
    "avgPrice": 0,
    "completedAmount": "0",
    "createDate": 1520938109560,
    "orderId": 3267112,
    "price": 0.00102,
    "status": 7,
    "type": "buy",
    "fee": 0,
    "cid": null,
    "source": "api"
  }
}

*/
