package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	pkgerrors "github.com/pkg/errors"
)

// eksNodeGroupVerifier verifies an AwsEksNodeGroup via DescribeNodegroup.
// DescribeNodegroup needs BOTH the cluster and node group names, and the
// harness carries exactly one identifier per component -- so this verifier
// keys on the node group ARN, which encodes both:
// arn:aws:eks:<region>:<account>:nodegroup/<cluster>/<group>/<uuid>.
// The typed ResourceNotFoundException is the absence signal (raised for a
// missing node group AND for a missing cluster, which covers teardown
// order either way).
type eksNodeGroupVerifier struct{}

func (*eksNodeGroupVerifier) IDOutputKey() string { return "nodegroup_arn" }

func (*eksNodeGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksNodeGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksnodegroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseksnodegroup %q not found after deploy", id)
	}
	return nil
}

func (*eksNodeGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := eksNodeGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseksnodegroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseksnodegroup %q still exists after destroy", id)
	}
	return nil
}

func eksNodeGroupExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	clusterName, nodeGroupName, err := parseNodeGroupArn(arn)
	if err != nil {
		return false, err
	}
	client := eks.NewFromConfig(cfg, func(o *eks.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err = client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   &clusterName,
		NodegroupName: &nodeGroupName,
	})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// parseNodeGroupArn extracts the cluster and node group names from a node
// group ARN's resource part: nodegroup/<cluster>/<group>/<uuid>.
func parseNodeGroupArn(arn string) (clusterName string, nodeGroupName string, err error) {
	const marker = ":nodegroup/"
	idx := strings.Index(arn, marker)
	if idx < 0 {
		return "", "", pkgerrors.Errorf("not an EKS node group ARN: %q", arn)
	}
	parts := strings.Split(arn[idx+len(marker):], "/")
	if len(parts) < 2 {
		return "", "", pkgerrors.Errorf("malformed EKS node group ARN: %q", arn)
	}
	return parts[0], parts[1], nil
}
