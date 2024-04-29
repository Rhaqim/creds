package router

import (
	"github.com/Rhaqim/creds/internal/api"
	"github.com/Rhaqim/creds/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() error {
	r := gin.Default()

	//Enable CORS
	CorsConfig := cors.DefaultConfig()
	CorsConfig.AllowAllOrigins = true

	r.Use(cors.New(CorsConfig))

	authenticationGroup := r.Group("/auth")
	{
		authenticationGroup.GET("/:provider/login", api.LoginHandler)
		authenticationGroup.GET("/:provider/callback", api.CallbackHandler)

		authenticationGroup.Use(middleware.AuthGuard())
		{
			authenticationGroup.GET("/logout", api.LogoutHandler)
		}
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthGuard())
	{
		organizationGroup := apiGroup.Group("/organization")
		{
			organizationGroup.POST("/create", api.CreateOrganization)

			credentialGroup := organizationGroup.Group("/credential")
			{
				credentialGroup.POST("/create", api.CreateCrendentials)
				credentialGroup.POST("/file/:cred_id", api.UploadFile)
			}
		}
	}

	return r.Run(":8080")
}
