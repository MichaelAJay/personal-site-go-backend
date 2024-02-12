package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/middleware"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/models"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/routes"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/auth"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/contact"
	db_client "github.com/MichaelAJay/personal-site-go-backend/pkg/services/db-client"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/secrets"
	"github.com/MichaelAJay/personal-site-go-backend/pkg/services/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func registerPprofRoutes(router *gin.Engine) {
	r := router.Group("/debug/pprof")

	r.GET("/", gin.WrapF(pprof.Index))
	r.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	r.GET("/profile", gin.WrapF(pprof.Profile))
	r.POST("/symbol", gin.WrapF(pprof.Symbol))
	r.GET("/symbol", gin.WrapF(pprof.Symbol))
	r.GET("/trace", gin.WrapF(pprof.Trace))
	r.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
	r.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	r.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	r.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	r.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	r.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println((".env not found"))
	}

	env := os.Getenv("ENV")
	if env == "" {
		log.Println("ENV variable not set, defaulting to \"development\"")
		env = "development"
	}

	secretManagerService, err := secrets.NewSecretManagerService()
	if err != nil {
		log.Fatalf("Error instantiating secret manager service: %v", err)
	}

	dbDsn, err := secretManagerService.GetSecret("DB_DSN")
	if err != nil {
		log.Fatalf("Error during configuration: %v", err)
	}

	// Database connection & automigration (non production)
	db := db_client.DbClient(dbDsn)

	if env == "local" {
		db.AutoMigrate(&models.Contact{}, &models.User{})
	}

	router := gin.Default()
	if env != "production" {
		registerPprofRoutes(router)
	}

	// Configure CORS
	config := cors.DefaultConfig()

	var origins []string
	if env == "production" {
		origins = []string{"https://michaelajay.github.io"}
	} else {
		origins = []string{"http://localhost:3000"}
	}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET"}
	router.Use(cors.New(config))

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	// Configure rate limiting
	limiter := rate.NewLimiter(20, 10)
	router.Use(middleware.RateLimitMiddleware(limiter))

	// Routes

	router.GET("/", routes.HomeHandler)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Health ok")
	})

	router.GET("/sierpinski", routes.SierpinskiHandler)

	userService := user.NewUserService(db)

	// Contact Form
	contactService := contact.NewContactService()
	router.POST("/contact", routes.PostContactFormHandler(contactService))
	router.GET("/contact/list", routes.GetContactFormListHandler(contactService))
	router.GET("/contact", routes.GetMessageHandler(contactService))
	router.PATCH("/contact/toggle-read-status/:id", routes.ToggleMessageReadStatus(contactService))

	// Auth
	secret, err := secretManagerService.GetSecret("JWT_SECRET")
	if err != nil {
		log.Fatalf("Error during configuration: %v", err)
	}
	authService, err := auth.NewAuthService(db_client.Db, userService, secret)
	if err != nil {
		log.Fatalf("Error instantiating auth service: %v", err)
	}
	router.POST("/sign-up", routes.SignUp(authService))
	router.POST("/sign-in", routes.SignIn(authService))

	// User
	// We shouldn't send the email plain - we should send a signed token
	// This way the email can't be circumvented - the correct token must be sent
	router.PATCH("/user/verify", routes.VerifyUser(userService))

	fmt.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}
