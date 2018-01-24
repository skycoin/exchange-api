package rpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Do does request to given addr and endpoint
func Do(addr, endpoint string, r Request) (*Response, error) {
	c          := http.Client{}
	requestURI := url.URL{}

	requestURI.Host = addr
	requestURI.Scheme = "http"
	requestURI.Path = "/" + endpoint
	requestData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", requestURI.String(), bytes.NewReader(requestData))
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("reesponse status code %d", resp.StatusCode)
	}
	respdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rpcResp Response
	err = json.Unmarshal(respdata, &rpcResp)
	if err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, rpcResp.Error
	}
	return &rpcResp, nil

}
