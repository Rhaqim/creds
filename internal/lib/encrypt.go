package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Simulate user account with encryption key association
type EncryptionService struct {
	EncryptionKey []byte
}

func (O *EncryptionService) Scramble(sensitiveData string) (string, error) {

	O.EncryptionKey = O.GenerateEncryptionKey()

	// Encrypt data
	encryptedData, err := O.Encrypt([]byte(sensitiveData), O.EncryptionKey)
	if err != nil {
		return "", err
	}

	// Encode encrypted data to string
	encodedData := O.EncodeToString(encryptedData)

	return encodedData, nil
}

func (O *EncryptionService) Unscramble(encodedData string) (string, error) {
	// Retrieve encoded data from secure location
	decodedData, err := O.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	// Decrypt data
	decryptedData, err := O.Decrypt(decodedData, O.EncryptionKey)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

// generateEncryptionKey generates a random encryption key
func (O *EncryptionService) GenerateEncryptionKey() []byte {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}
	return key
}

// encrypt encrypts plaintext using the given key
func (O *EncryptionService) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// decrypt decrypts ciphertext using the given key
func (O *EncryptionService) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func (O *EncryptionService) EncodeToString(data []byte) string {
	return hex.EncodeToString(data)
}

func (O *EncryptionService) DecodeString(data string) ([]byte, error) {
	return hex.DecodeString(data)
}
