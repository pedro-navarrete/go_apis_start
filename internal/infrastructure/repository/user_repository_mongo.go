package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pedro-navarrete/go_apis_start/internal/domain/user"
)

// userMongoRepository implementación del repositorio de usuarios para MongoDB
type userMongoRepository struct {
	collection *mongo.Collection
}

// NewUserMongoRepository crea una nueva instancia del repositorio de usuarios para MongoDB
func NewUserMongoRepository(collection *mongo.Collection) user.Repository {
	return &userMongoRepository{collection: collection}
}

// Create crea un nuevo usuario en MongoDB
func (r *userMongoRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.collection.InsertOne(ctx, u)
	return err
}

// GetByID obtiene un usuario por su ID de MongoDB
func (r *userMongoRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByUsername obtiene un usuario por su username de MongoDB
func (r *userMongoRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail obtiene un usuario por su email de MongoDB
func (r *userMongoRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetAll obtiene todos los usuarios con paginación de MongoDB
func (r *userMongoRepository) GetAll(ctx context.Context, limit, offset int) ([]user.User, int64, error) {
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

	var users []user.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update actualiza un usuario en MongoDB
func (r *userMongoRepository) Update(ctx context.Context, u *user.User) error {
	u.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"id": u.ID}, u)
	return err
}

// Delete elimina un usuario de MongoDB
func (r *userMongoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
