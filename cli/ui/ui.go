package ui

import (
	"fmt"

	"github.com/muesli/reflow/indent"
)

type ResultEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ResultTitle(title string) {
	fmt.Println(title)
}

func Result(entries []ResultEntry, indentation uint) {
	str := ""
	maxWidth := 0
	// This probably a better way to do this but this is quick and easy,
	for _, entry := range entries {
		if len(entry.Key) > maxWidth {
			maxWidth = len(entry.Key)
		}
	}
	for _, entry := range entries {
		padding := maxWidth - len(entry.Key) + 1
		str += fmt.Sprintf("%s:%*s%s\n", entry.Key, -padding, "", entry.Value)
	}
	str = indent.String(str, indentation)
	fmt.Println(str)
}
