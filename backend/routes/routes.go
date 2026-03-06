package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main.go/handlers"
	"main.go/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
	// CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Or use os.Getenv("VITE_API_URL") if strict
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Auth Routes
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)

	// Protected Routes group
	group := r.Group("/tests")
	group.Use(middlewares.RequireAuth())

	// Adding inside curly braces to define a local scope for the group variable.
	// Just for visual organization, not necessary for functionality.
	{
		group.GET("/", handlers.GetAllTests)
		group.POST("/", handlers.CreateTest)
		group.POST("/:id/run", handlers.CreateTestRun)
		group.GET("/:id", handlers.GetTestResult)
	}
}
