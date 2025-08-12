package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/db"
	"github.com/lousybear/bookcycle-go-back/routes"
	"github.com/lousybear/bookcycle-go-back/utils"
)

func main() {
	utils.LoadEnv()

	db.Init()
	defer db.Disconnect()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(router)

	port := utils.GetEnv("PORT", "8080")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
	log.Printf("Server started on port %s", port)
}
