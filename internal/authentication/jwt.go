package authentication

import (
	"errors"
	"sync"
	"time"

	"github.com/Rhaqim/creds/internal/config"
	"github.com/dgrijalva/jwt-go"
)

// TokenBlacklist is a map that stores blacklisted JWT tokens.
var TokenBlacklist map[string]bool = make(map[string]bool)
var TokenMutex = &sync.Mutex{}

type JWT struct {
	secret        string
	refreshSecret string
}

// NewJWT creates a new instance of JWT with the provided access and refresh secrets.
// It also initializes a MongoDB connection and returns the JWT instance.
func NewJWT() *JWT {
	secret := config.AccessSecret
	refreshSecret := config.RefreshSecret

	return &JWT{
		secret:        secret,
		refreshSecret: refreshSecret,
	}
}

// GenerateToken generates an access token and a refresh token for the given email, username, and userID.
// It returns the access token, refresh token, and any error encountered during the token generation process.
func (j *JWT) GenerateToken(email, username string, openID string) (string, string, error) {

	expirationTime := config.AccessTokenExpireTime
	refreshExpirationTime := config.RefreshTokenExpireTime

	claims := config.JWTClaims{
		Email:    email,
		Username: username,
		OpenID:   openID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	refreshClaims := config.JWTClaims{
		Email:    email,
		Username: username,
		OpenID:   openID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}

	refreshT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshT.SignedString([]byte(j.refreshSecret))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// ValidateToken validates the access token and returns the JWT claims if the token is valid.
// It takes a token string as input and returns the JWT claims and an error if any.
func (j *JWT) ValidateToken(token string) (claim *config.JWTClaims, err error) {
	claim = &config.JWTClaims{}

	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid signature")
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	if claim.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	return claim, nil
}

// ValidateRefreshToken validates the refresh token and returns the JWT claims if the token is valid.
// It takes a token string as input and returns the JWT claims and an error if any.
func (j *JWT) ValidateRefreshToken(token string) (claim *config.JWTClaims, err error) {
	claim = &config.JWTClaims{}

	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.refreshSecret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid signature")
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	if claim.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	return claim, nil
}

// DeleteToken adds the token to the blacklist
func (j *JWT) DeleteToken(token string) error {
	TokenMutex.Lock()
	defer TokenMutex.Unlock()

	TokenBlacklist[token] = true

	return nil
}

// IsTokenBlacklisted checks if the token is blacklisted
func (j *JWT) IsTokenBlacklisted(token string) bool {
	TokenMutex.Lock()
	defer TokenMutex.Unlock()

	return TokenBlacklist[token]
}

// ClearTokenBlacklist clears the token blacklist
func (j *JWT) ClearTokenBlacklist() {
	TokenMutex.Lock()
	defer TokenMutex.Unlock()

	TokenBlacklist = make(map[string]bool)
}
