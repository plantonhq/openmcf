package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	pkgerrors "github.com/pkg/errors"
)

// kinesisFirehoseVerifier verifies an AwsKinesisFirehose via
// DescribeDeliveryStream, keyed on the delivery_stream_name output. A
// delivery stream mid-deletion stays describable with a DELETING status
// before the typed ResourceNotFoundException appears, so existence is
// "described AND not deleting", and absence accepts either signal.
type kinesisFirehoseVerifier struct{}

func (*kinesisFirehoseVerifier) IDOutputKey() string { return "delivery_stream_name" }

func (*kinesisFirehoseVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kinesisFirehoseExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisfirehose verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awskinesisfirehose %q not found after deploy", id)
	}
	return nil
}

func (*kinesisFirehoseVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kinesisFirehoseExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisfirehose verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awskinesisfirehose %q still exists after destroy", id)
	}
	return nil
}

// kinesisFirehoseExists reports whether the delivery stream is present and
// not already on its way out. A ResourceNotFoundException is treated as
// absent.
func kinesisFirehoseExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := firehose.NewFromConfig(cfg, func(o *firehose.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDeliveryStream(ctx, &firehose.DescribeDeliveryStreamInput{DeliveryStreamName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.DeliveryStreamDescription == nil {
		return false, nil
	}
	status := out.DeliveryStreamDescription.DeliveryStreamStatus
	if status == types.DeliveryStreamStatusDeleting || status == types.DeliveryStreamStatusDeletingFailed {
		return false, nil
	}
	return true, nil
}
