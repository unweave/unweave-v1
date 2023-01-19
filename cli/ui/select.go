package ui

import (
	"github.com/manifoldco/promptui"
)

func Select(label string, items []string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return idx, nil
}
