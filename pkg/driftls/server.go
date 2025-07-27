package driftls

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/driftsl/driftls/pkg/jsonrpc"
)

type Server struct {
	rpcServer *jsonrpc.Server

	alive bool

	documents DocumentsVault
}

func NewServer(r *bufio.Reader, w *bufio.Writer) *Server {
	return &Server{
		rpcServer: jsonrpc.NewServer(r, w),

		documents: DocumentsVault{
			Documents: make(map[string]string),
		},
	}
}

func (s *Server) Serve() error {
	s.alive = true

	for s.alive {
		request, err := s.rpcServer.NextRequest()
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

func (s *Server) sendServerResponse(id any, r any) error {
	return s.sendJsonRpcResponse(id, r, nil)
}

func (s *Server) sendServerError(id any, code int, message string) error {
	return s.sendJsonRpcResponse(id, nil, &jsonrpc.Error{Code: code, Message: message})
}

func (s *Server) sendNotification(method string, params any) error {
	return s.sendJson(jsonrpc.Request[any]{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	})
}

func (s *Server) sendJsonRpcResponse(id any, result any, jsonRpcError *jsonrpc.Error) error {
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
