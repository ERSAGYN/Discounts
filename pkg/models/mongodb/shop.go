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

func (m *ShopModel) Insert(userId int, shopName, address string, products []models.Product) error {
	collection := m.DB.Collection("shops")

	shop := &models.Shop{
		OwnerID:  userId,
		ShopName: shopName,
		Address:  address,
		Products: products,
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
	log.Printf("Shop inserted with ID:\n", result.InsertedID.(primitive.ObjectID))
	return nil
}

func (m *ShopModel) GetAll() ([]*models.Shop, error) {
	collection := m.DB.Collection("shops")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error fetching all shops:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	shops := []*models.Shop{}

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

func (m *ShopModel) GetByOwner(ownerID int) ([]*models.Shop, error) {
	collection := m.DB.Collection("shops")

	// Fetch shops by owner ID
	cursor, err := collection.Find(context.TODO(), bson.D{{"owner_id", ownerID}})
	if err != nil {
		log.Println("Error fetching shops by owner ID:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	shops := []*models.Shop{}

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

func (m *ShopModel) AddProducts(shopID primitive.ObjectID, newProducts []models.Product) error {
	collection := m.DB.Collection("shops")

	// Find the existing shop
	var existingShop models.Shop
	err := collection.FindOne(context.TODO(), bson.D{{"_id", shopID}}).Decode(&existingShop)
	if err != nil {
		log.Println("Error fetching existing shop:", err)
		return err
	}

	// Find the highest existing ID in the products array
	maxProductID := 0
	for _, product := range existingShop.Products {
		if product.ID > maxProductID {
			maxProductID = product.ID
		}
	}

	// Increment the ID for the new products
	nextProductID := maxProductID + 1

	// Increment the ID for each new product
	for i := range newProducts {
		newProducts[i].ID = nextProductID
		nextProductID++
	}

	// Append new products to the existing shop's products array
	existingShop.Products = append(existingShop.Products, newProducts...)

	// Update the existing shop with the new products
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.D{{"_id", shopID}},
		bson.D{{"$set", bson.D{{"products", existingShop.Products}}}},
	)
	if err != nil {
		log.Println("Error updating existing shop:", err)
		return err
	}

	log.Printf("Products added to shop with ID: %v\n", shopID)
	return nil
}

/*var maxIDResult struct {
	ID int `bson:"id"`
}
err := collection.FindOne(context.TODO(), bson.D{}, &options.FindOneOptions{
	Sort: bson.D{{"id", -1}},
}).Decode(&maxIDResult)

if err != nil && err != mongo.ErrNoDocuments {
	log.Println("Error finding max product ID:", err)
	return err
}

// Use the next available ID
nextID := 1
if maxIDResult.ID > 0 {
	nextID = maxIDResult.ID + 1
}*/
