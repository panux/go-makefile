package makefile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//Writer is a tool which writes makefiles
type Writer struct {
	w      *bufio.Writer
	rw     *RuleWriter //current RuleWriter
	closed bool
}

//Write allows you to bypass the structured writing and write directly
func (w *Writer) Write(dat []byte) (int, error) {
	if w.closed {
		return 0, io.ErrClosedPipe
	}
	return w.w.Write(dat)
}

//Flush flushes the underlying write buffer
func (w *Writer) Flush() error {
	if w.closed {
		return io.ErrClosedPipe
	}
	return w.w.Flush()
}

func (w *Writer) clrw() {
	if w.rw != nil {
		w.rw.wd = true
		w.rw.cl = true
		w.rw = nil
	}
}

//Rule generates a new rule entry and returns a RuleWriter
func (w *Writer) Rule(name string) (*RuleWriter, error) {
	w.clrw()
	_, err := fmt.Fprintf(w, "\n%s:", name)
	if err != nil {
		return nil, err
	}
	rw := &RuleWriter{w: w}
	w.rw = rw
	return rw, nil
}

//QuickRule is a wrapper around Rule which does everything at once
func (w *Writer) QuickRule(name string, deps []string, cmd []string) error {
	rw, err := w.Rule(name)
	if err != nil {
		return err
	}
	for _, d := range deps {
		err = rw.AddDep(d)
		if err != nil {
			return err
		}
	}
	for _, c := range cmd {
		err = rw.AddCommand(c)
		if err != nil {
			return err
		}
	}
	return nil
}

//Comment writes a comment
//Multiline comments are supported
func (w *Writer) Comment(comment string) error {
	if strings.Contains(comment, "\n") {
		spl := strings.Split(comment, "\n")
		for _, l := range spl {
			err := w.Comment(l)
			if err != nil {
				return err
			}
		}
		return nil
	}
	_, err := fmt.Fprintf(w, "\n# %s", comment)
	return err
}

//Var writes out a variable assignment
func (w *Writer) Var(varname string, val []string) error {
	_, err := fmt.Fprintf(w, "\n%s = %s", varname, strings.Join(val, " "))
	return err
}

//VarAppend writes a variable append
func (w *Writer) VarAppend(varname string, val []string) error {
	_, err := fmt.Fprintf(w, "\n%s += %s", varname, strings.Join(val, " "))
	return err
}

//BlankLine writes out a blank line for spacing purposes
func (w *Writer) BlankLine() error {
	if w.closed {
		return io.ErrClosedPipe
	}
	_, err := w.w.WriteRune('\n')
	return err
}

//Close closes a writer
func (w *Writer) Close() error {
	err := w.Flush()
	if err != nil {
		return err
	}
	w.closed = true
	return nil
}

//NewWriter creates a new Makefile writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}
