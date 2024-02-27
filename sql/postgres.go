package sql

import (
	"fmt"
)

var Postgres postgresDialect

type (
	postgresDialect struct{
		DialectBase
	}
)

func (postgresDialect) WriteName(w Writer, name Name) {
	fmt.Fprintf(w, `"%s"`, name)
}

var _ Dialect = Postgres
