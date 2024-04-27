package lib

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// TOTP is a time-based one-time password algorithm.
type TOTP struct {
	// Secret is the shared secret key.
	Secret string
	// Digits is the number of digits in the OTP.
	Digits int
	// Period is the time step in seconds.
	Period int
}

// Generate generates a TOTP.
func (T *TOTP) Generate() (string, error) {
	secret, err := base32.StdEncoding.DecodeString(strings.ToUpper(T.Secret))
	if err != nil {
		return "", err
	}

	epoch := time.Now().Unix() / int64(T.Period)
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(epoch))

	hash := hmac.New(sha1.New, secret)
	hash.Write(msg)
	h := hash.Sum(nil)

	offset := h[len(h)-1] & 0x0F
	code := binary.BigEndian.Uint32(h[offset:offset+4]) & 0x7FFFFFFF
	code %= 1000000

	format := fmt.Sprintf("%%0%dd", T.Digits)
	return fmt.Sprintf(format, code), nil
}

// Verify verifies a TOTP.
func (T *TOTP) Verify(input string) (bool, error) {
	otp, err := T.Generate()
	if err != nil {
		return false, err
	}

	return otp == input, nil
}
