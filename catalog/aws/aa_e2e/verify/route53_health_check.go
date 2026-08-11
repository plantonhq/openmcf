package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	pkgerrors "github.com/pkg/errors"
)

// route53HealthCheckVerifier verifies an AwsRoute53HealthCheck via
// GetHealthCheck, keyed on the health_check_id output. Health checks are
// global objects with synchronous deletion — a deleted check returns the
// typed NoSuchHealthCheck immediately.
type route53HealthCheckVerifier struct{}

func (*route53HealthCheckVerifier) IDOutputKey() string { return "health_check_id" }

func (*route53HealthCheckVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := route53.NewFromConfig(cfg)
	out, err := client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{HealthCheckId: aws.String(id)})
	if err != nil {
		var notFound *route53types.NoSuchHealthCheck
		if errors.As(err, &notFound) {
			return pkgerrors.Errorf("awsroute53healthcheck %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awsroute53healthcheck verify-exists failed for %q", id)
	}

	// Posture coherence from the resource's own state: a CALCULATED check
	// exists only as an aggregation, so AWS must be storing the child set it
	// was sent (the threshold's exact value — including the always-healthy
	// explicit 0 — is lane evidence via GetHealthCheck snapshots).
	if hcConfig := out.HealthCheck.HealthCheckConfig; hcConfig != nil && hcConfig.Type == route53types.HealthCheckTypeCalculated {
		if len(hcConfig.ChildHealthChecks) == 0 {
			return pkgerrors.Errorf("awsroute53healthcheck %q is CALCULATED but AWS reports no child health checks", id)
		}
	}
	return nil
}

func (*route53HealthCheckVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := route53HealthCheckExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsroute53healthcheck verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsroute53healthcheck %q still exists after destroy", id)
	}
	return nil
}

func route53HealthCheckExists(ctx context.Context, cfg aws.Config, id string) (bool, error) {
	client := route53.NewFromConfig(cfg)
	_, err := client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{HealthCheckId: aws.String(id)})
	if err != nil {
		var notFound *route53types.NoSuchHealthCheck
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
