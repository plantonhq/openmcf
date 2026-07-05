package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	pkgerrors "github.com/pkg/errors"
)

// sqsQueueVerifier verifies an AwsSqsQueue via GetQueueUrl, keyed on queue_url.
type sqsQueueVerifier struct{}

func (*sqsQueueVerifier) IDOutputKey() string { return "queue_url" }

func (*sqsQueueVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sqsQueueExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssqsqueue verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssqsqueue %q not found after deploy", id)
	}
	return nil
}

func (*sqsQueueVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sqsQueueExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssqsqueue verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssqsqueue %q still exists after destroy", id)
	}
	return nil
}

func sqsQueueExists(ctx context.Context, cfg aws.Config, queueURL, region string) (bool, error) {
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
	})
	if err == nil {
		return true, nil
	}
	var notExist *types.QueueDoesNotExist
	if errors.As(err, &notExist) {
		return false, nil
	}
	return false, err
}
