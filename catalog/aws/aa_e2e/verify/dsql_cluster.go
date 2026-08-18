package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dsql"
	dsqltypes "github.com/aws/aws-sdk-go-v2/service/dsql/types"
	pkgerrors "github.com/pkg/errors"
)

// dsqlClusterVerifier verifies an AwsAuroraDsql via GetCluster, keyed
// on the AWS-generated cluster identifier (the provider's import ID).
// Exists demands ACTIVE - a cluster mid-setup reports CREATING or
// PENDING_SETUP. A deleted cluster answers ResourceNotFoundException
// (the absent signal); DELETING/DELETED count as absent too, since
// DSQL holds the record briefly after delete.
type dsqlClusterVerifier struct{}

func (*dsqlClusterVerifier) IDOutputKey() string { return "identifier" }

func (*dsqlClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := dsql.NewFromConfig(cfg).GetCluster(ctx, &dsql.GetClusterInput{
		Identifier: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsauroradsql verify-exists failed for %q", id)
	}
	if out.Status != dsqltypes.ClusterStatusActive {
		return pkgerrors.Errorf("awsauroradsql %q status %q, want ACTIVE", id, out.Status)
	}
	return nil
}

func (*dsqlClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := dsql.NewFromConfig(cfg).GetCluster(ctx, &dsql.GetClusterInput{
		Identifier: aws.String(id),
	})
	if err != nil {
		var notFound *dsqltypes.ResourceNotFoundException
		if pkgerrors.As(err, &notFound) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awsauroradsql verify-absent failed for %q", id)
	}
	if out.Status == dsqltypes.ClusterStatusDeleting || out.Status == dsqltypes.ClusterStatusDeleted {
		return nil
	}
	return pkgerrors.Errorf("awsauroradsql %q still exists after destroy (status %q)", id, out.Status)
}
