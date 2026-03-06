package main

import (
	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/db"
	"main.go/routes"
)

func main() {
	// Load configuration and initialize the database connection
	cfg := config.LoadConfig()
	db.InitDB(cfg.DBUrl)

	// Create a Gin router with default middleware (creating a http engine)
	router := gin.Default()

	// Passing router engine to the routes package to register all the routes
	routes.RegisterRoutes(router)
	router.Run(":8080") // Start the server on port 8080

}
