package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
	pkgerrors "github.com/pkg/errors"
)

// mskServerlessClusterVerifier verifies an AwsMskServerlessCluster via the
// Kafka DescribeClusterV2 API, keyed on the cluster_arn output (the ARN is the
// only describe key MSK accepts; DescribeClusterV2 serves both provisioned and
// serverless clusters). The lifecycle class matches the provisioned sibling: a
// cluster mid-deletion stays describable with a DELETING state before the
// service starts returning the typed NotFoundException, so existence is
// "described AND not deleting", and absence accepts either signal.
type mskServerlessClusterVerifier struct{}

func (*mskServerlessClusterVerifier) IDOutputKey() string { return "cluster_arn" }

func (*mskServerlessClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := mskServerlessClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmskserverlesscluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmskserverlesscluster %q not found after deploy", id)
	}
	return nil
}

func (*mskServerlessClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := mskServerlessClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmskserverlesscluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmskserverlesscluster %q still exists after destroy", id)
	}
	return nil
}

// mskServerlessClusterExists reports whether the cluster is present and not
// already on its way out. A NotFoundException is treated as absent.
func mskServerlessClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
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
