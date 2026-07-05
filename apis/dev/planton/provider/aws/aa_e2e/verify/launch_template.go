package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// launchTemplateVerifier verifies an AwsLaunchTemplate via
// DescribeLaunchTemplates, keyed on the template ID. A deleted template
// returns a typed InvalidLaunchTemplateId.* error (AWS uses .NotFound and
// .NotFoundException variants across API versions, so the family prefix is
// the "absent" signal); any other error is a genuine failure and must
// surface.
type launchTemplateVerifier struct{}

func (*launchTemplateVerifier) IDOutputKey() string { return "launch_template_id" }

func (*launchTemplateVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := launchTemplateExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslaunchtemplate verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awslaunchtemplate %q not found after deploy", id)
	}
	return nil
}

func (*launchTemplateVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := launchTemplateExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslaunchtemplate verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awslaunchtemplate %q still exists after destroy", id)
	}
	return nil
}

func launchTemplateExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeLaunchTemplates(ctx, &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && strings.HasPrefix(apiErr.ErrorCode(), "InvalidLaunchTemplateId.") {
			return false, nil
		}
		return false, err
	}
	return len(out.LaunchTemplates) > 0, nil
}
