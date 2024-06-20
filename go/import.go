package golang

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	code "github.com/trwk76/gocode"
)

func (i Imports) PkgPath(alias Identifier) string {
	for _, itm := range i {
		if itm.alias == alias {
			return itm.path
		}
	}

	return ""
}

func (i *Imports) Ensure(alias Identifier, pkgPath string, comment string) Pkg {
	if !isImportAlias(alias) {
		panic(fmt.Errorf("'%s' is not a valid package alias", alias))
	}

	alias = i.uniqueAlias(alias)

	res := Pkg{
		comm:  comment,
		alias: alias,
		path:  pkgPath,
	}

	*i = append(*i, res)
	return res
}

func (p Pkg) NamedType(id Identifier, args ...Type) NamedType {
	if p.alias == Discard {
		panic(fmt.Errorf("package '%s' is implicitely imported using _", p.path))
	} else if !id.Exported() {
		panic(fmt.Errorf("cannot reference a non-exported type from package '%s'", p.path))
	}

	return NamedType{pkg: p.alias, id: id, args: args}
}

func (p Pkg) Symbol(id Identifier, args ...Type) SymbolExpr {
	if p.alias == Discard {
		panic(fmt.Errorf("package '%s' is implicitely imported using _", p.path))
	} else if !id.Exported() {
		panic(fmt.Errorf("cannot reference a non-exported symbol from package '%s'", p.path))
	}

	return SymbolExpr{pkg: p.alias, id: id, args: args}
}

type (
	Imports []Pkg

	Pkg struct {
		comm  string
		alias Identifier
		path  string
	}
)

func (i Imports) uniqueAlias(base Identifier) Identifier {
	if base == Discard {
		return base
	}

	if i.PkgPath(base) == "" {
		return base
	}

	idx := 1
	alias := Identifier(fmt.Sprintf("%s%d", base, idx))

	for i.PkgPath(alias) != "" {
		idx++
		alias = Identifier(fmt.Sprintf("%s%d", base, idx))
	}

	return alias
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
		i[0].alias.write(w)
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

func (i Pkg) row() code.Row {
	var res code.Row

	res.Prefix = commentRenderer(i.comm)
	res.Cols = []string{string(i.alias), strconv.Quote(i.path)}

	return res
}

func isSysPkg(path string) bool {
	return !strings.ContainsAny(path, ".")
}

func isImportAlias(alias Identifier) bool {
	return !alias.Exported()
}
