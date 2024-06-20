package golang

import code "github.com/trwk76/gocode"

func Assign(dest []Expr, src []Expr) AssignStmt {
	return AssignStmt{op: simpleAssignOp, dest: dest, src: src}
}

func DeclAssign(dest []Expr, src []Expr) AssignStmt {
	return AssignStmt{op: declAssignOp, dest: dest, src: src}
}

func AddAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: addAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func SubAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: subAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func MulAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: mulAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func DivAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: divAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func ModAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: modAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func AndAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: andAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func OrAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: orAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func XorAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: xorAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func ShiftLeftAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: shlAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func ShiftRightAssign(dest Expr, src Expr) AssignStmt {
	return AssignStmt{op: shrAssignOp, dest: []Expr{dest}, src: []Expr{src}}
}

func Block(stmts ...Stmt) BlockStmt {
	return BlockStmt(stmts)
}

func If(init initStmt, cond Expr, thenStmt BlockStmt, elseStmt elseStmt) IfStmt {
	return IfStmt{cond: cond, tstmt: thenStmt, estmt: elseStmt}
}

func Return(val Expr) ReturnStmt {
	return ReturnStmt{val: val}
}

var (
	Break    BreakStmt
	Continue ContinueStmt
)

type (
	Stmt interface {
		item
		newLineAfter() bool
	}

	AssignStmt struct {
		op   assignOp
		dest []Expr
		src  []Expr
	}

	BlockStmt []Stmt

	BreakStmt    struct{}
	ContinueStmt struct{}

	IfStmt struct {
		init  initStmt
		cond  Expr
		tstmt BlockStmt
		estmt elseStmt
	}

	ReturnStmt struct {
		val Expr
	}

	initStmt interface {
		Stmt
		initStmt()
	}

	elseStmt interface {
		Stmt
		elseStmt()
	}

	assignOp string
)

const (
	simpleAssignOp assignOp = "="
	declAssignOp   assignOp = ":="
	addAssignOp    assignOp = "+="
	subAssignOp    assignOp = "-="
	mulAssignOp    assignOp = "*="
	divAssignOp    assignOp = "/="
	modAssignOp    assignOp = "%="
	andAssignOp    assignOp = "&="
	orAssignOp     assignOp = "|="
	xorAssignOp    assignOp = "^="
	shlAssignOp    assignOp = "<<="
	shrAssignOp    assignOp = ">>="
)

func (AssignStmt) newLineAfter() bool { return false }
func (AssignStmt) initStmt()          {}

func (s AssignStmt) write(w *code.Writer) {
	for idx, itm := range s.dest {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}

	w.Space()
	w.WriteString(string(s.op))
	w.Space()

	for idx, itm := range s.src {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}
}

func (s BlockStmt) newLineAfter() bool { return len(s) > 0 }
func (s BlockStmt) elseStmt()          {}

func (s BlockStmt) write(w *code.Writer) {
	w.WriteByte('{')

	if len(s) > 0 {
		w.Newline()
		w.Indent(func(w *code.Writer) {
			nla := false

			for _, itm := range s {
				if nla {
					w.Newline()
				}

				itm.write(w)
				w.Newline()
				nla = itm.newLineAfter()
			}
		})
	}

	w.WriteByte('}')
}

func (BreakStmt) newLineAfter() bool { return false }

func (BreakStmt) write(w *code.Writer) {
	w.WriteString("break")
}

func (ContinueStmt) newLineAfter() bool { return false }

func (ContinueStmt) write(w *code.Writer) {
	w.WriteString("continue")
}

func (s IfStmt) newLineAfter() bool { return true }
func (s IfStmt) elseStmt()          {}

func (s IfStmt) write(w *code.Writer) {
	w.WriteString("if ")

	if s.init != nil {
		s.init.write(w)
		w.WriteString("; ")
	}

	if s.cond != nil {
		s.cond.write(w)
		w.Space()
	}

	s.tstmt.write(w)

	if s.estmt != nil {
		w.WriteString(" else ")
		s.estmt.write(w)
	}
}

func (ReturnStmt) newLineAfter() bool { return false }

func (s ReturnStmt) write(w *code.Writer) {
	w.WriteString("return")

	if s.val != nil {
		w.Space()
		s.val.write(w)
	}
}

var (
	_ initStmt = AssignStmt{}
	_ elseStmt = BlockStmt{}
	_ Stmt     = BreakStmt{}
	_ Stmt     = ContinueStmt{}
	_ elseStmt = IfStmt{}
	_ Stmt     = ReturnStmt{}
)
