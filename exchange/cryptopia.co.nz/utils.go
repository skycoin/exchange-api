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
)

// init caches
func init() {
	_, _ = getCurrencyID("BTC")
	_, _ = getMarketID("LTC/BTC")
}

var currencyCache map[string]CurrencyInfo
var marketCache map[string]int

func getCurrencyID(currency string) (int, error) {
	if currencyCache == nil {

		currencies, err := GetCurrencies()
		if err != nil {
			return 0, err
		}
		currencyCache = make(map[string]CurrencyInfo)

		for _, v := range currencies {
			currencyCache[v.Symbol] = v
		}
	}
	if v, ok := currencyCache[normalize(currency)]; ok {
		return v.ID, nil
	}
	return 0, errors.Errorf("Currency %s does not found", currency)

}

func normalize(symbol string) string {
	symbol = strings.ToUpper(symbol)
	symbol = strings.Replace(symbol, "_", "/", -1)
	return symbol
}

func getMarketID(market string) (int, error) {
	if marketCache == nil {

		markets, err := GetTradePairs()
		if err != nil {
			return 0, err
		}
		marketCache = make(map[string]int)
		for _, v := range markets {
			marketCache[v.Label] = v.ID
		}
	}
	if v, ok := marketCache[normalize(market)]; ok {
		return v, nil
	}
	return 0, errors.Errorf("TradePair %s does not found", market)
}

func readResponse(r io.ReadCloser) (*response, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
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

	secretBytes, _ := base64.RawStdEncoding.DecodeString(secret)
	sign := sign(secretBytes, key, nuri, nonce, params)
	token := base64.StdEncoding.EncodeToString(sign)
	return "amx " + key + ":" + token + ":" + nonce
}

// Nonce creates random string
func Nonce() string {
	var b [8]byte
	rand.Read(b[:])
	return fmt.Sprintf("%x", b[:])
}
