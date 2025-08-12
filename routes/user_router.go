package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/handlers"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")

	user.POST("/signup", handlers.SignUpHandler)
	user.POST("/signin", handlers.SignInHandler)
}
