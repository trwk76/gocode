package gocode

import (
	"fmt"
	"strings"
)

type (
	Table struct {
		commpfx string
		colws   []int
		maxw    int
		rows    []tableRow
	}

	tableRow struct {
		comm  string
		cells []string
	}
)

func NewTable(commentPrefix string) Table {
	return Table{
		commpfx: commentPrefix,
		colws:   nil,
		maxw:    0,
		rows:    nil,
	}
}

func (t *Table) AddRow(comment string, cells ...any) {
	rcells := make([]string, len(cells))

	for idx, cell := range cells {
		var text string

		switch tcell := cell.(type) {
		case Writable:
			text = String(tcell)
		case fmt.Stringer:
			text = tcell.String()
		case string:
			text = tcell
		default:
			text = fmt.Sprintf("%v", tcell)
		}

		l := len(text)

		if idx < len(t.colws) {
			if l > t.colws[idx] {
				t.colws[idx] = l
			}
		} else {
			t.colws = append(t.colws, l)
		}

		if l > t.maxw {
			t.maxw = l
		}
	}

	t.rows = append(t.rows, tableRow{
		comm:  comment,
		cells: rcells,
	})
}

func (t *Table) Write(w *Writer) {
	pad := strings.Repeat(" ", t.maxw)

	for _, row := range t.rows {
		if row.comm != "" {
			// Render comment
			for _, line := range strings.Split(row.comm, "\n") {
				w.WriteString(t.commpfx)
				w.WriteString(line)
				w.NewLine()
			}
		}

		for idx, cell := range row.cells {
			if idx > 0 {
				w.Space()
			}

			w.WriteString(cell)
			w.WriteString(pad[:t.colws[idx]-len(cell)])
		}

		w.NewLine()
	}
}
