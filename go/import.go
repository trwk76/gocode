package golang

import (
	"fmt"
	"go/token"
	"slices"
	"sort"
	"strconv"
	"strings"

	code "github.com/trwk76/gocode"
)

type (
	Imports []Import

	Import struct {
		comm      string
		alias     string
		path      string
		aliasExpl bool
	}
)

func (i *Imports) Ensure(pkgPath string, alias string, comment string) Import {
	expl := false

	if alias != "" {
		expl = true
	} else {
		alias, expl = detectPkgAlias(pkgPath)
	}

	if !isImportAlias(alias) {
		panic(fmt.Errorf("'%s' is not a valid package alias", alias))
	}

	for _, itm := range *i {
		if itm.alias == alias {
			panic(fmt.Errorf("import alias '%s' is already used by package '%s'", alias, itm.path))
		}
	}

	res := Import{
		comm:      comment,
		alias:     alias,
		path:      pkgPath,
		aliasExpl: expl,
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

		if i[0].alias != "" {
			w.WriteString(i[0].alias)
			w.Space()
		}

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
	hasAlias := i.hasAlias()

	for _, itm := range i {
		if prvSys && !isSysPkg(itm.path) {
			// insert empty row as separator between system and non-system imports
			rows = append(rows, code.Row{})
		}

		rows = append(rows, itm.row(hasAlias))
	}

	w.WriteString("import (")
	w.Newline()

	w.Indent(func(w *code.Writer) {
		w.Table(rows...)
	})

	w.WriteByte(')')
	w.Newline()
}

func (i Imports) hasAlias() bool {
	return slices.ContainsFunc(i, func(i Import) bool { return i.aliasExpl })
}

func (i Import) row(hasAlias bool) code.Row {
	var res code.Row

	res.Prefix = commentRenderer(i.comm)

	if hasAlias {
		alias := ""

		if i.aliasExpl {
			alias = i.alias
		}

		res.Cols = []string{alias, strconv.Quote(i.path)}
	} else {
		res.Cols = []string{strconv.Quote(i.path)}
	}

	return res
}

func isSysPkg(path string) bool {
	return !strings.ContainsAny(path, ".")
}

func detectPkgAlias(path string) (string, bool) {
	res := path

	if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
		res = path[idx+1:]
	}

	if isImportAlias(res) {
		return res, false
	}

}

func isImportAlias(alias string) bool {
	return token.IsIdentifier(alias) && !token.IsExported(alias)
}
