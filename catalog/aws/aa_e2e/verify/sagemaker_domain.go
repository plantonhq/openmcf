package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	pkgerrors "github.com/pkg/errors"
)

// sageMakerDomainVerifier verifies an AwsSagemakerDomain via DescribeDomain,
// keyed on the domain_id output. The verifier is lifecycle-state-aware (the
// NAT-gateway class): domain deletion is asynchronous -- the domain stays
// describable in the Deleting state while SageMaker tears down apps and (with
// a Delete retention policy) the EFS home file system -- so Deleting/Deleted
// count as absent rather than racing the API's cleanup.
type sageMakerDomainVerifier struct{}

func (*sageMakerDomainVerifier) IDOutputKey() string { return "domain_id" }

func (v *sageMakerDomainVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sageMakerDomainExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssagemakerdomain verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssagemakerdomain %q not found after deploy", id)
	}
	return nil
}

func (v *sageMakerDomainVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := sageMakerDomainExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssagemakerdomain verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssagemakerdomain %q still exists after destroy", id)
	}
	return nil
}

func sageMakerDomainExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := sagemaker.NewFromConfig(cfg, func(o *sagemaker.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDomain(ctx, &sagemaker.DescribeDomainInput{
		DomainId: aws.String(id),
	})
	if err != nil {
		// SageMaker reports a missing domain as ResourceNotFound.
		var notFound *sagemakertypes.ResourceNotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	switch strings.ToLower(string(out.Status)) {
	case "deleting", "delete_failed":
		// Deleting is on its way out; DeleteFailed would leave a real orphan,
		// but claiming "exists" there turns the destroy lane red so the
		// failure is investigated instead of silently passing.
		return out.Status == sagemakertypes.DomainStatusDeleteFailed, nil
	}
	return true, nil
}
