package mongodb

import (
	"Discounts/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserModel struct {
	DB *mongo.Database
}

func (m *UserModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	collection := m.DB.Collection("users")
	user := &models.User{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		Role:           "user",
		Created:        time.Now(),
	}
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		if isDuplicateKeyError(err) {
			return models.ErrDuplicateEmail
		}

		log.Println("Error inserting user:", err)
		return err
	}
	log.Printf("User inserted with ID: %v\n", result.InsertedID)
	return nil

}

func (m *UserModel) GetAll() ([]*models.User, error) {
	collection := m.DB.Collection("users")

	// Fetch all users
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error fetching all users:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	users := []*models.User{}

	if err := cursor.All(context.TODO(), &users); err != nil {
		log.Println("Error decoding users:", err)
		return nil, err
	}

	return users, nil
}

func (m *UserModel) GetByID(id string) (*models.User, error) {
	collection := m.DB.Collection("users")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return nil, err
	}
	var user models.User
	err = collection.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&user)
	if err != nil {
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Authenticate(email, password string) (*models.User, error) {
	collection := m.DB.Collection("users")

	// Fetch user by email
	var user models.User
	err := collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		log.Println("Error fetching user by email:", err)
		return nil, models.ErrInvalidCredentials
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		log.Println("Invalid password:", err)
		return nil, models.ErrInvalidCredentials
	}

	return &user, nil
}

func isDuplicateKeyError(err error) bool {
	if writeException, ok := err.(mongo.WriteException); ok {
		for _, writeError := range writeException.WriteErrors {
			// MongoDB error code 11000 corresponds to duplicate key error
			if writeError.Code == 11000 {
				return true
			}
		}
	}
	return false
}
