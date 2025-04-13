package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)
import "github.com/pkg/errors"

type Config struct {
	Version         string
	ConfigPath      string
	Host            string `koanf:"Host"`
	Port            int    `koanf:"Port"`
	BaseUrl         string `koanf:"BaseUrl"`
	SessionSecret   string `koanf:"SessionSecret"`
	EncryptionKey   string `koanf:"EncryptionKey"`
	LogLevel        string `koanf:"LogLevel"`
	LogPath         string `koanf:"LogPath"`
	LogMaxSize      int    `koanf:"LogMaxSize"`
	LogMaxBackups   int    `koanf:"LogMaxBackups"`
	CheckForUpdates bool   `koanf:"CheckForUpdates"`
}

type ConfigUpdate struct {
	Host            *string `json:"host,omitempty"`
	Port            *int    `json:"port,omitempty"`
	LogLevel        *string `json:"log_level,omitempty"`
	LogPath         *string `json:"log_path,omitempty"`
	BaseURL         *string `json:"base_url,omitempty"`
	CheckForUpdates *bool   `json:"check_for_updates,omitempty"`
}

func (c *Config) getEncryptionKey() ([]byte, error) {
	key, err := hex.DecodeString(c.EncryptionKey)
	if err != nil {
		return nil, errors.New("invalid hex encryption key")
	}
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}
	return key, nil
}

func (c *Config) Decrypt(ciphertext, iv []byte) ([]byte, error) {
	key, err := c.getEncryptionKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *Config) Encrypt(plaintext, iv []byte) ([]byte, error) {
	key, err := c.getEncryptionKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, iv, plaintext, nil)
	return ciphertext, nil
}
