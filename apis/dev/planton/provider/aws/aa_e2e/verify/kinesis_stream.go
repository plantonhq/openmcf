package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pkgerrors "github.com/pkg/errors"
)

// kinesisStreamVerifier verifies an AwsKinesisStream via DescribeStreamSummary,
// keyed on the stream_name output. A stream mid-deletion stays describable
// with a DELETING status before the service starts returning the typed
// ResourceNotFoundException -- the DynamoDB/RDS lifecycle class -- so
// existence is "described AND not deleting", and absence accepts either
// signal.
type kinesisStreamVerifier struct{}

func (*kinesisStreamVerifier) IDOutputKey() string { return "stream_name" }

func (*kinesisStreamVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kinesisStreamExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisstream verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awskinesisstream %q not found after deploy", id)
	}
	return nil
}

func (*kinesisStreamVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kinesisStreamExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisstream verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awskinesisstream %q still exists after destroy", id)
	}
	return nil
}

// kinesisStreamExists reports whether the stream is present and not already
// on its way out. A ResourceNotFoundException is treated as absent.
func kinesisStreamExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := kinesis.NewFromConfig(cfg, func(o *kinesis.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{StreamName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.StreamDescriptionSummary == nil {
		return false, nil
	}
	if out.StreamDescriptionSummary.StreamStatus == types.StreamStatusDeleting {
		return false, nil
	}
	return true, nil
}
