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
	router.Use(cORSMiddleware())
	router.GET("/api/:url", controller.ResolveURL)
	//Adding Healthchecks for Render Deployment
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

func cORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
