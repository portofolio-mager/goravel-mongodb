package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tonidy/goravel-mongodb/contracts"
)

var _ contracts.QueryBuilder = &QueryBuilder{}

type QueryBuilder struct {
	collection *Collection
	filter     bson.M
	options    *options.FindOptions
	projection bson.M
}

func NewQueryBuilder(collection *Collection) *QueryBuilder {
	return &QueryBuilder{
		collection: collection,
		filter:     bson.M{},
		options:    options.Find(),
		projection: bson.M{},
	}
}

// Where conditions
func (q *QueryBuilder) Where(field string, value interface{}) contracts.QueryBuilder {
	q.filter[field] = value
	return q
}

func (q *QueryBuilder) WhereIn(field string, values []interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$in": values}
	return q
}

func (q *QueryBuilder) WhereNotIn(field string, values []interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$nin": values}
	return q
}

func (q *QueryBuilder) WhereExists(field string) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$exists": true}
	return q
}

func (q *QueryBuilder) WhereNotExists(field string) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$exists": false}
	return q
}

func (q *QueryBuilder) WhereGt(field string, value interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$gt": value}
	return q
}

func (q *QueryBuilder) WhereGte(field string, value interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$gte": value}
	return q
}

func (q *QueryBuilder) WhereLt(field string, value interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$lt": value}
	return q
}

func (q *QueryBuilder) WhereLte(field string, value interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$lte": value}
	return q
}

func (q *QueryBuilder) WhereNe(field string, value interface{}) contracts.QueryBuilder {
	q.filter[field] = bson.M{"$ne": value}
	return q
}

func (q *QueryBuilder) WhereRegex(field string, pattern string, options ...string) contracts.QueryBuilder {
	regexFilter := bson.M{"$regex": primitive.Regex{Pattern: pattern}}
	if len(options) > 0 {
		regexFilter["$options"] = options[0]
	}
	q.filter[field] = regexFilter
	return q
}

// Query modifiers
func (q *QueryBuilder) Limit(limit int64) contracts.QueryBuilder {
	q.options.SetLimit(limit)
	return q
}

func (q *QueryBuilder) Skip(skip int64) contracts.QueryBuilder {
	q.options.SetSkip(skip)
	return q
}

func (q *QueryBuilder) Sort(field string, order int) contracts.QueryBuilder {
	sort := q.options.Sort
	if sort == nil {
		sort = bson.D{}
	}

	// Convert to bson.D if needed
	var sortDoc bson.D
	if sortSlice, ok := sort.(bson.D); ok {
		sortDoc = sortSlice
	} else {
		sortDoc = bson.D{}
	}

	// Add or update the field
	found := false
	for i, elem := range sortDoc {
		if elem.Key == field {
			sortDoc[i].Value = order
			found = true
			break
		}
	}
	if !found {
		sortDoc = append(sortDoc, bson.E{Key: field, Value: order})
	}

	q.options.SetSort(sortDoc)
	return q
}

func (q *QueryBuilder) Select(fields ...string) contracts.QueryBuilder {
	if len(fields) > 0 {
		for _, field := range fields {
			q.projection[field] = 1
		}
		q.options.SetProjection(q.projection)
	}
	return q
}

// Result methods
func (q *QueryBuilder) Find(results interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := q.collection.collection.Find(ctx, q.filter, q.options)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

func (q *QueryBuilder) First(result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOneOpts := options.FindOne()
	if len(q.projection) > 0 {
		findOneOpts.SetProjection(q.projection)
	}
	if q.options.Sort != nil {
		findOneOpts.SetSort(q.options.Sort)
	}
	if q.options.Skip != nil {
		findOneOpts.SetSkip(*q.options.Skip)
	}

	return q.collection.collection.FindOne(ctx, q.filter, findOneOpts).Decode(result)
}

func (q *QueryBuilder) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	countOpts := options.Count()
	if q.options.Skip != nil {
		countOpts.SetSkip(*q.options.Skip)
	}
	if q.options.Limit != nil {
		countOpts.SetLimit(*q.options.Limit)
	}

	count, err := q.collection.collection.CountDocuments(ctx, q.filter, countOpts)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return count, nil
}
