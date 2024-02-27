package sql

import code "github.com/trwk76/gocode"

type (
	Item interface {
		write(w Writer)
	}
)

func Render(d Dialect, i Item) string {
	return code.WriteString("\t", func(w *code.Writer) {
		i.write(NewWriter(w, d))
	})
}
