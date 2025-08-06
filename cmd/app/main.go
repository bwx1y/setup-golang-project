package main

import (
	"log"
	"os"
	"setup-preoject/app/config"
	"setup-preoject/app/route"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connection database
	db := config.ConnectDatabase()

	// connection redis for cache
	cache, err := config.ConnectionRedist()

	// migration database
	config.MigrationDatabase(db)

	// setup gin
	engine := gin.Default()

	// register route
	route.SetupRoutes(engine, db, cache)

	// get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// running app
	err = engine.Run(":" + port)
	if err != nil {
		return
	}
}
