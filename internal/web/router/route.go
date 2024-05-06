package router

import (
	"github.com/Rhaqim/creds/internal/api"
	"github.com/Rhaqim/creds/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() error {
	r := gin.Default()

	// Set proxies
	r.SetTrustedProxies([]string{
		"http://localhost:3000",
	})

	//Enable CORS
	CorsConfig := cors.DefaultConfig()
	CorsConfig.AllowCredentials = true
	CorsConfig.AllowOrigins = []string{"http://localhost:3000"}

	r.Use(cors.New(CorsConfig))

	authenticationGroup := r.Group("/auth")
	{
		authenticationGroup.GET("/:provider/login", api.LoginHandler)
		authenticationGroup.GET("/:provider/callback", api.CallbackHandler)

		authenticationGroup.Use(middleware.AuthGuard())
		{
			authenticationGroup.GET("/me", api.AuthMeHandler)
			authenticationGroup.GET("/logout", api.LogoutHandler)
		}
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthGuard())
	{
		organizationGroup := apiGroup.Group("/organization")
		{
			organizationGroup.GET("", api.GetOrganizations)
			organizationGroup.GET("/:orgId", api.GetOrganization)
			organizationGroup.POST("/create", api.CreateOrganization)

			credentialGroup := organizationGroup.Group("/credential")
			{
				// credentialGroup.GET("", api.CreateCrendentials)
				credentialGroup.GET("/:credId", api.GetCredential)
				credentialGroup.POST("/:cred_id/upload", api.UploadFile)
				credentialGroup.POST("/create", api.CreateCrendentials)
				credentialGroup.POST("/fields/:credId", api.AddFields)
			}

			teamGroup := organizationGroup.Group("/team")
			{
				teamGroup.GET("", api.GetMembers)
				teamGroup.GET("/:memId", api.GetMember)
				teamGroup.POST("/create", api.CreateTeam)
				teamGroup.POST("/invite", api.InviteMember)
				teamGroup.POST("/add", api.CreateCrendentials)
			}
		}
	}

	return r.Run(":8080")
}
