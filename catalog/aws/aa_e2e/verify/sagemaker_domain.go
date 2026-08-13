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

// VerifyExistsFromOutputs raises the bar to the folded satellites: beyond the
// domain itself, every entry in the user_profile_arns and space_arns output
// maps (keyed by the spec names) must be describable and InService.
func (v *sageMakerDomainVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	domainId := stringOutputMap(outputs, "domain_id")
	if domainId == "" {
		return pkgerrors.New("awssagemakerdomain verify-exists: no domain_id in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, domainId, region); err != nil {
		return err
	}

	client := sagemaker.NewFromConfig(cfg, func(o *sagemaker.Options) {
		if region != "" {
			o.Region = region
		}
	})

	for name := range namedOutputMap(outputs, "user_profile_arns") {
		out, err := client.DescribeUserProfile(ctx, &sagemaker.DescribeUserProfileInput{
			DomainId:        aws.String(domainId),
			UserProfileName: aws.String(name),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awssagemakerdomain %q: user profile %q not describable", domainId, name)
		}
		if out.Status != sagemakertypes.UserProfileStatusInService {
			return pkgerrors.Errorf("awssagemakerdomain %q: user profile %q is %s, want InService", domainId, name, out.Status)
		}
	}

	for name := range namedOutputMap(outputs, "space_arns") {
		out, err := client.DescribeSpace(ctx, &sagemaker.DescribeSpaceInput{
			DomainId:  aws.String(domainId),
			SpaceName: aws.String(name),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awssagemakerdomain %q: space %q not describable", domainId, name)
		}
		if out.Status != sagemakertypes.SpaceStatusInService {
			return pkgerrors.Errorf("awssagemakerdomain %q: space %q is %s, want InService", domainId, name, out.Status)
		}
	}

	return nil
}

// VerifyAbsentFromOutputs delegates to the domain check: profiles and spaces
// cannot outlive their domain (the destroy path removes them first, and AWS
// refuses to delete a domain that still carries them), so domain absence
// proves the satellites are gone.
func (v *sageMakerDomainVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	domainId := stringOutputMap(outputs, "domain_id")
	if domainId == "" {
		return pkgerrors.New("awssagemakerdomain verify-absent: no domain_id in outputs")
	}
	return v.VerifyAbsent(ctx, cfg, domainId, region)
}

// namedOutputMap reads a map<string,string> output (spec name -> identifier).
// A missing or empty map means the spec declared no such satellites.
func namedOutputMap(outputs map[string]interface{}, key string) map[string]string {
	raw, ok := outputs[key]
	if !ok || raw == nil {
		return nil
	}
	named := map[string]string{}
	switch m := raw.(type) {
	case map[string]interface{}:
		for name, id := range m {
			if s, isStr := id.(string); isStr && s != "" {
				named[name] = s
			}
		}
	case map[string]string:
		for name, id := range m {
			if id != "" {
				named[name] = id
			}
		}
	}
	return named
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
