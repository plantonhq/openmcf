package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mwaa"
	"github.com/aws/aws-sdk-go-v2/service/mwaa/types"
	pkgerrors "github.com/pkg/errors"
)

// mwaaEnvironmentVerifier verifies an AwsMwaaEnvironment via the MWAA
// GetEnvironment API, keyed on the environment_name output (the name is the
// only get key MWAA accepts). An environment mid-deletion stays describable
// with a DELETING status before the service starts returning the typed
// ResourceNotFoundException -- the same lifecycle class as MSK/RDS -- so
// existence is "described AND not deleting", and absence accepts either
// signal.
type mwaaEnvironmentVerifier struct{}

func (*mwaaEnvironmentVerifier) IDOutputKey() string { return "environment_name" }

func (*mwaaEnvironmentVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := mwaaEnvironmentExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmwaaenvironment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmwaaenvironment %q not found after deploy", id)
	}
	return nil
}

func (*mwaaEnvironmentVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := mwaaEnvironmentExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmwaaenvironment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmwaaenvironment %q still exists after destroy", id)
	}
	return nil
}

// mwaaEnvironmentExists reports whether the environment is present and not
// already on its way out. A ResourceNotFoundException is treated as absent.
func mwaaEnvironmentExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := mwaa.NewFromConfig(cfg, func(o *mwaa.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.GetEnvironment(ctx, &mwaa.GetEnvironmentInput{Name: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.Environment == nil {
		return false, nil
	}
	if out.Environment.Status == types.EnvironmentStatusDeleting {
		return false, nil
	}
	return true, nil
}
