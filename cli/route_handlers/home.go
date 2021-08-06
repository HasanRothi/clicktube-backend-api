package route_handlers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	panic("Sentry Test Error")
	c.JSON(200, gin.H{
		"message": "Welcome to Link Home",
	})
}
func GuestHome(msg string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome " + msg,
		})
	}
}
