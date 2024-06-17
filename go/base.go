package golang

import (
	"strings"

	code "github.com/trwk76/gocode"
)

type (
	Comment string
	ID      string

	item interface {
		write(w *code.Writer, single bool)
	}
)

var (
	Discard ID = ID("_")
)

func (c Comment) write(w *code.Writer, single bool) {
	if single {
		w.WriteString("/*")
		w.WriteString(strings.ReplaceAll(string(c), "\n", " "))
		w.WriteString("*/")
	} else {
		for _, line := range strings.Split(string(c), "\n") {
			w.WriteString("//")
			w.WriteString(line)
			w.Newline()
		}
	}
}

func (i ID) write(w *code.Writer, single bool) {
	w.WriteString(string(i))
}

func writeList[I item](w *code.Writer, items []I, single bool) {
	if single {
		for idx, itm := range items {
			if idx > 0 {
				w.WriteString(", ")
			}

			itm.write(w, single)
		}
	}

	w.Indent(func(w *code.Writer) {
		for _, itm := range items {
			itm.write(w, single)
			w.WriteByte(',')
			w.Newline()
		}
	})
}

func toString(i item) string {
	buf := strings.Builder{}
	w := code.NewWriter(&buf, "")

	i.write(w, true)

	w.Close()
	return buf.String()
}

var (
	_ item = Comment("")
	_ item = Discard
)
