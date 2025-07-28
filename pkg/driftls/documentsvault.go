package driftls

import (
	"github.com/driftsl/driftls/pkg/lsp"
)

type DocumentsVault struct {
	Documents map[string]string
}

func (v *DocumentsVault) Open(params *lsp.DidOpenTextDocumentParams) error {
	v.Documents[params.TextDocument.Uri] = params.TextDocument.Text
	return nil
}

func (v *DocumentsVault) Change(params *lsp.DidChangeTextDocumentParams) error {
	v.Documents[params.TextDocument.Uri] = params.ContentChanges[0].Text
	return nil
}

func (v *DocumentsVault) Close(params *lsp.DidCloseTextDocumentParams) error {
	delete(v.Documents, params.TextDocument.Uri)
	return nil
}

func (v *DocumentsVault) Get(uri string) string {
	return v.Documents[uri]
}
