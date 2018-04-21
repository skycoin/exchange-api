package cryptopia

import (
	"bytes"
	"crypto/hmac"

	// the following is nolinted because it's part of cryptopia's authentication scheme
	"crypto/md5" // nolint: gas
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

func normalize(symbol string) string {
	symbol = strings.ToUpper(symbol)
	symbol = strings.Replace(symbol, "_", "/", -1)
	return symbol
}

func readResponse(r io.ReadCloser) (*response, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = r.Close(); err != nil {
			panic(err)
		}
	}()

	bb := bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	var resp response
	if err := json.Unmarshal(bb, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func encodeValues(vals map[string]interface{}) ([]byte, error) {
	if vals == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(vals)
}

func sign(secret []byte, key, uri, nonce string, params []byte) []byte {
	signer := hmac.New(sha256.New, secret)
	data := prepare(key, uri, nonce, params)
	if _, err := signer.Write(data[:]); err != nil {
		panic(err)
	}
	return signer.Sum(nil)
}

func prepare(key, uri, nonce string, params []byte) []byte {
	hash := md5.Sum(params) // nolint: gas
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
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b[:])
}
