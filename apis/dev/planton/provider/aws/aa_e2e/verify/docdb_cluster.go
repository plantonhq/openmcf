package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdb"
	"github.com/aws/aws-sdk-go-v2/service/docdb/types"
	pkgerrors "github.com/pkg/errors"
)

// docdbClusterVerifier verifies an AwsDocumentDb via the DocumentDB
// DescribeDBClusters, keyed on the cluster identifier output. A cluster
// mid-deletion stays describable with a "deleting" status before the
// service starts returning the typed DBClusterNotFoundFault -- the RDS
// lifecycle class -- so existence is "described AND not deleting", and
// absence accepts either signal. Status comparison is lowercased
// defensively; the service reports lowercase statuses but documents them
// without committing to case.
type docdbClusterVerifier struct{}

func (*docdbClusterVerifier) IDOutputKey() string { return "cluster_identifier" }

func (*docdbClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := docdbClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsdocumentdb verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsdocumentdb %q not found after deploy", id)
	}
	return nil
}

func (*docdbClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := docdbClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsdocumentdb verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsdocumentdb %q still exists after destroy", id)
	}
	return nil
}

// docdbClusterExists reports whether the cluster is present and not already
// on its way out. A DBClusterNotFoundFault is treated as absent.
func docdbClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := docdb.NewFromConfig(cfg, func(o *docdb.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDBClusters(ctx, &docdb.DescribeDBClustersInput{DBClusterIdentifier: &id})
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
