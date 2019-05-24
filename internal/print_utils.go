package internal

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

// printKind is a kind of information that we need to print.
type printKind int

const (
	// Available printKinds listed below.
	kindErr printKind = iota
	kindWarn
	kindSuccess
	kindInfo
)

// Printf formats given string with parameters and after calls Print with a result of formatting.
func Printf(kind printKind, format string, args ...interface{}) {
	Print(kind, fmt.Sprintf(format, args...))
}

// Print prints some value to stdout using coloring driven by kind.
func Print(kind printKind, v interface{}) {
	switch kind {
	case kindWarn:
		fmt.Println(Yellow(v))
	case kindErr:
		fmt.Println(Red(v))
	case kindSuccess:
		fmt.Println(Green(v))
	default:
		fmt.Println(v)
	}

	return
}
