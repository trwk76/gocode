package sql

import (
	"fmt"
)

var SQLServer sqlserverDialect

type (
	sqlserverDialect struct{
		DialectBase
	}
)

func (sqlserverDialect) WriteName(w Writer, n Name) {
	fmt.Fprintf(w, "[%s]", n)
}

var _ Dialect = SQLServer
