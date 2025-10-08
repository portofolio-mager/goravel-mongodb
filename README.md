# MongoDB Driver for Goravel

The MongoDB driver for Goravel framework, providing native MongoDB operations with direct driver access (bypassing GORM).

## Version

| goravel/mongodb | goravel/framework |
|-----------------|-------------------|
| v1.0.*          | v1.16.*           |

## Features

- **Direct MongoDB Access**: Bypass GORM for native MongoDB operations
- **Dual Interface**: Both native MongoDB client and ORM-like convenience methods
- **Full MongoDB Support**: Aggregation pipelines, indexes, transactions, etc.
- **Query Builder**: Familiar query building with MongoDB-native operations
- **Goravel Integration**: Works seamlessly with Goravel's service container and configuration

## Install

Run the command below in your project to install the package automatically:

```bash
./artisan package:install github.com/portofolio-mager/goravel-mongodb
```

Or check [the setup file](./setup/setup.go) to install the package manually.

## Configuration

Add MongoDB connection to your `config/database.go`:

```go
"mongodb": map[string]any{
    "uri":      config.Env("MONGODB_URI", "mongodb://localhost:27017"),
    "database": config.Env("MONGODB_DATABASE", "goravel"),
    "username": config.Env("MONGODB_USERNAME", ""),
    "password": config.Env("MONGODB_PASSWORD", ""),
    "auth_source": config.Env("MONGODB_AUTH_SOURCE", "admin"),
    "options": map[string]any{
        "max_pool_size": 100,
        "min_pool_size": 5,
        "connect_timeout": 10,
        "server_timeout": 30,
    },
    "via": func() (driver.Driver, error) {
        return mongodbfacades.MongoDBDriver("mongodb")
    },
}
```

## Usage

### Direct MongoDB Operations

```go
import "github.com/portofolio-mager/goravel-mongodb/facades"

// Get MongoDB client directly
client, _ := facades.MongoDB("mongodb")
database := client.Database("myapp")
collection := database.Collection("users")

// Native MongoDB operations
user := &User{}
err := collection.FindOne(bson.M{"email": "john@example.com"}, user)

// Insert document
result, err := collection.InsertOne(&User{
    Name:  "John Doe",
    Email: "john@example.com",
})

// Complex queries with MongoDB syntax
cursor, err := collection.Find(bson.M{
    "age": bson.M{"$gte": 18},
    "status": bson.M{"$in": []string{"active", "pending"}},
})
```

### Convenience Methods (ORM-like)

```go
// Get collection directly
collection, _ := facades.Collection("users")

// Simple queries
var users []User
err := collection.Where("status", "active").Limit(10).Find(&users)

// Single document
var user User
err := collection.Where("email", "john@example.com").First(&user)

// Create
err := collection.Create(&User{Name: "Jane", Email: "jane@example.com"})

// Query builder
users := []User{}
err := collection.
    Where("age", bson.M{"$gte": 18}).
    WhereIn("status", []interface{"active", "verified"}).
    Sort("created_at", -1).
    Limit(10).
    Find(&users)
```

### Native Facade Helpers

```go
import (
    // If you also import Goravel core facades, consider aliasing:
    // mongodbfacades "github.com/portofolio-mager/goravel-mongodb/facades"
    "github.com/portofolio-mager/goravel-mongodb/facades"
)

// Native MongoDB client (*mongo.Client)
nativeClient, err := facades.NativeClient()                    // uses connection "mongodb"
nativeClientAlt, err := facades.NativeClient("analytics")     // specify connection

// Native MongoDB database (*mongo.Database)
nativeDB, err := facades.NativeDatabase("myapp")              // default connection
nativeDBAlt, err := facades.NativeDatabase("myapp", "analytics")

// Native MongoDB collection (*mongo.Collection)
nativeCol, err := facades.NativeCollection("users")           // uses configured default database
nativeColAlt, err := facades.NativeCollection("users", "other_db")
```

### Advanced Features

```go
// Access native MongoDB client for advanced operations
client := collection.Native().Database().Client()

// Aggregation pipelines
pipeline := []bson.M{
    {"$match": bson.M{"status": "active"}},
    {"$group": bson.M{"_id": "$department", "count": bson.M{"$sum": 1}}},
    {"$sort": bson.M{"count": -1}},
}
cursor, err := collection.Native().Aggregate(context.TODO(), pipeline)

// Transactions (MongoDB 4.0+ with replica sets)
session, err := client.StartSession()
defer session.EndSession(context.TODO())

err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {
    // Transactional operations here
    return nil
})
```

## Query Builder Methods

### Where Conditions
- `Where(field, value)` - Exact match
- `WhereIn(field, values)` - Match any value in array
- `WhereNotIn(field, values)` - Don't match any value in array
- `WhereExists(field)` - Field exists
- `WhereNotExists(field)` - Field doesn't exist
- `WhereGt(field, value)` - Greater than
- `WhereGte(field, value)` - Greater than or equal
- `WhereLt(field, value)` - Less than
- `WhereLte(field, value)` - Less than or equal
- `WhereNe(field, value)` - Not equal
- `WhereRegex(field, pattern, options...)` - Regular expression

### Query Modifiers
- `Limit(limit)` - Limit results
- `Skip(skip)` - Skip results
- `Sort(field, order)` - Sort by field (1 = ascending, -1 = descending)
- `Select(fields...)` - Select specific fields

### Result Methods
- `Find(results)` - Find multiple documents
- `First(result)` - Find first document
- `Count()` - Count documents

## Environment Variables

```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=goravel
MONGODB_USERNAME=
MONGODB_PASSWORD=
MONGODB_AUTH_SOURCE=admin
```

## Testing

Run command below to run test:

```bash
go test ./...
```

## Docker Support

The package includes Docker support for testing:

```yaml
# docker-compose.yml
services:
  mongodb:
    image: mongo:7.0
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: 123123
      MONGO_INITDB_DATABASE: goravel_test
```
