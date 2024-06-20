package golang

import code "github.com/trwk76/gocode"

func Named(id Identifier, args ...Type) NamedType {
	return NamedType{id: id, args: args}
}

func Ptr(target Type) PtrType {
	return PtrType{tgt: target}
}

func Slice(size Expr, item Type) SliceType {
	return SliceType{size: size, itm: item}
}

func Map(key Type, item Type) MapType {
	return MapType{key: key, itm: item}
}

var (
	Any     NamedType = NamedType{id: "any"}
	Bool    NamedType = NamedType{id: "bool"}
	Byte    NamedType = NamedType{id: "byte"}
	Rune    NamedType = NamedType{id: "rune"}
	Int     NamedType = NamedType{id: "int"}
	Int8    NamedType = NamedType{id: "int8"}
	Int16   NamedType = NamedType{id: "int16"}
	Int32   NamedType = NamedType{id: "int32"}
	Int64   NamedType = NamedType{id: "int64"}
	Uint    NamedType = NamedType{id: "uint"}
	Uint8   NamedType = NamedType{id: "uint8"}
	Uint16  NamedType = NamedType{id: "uint16"}
	Uint32  NamedType = NamedType{id: "uint32"}
	Uint64  NamedType = NamedType{id: "uint64"}
	Float32 NamedType = NamedType{id: "float32"}
	Float64 NamedType = NamedType{id: "float64"}
	String  NamedType = NamedType{id: "string"}
)

type (
	Type interface {
		item
	}

	NamedType struct {
		pkg  Identifier
		id   Identifier
		args []Type
	}

	PtrType struct {
		tgt Type
	}

	SliceType struct {
		size Expr
		itm  Type
	}

	MapType struct {
		key Type
		itm Type
	}
)

func (t NamedType) write(w *code.Writer) {
	if t.pkg != "" {
		t.pkg.write(w)
		w.WriteByte('.')
	}

	t.id.write(w)

	if len(t.args) > 0 {
		w.WriteByte('[')

		for idx, arg := range t.args {
			if idx > 0 {
				w.WriteString(", ")
			}

			arg.write(w)
		}

		w.WriteByte(']')
	}
}

func (t PtrType) write(w *code.Writer) {
	w.WriteByte('*')
	t.tgt.write(w)
}

func (t SliceType) write(w *code.Writer) {
	w.WriteByte('[')

	if t.size != nil {
		t.size.write(w)
	}

	w.WriteByte(']')
	t.itm.write(w)
}

func (t MapType) write(w *code.Writer) {
	w.WriteString("map[")
	t.key.write(w)
	w.WriteByte(']')
	t.itm.write(w)
}

var (
	_ Type = NamedType{}
	_ Type = PtrType{}
	_ Type = SliceType{}
	_ Type = MapType{}
)
