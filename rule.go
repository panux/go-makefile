package makefile

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

//RuleWriter is a writer for a Makefile rule
type RuleWriter struct {
	wd bool //are we done writing deps?
	cl bool //is this closed?
	w  *Writer
}

//ErrDepsWritten indicates that the user attempted to write another dep after a command had been written
var ErrDepsWritten = errors.New("deps already written")

//AddDep adds a dependency to a rule
func (rw *RuleWriter) AddDep(dep string) error {
	if rw.wd {
		return ErrDepsWritten
	}
	_, err := fmt.Fprintf(rw.w, " %s", dep)
	if err != nil {
		return err
	}
	return nil
}

//AddCommand writes a command to the RuleWriter
func (rw *RuleWriter) AddCommand(cmd string) error {
	rw.wd = true
	if rw.cl {
		return io.ErrClosedPipe
	}
	if strings.Contains(cmd, "\n") {
		spl := strings.Split(cmd, "\n")
		for _, c := range spl {
			err := rw.AddCommand(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
	_, err := fmt.Fprintf(rw.w, "\n\t%s", cmd)
	if err != nil {
		return err
	}
	return nil
}
