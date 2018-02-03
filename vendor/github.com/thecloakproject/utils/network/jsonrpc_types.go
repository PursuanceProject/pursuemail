// Steve Phillips / elimisteve
// 2013.01.13
// Originally part of Decentra prototype, then github.com/elimisteve/fun

package network

import (
	"encoding/json"
)

// From http://golang.org/src/pkg/net/rpc/jsonrpc/client.go

type ClientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	Id     uint64         `json:"id"`
}

type ClientResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     uint64           `json:"id"`
}
