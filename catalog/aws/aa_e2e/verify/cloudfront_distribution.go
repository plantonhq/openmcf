package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	pkgerrors "github.com/pkg/errors"
)

// cloudFrontDistributionVerifier verifies an AwsCloudFront via
// GetDistribution, keyed on the distribution_id output. CloudFront is a
// global service (region is irrelevant to the API call). A distribution
// that is still propagating reports status "InProgress" -- it exists and is
// converging, so existence is simply "gettable"; the modules' own
// wait_for_deployment handles propagation. Deletion is terminal: once the
// disable-then-delete dance completes, GetDistribution returns NoSuchDistribution.
//
// As an OutputsVerifier it also asserts the stack's own status claim
// against the live distribution: wait_for_deployment is spec-driven, so
// "Deployed" is not an unconditional invariant -- but when the stack's
// `status` output SAYS Deployed, the live distribution must report
// Deployed too (propagation never regresses), and a stack that exported
// InProgress under the default wait is a waiter defect worth failing on.
type cloudFrontDistributionVerifier struct{}

func (*cloudFrontDistributionVerifier) IDOutputKey() string { return "distribution_id" }

func (v *cloudFrontDistributionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	dist, err := cloudFrontDistributionGet(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudfront verify-exists failed for %q", id)
	}
	if dist == nil {
		return pkgerrors.Errorf("awscloudfront %q not found after deploy", id)
	}
	return nil
}

func (*cloudFrontDistributionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	dist, err := cloudFrontDistributionGet(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudfront verify-absent failed for %q", id)
	}
	if dist != nil {
		return pkgerrors.Errorf("awscloudfront %q still exists after destroy", id)
	}
	return nil
}

// VerifyExistsFromOutputs additionally cross-checks the exported status and
// domain name against the distribution's own state.
func (v *cloudFrontDistributionVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	id, _ := outputs[v.IDOutputKey()].(string)
	if id == "" {
		return pkgerrors.Errorf("awscloudfront outputs carry no %q", v.IDOutputKey())
	}
	dist, err := cloudFrontDistributionGet(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudfront verify-exists failed for %q", id)
	}
	if dist == nil {
		return pkgerrors.Errorf("awscloudfront %q not found after deploy", id)
	}
	if claimed, _ := outputs["status"].(string); claimed == "Deployed" {
		if live := aws.ToString(dist.Status); live != "Deployed" {
			return pkgerrors.Errorf("awscloudfront %q: status output claims Deployed but the distribution reports %q", id, live)
		}
	}
	if wantDomain, _ := outputs["domain_name"].(string); wantDomain != "" {
		if live := aws.ToString(dist.DomainName); live != wantDomain {
			return pkgerrors.Errorf("awscloudfront %q: domain_name output %q does not match the distribution's %q", id, wantDomain, live)
		}
	}
	return nil
}

func (v *cloudFrontDistributionVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	id, _ := outputs[v.IDOutputKey()].(string)
	if id == "" {
		return pkgerrors.Errorf("awscloudfront outputs carry no %q", v.IDOutputKey())
	}
	return v.VerifyAbsent(ctx, cfg, id, region)
}

// cloudFrontDistributionGet returns the distribution, or nil when AWS
// reports NoSuchDistribution.
func cloudFrontDistributionGet(ctx context.Context, cfg aws.Config, id string) (*types.Distribution, error) {
	client := cloudfront.NewFromConfig(cfg)
	out, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{Id: &id})
	if err != nil {
		var notFound *types.NoSuchDistribution
		if errors.As(err, &notFound) {
			return nil, nil
		}
		return nil, err
	}
	return out.Distribution, nil
}
