package postgres

import "github.com/trwk76/gocode"

type (
	Name struct {
		Value    string
		Unquoted bool
	}

	ObjectName struct {
		Schema *Name
		Name   Name
	}
)

func (n Name) Write(w *gocode.Writer) {
	if n.Unquoted {
		w.WriteString(n.Value)
	} else {
		w.WriteByte('"')
		w.WriteString(n.Value)
		w.WriteByte('"')
	}
}

func (n ObjectName) Write(w *gocode.Writer) {
	if n.Schema != nil {
		n.Schema.Write(w)
		w.WriteByte('.')
	}

	n.Name.Write(w)
}

var (
	_ gocode.Writable = Name{}
	_ gocode.Writable = ObjectName{}
)
