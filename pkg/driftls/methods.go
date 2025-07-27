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
		s.initialize(request.Id)

	case "textDocument/didOpen":
		if err := s.documents.Open(request.Params); err != nil {
			return err
		}
	case "textDocument/didChange":
		if err := s.documents.Change(request.Params); err != nil {
			return err
		}
	case "textDocument/didClose":
		if err := s.documents.Close(request.Params); err != nil {
			return err
		}
	case "textDocument/semanticTokens/full":
		if err := s.sendTokens(request.Id, request.Params); err != nil {
			return err
		}
	}

	return nil
}

func parseAndHandle[T any](s *Server, params json.RawMessage, handler func(T) error) {
	var parsedParams T
	if err := json.Unmarshal(params, &parsedParams); err != nil {
		return
	}

	if err := handler(parsedParams); err != nil {
		return
	}
}

func (s *Server) initialize(id any) error {
	return s.sendServerResponse(id, lsp.InitializeResult{
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

func (s *Server) sendTokens(id any, rawParams json.RawMessage) error {
	var params DocumentParams[lsp.TextDocumentIdentifier]

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	lexer := driftc.Lexer{ParseAllErrors: true, ParseComments: true}

	tokens, errors := lexer.Tokenize([]rune(s.documents.Get(params.TextDocument.Uri)))

	var notification struct {
		lsp.TextDocumentIdentifier
		Diagnostics []lsp.Diagnostic `json:"diagnostics"`
	}

	notification.Uri = params.TextDocument.Uri
	notification.Diagnostics = make([]lsp.Diagnostic, 0)

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

	s.sendNotification("textDocument/publishDiagnostics", notification)

	var result struct {
		Data []uint `json:"data"`
	}

	result.Data = make([]uint, 0, len(tokens)*5)

	prevLine := -1
	prevColumn := 0

	for _, tok := range tokens {
		tokenType := mapTokenType(tok.Type)
		if tokenType < 0 {
			continue
		}

		line := tok.Line - 1
		column := tok.Column - 1

		deltaLine := line
		if prevLine != -1 { // absoulute, if first token
			deltaLine -= prevLine
		}

		deltaStart := column
		if prevLine == line { // delta, if on another line
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

	return s.sendServerResponse(id, result)
}
