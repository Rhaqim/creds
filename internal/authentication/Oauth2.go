package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Rhaqim/creds/internal/config"
	"github.com/Rhaqim/creds/internal/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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

	githubOauthConf := &oauth2.Config{
		ClientID:     config.GithubOAuth2ClientID,
		ClientSecret: config.GithubOAuth2ClientSecret,
		RedirectURL:  config.GithubOAuth2RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     github.Endpoint,
	}

	var oauthConf *oauth2.Config
	var userInfoURL string

	switch provider {
	case "google":
		oauthConf = googleOauthConf
		userInfoURL = config.GoogleOAuth2UserinfoURL
	case "github":
		oauthConf = githubOauthConf
		userInfoURL = config.GithubOAuth2UserinfoURL
	default:
		return nil
	}

	return &OAuth2{
		oauthConf:   oauthConf,
		userInfoURL: userInfoURL,
	}
}

func (OA *OAuth2) OAuth2LoginHandler(c *gin.Context) {

	// Generate the OAuth2 URL
	url := OA.oauthConf.AuthCodeURL(config.GoogleOAuth2OauthStateString)

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

	OA.UserExistsHandler(c, userInfo)

}

func (OA *OAuth2) UserExistsHandler(c *gin.Context, userInfo map[string]interface{}) {
	var user models.User

	var email string = userInfo["email"].(string)
	var name string = userInfo["name"].(string)
	var oauthID string = userInfo["sub"].(string)

	// Get User
	err = user.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new user
			user.Email = email
			user.DisplayName = name
			user.OAuthID = oauthID

			err := user.Insert()
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

	access, refresh, err := NewJWT().GenerateToken(email, name, oauthID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", access, config.CookieTokenExpireTime, "/", config.Domain, true, true)
	c.SetCookie("refreshToken", refresh, config.CookieRefreshTokenExpireTime, "/", config.Domain, true, true)

	c.Redirect(http.StatusFound, config.FrontendURL)
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

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", config.Domain, true, true)
	c.SetCookie("refreshToken", "", -1, "/", config.Domain, true, true)
}
