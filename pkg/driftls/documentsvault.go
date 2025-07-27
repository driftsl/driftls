package driftls

import (
	"encoding/json"

	"github.com/driftsl/driftls/pkg/lsp"
)

type DocumentsVault struct {
	Documents map[string]string
}

type DocumentParams[T any] struct {
	TextDocument T `json:"textDocument"`
}

func (v *DocumentsVault) Open(rawParams json.RawMessage) error {
	var params DocumentParams[struct {
		lsp.TextDocumentIdentifier
		Text string `json:"text"`
	}]

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	v.Documents[params.TextDocument.Uri] = params.TextDocument.Text
	return nil
}

func (v *DocumentsVault) Change(rawParams json.RawMessage) error {
	var params struct {
		DocumentParams[lsp.TextDocumentIdentifier]
		ContentChanges []struct {
			Text string `json:"text"`
		} `json:"contentChanges"`
	}

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	v.Documents[params.TextDocument.Uri] = params.ContentChanges[0].Text
	return nil
}

func (v *DocumentsVault) Close(rawParams json.RawMessage) error {
	var params DocumentParams[lsp.TextDocumentIdentifier]

	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	delete(v.Documents, params.TextDocument.Uri)
	return nil
}

func (v *DocumentsVault) Get(uri string) string {
	return v.Documents[uri]
}
