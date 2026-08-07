package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	pkgerrors "github.com/pkg/errors"
)

// snsSubscriptionVerifier verifies an AwsSnsSubscription via
// GetSubscriptionAttributes, keyed on subscription_arn.
type snsSubscriptionVerifier struct{}

func (*snsSubscriptionVerifier) IDOutputKey() string { return "subscription_arn" }

func (*snsSubscriptionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := snsSubscriptionExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssnssubscription verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssnssubscription %q not found after deploy", id)
	}
	return nil
}

func (*snsSubscriptionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := snsSubscriptionExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssnssubscription verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssnssubscription %q still exists after destroy", id)
	}
	return nil
}

func snsSubscriptionExists(ctx context.Context, cfg aws.Config, subscriptionARN, region string) (bool, error) {
	client := sns.NewFromConfig(cfg, func(o *sns.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetSubscriptionAttributes(ctx, &sns.GetSubscriptionAttributesInput{
		SubscriptionArn: aws.String(subscriptionARN),
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
