package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Rhaqim/creds/internal/config"
	"github.com/Rhaqim/creds/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2 struct {
	oauthConf   *oauth2.Config
	userInfoURL string
}

func NewOAuth2Adapter(provider string) *OAuth2 {
	// Set up OAuth2 configiguration
	googleOauthConf := &oauth2.Config{
		ClientID:     config.GoogleOAuth2ClientID,
		ClientSecret: config.GoogleOAuth2ClientSecret,
		RedirectURL:  config.GoogleOAuth2RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	var oauthConf *oauth2.Config
	var userInfoURL string

	if provider == "google" {
		oauthConf = googleOauthConf
		userInfoURL = config.GoogleOAuth2UserinfoURL
	} else {
		return nil
	}

	return &OAuth2{
		oauthConf:   oauthConf,
		userInfoURL: userInfoURL,
	}
}

func (OA *OAuth2) OAuth2LoginHandler(c *gin.Context) {

	// Generate the OAuth2 URL
	url := OA.oauthConf.AuthCodeURL(config.GoogleOAuth2OauthStateString, oauth2.AccessTypeOffline)

	// Redirect the user to the OAuth2 provider for authentication
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (OA *OAuth2) OAuth2CallbackHandler(c *gin.Context) {
	state := c.Query("state")
	if state != config.GoogleOAuth2OauthStateString {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid state",
		})
		return
	}

	// Get the authorization code from the query parameters
	code := c.Query("code")

	// Exchange the authorization code for an access token
	token, err := OA.oauthConf.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Failed to exchange the authorization code for an access token",
		})
		return
	}

	// Get the user information
	userInfo, err := OA.getUserInfo(c, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user information",
		})
		return
	}

	access, refresh, err := NewJWT().GenerateToken(userInfo["email"].(string), userInfo["first_name"].(string), userInfo["openid"].(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", access, config.CookieTokenExpireTime, "/", config.Domain, true, true)
	c.SetCookie("refreshToken", refresh, config.CookieRefreshTokenExpireTime, "/", config.Domain, true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user":          userInfo,
	})

	OA.UserExistsHandler(c, userInfo)

}

func (OA *OAuth2) UserExistsHandler(c *gin.Context, userInfo map[string]interface{}) {
	var user models.User

	var email string = userInfo["email"].(string)

	access, refresh, err := NewJWT().GenerateToken(email, userInfo["first_name"].(string), userInfo["openid"].(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	// Get User
	err = user.GetByEmail(email)
	if err != nil {
		if err.Error() == "record not found" {
			// Create new user
			user.Email = email
			user.DisplayName = userInfo["first_name"].(string)
			user.OAuthID = userInfo["openid"].(string)
			user.Role = models.Member

			err := user.Register(user)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "Error creating user",
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Error getting user",
			})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "User exists",
		"access_token":  access,
		"refresh_token": refresh,
		"user":          user,
	})
}

func (OA *OAuth2) getUserInfo(ctx context.Context, token *oauth2.Token) (map[string]interface{}, error) {
	var err error
	var userInfo map[string]interface{}

	// Use the access token to make API requests
	client := OA.oauthConf.Client(ctx, token)
	resp, err := client.Get(OA.userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get the user information
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
