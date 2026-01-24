package repositories

import (
	"context"
	"fmt"
	"time"
	"tuneslap/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

// UpdateNestedArray updates a specific element in a nested array (on Repository)
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

// DeleteFromNestedArray removes an element from a nested array (on Repository)
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

// AddToNestedArray adds an element to a nested array (on Repository)
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

// getCollection returns the MongoDB collection (helper to access parent method)
func (r *ArrayRepository[T]) getCollection() *mongo.Collection {
	return database.GetCollection(r.collectionName)
}

// getContext returns a context with timeout (helper to access parent method)
func (r *ArrayRepository[T]) getContext() (context.Context, context.CancelFunc) {
	return database.GetContext(r.timeout)
}
