package sql

import "fmt"

type (
	Dialect interface {
		WriteName(w Writer, name Name)
		WriteQuerySet(w Writer, q QuerySet)
		WriteSelect(w Writer, s Select)
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
