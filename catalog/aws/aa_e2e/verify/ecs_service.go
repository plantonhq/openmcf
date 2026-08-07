package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	pkgerrors "github.com/pkg/errors"
)

// ecsServiceVerifier verifies an AwsEcsService via DescribeServices.
// DescribeServices needs BOTH the cluster and service names, and the
// harness carries exactly one identifier per component -- so this verifier
// keys on the service ARN, which encodes both:
// arn:aws:ecs:<region>:<account>:service/<cluster>/<service>.
//
// A deleted service is not merely missing: ECS keeps it describable as
// INACTIVE for a while after deletion, so existence is "described AND
// ACTIVE" and absence accepts a MISSING failure reason, a missing cluster,
// or any non-ACTIVE status.
type ecsServiceVerifier struct{}

func (*ecsServiceVerifier) IDOutputKey() string { return "service_arn" }

func (*ecsServiceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsServiceActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecsservice verify-exists failed for %q", id)
	}
	if !active {
		return pkgerrors.Errorf("awsecsservice %q not ACTIVE after deploy", id)
	}
	return nil
}

func (*ecsServiceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsServiceActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecsservice verify-absent failed for %q", id)
	}
	if active {
		return pkgerrors.Errorf("awsecsservice %q still ACTIVE after destroy", id)
	}
	return nil
}

func ecsServiceActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	clusterName, serviceName, err := parseEcsServiceArn(arn)
	if err != nil {
		return false, err
	}
	client := ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		if region != "" {
			o.Region = region
		}
	})
	described, err := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: []string{serviceName},
	})
	if err != nil {
		// The cluster being gone covers teardown order: prerequisites are
		// destroyed after the component, so a vanished cluster means the
		// service is gone too.
		var clusterNotFound *types.ClusterNotFoundException
		if errors.As(err, &clusterNotFound) {
			return false, nil
		}
		return false, err
	}
	// DescribeServices reports unknown services in Failures (reason
	// MISSING) rather than erroring.
	if len(described.Services) == 0 {
		return false, nil
	}
	service := described.Services[0]
	if service.Status == nil {
		return false, nil
	}
	return *service.Status == "ACTIVE", nil
}

// parseEcsServiceArn extracts the cluster and service names from a service
// ARN's resource part: service/<cluster>/<service>.
func parseEcsServiceArn(arn string) (clusterName string, serviceName string, err error) {
	const marker = ":service/"
	idx := strings.Index(arn, marker)
	if idx < 0 {
		return "", "", pkgerrors.Errorf("not an ECS service ARN: %q", arn)
	}
	parts := strings.Split(arn[idx+len(marker):], "/")
	if len(parts) < 2 {
		return "", "", pkgerrors.Errorf("malformed ECS service ARN (expected service/<cluster>/<service>): %q", arn)
	}
	return parts[0], parts[1], nil
}
