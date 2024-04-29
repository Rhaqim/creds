package config

import (
	"time"

	ut "github.com/Rhaqim/creds/internal/utils"
	"github.com/dgrijalva/jwt-go"
)

var Domain = "localhost"

// Authentication
var (
	AccessSecret  = ut.Env("ACCESS_SECRET", "")
	RefreshSecret = ut.Env("REFRESH_SECRET", "")

	CookieTokenExpireTime        = 60 * 60 * 24 * 7     // 1 week
	CookieRefreshTokenExpireTime = 60 * 60 * 24 * 7 * 4 // 4 weeks (1 month)

	AccessTokenExpireTime  = (time.Now().Add(time.Hour * 24 * 7).Unix())   // 1 week
	RefreshTokenExpireTime = time.Now().Add(time.Hour * 24 * 7 * 4).Unix() // 4 weeks (1 month)
)

// JWTClaims
type JWTClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	OpenID   string `json:"open_id"`
	jwt.StandardClaims
}

// OAuth2
var (
	OauthStateString = ut.Env("OAUTH_STATE_STRING", "random")
)

// OAuth2 Google
var (
	GoogleOAuth2ClientID         = ut.Env("GOOGLE_OAUTH2_CLIENT_ID", "")
	GoogleOAuth2ClientSecret     = ut.Env("GOOGLE_OAUTH2_CLIENT_SECRET", "")
	GoogleOAuth2RedirectURL      = ut.Env("GOOGLE_OAUTH2_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
	GoogleOAuth2AuthURL          = ut.Env("GOOGLE_OAUTH2_AUTH_URL", "https://accounts.google.com/o/oauth2/auth")
	GoogleOAuth2TokenURL         = ut.Env("GOOGLE_OAUTH2_TOKEN_URL", "https://accounts.google.com/o/oauth2/token")
	GoogleOAuth2Scope            = ut.Env("GOOGLE_OAUTH2_SCOPE", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile, openid")
	GoogleOAuth2AccessType       = ut.Env("GOOGLE_OAUTH2_ACCESS_TYPE", "offline")
	GoogleOAuth2UserinfoURL      = ut.Env("GOOGLE_USERINFO_URL", "https://www.googleapis.com/oauth2/v3/userinfo")
	GoogleOAuth2OauthStateString = ut.Env("GOOGLE_OAUTH2_OAUTH_STATE_STRING", "random")
)

// OAuth2 Apple
var (
	GithubOAuth2ClientID     = ut.Env("GITHUB_OAUTH2_CLIENT_ID", "")
	GithubOAuth2ClientSecret = ut.Env("GITHUB_OAUTH2_CLIENT_SECRET", "")
	GithubOAuth2RedirectURL  = ut.Env("GITHUB_OAUTH2_REDIRECT_URL", "https://localhost:8080/auth/github/callback")
	GithubOAuth2AuthURL      = ut.Env("GITHUB_OAUTH2_AUTH_URL", "https://www.github.com/auth/authorize")
	GithubOAuth2TokenURL     = ut.Env("GITHUB_OAUTH2_TOKEN_URL", "https://www.github.com/auth/token")
	GithubOAuth2Scope        = ut.Env("GITHUB_OAUTH2_SCOPE", "email name")
	GithubOAuth2ResponseType = ut.Env("GITHUB_OAUTH2_RESPONSE_TYPE", "code")
	GithubOAuth2UserinfoURL  = ut.Env("GITHUB_USERINFO_URL", "https://www.github.com/auth/token")
)

// OAuth2 Apple
var (
	AppleOAuth2ClientID     = ut.Env("APPLE_OAUTH2_CLIENT_ID", "")
	AppleOAuth2ClientSecret = ut.Env("APPLE_OAUTH2_CLIENT_SECRET", "")
	AppleOAuth2RedirectURL  = ut.Env("APPLE_OAUTH2_REDIRECT_URL", "https://localhost:8080/api/v2/auth/apple/callback")
	AppleOAuth2AuthURL      = ut.Env("APPLE_OAUTH2_AUTH_URL", "https://appleid.apple.com/auth/authorize")
	AppleOAuth2TokenURL     = ut.Env("APPLE_OAUTH2_TOKEN_URL", "https://appleid.apple.com/auth/token")
	AppleOAuth2Scope        = ut.Env("APPLE_OAUTH2_SCOPE", "email name")
	AppleOAuth2ResponseType = ut.Env("APPLE_OAUTH2_RESPONSE_TYPE", "code")
	AppleOAuth2UserinfoURL  = ut.Env("APPLE_USERINFO_URL", "https://appleid.apple.com/auth/token")
)
