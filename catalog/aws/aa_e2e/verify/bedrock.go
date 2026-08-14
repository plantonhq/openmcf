package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrock/types"
	pkgerrors "github.com/pkg/errors"
)

func bedrockClient(cfg aws.Config, region string) *bedrock.Client {
	return bedrock.NewFromConfig(cfg, func(o *bedrock.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// bedrockGuardrailVerifier verifies an AwsBedrockGuardrail via
// GetGuardrail, keyed on guardrail_id (the DRAFT view). Exists requires
// the guardrail past its brief CREATING window (READY, or VERSIONING
// during a same-deploy publish).
type bedrockGuardrailVerifier struct{}

func (*bedrockGuardrailVerifier) IDOutputKey() string { return "guardrail_id" }

func (*bedrockGuardrailVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockClient(cfg, region).GetGuardrail(ctx, &bedrock.GetGuardrailInput{
		GuardrailIdentifier: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbedrockguardrail verify-exists failed for %q", id)
	}
	if out.Status == types.GuardrailStatusFailed || out.Status == types.GuardrailStatusDeleting {
		return pkgerrors.Errorf("awsbedrockguardrail %q in unexpected status %s after deploy", id, out.Status)
	}
	return nil
}

func (*bedrockGuardrailVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockClient(cfg, region).GetGuardrail(ctx, &bedrock.GetGuardrailInput{
		GuardrailIdentifier: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awsbedrockguardrail %q still exists after destroy", id)
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awsbedrockguardrail verify-absent failed for %q", id)
}

// bedrockCustomModelVerifier verifies an AwsBedrockCustomModel via
// GetModelCustomizationJob, keyed on job_arn (the job is the tracked
// object; the model materializes only when the job completes). A Failed
// job after deploy is a verification failure -- apply success alone never
// proves the training outcome.
type bedrockCustomModelVerifier struct{}

func (*bedrockCustomModelVerifier) IDOutputKey() string { return "job_arn" }

func (*bedrockCustomModelVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockClient(cfg, region).GetModelCustomizationJob(ctx, &bedrock.GetModelCustomizationJobInput{
		JobIdentifier: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbedrockcustommodel verify-exists failed for %q", id)
	}
	if out.Status == types.ModelCustomizationJobStatusFailed {
		return pkgerrors.Errorf("awsbedrockcustommodel job %q Failed after deploy", id)
	}
	return nil
}

func (*bedrockCustomModelVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	// The job RECORD survives model deletion (job history is immutable);
	// absence means the custom model the job produced is gone. A job still
	// describable with its output model deleted is the destroyed state.
	out, err := bedrockClient(cfg, region).GetModelCustomizationJob(ctx, &bedrock.GetModelCustomizationJobInput{
		JobIdentifier: aws.String(id),
	})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awsbedrockcustommodel verify-absent failed for %q", id)
	}
	if out.OutputModelArn == nil {
		return nil
	}
	if _, err := bedrockClient(cfg, region).GetCustomModel(ctx, &bedrock.GetCustomModelInput{
		ModelIdentifier: out.OutputModelArn,
	}); err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awsbedrockcustommodel verify-absent model check failed for %q", id)
	}
	return pkgerrors.Errorf("awsbedrockcustommodel %q custom model still exists after destroy", id)
}

// bedrockInferenceProfileVerifier verifies an AwsBedrockInferenceProfile
// via GetInferenceProfile, keyed on inference_profile_id (ACTIVE, type
// APPLICATION).
type bedrockInferenceProfileVerifier struct{}

func (*bedrockInferenceProfileVerifier) IDOutputKey() string { return "inference_profile_id" }

func (*bedrockInferenceProfileVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockClient(cfg, region).GetInferenceProfile(ctx, &bedrock.GetInferenceProfileInput{
		InferenceProfileIdentifier: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbedrockinferenceprofile verify-exists failed for %q", id)
	}
	if out.Status != types.InferenceProfileStatusActive {
		return pkgerrors.Errorf("awsbedrockinferenceprofile %q not ACTIVE after deploy (status %s)", id, out.Status)
	}
	if out.Type != types.InferenceProfileTypeApplication {
		return pkgerrors.Errorf("awsbedrockinferenceprofile %q has type %s, want APPLICATION", id, out.Type)
	}
	return nil
}

func (*bedrockInferenceProfileVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockClient(cfg, region).GetInferenceProfile(ctx, &bedrock.GetInferenceProfileInput{
		InferenceProfileIdentifier: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awsbedrockinferenceprofile %q still exists after destroy", id)
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awsbedrockinferenceprofile verify-absent failed for %q", id)
}

// bedrockProvisionedThroughputVerifier verifies an
// AwsBedrockProvisionedThroughput via GetProvisionedModelThroughput,
// keyed on provisioned_model_arn.
type bedrockProvisionedThroughputVerifier struct{}

func (*bedrockProvisionedThroughputVerifier) IDOutputKey() string { return "provisioned_model_arn" }

func (*bedrockProvisionedThroughputVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockClient(cfg, region).GetProvisionedModelThroughput(ctx, &bedrock.GetProvisionedModelThroughputInput{
		ProvisionedModelId: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbedrockprovisionedthroughput verify-exists failed for %q", id)
	}
	if out.Status == types.ProvisionedModelStatusFailed {
		return pkgerrors.Errorf("awsbedrockprovisionedthroughput %q Failed after deploy", id)
	}
	return nil
}

func (*bedrockProvisionedThroughputVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := bedrockClient(cfg, region).GetProvisionedModelThroughput(ctx, &bedrock.GetProvisionedModelThroughputInput{
		ProvisionedModelId: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awsbedrockprovisionedthroughput %q still exists after destroy", id)
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awsbedrockprovisionedthroughput verify-absent failed for %q", id)
}

// bedrockModelAccessVerifier verifies an AwsBedrockModelAccess via
// GetFoundationModelAvailability, keyed on model_id: the agreement's
// availability status is AVAILABLE after deploy and NOT_AVAILABLE (or the
// model gone entirely) after destroy -- the create waiter's contract,
// checked in both directions.
type bedrockModelAccessVerifier struct{}

func (*bedrockModelAccessVerifier) IDOutputKey() string { return "model_id" }

func (*bedrockModelAccessVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockClient(cfg, region).GetFoundationModelAvailability(ctx, &bedrock.GetFoundationModelAvailabilityInput{
		ModelId: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbedrockmodelaccess verify-exists failed for %q", id)
	}
	if out.AgreementAvailability == nil || out.AgreementAvailability.Status != types.AgreementStatusAvailable {
		status := "unknown"
		if out.AgreementAvailability != nil {
			status = string(out.AgreementAvailability.Status)
		}
		return pkgerrors.Errorf("awsbedrockmodelaccess %q agreement not AVAILABLE after deploy (status %s)", id, status)
	}
	return nil
}

func (*bedrockModelAccessVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrockClient(cfg, region).GetFoundationModelAvailability(ctx, &bedrock.GetFoundationModelAvailabilityInput{
		ModelId: aws.String(id),
	})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awsbedrockmodelaccess verify-absent failed for %q", id)
	}
	if out.AgreementAvailability != nil && out.AgreementAvailability.Status == types.AgreementStatusAvailable {
		return pkgerrors.Errorf("awsbedrockmodelaccess %q agreement still AVAILABLE after destroy", id)
	}
	return nil
}
