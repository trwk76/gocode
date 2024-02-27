package sql

import code "github.com/trwk76/gocode"

func NewWriter(w *code.Writer, d Dialect) Writer {
	return Writer{
		Writer: w,
		d:      d,
	}
}

type (
	Writer struct {
		*code.Writer
		d Dialect
	}
)
