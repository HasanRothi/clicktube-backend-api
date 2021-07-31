package main

import (
	"linkbook/cli/db"
	"linkbook/cli/middlewares"
	"linkbook/cli/route_handlers"

	"github.com/gin-gonic/gin"
)

// var DB string

func require() {
	middlewares.Logger()
	db.Connect()
	// fmt.Println(db.DatabaseList[0])
	// DB = db.DatabaseList[0]
}
func main() {
	require()

	server := gin.New()
	server.Use(gin.Recovery())
	server.GET("/", route_handlers.Home)
	server.GET("/guest", route_handlers.GuestHome("ROtHi"))
	server.GET("/links", route_handlers.GetAllLink)
	server.Run()
}
