package c2cx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"

	// the following is nolinted because it's part of c2cx' authentication scheme
	"crypto/md5" // nolint: gas

	"github.com/skycoin/exchange-api/exchange"
)

func sign(secret string, params url.Values) string {
	var paramString = encodeParamsSorted(params)
	if len(paramString) > 0 {
		paramString += "&secretKey=" + secret
	} else {
		paramString += "secretKey=" + secret
	}

	sum := md5.Sum([]byte(paramString)) // nolint: gas
	return strings.ToUpper(fmt.Sprintf("%x", sum))
}

// returns sorted string for signing
func encodeParamsSorted(params url.Values) string {
	if params == nil {
		return ""
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	result := bytes.Buffer{}
	for i, k := range keys {
		result.WriteString(url.QueryEscape(k))
		result.WriteString("=")
		result.WriteString(url.QueryEscape(params.Get(k)))

		if i != len(keys)-1 {
			result.WriteString("&")
		}
	}

	return result.String()
}

// normalize tradepair symbol
func normalize(sym string) (string, error) {
	sym = strings.ToUpper(strings.Replace(sym, "/", "_", -1))
	for _, v := range Markets {
		if v == sym {
			return sym, nil
		}
	}
	return "", fmt.Errorf("Market pair %s does not exists", sym)
}

func readResponse(r io.ReadCloser) (*response, error) {
	var tmp struct {
		Fail    []json.RawMessage `json:"fail,omitempty"`
		Success json.RawMessage   `json:"success,omitempty"`
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := json.Unmarshal(b, &tmp); err != nil {
		return nil, err
	}

	if len(tmp.Fail) != 0 {
		return nil, errors.New(string(tmp.Fail[0]))
	}

	var resp response
	if err := json.Unmarshal(tmp.Success, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func convert(order Order) exchange.Order {
	status := lookupStatus(order.Status)
	var completed time.Time
	accepted := time.Unix(order.CreateDate, 0)
	if order.CompleteDate != 0 {
		completed = time.Unix(order.CompleteDate, 0)
	}

	return exchange.Order{
		OrderID: order.OrderID,
		Status:  status,
		Type:    order.Type,

		CompletedAmount: order.CompletedAmount,
		Fee:             order.Fee,

		Price:  order.Price,
		Amount: order.Amount,

		Accepted:  accepted,
		Completed: completed,
	}
}

func lookupStatus(statusID int) string {
	for k, v := range statuses {
		if v == statusID {
			return k
		}
	}
	return "unknown"
}
