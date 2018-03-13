# Exchange-API

[![GoDoc](https://godoc.org/github.com/skycoin/exchange-api?status.svg)](https://godoc.org/github.com/skycoin/exchange-api)
[![Build Status](https://travis-ci.org/skycoin/exchange-api.svg?branch=master)](https://travis-ci.org/skycoin/exchange-api)

Exchange-API implements an interface to various cryptocurrency exchanges APIs in Go.

<!-- MarkdownTOC autolink="true" bracket="round" depth="5" -->

- [Integration Tests](#integration-tests)

<!-- /MarkdownTOC -->


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


