package makefile

import "fmt"

//Glob creates a glob substitution using the provided pattern
func Glob(pattern string) string {
	return fmt.Sprintf("$(glob %s)", pattern)
}

//ExtSub returns a substitution that replaces an extension
//NOTE: inext and outext should not have a leading dot (auto-added)
func ExtSub(varname string, inext string, outext string) string {
	return fmt.Sprintf("$(%s:.%s=.%s)", varname, inext, outext)
}

//ShellSub does a shell substitution
func ShellSub(cmd string) string {
	return fmt.Sprintf("$(shell %s)", cmd)
}

//VarSub does a makefile variable substitution
func VarSub(varname string) string {
	return fmt.Sprintf("$(%s)", varname)
}
