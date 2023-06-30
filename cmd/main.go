package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tolumadamori/scissor/pkg/config"
	"github.com/tolumadamori/scissor/pkg/controller"
	"github.com/tolumadamori/scissor/pkg/middleware"
	"github.com/tolumadamori/scissor/pkg/models"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	loadDatabase()
	serveApplication()

}

func loadDatabase() {
	config.ConnectDB()
	config.Database.AutoMigrate(&models.User{}, &models.URL{})

}

func serveApplication() {
	router := gin.Default()
	router.GET("/:url", controller.ResolveURL)
	//Adding Healthchecks for Render Deployment
	router.GET("/healthz", controller.Healthchecks)
	router.GET("/", controller.Healthchecks)

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/cut")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("", controller.ShortenURL)
	protectedRoutes.GET("", controller.GetAllURLs)

	router.Run(":10000")
	fmt.Println("server is running on port")
}
