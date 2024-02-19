package code

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	w  *bufio.Writer
	ts string
	in int
	nl bool
}

func WriteFile(path string, tabString string, action func(w *Writer)) error {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0777)); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	w := NewWriter(f, tabString)

	action(&w)

	w.Flush()
	return nil
}

func WriteString(tabString string, action func(w *Writer)) string {
	buf := strings.Builder{}
	w := NewWriter(&buf, tabString)

	action(&w)

	w.Flush()
	return buf.String()
}

func NewWriter(w io.Writer, tabString string) Writer {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}

	if tabString == "" {
		tabString = "\t"
	}

	return Writer{
		w:  bw,
		ts: tabString,
		in: 0,
		nl: true,
	}
}

func (w *Writer) WriteByte(b byte) error {
	if b == '\n' {
		w.nl = true
	} else {
		w.ensureIndented()
	}

	return w.w.WriteByte(b)
}

func (w *Writer) WriteString(str string) (int, error) {
	for idx, line := range strings.Split(str, "\n") {
		if idx > 0 {
			w.WriteByte('\n')
		}

		w.ensureIndented()
		w.w.WriteString(line)
	}

	return len(str), nil
}

func (w *Writer) Write(data []byte) (int, error) {
	res := 0

	for len(data) > 0 {
		if err := w.WriteByte(data[0]); err != nil {
			return res, err
		}

		data = data[1:]
		res++
	}

	return res, nil
}

func (w *Writer) Space() {
	w.WriteByte(' ')
}

func (w *Writer) Newline() {
	w.WriteByte('\n')
}

func (w *Writer) Indent(action func()) {
	w.in++
	action()
	w.in--
}

func (w *Writer) Flush() error {
	return w.w.Flush()
}

func (w *Writer) ensureIndented() {
	if !w.nl {
		return
	}

	for i := 0; i < w.in; i++ {
		w.w.WriteString(w.ts)
	}

	w.nl = false
}
