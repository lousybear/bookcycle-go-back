package main

import (
	"log"

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
	routes.RegisterRoutes(router)

	port := utils.GetEnv("PORT", "8080")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
	log.Printf("Server started on port %s", port)
}
