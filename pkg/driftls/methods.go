package driftls

import (
	"encoding/json"
	"fmt"

	"github.com/driftsl/driftc/pkg/driftc"
	"github.com/driftsl/driftls/pkg/jsonrpc"
	"github.com/driftsl/driftls/pkg/lsp"
)

func (s *Server) handleRequest(request *jsonrpc.Request[json.RawMessage]) error {
	switch request.Method {
	case "initialize":
		return s.initialize(request.Id)

	case "textDocument/didOpen":
		if ok, params := tryParse[lsp.DidOpenTextDocumentParams](s, request); ok {
			return s.documents.Open(params)
		}
	case "textDocument/didChange":
		if ok, params := tryParse[lsp.DidChangeTextDocumentParams](s, request); ok {
			return s.documents.Change(params)
		}
	case "textDocument/didClose":
		if ok, params := tryParse[lsp.DidCloseTextDocumentParams](s, request); ok {
			return s.documents.Close(params)
		}

	case "textDocument/semanticTokens/full":
		if ok, params := tryParse[lsp.SemanticTokensParams](s, request); ok {
			return s.sendTokens(request.Id, params)
		}
	}

	return nil
}

func tryParse[T any](s *Server, request *jsonrpc.Request[json.RawMessage]) (bool, *T) {
	var params T
	if err := json.Unmarshal(request.Params, &params); err != nil {
		s.sendRpcError(request.Id, jsonrpc.CodeInvalidParams, err.Error())
		return false, &params
	}

	return true, &params
}

func (s *Server) initialize(id any) error {
	return s.sendRpcResponse(id, lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: ptr(lsp.TextDocumentSyncKindFull),

			SemanticTokensProvider: &lsp.SemanticTokensOptions{
				Legend: lsp.SemanticTokensLegend{
					TokenTypes:     tokensArray[:],
					TokenModifiers: make([]string, 0),
				},

				Full: ptr(true),
			},
		},

		ServerInfo: &lsp.ServerInfo{
			Name:    "driftls",
			Version: ptr(fmt.Sprintf("v%s (driftc v%s)", VERSION, driftc.VERSION)),
		},
	})
}

func (s *Server) sendTokens(id any, params *lsp.SemanticTokensParams) error {
	lexer := driftc.Lexer{ParseAllErrors: true, ParseComments: true}

	tokens, errors := lexer.Tokenize([]rune(s.documents.Get(params.TextDocument.Uri)))

	notification := lsp.PublishDiagnosticsParams{
		Uri:         params.TextDocument.Uri,
		Diagnostics: make([]lsp.Diagnostic, 0),
	}

	for _, err := range errors {
		notification.Diagnostics = append(notification.Diagnostics, lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      err.Token.Line - 1,
					Character: err.Token.Column - 1,
				},
				End: lsp.Position{
					Line:      err.Token.Line - 1,
					Character: err.Token.Column - 1 + len(err.Token.Value),
				},
			},
			Severity: 1,
			Source:   "lexer",
			Message:  err.Err.Error(),
		})
	}

	s.sendRpcNotification("textDocument/publishDiagnostics", notification)

	result := lsp.SemanticTokens{
		Data: make([]uint, 0),
	}

	prevLine := -1
	prevColumn := 0

	for _, tok := range tokens {
		tokenType := mapTokenType(tok.Type)
		if tokenType < 0 {
			continue
		}

		line := tok.Line - 1
		column := tok.Column - 1

		deltaLine := line // absoulute, if first token
		if prevLine != -1 {
			deltaLine -= prevLine
		}

		deltaStart := column // absoulute, if on another line
		if prevLine == line {
			deltaStart -= prevColumn
		}

		length := len(tok.Value)

		result.Data = append(result.Data,
			uint(deltaLine),
			uint(deltaStart),
			uint(length),
			uint(tokenType),
			0,
		)

		prevLine, prevColumn = line, column
	}

	return s.sendRpcResponse(id, result)
}
