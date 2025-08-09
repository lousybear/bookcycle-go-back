package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/db"
	"github.com/lousybear/bookcycle-go-back/models"
	"go.mongodb.org/mongo-driver/bson"
)

func AddBookHandler(c *gin.Context) {
	var input models.Book

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	collection := db.GetCollection("books")

	filter := bson.M{
		"$and": []bson.M{
			{"title": input.Title},
			{"author": input.Author},
		},
	}

	var existing models.Book
	err := collection.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "book already exists"})
		return
	}

	_, err = collection.InsertOne(context.TODO(), input)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add new book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Added book successfully",
	})
}

func GetAllBooksHandler(c *gin.Context) {
	collection := db.GetCollection("books")

	var books []models.Book

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	c.JSON(http.StatusOK, books)
}
