package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/opensearch/types"
	pkgerrors "github.com/pkg/errors"
)

// opensearchDomainVerifier verifies an AwsOpenSearchDomain via DescribeDomain,
// keyed on the domain_name output. A domain mid-deletion stays describable
// with the Deleted flag raised (deletion takes several minutes) before the
// service starts returning the typed ResourceNotFoundException, so existence
// is "described AND not flagged deleted", and absence accepts either signal.
type opensearchDomainVerifier struct{}

func (*opensearchDomainVerifier) IDOutputKey() string { return "domain_name" }

func (*opensearchDomainVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := opensearchDomainExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsopensearchdomain verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsopensearchdomain %q not found after deploy", id)
	}
	return nil
}

func (*opensearchDomainVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := opensearchDomainExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsopensearchdomain verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsopensearchdomain %q still exists after destroy", id)
	}
	return nil
}

// opensearchDomainExists reports whether the domain is present and not
// already on its way out. A ResourceNotFoundException is treated as absent.
func opensearchDomainExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := opensearch.NewFromConfig(cfg, func(o *opensearch.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDomain(ctx, &opensearch.DescribeDomainInput{DomainName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.DomainStatus == nil {
		return false, nil
	}
	if out.DomainStatus.Deleted != nil && *out.DomainStatus.Deleted {
		return false, nil
	}
	return true, nil
}
