package client

import (
	"encoding/json"
	"fmt"
)

var requestID = int64(0)

const JsonRpcVersion = "2.0"

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

	RpcResponseRaw struct {
		JsonRpc string          `json:"jsonrpc"`
		Error   *RpcError       `json:"error,omitempty"`
		Result  json.RawMessage `json:"result,omitempty"`
		Id      int64           `json:"id,omitempty"`
	}

	RpcError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}
)

func (r *Request) RpcCall(result interface{}, method string, params interface{}) error {
	req := &RpcRequest{JsonRpc: JsonRpcVersion, Method: method, Params: params, Id: genID()}
	var resp *RpcResponse
	err := r.Post(&resp, "", req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}
	return resp.GetObject(result)
}

func (r *Request) RpcCallRaw(method string, params interface{}) ([]byte, error) {
	req := &RpcRequest{JsonRpc: JsonRpcVersion, Method: method, Params: params, Id: genID()}
	var resp *RpcResponseRaw
	err := r.Post(&resp, "", req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return []byte(resp.Result), nil
}

func (r *Request) RpcBatchCall(requests RpcRequests) ([]RpcResponse, error) {
	var resp []RpcResponse
	err := r.Post(&resp, "", requests.fillDefaultValues())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *RpcError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
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

func (rs RpcRequests) fillDefaultValues() RpcRequests {
	for _, r := range rs {
		r.JsonRpc = JsonRpcVersion
		r.Id = genID()
	}
	return rs
}

func genID() int64 {
	requestID++
	return requestID
}
