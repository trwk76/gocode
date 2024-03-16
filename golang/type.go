package golang

import (
	"strconv"

	code "github.com/trwk76/gocode"
)

type (
	Type interface {
		Item
		simpleType() bool
	}

	NamedType struct {
		Pkg  string
		ID   ID
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
		ID   ID
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
		ID ID
	}
)

var (
	Any     NamedType = NamedType{ID: "any"}
	Bool    NamedType = NamedType{ID: "bool"}
	Byte    NamedType = NamedType{ID: "byte"}
	Rune    NamedType = NamedType{ID: "rune"}
	Int     NamedType = NamedType{ID: "int"}
	Int8    NamedType = NamedType{ID: "int8"}
	Int16   NamedType = NamedType{ID: "int16"}
	Int32   NamedType = NamedType{ID: "int32"}
	Int64   NamedType = NamedType{ID: "int64"}
	Uint    NamedType = NamedType{ID: "uint"}
	Uint8   NamedType = NamedType{ID: "uint8"}
	Uint16  NamedType = NamedType{ID: "uint16"}
	Uint32  NamedType = NamedType{ID: "uint32"}
	Uint64  NamedType = NamedType{ID: "uint64"}
	Float32 NamedType = NamedType{ID: "float32"}
	Float64 NamedType = NamedType{ID: "float64"}
	String  NamedType = NamedType{ID: "string"}
)

func (t NamedType) simpleType() bool {
	return t.Args.simple()
}

func (t NamedType) write(w *code.Writer) {
	if t.Pkg != "" {
		w.WriteString(t.Pkg)
		w.WriteByte('.')
	}

	t.ID.write(w)
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
	w.WriteString("struct {")

	if len(t.Fields) > 0 {
		w.Newline()
		w.Indent(func() {
			tbl := code.Table{}

			for _, fld := range t.Fields {
				cols := make([]string, 0)
				col1 := Render(fld.ID)
				col2 := Render(fld.Type)
				col3 := Render(fld.Tags)

				if col1 == "" {
					cols = append(cols, col2)

					if col3 != "" {
						cols = append(cols, "", col3)
					}
				} else {
					cols = append(cols, col1, col2)

					if col3 != "" {
						cols = append(cols, col3)
					}
				}

				tbl.AddRow(Render(fld.Doc), cols...)
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
	w.WriteString("interface {")

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
	i.ID.write(w)
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
