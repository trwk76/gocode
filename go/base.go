package golang

import (
	"strings"

	code "github.com/trwk76/gocode"
)

type (
	item interface {
		spaceAfter() bool
		write(w *code.Writer)
	}

	row interface {
		row() code.Row
	}

	items[T item] []T
	rows[T row]   []T
)

func (r items[T]) write(w *code.Writer) {
	prvSpaceAfter := false

	for _, itm := range r {
		if prvSpaceAfter {
			w.Newline()
		}

		itm.write(w)
		prvSpaceAfter = itm.spaceAfter()
	}
}

func (r rows[T]) write(w *code.Writer) {
	rows := make([]code.Row, len(r))

	for idx, itm := range r {
		rows[idx] = itm.row()
	}

	w.Table(rows...)
}

func commentRenderer(text string) code.Renderer {
	if text == "" {
		return nil
	}

	return func(w *code.Writer) {
		for _, line := range strings.Split(text, "\n") {
			w.WriteString("// ")
			w.WriteString(line)
			w.Newline()
		}
	}
}
