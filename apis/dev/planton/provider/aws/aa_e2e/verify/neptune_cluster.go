package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
	pkgerrors "github.com/pkg/errors"
)

// neptuneClusterVerifier verifies an AwsNeptuneCluster via the Neptune
// DescribeDBClusters, keyed on the cluster identifier output. A cluster
// mid-deletion stays describable with a "deleting" status before the
// service starts returning the typed DBClusterNotFoundFault -- the RDS
// lifecycle class -- so existence is "described AND not deleting", and
// absence accepts either signal. Status comparison is lowercased
// defensively; the service reports lowercase statuses but documents them
// without committing to case.
type neptuneClusterVerifier struct{}

func (*neptuneClusterVerifier) IDOutputKey() string { return "cluster_identifier" }

func (*neptuneClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := neptuneClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsneptunecluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsneptunecluster %q not found after deploy", id)
	}
	return nil
}

func (*neptuneClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := neptuneClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsneptunecluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsneptunecluster %q still exists after destroy", id)
	}
	return nil
}

// neptuneClusterExists reports whether the cluster is present and not
// already on its way out. A DBClusterNotFoundFault is treated as absent.
func neptuneClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := neptune.NewFromConfig(cfg, func(o *neptune.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{DBClusterIdentifier: &id})
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
