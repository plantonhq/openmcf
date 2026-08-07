package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
	pkgerrors "github.com/pkg/errors"
)

// mskClusterVerifier verifies an AwsMskCluster via the Kafka DescribeClusterV2
// API, keyed on the cluster_arn output (the ARN is the only describe key MSK
// accepts). A cluster mid-deletion stays describable with a DELETING state
// before the service starts returning the typed NotFoundException -- the same
// lifecycle class as RDS/Redshift -- so existence is "described AND not
// deleting", and absence accepts either signal.
type mskClusterVerifier struct{}

func (*mskClusterVerifier) IDOutputKey() string { return "cluster_arn" }

func (*mskClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := mskClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmskcluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmskcluster %q not found after deploy", id)
	}
	return nil
}

func (*mskClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := mskClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmskcluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmskcluster %q still exists after destroy", id)
	}
	return nil
}

// mskClusterExists reports whether the cluster is present and not already on
// its way out. A NotFoundException is treated as absent.
func mskClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := kafka.NewFromConfig(cfg, func(o *kafka.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeClusterV2(ctx, &kafka.DescribeClusterV2Input{ClusterArn: &id})
	if err != nil {
		var notFound *types.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.ClusterInfo == nil {
		return false, nil
	}
	if out.ClusterInfo.State == types.ClusterStateDeleting {
		return false, nil
	}
	return true, nil
}
