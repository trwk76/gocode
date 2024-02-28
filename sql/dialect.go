package sql

import "fmt"

type (
	Dialect interface {
		WriteName(w Writer, name Name)
		WriteQuerySet(w Writer, q QuerySet)
		WriteSelect(w Writer, s Select)

		WriteOrCond(w Writer, c OrCond)
	}

	DialectBase struct{}
)

func (DialectBase) WriteQuerySet(w Writer, q QuerySet) {
	q.Lhs.write(w)

	switch q.Op {
	case Union:
		w.WriteString("UNION")
	case Intersect:
		w.WriteString("INTERSECT")
	case Except:
		w.WriteString("EXCEPT")
	default:
		panic(fmt.Errorf("query set operator '%s' is not supported", q.Op))
	}

	if q.All {
		w.WriteString(" ALL")
	}

	w.Newline()

	q.Rhs.write(w)
}

func (d DialectBase) WriteOrCond(w Writer, c OrCond) {
	for idx, op := range c.ops {
		if idx > 0 {
			w.WriteString(" OR ")
		}

		d.WriteCondOperand(w, c.condPrec(), op)
	}
}

func (d DialectBase) WriteCondOperand(w Writer, prec uint8, op Cond) {
	if op.condPred() < prec {
		w.WriteByte('(')
	}

	op.write(w)

	if op.condPred() < prec {
		w.WriteByte(')')
	}
}
