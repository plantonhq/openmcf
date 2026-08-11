package verify

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	pkgerrors "github.com/pkg/errors"
)

// ecsClusterVerifier verifies an AwsEcsCluster via DescribeClusters, keyed
// on the cluster ARN output. A deleted cluster is not merely missing: ECS
// keeps it describable as INACTIVE for a while after deletion, so existence
// is "described AND ACTIVE" and absence accepts a MISSING failure or any
// non-ACTIVE status. Existence additionally asserts every CUSTOM capacity
// provider attached to the cluster is ACTIVE via DescribeCapacityProviders
// -- the association PUT succeeding says nothing about the providers
// themselves, and a provider stuck outside ACTIVE is invisible to the
// cluster-status check. The FARGATE built-ins are skipped (AWS-owned,
// always active, never describable as customer resources in some
// partitions). Absence of the custom providers after destroy is covered by
// the account-level orphan sweep: a deleted cluster no longer enumerates
// its providers reliably.
type ecsClusterVerifier struct{}

func (*ecsClusterVerifier) IDOutputKey() string { return "cluster_arn" }

func (*ecsClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := ecsClient(cfg, region)
	cluster, err := ecsDescribeCluster(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecscluster verify-exists failed for %q", id)
	}
	if cluster == nil || cluster.Status == nil || *cluster.Status != "ACTIVE" {
		return pkgerrors.Errorf("awsecscluster %q not ACTIVE after deploy", id)
	}
	var custom []string
	for _, name := range cluster.CapacityProviders {
		if name == "FARGATE" || name == "FARGATE_SPOT" {
			continue
		}
		custom = append(custom, name)
	}
	if len(custom) == 0 {
		return nil
	}
	// Managed-instances providers provision asynchronously after create
	// (a seconds-wide PROVISIONING window observed live), so the ACTIVE
	// assertion polls briefly instead of failing on first read.
	var lastErr error
	for attempt := 0; attempt < 12; attempt++ {
		if attempt > 0 {
			time.Sleep(10 * time.Second)
		}
		lastErr = ecsCapacityProvidersActive(ctx, client, id, custom)
		if lastErr == nil {
			return nil
		}
	}
	return lastErr
}

// ecsCapacityProvidersActive asserts every named capacity provider is
// describable and ACTIVE.
func ecsCapacityProvidersActive(ctx context.Context, client *ecs.Client, clusterARN string, names []string) error {
	described, err := client.DescribeCapacityProviders(ctx, &ecs.DescribeCapacityProvidersInput{
		CapacityProviders: names,
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecscluster %q: describing attached capacity providers %v failed", clusterARN, names)
	}
	seen := map[string]string{}
	for _, cp := range described.CapacityProviders {
		if cp.Name == nil {
			continue
		}
		seen[*cp.Name] = string(cp.Status)
	}
	for _, name := range names {
		status, ok := seen[name]
		if !ok {
			return pkgerrors.Errorf("awsecscluster %q: attached capacity provider %q not returned by DescribeCapacityProviders", clusterARN, name)
		}
		if status != "ACTIVE" {
			return pkgerrors.Errorf("awsecscluster %q: capacity provider %q is %s, want ACTIVE", clusterARN, name, status)
		}
	}
	return nil
}

func (*ecsClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	cluster, err := ecsDescribeCluster(ctx, ecsClient(cfg, region), id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecscluster verify-absent failed for %q", id)
	}
	if cluster != nil && cluster.Status != nil && *cluster.Status == "ACTIVE" {
		return pkgerrors.Errorf("awsecscluster %q still ACTIVE after destroy", id)
	}
	return nil
}

func ecsClient(cfg aws.Config, region string) *ecs.Client {
	return ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// ecsDescribeCluster returns the cluster, or nil when ECS reports it
// unknown -- DescribeClusters reports unknown clusters in Failures (reason
// MISSING) rather than erroring.
func ecsDescribeCluster(ctx context.Context, client *ecs.Client, arn string) (*types.Cluster, error) {
	described, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: []string{arn},
	})
	if err != nil {
		return nil, err
	}
	if len(described.Clusters) == 0 {
		return nil, nil
	}
	return &described.Clusters[0], nil
}
