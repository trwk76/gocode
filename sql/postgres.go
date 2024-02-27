package sql

import (
	"fmt"

	code "github.com/trwk76/gocode"
)

var Postgres postgresDialect

type postgresDialect struct{}

func (postgresDialect) WriteName(w *code.Writer, name Name) {
	fmt.Fprintf(w, `"%s"`, name)
}

var _ Dialect = Postgres
