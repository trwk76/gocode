package sql

type (
	Predicate interface {
		Item
		pred()
	}
)
