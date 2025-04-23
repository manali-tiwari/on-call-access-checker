// main.go (updated)
package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/manali-tiwari/on-call-access-checker/backend/internal/api"
	"github.com/manali-tiwari/on-call-access-checker/backend/internal/auth"
)

func main() {
	// Initialize dependencies
	oktaAuth, err := auth.NewOktaAuth()
	if err != nil {
		log.Fatalf("Failed to initialize Okta client: %v", err)
	}

	awsAuth, err := auth.NewAWSAuth()
	if err != nil {
		log.Fatalf("Failed to initialize AWS client: %v", err)
	}

	// Set up router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Log incoming requests
		if c.Request.Method == "POST" {
			body, _ := io.ReadAll(c.Request.Body)
			log.Printf("Incoming request: %s %s\nBody: %s", c.Request.Method, c.Request.URL.Path, string(body))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		c.Next()
	})

	// Register routes
	apiHandler := api.NewHandler(oktaAuth, awsAuth)
	apiHandler.RegisterRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}
