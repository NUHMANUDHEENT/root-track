package main

import (
	"fmt"
	"log"
	"os"
	"roottrack-backend/config"
	"roottrack-backend/routes"
)

func main() {
	// 1. Connect to Database
	config.ConnectDatabase()

	// 2. Setup Router
	r := routes.SetupRouter()

	
	// 3. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	fmt.Printf("Server is running on port %s\n", port)
	r.Run(":" + port)
}
