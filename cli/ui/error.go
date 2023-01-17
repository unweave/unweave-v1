package ui

import (
	"fmt"

	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/unweave/unweave/api/types"
)

type Error struct {
	*types.HTTPError
}

func (e *Error) Short() string {
	str := fmt.Sprintf("%s API error: %s", e.Provider, e.Message)
	return str
}

func (e *Error) Verbose() string {
	header := "API error:\n"
	if e.Provider != "" {
		header = fmt.Sprintf("%s API error:\n", e.Provider)
	}
	body := ""
	if e.Code != 0 {
		body += wordwrap.String(fmt.Sprintf("Code: %d", e.Code), MaxOutputLineLength-IndentWidth)
		body += "\n"
	}
	if e.Message != "" {
		body += wordwrap.String(fmt.Sprintf("Message: %s", e.Message), MaxOutputLineLength-IndentWidth)
		body += "\n"
	}
	if e.Suggestion != "" {
		body += wordwrap.String(fmt.Sprintf("Suggestion: %s", e.Suggestion), MaxOutputLineLength-IndentWidth)
		body += "\n"
	}
	str := header + indent.String(body, IndentWidth)
	return str
}
