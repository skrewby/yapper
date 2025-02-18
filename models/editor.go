package models

import (
	"bytes"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Editor struct {
	db *pgxpool.Pool
}

func NewEditorModel(db *pgxpool.Pool) *Editor {
	return &Editor{
		db,
	}
}

func (m Editor) ConvertMarkdownToHTML(source []byte) (string, error) {
	var buf bytes.Buffer
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	if err := gm.Convert(source, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
