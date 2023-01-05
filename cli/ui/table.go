package ui

import (
	"fmt"
	"strings"
)

type Column struct {
	Title string
	Width int
}

type Row []string

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func Table(title string, cols []Column, rows []Row) {
	totalWidth := 0
	header := ""
	body := ""

	for _, col := range cols {
		totalWidth += col.Width
		header += fmt.Sprintf("|%s", center(col.Title, col.Width))
	}
	header += "|\n"
	title = center(title, totalWidth)
	separator := strings.Repeat("-", totalWidth+len(cols)+1) + "\n"

	for _, row := range rows {
		for idx, col := range cols {
			// Truncate the row to the column width
			if len(row[idx]) > col.Width {
				row[idx] = row[idx][:col.Width]
			}
			body += fmt.Sprintf("%s", center(row[idx], col.Width+1))
		}
		body += "\n"
	}

	fmt.Printf("%s\n%s%s%s%s", title, separator, header, separator, body)
}
