package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3tables"
	s3tablestypes "github.com/aws/aws-sdk-go-v2/service/s3tables/types"
	"github.com/aws/aws-sdk-go-v2/service/s3vectors"
	s3vectorstypes "github.com/aws/aws-sdk-go-v2/service/s3vectors/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// --- AwsS3DirectoryBucket -------------------------------------------------------

// s3DirectoryBucketVerifier verifies an AwsS3DirectoryBucket via
// HeadBucket on the module-derived full name
// ("{base}--{zone_id}--x-s3"), keyed on bucket_name. The standard S3
// client routes directory-bucket names transparently.
type s3DirectoryBucketVerifier struct{}

func (*s3DirectoryBucketVerifier) IDOutputKey() string { return "bucket_name" }

func (*s3DirectoryBucketVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := directoryBucketExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awss3directorybucket verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awss3directorybucket %q not found after deploy", id)
	}
	return nil
}

func (*s3DirectoryBucketVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := directoryBucketExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awss3directorybucket verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awss3directorybucket %q still exists after destroy", id)
	}
	return nil
}

func directoryBucketExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(id)})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "NotFound", "NoSuchBucket":
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

// --- AwsS3TableBucket -----------------------------------------------------------

// s3TableBucketVerifier verifies an AwsS3TableBucket via
// GetTableBucket, keyed on table_bucket_arn. The namespaces and
// tables are children of the bucket (force_destroy drains them) and
// cannot outlive it.
type s3TableBucketVerifier struct{}

func (*s3TableBucketVerifier) IDOutputKey() string { return "table_bucket_arn" }

func (*s3TableBucketVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := tableBucketExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awss3tablebucket verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awss3tablebucket %q not found after deploy", id)
	}
	return nil
}

func (*s3TableBucketVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := tableBucketExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awss3tablebucket verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awss3tablebucket %q still exists after destroy", id)
	}
	return nil
}

func tableBucketExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := s3tables.NewFromConfig(cfg, func(o *s3tables.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetTableBucket(ctx, &s3tables.GetTableBucketInput{
		TableBucketARN: aws.String(id),
	})
	if err != nil {
		var notFound *s3tablestypes.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// --- AwsS3VectorBucket ----------------------------------------------------------

// s3VectorBucketVerifier verifies an AwsS3VectorBucket via
// GetVectorBucket, keyed on vector_bucket_arn. The indexes are
// children of the bucket (force_destroy drains them) and cannot
// outlive it.
type s3VectorBucketVerifier struct{}

func (*s3VectorBucketVerifier) IDOutputKey() string { return "vector_bucket_arn" }

func (*s3VectorBucketVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := vectorBucketExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awss3vectorbucket verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awss3vectorbucket %q not found after deploy", id)
	}
	return nil
}

func (*s3VectorBucketVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := vectorBucketExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awss3vectorbucket verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awss3vectorbucket %q still exists after destroy", id)
	}
	return nil
}

func vectorBucketExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := s3vectors.NewFromConfig(cfg, func(o *s3vectors.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetVectorBucket(ctx, &s3vectors.GetVectorBucketInput{
		VectorBucketArn: aws.String(id),
	})
	if err != nil {
		var notFound *s3vectorstypes.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
