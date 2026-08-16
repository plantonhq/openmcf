package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// trackingServer creates the MLflow tracking server and exports
// outputs.
//
// Lifecycle facts the renders below depend on:
//   - create and delete each take ~25 minutes (the provider's own
//     timeouts are 45m per operation) - budget accordingly;
//   - the server bills hourly from Created onward (Small ~$0.6/hour);
//   - automatic_model_registration CANNOT be turned back off through
//     the provider (a true-to-false change is silently not transmitted
//   - an upstream update-guard gap taught on the spec field); the
//     module always renders the spec value so the intent is visible in
//     the preview;
//   - AWS normalizes mlflow_version to major.minor (the spec's pattern
//     already forbids patch-level values, so no drift).
func trackingServer(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &sagemaker.MlflowTrackingServerArgs{
		// The component's name IS the tracking server name.
		TrackingServerName: pulumi.String(locals.TrackingServerName),
		ArtifactStoreUri:   pulumi.String(spec.ArtifactStoreUri),
		// Changing the role replaces the server (provider-enforced).
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
		Tags:    pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Size != "" {
		args.TrackingServerSize = pulumi.String(spec.Size)
	}
	// Changing the version replaces the server (provider-enforced).
	if spec.MlflowVersion != "" {
		args.MlflowVersion = pulumi.String(spec.MlflowVersion)
	}
	if spec.AutomaticModelRegistration {
		args.AutomaticModelRegistration = pulumi.Bool(true)
	}
	if spec.WeeklyMaintenanceWindowStart != "" {
		args.WeeklyMaintenanceWindowStart = pulumi.String(spec.WeeklyMaintenanceWindowStart)
	}

	createdServer, err := sagemaker.NewMlflowTrackingServer(ctx, locals.TrackingServerName, args,
		pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create mlflow tracking server")
	}

	ctx.Export(OpTrackingServerName, createdServer.TrackingServerName)
	ctx.Export(OpTrackingServerArn, createdServer.Arn)
	ctx.Export(OpTrackingServerUrl, createdServer.TrackingServerUrl)

	return nil
}
