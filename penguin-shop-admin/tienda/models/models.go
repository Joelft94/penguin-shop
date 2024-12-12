package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string            `bson:"name"`
	Description string            `bson:"description"`
	Price       float64           `bson:"price"`
	Stock       int               `bson:"stock"`
	Type        string            `bson:"type"`
	CreatedAt   time.Time         `bson:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt"`
	V           int               `bson:"__v"`
}

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"productId"`
	Quantity  int               `bson:"quantity"`
}

type Order struct {
	ID           primitive.ObjectID `bson:"_id"`
	OrderNumber  string            `bson:"orderNumber"`
	IglooAddress string            `bson:"iglooAddress"`
	Items        []OrderItem       `bson:"items"`
	TotalAmount  float64           `bson:"totalAmount"`
	Status       string            `bson:"status"`
	CreatedAt    time.Time         `bson:"createdAt"`
}