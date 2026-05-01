package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pedro-navarrete/go_apis_start/internal/domain/product"
)

// productMongoRepository implementación del repositorio de productos para MongoDB
type productMongoRepository struct {
	collection *mongo.Collection
}

// NewProductMongoRepository crea una nueva instancia del repositorio de productos para MongoDB
func NewProductMongoRepository(collection *mongo.Collection) product.Repository {
	return &productMongoRepository{collection: collection}
}

// Create crea un nuevo producto en MongoDB
func (r *productMongoRepository) Create(ctx context.Context, p *product.Product) error {
	_, err := r.collection.InsertOne(ctx, p)
	return err
}

// GetByID obtiene un producto por su ID de MongoDB
func (r *productMongoRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	var p product.Product
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetAll obtiene todos los productos con paginación de MongoDB
func (r *productMongoRepository) GetAll(ctx context.Context, limit, offset int) ([]product.Product, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []product.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update actualiza un producto en MongoDB
func (r *productMongoRepository) Update(ctx context.Context, p *product.Product) error {
	p.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"id": p.ID}, p)
	return err
}

// Delete elimina un producto de MongoDB
func (r *productMongoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
