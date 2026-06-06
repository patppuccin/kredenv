package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/patppuccin/kredenv/src/consts"
	"github.com/patppuccin/kredenv/src/helpers"
)

type Store struct {
	password string
	data     map[string]string
	changed  bool
}

func Open(password string) (*Store, error) {
	s := &Store{password: password, data: map[string]string{}}

	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return nil, err
	}

	encFilePath := filepath.Join(rootDir, consts.EncFileName)

	content, err := os.ReadFile(encFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, fmt.Errorf("could not read store: %w", err)
	}

	// Empty file likely means no secrets (valid as well)
	if len(content) == 0 {
		return s, nil
	}

	plaintext, err := helpers.Decrypt(string(content), password)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt store: %w", err)
	}

	if err := json.Unmarshal(plaintext, &s.data); err != nil {
		return nil, fmt.Errorf("could not parse secrets store: %w", err)
	}

	return s, nil
}

func (s *Store) Get(key string) (string, error) {
	val, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

func (s *Store) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	s.data[key] = value
	s.changed = true
	return nil
}

func (s *Store) List() (map[string]string, error) {
	result := make(map[string]string, len(s.data))
	maps.Copy(result, s.data)
	return result, nil
}

func (s *Store) Delete(key string) error {
	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	delete(s.data, key)
	s.changed = true
	return nil
}

func (s *Store) Close() error {
	if !s.changed {
		return nil
	}

	plaintext, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("could not serialize store: %w", err)
	}

	encrypted, err := helpers.Encrypt(plaintext, s.password)
	if err != nil {
		return fmt.Errorf("could not encrypt store: %w", err)
	}

	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return err
	}

	encFilePath := filepath.Join(rootDir, consts.EncFileName)

	// Atomic write to prevent partial writes
	tmp := encFilePath + ".tmp"
	if err := os.WriteFile(tmp, []byte(encrypted), 0600); err != nil {
		return fmt.Errorf("could not write store: %w", err)
	}

	if err := os.Rename(tmp, encFilePath); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("could not finalize store: %w", err)
	}

	return nil
}

func Exists() bool {
	rootDir, err := helpers.GetRootDir()
	if err != nil {
		return false
	}

	encFilePath := filepath.Join(rootDir, consts.EncFileName)

	info, err := os.Stat(encFilePath)
	return err == nil && !info.IsDir()
}

func Migrate(oldPassword, newPassword string) error {
	s, err := Open(oldPassword)
	if err != nil {
		return fmt.Errorf("could not open store: %w", err)
	}
	s.password = newPassword
	s.changed = true
	return s.Close()
}
