package golang

import "gitlab.trwk.com/go/code"

type (
	Stmt interface {
		Item
		simpleStmt() bool
	}

	ElseStmt interface {
		Stmt
		elseStmt()
	}

	AssignStmt struct {
		Dest Exprs
		Src  Exprs
		Op   AssignOp
	}

	AssignOp string

	BlockStmt []Stmt

	BreakStmt struct{}

	ContinueStmt struct{}

	ExprStmt struct {
		Expr Expr
	}

	IfStmt struct {
		Init *AssignStmt
		Cond Expr
		Then BlockStmt
		Else ElseStmt
	}

	ReturnStmt struct {
		Value Expr
	}
)

const (
	Assign     AssignOp = "="
	DeclAssign AssignOp = ":="
	AddAssign  AssignOp = "+="
	SubAssign  AssignOp = "-="
)

func (s AssignStmt) simpleStmt() bool {
	return s.Dest.simple() && s.Src.simple()
}

func (s AssignStmt) write(w *code.Writer) {
	s.Dest.write(w)
	w.Space()
	w.WriteString(string(s.Op))
	w.Space()
	s.Src.write(w)
}

func (s BlockStmt) simpleStmt() bool {
	return len(s) == 0 || s[0].simpleStmt()
}

func (s BlockStmt) write(w *code.Writer) {
	w.WriteByte('{')

	if l := len(s); l > 0 {
		if l == 1 && s[0].simpleStmt() {
			w.Space()
			s[0].write(w)
			w.Space()
		} else {
			w.Newline()
			w.Indent(func() {
				for idx, itm := range s {
					if idx > 0 {
						w.Newline()
					}

					itm.write(w)
					w.Newline()
				}
			})
		}
	}

	w.WriteByte('}')
}

func (BreakStmt) simpleStmt() bool {
	return true
}

func (BreakStmt) write(w *code.Writer) {
	w.WriteString("break")
}

func (ContinueStmt) simpleStmt() bool {
	return true
}

func (ContinueStmt) write(w *code.Writer) {
	w.WriteString("continue")
}

func (s ExprStmt) simpleStmt() bool {
	return s.Expr.simpleExpr()
}

func (s ExprStmt) write(w *code.Writer) {
	s.Expr.write(w)
}

func (s IfStmt) simpleStmt() bool {
	if s.Init != nil && !s.Init.simpleStmt() {
		return false
	}

	if !s.Cond.simpleExpr() || !s.Then.simpleStmt() {
		return false
	}

	if s.Else != nil && !s.Else.simpleStmt() {
		return false
	}

	return true
}

func (s IfStmt) write(w *code.Writer) {
	w.WriteString("if ")

	if s.Init != nil {
		s.Init.write(w)
		w.WriteString("; ")
	}

	s.Cond.write(w)
	w.Space()
	s.Then.write(w)

	if s.Else != nil {
		w.WriteString(" else ")
		s.Else.write(w)
	}
}

func (s ReturnStmt) simpleStmt() bool {
	return s.Value == nil || s.Value.simpleExpr()
}

func (s ReturnStmt) write(w *code.Writer) {
	w.WriteString("return")

	if s.Value != nil {
		w.Space()
		s.Value.write(w)
	}
}

func (BlockStmt) elseStmt() {}
func (IfStmt) elseStmt()    {}

var (
	_ Stmt     = AssignStmt{}
	_ ElseStmt = BlockStmt{}
	_ Stmt     = BreakStmt{}
	_ Stmt     = ContinueStmt{}
	_ Stmt     = ExprStmt{}
	_ ElseStmt = IfStmt{}
	_ Stmt     = ReturnStmt{}
)
