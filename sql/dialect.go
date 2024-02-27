package sql

import code "github.com/trwk76/gocode"

type (
	Dialect interface {
		WriteName(w *code.Writer, name Name)
	}
)
