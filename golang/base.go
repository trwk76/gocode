package golang

import (
	"strings"

	"gitlab.trwk.com/go/code"
)

type (
	Item interface {
		write(w *code.Writer)
	}

	Comment string

	ID string
)

func Render(item Item) string {
	if item == nil {
		return ""
	}

	return code.WriteString("", func(w *code.Writer) { item.write(w) })
}

func (i Comment) write(w *code.Writer) {
	if i == "" {
		return
	}

	for _, line := range strings.Split(string(i), "\n") {
		w.WriteString("// ")
		w.WriteString(line)
		w.Newline()
	}
}

func (i ID) write(w *code.Writer) {
	w.WriteString(string(i))
}

var (
	_ Item = Comment("")
	_ Item = ID("")
)
