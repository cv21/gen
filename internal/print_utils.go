package internal

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

// printKind is a kind of information that we need to print.
type printKind int

const (
	// Available printKinds listed below.
	kindErr printKind = iota
	kindWarn
	kindSuccess
	kindInfo
	kindInfoFaint
)

func Sprint(kind printKind, args ...interface{}) string {
	return fmt.Sprint(colorize(kind, fmt.Sprint(args...)))
}

func Sprintf(kind printKind, format string, args ...interface{}) string {
	return fmt.Sprint(colorize(kind, fmt.Sprintf(format, args...)))
}

// Printf formats given string with parameters and after calls Println with a result of formatting.
func Printf(kind printKind, format string, args ...interface{}) {
	Print(kind, fmt.Sprintf(format, args...))
}

// Println prints some value to stdout using coloring driven by kind.
func Print(kind printKind, v interface{}) {
	fmt.Print(colorize(kind, v))
}

// Println prints some value to stdout using coloring driven by kind.
func Println(kind printKind, v interface{}) {
	fmt.Println(colorize(kind, v))
}

func colorize(kind printKind, v interface{}) string {
	switch kind {
	case kindWarn:
		return fmt.Sprint(aurora.BrightYellow(v))
	case kindErr:
		return fmt.Sprint(aurora.BrightRed(v))
	case kindSuccess:
		return fmt.Sprint(aurora.BrightGreen(v))
	case kindInfoFaint:
		return fmt.Sprint(aurora.Faint(v))
	default:
		return fmt.Sprint(v)
	}
}
