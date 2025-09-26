package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tonidy/goravel-mongodb/contracts"
)

var _ contracts.Database = &Database{}

type Database struct {
	client   *mongo.Client
	database *mongo.Database
	config   contracts.ConfigBuilder
}

func NewDatabase(client *mongo.Client, config contracts.ConfigBuilder, name string) *Database {
	return &Database{
		client:   client,
		database: client.Database(name),
		config:   config,
	}
}

func (d *Database) Native() *mongo.Database {
	return d.database
}

func (d *Database) Collection(name string) contracts.Collection {
	return NewCollection(d.client, d.config, name, d.database.Name())
}

func (d *Database) CreateCollection(name string, opts ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var createOpts *options.CreateCollectionOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.CreateCollectionOptions); ok {
			createOpts = opt
		}
	}

	return d.database.CreateCollection(ctx, name, createOpts)
}

func (d *Database) ListCollections() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := d.database.ListCollections(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var collections []string
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		if name, ok := result["name"].(string); ok {
			collections = append(collections, name)
		}
	}

	return collections, cursor.Err()
}

func (d *Database) Drop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return d.database.Drop(ctx)
}

func (d *Database) Name() string {
	return d.database.Name()
}
