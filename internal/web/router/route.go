package router

import (
	"github.com/Rhaqim/creds/internal/api"
	"github.com/gin-gonic/gin"
)

func Init() error {
	r := gin.Default()

	authentication := r.Group("/auth")
	{
		authentication.POST("/:provider/login", api.LoginHandler)
		authentication.POST("/:provider/callback", api.CallbackHandler)
	}

	return r.Run(":8080")
}
