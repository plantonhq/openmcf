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
	// StackUpdateStreamingResponseCollectionName is the name of the MongoDB collection for streaming responses.
	StackUpdateStreamingResponseCollectionName = "stackupdate_streaming_responses"
)

// StackUpdateStreamingResponseRepository provides data access methods for streaming responses.
type StackUpdateStreamingResponseRepository struct {
	collection *mongo.Collection
}

// NewStackUpdateStreamingResponseRepository creates a new repository instance.
func NewStackUpdateStreamingResponseRepository(db *MongoDB) *StackUpdateStreamingResponseRepository {
	return &StackUpdateStreamingResponseRepository{
		collection: db.Database.Collection(StackUpdateStreamingResponseCollectionName),
	}
}

// Create inserts a new streaming response chunk into MongoDB.
func (r *StackUpdateStreamingResponseRepository) Create(ctx context.Context, response *models.StackUpdateStreamingResponse) (*models.StackUpdateStreamingResponse, error) {
	now := time.Now()
	response.ID = primitive.NewObjectID()
	response.CreatedAt = now

	result, err := r.collection.InsertOne(ctx, response)
	if err != nil {
		return nil, fmt.Errorf("failed to insert streaming response: %w", err)
	}

	response.ID = result.InsertedID.(primitive.ObjectID)
	return response, nil
}

// CreateBatch inserts multiple streaming response chunks in a single operation.
func (r *StackUpdateStreamingResponseRepository) CreateBatch(ctx context.Context, responses []*models.StackUpdateStreamingResponse) error {
	if len(responses) == 0 {
		return nil
	}

	now := time.Now()
	docs := make([]interface{}, len(responses))
	for i, resp := range responses {
		resp.ID = primitive.NewObjectID()
		resp.CreatedAt = now
		docs[i] = resp
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert streaming responses batch: %w", err)
	}

	return nil
}

// FindByStackUpdateID retrieves all streaming responses for a specific stack-update, ordered by sequence number.
func (r *StackUpdateStreamingResponseRepository) FindByStackUpdateID(ctx context.Context, stackUpdateID string) ([]*models.StackUpdateStreamingResponse, error) {
	filter := bson.M{"stack_update_id": stackUpdateID}
	// Use bson.D for ordered sort (sequence_num first, then created_at)
	opts := options.Find().SetSort(bson.D{
		{Key: "sequence_num", Value: 1},
		{Key: "created_at", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find streaming responses: %w", err)
	}
	defer cursor.Close(ctx)

	var responses []*models.StackUpdateStreamingResponse
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, fmt.Errorf("failed to decode streaming responses: %w", err)
	}

	return responses, nil
}

// DeleteByStackUpdateID deletes all streaming responses for a specific stack-update.
func (r *StackUpdateStreamingResponseRepository) DeleteByStackUpdateID(ctx context.Context, stackUpdateID string) error {
	filter := bson.M{"stack_update_id": stackUpdateID}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete streaming responses: %w", err)
	}

	_ = result // Result contains deleted count, but we don't need it for now
	return nil
}

// GetNextSequenceNum returns the next sequence number for a stack-update.
func (r *StackUpdateStreamingResponseRepository) GetNextSequenceNum(ctx context.Context, stackUpdateID string) (int, error) {
	filter := bson.M{"stack_update_id": stackUpdateID}
	opts := options.FindOne().SetSort(bson.M{"sequence_num": -1})

	var lastResponse models.StackUpdateStreamingResponse
	err := r.collection.FindOne(ctx, filter, opts).Decode(&lastResponse)
	if err == mongo.ErrNoDocuments {
		// No previous responses, start at 0
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get next sequence number: %w", err)
	}

	return lastResponse.SequenceNum + 1, nil
}

// FindByStackUpdateIDAfterSequence retrieves streaming responses for a stack-update after a specific sequence number.
// Used for resuming streams from a specific point.
func (r *StackUpdateStreamingResponseRepository) FindByStackUpdateIDAfterSequence(ctx context.Context, stackUpdateID string, afterSequenceNum int) ([]*models.StackUpdateStreamingResponse, error) {
	filter := bson.M{
		"stack_update_id": stackUpdateID,
		"sequence_num":    bson.M{"$gt": afterSequenceNum},
	}
	// Use bson.D for ordered sort (sequence_num first, then created_at)
	opts := options.Find().SetSort(bson.D{
		{Key: "sequence_num", Value: 1},
		{Key: "created_at", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find streaming responses: %w", err)
	}
	defer cursor.Close(ctx)

	var responses []*models.StackUpdateStreamingResponse
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, fmt.Errorf("failed to decode streaming responses: %w", err)
	}

	return responses, nil
}
