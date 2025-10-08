package facades

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/driver"

	mongodb "github.com/portofolio-mager/goravel-mongodb"
	"github.com/portofolio-mager/goravel-mongodb/contracts"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB returns a MongoDB client instance
func MongoDB(connection string) (contracts.Client, error) {
	if mongodb.App == nil {
		return nil, fmt.Errorf("please register mongodb service provider")
	}

	instance, err := mongodb.App.MakeWith(mongodb.Binding, map[string]any{
		"connection": connection,
	})
	if err != nil {
		return nil, err
	}

	return instance.(contracts.Client), nil
}

// MongoDBDriver returns the MongoDB driver (for Goravel compatibility)
func MongoDBDriver(connection string) (driver.Driver, error) {
	if mongodb.App == nil {
		return nil, fmt.Errorf("please register mongodb service provider")
	}

	instance, err := mongodb.App.MakeWith(mongodb.Binding, map[string]any{
		"connection": connection,
	})
	if err != nil {
		return nil, err
	}

	return instance.(*mongodb.MongoDB), nil
}

// Database returns a MongoDB database instance
func Database(name string, connection ...string) (contracts.Database, error) {
	conn := "mongodb"
	if len(connection) > 0 {
		conn = connection[0]
	}

	client, err := MongoDB(conn)
	if err != nil {
		return nil, err
	}

	return client.Database(name), nil
}

// Collection returns a MongoDB collection instance
func Collection(collection string, database ...string) (contracts.Collection, error) {
	client, err := MongoDB("mongodb")
	if err != nil {
		return nil, err
	}

	return client.Collection(collection, database...), nil
}

// NativeCollection returns the native *mongo.Collection for direct driver usage
func NativeCollection(collection string, database ...string) (*mongo.Collection, error) {
	col, err := Collection(collection, database...)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("mongodb collection is nil")
	}
	return col.Native(), nil
}

// NativeClient returns the native *mongo.Client for the given connection
func NativeClient(connection ...string) (*mongo.Client, error) {
	conn := "mongodb"
	if len(connection) > 0 && connection[0] != "" {
		conn = connection[0]
	}

	client, err := MongoDB(conn)
	if err != nil {
		return nil, err
	}
	native := client.Native()
	if native == nil {
		return nil, fmt.Errorf("mongodb client is nil")
	}
	return native, nil
}

// NativeDatabase returns the native *mongo.Database for a database name and optional connection
func NativeDatabase(name string, connection ...string) (*mongo.Database, error) {
	db, err := Database(name, connection...)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("mongodb database is nil")
	}
	native := db.Native()
	if native == nil {
		return nil, fmt.Errorf("mongodb database native is nil")
	}
	return native, nil
}
