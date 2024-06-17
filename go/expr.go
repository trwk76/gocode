package golang

import (
	"strconv"

	code "github.com/trwk76/gocode"
)

type (
	Expr interface {
		item
		simpleExpr() bool
	}

	NilVal    struct{}
	BoolVal   bool
	IntVal    int64
	UintVal   uint64
	FloatVal  float64
	RuneVal   rune
	StringVal string

	Symbol struct {
		ID      ID
		Pkg     ID
		GenArgs []Type
	}
)

func (NilVal) simpleExpr() bool { return true }
func (NilVal) write(w *code.Writer, single bool) {
	w.WriteString("nil")
}

func (BoolVal) simpleExpr() bool { return true }
func (v BoolVal) write(w *code.Writer, single bool) {
	w.WriteString(strconv.FormatBool(bool(v)))
}

func (IntVal) simpleExpr() bool { return true }
func (v IntVal) write(w *code.Writer, single bool) {
	w.WriteString(strconv.FormatInt(int64(v), 10))
}

func (UintVal) simpleExpr() bool { return true }
func (v UintVal) write(w *code.Writer, single bool) {
	w.WriteString(strconv.FormatUint(uint64(v), 10))
}

func (FloatVal) simpleExpr() bool { return true }
func (v FloatVal) write(w *code.Writer, single bool) {
	w.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 64))
}

func (RuneVal) simpleExpr() bool { return true }
func (v RuneVal) write(w *code.Writer, single bool) {
	w.WriteString(strconv.QuoteRune(rune(v)))
}

func (StringVal) simpleExpr() bool { return true }
func (v StringVal) write(w *code.Writer, single bool) {
	w.WriteString(strconv.Quote(string(v)))
}

func (Symbol) simpleExpr() bool { return true }
func (e Symbol) write(w *code.Writer, single bool) {
	if e.Pkg != "" {
		e.Pkg.write(w, single)
		w.WriteByte('.')
	}

	e.ID.write(w, single)

	if len(e.GenArgs) > 0 {
		w.WriteByte('<')
		writeList(w, e.GenArgs, true)
		w.WriteByte('>')
	}
}
