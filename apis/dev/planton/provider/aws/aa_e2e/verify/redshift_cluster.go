package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	pkgerrors "github.com/pkg/errors"
)

// redshiftClusterVerifier verifies an AwsRedshiftCluster via the Redshift
// DescribeClusters, keyed on the cluster identifier output. A cluster
// mid-deletion stays describable with a "deleting" status before the
// service starts returning the typed ClusterNotFoundFault -- the RDS
// lifecycle class -- so existence is "described AND not deleting", and
// absence accepts either signal. Status comparison is lowercased
// defensively; the service reports lowercase statuses but documents them
// without committing to case.
type redshiftClusterVerifier struct{}

func (*redshiftClusterVerifier) IDOutputKey() string { return "cluster_identifier" }

func (*redshiftClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := redshiftClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsredshiftcluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsredshiftcluster %q not found after deploy", id)
	}
	return nil
}

func (*redshiftClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := redshiftClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsredshiftcluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsredshiftcluster %q still exists after destroy", id)
	}
	return nil
}

// redshiftClusterExists reports whether the cluster is present and not
// already on its way out. A ClusterNotFoundFault is treated as absent.
func redshiftClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := redshift.NewFromConfig(cfg, func(o *redshift.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeClusters(ctx, &redshift.DescribeClustersInput{ClusterIdentifier: &id})
	if err != nil {
		var notFound *types.ClusterNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, cluster := range out.Clusters {
		if cluster.ClusterStatus == nil {
			continue
		}
		switch strings.ToLower(*cluster.ClusterStatus) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}
