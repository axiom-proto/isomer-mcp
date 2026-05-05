package schema

import "encoding/json"

// Request is a JSON-RPC request received from an MCP client.
type Request struct {
	JsonRPC string          `json:"jsonrpc"`
	Id      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a JSON-RPC response sent to an MCP client.
type Response struct {
	JsonRPC string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError describes a JSON-RPC error payload.
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HandshakeResponse is returned from the MCP initialize method.
type HandshakeResponse struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      map[string]string      `json:"serverInfo"`
}

// ResourceListResponseItem describes a resource available through MCP.
type ResourceListResponseItem struct {
	Uri         string `json:"uri"`
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// ResourceListResponse is returned from the MCP resources/list method.
type ResourceListResponse struct {
	Resources []ResourceListResponseItem `json:"resources"`
}
