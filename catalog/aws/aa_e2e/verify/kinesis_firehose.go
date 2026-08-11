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
// "described AND ACTIVE", and absence accepts either signal. ACTIVE is a
// post-apply invariant: both engines' create waiters return only once the
// stream reaches ACTIVE, so a CREATING or CREATING_FAILED stream after
// deploy means the waiter lied, never a timing artifact.
type kinesisFirehoseVerifier struct{}

func (*kinesisFirehoseVerifier) IDOutputKey() string { return "delivery_stream_name" }

func (*kinesisFirehoseVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	desc, err := kinesisFirehoseDescribe(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisfirehose verify-exists failed for %q", id)
	}
	if desc == nil || desc.DeliveryStreamStatus == types.DeliveryStreamStatusDeleting ||
		desc.DeliveryStreamStatus == types.DeliveryStreamStatusDeletingFailed {
		return pkgerrors.Errorf("awskinesisfirehose %q not found after deploy", id)
	}
	if desc.DeliveryStreamStatus != types.DeliveryStreamStatusActive {
		return pkgerrors.Errorf("awskinesisfirehose %q exists but reports status %q, want ACTIVE", id, desc.DeliveryStreamStatus)
	}
	return nil
}

func (*kinesisFirehoseVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	desc, err := kinesisFirehoseDescribe(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisfirehose verify-absent failed for %q", id)
	}
	if desc != nil && desc.DeliveryStreamStatus != types.DeliveryStreamStatusDeleting &&
		desc.DeliveryStreamStatus != types.DeliveryStreamStatusDeletingFailed {
		return pkgerrors.Errorf("awskinesisfirehose %q still exists after destroy", id)
	}
	return nil
}

// kinesisFirehoseDescribe returns the delivery stream description, or nil
// when AWS reports ResourceNotFoundException.
func kinesisFirehoseDescribe(ctx context.Context, cfg aws.Config, id, region string) (*types.DeliveryStreamDescription, error) {
	client := firehose.NewFromConfig(cfg, func(o *firehose.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDeliveryStream(ctx, &firehose.DescribeDeliveryStreamInput{DeliveryStreamName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil, nil
		}
		return nil, err
	}
	return out.DeliveryStreamDescription, nil
}
