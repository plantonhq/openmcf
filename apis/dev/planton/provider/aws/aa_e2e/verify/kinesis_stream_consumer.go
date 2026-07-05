package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pkgerrors "github.com/pkg/errors"
)

// kinesisStreamConsumerVerifier verifies an AwsKinesisStreamConsumer via
// DescribeStreamConsumer, keyed on the consumer_arn output (the ARN is the
// consumer's unique identity -- names are only unique per stream). A
// consumer mid-deregistration stays describable with a DELETING status
// before the typed ResourceNotFoundException appears, so existence is
// "described AND not deleting", and absence accepts either signal.
type kinesisStreamConsumerVerifier struct{}

func (*kinesisStreamConsumerVerifier) IDOutputKey() string { return "consumer_arn" }

func (*kinesisStreamConsumerVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kinesisStreamConsumerExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisstreamconsumer verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awskinesisstreamconsumer %q not found after deploy", id)
	}
	return nil
}

func (*kinesisStreamConsumerVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kinesisStreamConsumerExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskinesisstreamconsumer verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awskinesisstreamconsumer %q still exists after destroy", id)
	}
	return nil
}

// kinesisStreamConsumerExists reports whether the consumer is registered and
// not already deregistering. A ResourceNotFoundException is treated as
// absent.
func kinesisStreamConsumerExists(ctx context.Context, cfg aws.Config, consumerARN, region string) (bool, error) {
	client := kinesis.NewFromConfig(cfg, func(o *kinesis.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeStreamConsumer(ctx, &kinesis.DescribeStreamConsumerInput{ConsumerARN: &consumerARN})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.ConsumerDescription == nil {
		return false, nil
	}
	if out.ConsumerDescription.ConsumerStatus == types.ConsumerStatusDeleting {
		return false, nil
	}
	return true, nil
}
