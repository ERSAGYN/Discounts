package main

import (
	"Discounts/pkg/models/mongodb"
	"context"
	"fmt"
	"github.com/golangcollege/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	user          *mongodb.UserModel
	shop          *mongodb.ShopModel
	product       *mongodb.ProductModel
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	database := client.Database("discounts")

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = database.Collection("users").Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		log.Fatal(err)
	}

	test(database)

	defer client.Disconnect(context.TODO())

}
