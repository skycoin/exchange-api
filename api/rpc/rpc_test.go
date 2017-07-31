package rpc

import "testing"
import "encoding/json"

func TestPackageHandler_AddFunction(t *testing.T) {
	var request = Request{
		ID:      new(string),
		JSONRPC: JSONRPC,
		Method:  "test",
		Params:  json.RawMessage("{\"param\": \"value\"}"),
	}
	var f = func(r Request, env map[string]string) Response {
		params, err := DecodeParams(r)
		if err != nil {
			t.Fatal(err)
		}
		param, err := GetStringParam(params, "param")
		if err != nil {
			t.Fatal(err)
		}
		if param != "value" {
			t.Fatal("want param value \"value\", expected", param)
		}
		return MakeSuccessResponse(r, env["envparam"])
	}
	var handler = PackageHandler{
		Client: new(ex),
		Handlers: map[string]PackageFunc{
			"test": f,
		},
		Env: map[string]string{
			"envparam": "envval",
		},
	}
	response := handler.Process(request)
	if response == nil {
		t.Fatal("response is nil")
	}
	if string(response.Result) != "\"envval\"" {
		t.Fatal("response result want envval, expected", string(response.Result))
	}
}
