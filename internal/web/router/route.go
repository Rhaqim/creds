package router

import (
	"github.com/Rhaqim/creds/internal/api"
	"github.com/Rhaqim/creds/internal/web/middleware"
	"github.com/gin-gonic/gin"
)

func Init() error {
	r := gin.Default()

	authenticationGroup := r.Group("/auth")
	{
		authenticationGroup.POST("/:provider/login", api.LoginHandler)
		authenticationGroup.POST("/:provider/callback", api.CallbackHandler)
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthGuard())
	{
		organizationGroup := apiGroup.Group("/organization")
		{
			organizationGroup.POST("/organization", api.CreateOrganization)

			credentialGroup := organizationGroup.Group("/credential")
			{
				credentialGroup.POST("/credentials", api.CreateCrendentials)
				credentialGroup.POST("/file/:cred_id", api.UploadFile)
			}
		}
	}

	return r.Run(":8080")
}
