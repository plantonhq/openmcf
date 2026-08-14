package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
	pkgerrors "github.com/pkg/errors"
)

// openSearchServerlessCollectionVerifier verifies an
// AwsOpenSearchServerlessCollection via BatchGetCollection, keyed on
// collection_name. A collection reported DELETING counts as absent (the
// deletion is already irreversible); BatchGetCollection returns an empty
// list rather than a not-found error for unknown names.
type openSearchServerlessCollectionVerifier struct{}

func (*openSearchServerlessCollectionVerifier) IDOutputKey() string { return "collection_name" }

func (*openSearchServerlessCollectionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := openSearchServerlessCollectionExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsopensearchserverlesscollection verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsopensearchserverlesscollection %q not found after deploy", id)
	}
	return nil
}

func (*openSearchServerlessCollectionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := openSearchServerlessCollectionExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsopensearchserverlesscollection verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsopensearchserverlesscollection %q still exists after destroy", id)
	}
	return nil
}

func openSearchServerlessCollectionExists(ctx context.Context, cfg aws.Config, collectionName, region string) (bool, error) {
	client := opensearchserverless.NewFromConfig(cfg, func(o *opensearchserverless.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.BatchGetCollection(ctx, &opensearchserverless.BatchGetCollectionInput{
		Names: []string{collectionName},
	})
	if err != nil {
		return false, err
	}
	for _, detail := range out.CollectionDetails {
		if aws.ToString(detail.Name) == collectionName {
			// DELETING means the destroy already happened and AWS is
			// tearing the collection down -- absent for our purposes.
			if detail.Status == types.CollectionStatusDeleting {
				return false, nil
			}
			return true, nil
		}
	}
	return false, nil
}
