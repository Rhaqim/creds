package tests

import (
	"testing"

	"github.com/Rhaqim/creds/internal/lib"
)

func TestEncryptionService_Decrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	plaintext := []byte("Hello, World!")

	encrypter := lib.EncryptionService{}
	ciphertext, err := encrypter.Encrypt(plaintext, key)
	if err != nil {
		t.Errorf("encryption error: %v", err)
	}

	decrypted, err := encrypter.Decrypt(ciphertext, key)
	if err != nil {
		t.Errorf("decryption error: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted text does not match plaintext")
	}
}

func TestEncryptionService_Scramble(t *testing.T) {
	encrypter := lib.EncryptionService{}

	sensitiveData := "Hello, World!"

	encodedData, err := encrypter.Scramble(sensitiveData)
	if err != nil {
		t.Errorf("scramble error: %v", err)
	}

	// Decode encrypted data from string
	encryptedData, err := encrypter.DecodeString(encodedData)
	if err != nil {
		t.Errorf("decode error: %v", err)
	}

	// Decrypt data
	decryptedData, err := encrypter.Decrypt(encryptedData, encrypter.EncryptionKey)
	if err != nil {
		t.Errorf("decryption error: %v", err)
	}

	if string(decryptedData) != sensitiveData {
		t.Errorf("decrypted data does not match sensitive data")
	}
}
