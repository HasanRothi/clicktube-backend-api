package main

import (
	"linkbook/cli/route_handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/", route_handlers.Home)
	server.GET("/guest", route_handlers.GuestHome("ROtHi"))
	server.Run()
}
