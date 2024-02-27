package golang

import (
	"strconv"

	code "github.com/trwk76/gocode"
)

type (
	Decls []Decl

	Decl interface {
		Item
		simpleDecl() bool
	}

	PackageDecl string

	ImportDecls []ImportDecl

	ImportDecl struct {
		Path  string
		Alias string
	}

	TypeDecls []TypeDecl

	TypeDecl struct {
		Doc    Comment
		ID     ID
		Params GenParams
		Type   Type
	}

	FuncDecl struct {
		Doc       Comment
		Rcv       Receiver
		ID        ID
		GenParams GenParams
		Params    Params
		Return    Return
		Body      BlockStmt
	}

	ConstDecls []VarDecl
	VarDecls   []VarDecl

	VarDecl struct {
		Doc   Comment
		ID    ID
		Type  Type
		Value Expr
	}

	GenParams []GenParam

	GenParam struct {
		ID         ID
		Constraint TypeConst
	}

	Receiver = *Param
	Params   []Param
	Return   []Param

	Param struct {
		ID   ID
		Var  bool
		Type Type
	}
)

func (i Comment) simpleDecl() bool {
	return false
}

func (i PackageDecl) simpleDecl() bool {
	return true
}

func (i PackageDecl) write(w *code.Writer) {
	w.WriteString("package ")
	w.WriteString(string(i))
}

func (i ImportDecls) simpleDecl() bool {
	return len(i) < 2
}

func (i ImportDecls) write(w *code.Writer) {
	switch len(i) {
	case 0:
	case 1:
		w.WriteString("import ")

		if i[0].Alias != "" {
			w.WriteString(i[0].Alias)
			w.Space()
		}

		w.WriteString(strconv.Quote(i[0].Path))
	default:
		alias := false

		for _, itm := range i {
			if itm.Alias != "" {
				alias = true
			}
		}

		w.WriteString("import (")
		w.Newline()
		w.Indent(func() {
			tbl := code.Table{}

			for _, itm := range i {
				if alias {
					tbl.AddRow("", itm.Alias, strconv.Quote(itm.Path))
				} else {
					tbl.AddRow("", strconv.Quote(itm.Path))
				}
			}

			tbl.Write(w)
		})
		w.WriteByte(')')
	}
}

func (i TypeDecls) simpleDecl() bool {
	switch len(i) {
	case 0:
		return true
	case 1:
		return i[0].Type.simpleType()
	}

	return false
}

func (i TypeDecls) simpleStmt() bool {
	return i.simpleDecl()
}

func (i TypeDecls) write(w *code.Writer) {
	switch len(i) {
	case 0:
	case 1:
		i[0].Doc.write(w)
		w.WriteString("type ")
		i[0].ID.write(w)
		i[0].Params.write(w)
		w.Space()
		i[0].Type.write(w)
	default:
		w.WriteString("type (")
		w.Newline()
		w.Indent(func() {
			for idx, itm := range i {
				if idx > 0 {
					w.Newline()
				}

				itm.Doc.write(w)
				itm.ID.write(w)
				itm.Params.write(w)
				w.Space()
				itm.Type.write(w)
				w.Newline()
			}
		})
		w.WriteByte(')')
	}
}

func (i FuncDecl) simpleDecl() bool {
	return i.Rcv.simple() && i.Params.simple() && i.Return.simple() && i.Body.simpleStmt()
}

func (i FuncDecl) simpleStmt() bool {
	return i.simpleDecl()
}

func (i FuncDecl) write(w *code.Writer) {
	i.Doc.write(w)
	w.WriteString("func ")

	if i.Rcv != nil {
		w.WriteByte('(')
		i.Rcv.write(w)
		w.WriteString(") ")
	}

	i.ID.write(w)
	i.GenParams.write(w)
	i.Params.write(w)
	i.Return.write(w)
	w.Space()
	i.Body.write(w)
}

func (i VarDecls) simpleDecl() bool {
	switch len(i) {
	case 0:
		return true
	case 1:
		return i[0].simple()
	}

	return false
}

func (i VarDecls) simpleStmt() bool {
	return i.simpleDecl()
}

func (i VarDecls) write(w *code.Writer) {
	switch len(i) {
	case 0:
	case 1:
		i[0].Doc.write(w)
		w.WriteString("var ")
		i[0].ID.write(w)
		w.Space()
		i[0].Type.write(w)

		if i[0].Value != nil {
			w.WriteString(" = ")
			i[0].Value.write(w)
		}
	default:
		w.WriteString("var (")
		w.Newline()
		w.Indent(func() {
			tbl := code.Table{}

			for _, itm := range i {
				cols := []string{Render(itm.ID), Render(itm.Type)}

				if itm.Value != nil {
					cols = append(cols, "= "+Render(itm.Value))
				}

				tbl.AddRow(Render(itm.Doc), cols...)
			}

			tbl.Write(w)
		})
		w.WriteByte(')')
	}
}

func (i ConstDecls) simpleDecl() bool {
	switch len(i) {
	case 0:
		return true
	case 1:
		return i[0].simple()
	}

	return false
}

func (i ConstDecls) simpleStmt() bool {
	return i.simpleDecl()
}

func (i ConstDecls) write(w *code.Writer) {
	switch len(i) {
	case 0:
	case 1:
		i[0].Doc.write(w)
		w.WriteString("const ")
		i[0].ID.write(w)
		w.Space()
		i[0].Type.write(w)

		if i[0].Value != nil {
			w.WriteString(" = ")
			i[0].Value.write(w)
		}
	default:
		w.WriteString("const (")
		w.Newline()
		w.Indent(func() {
			tbl := code.Table{}

			for _, itm := range i {
				cols := []string{Render(itm.ID), Render(itm.Type)}

				if itm.Value != nil {
					cols = append(cols, "= "+Render(itm.Value))
				}

				tbl.AddRow(Render(itm.Doc), cols...)
			}

			tbl.Write(w)
		})
		w.WriteByte(')')
	}
}

func (i VarDecl) simple() bool {
	return i.Type.simpleType() && (i.Value == nil || i.Value.simpleExpr())
}

func (i *Param) simple() bool {
	return i == nil || i.Type.simpleType()
}

func (p GenParams) write(w *code.Writer) {
	if len(p) < 1 {
		return
	}

	w.WriteByte('[')

	for idx, itm := range p {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}

	w.WriteByte(']')
}

func (p GenParam) write(w *code.Writer) {
	p.ID.write(w)
	w.Space()
	p.Constraint.write(w)
}

func (p Params) simple() bool {
	for _, itm := range p {
		if !itm.Type.simpleType() {
			return false
		}
	}

	return true
}

func (p Params) write(w *code.Writer) {
	w.WriteByte('(')

	for idx, itm := range p {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}

	w.WriteByte(')')
}

func (p Return) simple() bool {
	for _, itm := range p {
		if !itm.Type.simpleType() {
			return false
		}
	}

	return true
}

func (p Return) write(w *code.Writer) {
	if len(p) < 1 {
		return
	}

	w.Space()

	if len(p) == 1 && p[0].ID == "" {
		p[0].Type.write(w)
		return
	}

	w.WriteByte('(')

	for idx, itm := range p {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}

	w.WriteByte(')')
}

func (i Param) write(w *code.Writer) {
	i.ID.write(w)

	if i.Var {
		w.WriteString("...")
	}

	w.Space()
	i.Type.write(w)
}

var (
	_ Decl = Comment("")
	_ Decl = PackageDecl("")
	_ Decl = ImportDecls{}
	_ Decl = TypeDecls{}
	_ Stmt = TypeDecls{}
	_ Decl = FuncDecl{}
	_ Stmt = FuncDecl{}
	_ Decl = VarDecls{}
	_ Stmt = VarDecls{}
	_ Decl = ConstDecls{}
	_ Stmt = ConstDecls{}
)
