package c2cx

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/uberfurrer/tradebot/exchange"
)

func sign(secret string, params url.Values) string {
	var paramString = abcdsort(params)
	if len(paramString) > 0 {
		paramString += "&secretKey=" + secret
	} else {
		paramString += "secretKey=" + secret
	}

	sum := md5.Sum([]byte(paramString))
	return strings.ToUpper(fmt.Sprintf("%x", sum))
}

// returns sorted string for signing
func abcdsort(params url.Values) string {
	if params == nil {
		return ""
	}
	var sortedKeys = make([]string, 0, len(params))
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}

	for i := 0; i < len(sortedKeys); i++ {
		for j := 0; j < i; j++ {
			var wordLen int
			if len(sortedKeys[i]) < len(sortedKeys[j]) {
				wordLen = len(sortedKeys[i])
			} else {
				wordLen = len(sortedKeys[j])
			}
			for let := 0; let < wordLen; let++ {
				switch {
				case sortedKeys[i][let] < sortedKeys[j][let]:
					sortedKeys[i], sortedKeys[j] = sortedKeys[j], sortedKeys[i]
					goto next
				case sortedKeys[i][let] == sortedKeys[j][let]:
					continue
				case sortedKeys[i][let] > sortedKeys[j][let]:
					goto next
				}
			}
			if len(sortedKeys[i]) < len(sortedKeys[j]) {
				sortedKeys[i], sortedKeys[j] = sortedKeys[j], sortedKeys[i]
			}
		next:
		}
	}

	var result = bytes.NewBuffer(nil)
	for i := 0; i < len(sortedKeys); i++ {
		result.WriteString(sortedKeys[i])
		result.WriteString("=")
		result.WriteString(params.Get(sortedKeys[i]))
		result.WriteString("&")
	}
	// delete last &
	return result.String()[:result.Len()-1]

}

// normalize tradepair symbol
func normalize(sym string) (string, error) {
	sym = strings.ToUpper(strings.Replace(sym, "/", "_", -1))
	for _, v := range Markets {
		if v == sym {
			return sym, nil
		}
	}
	return "", errors.Errorf("Market pair %s does not exists", sym)
}

func unix(unix int64) time.Time {
	var secs = int64(unix / 10e2)
	var nanos = int64((unix % 10e2) * 10e5)
	return time.Unix(secs, nanos)
}

func readResponse(r io.ReadCloser) (*response, error) {
	var (
		tmp struct {
			Fail    json.RawMessage `json:"fail,omitempty"`
			Success json.RawMessage `json:"success,omitempty"`
		}
		resp response
	)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	if tmp.Fail != nil {
		return nil, errors.New(string(tmp.Fail))
	}
	if tmp.Success != nil {
		err := json.Unmarshal(tmp.Success, &resp)
		return &resp, err
	}
	return &resp, json.Unmarshal(b, &resp)
}

func convert(order Order) exchange.Order {
	var (
		status    = lookupStatus(order.Status)
		accepted  time.Time
		completed time.Time
	)
	accepted = unix(order.CreateDate)
	if order.CompleteDate != 0 {
		completed = unix(order.CompleteDate)
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
