package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	pkgerrors "github.com/pkg/errors"
)

// snsTopicVerifier verifies an AwsSnsTopic via GetTopicAttributes, keyed on
// topic_arn.
type snsTopicVerifier struct{}

func (*snsTopicVerifier) IDOutputKey() string { return "topic_arn" }

func (*snsTopicVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := snsTopicExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssnstopic verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssnstopic %q not found after deploy", id)
	}
	return nil
}

func (*snsTopicVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := snsTopicExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssnstopic verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssnstopic %q still exists after destroy", id)
	}
	return nil
}

func snsTopicExists(ctx context.Context, cfg aws.Config, topicARN, region string) (bool, error) {
	client := sns.NewFromConfig(cfg, func(o *sns.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: aws.String(topicARN),
	})
	if err == nil {
		return true, nil
	}
	var notFound *types.NotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}
