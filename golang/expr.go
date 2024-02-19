package golang

import (
	"strconv"

	"gitlab.trwk.com/go/code"
)

type (
	Exprs []Expr

	Expr interface {
		Item
		simpleExpr() bool
	}

	NilExpr    struct{}
	BoolExpr   bool
	RuneExpr   rune
	IntExpr    int64
	UintExpr   uint64
	FloatExpr  float64
	StringExpr string

	StructExpr struct {
		Type   Type
		Fields StructExprFields
	}

	StructExprFields []StructExprField

	StructExprField struct {
		Name  string
		Value Expr
	}

	MapExpr struct {
		Type   Type
		Values MapExprValues
	}

	MapExprValues []MapExprValue

	MapExprValue struct {
		Key   Expr
		Value Expr
	}

	FuncExpr struct {
		Params Params
		Return Return
		Body   BlockStmt
	}

	SymbolExpr struct {
		Name string
		Pkg  string
		Args GenArgs
	}

	MemberExpr struct {
		Expr Expr
		Name string
	}

	CallExpr struct {
		Func Expr
		Args Args
	}

	Args []Expr

	UnOperator string

	UnaryExpr struct {
		Op   UnOperator
		Expr Expr
	}

	BinOperator string

	BinaryExpr struct {
		Op  BinOperator
		Lhs Expr
		Rhs Expr
	}
)

const (
	Identity UnOperator = "+"
	Negate   UnOperator = "-"
	Not      UnOperator = "!"
	Compl    UnOperator = "^"
	AddrOf   UnOperator = "&"
	Deref    UnOperator = "*"
)

const (
	Multiply   BinOperator = "*"
	Divide     BinOperator = "/"
	Remainder  BinOperator = "%"
	Add        BinOperator = "+"
	Subtract   BinOperator = "-"
	ShiftLeft  BinOperator = "<<"
	ShiftRight BinOperator = ">>"
	BitAnd     BinOperator = "&"
	BitXor     BinOperator = "^"
	BitOr      BinOperator = "|"
	Equal      BinOperator = "=="
	NotEqual   BinOperator = "!="
	LessThan   BinOperator = "<"
	LessEqual  BinOperator = "<="
	MoreThan   BinOperator = ">"
	MoreEqual  BinOperator = ">="
	LogAnd     BinOperator = "&&"
	LogOr      BinOperator = "||"
)

var (
	Nil    NilExpr    = NilExpr{}
	False  BoolExpr   = BoolExpr(false)
	True   BoolExpr   = BoolExpr(true)
	Ignore SymbolExpr = SymbolExpr{Name: "_"}
)

func (e Exprs) simple() bool {
	for _, itm := range e {
		if !itm.simpleExpr() {
			return false
		}
	}

	return true
}

func (e Exprs) write(w *code.Writer) {
	for idx, itm := range e {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}
}

func (e NilExpr) write(w *code.Writer) {
	w.WriteString("nil")
}

func (e BoolExpr) write(w *code.Writer) {
	if bool(e) {
		w.WriteString("true")
	} else {
		w.WriteString("false")
	}
}

func (e RuneExpr) write(w *code.Writer) {
	w.WriteString(strconv.QuoteRune(rune(e)))
}

func (e IntExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatInt(int64(e), 10))
}

func (e UintExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatUint(uint64(e), 10))
}

func (e FloatExpr) write(w *code.Writer) {
	w.WriteString(strconv.FormatFloat(float64(e), 'g', -1, 64))
}

func (e StringExpr) write(w *code.Writer) {
	w.WriteString(strconv.Quote(string(e)))
}

func (e StructExpr) write(w *code.Writer) {
	if e.Type != nil {
		e.Type.write(w)
	}

	w.WriteByte('{')

	if len(e.Fields) > 0 {
		w.Newline()
		w.Indent(func() {
			for _, itm := range e.Fields {
				w.WriteString(itm.Name)
				w.WriteString(": ")
				itm.Value.write(w)
				w.WriteByte(',')
				w.Newline()
			}
		})
	}

	w.WriteByte('}')
}

func (e MapExpr) write(w *code.Writer) {
	if e.Type != nil {
		e.Type.write(w)
	}

	w.WriteByte('{')

	if len(e.Values) > 0 {
		w.Newline()
		w.Indent(func() {
			for _, itm := range e.Values {
				itm.Key.write(w)
				w.WriteString(": ")
				itm.Value.write(w)
				w.WriteByte(',')
				w.Newline()
			}
		})
	}

	w.WriteByte('}')
}

func (e FuncExpr) write(w *code.Writer) {
	w.WriteString("func ")
	e.Params.write(w)
	e.Return.write(w)
	w.Space()
	e.Body.write(w)
}

func (e SymbolExpr) write(w *code.Writer) {
	if e.Pkg != "" {
		w.WriteString(e.Pkg)
		w.WriteByte('.')
	}

	w.WriteString(e.Name)
	e.Args.write(w)
}

func (e MemberExpr) write(w *code.Writer) {
	e.Expr.write(w)
	w.WriteByte('.')
	w.WriteString(e.Name)
}

func (e CallExpr) write(w *code.Writer) {
	e.Func.write(w)
	e.Args.write(w)
}

func (e UnaryExpr) write(w *code.Writer) {
	w.WriteString(string(e.Op))
	e.Expr.write(w)
}

func (e BinaryExpr) write(w *code.Writer) {
	e.Lhs.write(w)
	w.Space()
	w.WriteString(string(e.Op))
	w.Space()
	e.Rhs.write(w)
}

func (a Args) simple() bool {
	for _, itm := range a {
		if !itm.simpleExpr() {
			return false
		}
	}

	return true
}

func (a Args) write(w *code.Writer) {
	w.WriteByte('(')

	if a.simple() {
		for idx, itm := range a {
			if idx > 0 {
				w.WriteString(", ")
			}

			itm.write(w)
		}
	} else {
		w.Newline()
		w.Indent(func() {
			for _, itm := range a {
				itm.write(w)
				w.WriteByte(',')
				w.Newline()
			}
		})
	}

	w.WriteByte(')')
}

func (NilExpr) simpleExpr() bool      { return true }
func (BoolExpr) simpleExpr() bool     { return true }
func (RuneExpr) simpleExpr() bool     { return true }
func (IntExpr) simpleExpr() bool      { return true }
func (UintExpr) simpleExpr() bool     { return true }
func (FloatExpr) simpleExpr() bool    { return true }
func (StringExpr) simpleExpr() bool   { return true }
func (SymbolExpr) simpleExpr() bool   { return true }
func (e StructExpr) simpleExpr() bool { return len(e.Fields) < 1 }
func (e MapExpr) simpleExpr() bool    { return len(e.Values) < 1 }
func (e FuncExpr) simpleExpr() bool   { return e.Body.simpleStmt() }
func (e MemberExpr) simpleExpr() bool { return e.Expr.simpleExpr() }
func (e CallExpr) simpleExpr() bool   { return e.Func.simpleExpr() && e.Args.simple() }
func (e UnaryExpr) simpleExpr() bool  { return e.Expr.simpleExpr() }
func (e BinaryExpr) simpleExpr() bool { return e.Lhs.simpleExpr() && e.Rhs.simpleExpr() }

var (
	_ Expr = NilExpr{}
	_ Expr = False
	_ Expr = RuneExpr(' ')
	_ Expr = IntExpr(1)
	_ Expr = UintExpr(1)
	_ Expr = FloatExpr(1)
	_ Expr = StringExpr("")
	_ Expr = SymbolExpr{}
	_ Expr = StructExpr{}
	_ Expr = MapExpr{}
	_ Expr = FuncExpr{}
	_ Expr = MemberExpr{}
	_ Expr = CallExpr{}
	_ Expr = UnaryExpr{}
	_ Expr = BinaryExpr{}
)
