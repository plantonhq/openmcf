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
// signal. Existence is "described AND ACTIVE" -- a cluster describes during
// CREATING and DELETING too, and only ACTIVE proves a usable control plane.
// When the cluster reports a remote-network configuration (EKS Hybrid
// Nodes), existence additionally asserts the configuration carries at least
// one CIDR -- the posture is read from the cluster's own state, so any
// scenario declaring remote networks inherits the check. Absence stays
// describability-based: destroy waits for full deletion before verification
// runs, so a still-describable cluster after destroy is a real leak.
type eksClusterVerifier struct{}

func (*eksClusterVerifier) IDOutputKey() string { return "name" }

func (*eksClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	cluster, err := eksDescribeCluster(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsekscluster verify-exists failed for %q", id)
	}
	if cluster == nil {
		return pkgerrors.Errorf("awsekscluster %q not found after deploy", id)
	}
	if cluster.Status != types.ClusterStatusActive {
		return pkgerrors.Errorf("awsekscluster %q is %s after deploy, want ACTIVE", id, cluster.Status)
	}
	if rnc := cluster.RemoteNetworkConfig; rnc != nil {
		var cidrs int
		for _, n := range rnc.RemoteNodeNetworks {
			cidrs += len(n.Cidrs)
		}
		for _, p := range rnc.RemotePodNetworks {
			cidrs += len(p.Cidrs)
		}
		if cidrs == 0 {
			return pkgerrors.Errorf("awsekscluster %q reports a remote network config with no CIDRs", id)
		}
	}
	return nil
}

func (*eksClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	cluster, err := eksDescribeCluster(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsekscluster verify-absent failed for %q", id)
	}
	if cluster != nil {
		return pkgerrors.Errorf("awsekscluster %q still exists after destroy", id)
	}
	return nil
}

// eksDescribeCluster returns the cluster, or nil when EKS reports it unknown.
func eksDescribeCluster(ctx context.Context, cfg aws.Config, name, region string) (*types.Cluster, error) {
	client := eks.NewFromConfig(cfg, func(o *eks.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: &name})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil, nil
		}
		return nil, err
	}
	return out.Cluster, nil
}
