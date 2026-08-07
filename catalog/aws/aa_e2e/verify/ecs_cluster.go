package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	pkgerrors "github.com/pkg/errors"
)

// ecsClusterVerifier verifies an AwsEcsCluster via DescribeClusters, keyed
// on the cluster ARN output. A deleted cluster is not merely missing: ECS
// keeps it describable as INACTIVE for a while after deletion, so existence
// is "described AND ACTIVE" and absence accepts a MISSING failure or any
// non-ACTIVE status.
type ecsClusterVerifier struct{}

func (*ecsClusterVerifier) IDOutputKey() string { return "cluster_arn" }

func (*ecsClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsClusterActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecscluster verify-exists failed for %q", id)
	}
	if !active {
		return pkgerrors.Errorf("awsecscluster %q not ACTIVE after deploy", id)
	}
	return nil
}

func (*ecsClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsClusterActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecscluster verify-absent failed for %q", id)
	}
	if active {
		return pkgerrors.Errorf("awsecscluster %q still ACTIVE after destroy", id)
	}
	return nil
}

func ecsClusterActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	client := ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		if region != "" {
			o.Region = region
		}
	})
	described, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: []string{arn},
	})
	if err != nil {
		return false, err
	}
	// DescribeClusters reports unknown clusters in Failures (reason
	// MISSING) rather than erroring.
	if len(described.Clusters) == 0 {
		return false, nil
	}
	cluster := described.Clusters[0]
	if cluster.Status == nil {
		return false, nil
	}
	return *cluster.Status == "ACTIVE", nil
}
