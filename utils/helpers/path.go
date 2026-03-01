package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patppuccin/kredenv/consts"
)

func GetRootDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, consts.RootDirName), nil
}
