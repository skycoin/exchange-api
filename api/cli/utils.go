package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"

	"fmt"

	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/uberfurrer/tradebot/exchange"
)

var (
	addr = "localhost:12345"
)

func makeRPCCall(endpoint, method string, params json.RawMessage) (json.RawMessage, error) {
	var id = strconv.Itoa(rand.Int())
	var r = rpc.Request{
		ID:      &id,
		JSONRPC: rpc.JSONRPC,
		Method:  method,
		Params:  params,
	}
	response, err := rpc.Do(addr, endpoint, r)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error
	}
	return response.Result, nil
}

// ErrInvalidAgrs returns if command arguments incorrect
var (
	ErrInvalidArgs = errors.New("invalid command arguments")
)

func printOrderInfo(m json.RawMessage) error {
	var result exchange.OrderInfo
	err := json.Unmarshal(m, &result)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	_, err = fmt.Println(result)
	return err
}
func printOrderID(m json.RawMessage) error {
	var result int
	err := json.Unmarshal(m, &result)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	_, err = fmt.Println(result)
	return err
}
func printOrderInfoArr(m json.RawMessage) error {
	var result []exchange.OrderInfo
	err := json.Unmarshal(m, &result)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	for _, v := range result {
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}
		_, err = fmt.Println(v)
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	return nil
}
func printString(m json.RawMessage) error {
	var result string
	err := json.Unmarshal(m, &result)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	_, err = fmt.Println(result)
	return err
}
func printOrderbook(m json.RawMessage) error {
	var result exchange.Orderbook
	err := json.Unmarshal(m, &result)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	_, err = fmt.Println(result)
	return err
}
