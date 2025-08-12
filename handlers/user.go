package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/db"
	"github.com/lousybear/bookcycle-go-back/models"
	"github.com/lousybear/bookcycle-go-back/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUpHandler(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	collection := db.GetCollection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"username": input.Username},
			{"email": input.Email},
			{"phone": input.Phone},
		},
	}

	var existing models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username, email or phone already in use"})
		return
	}

	hashedPass, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	input.Password = hashedPass

	result, err := collection.InsertOne(context.TODO(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign up failed"})
		return
	}

	userID := result.InsertedID.(primitive.ObjectID).Hex()
	token, err := utils.GenerateJWT(userID, input.Username, input.Email, input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":       userID,
			"username": input.Username,
			"email":    input.Email,
			"phone":    input.Phone,
		},
		"token": token,
	})
}

func SignInHandler(c *gin.Context) {
	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	collection := db.GetCollection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"username": input.Identifier},
			{"email": input.Identifier},
			{"phone": input.Identifier},
		},
	}

	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := utils.CheckPasswordHash(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Username, user.Email, user.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user logged in successfully",
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
		},
		"token": token,
	})
}
