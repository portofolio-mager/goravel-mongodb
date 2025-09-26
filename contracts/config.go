package contracts

import (
	contractsconfig "github.com/goravel/framework/contracts/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigBuilder interface {
	Config() contractsconfig.Config
	Connection() string
	Writers() []FullConfig
}

// Config Used in config/database.go for MongoDB
type Config struct {
	URI            string                 `json:"uri"`
	Database       string                 `json:"database"`
	Username       string                 `json:"username"`
	Password       string                 `json:"password"`
	AuthSource     string                 `json:"auth_source"`
	ReplicaSet     string                 `json:"replica_set"`
	TLS            bool                   `json:"tls"`
	TLSCAFile      string                 `json:"tls_ca_file"`
	TLSCertFile    string                 `json:"tls_cert_file"`
	TLSKeyFile     string                 `json:"tls_key_file"`
	MaxPoolSize    *uint64                `json:"max_pool_size"`
	MinPoolSize    *uint64                `json:"min_pool_size"`
	ConnectTimeout *int                   `json:"connect_timeout"`
	ServerTimeout  *int                   `json:"server_timeout"`
	Options        map[string]interface{} `json:"options"`
}

// FullConfig Fill the default value for Config
type FullConfig struct {
	Config
	Connection string
	Driver     string
}

// Client represents the MongoDB client interface
type Client interface {
	// Native MongoDB client access
	Native() *mongo.Client

	// Database operations
	Database(name ...string) Database

	// Collection operations
	Collection(collection string, database ...string) Collection

	// Connection management
	Ping() error
	Close() error
}

// Database represents a MongoDB database interface
type Database interface {
	// Native database access
	Native() *mongo.Database

	// Collection operations
	Collection(name string) Collection

	// Database operations
	CreateCollection(name string, opts ...interface{}) error
	ListCollections() ([]string, error)
	Drop() error
	Name() string
}

// Collection represents a MongoDB collection interface
type Collection interface {
	// Native collection access
	Native() *mongo.Collection

	// Basic CRUD operations
	FindOne(filter interface{}, result interface{}, opts ...interface{}) error
	Find(filter interface{}, opts ...interface{}) (*mongo.Cursor, error)
	InsertOne(document interface{}, opts ...interface{}) (*mongo.InsertOneResult, error)
	InsertMany(documents []interface{}, opts ...interface{}) (*mongo.InsertManyResult, error)
	UpdateOne(filter interface{}, update interface{}, opts ...interface{}) (*mongo.UpdateResult, error)
	UpdateMany(filter interface{}, update interface{}, opts ...interface{}) (*mongo.UpdateResult, error)
	DeleteOne(filter interface{}, opts ...interface{}) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...interface{}) (*mongo.DeleteResult, error)

	// ORM-like convenience methods
	Create(document interface{}) error
	First(result interface{}, filter ...interface{}) error
	Where(field string, value interface{}) QueryBuilder

	// Collection management
	Drop() error
	Name() string
	CountDocuments(filter interface{}, opts ...interface{}) (int64, error)
}

// QueryBuilder represents a query builder interface for MongoDB
type QueryBuilder interface {
	Where(field string, value interface{}) QueryBuilder
	WhereIn(field string, values []interface{}) QueryBuilder
	WhereNotIn(field string, values []interface{}) QueryBuilder
	WhereExists(field string) QueryBuilder
	WhereNotExists(field string) QueryBuilder
	WhereGt(field string, value interface{}) QueryBuilder
	WhereGte(field string, value interface{}) QueryBuilder
	WhereLt(field string, value interface{}) QueryBuilder
	WhereLte(field string, value interface{}) QueryBuilder
	WhereNe(field string, value interface{}) QueryBuilder
	WhereRegex(field string, pattern string, options ...string) QueryBuilder

	// Result methods
	Find(results interface{}) error
	First(result interface{}) error
	Count() (int64, error)

	// Query modifiers
	Limit(limit int64) QueryBuilder
	Skip(skip int64) QueryBuilder
	Sort(field string, order int) QueryBuilder
	Select(fields ...string) QueryBuilder
}
