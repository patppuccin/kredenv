package keyring

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	service      = "kredenv"
	indexPrefix  = "__kredenv-idx"
	maxPageBytes = 1024
)

func addToIndex(key string) error {
	// loop through pages, find one with space, add key to it
	pageIdx := 1
	for {
		pageKey := fmt.Sprintf("%s%d", indexPrefix, pageIdx)
		value, err := keyring.Get(service, pageKey)
		if err != nil {
			if errors.Is(err, keyring.ErrNotFound) { // likely no page yet, start fresh
				data, _ := json.Marshal([]string{key})
				return keyring.Set(service, pageKey, string(data))
			}
			return fmt.Errorf("could not fetch index page %d: %w", pageIdx, err)
		}

		var keys []string
		if err := json.Unmarshal([]byte(value), &keys); err != nil {
			return fmt.Errorf("could not parse index page %d: %w", pageIdx, err)
		}

		// check if key already exists
		if slices.Contains(keys, key) {
			return nil
		}

		// check if there's room on this page
		candidate, _ := json.Marshal(append(keys, key))
		if len(candidate) <= maxPageBytes {
			return keyring.Set(service, pageKey, string(candidate))
		}

		// page full, try next
		pageIdx++
	}
}

func removeFromIndex(key string) error {
	pageIdx := 1
	for {
		pageKey := fmt.Sprintf("%s%d", indexPrefix, pageIdx)
		value, err := keyring.Get(service, pageKey)
		if err != nil {
			if errors.Is(err, keyring.ErrNotFound) {
				break // key not found in any page
			}
			return fmt.Errorf("could not fetch index page %d: %w", pageIdx, err)
		}

		var keys []string
		if err := json.Unmarshal([]byte(value), &keys); err != nil {
			return fmt.Errorf("could not parse index page %d: %w", pageIdx, err)
		}

		filtered := keys[:0]
		found := false
		for _, k := range keys {
			if k == key {
				found = true
				continue
			}
			filtered = append(filtered, k)
		}

		if found {
			data, _ := json.Marshal(filtered)
			return keyring.Set(service, pageKey, string(data))
		}

		pageIdx++
	}
	return nil
}

func List() ([]string, error) {
	var allKeys []string
	pageIdx := 1
	for {
		value, err := keyring.Get(service, fmt.Sprintf("%s%d", indexPrefix, pageIdx))
		if err != nil {
			if errors.Is(err, keyring.ErrNotFound) {
				break // no more pages
			}
			return nil, fmt.Errorf("could not fetch index page %d: %w", pageIdx, err)
		}

		var keys []string
		if err := json.Unmarshal([]byte(value), &keys); err != nil {
			return nil, fmt.Errorf("could not parse index page %d: %w", pageIdx, err)
		}

		allKeys = append(allKeys, keys...)
		pageIdx++
	}
	return allKeys, nil
}

func Set(key, value string) error {
	if strings.HasPrefix(key, indexPrefix) {
		return fmt.Errorf("key %q is reserved for internal use", key)
	}
	if err := keyring.Set(service, key, value); err != nil {
		return fmt.Errorf("could not store key %q: %w", key, err)
	}
	return addToIndex(key)
}

func Get(key string) (string, error) {
	if strings.HasPrefix(key, indexPrefix) {
		return "", fmt.Errorf("key %q is reserved for internal use", key)
	}
	value, err := keyring.Get(service, key)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", fmt.Errorf("key %q not found", key)
		}
		return "", fmt.Errorf("could not retrieve key %q: %w", key, err)
	}
	return value, nil
}

func Delete(key string) error {
	if strings.HasPrefix(key, indexPrefix) {
		return fmt.Errorf("key %q is reserved for internal use", key)
	}
	if err := keyring.Delete(service, key); err != nil {
		if err == keyring.ErrNotFound {
			return fmt.Errorf("key %q not found", key)
		}
		return fmt.Errorf("could not delete key %q: %w", key, err)
	}
	return removeFromIndex(key)
}

func Exists(key string) bool {
	_, err := keyring.Get(service, key)
	return err == nil
}
