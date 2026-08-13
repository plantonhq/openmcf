package aa_e2e

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3vectors"
	s3vectorstypes "github.com/aws/aws-sdk-go-v2/service/s3vectors/types"
	"github.com/pkg/errors"
)

// Standing S3 Vectors fixture for knowledge-base lanes. A Bedrock
// knowledge base of the VECTOR type needs a vector store, and S3 Vectors
// is the only self-contained one at the pin: the OpenSearch Serverless
// path needs a vector index no AWS-provider resource can create (upstream
// tests that path with the third-party opensearch provider), and no
// catalog kind wraps aws_s3vectors_* yet (a planned AwsS3VectorBucket
// kind carries the ledger judgment). Until that kind exists, the harness
// ensures the fixture through the SDK -- exactly how the upstream
// provider's own acceptance tests provision it -- and exports its index
// ARN as a PLANTON_E2E_ environment token for scenario manifests.
//
// The fixture is a STANDING account fixture (the seeded-S3-bucket class):
// S3 Vectors bills per use, so an empty bucket+index costs nothing at
// rest, and reusing one across runs avoids create/delete churn. It is
// intentionally NOT torn down per run; the account janitor policy governs
// it like the other seeded fixtures.
const (
	// S3VectorsIndexArnEnvVar is the ${E2E_ENV:...} token variable the
	// knowledge-base scenarios reference for the fixture index's ARN.
	S3VectorsIndexArnEnvVar = "PLANTON_E2E_S3VECTORS_INDEX_ARN"

	s3VectorsFixtureBucket = "planton-e2e-s3vectors"
	s3VectorsFixtureIndex  = "planton-e2e-kb-index"

	// The index shape pairs with Titan Text Embeddings V2 at 256
	// dimensions (the scenario pins dimensions: 256) -- the same shape the
	// upstream provider's own S3 Vectors acceptance fixture uses.
	s3VectorsFixtureDimension = 256
)

// EnsureS3VectorsKnowledgeBaseFixture idempotently creates the standing
// S3 Vectors bucket+index and exports S3VectorsIndexArnEnvVar for
// scenario token expansion. Call it from a component's test entrypoint
// BEFORE running scenarios that reference the token.
func EnsureS3VectorsKnowledgeBaseFixture(ctx context.Context) error {
	if os.Getenv(S3VectorsIndexArnEnvVar) != "" {
		return nil
	}

	region := firstNonEmpty(os.Getenv("E2E_AWS_REGION"), os.Getenv("AWS_REGION"), defaultRegion)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return errors.Wrap(err, "failed to load AWS config for the S3 Vectors fixture")
	}
	client := s3vectors.NewFromConfig(cfg)

	if _, err := client.CreateVectorBucket(ctx, &s3vectors.CreateVectorBucketInput{
		VectorBucketName: aws.String(s3VectorsFixtureBucket),
	}); err != nil && !isS3VectorsConflict(err) {
		return errors.Wrapf(err, "create vector bucket %s", s3VectorsFixtureBucket)
	}

	if _, err := client.CreateIndex(ctx, &s3vectors.CreateIndexInput{
		VectorBucketName: aws.String(s3VectorsFixtureBucket),
		IndexName:        aws.String(s3VectorsFixtureIndex),
		DataType:         s3vectorstypes.DataTypeFloat32,
		Dimension:        aws.Int32(s3VectorsFixtureDimension),
		DistanceMetric:   s3vectorstypes.DistanceMetricEuclidean,
	}); err != nil && !isS3VectorsConflict(err) {
		return errors.Wrapf(err, "create vector index %s", s3VectorsFixtureIndex)
	}

	index, err := client.GetIndex(ctx, &s3vectors.GetIndexInput{
		VectorBucketName: aws.String(s3VectorsFixtureBucket),
		IndexName:        aws.String(s3VectorsFixtureIndex),
	})
	if err != nil {
		return errors.Wrapf(err, "get vector index %s", s3VectorsFixtureIndex)
	}

	arn := aws.ToString(index.Index.IndexArn)
	if err := os.Setenv(S3VectorsIndexArnEnvVar, arn); err != nil {
		return errors.Wrap(err, "export the fixture index ARN")
	}
	fmt.Printf("  [aws] S3 Vectors knowledge-base fixture ready: %s\n", arn)
	return nil
}

// isS3VectorsConflict matches the already-exists class of both fixture
// creates (idempotent ensure).
func isS3VectorsConflict(err error) bool {
	var conflict *s3vectorstypes.ConflictException
	return errors.As(err, &conflict) || strings.Contains(err.Error(), "ConflictException")
}
