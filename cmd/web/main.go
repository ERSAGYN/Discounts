package main

import (
	"Discounts/pkg/models/mongodb"
	"context"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	users         *mongodb.UserModel
	products      *mongodb.ProductModel
	shops         *mongodb.ShopModel
	templateCache map[string]*template.Template
}

func main() {
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	addr := flag.String("addr", ":4000", "HTTP network address")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	database := client.Database("discounts")

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		users:         &mongodb.UserModel{DB: database},
		products:      &mongodb.ProductModel{DB: database},
		shops:         &mongodb.ShopModel{DB: database},
		templateCache: templateCache,
	}

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = database.Collection("users").Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		log.Fatal(err)
	}

	test(database)

	srv := &http.Server{
		Addr:           *addr,
		MaxHeaderBytes: 524288,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	defer client.Disconnect(context.TODO())

	errorLog.Fatal(err)

}
