package golang

import (
	"fmt"

	code "github.com/trwk76/gocode"
)

func NewFile(pkgName Identifier) *File {
	if pkgName.Exported() {
		panic(fmt.Errorf("'%s' is not a valid package name", pkgName))
	}

	return &File{
		pkgName: pkgName,
	}
}

type (
	File struct {
		comms   []Comment
		pkgName Identifier
		imps    Imports
		decls   []decls
	}
)

func (f *File) CommentPrefix(comment Comment) {
	if len(comment) > 0 {
		f.comms = append(f.comms, comment)
	}
}

func (f *File) Imports() *Imports {
	return &f.imps
}

func (f *File) Write(w *code.Writer) {
	for _, comm := range f.comms {
		comm.write(w)
		w.Newline()
	}

	fmt.Fprintf(w, "package %s\n\n", f.pkgName)

	for _, decl := range f.decls {
		decl.write(w)
		w.Newline()
	}
}
