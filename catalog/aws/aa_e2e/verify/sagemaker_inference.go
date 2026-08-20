package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// SageMaker's control plane reports missing serving resources as
// ValidationException with a "Could not find ..." (or "RecordNotFound")
// message rather than a typed not-found error - the family's shared
// absence check.
func sageMakerValidationNotFound(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ValidationException" {
		message := apiErr.ErrorMessage()
		return strings.Contains(message, "Could not find") || strings.Contains(message, "RecordNotFound") || strings.Contains(message, "does not exist")
	}
	return false
}

func sageMakerClient(cfg aws.Config, region string) *sagemaker.Client {
	return sagemaker.NewFromConfig(cfg, func(o *sagemaker.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// sageMakerModelVerifier verifies an AwsSagemakerModel via
// DescribeModel, keyed on the model_name output. Models have no status
// machine - describable IS created.
type sageMakerModelVerifier struct{}

func (*sageMakerModelVerifier) IDOutputKey() string { return "model_name" }

func (v *sageMakerModelVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := sageMakerClient(cfg, region).DescribeModel(ctx, &sagemaker.DescribeModelInput{ModelName: aws.String(id)})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return pkgerrors.Errorf("awssagemakermodel %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakermodel verify-exists failed for %q", id)
	}
	return nil
}

func (v *sageMakerModelVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := sageMakerClient(cfg, region).DescribeModel(ctx, &sagemaker.DescribeModelInput{ModelName: aws.String(id)})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakermodel verify-absent failed for %q", id)
	}
	return pkgerrors.Errorf("awssagemakermodel %q still exists after destroy", id)
}

// sageMakerEndpointVerifier verifies an AwsSagemakerEndpoint via
// DescribeEndpoint, keyed on the endpoint_name output. Exists means
// InService (Failed surfaces the reason); an endpoint in Deleting
// counts as absent (the provider's own read semantics).
type sageMakerEndpointVerifier struct{}

func (*sageMakerEndpointVerifier) IDOutputKey() string { return "endpoint_name" }

func (v *sageMakerEndpointVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeEndpoint(ctx, &sagemaker.DescribeEndpointInput{EndpointName: aws.String(id)})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return pkgerrors.Errorf("awssagemakerendpoint %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakerendpoint verify-exists failed for %q", id)
	}
	if out.EndpointStatus != sagemakertypes.EndpointStatusInService {
		reason := ""
		if out.FailureReason != nil {
			reason = " (" + *out.FailureReason + ")"
		}
		return pkgerrors.Errorf("awssagemakerendpoint %q is %s, want InService%s", id, out.EndpointStatus, reason)
	}
	return nil
}

func (v *sageMakerEndpointVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeEndpoint(ctx, &sagemaker.DescribeEndpointInput{EndpointName: aws.String(id)})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakerendpoint verify-absent failed for %q", id)
	}
	// Deletion is asynchronous - Deleting is on its way out.
	if out.EndpointStatus == sagemakertypes.EndpointStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakerendpoint %q still exists after destroy (status %s)", id, out.EndpointStatus)
}

// VerifyExistsFromOutputs raises the bar to the folded configuration:
// the suffixed endpoint configuration in the endpoint_config_name
// output must be describable and referenced by the live endpoint.
func (v *sageMakerEndpointVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	endpointName := stringOutputMap(outputs, "endpoint_name")
	if endpointName == "" {
		return pkgerrors.New("awssagemakerendpoint verify-exists: no endpoint_name in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, endpointName, region); err != nil {
		return err
	}

	configName := stringOutputMap(outputs, "endpoint_config_name")
	if configName == "" {
		return pkgerrors.New("awssagemakerendpoint verify-exists: no endpoint_config_name in outputs")
	}
	client := sageMakerClient(cfg, region)
	if _, err := client.DescribeEndpointConfig(ctx, &sagemaker.DescribeEndpointConfigInput{
		EndpointConfigName: aws.String(configName),
	}); err != nil {
		return pkgerrors.Wrapf(err, "awssagemakerendpoint %q: configuration %q not describable", endpointName, configName)
	}
	endpoint, err := client.DescribeEndpoint(ctx, &sagemaker.DescribeEndpointInput{EndpointName: aws.String(endpointName)})
	if err != nil {
		return pkgerrors.Wrapf(err, "awssagemakerendpoint verify-exists failed for %q", endpointName)
	}
	if aws.ToString(endpoint.EndpointConfigName) != configName {
		return pkgerrors.Errorf("awssagemakerendpoint %q references configuration %q, want %q",
			endpointName, aws.ToString(endpoint.EndpointConfigName), configName)
	}
	return nil
}

// VerifyAbsentFromOutputs checks both the endpoint and its folded
// configuration are gone.
func (v *sageMakerEndpointVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	endpointName := stringOutputMap(outputs, "endpoint_name")
	if endpointName == "" {
		return pkgerrors.New("awssagemakerendpoint verify-absent: no endpoint_name in outputs")
	}
	if err := v.VerifyAbsent(ctx, cfg, endpointName, region); err != nil {
		return err
	}
	configName := stringOutputMap(outputs, "endpoint_config_name")
	if configName == "" {
		return nil
	}
	_, err := sageMakerClient(cfg, region).DescribeEndpointConfig(ctx, &sagemaker.DescribeEndpointConfigInput{
		EndpointConfigName: aws.String(configName),
	})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakerendpoint verify-absent failed for configuration %q", configName)
	}
	return pkgerrors.Errorf("awssagemakerendpoint configuration %q still exists after destroy", configName)
}

// sageMakerNotebookInstanceVerifier verifies an
// AwsSagemakerNotebookInstance via DescribeNotebookInstance, keyed on
// the notebook_instance_name output. Exists means InService; Stopping/
// Deleting count as absent (deletion stops first).
type sageMakerNotebookInstanceVerifier struct{}

func (*sageMakerNotebookInstanceVerifier) IDOutputKey() string { return "notebook_instance_name" }

func (v *sageMakerNotebookInstanceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeNotebookInstance(ctx, &sagemaker.DescribeNotebookInstanceInput{
		NotebookInstanceName: aws.String(id),
	})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return pkgerrors.Errorf("awssagemakernotebookinstance %q not found after deploy", id)
		}
		return pkgerrors.Wrapf(err, "awssagemakernotebookinstance verify-exists failed for %q", id)
	}
	if out.NotebookInstanceStatus != sagemakertypes.NotebookInstanceStatusInService {
		reason := ""
		if out.FailureReason != nil {
			reason = " (" + *out.FailureReason + ")"
		}
		return pkgerrors.Errorf("awssagemakernotebookinstance %q is %s, want InService%s", id, out.NotebookInstanceStatus, reason)
	}
	return nil
}

func (v *sageMakerNotebookInstanceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sageMakerClient(cfg, region).DescribeNotebookInstance(ctx, &sagemaker.DescribeNotebookInstanceInput{
		NotebookInstanceName: aws.String(id),
	})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakernotebookinstance verify-absent failed for %q", id)
	}
	if out.NotebookInstanceStatus == sagemakertypes.NotebookInstanceStatusDeleting {
		return nil
	}
	return pkgerrors.Errorf("awssagemakernotebookinstance %q still exists after destroy (status %s)", id, out.NotebookInstanceStatus)
}

// VerifyExistsFromOutputs raises the bar to the folded lifecycle
// configuration when the lifecycle_config_name output carries one.
func (v *sageMakerNotebookInstanceVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	instanceName := stringOutputMap(outputs, "notebook_instance_name")
	if instanceName == "" {
		return pkgerrors.New("awssagemakernotebookinstance verify-exists: no notebook_instance_name in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, instanceName, region); err != nil {
		return err
	}
	lifecycleName := stringOutputMap(outputs, "lifecycle_config_name")
	if lifecycleName == "" {
		return nil
	}
	if _, err := sageMakerClient(cfg, region).DescribeNotebookInstanceLifecycleConfig(ctx,
		&sagemaker.DescribeNotebookInstanceLifecycleConfigInput{
			NotebookInstanceLifecycleConfigName: aws.String(lifecycleName),
		}); err != nil {
		return pkgerrors.Wrapf(err, "awssagemakernotebookinstance %q: lifecycle config %q not describable", instanceName, lifecycleName)
	}
	return nil
}

// VerifyAbsentFromOutputs checks the instance and its folded lifecycle
// configuration are gone.
func (v *sageMakerNotebookInstanceVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	instanceName := stringOutputMap(outputs, "notebook_instance_name")
	if instanceName == "" {
		return pkgerrors.New("awssagemakernotebookinstance verify-absent: no notebook_instance_name in outputs")
	}
	if err := v.VerifyAbsent(ctx, cfg, instanceName, region); err != nil {
		return err
	}
	lifecycleName := stringOutputMap(outputs, "lifecycle_config_name")
	if lifecycleName == "" {
		return nil
	}
	_, err := sageMakerClient(cfg, region).DescribeNotebookInstanceLifecycleConfig(ctx,
		&sagemaker.DescribeNotebookInstanceLifecycleConfigInput{
			NotebookInstanceLifecycleConfigName: aws.String(lifecycleName),
		})
	if err != nil {
		if sageMakerValidationNotFound(err) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awssagemakernotebookinstance verify-absent failed for lifecycle config %q", lifecycleName)
	}
	return pkgerrors.Errorf("awssagemakernotebookinstance lifecycle config %q still exists after destroy", lifecycleName)
}
