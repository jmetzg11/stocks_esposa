package routes

import (
	"stocks/backend/api"
	"stocks/backend/database"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(router *gin.Engine) {
	setupCORS(router)

	handler := api.NewHandler(database.DB)

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from Go!"})
		})
		apiRouter.POST("/start_simulation", handler.StartSimulation)
	}
}

func setupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow local frontend in dev
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
