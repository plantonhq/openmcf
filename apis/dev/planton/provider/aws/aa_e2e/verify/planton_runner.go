package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	pkgerrors "github.com/pkg/errors"
)

// plantonRunnerVerifier verifies an AwsPlantonRunner via the ECS service
// that keeps the runner running, keyed on the service ARN (which encodes
// both the dedicated cluster and the service name). Verification is
// deliberately service-level, independent of task health: ECS reports the
// service ACTIVE regardless of whether the container's control-plane
// handshake succeeds, so the lane proves the appliance's full
// provisioning surface (secret, roles, log group, security group,
// cluster, task definition, service) with synthetic credentials and
// without a live control plane.
type plantonRunnerVerifier struct{}

func (*plantonRunnerVerifier) IDOutputKey() string { return "service_arn" }

func (*plantonRunnerVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsServiceActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsplantonrunner verify-exists failed for %q", id)
	}
	if !active {
		return pkgerrors.Errorf("awsplantonrunner %q not ACTIVE after deploy", id)
	}
	return nil
}

func (*plantonRunnerVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsServiceActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsplantonrunner verify-absent failed for %q", id)
	}
	if active {
		return pkgerrors.Errorf("awsplantonrunner %q still ACTIVE after destroy", id)
	}
	return nil
}
