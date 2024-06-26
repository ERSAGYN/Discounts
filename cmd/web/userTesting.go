package main

import (
	"Discounts/pkg/models"
	"Discounts/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func test(database *mongo.Database) {
	userModel := &mongodb.UserModel{DB: database}
	err := userModel.Insert("testuser", "test@example.com", "testpassword")
	if err != nil {
		log.Printf("Failed to insert user: %v\n", err)
	}

	err = userModel.Insert("duplicateuser", "test@example.com", "testpassword2")
	if err != nil {
		log.Printf("Failed to insert user: %v\n", err)
	}

	users, err := userModel.GetAll()
	if err != nil {
		log.Printf("Failed to get all users: %v", err)
	}
	for i, user := range users {
		log.Printf("User %d: ID=%d, Username=%s, Email=%s, Role=%s, Created=%s",
			i+1, user.ID, user.Username, user.Email, user.Role, user.Created.Format(time.RFC3339))
	}
	log.Print(users[0].ID)
	userByID, err := userModel.GetByID("65e70f157d38634495c2893f")
	if err != nil {
		log.Printf("Failed to get user by ID: %v", err)
	}
	log.Printf("Fetched user by ID: ID=%d, Username=%s, Email=%s, Role=%s, Created=%s",
		userByID.ID, userByID.Username, userByID.Email, userByID.Role, userByID.Created.Format(time.RFC3339))
	shopName := "Test Shop"
	address := "Test Address"
	products := []models.Product{
		{
			ProductName: "Product 1",
			Price:       10,
		},
		{
			ProductName: "Product 2",
			Price:       20,
		}}
	shopModel := &mongodb.ShopModel{DB: database}
	err = shopModel.Insert(userByID.ID, shopName, address, products)
	if err != nil {
		log.Printf("Error inserting shop: %v", err)
	}

	// Test case: Authenticate with correct credentials
	authenticatedUser, err := userModel.Authenticate("test@example.com", "testpassword")
	if err != nil {
		log.Printf("Failed to authenticate user: %v", err)
	}

	log.Printf("Authenticated user: ID=%s, Username=%s, Email=%s, Role=%s, Created=%s",
		authenticatedUser.ID, authenticatedUser.Username, authenticatedUser.Email,
		authenticatedUser.Role, authenticatedUser.Created.Format(time.RFC3339))

	if authenticatedUser.Username != "testuser" {
		log.Printf("Expected authenticated username 'testuser', got '%s'", authenticatedUser.Username)
	}

	// Test case: Authenticate with incorrect password
	_, err = userModel.Authenticate("test@example.com", "wrongpassword")
	if err != models.ErrInvalidCredentials {
		log.Printf("Expected ErrInvalidCredentials, got %v", err)
	}

	// Test case: Authenticate with incorrect email
	_, err = userModel.Authenticate("wrong@example.com", "testpassword")
	if err != models.ErrInvalidCredentials {
		log.Printf("Expected ErrInvalidCredentials, got %v", err)
	}
}
