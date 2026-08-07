package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	pkgerrors "github.com/pkg/errors"
)

// dynamodbTableVerifier verifies an AwsDynamodb table via DescribeTable,
// keyed on the table_name output. A table mid-deletion stays describable
// with a DELETING status before the service starts returning the typed
// ResourceNotFoundException -- the same lifecycle class as the RDS-shaped
// kinds -- so existence is "described AND not deleting", and absence
// accepts either signal.
type dynamodbTableVerifier struct{}

func (*dynamodbTableVerifier) IDOutputKey() string { return "table_name" }

func (*dynamodbTableVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := dynamodbTableExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsdynamodb verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsdynamodb table %q not found after deploy", id)
	}
	return nil
}

func (*dynamodbTableVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := dynamodbTableExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsdynamodb verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsdynamodb table %q still exists after destroy", id)
	}
	return nil
}

// dynamodbTableExists reports whether the table is present and not
// already on its way out. A ResourceNotFoundException is treated as
// absent.
func dynamodbTableExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.Table == nil {
		return false, nil
	}
	if out.Table.TableStatus == types.TableStatusDeleting {
		return false, nil
	}
	return true, nil
}
