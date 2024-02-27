package sql

type (
	Names []Name

	Name string

	ObjectName struct {
		Schema Name
		Name   Name
	}
)

func (n Names) write(w Writer) {
	for idx, name := range n {
		if idx > 0 {
			w.WriteString(", ")
		}

		name.write(w)
	}
}

func (n Name) write(w Writer) {
	w.d.WriteName(w, n)
}

func (n ObjectName) write(w Writer) {
	if n.Schema != "" {
		n.Schema.write(w)
		w.WriteByte('.')
	}

	n.Name.write(w)
}

func (ObjectName) selectSrc() {}

var (
	_ Item         = Name("")
	_ Item         = ObjectName{}
	_ SelectSource = ObjectName{}
)
