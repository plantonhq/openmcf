package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	pkgerrors "github.com/pkg/errors"
)

// sageMakerResourceNotFound covers the typed not-found the newer
// SageMaker MLOps APIs return (feature groups, pipelines, images,
// MLflow), alongside the family's ValidationException message class.
func sageMakerResourceNotFound(err error) bool {
	var notFound *sagemakertypes.ResourceNotFound
	if errors.As(err, &notFound) {
		return true
	}
	return sageMakerValidationNotFound(err)
}

// sageMakerFeatureGroupVerifier verifies an AwsSagemakerFeatureGroup
// via DescribeFeatureGroup, keyed on the feature_group_name output.
// Exists means Created; Deleting counts as absent.
type sageMakerFeatureGroupVerifier struct{}

func (*sageMakerFeatureGroupVerifier) IDOutputKey() string { return "feature_group_name" }

func (v *sageMakerFeatureGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeFeatureGroup(ctx, &sagemaker.DescribeFeatureGroupInput{
		FeatureGroupName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return pkgerrors.Errorf("awssagemakerfeaturegroup %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakerfeaturegroup verify-exists failed for %q", id)
	}
	if out.FeatureGroupStatus != sagemakertypes.FeatureGroupStatusCreated {
		reason := ""
		if out.FailureReason != nil {
			reason = " (" + *out.FailureReason + ")"
		}
		return pkgerrors.Errorf("awssagemakerfeaturegroup %q is %s, want Created%s", id, out.FeatureGroupStatus, reason)
	}
	return nil
}

func (v *sageMakerFeatureGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeFeatureGroup(ctx, &sagemaker.DescribeFeatureGroupInput{
		FeatureGroupName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakerfeaturegroup verify-absent failed for %q", id)
	}
	if out.FeatureGroupStatus == sagemakertypes.FeatureGroupStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakerfeaturegroup %q still exists after destroy (status %s)", id, out.FeatureGroupStatus)
}

// sageMakerModelRegistryVerifier verifies an AwsSagemakerModelRegistry
// via DescribeModelPackageGroup, keyed on the
// model_package_group_name output. Exists means Completed; Deleting
// counts as absent.
type sageMakerModelRegistryVerifier struct{}

func (*sageMakerModelRegistryVerifier) IDOutputKey() string { return "model_package_group_name" }

func (v *sageMakerModelRegistryVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeModelPackageGroup(ctx, &sagemaker.DescribeModelPackageGroupInput{
		ModelPackageGroupName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return pkgerrors.Errorf("awssagemakermodelregistry %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakermodelregistry verify-exists failed for %q", id)
	}
	if out.ModelPackageGroupStatus != sagemakertypes.ModelPackageGroupStatusCompleted {
		return pkgerrors.Errorf("awssagemakermodelregistry %q is %s, want Completed", id, out.ModelPackageGroupStatus)
	}
	return nil
}

func (v *sageMakerModelRegistryVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeModelPackageGroup(ctx, &sagemaker.DescribeModelPackageGroupInput{
		ModelPackageGroupName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakermodelregistry verify-absent failed for %q", id)
	}
	if out.ModelPackageGroupStatus == sagemakertypes.ModelPackageGroupStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakermodelregistry %q still exists after destroy (status %s)", id, out.ModelPackageGroupStatus)
}

// sageMakerPipelineVerifier verifies an AwsSagemakerPipeline via
// DescribePipeline, keyed on the pipeline_name output. Exists means
// Active (a green create IS the definition-validity claim - AWS
// validates the DAG server-side).
type sageMakerPipelineVerifier struct{}

func (*sageMakerPipelineVerifier) IDOutputKey() string { return "pipeline_name" }

func (v *sageMakerPipelineVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribePipeline(ctx, &sagemaker.DescribePipelineInput{
		PipelineName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return pkgerrors.Errorf("awssagemakerpipeline %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakerpipeline verify-exists failed for %q", id)
	}
	if out.PipelineStatus != sagemakertypes.PipelineStatusActive {
		return pkgerrors.Errorf("awssagemakerpipeline %q is %s, want Active", id, out.PipelineStatus)
	}
	return nil
}

func (v *sageMakerPipelineVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := sageMakerClient(cfg, region).DescribePipeline(ctx, &sagemaker.DescribePipelineInput{
		PipelineName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakerpipeline verify-absent failed for %q", id)
	}
	return pkgerrors.Errorf("awssagemakerpipeline %q still exists after destroy", id)
}

// sageMakerImageVerifier verifies an AwsSagemakerImage via
// DescribeImage, keyed on the image_name output. Exists means CREATED;
// DELETING counts as absent.
type sageMakerImageVerifier struct{}

func (*sageMakerImageVerifier) IDOutputKey() string { return "image_name" }

func (v *sageMakerImageVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeImage(ctx, &sagemaker.DescribeImageInput{
		ImageName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return pkgerrors.Errorf("awssagemakerimage %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakerimage verify-exists failed for %q", id)
	}
	if out.ImageStatus != sagemakertypes.ImageStatusCreated {
		reason := ""
		if out.FailureReason != nil {
			reason = " (" + *out.FailureReason + ")"
		}
		return pkgerrors.Errorf("awssagemakerimage %q is %s, want CREATED%s", id, out.ImageStatus, reason)
	}
	return nil
}

func (v *sageMakerImageVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeImage(ctx, &sagemaker.DescribeImageInput{
		ImageName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakerimage verify-absent failed for %q", id)
	}
	if out.ImageStatus == sagemakertypes.ImageStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakerimage %q still exists after destroy (status %s)", id, out.ImageStatus)
}

// sageMakerMlflowServerVerifier verifies an AwsSagemakerMlflowServer
// via DescribeMlflowTrackingServer, keyed on the tracking_server_name
// output. Exists means Created; Deleting counts as absent.
type sageMakerMlflowServerVerifier struct{}

func (*sageMakerMlflowServerVerifier) IDOutputKey() string { return "tracking_server_name" }

func (v *sageMakerMlflowServerVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeMlflowTrackingServer(ctx, &sagemaker.DescribeMlflowTrackingServerInput{
		TrackingServerName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return pkgerrors.Errorf("awssagemakermlflowserver %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakermlflowserver verify-exists failed for %q", id)
	}
	if out.TrackingServerStatus != sagemakertypes.TrackingServerStatusCreated {
		return pkgerrors.Errorf("awssagemakermlflowserver %q is %s, want Created", id, out.TrackingServerStatus)
	}
	return nil
}

func (v *sageMakerMlflowServerVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeMlflowTrackingServer(ctx, &sagemaker.DescribeMlflowTrackingServerInput{
		TrackingServerName: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakermlflowserver verify-absent failed for %q", id)
	}
	if out.TrackingServerStatus == sagemakertypes.TrackingServerStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakermlflowserver %q still exists after destroy (status %s)", id, out.TrackingServerStatus)
}

// sageMakerMlflowAppVerifier verifies an AwsSagemakerMlflowApp via
// DescribeMlflowApp, keyed on the app_arn output (the app's AWS
// identity). DELETED is a terminal soft-delete status - accepted as
// absent (the provider's own read semantics).
type sageMakerMlflowAppVerifier struct{}

func (*sageMakerMlflowAppVerifier) IDOutputKey() string { return "app_arn" }

func (v *sageMakerMlflowAppVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeMlflowApp(ctx, &sagemaker.DescribeMlflowAppInput{
		Arn: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return pkgerrors.Errorf("awssagemakermlflowapp %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakermlflowapp verify-exists failed for %q", id)
	}
	if out.Status != sagemakertypes.MlflowAppStatusCreated {
		return pkgerrors.Errorf("awssagemakermlflowapp %q is %s, want Created", id, out.Status)
	}
	return nil
}

func (v *sageMakerMlflowAppVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeMlflowApp(ctx, &sagemaker.DescribeMlflowAppInput{
		Arn: aws.String(id),
	})
	if err != nil {
		if sageMakerResourceNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakermlflowapp verify-absent failed for %q", id)
	}
	// DELETED is terminal soft-delete; Deleting is on its way out.
	if out.Status == sagemakertypes.MlflowAppStatusDeleted || out.Status == sagemakertypes.MlflowAppStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakermlflowapp %q still exists after destroy (status %s)", id, out.Status)
}
