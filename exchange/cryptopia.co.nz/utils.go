package cryptopia

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/skycoin/exchange-api/exchange"
)

var incr int

func getCurrencyID(currency string) (int, error) {
	if v, ok := currencyCache[normalize(currency)]; ok {
		return v.ID, nil
	}
	// If not found, try update first
	err := updateCaches()
	if err != nil {
		return 0, err
	}
	if v, ok := currencyCache[normalize(currency)]; ok {
		return v.ID, nil
	}

	return 0, errors.Errorf("Currency %s does not found", currency)
}

func updateCaches() error {
	if currencyCache == nil {
		currencyCache = make(map[string]CurrencyInfo)
	}
	if marketCache == nil {
		marketCache = make(map[string]int)
	}
	crs, err := getCurrencies()
	if err != nil {
		return err
	}
	for _, v := range crs {
		currencyCache[v.Symbol] = v
	}
	mrkts, err := getTradePairs()
	if err != nil {
		return err
	}
	for _, v := range mrkts {
		marketCache[v.Label] = v.ID
	}
	return nil
}

func normalize(symbol string) string {
	symbol = strings.ToUpper(symbol)
	symbol = strings.Replace(symbol, "_", "/", -1)
	return symbol
}

func getMarketID(market string) (int, error) {
	if v, ok := marketCache[normalize(market)]; ok {
		return v, nil
	}
	// If not found, try update first
	err := updateCaches()
	if err != nil {
		return 0, err
	}
	if v, ok := marketCache[normalize(market)]; ok {
		return v, nil
	}
	return 0, errors.Errorf("TradePair %s does not found", market)
}

func readResponse(r io.ReadCloser) (*response, error) {
	incr++
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	//ioutil.WriteFile(fmt.Sprintf("responses/response%d.json", incr), b, os.ModePerm)
	defer r.Close()
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	var resp response
	err = json.Unmarshal(b, &resp)
	return &resp, err
}

func encodeValues(vals map[string]interface{}) []byte {
	if vals == nil {
		return []byte("{}")
	}
	b, _ := json.Marshal(vals)
	return b
}

func sign(secret []byte, key, uri, nonce string, params []byte) []byte {
	signer := hmac.New(sha256.New, secret)
	data := prepare(key, uri, nonce, params)
	signer.Write(data[:])
	return signer.Sum(nil)
}
func prepare(key, uri, nonce string, params []byte) []byte {
	hash := md5.Sum(params)
	encodedParams := base64.StdEncoding.EncodeToString(hash[:])
	var signData []byte
	signData = append(signData, key...)
	signData = append(signData, "POST"...)
	signData = append(signData, uri...)
	signData = append(signData, nonce...)
	signData = append(signData, encodedParams...)
	return signData[:]
}
func header(key, secret, nonce string, uri url.URL, params []byte) string {
	nuri := strings.ToLower(url.QueryEscape(uri.String()))
	secretBytes, _ := base64.StdEncoding.DecodeString(secret)
	sign := sign(secretBytes, key, nuri, nonce, params)
	token := base64.StdEncoding.EncodeToString(sign)
	return "amx " + key + ":" + token + ":" + nonce
}

// nonce creates random string
func nonce() string {
	var b [8]byte
	rand.Read(b[:])
	return fmt.Sprintf("%x", b[:])
}

func convert(order Order) exchange.Order {
	return exchange.Order{
		OrderID:         order.OrderID,
		Market:          order.Market,
		Price:           order.Rate,
		Amount:          order.Amount,
		Accepted:        order.Timestamp,
		Fee:             order.Fee,
		CompletedAmount: order.Total.Sub(order.Remaining),
	}
}
