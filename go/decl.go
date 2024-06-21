package golang

import code "github.com/trwk76/gocode"

func (f *File) Types(cfg func(t *TypeDecls)) *TypeDecls {
	res := &TypeDecls{}

	f.decls = append(f.decls, res)

	return config(res, cfg)
}

func (t *TypeDecls) Add(comment Comment, id Identifier, genParams GenParamDecls, def Type, alias bool) {
	t.items = append(t.items, TypeDecl{
		comm:   comment,
		id:     id,
		params: genParams,
		def:    def,
		alias:  alias,
	})
}

func (f *File) Funcs(cfg func(f *FuncDecls)) *FuncDecls {
	res := &FuncDecls{}

	f.decls = append(f.decls, res)

	return config(res, cfg)
}

func (f *FuncDecls) Add(comm Comment, recv ParamDecl, id Identifier, genParams GenParamDecls, params ParamDecls, ret ParamDecls, body func(b *BlockStmt)) {
	var (
		rcv *RecvDecl
		b   BlockStmt
	)

	if recv.typ != nil {
		tmp := RecvDecl(recv)
		rcv = &tmp
	}

	config(&b, body)

	f.items = append(f.items, FuncDecl{
		comm:    comm,
		rcv:     rcv,
		id:      id,
		gparams: genParams,
		params:  InParamDecls(params),
		ret:     OutParamDecls(ret),
		body:    b,
	})
}

func (f *File) Vars(cfg func(v *VarDecls)) *VarDecls {
	res := &VarDecls{}

	f.decls = append(f.decls, res)

	return config(res, cfg)
}

func (v *VarDecls) Add(comm Comment, id Identifier, typ Type, init Expr) {
	v.items = append(v.items, VarDecl{
		comm: comm,
		id:   id,
		typ:  typ,
		init: init,
	})
}

func (f *File) Consts(cfg func(c *ConstDecls)) *ConstDecls {
	res := &ConstDecls{}

	f.decls = append(f.decls, res)

	return config(res, cfg)
}

func (c *ConstDecls) Add(comm Comment, id Identifier, typ Type, init Expr) {
	c.items = append(c.items, VarDecl{
		comm: comm,
		id:   id,
		typ:  typ,
		init: init,
	})
}

func GenParam(id Identifier, ext bool, constraint Type) GenParamDecl {
	return GenParamDecl{id: id, ext: ext, cnst: constraint}
}

func Param(id Identifier, t Type, variadic bool) ParamDecl {
	return ParamDecl{id: id, typ: t, vargs: variadic}
}

type (
	TypeDecls struct {
		items []TypeDecl
	}

	FuncDecls struct {
		items []FuncDecl
	}

	VarDecls struct {
		items []VarDecl
	}

	ConstDecls struct {
		items []VarDecl
	}

	TypeDecl struct {
		comm   Comment
		id     Identifier
		params GenParamDecls
		def    Type
		alias  bool
	}

	FuncDecl struct {
		comm    Comment
		rcv     *RecvDecl
		id      Identifier
		gparams GenParamDecls
		params  InParamDecls
		ret     OutParamDecls
		body    BlockStmt
	}

	VarDecl struct {
		comm Comment
		id   Identifier
		typ  Type
		init Expr
	}

	GenParamDecls []GenParamDecl
	ParamDecls    []ParamDecl
	RecvDecl      ParamDecl
	InParamDecls  []ParamDecl
	OutParamDecls []ParamDecl

	GenParamDecl struct {
		id   Identifier
		ext  bool
		cnst Type
	}

	ParamDecl struct {
		id    Identifier
		typ   Type
		vargs bool
	}

	decls interface {
		write(w *code.Writer)
	}

	decl interface {
		item
		simpleDecl() bool
		row() code.Row
	}
)

func (*TypeDecls) decl()  {}
func (*FuncDecls) decl()  {}
func (*VarDecls) decl()   {}
func (*ConstDecls) decl() {}

func (d *TypeDecls) write(w *code.Writer) {
	writeSection(w, d.items, "type")
}

func (d TypeDecl) simpleDecl() bool {
	return d.def.simpleType()
}

func (d TypeDecl) row() code.Row {
	pfx := ""

	if d.alias {
		pfx = "= "
	}

	return code.Row{
		Prefix: commentRenderer(string(d.comm)),
		Cols:   []string{renderLine(d.id, d.params), pfx + renderLine(d.def)},
	}
}

func (d TypeDecl) write(w *code.Writer, line bool) {
	d.id.write(w, line)
	d.params.write(w, line)
	w.Space()

	if d.alias {
		w.WriteString("= ")
	}

	d.def.write(w, line)
	w.Newline()
}

func (d *FuncDecls) write(w *code.Writer) {
	first := true
	items := d.items

	for len(items) > 0 {
		if first {
			first = false
		} else {
			w.Newline()
		}

		if cnt := countSimple(items); cnt > 1 {
			rows := make([]code.Row, cnt)

			for i := 0; i < cnt; i++ {
				rows[i] = items[i].row()
			}

			w.Table(rows...)
			items = items[cnt:]
		} else {
			items[0].write(w, false)
			items = items[1:]
		}
	}
}

func (d FuncDecl) simpleDecl() bool {
	return d.body.simpleStmt()
}

func (d FuncDecl) row() code.Row {
	col0 := "func"

	if d.rcv != nil {
		col0 += " " + renderLine(*d.rcv)
	}

	return code.Row{
		Prefix: commentRenderer(string(d.comm)),
		Cols: []string{
			col0,
			renderLine(d.id, d.gparams, d.params, d.ret),
			renderLine(d.body),
		},
	}
}

func (d FuncDecl) write(w *code.Writer, line bool) {
	w.WriteString("func ")

	if d.rcv != nil {
		d.rcv.write(w, line)
		w.Space()
	}

	d.id.write(w, line)
	d.gparams.write(w, line)
	d.params.write(w, line)
	d.ret.write(w, line)
	w.Space()
	d.body.write(w, line)
	w.Newline()
}

func (d *VarDecls) write(w *code.Writer) {
	writeSection(w, d.items, "var")
}

func (d *ConstDecls) write(w *code.Writer) {
	writeSection(w, d.items, "const")
}

func (d GenParamDecls) write(w *code.Writer, line bool) {
	if len(d) < 1 {
		return
	}

	w.WriteByte('[')

	for idx, itm := range d {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w, true)
	}

	w.WriteByte(']')
}

func (d GenParamDecl) write(w *code.Writer, line bool) {
	d.id.write(w, line)
	w.Space()

	if d.ext {
		w.WriteByte('~')
	}

	d.cnst.write(w, line)
}

func (d ParamDecl) write(w *code.Writer, line bool) {
	d.id.write(w, line)

	if d.id != "" && d.typ != nil {
		w.Space()
	}

	if d.vargs {
		w.WriteString("...")
	}

	if d.typ != nil {
		d.typ.write(w, line)
	}
}

func (d RecvDecl) write(w *code.Writer, line bool) {
	w.WriteByte('(')
	ParamDecl(d).write(w, line)
	w.WriteByte(')')
}

func (d InParamDecls) write(w *code.Writer, line bool) {
	w.WriteByte('(')

	for idx, itm := range d {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w, line)
	}

	w.WriteByte(')')
}

func (d OutParamDecls) write(w *code.Writer, line bool) {
	switch len(d) {
	case 0:
		return
	case 1:
		if d[0].id == "" {
			w.Space()
			d[0].typ.write(w, line)
			return
		}
	}

	w.Space()
	w.WriteByte('(')

	for idx, itm := range d {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w, line)
	}

	w.WriteByte(')')
}

func (d VarDecl) simpleDecl() bool {
	return d.init == nil || d.init.simpleExpr()
}

func config[T any](item T, cfg func(i T)) T {
	if cfg != nil {
		cfg(item)
	}

	return item
}

func writeSection[I decl](w *code.Writer, items []I, name string) {
	switch len(items) {
	case 0:
		return
	case 1:
		w.WriteString(name)
		w.Space()
		items[0].write(w, false)
		w.Newline()
		return
	}

	w.WriteString(name)
	w.WriteString(" (")
	w.Newline()

	w.Indent(func(w *code.Writer) {
		first := true

		for len(items) > 0 {
			if first {
				first = false
			} else {
				w.Newline()
			}

			if cnt := countSimple(items); cnt > 1 {
				rows := make([]code.Row, cnt)

				for i := 0; i < cnt; i++ {
					rows[i] = items[i].row()
				}

				w.Table(rows...)
				items = items[cnt:]
			} else {
				items[0].write(w, false)
				items = items[1:]
			}
		}
	})

	w.WriteByte(')')
	w.Newline()
}

func countSimple[I decl](items []I) int {
	res := 0

	for _, itm := range items {
		if itm.simpleDecl() {
			res++
		} else {
			break
		}
	}

	return res
}
