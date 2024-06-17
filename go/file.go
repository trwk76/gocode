package golang

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	code "github.com/trwk76/gocode"
)

func NewFile(comments []Comment, pkgName string) *File {
	return &File{
		comm: comments,
		pkg:  pkgName,
	}
}

type (
	File struct {
		comm []Comment
		pkg  string

		Imports Imports
	}

	Imports []ImportItem

	ImportItem struct {
		Path  StringVal
		Alias ID
	}
)

func (f *File) Write(w *code.Writer) {
	for _, itm := range f.comm {
		itm.write(w, false)
		w.Newline()
	}

	fmt.Fprintf(w, "package %s", f.pkg)
	w.Newline()
	w.Newline()

	f.Imports.write(w)
}

func (i Imports) write(w *code.Writer) {
	switch len(i) {
	case 0:
		return
	case 1:
		w.WriteString("import")

		if i[0].Alias != "" {
			w.Space()
			i[0].Alias.write(w, true)
		}

		w.Space()
		i[0].Path.write(w, true)
		w.Newline()
		return
	}

	sort.Slice(i, func(li, ri int) bool {
		lsys := isSysImport(i[li].Path)
		rsys := isSysImport(i[ri].Path)

		if lsys != rsys {
			return lsys
		}

		return i[li].Path < i[ri].Path
	})

	hasAlias := slices.ContainsFunc(i, func(itm ImportItem) bool { return itm.Alias != "" })
	prvSys := true
	rows := make([]code.Row, 0)

	for _, itm := range i {
		curSys := isSysImport(itm.Path)

		if !curSys && prvSys {
			rows = append(rows, code.Row{})
		}

		prvSys = curSys

		if hasAlias {
			rows = append(rows, code.Row{Cols: []string{string(itm.Alias), toString(itm.Path)}})
		} else {
			rows = append(rows, code.Row{Cols: []string{toString(itm.Path)}})
		}
	}

	w.WriteString("import (")
	w.Newline()
	w.Indent(func(w *code.Writer) {
		w.Table(rows...)
	})
	w.WriteByte(')')
	w.Newline()
}

func isSysImport(path StringVal) bool {
	return strings.ContainsAny(string(path), ".")
}
