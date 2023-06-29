package main

import (
	"fmt"
	"os"

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

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/cut")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("", controller.ShortenURL)
	protectedRoutes.GET("", controller.GetAllURLs)

	domain := os.Getenv("SERVE_PORT")
	router.Run(domain)
	fmt.Println("server is running at domain: " + domain)
}
