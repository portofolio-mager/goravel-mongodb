package facades

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/driver"

	mongodb "github.com/tonidy/goravel-mongodb"
	"github.com/tonidy/goravel-mongodb/contracts"
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
