package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pipes"
	pipestypes "github.com/aws/aws-sdk-go-v2/service/pipes/types"
	pkgerrors "github.com/pkg/errors"
)

// eventBridgePipeVerifier verifies an AwsEventBridgePipe via
// DescribePipe, keyed on pipe_name. Pipe deletes drain through a
// DELETING state for up to minutes - a describable DELETING pipe counts
// as absent.
type eventBridgePipeVerifier struct{}

func (*eventBridgePipeVerifier) IDOutputKey() string { return "pipe_name" }

func (*eventBridgePipeVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eventBridgePipeExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgepipe verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseventbridgepipe %q not found after deploy", id)
	}
	return nil
}

func (*eventBridgePipeVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eventBridgePipeExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgepipe verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseventbridgepipe %q still exists after destroy", id)
	}
	return nil
}

func eventBridgePipeExists(ctx context.Context, cfg aws.Config, pipeName, region string) (bool, error) {
	client := pipes.NewFromConfig(cfg, func(o *pipes.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribePipe(ctx, &pipes.DescribePipeInput{
		Name: aws.String(pipeName),
	})
	if err != nil {
		var notFound *pipestypes.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	switch out.CurrentState {
	case pipestypes.PipeStateDeleting:
		// Deletes drain - a DELETING pipe counts as absent.
		return false, nil
	default:
		return true, nil
	}
}
