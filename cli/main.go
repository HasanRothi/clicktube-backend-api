package main

import (
	"linkbook/cli/db"
	"linkbook/cli/middlewares"
	"linkbook/cli/route_handlers"

	"github.com/gin-gonic/gin"
)

func require() {
	middlewares.Logger()
	db.Connect()
}
func main() {
	require()

	server := gin.New()
	server.Use(gin.Recovery())
	server.GET("/", route_handlers.Home)
	server.GET("/guest", route_handlers.GuestHome("ROtHi"))
	server.Run()
}
