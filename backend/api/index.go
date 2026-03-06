package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/db"
	"main.go/routes"
)

// This file is only created for Vercel serverless deployment
// It is not used in local development
var app *gin.Engine

func init() {
	// Initialize configuration and DB
	cfg := config.LoadConfig()
	db.InitDB(cfg.DBUrl)

	// Setup Gin app
	app = gin.Default()
	routes.RegisterRoutes(app)
}

// Handler is the exact entrypoint Vercel expects
// Takes the request from Vercel and forwards it to the Gin app
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
