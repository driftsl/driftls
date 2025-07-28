package lsp

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/

// general

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// initialization

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

const (
	TextDocumentSyncKindNone        int = 0
	TextDocumentSyncKindFull        int = 1
	TextDocumentSyncKindIncremental int = 2
)

type ServerCapabilities struct {
	TextDocumentSync *int `json:"textDocumentSync,omitempty"`

	SemanticTokensProvider *SemanticTokensOptions `json:"semanticTokensProvider,omitempty"`
}

type SemanticTokensOptions struct {
	Legend SemanticTokensLegend `json:"legend"`
	Full   *bool                `json:"full"`
}

type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

// text documents

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}

type TextDocumentItem struct {
	Uri  string `json:"uri"`
	Text string `json:"text"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}
type DidChangeTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// tokens

type SemanticTokensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SemanticTokens struct {
	Data []uint `json:"data"`
}

// diagnostics

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Code     *int   `json:"code,omitempty"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}

type PublishDiagnosticsParams struct {
	Uri string `json:"uri"`

	Diagnostics []Diagnostic `json:"diagnostics"`
}
