package makefile

import (
	"strconv"
	"strings"
)

//NoPrint formats a command to be executed without being printed
func NoPrint(cmd string) string {
	if strings.Contains(cmd, "\n") {
		spl := strings.Split(cmd, "\n")
		for i, c := range spl {
			spl[i] = NoPrint(c)
		}
		return strings.Join(spl, "\n")
	}
	return "@" + cmd
}

//Print returns a command that prints something when executed
//NOTE: handles multi-line strings
func Print(str string) string {
	if strings.Contains(str, "\n") {
		spl := strings.Split(str, "\n")
		for i, c := range spl {
			spl[i] = Print(c)
		}
		return strings.Join(spl, "\n")
	}
	return NoPrint("echo " + strconv.Quote(str))
}
