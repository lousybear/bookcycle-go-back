package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/handlers"
)

func RegisterBookRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/books")

	user.POST("/addbook", handlers.AddBookHandler)
	user.GET("/getallbooks", handlers.GetAllBooksHandler)
}
