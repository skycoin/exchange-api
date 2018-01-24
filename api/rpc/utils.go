package rpc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// validateRequest checks that request has correct jsonrpc version and has id
func validateRequest(r Request) (resp Response, err error) {
	resp.JSONRPC = JSONRPC
	if r.JSONRPC != JSONRPC {
		resp.Error = makeError(InvalidRequest, invalidRequestMsg, nil)
		return resp, resp.Error
	}
	if r.ID == nil {
		resp.Error = makeError(InvalidRequest, invalidRequestMsg, nil)
		return resp, resp.Error
	}
	resp.ID = (*r.ID)
	return
}

// DecodeParams extract params from request to map[string]interface{}
// You need explictly vaildate all params, that needed for you
func DecodeParams(r Request) (map[string]interface{}, error) {
	if r.Params == nil {
		return nil, errEmptyParams
	}
	result := make(map[string]interface{})
	err := json.Unmarshal(r.Params, &result)
	return result, err
}

var errEmptyParams = errors.New("params empty")

// GetIntParam extract int param from request params
func GetIntParam(params map[string]interface{}, key string) (value int, err error) {
	if v, ok := params[key]; ok {
		if val, ok := v.(float64); ok {
			return int(val), nil
		}
		return 0, errInvalidType
	}
	return 0, errParamNotFound
}

// GetFloatParam extract float64 params from request params
func GetFloatParam(params map[string]interface{}, key string) (value float64, err error) {
	if v, ok := params[key]; ok {
		if value, ok = v.(float64); ok {
			return value, nil
		}
		return 0, errInvalidType
	}
	return 0, errParamNotFound
}

// GetStringParam extract string param from request params
func GetStringParam(params map[string]interface{}, key string) (value string, err error) {
	if v, ok := params[key]; ok {
		if value, ok = v.(string); ok {
			return value, nil
		}
		return "", errInvalidType
	}
	return "", errParamNotFound
}

var errParamNotFound = errors.New("param does not found")
var errInvalidType = errors.New("invalid param type")

// MakeErrorResponse creates error response
func MakeErrorResponse(r Request, errortype int, err error) Response {
	errorcode := 0
	errormsg  := ""
	switch errortype {
	case InvalidParams:
		errorcode = InvalidParams
		errormsg = invalidParamsMsg
	case InvalidRequest:
		errorcode = InvalidRequest
		errormsg = invalidRequestMsg
	case MethodNotFound:
		errorcode = MethodNotFound
		errormsg = methodNotFoundMsg
	case ParseError:
		errorcode = ParseError
		errormsg = parseErrorMsg
	default:
		errorcode = InternalError
		errormsg = internalErrorMsg
	}
	return Response{
		ID:      *r.ID,
		JSONRPC: JSONRPC,
		Error:   makeError(errorcode, errormsg, err),
	}
}

// MakeSuccessResponse creates success response
func MakeSuccessResponse(r Request, result interface{}) Response {
	data, err := json.Marshal(result)
	if err != nil {
		return Response{
			ID:      *r.ID,
			JSONRPC: r.JSONRPC,
			Error:   makeError(InternalError, internalErrorMsg, err),
		}
	}
	return Response{
		ID:      *r.ID,
		JSONRPC: JSONRPC,
		Result:  data,
	}
}
