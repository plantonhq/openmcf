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
	// CloudResourceCollectionName is the name of the MongoDB collection for cloud resources.
	CloudResourceCollectionName = "cloud_resources"
)

// CloudResourceRepository provides data access methods for cloud resources.
type CloudResourceRepository struct {
	collection *mongo.Collection
}

// NewCloudResourceRepository creates a new repository instance.
func NewCloudResourceRepository(db *MongoDB) *CloudResourceRepository {
	return &CloudResourceRepository{
		collection: db.Database.Collection(CloudResourceCollectionName),
	}
}

// FindByName retrieves a cloud resource by name.
func (r *CloudResourceRepository) FindByName(ctx context.Context, name string) (*models.CloudResource, error) {
	var resource models.CloudResource
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&resource)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, but not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query cloud resource by name: %w", err)
	}
	return &resource, nil
}

// FindByNameAndKind retrieves a cloud resource by name and kind.
func (r *CloudResourceRepository) FindByNameAndKind(ctx context.Context, name string, kind string) (*models.CloudResource, error) {
	var resource models.CloudResource
	err := r.collection.FindOne(ctx, bson.M{"name": name, "kind": kind}).Decode(&resource)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, but not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query cloud resource by name and kind: %w", err)
	}
	return &resource, nil
}

// FindByID retrieves a cloud resource by ID.
func (r *CloudResourceRepository) FindByID(ctx context.Context, id string) (*models.CloudResource, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var resource models.CloudResource
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&resource)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Not found, but not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query cloud resource by ID: %w", err)
	}
	return &resource, nil
}

// Create inserts a new cloud resource into MongoDB.
func (r *CloudResourceRepository) Create(ctx context.Context, resource *models.CloudResource) (*models.CloudResource, error) {
	now := time.Now()
	resource.ID = primitive.NewObjectID()
	resource.CreatedAt = now
	resource.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to insert cloud resource: %w", err)
	}

	resource.ID = result.InsertedID.(primitive.ObjectID)
	return resource, nil
}

// CloudResourceListOptions contains options for listing cloud resources.
type CloudResourceListOptions struct {
	Kind     *string
	PageNum  *int32
	PageSize *int32
}

// List retrieves cloud resources from MongoDB with optional filters and pagination.
func (r *CloudResourceRepository) List(ctx context.Context, opts *CloudResourceListOptions) ([]*models.CloudResource, error) {
	filter := bson.M{}

	if opts != nil {
		if opts.Kind != nil && *opts.Kind != "" {
			filter["kind"] = *opts.Kind
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
		return nil, fmt.Errorf("failed to query cloud resources: %w", err)
	}
	defer cursor.Close(ctx)

	var resources []*models.CloudResource
	if err := cursor.All(ctx, &resources); err != nil {
		return nil, fmt.Errorf("failed to decode cloud resources: %w", err)
	}

	return resources, nil
}

// Update updates an existing cloud resource in MongoDB.
func (r *CloudResourceRepository) Update(ctx context.Context, id string, resource *models.CloudResource) (*models.CloudResource, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	resource.ID = objectID
	resource.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       resource.Name,
			"kind":       resource.Kind,
			"manifest":   resource.Manifest,
			"updated_at": resource.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // Return the document after update

	var updatedResource models.CloudResource
	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		opts,
	).Decode(&updatedResource)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("cloud resource not found")
		}
		return nil, fmt.Errorf("failed to update cloud resource: %w", err)
	}

	return &updatedResource, nil
}

// Delete removes a cloud resource from MongoDB.
func (r *CloudResourceRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete cloud resource: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("cloud resource not found")
	}

	return nil
}

// Count returns the total count of cloud resources with optional filters.
func (r *CloudResourceRepository) Count(ctx context.Context, opts *CloudResourceListOptions) (int64, error) {
	filter := bson.M{}

	if opts != nil {
		if opts.Kind != nil && *opts.Kind != "" {
			filter["kind"] = *opts.Kind
		}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count cloud resources: %w", err)
	}

	return count, nil
}
