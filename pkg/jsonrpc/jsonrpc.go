package jsonrpc

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

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

type Server struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewServer(r *bufio.Reader, w *bufio.Writer) *Server {
	return &Server{reader: r, writer: w}
}

func (s *Server) NextRequest() (*Request[json.RawMessage], error) {
	var contentLength int

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = trimLineBreak(line)

		if line == "" {
			break
		}

		headerParts := strings.SplitN(line, ": ", 2)

		switch strings.ToLower(headerParts[0]) {
		case "content-length":
			contentLength, err = strconv.Atoi(headerParts[1])
			if err != nil {
				return nil, err
			}
		}
	}

	buf := make([]byte, contentLength)
	_, err := io.ReadFull(s.reader, buf)
	if err != nil {
		return nil, err
	}

	var result Request[json.RawMessage]

	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func trimLineBreak(line string) string {
	length := len(line)
	if length >= 2 && line[length-2] == '\r' {
		return line[:length-2]
	}
	return line[:length-1]
}
