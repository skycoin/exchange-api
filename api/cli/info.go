package cli

import (
	"strconv"
	"strings"

	"encoding/json"

	"fmt"

	"github.com/uberfurrer/tradebot/api/rpc"
	"github.com/uberfurrer/tradebot/exchange/tracker"
	"github.com/urfave/cli"
)

func orderCmd() cli.Command {
	var name = "order"
	return cli.Command{
		Name:      name,
		Usage:     "Get info about order",
		ArgsUsage: "[subcommand orderid]",
		Subcommands: []cli.Command{
			infoCmd(),
			statusCmd(),
		},
	}
}

// endpoint order_info
// params: {"orderid":1}
func infoCmd() cli.Command {
	var name = "info"
	return cli.Command{
		Name:  name,
		Usage: "Gets detailed information about order",
		Action: func(c *cli.Context) error {
			var (
				args    = c.Args()
				orderid int
				err     error
			)
			if len(args) != 1 {
				return errInvalidParams
			}
			if orderid, err = strconv.Atoi(args[0]); err != nil {
				return errInvalidParams
			}
			params, _ := json.Marshal(map[string]int{"orderid": orderid})
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "...",
				Params:  params,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request %s", err)
				return errRPC
			}
			var result tracker.Order
			if err = json.Unmarshal(resp.Result, &result); err != nil {
				fmt.Printf("Error: invalid response, %s", err)
				return errInvalidResponse
			}
			b, _ := json.MarshalIndent(result, "", "    ")
			fmt.Println(string(b))
			return nil
		},
	}
}

// endpoint order_status
// params: {"orderid":1}
func statusCmd() cli.Command {
	var name = "status"
	return cli.Command{
		Name:  name,
		Usage: "Gets status of order",
		Action: func(c *cli.Context) error {
			var (
				args    = c.Args()
				orderid int
				err     error
			)
			if len(args) != 1 {
				return errInvalidParams
			}
			if orderid, err = strconv.Atoi(args[0]); err != nil {
				return errInvalidParams
			}
			params, _ := json.Marshal(map[string]int{"orderid": orderid})
			var req = rpc.Request{
				ID:      reqID(),
				JSONRPC: rpc.JSONRPC,
				Method:  "...",
				Params:  params,
			}
			resp, err := rpc.Do(rpcaddr, endpoint, req)
			if err != nil {
				fmt.Printf("Error processing request %s", err)
				return errRPC
			}
			var result string
			if err = json.Unmarshal(resp.Result, &result); err != nil {
				fmt.Printf("Error: invalid response, %s", err)
				return errInvalidResponse
			}
			b, _ := json.MarshalIndent(map[string]interface{}{
				"orderid": orderid,
				"status":  result,
			}, "", "    ")
			fmt.Println(string(b))
			return nil
		},
	}
}
func completedCmd() cli.Command {
	return cli.Command{
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "all",
				Usage: "Gets details about all completed orders",
				Action: func(c *cli.Context) error {
					var args = c.Args()
					if len(args) != 0 {
						return errInvalidParams
					}
					orders, err := getCompleted()
					if err != nil {
						return err
					}
					for _, v := range orders {
						b, _ := json.MarshalIndent(v, "", "    ")
						fmt.Println(string(b))
					}
					return nil
				},
			},
			cli.Command{
				Name:      "market",
				Usage:     "Gets details about all completed orders in markets",
				ArgsUsage: "[tradepairs...]",
				Action: func(c *cli.Context) error {
					var args = c.Args()
					if len(args) == 0 {
						return errInvalidParams
					}
					orders, err := getCompleted()
					if err != nil {
						return err
					}
					for _, v := range args {
						v = strings.ToUpper(v)
						for _, order := range orders {
							if order.Market == v {
								b, _ := json.MarshalIndent(order, "", "    ")
								fmt.Println(string(b))
							}
						}
					}
					return nil
				},
			},
			cli.Command{
				Name:      "order",
				Usage:     "Gets details about completed orders",
				ArgsUsage: "[orderids...]",
				Action: func(c *cli.Context) error {
					var args = c.Args()
					if len(args) == 0 {
						return errInvalidParams
					}
					orders, err := getCompleted()
					if err != nil {
						return err
					}
					for i, v := range args {
						var (
							orderid int
							err     error
						)
						if orderid, err = strconv.Atoi(v); err != nil {
							fmt.Printf("Invalid orderid %v, position %d", v, i)
							return errInvalidParams
						}
						for i := 0; i < len(orders); i++ {
							if orders[i].OrderID == orderid {
								b, _ := json.MarshalIndent(orders[i], "", "    ")
								fmt.Println(string(b))
							}
						}
					}
					return nil
				},
			},
		},
	}
}
func executedCmd() cli.Command {
	return cli.Command{
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "all",
				Usage: "Gets details about all executed orders",
				Action: func(c *cli.Context) error {
					var args = c.Args()
					if len(args) != 0 {
						return errInvalidParams
					}
					orders, err := getExecuted()
					if err != nil {
						return err
					}
					for _, v := range orders {
						b, _ := json.MarshalIndent(v, "", "    ")
						fmt.Println(string(b))
					}
					return nil
				},
			},
			cli.Command{
				Name:      "market",
				Usage:     "Gets details about all executed orders in markets",
				ArgsUsage: "[tradepairs...]",
				Action: func(c *cli.Context) error {
					var args = c.Args()
					if len(args) == 0 {
						return errInvalidParams
					}
					orders, err := getExecuted()
					if err != nil {
						return err
					}
					for _, v := range args {
						v = strings.ToUpper(v)
						for _, order := range orders {
							if order.Market == v {
								b, _ := json.MarshalIndent(order, "", "    ")
								fmt.Println(string(b))
							}
						}
					}
					return nil
				},
			},
			cli.Command{
				Name:      "order",
				Usage:     "Gets details about executed orders",
				ArgsUsage: "[orderids...]",
				Action: func(c *cli.Context) error {
					var args = c.Args()
					if len(args) == 0 {
						return errInvalidParams
					}
					orders, err := getExecuted()
					if err != nil {
						return err
					}
					for i, v := range args {
						var (
							orderid int
							err     error
						)
						if orderid, err = strconv.Atoi(v); err != nil {
							fmt.Printf("Invalid orderid %v, position %d", v, i)
							return errInvalidParams
						}
						for i := 0; i < len(orders); i++ {
							if orders[i].OrderID == orderid {
								b, _ := json.MarshalIndent(orders[i], "", "    ")
								fmt.Println(string(b))
							}
						}
					}
					return nil
				},
			},
		},
	}
}

func getCompleted() ([]tracker.Order, error) {
	var req = rpc.Request{
		ID:      reqID(),
		JSONRPC: rpc.JSONRPC,
		Method:  "completed",
		Params:  nil,
	}
	resp, err := rpc.Do(rpcaddr, endpoint, req)
	if err != nil {
		fmt.Printf("Error processing request %s", err)
		return nil, errRPC
	}
	var result []tracker.Order
	if err = json.Unmarshal(resp.Result, &result); err != nil {
		fmt.Printf("Error: invalid response, %s", err)
		return nil, errInvalidResponse
	}
	return result, nil
}

func getExecuted() ([]tracker.Order, error) {
	var req = rpc.Request{
		ID:      reqID(),
		JSONRPC: rpc.JSONRPC,
		Method:  "executed",
		Params:  nil,
	}
	resp, err := rpc.Do(rpcaddr, endpoint, req)
	if err != nil {
		fmt.Printf("Error processing request %s", err)
		return nil, errRPC
	}
	var result []tracker.Order
	if err = json.Unmarshal(resp.Result, &result); err != nil {
		fmt.Printf("Error: invalid response, %s", err)
		return nil, errInvalidResponse
	}
	return result, nil
}
