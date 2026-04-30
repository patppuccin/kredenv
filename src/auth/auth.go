package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/helpers"
	"github.com/zalando/go-keyring"
)

func Retrieve() (string, error) {
	if pwd := os.Getenv(consts.AuthEnvVar); pwd != "" {
		return pwd, nil
	}

	if pwd, err := keyring.Get(consts.AppName, consts.KeyringKey); err == nil && pwd != "" {
		return pwd, nil
	}

	if pwd, err := retrieveFromFile(); err == nil && pwd != "" {
		return pwd, nil
	}

	return "", fmt.Errorf("no master password found")
}

func Store(password string) error {
	if err := keyring.Set(consts.AppName, consts.KeyringKey, password); err == nil {
		return nil
	}
	return storeInFile(password)
}

func Delete() error {
	return keyring.Delete(consts.AppName, consts.KeyringKey)
}

func retrieveFromFile() (string, error) {
	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return "", err
	}

	kredmasterFilePath := filepath.Join(rootDir, consts.AuthMasterFile)

	data, err := os.ReadFile(kredmasterFilePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func storeInFile(password string) error {
	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return err
	}

	kredmasterFilePath := filepath.Join(rootDir, consts.AuthMasterFile)

	if err := os.WriteFile(kredmasterFilePath, []byte(password), 0600); err != nil {
		return fmt.Errorf("could not write to %s: %w", filepath.Base(kredmasterFilePath), err)
	}
	return nil
}
