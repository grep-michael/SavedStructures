package savedstructures

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type EncryptedSaveable struct {
	mu       sync.RWMutex
	filepath string
	data     interface{}
	key      []byte // 32 bytes for AES-256
}

func NewEncryptedSaveable(key []byte) *EncryptedSaveable {
	if len(key) != 32 {
		log.Fatal("Encryption key must be exactly 32 bytes for AES-256")
	}
	return &EncryptedSaveable{key: key}
}

func (ep *EncryptedSaveable) InitPersistent(filepath string, data interface{}) {
	ep.filepath = filepath
	ep.data = data
}

func (ep *EncryptedSaveable) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ep.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (ep *EncryptedSaveable) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ep.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (ep *EncryptedSaveable) Save() error {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	jsonData, err := json.MarshalIndent(ep.data, "", "  ")
	if err != nil {
		return err
	}

	encrypted, err := ep.encrypt(jsonData)
	if err != nil {
		return err
	}

	return os.WriteFile(ep.filepath, encrypted, 0644)
}

func (es *EncryptedSaveable) Load() error {
	es.mu.Lock()
	defer es.mu.Unlock()

	file, err := os.ReadFile(es.filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// File doesn't exist, save encrypted defaults
		jsonData, err := json.MarshalIndent(es.data, "", "  ")
		if err != nil {
			return err
		}
		encrypted, err := es.encrypt(jsonData)
		if err != nil {
			return err
		}
		return os.WriteFile(es.filepath, encrypted, 0644)
	}

	decrypted, err := es.decrypt(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(decrypted, es.data)
}
