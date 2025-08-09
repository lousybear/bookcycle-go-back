package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/handlers"
)

func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api")
	api.GET("/ping", handlers.HealthCheck)

	RegisterUserRoutes(api)
	RegisterBookRoutes(api)
}
