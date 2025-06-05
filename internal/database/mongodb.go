package database

import (
	"context"
	"fmt"
	"time"

	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB represents a MongoDB database connection
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// MongoDBConfig represents MongoDB configuration
type MongoDBConfig struct {
	URI      string
	Database string
	Username string
	Password string
	Timeout  time.Duration
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg *MongoDBConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Create client options
	clientOptions := options.Client().ApplyURI(cfg.URI)
	if cfg.Username != "" && cfg.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Connected to MongoDB database")
	return &MongoDB{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	if m.Client != nil {
		return m.Client.Disconnect(ctx)
	}
	return nil
}

// Ping checks if the MongoDB connection is alive
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.Primary())
}

// Collection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// WithTransaction executes a function within a transaction
func (m *MongoDB) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := m.Client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn)
}

// ExecuteWithRetry executes a function with retry logic
func (m *MongoDB) ExecuteWithRetry(ctx context.Context, fn func(context.Context) error, maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn(ctx)
		if err == nil {
			return nil
		}
		
		logger.Warnf("MongoDB operation failed (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}
	return err
}

