package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	pkgerrors "github.com/pkg/errors"
)

// eksClusterVerifier verifies an AwsEksCluster via DescribeCluster, keyed on
// the cluster name. The typed ResourceNotFoundException is the absence
// signal. A cluster mid-deletion still describes successfully (status
// DELETING), but destroy waits for full deletion before verification runs,
// so existence here means a real control plane.
type eksClusterVerifier struct{}

func (*eksClusterVerifier) IDOutputKey() string { return "name" }

func (*eksClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsekscluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsekscluster %q not found after deploy", id)
	}
	return nil
}

func (*eksClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsekscluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsekscluster %q still exists after destroy", id)
	}
	return nil
}

func eksClusterExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := eks.NewFromConfig(cfg, func(o *eks.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: &name})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
