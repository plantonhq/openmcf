package database

import (
	"context"
	"fmt"
	"time"

	"github.com/plantonhq/openmcf/app/backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// StackUpdateCollectionName is the name of the MongoDB collection for stack-updates.
	StackUpdateCollectionName = "stackupdates"
)

// StackUpdateRepository provides data access methods for stack-updates.
type StackUpdateRepository struct {
	collection *mongo.Collection
}

// NewStackUpdateRepository creates a new repository instance.
func NewStackUpdateRepository(db *MongoDB) *StackUpdateRepository {
	return &StackUpdateRepository{
		collection: db.Database.Collection(StackUpdateCollectionName),
	}
}

// Create inserts a new stack-update into MongoDB.
func (r *StackUpdateRepository) Create(ctx context.Context, stackUpdate *models.StackUpdate) (*models.StackUpdate, error) {
	now := time.Now()
	stackUpdate.ID = primitive.NewObjectID()
	stackUpdate.CreatedAt = now
	stackUpdate.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, stackUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to insert stack-update: %w", err)
	}

	stackUpdate.ID = result.InsertedID.(primitive.ObjectID)
	return stackUpdate, nil
}

// FindByID retrieves a stack-update by ID.
func (r *StackUpdateRepository) FindByID(ctx context.Context, id string) (*models.StackUpdate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var stackUpdate models.StackUpdate
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&stackUpdate)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, but not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query stack-update by ID: %w", err)
	}
	return &stackUpdate, nil
}

// FindByCloudResourceID retrieves stack-updates by cloud resource ID, sorted by created_at descending.
func (r *StackUpdateRepository) FindByCloudResourceID(ctx context.Context, cloudResourceID string) ([]*models.StackUpdate, error) {
	filter := bson.M{"cloud_resource_id": cloudResourceID}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Newest first

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to query stack-updates: %w", err)
	}
	defer cursor.Close(ctx)

	var stackUpdates []*models.StackUpdate
	if err := cursor.All(ctx, &stackUpdates); err != nil {
		return nil, fmt.Errorf("failed to decode stack-updates: %w", err)
	}

	return stackUpdates, nil
}

// Update updates an existing stack-update in MongoDB.
func (r *StackUpdateRepository) Update(ctx context.Context, id string, stackUpdate *models.StackUpdate) (*models.StackUpdate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	stackUpdate.ID = objectID
	stackUpdate.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{},
	}

	// Only update fields that are provided (non-empty for strings)
	if stackUpdate.Status != "" {
		update["$set"].(bson.M)["status"] = stackUpdate.Status
	}
	if stackUpdate.Output != "" {
		update["$set"].(bson.M)["output"] = stackUpdate.Output
	}
	update["$set"].(bson.M)["updated_at"] = stackUpdate.UpdatedAt

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedStackUpdate models.StackUpdate
	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		opts,
	).Decode(&updatedStackUpdate)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("stack-update not found")
		}
		return nil, fmt.Errorf("failed to update stack-update: %w", err)
	}

	return &updatedStackUpdate, nil
}

// StackUpdateListOptions contains options for listing stack-updates.
type StackUpdateListOptions struct {
	CloudResourceID *string
	Status          *string
	PageNum         *int32
	PageSize        *int32
}

// List retrieves stack-updates with optional filters and pagination.
func (r *StackUpdateRepository) List(ctx context.Context, opts *StackUpdateListOptions) ([]*models.StackUpdate, error) {
	filter := bson.M{}

	if opts != nil {
		if opts.CloudResourceID != nil && *opts.CloudResourceID != "" {
			filter["cloud_resource_id"] = *opts.CloudResourceID
		}

		if opts.Status != nil && *opts.Status != "" {
			filter["status"] = *opts.Status
		}
	}

	findOptions := options.Find()

	// Apply pagination if provided
	if opts != nil && opts.PageNum != nil && opts.PageSize != nil {
		pageNum := *opts.PageNum
		pageSize := *opts.PageSize
		if pageSize > 0 {
			skip := int64(pageNum) * int64(pageSize)
			findOptions.SetSkip(skip)
			findOptions.SetLimit(int64(pageSize))
		}
	}

	// Sort by created_at descending (newest first)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to query stack-updates: %w", err)
	}
	defer cursor.Close(ctx)

	var stackUpdates []*models.StackUpdate
	if err := cursor.All(ctx, &stackUpdates); err != nil {
		return nil, fmt.Errorf("failed to decode stack-updates: %w", err)
	}

	return stackUpdates, nil
}

// Count returns the total count of stack-updates with optional filters.
func (r *StackUpdateRepository) Count(ctx context.Context, opts *StackUpdateListOptions) (int64, error) {
	filter := bson.M{}

	if opts != nil {
		if opts.CloudResourceID != nil && *opts.CloudResourceID != "" {
			filter["cloud_resource_id"] = *opts.CloudResourceID
		}

		if opts.Status != nil && *opts.Status != "" {
			filter["status"] = *opts.Status
		}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count stack-updates: %w", err)
	}

	return count, nil
}
