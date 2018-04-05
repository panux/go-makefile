package makefile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type makeWriter struct {
	w *bufio.Writer
}

func (mw makeWriter) writeRaw(raw RawText) error {
	_, err := mw.w.WriteString(string(raw))
	return err
}

func (mw makeWriter) writeMakeVarAssign(mva makeVarAssign) (err error) {
	defer func() { //catch errors
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()
	return mw.writeRaw(RawText("\n" + mva.String()))
}

func (mw makeWriter) writeComment(c *Comment) error {
	c.afix()
	for _, l := range *c {
		_, err := fmt.Fprintf(mw.w, "\n# %s", l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mw makeWriter) writeRule(r *Rule) (err error) {
	defer func() { //catch errors
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()
	if r.Attr.OneShell {
		_, err = mw.w.WriteString("\n.ONESHELL:")
		if err != nil {
			return
		}
	}
	_, err = fmt.Fprintf(mw.w, "\n%s:", r.Name.Convert())
	if err != nil {
		return
	}
	if r.Attr.Shell != nil {
		_, err = fmt.Fprintf(mw.w, " SHELL:=%s", r.Attr.Shell.Convert())
		if err != nil {
			return
		}
	}
	if len(r.Deps) != 0 {
		for _, v := range r.Deps {
			_, err = fmt.Fprintf(mw.w, " %s", v.Convert())
			if err != nil {
				return
			}
		}
	}
	if len(r.Command) != 0 {
		for _, v := range r.Command {
			_, err = fmt.Fprintf(mw.w, "\n\t%s", v.String())
			if err != nil {
				return
			}
		}
	}
	return
}

func (mw makeWriter) writeSomething(i interface{}) error {
	switch v := i.(type) {
	case *Rule:
		return mw.writeRule(v)
	case RawText:
		return mw.writeRaw(v)
	case makeVarAssign:
		return mw.writeMakeVarAssign(v)
	case *Comment:
		return mw.writeComment(v)
	default:
		return errors.New("Unsupported type write")
	}
}

type countWriter struct {
	w io.Writer
	n int64
}

func (cw *countWriter) Write(dat []byte) (int, error) {
	n, err := cw.w.Write(dat)
	cw.n += int64(n)
	return n, err
}
