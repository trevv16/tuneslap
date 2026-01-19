package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tuneslap/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides generic database operations with caching
type Repository[T any] struct {
	collectionName string
	timeout        time.Duration
	cacheTTL       time.Duration
}

// NewRepository creates a new repository instance
func NewRepository[T any](collectionName string) *Repository[T] {
	return &Repository[T]{
		collectionName: collectionName,
		timeout:        10 * time.Second,
		cacheTTL:       5 * time.Minute,
	}
}

// getCollection returns the MongoDB collection
func (r *Repository[T]) getCollection() *mongo.Collection {
	return database.GetCollection(r.collectionName)
}

// getContext returns a context with timeout
func (r *Repository[T]) getContext() (context.Context, context.CancelFunc) {
	return database.GetContext(r.timeout)
}

// getCacheKey generates a cache key for the given operation and parameters
func (r *Repository[T]) getCacheKey(operation string, params ...interface{}) string {
	key := fmt.Sprintf("%s:%s", r.collectionName, operation)
	for _, param := range params {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}

// getFromCache retrieves data from cache
func (r *Repository[T]) getFromCache(key string) (T, bool) {
	data, err := database.GetCache(key)
	if err != nil {
		var zero T
		return zero, false
	}

	var result T
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		var zero T
		return zero, false
	}

	return result, true
}

// setCache stores data in cache
func (r *Repository[T]) setCache(key string, data T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return database.SetCache(key, string(jsonData), r.cacheTTL)
}

// clearCache clears all cache for this collection
func (r *Repository[T]) clearCache() error {
	// Get Redis client for pattern-based deletion
	client := database.GetRedisClient()
	if client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	ctx := context.Background()
	pattern := fmt.Sprintf("%s:*", r.collectionName)

	// Use SCAN to find all keys matching the pattern
	var keys []string
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	// Delete all found keys
	if len(keys) > 0 {
		return client.Del(ctx, keys...).Err()
	}

	return nil
}

// FindAll retrieves all documents with optional filtering, pagination, and projection
func (r *Repository[T]) FindAll(filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	// Generate cache key based on filter and options
	filterBytes, _ := json.Marshal(filter)
	optsBytes, _ := json.Marshal(opts)
	findAllKey := r.getCacheKey("findAllResults", string(filterBytes), string(optsBytes))

	// Try to get from cache first
	if findAllStr, err := database.GetCache(findAllKey); err == nil {
		var results []T
		if err := json.Unmarshal([]byte(findAllStr), &results); err == nil {
			return results, nil
		}
	}

	// If not in cache, query database
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		// If error is "no documents", treat as empty result (defensive check)
		// FindAll shouldn't return this error, but handle it gracefully
		if strings.Contains(err.Error(), "no documents") {
			return []T{}, nil
		}
		return nil, err
	}

	// Ensure results is not nil (defensive check)
	if results == nil {
		results = []T{}
	}

	// Store results in cache (even empty results to avoid repeated queries)
	if resultsBytes, err := json.Marshal(results); err == nil {
		database.SetCache(findAllKey, string(resultsBytes), r.cacheTTL)
	}

	return results, nil
}

// FindOne retrieves a single document by filter with optional projection
func (r *Repository[T]) FindOne(filter bson.M, opts ...*options.FindOneOptions) (T, error) {
	// Generate cache key based on filter and options
	filterBytes, _ := json.Marshal(filter)
	optsBytes, _ := json.Marshal(opts)
	cacheKey := r.getCacheKey("findOne", string(filterBytes), string(optsBytes))

	// Try to get from cache first
	if cached, found := r.getFromCache(cacheKey); found {
		return cached, nil
	}

	// If not in cache, query database
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	var result T
	if len(opts) > 0 {
		err := collection.FindOne(ctx, filter, opts...).Decode(&result)
		if err != nil {
			return result, err
		}
	} else {
		err := collection.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			return result, err
		}
	}

	// Store in cache
	r.setCache(cacheKey, result)

	return result, nil
}

// FindByID retrieves a document by its ID with optional projection
func (r *Repository[T]) FindByID(id primitive.ObjectID, opts ...*options.FindOneOptions) (T, error) {
	return r.FindOne(bson.M{"_id": id}, opts...)
}

// Create inserts a new document
func (r *Repository[T]) Create(document T) (T, error) {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		var zero T
		return zero, err
	}

	// Clear cache after creation
	r.clearCache()

	// Fetch the created document
	return r.FindByID(result.InsertedID.(primitive.ObjectID))
}

// Update updates a document by ID
func (r *Repository[T]) Update(id primitive.ObjectID, update bson.M) (T, error) {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Add updatedAt timestamp to existing $set operation
	if setMap, ok := update["$set"].(bson.M); ok {
		setMap["updatedAt"] = time.Now()
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		var zero T
		return zero, err
	}

	if result.MatchedCount == 0 {
		var zero T
		return zero, mongo.ErrNoDocuments
	}

	// Clear cache after update
	r.clearCache()

	return r.FindByID(id)
}

// Delete removes a document by ID
func (r *Repository[T]) Delete(id primitive.ObjectID) error {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Clear cache after deletion
	r.clearCache()

	return nil
}

// Count returns the number of documents matching the filter
func (r *Repository[T]) Count(filter bson.M) (int64, error) {
	// Generate cache key based on filter
	filterBytes, _ := json.Marshal(filter)
	cacheKey := r.getCacheKey("count", string(filterBytes))

	// Try to get from cache first
	if _, found := r.getFromCache(cacheKey); found {
		// For count, we need to handle it differently since we're storing a single item
		// We'll use a special count cache key format
		countKey := r.getCacheKey("countValue", string(filterBytes))
		if countStr, err := database.GetCache(countKey); err == nil {
			if count, err := strconv.ParseInt(countStr, 10, 64); err == nil {
				return count, nil
			}
		}
	}

	// If not in cache, query database
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	// Store count in cache with a special key
	countKey := r.getCacheKey("countValue", string(filterBytes))
	database.SetCache(countKey, fmt.Sprintf("%d", count), r.cacheTTL)

	return count, nil
}

// Aggregate performs MongoDB aggregation pipeline
func (r *Repository[T]) Aggregate(pipeline interface{}) ([]T, error) {
	// Generate cache key based on pipeline
	pipelineBytes, _ := json.Marshal(pipeline)
	cacheKey := r.getCacheKey("aggregate", string(pipelineBytes))

	// Try to get from cache first
	if _, found := r.getFromCache(cacheKey); found {
		// For aggregation, we need to handle multiple results
		// We'll use a special aggregation cache key format
		aggKey := r.getCacheKey("aggregateResults", string(pipelineBytes))
		if aggStr, err := database.GetCache(aggKey); err == nil {
			var results []T
			if err := json.Unmarshal([]byte(aggStr), &results); err == nil {
				return results, nil
			}
		}
	}

	// If not in cache, query database
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Store results in cache with a special key
	if len(results) > 0 {
		aggKey := r.getCacheKey("aggregateResults", string(pipelineBytes))
		if resultsBytes, err := json.Marshal(results); err == nil {
			database.SetCache(aggKey, string(resultsBytes), r.cacheTTL)
		}
	}

	return results, nil
}

// AggregateWithCount performs aggregation and returns results with total count
func (r *Repository[T]) AggregateWithCount(pipeline interface{}, countFilter bson.M) ([]T, int64, error) {
	results, err := r.Aggregate(pipeline)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := r.Count(countFilter)
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

// UpdateNestedArray updates a specific element in a nested array
func (r *Repository[T]) UpdateNestedArray(
	parentID primitive.ObjectID,
	arrayField string,
	elementID primitive.ObjectID,
	update bson.M,
) error {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Build the update with positional operator
	updateFields := bson.M{}
	for key, value := range update {
		updateFields[arrayField+".$."+key] = value
	}
	updateFields["updatedAt"] = time.Now()

	filter := bson.M{
		"_id":               parentID,
		arrayField + "._id": elementID,
	}

	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateFields})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Clear cache after update
	r.clearCache()

	return nil
}

// DeleteFromNestedArray removes an element from a nested array
func (r *Repository[T]) DeleteFromNestedArray(
	parentID primitive.ObjectID,
	arrayField string,
	elementID primitive.ObjectID,
) error {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	filter := bson.M{"_id": parentID}
	update := bson.M{
		"$pull": bson.M{
			arrayField: bson.M{"_id": elementID},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Clear cache after deletion
	r.clearCache()

	return nil
}

// AddToNestedArray adds an element to a nested array
func (r *Repository[T]) AddToNestedArray(
	parentID primitive.ObjectID,
	arrayField string,
	element interface{},
) error {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	filter := bson.M{"_id": parentID}
	update := bson.M{
		"$push": bson.M{
			arrayField: element,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Clear cache after addition
	r.clearCache()

	return nil
}

// ArrayRepository provides methods for handling embedded arrays in documents
type ArrayRepository[T any] struct {
	*Repository[T]
}

// NewArrayRepository creates a new array repository instance
func NewArrayRepository[T any](collectionName string) *ArrayRepository[T] {
	return &ArrayRepository[T]{
		Repository: NewRepository[T](collectionName),
	}
}

// CreateInArray adds a new element to an embedded array field
func (r *ArrayRepository[T]) CreateInArray(
	parentID primitive.ObjectID,
	arrayField string,
	element T,
) (T, error) {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	filter := bson.M{
		"_id": parentID,
	}

	// First, check if the array field exists and is null, initialize it if needed
	var parentDoc bson.M
	err := collection.FindOne(ctx, filter).Decode(&parentDoc)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("parent document not found")
	}

	// If the array field is null or doesn't exist, initialize it as an empty array
	if parentDoc[arrayField] == nil {
		initUpdate := bson.M{
			"$set": bson.M{
				arrayField: []interface{}{},
			},
		}
		_, err = collection.UpdateOne(ctx, filter, initUpdate)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to initialize array field: %w", err)
		}
	}

	// Add the element to the array
	update := bson.M{
		"$push": bson.M{
			arrayField: element,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to add element to array: %w", err)
	}

	if result.MatchedCount == 0 {
		var zero T
		return zero, fmt.Errorf("parent document not found")
	}

	// Clear cache for this collection since we modified the document
	r.clearCache()

	return element, nil
}

// UpdateInArray updates a specific element in an embedded array field
func (r *ArrayRepository[T]) UpdateInArray(
	parentID primitive.ObjectID,
	arrayField string,
	elementID primitive.ObjectID,
	update bson.M,
) error {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Build the update fields for the array element
	setFields := bson.M{}
	for key, value := range update {
		setFields[arrayField+".$."+key] = value
	}
	setFields["updatedAt"] = time.Now()

	// Update the specific element in the array using positional operator
	filter := bson.M{
		"_id":               parentID,
		arrayField + "._id": elementID,
	}

	updateDoc := bson.M{
		"$set": setFields,
	}

	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("failed to update element in array: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("parent document or element not found")
	}

	// Clear cache for this collection since we modified the document
	r.clearCache()

	return nil
}

// DeleteFromArray removes a specific element from an embedded array field
func (r *ArrayRepository[T]) DeleteFromArray(
	parentID primitive.ObjectID,
	arrayField string,
	elementID primitive.ObjectID,
) error {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Remove the element from the array
	update := bson.M{
		"$pull": bson.M{
			arrayField: bson.M{
				"_id": elementID,
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	filter := bson.M{
		"_id": parentID,
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove element from array: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("parent document not found")
	}

	// Clear cache for this collection since we modified the document
	r.clearCache()

	return nil
}

// GetArrayElement retrieves a specific element from an embedded array field
func (r *ArrayRepository[T]) GetArrayElement(
	parentID primitive.ObjectID,
	arrayField string,
	elementID primitive.ObjectID,
) (T, error) {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Use aggregation to extract the specific element from the array
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": parentID,
			},
		},
		{
			"$unwind": "$" + arrayField,
		},
		{
			"$match": bson.M{
				arrayField + "._id": elementID,
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$" + arrayField,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to get element from array: %w", err)
	}
	defer cursor.Close(ctx)

	var result T
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			var zero T
			return zero, fmt.Errorf("failed to decode element: %w", err)
		}
		return result, nil
	}

	var zero T
	return zero, fmt.Errorf("element not found")
}

// GetAllArrayElements retrieves all elements from an embedded array field
func (r *ArrayRepository[T]) GetAllArrayElements(
	parentID primitive.ObjectID,
	arrayField string,
) ([]T, error) {
	collection := r.getCollection()
	ctx, cancel := r.getContext()
	defer cancel()

	// Use aggregation to extract all elements from the array
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": parentID,
			},
		},
		{
			"$unwind": "$" + arrayField,
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$" + arrayField,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get elements from array: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode elements: %w", err)
	}

	return results, nil
}
