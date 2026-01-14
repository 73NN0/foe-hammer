package stdio

import "encoding/json"

type Request struct {
	ID     string          `json:"id,omitempty"`
	Op     string          `json:"op"`
	Params json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	ID     string `json:"id,omitempty"`
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result any    `json:"result,omitempty"`
}

type GetParams struct {
	ID int `json:"id"`
}
