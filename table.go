package code

import "strings"

type (
	Table struct {
		rows []tableRow
		colw []int
		maxw int
	}

	tableRow struct {
		pfx  string
		cols []string
	}
)

func (t *Table) AddRow(prefix string, columns... string) {
	for idx, col := range columns {
		l := len(col)

		if idx < len(t.colw) {
			if l > t.colw[idx] {
				t.colw[idx] = l
			}
		} else {
			t.colw = append(t.colw, l)
		}

		if l > t.maxw {
			t.maxw = l
		}
	}

	t.rows = append(t.rows, tableRow{
		pfx:  prefix,
		cols: columns,
	})
}

func (t *Table) Write(w *Writer) {
	pad := strings.Repeat(" ", t.maxw)

	for _, row := range t.rows {
		if row.pfx != "" {
			w.WriteString(row.pfx)
		}

		for idx, col := range row.cols {
			w.WriteString(col)

			if idx < len(row.cols) - 1 {
				w.WriteString(pad[:t.colw[idx] - len(col)])
				w.WriteByte(' ')
			}
		}

		w.Newline()
	}
}
