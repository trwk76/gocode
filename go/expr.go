package golang

import (
	"strconv"

	code "github.com/trwk76/gocode"
)

func LitBool(value bool) LitBoolExpr {
	return LitBoolExpr(value)
}

func LitRune(value rune) LitRuneExpr {
	return LitRuneExpr(value)
}

func LitInt(value int64) LitIntExpr {
	return LitIntExpr(value)
}

func LitUint(value uint64) LitUintExpr {
	return LitUintExpr(value)
}

func LitFloat(value float64) LitFloatExpr {
	return LitFloatExpr(value)
}

func LitString(value string) LitStringExpr {
	return LitStringExpr(value)
}

func Range(lower Expr, upper Expr) RangeExpr {
	return RangeExpr{low: lower, upp: upper}
}

func Symbol(id Identifier, args ...Type) SymbolExpr {
	return SymbolExpr{id: id, args: args}
}

func Paren(expr Expr) ParenExpr {
	return ParenExpr{expr: expr}
}

func Member(expr Expr, id Identifier) MemberExpr {
	return MemberExpr{expr: expr, id: id}
}

func Call(callee Expr, args ...Expr) CallExpr {
	return CallExpr{callee: callee, args: args}
}

func Index(expr Expr, index Expr) IndexExpr {
	return IndexExpr{expr: expr, idx: index}
}

func Identity(expr Expr) UnaryExpr {
	return UnaryExpr{op: identityOp, expr: expr}
}

func Negate(expr Expr) UnaryExpr {
	return UnaryExpr{op: negateOp, expr: expr}
}

func Not(expr Expr) UnaryExpr {
	return UnaryExpr{op: notOp, expr: expr}
}

func Compl(expr Expr) UnaryExpr {
	return UnaryExpr{op: complOp, expr: expr}
}

func AddrOf(expr Expr) UnaryExpr {
	return UnaryExpr{op: addrOp, expr: expr}
}

func Indirect(expr Expr) UnaryExpr {
	return UnaryExpr{op: indirOp, expr: expr}
}

var (
	Nil   NilExpr
	False LitBoolExpr = LitBoolExpr(false)
	True  LitBoolExpr = LitBoolExpr(true)
)

type (
	Expr interface {
		item
	}

	NilExpr       struct{}
	LitBoolExpr   bool
	LitRuneExpr   rune
	LitIntExpr    int64
	LitUintExpr   uint64
	LitFloatExpr  float64
	LitStringExpr string

	RangeExpr struct {
		low Expr
		upp Expr
	}

	SymbolExpr struct {
		pkg  Identifier
		id   Identifier
		args []Type
	}

	ParenExpr struct {
		expr Expr
	}

	MemberExpr struct {
		expr Expr
		id   Identifier
	}

	CallExpr struct {
		callee Expr
		args   []Expr
	}

	IndexExpr struct {
		expr Expr
		idx  Expr
	}

	UnaryExpr struct {
		op   unaryOp
		expr Expr
	}

	unaryOp string
)

const (
	identityOp unaryOp = "+"
	negateOp   unaryOp = "-"
	notOp      unaryOp = "!"
	complOp    unaryOp = "^"
	addrOp     unaryOp = "&"
	indirOp    unaryOp = "*"
)

func (NilExpr) write(w *code.Writer) {
	w.WriteString("nil")
}

func (e LitBoolExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatBool(bool(e)))
}

func (e LitRuneExpr) write(w *code.Writer) {
	w.WriteString(strconv.QuoteRune(rune(e)))
}

func (e LitIntExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatInt(int64(e), 10))
}

func (e LitUintExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatUint(uint64(e), 10))
}

func (e LitFloatExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatFloat(float64(e), 'g', -1, 64))
}

func (e LitStringExpr) write(w *code.Writer) {
	w.WriteString(strconv.Quote(string(e)))
}

func (e RangeExpr) write(w *code.Writer) {
	if e.low != nil {
		e.low.write(w)
	}

	w.WriteByte(':')

	if e.upp != nil {
		e.upp.write(w)
	}
}

func (e SymbolExpr) write(w *code.Writer) {
	if e.pkg != "" {
		e.pkg.write(w)
		w.WriteByte('.')
	}

	e.id.write(w)

	if len(e.args) > 0 {
		w.WriteByte('[')

		for idx, arg := range e.args {
			if idx > 0 {
				w.WriteString(", ")
			}

			arg.write(w)
		}

		w.WriteByte(']')
	}
}

func (e ParenExpr) write(w *code.Writer) {
	w.WriteByte('(')
	e.expr.write(w)
	w.WriteByte(')')
}

func (e MemberExpr) write(w *code.Writer) {
	e.expr.write(w)
	w.WriteByte('.')
	e.id.write(w)
}

func (e CallExpr) write(w *code.Writer) {
	e.callee.write(w)
	w.WriteByte('(')

	for idx, itm := range e.args {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}

	w.WriteByte(')')
}

func (e IndexExpr) write(w *code.Writer) {
	e.expr.write(w)
	w.WriteByte('[')
	e.idx.write(w)
	w.WriteByte(']')
}

func (e UnaryExpr) write(w *code.Writer) {
	w.WriteString(string(e.op))
	e.expr.write(w)
}

var (
	_ Expr = Nil
	_ Expr = LitBoolExpr(false)
	_ Expr = LitRuneExpr(0)
	_ Expr = LitIntExpr(0)
	_ Expr = LitUintExpr(0)
	_ Expr = LitFloatExpr(0)
	_ Expr = LitStringExpr("")
	_ Expr = RangeExpr{}
	_ Expr = SymbolExpr{}
	_ Expr = ParenExpr{}
	_ Expr = MemberExpr{}
	_ Expr = CallExpr{}
	_ Expr = IndexExpr{}
	_ Expr = UnaryExpr{}
)
