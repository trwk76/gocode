package golang

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	code "github.com/trwk76/gocode"
)

type (
	Imports []Import

	Import struct {
		comm  string
		alias string
		path  string
	}
)

func (i *Imports) Ensure(pkgPath string, alias string, comment string) Import {
	if !isImportAlias(alias) {
		panic(fmt.Errorf("'%s' is not a valid package alias", alias))
	}

	alias = i.uniqueAlias(alias)

	res := Import{
		comm:      comment,
		alias:     alias,
		path:      pkgPath,
	}

	*i = append(*i, res)
	return res
}

func (i Imports) write(w *code.Writer) {
	switch len(i) {
	case 0:
		return
	case 1:
		if pfx := commentRenderer(i[0].comm); pfx != nil {
			pfx(w)
		}

		w.WriteString("import ")
		w.WriteString(i[0].alias)
		w.Space()
		w.WriteString(strconv.Quote(i[0].path))
		w.Newline()
		return
	}

	sort.Slice(i, func(iidx, jidx int) bool {
		is := isSysPkg(i[iidx].path)
		js := isSysPkg(i[jidx].path)

		if is != js {
			return is
		}

		return i[iidx].path < i[jidx].path
	})

	prvSys := false
	rows := make([]code.Row, 0, len(i))

	for _, itm := range i {
		if prvSys && !isSysPkg(itm.path) {
			// insert empty row as separator between system and non-system imports
			rows = append(rows, code.Row{})
		}

		rows = append(rows, itm.row())
	}

	w.WriteString("import (")
	w.Newline()

	w.Indent(func(w *code.Writer) {
		w.Table(rows...)
	})

	w.WriteByte(')')
	w.Newline()
}

func (i Import) row() code.Row {
	var res code.Row

	res.Prefix = commentRenderer(i.comm)
	res.Cols = []string{i.alias, strconv.Quote(i.path)}

	return res
}

func isSysPkg(path string) bool {
	return !strings.ContainsAny(path, ".")
}

func isImportAlias(alias string) bool {
	if alias == "_" {
		return true
	}

	return token.IsIdentifier(alias) && !token.IsExported(alias)
}
