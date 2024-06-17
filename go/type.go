package golang

type (
	Type interface {
		item
		simpleType() bool
	}

	NamedType struct {
		ID      ID
		Pkg     ID
		GenArgs []Type
	}

	PtrType struct {
		Elem Type
	}

	SliceType struct {
		Size Expr
		Elem Type
	}

	MapType struct {
		Key  Type
		Elem Type
	}
)

