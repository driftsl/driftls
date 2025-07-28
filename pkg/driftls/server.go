package driftls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/driftsl/driftls/pkg/jsonrpc"
)

type Server struct {
	reader *bufio.Reader
	writer *bufio.Writer

	alive bool

	documents DocumentsVault
}

func NewServer(r *bufio.Reader, w *bufio.Writer) *Server {
	return &Server{
		reader: r,
		writer: w,

		documents: DocumentsVault{
			Documents: make(map[string]string),
		},
	}
}

func (s *Server) nextRequest() (*jsonrpc.Request[json.RawMessage], error) {
	var contentLength int

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// trim '\n' or '\r\n'
		length := len(line)
		if length >= 2 && line[length-2] == '\r' {
			line = line[:length-2]
		} else {
			line = line[:length-1]
		}

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

	var result jsonrpc.Request[json.RawMessage]

	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Server) Serve() error {
	s.alive = true

	for s.alive {
		request, err := s.nextRequest()
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, request.Method, request.Id)

		if err := s.handleRequest(request); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) sendRpcResponse(id any, r any) error {
	return s.sendRpc(id, r, nil)
}

func (s *Server) sendRpcError(id any, code int, message string) error {
	return s.sendRpc(id, nil, &jsonrpc.Error{Code: code, Message: message})
}

func (s *Server) sendRpcNotification(method string, params any) error {
	return s.sendJson(jsonrpc.Request[any]{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	})
}

func (s *Server) sendRpc(id any, result any, jsonRpcError *jsonrpc.Error) error {
	return s.sendJson(jsonrpc.Response{
		JsonRpc: "2.0",
		Id:      id,
		Result:  result,
		Error:   jsonRpcError,
	})
}

func (s *Server) sendJson(object any) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return s.send(data)
}

func (s *Server) send(data []byte) error {
	if _, err := s.writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))); err != nil {
		return err
	}

	if _, err := s.writer.Write(data); err != nil {
		return err
	}

	return s.writer.Flush()
}
