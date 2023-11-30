package gocode

import (
	"bufio"
	"io"
	"strings"
)

type (
	Writer struct {
		w   *bufio.Writer
		ts  string
		nl  bool
		ind int
	}

	Writable interface {
		Write(w *Writer)
	}
)

func NewWriter(w io.Writer, tabString string) Writer {
	if tabString == "" {
		tabString = "\t"
	}

	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}

	return Writer{
		w:   bw,
		ts:  tabString,
		nl:  true,
		ind: 0,
	}
}

func (w *Writer) Close() error {
	return w.w.Flush()
}

func (w *Writer) WriteByte(c byte) error {
	if c == '\n' {
		w.nl = true
	} else {
		if err := w.ensureIndented(); err != nil {
			return err
		}
	}

	return w.w.WriteByte(c)
}

func (w *Writer) WriteRune(c rune) (int, error) {
	if c == '\n' {
		w.nl = true
	} else {
		if err := w.ensureIndented(); err != nil {
			return 0, err
		}
	}

	return w.w.WriteRune(c)
}

func (w *Writer) WriteString(str string) (int, error) {
	res := 0

	for idx, line := range strings.Split(str, "\n") {
		if idx > 0 {
			if err := w.NewLine(); err != nil {
				return res, err
			}
		}

		if err := w.ensureIndented(); err != nil {
			return res, err
		}

		if cnt, err := w.w.WriteString(line); err != nil {
			return res, err
		} else {
			res += cnt
		}
	}

	return res, nil
}

func (w *Writer) Space() error {
	return w.WriteByte(' ')
}

func (w *Writer) NewLine() error {
	return w.WriteByte('\n')
}

func String(w Writable) string {
	buf := strings.Builder{}
	wr := NewWriter(&buf, "")

	w.Write(&wr)
	wr.Close()

	return buf.String()
}

func (w *Writer) ensureIndented() error {
	if !w.nl {
		return nil
	}

	for i := 0; i < w.ind; i += 1 {
		if _, err := w.w.WriteString(w.ts); err != nil {
			return err
		}
	}

	w.nl = false
	return nil
}
