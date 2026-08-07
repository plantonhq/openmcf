package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	pkgerrors "github.com/pkg/errors"
)

// rdsClusterVerifier verifies an AwsRdsCluster via DescribeDBClusters, keyed
// on the cluster identifier output. A cluster mid-deletion stays describable
// with a "deleting" status before RDS starts returning the typed
// DBClusterNotFoundFault -- the NAT-gateway lifecycle class -- so existence is
// "described AND not deleting", and absence accepts either signal. Status
// comparison is lowercased defensively; RDS reports lowercase statuses but
// documents them without committing to case.
type rdsClusterVerifier struct{}

func (*rdsClusterVerifier) IDOutputKey() string { return "cluster_identifier" }

func (*rdsClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := rdsClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrdscluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsrdscluster %q not found after deploy", id)
	}
	return nil
}

func (*rdsClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := rdsClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrdscluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsrdscluster %q still exists after destroy", id)
	}
	return nil
}

// rdsClusterExists reports whether the cluster is present and not already on
// its way out. A DBClusterNotFoundFault is treated as absent.
func rdsClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := rds.NewFromConfig(cfg, func(o *rds.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{DBClusterIdentifier: &id})
	if err != nil {
		var notFound *types.DBClusterNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, cluster := range out.DBClusters {
		if cluster.Status == nil {
			continue
		}
		switch strings.ToLower(*cluster.Status) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}
