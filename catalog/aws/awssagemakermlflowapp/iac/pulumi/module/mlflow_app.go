package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// mlflowApp creates the serverless MLflow app and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - the ARN is the app's identity (all API operations key on it); the
//     name updates in place;
//   - role_arn is the ONE replace-on-change argument
//     (provider-enforced);
//   - the app is standalone - it associates with SageMaker DOMAINS, not
//     with tracking servers;
//   - a soft-deleted app (status DELETED) reads as absent upstream.
func mlflowApp(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.MlflowAppArgs{
		// The component's name IS the app name.
		Name:             pulumi.String(locals.AppName),
		ArtifactStoreUri: pulumi.String(spec.ArtifactStoreUri),
		// The one replace-on-change argument.
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.AccountDefaultStatus != "" {
		args.AccountDefaultStatus = pulumi.String(spec.AccountDefaultStatus)
	}
	if len(spec.DefaultDomainIds) > 0 {
		var domainIds pulumi.StringArray
		for _, d := range spec.DefaultDomainIds {
			domainIds = append(domainIds, pulumi.String(d.GetValue()))
		}
		// The bridged field name pluralizes the provider's
		// default_domain_id_list.
		args.DefaultDomainIdLists = domainIds
	}
	if spec.ModelRegistrationMode != "" {
		args.ModelRegistrationMode = pulumi.String(spec.ModelRegistrationMode)
	}
	if spec.WeeklyMaintenanceWindowStart != "" {
		args.WeeklyMaintenanceWindowStart = pulumi.String(spec.WeeklyMaintenanceWindowStart)
	}

	createdApp, err := sagemaker.NewMlflowApp(ctx, locals.AppName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create mlflow app")
	}

	ctx.Export(OpAppArn, createdApp.Arn)
	ctx.Export(OpAppName, createdApp.Name)

	return nil
}
