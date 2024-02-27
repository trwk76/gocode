package sql

type (
	Exprs []Expr

	Expr interface {
		Item
		expr()
	}
)

func (e Exprs) write(w Writer) {
	for idx, itm := range e {
		if idx > 0 {
			w.WriteString(", ")
		}

		itm.write(w)
	}
}
