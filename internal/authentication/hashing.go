package authentication

import "golang.org/x/crypto/bcrypt"

type Hash struct {
}

func NewHash() *Hash {
	return &Hash{}
}

func (h *Hash) HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func (h *Hash) ComparePassword(hashedPassword, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func (h *Hash) CheckPasswordStrength(password string) bool {

	return len(password) >= 8
}
