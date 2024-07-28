package routes

import "github.com/gin-gonic/gin"

var (
	Router = gin.Default()
)

func init() {
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
