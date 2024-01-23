package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/config"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/routes"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/auth"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/contact"
	db_client "github.com/MichaelAJay/personal-site-go-backend/pkg/services/db-client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadConfig(path string) (config.Config, error) {
	var cfg config.Config
	configFile, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	return cfg, err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(("Error loading .env file"))
	}

	// Get rid of config
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	env := cfg.Env
	log.Printf("ENV: %s\n", env)

	// Database connection & automigration (non production)
	db := db_client.DbClient(cfg.DbDsn)

	if env == "local" {
		db.AutoMigrate(&models.Contact{})
	}

	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET"}
	router.Use(cors.New(config))

	router.GET("/", routes.HomeHandler)
	router.GET("/sierpinski", routes.SierpinskiHandler)

	// Contact Form
	contactService := contact.NewContactService()
	router.POST("/contact", routes.PostContactFormHandler(contactService))
	router.GET("/contact/list", routes.GetUnreadContactFormListHandler(contactService))
	router.GET("/contact", routes.GetMessageHandler(contactService))
	router.PATCH("/contact/toggle-read-status/:id", routes.ToggleMessageReadStatus(contactService))

	// Auth
	authService, err := auth.NewAuthService()
	if err != nil {
		log.Fatalf("Error instantiating auth service: %v", err)
	}

	// Handle err
	router.POST("/sign-up", routes.SignUp(authService))
	router.POST("/sign-in", routes.SignIn(authService))

	fmt.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}
