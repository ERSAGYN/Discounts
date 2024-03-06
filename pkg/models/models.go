package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateKey       = errors.New("models: duplicate key")
)

type Product struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	//ShopID      primitive.ObjectID `bson:"shop_id"`
	ProductName string    `bson:"product_name"`
	Category    string    `bson:"category"`
	Price       int       `bson:"price"`
	Discount    int       `bson:"discount"`
	Created     time.Time `bson:"created"`
}

type Shop struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID  primitive.ObjectID `bson:"owner_id"`
	ShopName string             `bson:"shop_name"`
	Address  string             `bson:"address"`
	Created  time.Time          `bson:"created"`
	Products []Product          `bson:"products"`
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username"`
	Email          string             `bson:"email"`
	HashedPassword []byte             `bson:"hashed_password"`
	Role           string             `bson:"role"`
	Created        time.Time          `bson:"created"`
}
