package c2cx

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func ExampleMarketBuy() { // nolint: vet
	c := NewAPIClient("your-key-here", "your-secret-here")

	amount, err := decimal.NewFromString("2.12345")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	orderID, err := c.MarketBuy(BtcSky, amount, nil)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("Placed market buy order:", orderID)
}

func TestSignParams(t *testing.T) {
	key := "C821DB84-6FBD-11E4-A9E3-C86000D26D7C"
	secret := "12D857DE-7A92-F555-10AC-7566A0D84D1B"

	// No params
	expected := "BC0DE7EBA50C730BDFC575FE2CD54082"
	signed := signParams(key, secret, nil)
	require.Equal(t, expected, signed)

	// With params
	params := url.Values{}
	params.Set("foo", "123")
	params.Set("bar", "345")
	expected = "8744DB365DE080C0A45C9940879E9644"
	signed = signParams(key, secret, params)
	require.Equal(t, expected, signed)

	// With params whose value has a character that would change with query encoding
	params = url.Values{}
	params.Set("foo", "a:b")
	expected = "965DAC47B580945BFB97338AA7C34218"
	signed = signParams(key, secret, params)
	require.Equal(t, expected, signed)
}

func TestErrorInterfaces(t *testing.T) {
	otherErr := NewOtherError(errors.New("foo"))
	require.Equal(t, "foo", otherErr.Error())
	require.Implements(t, (*Error)(nil), otherErr)

	apiErr := NewAPIError(getOrderbookEndpoint, http.StatusBadRequest, "failed")
	require.Implements(t, (*Error)(nil), apiErr)
}
