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

type ProductModel struct {
	DB *mongo.Database
}

func (m *ProductModel) Insert(shopID int, productName, category string, price, discount int) error {
	collection := m.DB.Collection("products")
	product := &models.Product{
		ShopID:      shopID,
		ProductName: productName,
		Category:    category,
		Price:       price,
		Discount:    discount,
		Created:     time.Now(),
	}
	result, err := collection.InsertOne(context.TODO(), product)
	if err != nil {
		if isDuplicateKeyError(err) {
			return models.ErrDuplicateKey
		}

		log.Println("Error inserting product:", err)
		return err
	}
	log.Printf("Product inserted with ID: %v\n", result.InsertedID)
	return nil
}

func (m *ProductModel) GetAll() ([]models.Product, error) {
	collection := m.DB.Collection("products")

	// Fetch all products
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error fetching all products:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		log.Println("Error decoding products:", err)
		return nil, err
	}

	return products, nil
}

func (m *ProductModel) GetByID(id string) (*models.Product, error) {
	collection := m.DB.Collection("products")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return nil, err
	}
	var product models.Product
	err = collection.FindOne(context.TODO(), bson.D{{"_id", objectID}}).Decode(&product)
	if err != nil {
		log.Println("Error fetching product by ID:", err)
		return nil, err
	}

	return &product, nil
}

func (m *ProductModel) GetByShop(shopID int) ([]models.Product, error) {
	collection := m.DB.Collection("products")

	// Fetch products by shop ID
	cursor, err := collection.Find(context.TODO(), bson.D{{"shop_id", shopID}})
	if err != nil {
		log.Println("Error fetching products by shop ID:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		log.Println("Error decoding products:", err)
		return nil, err
	}

	return products, nil
}

func (m *ProductModel) GetByCategory(category string) ([]models.Product, error) {
	collection := m.DB.Collection("products")

	// Fetch products by category
	cursor, err := collection.Find(context.TODO(), bson.D{{"category", category}})
	if err != nil {
		log.Println("Error fetching products by category:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		log.Println("Error decoding products:", err)
		return nil, err
	}

	return products, nil
}

// GetByPrice fetches products by price
func (m *ProductModel) GetByPrice(price int) ([]models.Product, error) {
	collection := m.DB.Collection("products")

	// Fetch products by price
	cursor, err := collection.Find(context.TODO(), bson.D{{"price", price}})
	if err != nil {
		log.Println("Error fetching products by price:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		log.Println("Error decoding products:", err)
		return nil, err
	}

	return products, nil
}

// GetByDiscount fetches products by discount
func (m *ProductModel) GetByDiscount(discount int) ([]models.Product, error) {
	collection := m.DB.Collection("products")

	// Fetch products by discount
	cursor, err := collection.Find(context.TODO(), bson.D{{"discount", discount}})
	if err != nil {
		log.Println("Error fetching products by discount:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		log.Println("Error decoding products:", err)
		return nil, err
	}

	return products, nil
}

// GetByCreated fetches products by created time
func (m *ProductModel) GetByCreated(created time.Time) ([]models.Product, error) {
	collection := m.DB.Collection("products")

	// Fetch products by created time
	cursor, err := collection.Find(context.TODO(), bson.D{{"created", bson.D{{"$gte", created}}}})
	if err != nil {
		log.Println("Error fetching products by created time:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		log.Println("Error decoding products:", err)
		return nil, err
	}

	return products, nil
}
