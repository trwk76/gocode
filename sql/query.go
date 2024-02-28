package sql

import code "github.com/trwk76/gocode"

type (
	QueryStmt struct {
		With  CTEs
		Query Query
	}

	Query interface {
		SelectSource
		qry()
	}

	QuerySet struct {
		Op  QuerySetOp
		All bool
		Lhs Query
		Rhs Query
	}

	QuerySetOp string

	Select struct {
		Distinct bool
		Columns  SelectColumns
		From     SelectFrom
		Joins    Joins
		GroupBy  Exprs
		OrderBy  OrderClauses
		Where    Cond
		Offset   uint64
		Limit    uint64
	}

	SelectColumns []SelectColumn

	SelectColumn struct {
		Expr  Expr
		Alias Name
	}

	SelectFrom struct {
		Source SelectSource
		Alias  Name
	}

	SelectSource interface {
		Item
		selectSrc()
	}

	Joins []Join

	Join struct {
		Type   JoinType
		Source SelectSource
		Alias  Name
		On     Cond
	}

	JoinType string

	OrderClauses []OrderClause

	OrderClause struct {
		Expr Expr
		Desc bool
	}
)

const (
	Union     QuerySetOp = "union"
	Intersect QuerySetOp = "intersect"
	Except    QuerySetOp = "except"
)

const (
	InnerJoin      JoinType = "inner"
	LeftOuterJoin  JoinType = "left outer"
	RightOuterJoin JoinType = "right outer"
)

func (q QuerySet) write(w Writer) {
	w.d.WriteQuerySet(w, q)
}

func (s Select) write(w Writer) {
	w.d.WriteSelect(w, s)
}

func (DialectBase) WriteSelect(w Writer, s Select) {
	tbl := code.Table{}

	kw := "SELECT"

	if s.Distinct {
		kw += " DISTINCT"
	}

	for idx, col := range s.Columns {
		if idx > 0 {
			kw = ""
		}

		expr := Render(w.d, col.Expr)

		if col.Alias != "" {
			expr += " AS " + Render(w.d, col.Alias)
		}

		if idx < len(s.Columns) - 1 {
			expr += ","
		}

		tbl.AddRow("", kw, expr)
	}

	if s.From.Source != nil {
		src := Render(w.d, s.From.Source)

		if s.From.Alias != "" {
			src += " AS " + Render(w.d, s.From.Alias)
		}

		tbl.AddRow("", "FROM", src)
	}
}
