package ui

import (
	"fmt"

	"github.com/muesli/reflow/wordwrap"
)

func Errorf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Println(wordwrap.String(s, MaxOutputLineLength))
}

func Infof(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Println(wordwrap.String(s, MaxOutputLineLength))
}

func Successf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Println(wordwrap.String(s, MaxOutputLineLength))
}
