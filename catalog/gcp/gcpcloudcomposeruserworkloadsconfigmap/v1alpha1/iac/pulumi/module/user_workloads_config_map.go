package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/composer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// userWorkloadsConfigMap provisions the Kubernetes ConfigMap Composer
// manages in the environment's GKE cluster. Airflow DAGs
// (KubernetesPodOperator tasks) consume it by name for non-secret
// configuration: feature flags, endpoints, tuning parameters.
//
// The data updates in place; name, environment, region, and project are
// immutable. The environment must already exist — it is a first-class
// resource this one composes against by reference.
//
// No API enablement here: the Composer API is enabled by the
// environment this ConfigMap is delivered into (a ConfigMap cannot
// exist without one).
func userWorkloadsConfigMap(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpCloudComposerUserWorkloadsConfigMap.Spec

	args := &composer.UserWorkloadsConfigMapArgs{
		Name:        pulumi.StringPtr(spec.ConfigMapName),
		Environment: pulumi.String(spec.Environment.GetValue()),
		Region:      pulumi.StringPtr(spec.Region),
		Data:        pulumi.ToStringMap(spec.Data),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	createdConfigMap, err := composer.NewUserWorkloadsConfigMap(ctx, "user-workloads-config-map", args,
		pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create user workloads config map")
	}

	// The resource ID is the fully qualified resource name — the bridged
	// provider inherits Terraform's ID format, identical to the Terraform
	// module's output.
	ctx.Export(OpName, createdConfigMap.ID())
	ctx.Export(OpConfigMapName, createdConfigMap.Name)

	return nil
}
