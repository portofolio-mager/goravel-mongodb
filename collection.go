package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/portofolio-mager/goravel-mongodb/contracts"
)

var _ contracts.Collection = &Collection{}

type Collection struct {
	client     *mongo.Client
	collection *mongo.Collection
	config     contracts.ConfigBuilder
}

func NewCollection(client *mongo.Client, config contracts.ConfigBuilder, name string, database string) *Collection {
	return &Collection{
		client:     client,
		collection: client.Database(database).Collection(name),
		config:     config,
	}
}

func (c *Collection) Native() *mongo.Collection {
	return c.collection
}

// Basic CRUD operations
func (c *Collection) FindOne(filter interface{}, result interface{}, opts ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var findOpts *options.FindOneOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.FindOneOptions); ok {
			findOpts = opt
		}
	}

	return c.collection.FindOne(ctx, filter, findOpts).Decode(result)
}

func (c *Collection) Find(filter interface{}, opts ...interface{}) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var findOpts *options.FindOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.FindOptions); ok {
			findOpts = opt
		}
	}

	return c.collection.Find(ctx, filter, findOpts)
}

func (c *Collection) InsertOne(document interface{}, opts ...interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var insertOpts *options.InsertOneOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.InsertOneOptions); ok {
			insertOpts = opt
		}
	}

	return c.collection.InsertOne(ctx, document, insertOpts)
}

func (c *Collection) InsertMany(documents []interface{}, opts ...interface{}) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var insertOpts *options.InsertManyOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.InsertManyOptions); ok {
			insertOpts = opt
		}
	}

	return c.collection.InsertMany(ctx, documents, insertOpts)
}

func (c *Collection) UpdateOne(filter interface{}, update interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updateOpts *options.UpdateOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.UpdateOptions); ok {
			updateOpts = opt
		}
	}

	return c.collection.UpdateOne(ctx, filter, update, updateOpts)
}

func (c *Collection) UpdateMany(filter interface{}, update interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updateOpts *options.UpdateOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.UpdateOptions); ok {
			updateOpts = opt
		}
	}

	return c.collection.UpdateMany(ctx, filter, update, updateOpts)
}

func (c *Collection) DeleteOne(filter interface{}, opts ...interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var deleteOpts *options.DeleteOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.DeleteOptions); ok {
			deleteOpts = opt
		}
	}

	return c.collection.DeleteOne(ctx, filter, deleteOpts)
}

func (c *Collection) DeleteMany(filter interface{}, opts ...interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var deleteOpts *options.DeleteOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.DeleteOptions); ok {
			deleteOpts = opt
		}
	}

	return c.collection.DeleteMany(ctx, filter, deleteOpts)
}

// ORM-like convenience methods
func (c *Collection) Create(document interface{}) error {
	_, err := c.InsertOne(document)
	return err
}

func (c *Collection) First(result interface{}, filter ...interface{}) error {
	var f interface{} = bson.M{}
	if len(filter) > 0 {
		f = filter[0]
	}
	return c.FindOne(f, result)
}

func (c *Collection) Where(field string, value interface{}) contracts.QueryBuilder {
	return NewQueryBuilder(c).Where(field, value)
}

// Collection management
func (c *Collection) Drop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return c.collection.Drop(ctx)
}

func (c *Collection) Name() string {
	return c.collection.Name()
}

func (c *Collection) CountDocuments(filter interface{}, opts ...interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var countOpts *options.CountOptions
	if len(opts) > 0 {
		if opt, ok := opts[0].(*options.CountOptions); ok {
			countOpts = opt
		}
	}

	return c.collection.CountDocuments(ctx, filter, countOpts)
}
