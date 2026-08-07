package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	pkgerrors "github.com/pkg/errors"
)

// eventBridgeBusVerifier verifies an AwsEventBridgeBus via DescribeEventBus,
// keyed on bus_name.
type eventBridgeBusVerifier struct{}

func (*eventBridgeBusVerifier) IDOutputKey() string { return "bus_name" }

func (*eventBridgeBusVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eventBridgeBusExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgebus verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseventbridgebus %q not found after deploy", id)
	}
	return nil
}

func (*eventBridgeBusVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eventBridgeBusExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgebus verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseventbridgebus %q still exists after destroy", id)
	}
	return nil
}

func eventBridgeBusExists(ctx context.Context, cfg aws.Config, busName, region string) (bool, error) {
	client := eventbridge.NewFromConfig(cfg, func(o *eventbridge.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
		Name: aws.String(busName),
	})
	if err == nil {
		return true, nil
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}
