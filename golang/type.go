package golang

import (
	"strconv"

	"gitlab.trwk.com/go/code"
)

type (
	Type interface {
		Item
		simpleType() bool
	}

	NamedType struct {
		Pkg  string
		Name string
		Args GenArgs
	}

	PtrType struct {
		Target Type
	}

	SliceType struct {
		Len  Expr
		Item Type
	}

	MapType struct {
		Key Type
		Val Type
	}

	StructType struct {
		Fields Fields
	}

	InterfaceType struct {
		Consts TypeConsts
		Funcs  InterfaceFuncs
	}

	FuncType struct {
		Params Params
		Return Return
	}

	GenArgs []Type

	Fields []Field

	Field struct {
		Doc  Comment
		Name string
		Type Type
		Tags Tags
	}

	Tags []Tag

	Tag struct {
		Name string
		Val  string
	}

	TypeConsts []TypeConst

	TypeConst struct {
		Derived bool
		Type    Type
	}

	InterfaceFuncs []InterfaceFunc

	InterfaceFunc struct {
		FuncType
		Name string
	}
)

var (
	Any     NamedType = NamedType{Name: "any"}
	Bool    NamedType = NamedType{Name: "bool"}
	Byte    NamedType = NamedType{Name: "byte"}
	Rune    NamedType = NamedType{Name: "rune"}
	Int     NamedType = NamedType{Name: "int"}
	Int8    NamedType = NamedType{Name: "int8"}
	Int16   NamedType = NamedType{Name: "int16"}
	Int32   NamedType = NamedType{Name: "int32"}
	Int64   NamedType = NamedType{Name: "int64"}
	Uint    NamedType = NamedType{Name: "uint"}
	Uint8   NamedType = NamedType{Name: "uint8"}
	Uint16  NamedType = NamedType{Name: "uint16"}
	Uint32  NamedType = NamedType{Name: "uint32"}
	Uint64  NamedType = NamedType{Name: "uint64"}
	Float32 NamedType = NamedType{Name: "float32"}
	Float64 NamedType = NamedType{Name: "float64"}
	String  NamedType = NamedType{Name: "string"}
)

func (t NamedType) simpleType() bool {
	return t.Args.simple()
}

func (t NamedType) write(w *code.Writer) {
	if t.Pkg != "" {
		w.WriteString(t.Pkg)
		w.WriteByte('.')
	}

	w.WriteString(t.Name)
	t.Args.write(w)
}

func (t PtrType) simpleType() bool {
	return t.Target.simpleType()
}

func (t PtrType) write(w *code.Writer) {
	w.WriteByte('*')
	t.Target.write(w)
}

func (t SliceType) simpleType() bool {
	return (t.Len == nil || t.Len.simpleExpr()) && t.Item.simpleType()
}

func (t SliceType) write(w *code.Writer) {
	w.WriteByte('[')

	if t.Len != nil {
		t.Len.write(w)
	}

	w.WriteByte(']')
	t.Item.write(w)
}

func (t MapType) simpleType() bool {
	return t.Key.simpleType() && t.Val.simpleType()
}

func (t MapType) write(w *code.Writer) {
	w.WriteString("map[")
	t.Key.write(w)
	w.WriteByte(']')
	t.Val.write(w)
}

func (t StructType) simpleType() bool {
	return len(t.Fields) < 1
}

func (t StructType) write(w *code.Writer) {
	w.WriteString("struct{")

	if len(t.Fields) > 0 {
		w.Newline()
		w.Indent(func() {
			tbl := code.Table{}

			for _, fld := range t.Fields {
				col1 := fld.Name
				col2 := Render(fld.Type)

				if col1 == "" {
					col1 = col2
					col2 = ""
				}

				tbl.AddRow(
					Render(fld.Doc),
					col1,
					col2,
					Render(fld.Tags),
				)
			}

			tbl.Write(w)
		})
	}

	w.WriteByte('}')
}

func (t InterfaceType) simpleType() bool {
	return len(t.Consts) < 1 && len(t.Funcs) < 1
}

func (t InterfaceType) write(w *code.Writer) {
	w.WriteString("interface{")

	if len(t.Consts) > 0 || len(t.Funcs) > 0 {
		w.Newline()
		w.Indent(func() {
			t.Consts.write(w)

			if len(t.Consts) > 0 && len(t.Funcs) > 0 {
				w.Newline()
			}

			t.Funcs.write(w)
		})
	}

	w.WriteByte('}')
}

func (t FuncType) simpleType() bool {
	return t.Params.simple() && t.Return.simple()
}

func (t FuncType) write(w *code.Writer) {
	w.WriteString("func")
	t.Params.write(w)
	t.Return.write(w)
}

func (i GenArgs) simple() bool {
	for _, itm := range i {
		if !itm.simpleType() {
			return false
		}
	}

	return true
}

func (i GenArgs) write(w *code.Writer) {
	if len(i) < 1 {
		return
	}

	w.WriteByte('[')

	for idx, itm := range i {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}

	w.WriteByte(']')
}

func (i TypeConsts) write(w *code.Writer) {
	for idx, itm := range i {
		if idx > 0 {
			w.WriteString(" | ")
		}

		itm.write(w)
	}
}

func (i TypeConst) write(w *code.Writer) {
	if i.Derived {
		w.WriteByte('~')
	}

	i.Type.write(w)
}

func (i Tags) write(w *code.Writer) {
	if len(i) < 1 {
		return
	}

	w.WriteByte('`')

	for idx, itm := range i {
		if idx > 0 {
			w.Space()
		}

		itm.write(w)
	}

	w.WriteByte('`')
}

func (i Tag) write(w *code.Writer) {
	w.WriteString(i.Name)
	w.WriteByte(':')
	w.WriteString(strconv.Quote(i.Val))
}

func (i InterfaceFuncs) write(w *code.Writer) {
	for _, itm := range i {
		itm.write(w)
		w.Newline()
	}
}

func (i InterfaceFunc) write(w *code.Writer) {
	w.WriteString(i.Name)
	i.Params.write(w)
	i.Return.write(w)
}

var (
	_ Type = NamedType{}
	_ Type = PtrType{}
	_ Type = SliceType{}
	_ Type = MapType{}
	_ Type = StructType{}
	_ Type = InterfaceType{}
	_ Type = FuncType{}
)
