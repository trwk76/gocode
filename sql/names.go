package sql

type (
	Name string

	ObjectName struct {
		Schema Name
		Name   Name
	}
)

func (n Name) write(w Writer) {
	w.d.WriteName(w.Writer, n)
}

func (n ObjectName) write(w Writer) {
	if n.Schema != "" {
		n.Schema.write(w)
		w.WriteByte('.')
	}

	n.Name.write(w)
}

var (
	_ Item = Name("")
	_ Item = ObjectName{}
)
