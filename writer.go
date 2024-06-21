package code

import (
	"bufio"
	"io"
	"strings"
)

func NewWriter(w io.Writer, tab string) *Writer {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}

	if tab == "" {
		tab = "\t"
	}

	return &Writer{
		w:   bw,
		tab: tab,
		ind: 0,
		nl:  true,
	}
}

func String(f func(w *Writer)) string {
	buf := strings.Builder{}
	w := NewWriter(&buf, "")

	if f != nil {
		f(w)
	}

	w.Close()
	return buf.String()
}

type (
	Writer struct {
		w   *bufio.Writer
		tab string
		ind uint16
		nl  bool
	}

	Row struct {
		Prefix Renderer
		Cols   []string
	}

	Renderer func(w *Writer)
)

func (w *Writer) Close() {
	w.w.Flush()
}

func (w *Writer) WriteByte(b byte) error {
	if b == '\n' {
		w.nl = true
	} else {
		w.ensureIndented()
	}

	return w.w.WriteByte(b)
}

func (w *Writer) WriteRune(r rune) (int, error) {
	if r == '\n' {
		w.nl = true
	} else {
		w.ensureIndented()
	}

	return w.w.WriteRune(r)
}

func (w *Writer) WriteString(s string) (int, error) {
	res := 0

	for idx, line := range strings.Split(s, "\n") {
		if idx > 0 {
			w.WriteByte('\n')
			res++
		}

		w.ensureIndented()
		cnt, err := w.w.WriteString(line)
		if err != nil {
			return res, err
		}

		res += cnt
	}

	return res, nil
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.WriteString(string(p))
}

func (w *Writer) Space() {
	w.WriteByte(' ')
}

func (w *Writer) Newline() {
	w.WriteByte('\n')
}

func (w *Writer) Indent(f func(w *Writer)) {
	w.ind++
	f(w)
	w.ind--
}

func (w *Writer) Table(rows ...Row) {
	colw := make([]int, 0)
	maxw := 0

	for _, row := range rows {
		for idx, col := range row.Cols {
			l := len(col)

			if idx < len(colw) {
				if l > colw[idx] {
					colw[idx] = l
				}
			} else {
				colw = append(colw, l)
			}

			if l > maxw {
				maxw = l
			}
		}
	}

	pad := strings.Repeat(" ", maxw)

	for _, row := range rows {
		if row.Prefix != nil {
			row.Prefix(w)
		}

		for idx, col := range row.Cols {
			if idx > 0 {
				w.Space()
			}

			w.WriteString(col)

			if idx < len(row.Cols)-1 {
				w.WriteString(pad[:colw[idx]-len(col)])
			}
		}

		w.Newline()
	}
}

func (w *Writer) ensureIndented() {
	if !w.nl {
		return
	}

	for i := uint16(0); i < w.ind; i++ {
		w.w.WriteString(w.tab)
	}

	w.nl = false
}
