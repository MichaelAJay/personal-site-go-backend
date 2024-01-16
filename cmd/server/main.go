package main

import (
	"fmt"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET"}
	router.Use(cors.New(config))

	router.GET("/", routes.HomeHandler)
	router.GET("/sierpinski", routes.SierpinskiHandler)

	fmt.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}
