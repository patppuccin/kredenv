package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	saltSize  = 16
	keySize   = 32
	timeParam = 1
	memory    = 64 * 1024
	threads   = 4
)

// Encrypt encrypts plaintext with a password using AES-256-GCM
// Output as base64(salt + nonce + ciphertext)
func Encrypt(data []byte, password string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("could not generate salt: %w", err)
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("could not create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("could not create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("could not generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	payload := append(salt, ciphertext...)

	return base64.StdEncoding.EncodeToString(payload), nil
}

// Decrypt decrypts a base64 encoded payload with a password
func Decrypt(encoded, password string) ([]byte, error) {
	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("could not decode payload: %w", err)
	}

	if len(payload) < saltSize {
		return nil, fmt.Errorf("invalid payload")
	}

	salt := payload[:saltSize]
	ciphertext := payload[saltSize:]
	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt — wrong password?")
	}

	return plaintext, nil
}

// deriveKey derives a 32 byte key from a password and salt using argon2id
func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, timeParam, memory, threads, keySize)
}
