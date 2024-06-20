package golang

import (
	"fmt"
	"go/token"
	"strings"

	code "github.com/trwk76/gocode"
)

func ID(id string) Identifier {
	if !token.IsIdentifier(id) {
		panic(fmt.Errorf("'%s' is not a valid identifier", id))
	}

	return Identifier(id)
}

func (i Identifier) Exported() bool {
	return token.IsExported(string(i))
}

var Discard Identifier = Identifier("_")

type (
	Identifier string
	Comment    string

	item interface {
		write(w *code.Writer)
	}
)

func (i Identifier) write(w *code.Writer) {
	w.WriteString(string(i))
}

func (c Comment) write(w *code.Writer) {
	if len(c) < 1 {
		return
	}

	for _, line := range strings.Split(string(c), "\n") {
		w.WriteString("// ")
		w.WriteString(line)
		w.Newline()
	}
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
