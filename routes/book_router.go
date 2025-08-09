package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/handlers"
)

func RegisterBookRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/books")

	user.POST("/", handlers.AddBookHandler)
	user.GET("/", handlers.GetAllBooksHandler)
}
