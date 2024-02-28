package sql

type (
	Cond interface {
		Item
		condPrec() uint8
	}

	OrCond struct {
		ops []Cond
	}

	AndCond struct {
		ops []Cond
	}

	NotCond struct {
		op Cond
	}

	BinCond struct {
		op  BinOp
		lhs Expr
		rhs Expr
	}

	BinOp string

	UnCond struct {
		op   UnOp
		expr Expr
	}

	UnOp string
)

const (
	EqOp  BinOp = "="
	NeqOp BinOp = "<>"
	LtOp  BinOp = "<"
	LeqOp BinOp = "<="
	MtOp  BinOp = ">"
	MeqOp BinOp = ">="
)

const (
	IsNullOp    UnOp = "isnull"
	IsNotNullOp UnOp = "isnotnull"
)

func Or(ops ...Cond) Cond {
	o := make([]Cond, 0, len(ops))

	for _, op := range ops {
		if op != nil {
			if tgt, ok := op.(OrCond); ok {
				o = append(o, tgt.ops...)
			} else {
				o = append(o, op)
			}
		}
	}

	switch len(o) {
	case 0:
		return nil
	case 1:
		return o[0]
	}

	return OrCond{ops: o}
}

func (c OrCond) Operands() []Cond {
	return c.ops
}

func (c OrCond) write(w Writer) {
	w.d.WriteOrCond(w, c)
}

func (OrCond)  condPrec() uint8 { return 0 }
func (AndCond) condPrec() uint8 { return 1 }
func (NotCond) condPrec() uint8 { return 2 }
func (BinCond) condPrec() uint8 { return 3 }
func (UnCond)  condPrec() uint8 { return 3 }

var (
	_ Cond = OrCond{}
	_ Cond = AndCond{}
	_ Cond = NotCond{}
	_ Cond = BinCond{}
	_ Cond = UnCond{}
)
