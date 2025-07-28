package jsonrpc

type Request[T any] struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  T      `json:"params"`
	Id      any    `json:"id,omitempty"`
}

type Response struct {
	JsonRpc string `json:"jsonrpc"`
	Id      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	CodeParseError     = -32700 // Invalid JSON was received by the server.
	CodeInvalidRequest = -32600 // The JSON sent is not a valid Request object.
	CodeMethodNotFound = -32601 // The method does not exist / is not available.
	CodeInvalidParams  = -32602 // Invalid method parameter(s).
	CodeInternalError  = -32603 // Internal JSON-RPC error.

	CodeImplementationDefined = -32000 // Implementation-defined server-errors: -32000 to -32099
)
