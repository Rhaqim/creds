package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"gorm.io/gorm"
)

// Simulate user account with encryption key association
type UserCredentialEncrypt struct {
	gorm.Model
	UserID        uint
	EncryptionKey []byte
}

func Check() {
	// Simulated user account with encryption key
	user := UserCredentialEncrypt{
		UserID: 1,
	}

	user.EncryptionKey = user.generateEncryptionKey()

	// Data to be encrypted
	sensitiveData := "Sensitive information"

	// Encrypt data
	encryptedData, err := user.encrypt([]byte(sensitiveData), user.EncryptionKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	// Encode encrypted data to string
	encodedData := user.encodeToString(encryptedData)

	// Store encrypted data in a secure location (e.g., database)
	user.saveEncodedData(encodedData)

	// Retrieve encoded data from secure location
	decodedData, err := user.decodeString(encodedData)
	if err != nil {
		return
	}

	// Decrypt data
	decryptedData, err := user.decrypt(decodedData, user.EncryptionKey)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	fmt.Println("Decrypted data:", string(decryptedData))
}

// generateEncryptionKey generates a random encryption key
func (O *UserCredentialEncrypt) generateEncryptionKey() []byte {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}
	return key
}

// encrypt encrypts plaintext using the given key
func (O *UserCredentialEncrypt) encrypt(plaintext []byte, key []byte) ([]byte, error) {
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
func (O *UserCredentialEncrypt) decrypt(ciphertext []byte, key []byte) ([]byte, error) {
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

func (O *UserCredentialEncrypt) encodeToString(data []byte) string {
	return hex.EncodeToString(data)
}

func (O *UserCredentialEncrypt) decodeString(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

func (O *UserCredentialEncrypt) saveEncodedData(data string) {
	// Save encoded data
}
