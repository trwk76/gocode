package sql

type (
	CTEs []CTE

	CTE struct {
		Name      Name
		Recursive bool
		Columns   Names
		Query     Query
	}
)
