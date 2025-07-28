package driftls

import (
	"encoding/json"

	"github.com/driftsl/driftls/pkg/lsp"
)

type DocumentsVault struct {
	Documents map[string]string
}

func (v *DocumentsVault) Open(rawParams json.RawMessage) error {
	var params lsp.DidOpenTextDocumentParams

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	v.Documents[params.TextDocument.Uri] = params.TextDocument.Text
	return nil
}

func (v *DocumentsVault) Change(rawParams json.RawMessage) error {
	var params lsp.DidChangeTextDocumentParams

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	v.Documents[params.TextDocument.Uri] = params.ContentChanges[0].Text
	return nil
}

func (v *DocumentsVault) Close(rawParams json.RawMessage) error {
	var params lsp.DidCloseTextDocumentParams

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	delete(v.Documents, params.TextDocument.Uri)
	return nil
}

func (v *DocumentsVault) Get(uri string) string {
	return v.Documents[uri]
}
