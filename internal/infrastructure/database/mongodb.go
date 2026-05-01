package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pedro-navarrete/go_apis_start/internal/config"
	applogger "github.com/pedro-navarrete/go_apis_start/internal/utils/logger"
)

// MongoDB wraps the MongoDB client and database
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDBConnection crea una nueva conexión a MongoDB
func NewMongoDBConnection(cfg *config.MongoDBConfig) (*MongoDB, error) {
	// Crear contexto con timeout para la conexión
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configurar opciones del cliente
	clientOptions := options.Client().ApplyURI(cfg.URI)

	// Conectar a MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	// Verificar la conexión con ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error haciendo ping a MongoDB: %w", err)
	}

	applogger.Info("✅ Conexión a MongoDB establecida")

	return &MongoDB{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

// Close cierra la conexión a MongoDB
func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// Ping verifica que la conexión a MongoDB esté activa
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

// Collection retorna una colección específica de la base de datos
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
