package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	Description string             `bson:"description" json:"description"`
	Price       int                `bson:"price" json:"price"`
	Genre       string             `bson:"genre" json:"genre"`
	Image       string             `bson:"image" json:"image"`
}
