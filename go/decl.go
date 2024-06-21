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
	Decl interface {
		item
		decl()
	}

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
)

func (*TypeDecls)  decl() {}
func (*FuncDecls)  decl() {}
func (*VarDecls)   decl() {}
func (*ConstDecls) decl() {}

func (d *TypeDecls) write(w *code.Writer) {

}

func (d *FuncDecls) write(w *code.Writer) {

}

func (d *VarDecls) write(w *code.Writer) {

}

func (d *ConstDecls) write(w *code.Writer) {
	
}

func config[T any](item T, cfg func(i T)) T {
	if cfg != nil {
		cfg(item)
	}

	return item
}
