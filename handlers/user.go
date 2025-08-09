package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lousybear/bookcycle-go-back/db"
	"github.com/lousybear/bookcycle-go-back/models"
	"github.com/lousybear/bookcycle-go-back/utils"
	"go.mongodb.org/mongo-driver/bson"
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

	_, err = collection.InsertOne(context.TODO(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign up failed"})
		return
	}

	// Send OTP via Twilio
	if err := utils.SendOTP(input.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully, OTP sent to phone",
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

	c.JSON(http.StatusOK, gin.H{
		"message": "signed in successfully",
		"user":    user.Username,
		"token":   token,
	})
}

func VerifyOTPHandler(c *gin.Context) {
	var input struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ok, err := utils.VerifyOTP(input.Phone, input.Code)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})

}
