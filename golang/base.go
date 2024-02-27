package golang

import (
	"strings"

	code "github.com/trwk76/gocode"
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

func (i ID) Exported() bool {
	return i.Valid() && isUpper(i[0])
}

func (i ID) Internal() bool {
	return i.Valid() && !isUpper(i[0])
}

func (i ID) Valid() bool {
	if len(i) < 1 {
		return false
	}

	if !isLetter(i[0]) && i[0] != '_' {
		return false
	}

	for idx := 1; idx < len(i); idx++ {
		if !isLetter(i[idx]) && !isDigit(i[idx]) && i[idx] != '_' {
			return false
		}
	}

	return true
}

func (i ID) write(w *code.Writer) {
	w.WriteString(string(i))
}

func isLetter(c byte) bool {
	return isUpper(c) || isLower(c)
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

var (
	_ Item = Comment("")
	_ Item = ID("")
)
