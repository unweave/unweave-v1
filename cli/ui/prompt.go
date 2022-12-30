package ui

import (
	"strings"

	"github.com/manifoldco/promptui"
)

func Confirm(message string) bool {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil || !strings.EqualFold(result, "y") {
		return false
	}
	return true
}
