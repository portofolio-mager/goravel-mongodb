package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database"
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"

	"github.com/tonidy/goravel-mongodb/contracts"
)

var _ driver.Driver = &MongoDB{}
var _ contracts.Client = &MongoDB{}

type MongoDB struct {
	config contracts.ConfigBuilder
	client *mongo.Client
	log    log.Log
}

func NewMongoDB(config config.Config, log log.Log, connection string) *MongoDB {
	return &MongoDB{
		config: NewConfig(config, connection),
		log:    log,
	}
}

func (m *MongoDB) connect() error {
	if m.client != nil {
		return nil
	}

	writers := m.config.Writers()
	if len(writers) == 0 {
		return errors.DatabaseConfigNotFound
	}

	fullConfig := writers[0]
	uri := fullConfig.URI
	if uri == "" {
		return fmt.Errorf("MongoDB URI is required")
	}

	clientOptions := options.Client().ApplyURI(uri)

	// Apply authentication if provided
	if fullConfig.Username != "" && fullConfig.Password != "" {
		credential := options.Credential{
			Username: fullConfig.Username,
			Password: fullConfig.Password,
		}
		if fullConfig.AuthSource != "" {
			credential.AuthSource = fullConfig.AuthSource
		}
		clientOptions.SetAuth(credential)
	}

	// Apply connection pool settings
	if fullConfig.MaxPoolSize != nil {
		clientOptions.SetMaxPoolSize(*fullConfig.MaxPoolSize)
	}
	if fullConfig.MinPoolSize != nil {
		clientOptions.SetMinPoolSize(*fullConfig.MinPoolSize)
	}

	// Apply timeouts
	if fullConfig.ConnectTimeout != nil {
		clientOptions.SetConnectTimeout(time.Duration(*fullConfig.ConnectTimeout) * time.Second)
	}
	if fullConfig.ServerTimeout != nil {
		clientOptions.SetServerSelectionTimeout(time.Duration(*fullConfig.ServerTimeout) * time.Second)
	}

	// Create client
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.client = client
	return nil
}

// Driver interface implementation (required for Goravel)
func (m *MongoDB) Docker() (docker.DatabaseDriver, error) {
	writers := m.config.Writers()
	if len(writers) == 0 {
		return nil, errors.DatabaseConfigNotFound
	}

	return NewDocker(writers[0].Database), nil
}

func (m *MongoDB) Grammar() driver.Grammar {
	return NewGrammar(m.log)
}

func (m *MongoDB) Pool() database.Pool {
	return database.Pool{
		Writers: m.fullConfigsToConfigs(m.config.Writers()),
	}
}

func (m *MongoDB) Processor() driver.Processor {
	return NewProcessor()
}

// Client interface implementation (MongoDB-specific)
func (m *MongoDB) Native() *mongo.Client {
	if err := m.connect(); err != nil {
		m.log.Errorf("Failed to connect to MongoDB: %v", err)
		return nil
	}
	return m.client
}

func (m *MongoDB) Database(name ...string) contracts.Database {
	if err := m.connect(); err != nil {
		m.log.Errorf("Failed to connect to MongoDB: %v", err)
		return nil
	}

	var dbName string
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	} else {
		writers := m.config.Writers()
		if len(writers) > 0 {
			dbName = writers[0].Database
		}
	}

	if dbName == "" {
		m.log.Error("No database name specified")
		return nil
	}

	return NewDatabase(m.client, m.config, dbName)
}

func (m *MongoDB) Collection(collection string, database ...string) contracts.Collection {
	db := m.Database(database...)
	if db == nil {
		return nil
	}
	return db.Collection(collection)
}

func (m *MongoDB) Ping() error {
	if err := m.connect(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Ping(ctx, nil)
}

func (m *MongoDB) Close() error {
	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// fullConfigToDialector creates a GORM dialector for MongoDB (similar to PostgreSQL implementation)
func (m *MongoDB) fullConfigToDialector(config contracts.FullConfig) gorm.Dialector {
	// Create DSN from MongoDB URI
	dsn := config.URI
	if dsn == "" {
		dsn = "mongodb://localhost:27017/" + config.Database
	}
	return Open(dsn)
}

func (m *MongoDB) fullConfigsToConfigs(fullConfigs []contracts.FullConfig) []database.Config {
	configs := make([]database.Config, len(fullConfigs))
	for i, fullConfig := range fullConfigs {
		configs[i] = database.Config{
			Connection: fullConfig.Connection,
			Database:   fullConfig.Database,
			Driver:     Name,
			Dialector:  m.fullConfigToDialector(fullConfig),
		}
	}
	return configs
}
