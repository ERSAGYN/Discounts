package mongodb

import (
	"Discounts/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ShopModel struct {
	DB *mongo.Database
}

func (m *ShopModel) Insert(ownerID int, shopName, address string) error {
	collection := m.DB.Collection("shops")
	shop := &models.Shop{
		OwnerID:  ownerID,
		ShopName: shopName,
		Address:  address,
		Created:  time.Now(),
	}
	result, err := collection.InsertOne(context.TODO(), shop)
	if err != nil {
		if isDuplicateKeyError(err) {
			return models.ErrDuplicateKey
		}

		log.Println("Error inserting shop:", err)
		return err
	}
	log.Printf("Shop inserted with ID: %v\n", result.InsertedID)
	return nil
}

func (m *ShopModel) GetAll() ([]models.Shop, error) {
	collection := m.DB.Collection("shops")

	// Fetch all shops
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error fetching all shops:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var shops []models.Shop
	if err := cursor.All(context.TODO(), &shops); err != nil {
		log.Println("Error decoding shops:", err)
		return nil, err
	}

	return shops, nil
}

func (m *ShopModel) GetByID(id string) (*models.Shop, error) {
	collection := m.DB.Collection("shops")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return nil, err
	}
	var shop models.Shop
	err = collection.FindOne(context.TODO(), bson.D{{"_id", objectID}}).Decode(&shop)
	if err != nil {
		log.Println("Error fetching shop by ID:", err)
		return nil, err
	}

	return &shop, nil
}

func (m *ShopModel) GetByOwner(ownerID int) ([]models.Shop, error) {
	collection := m.DB.Collection("shops")

	// Fetch shops by owner ID
	cursor, err := collection.Find(context.TODO(), bson.D{{"owner_id", ownerID}})
	if err != nil {
		log.Println("Error fetching shops by owner ID:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var shops []models.Shop
	if err := cursor.All(context.TODO(), &shops); err != nil {
		log.Println("Error decoding shops:", err)
		return nil, err
	}

	return shops, nil
}

func (m *ShopModel) GetCreated(created time.Time) ([]models.Shop, error) {
	collection := m.DB.Collection("shops")

	// Fetch shops by created time
	cursor, err := collection.Find(context.TODO(), bson.D{{"created", bson.D{{"$gte", created}}}})
	if err != nil {
		log.Println("Error fetching shops by created time:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var shops []models.Shop
	if err := cursor.All(context.TODO(), &shops); err != nil {
		log.Println("Error decoding shops:", err)
		return nil, err
	}

	return shops, nil
}
