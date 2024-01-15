package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/routes"
)

func main() {
	http.HandleFunc("/", routes.HomeHandler)
	http.HandleFunc("/sierpinski", routes.SierpinskiHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
