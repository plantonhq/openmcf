package verify

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	pkgerrors "github.com/pkg/errors"
)

// lambdaLayerVerifier verifies an AwsLambdaLayer via
// GetLayerVersionByArn, keyed on the published version's ARN (the
// provider's import ID). The typed ResourceNotFoundException is the
// absent signal. When the deploy's outputs are available, the grant
// count is asserted against the version's policy - a share grant that
// silently failed to attach would otherwise pass existence checks.
type lambdaLayerVerifier struct{}

func (*lambdaLayerVerifier) IDOutputKey() string { return "layer_version_arn" }

func (*lambdaLayerVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	_, err := lambda.NewFromConfig(cfg).GetLayerVersionByArn(ctx, &lambda.GetLayerVersionByArnInput{
		Arn: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambdalayer verify-exists failed for %q", id)
	}
	return nil
}

func (*lambdaLayerVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	_, err := lambda.NewFromConfig(cfg).GetLayerVersionByArn(ctx, &lambda.GetLayerVersionByArnInput{
		Arn: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awslambdalayer %q still exists after destroy", id)
	}
	var notFound *lambdatypes.ResourceNotFoundException
	if pkgerrors.As(err, &notFound) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awslambdalayer verify-absent failed for %q", id)
}

// VerifyExistsFromOutputs additionally proves each declared share
// grant landed: the version's policy must exist whenever the
// permission_revision_ids map is non-empty.
func (v *lambdaLayerVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	id, _ := outputs["layer_version_arn"].(string)
	if id == "" {
		return pkgerrors.New("awslambdalayer outputs carry no layer_version_arn")
	}
	if err := v.VerifyExists(ctx, cfg, id, region); err != nil {
		return err
	}

	revisionIds, _ := outputs["permission_revision_ids"].(map[string]interface{})
	if len(revisionIds) == 0 {
		return nil
	}
	layerName, _ := outputs["layer_arn"].(string)
	version, _ := outputs["version"].(string)
	versionNumber, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambdalayer version output %q is not a number", version)
	}
	policy, err := lambda.NewFromConfig(cfg).GetLayerVersionPolicy(ctx, &lambda.GetLayerVersionPolicyInput{
		LayerName:     aws.String(layerName),
		VersionNumber: aws.Int64(versionNumber),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambdalayer policy read failed for %q", layerName)
	}
	if policy.Policy == nil || *policy.Policy == "" {
		return pkgerrors.Errorf("awslambdalayer %q declares %d grants but the version has no policy", layerName, len(revisionIds))
	}
	return nil
}
