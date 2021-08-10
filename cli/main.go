package main

import (
	"linkbook/cli/db"
	"linkbook/cli/middlewares"
	"linkbook/cli/route_handlers"
	"os"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func require() {
	// middlewares.Logger()
	db.Connect()
	middlewares.SentryInit()
}
func main() {
	require()
	// var Level1 = []string{"User", "Admin", "SuperAdmin"}
	// var Level2 = []string{"Admin", "SuperAdmin"}
	// var Level3 = []string{"SuperAdmin"}
	server := gin.New()
	server.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	server.Use(gin.Recovery())
	server.Use(middlewares.Recover)
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("UI")},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Welcome routes
	server.GET("/", route_handlers.Home)
	server.GET("/guest", route_handlers.GuestHome("ROtHi"))

	//link routes
	server.GET("/links", route_handlers.GetAllLink)
	server.GET("/link/:id", route_handlers.GetSingleLink)
	server.GET("/links/pending", route_handlers.GetPendingLinks)
	server.GET("/links/popular", route_handlers.GetPopularLinks)
	server.POST("/link", route_handlers.PostSingleLink)
	server.POST("/link/published", route_handlers.PublishedSingleLink)

	//user routes
	server.GET("/users", route_handlers.GetAllUser)
	server.GET("/user/:id", route_handlers.SingleUserLinks)
	// server.GET("/links/user/:id", route_handlers.GetSingleUserLinks)/
	server.POST("/user", route_handlers.PostSingleUser)

	//Login routes
	server.POST("/credentials/login", route_handlers.Login)

	server.Run()
}
