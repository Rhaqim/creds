package router

import "github.com/gin-gonic/gin"

func Init() error {
	r := gin.Default()

	return r.Run(":8080")
}
