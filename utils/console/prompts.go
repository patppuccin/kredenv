package console

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"golang.org/x/term"
)

var fmtPrompt = func(content string) string { return color.New(color.FgYellow).Sprint(content) }

func PromptSecret(prompt string) (string, error) {
	fmt.Printf("%s %s: ", fmtPrompt("(*)"), prompt)
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

func PromptConfirm(prompt string) (bool, error) {
	fmt.Printf("%s %s [y/N]: ", fmtPrompt("(?)"), prompt)
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return false, err
	}
	return response == "y" || response == "Y", nil
}

func PromptAndConfirmPassword(prompt, confirmPrompt string) (string, error) {
	pwd, err := PromptSecret(prompt)
	if err != nil {
		return "", err
	}

	confirm, err := PromptSecret(confirmPrompt)
	if err != nil {
		return "", err
	}

	if pwd != confirm {
		return "", fmt.Errorf("passwords do not match")
	}

	return pwd, nil
}
