package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	pkgerrors "github.com/pkg/errors"
)

// eventBridgeApiDestinationVerifier verifies an
// AwsEventBridgeApiDestination arm-for-arm from its outputs: the
// connection via DescribeConnection and the destination via
// DescribeApiDestination (names parsed from the ARNs -
// "arn:...:connection/{name}/{uuid}" /
// "arn:...:api-destination/{name}/{uuid}"). Either arm may be absent
// (the spec's at-least-one-arm contract) - empty outputs skip their
// checks. A DELETING connection counts as absent.
type eventBridgeApiDestinationVerifier struct{}

func (*eventBridgeApiDestinationVerifier) IDOutputKey() string { return "api_destination_arn" }

func (v *eventBridgeApiDestinationVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	if id == "" {
		return nil
	}
	exists, err := apiDestinationExistsByArn(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgeapidestination verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseventbridgeapidestination %q not found after deploy", id)
	}
	return nil
}

func (v *eventBridgeApiDestinationVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	if id == "" {
		return nil
	}
	exists, err := apiDestinationExistsByArn(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgeapidestination verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseventbridgeapidestination %q still exists after destroy", id)
	}
	return nil
}

func (v *eventBridgeApiDestinationVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	verifiedAnArm := false
	if connectionArn := stringOutput(outputs, "connection_arn"); connectionArn != "" {
		exists, err := eventConnectionExistsByArn(ctx, cfg, connectionArn, region)
		if err != nil {
			return pkgerrors.Wrapf(err, "awseventbridgeapidestination verify-exists failed for connection %q", connectionArn)
		}
		if !exists {
			return pkgerrors.Errorf("awseventbridgeapidestination connection %q not found after deploy", connectionArn)
		}
		verifiedAnArm = true
	}
	if destinationArn := stringOutput(outputs, "api_destination_arn"); destinationArn != "" {
		if err := v.VerifyExists(ctx, cfg, destinationArn, region); err != nil {
			return err
		}
		verifiedAnArm = true
	}
	if !verifiedAnArm {
		return pkgerrors.New("awseventbridgeapidestination verify-exists: neither connection_arn nor api_destination_arn in outputs")
	}
	return nil
}

func (v *eventBridgeApiDestinationVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	if connectionArn := stringOutput(outputs, "connection_arn"); connectionArn != "" {
		exists, err := eventConnectionExistsByArn(ctx, cfg, connectionArn, region)
		if err != nil {
			return pkgerrors.Wrapf(err, "awseventbridgeapidestination verify-absent failed for connection %q", connectionArn)
		}
		if exists {
			return pkgerrors.Errorf("awseventbridgeapidestination connection %q still exists after destroy", connectionArn)
		}
	}
	if destinationArn := stringOutput(outputs, "api_destination_arn"); destinationArn != "" {
		if err := v.VerifyAbsent(ctx, cfg, destinationArn, region); err != nil {
			return err
		}
	}
	return nil
}

func eventsClientForRegion(cfg aws.Config, region string) *eventbridge.Client {
	return eventbridge.NewFromConfig(cfg, func(o *eventbridge.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// eventsNameFromArn extracts {name} from
// "arn:...:{prefix}/{name}/{uuid}".
func eventsNameFromArn(arn, prefix string) (string, error) {
	resource := strings.TrimPrefix(arnResourcePart(arn), prefix+"/")
	parts := strings.SplitN(resource, "/", 2)
	if parts[0] == "" {
		return "", pkgerrors.Errorf("unexpected events ARN shape %q", arn)
	}
	return parts[0], nil
}

func eventConnectionExistsByArn(ctx context.Context, cfg aws.Config, connectionArn, region string) (bool, error) {
	name, err := eventsNameFromArn(connectionArn, "connection")
	if err != nil {
		return false, err
	}
	out, err := eventsClientForRegion(cfg, region).DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
		Name: aws.String(name),
	})
	if err != nil {
		var notFound *eventbridgetypes.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.ConnectionState == eventbridgetypes.ConnectionStateDeleting {
		return false, nil
	}
	return true, nil
}

func apiDestinationExistsByArn(ctx context.Context, cfg aws.Config, destinationArn, region string) (bool, error) {
	name, err := eventsNameFromArn(destinationArn, "api-destination")
	if err != nil {
		return false, err
	}
	_, err = eventsClientForRegion(cfg, region).DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
		Name: aws.String(name),
	})
	if err == nil {
		return true, nil
	}
	var notFound *eventbridgetypes.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}
