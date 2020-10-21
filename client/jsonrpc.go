package client

import (
	"encoding/json"
	"errors"
)

var (
	requestId = int64(0)
)

const (
	JsonRpcVersion = "2.0"
)

type (
	RpcRequests []*RpcRequest

	RpcRequest struct {
		JsonRpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
		Id      int64       `json:"id,omitempty"`
	}

	RpcResponse struct {
		JsonRpc string      `json:"jsonrpc"`
		Error   *RpcError   `json:"error,omitempty"`
		Result  interface{} `json:"result,omitempty"`
		Id      int64       `json:"id,omitempty"`
	}

	RpcError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func (r *Request) RpcCall(result interface{}, method string, params interface{}) error {

	req := &RpcRequest{JsonRpc: JsonRpcVersion, Method: method, Params: params, Id: genId()}
	var resp *RpcResponse
	err := r.Post(&resp, "", req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.New("RPC Call error")
	}
	return resp.GetObject(result)
}

func (r *RpcResponse) GetObject(toType interface{}) error {
	js, err := json.Marshal(r.Result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, toType)
	if err != nil {
		return err
	}
	return nil
}

func genId() int64 {
	requestId += 1
	return requestId
}