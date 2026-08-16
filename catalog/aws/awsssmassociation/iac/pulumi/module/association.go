package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// association creates the State Manager association and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - the document reference (Name) forces replacement; every other
//     change creates a new association version in place, and the
//     provider sends the FULL argument set on update (AWS versions
//     associations whole);
//   - AWS identifies the association by a generated UUID, not a name;
//   - AWS materializes the document's declared parameter defaults into
//     the parameters map server-side, and
//     wait_for_success_timeout_seconds is a create-time wait never
//     read back - the import map declares both accordingly;
//   - the association name (association_name) is display metadata, not
//     identity.
func association(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &ssm.AssociationArgs{
		Name: pulumi.String(spec.DocumentName.GetValue()),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.AssociationName != "" {
		args.AssociationName = pulumi.String(spec.AssociationName)
	}
	if spec.DocumentVersion != "" {
		args.DocumentVersion = pulumi.String(spec.DocumentVersion)
	}
	if len(spec.Parameters) > 0 {
		args.Parameters = pulumi.ToStringMap(spec.Parameters)
	}

	var targets ssm.AssociationTargetArray
	for _, t := range spec.Targets {
		targets = append(targets, &ssm.AssociationTargetArgs{
			Key:    pulumi.String(t.Key),
			Values: pulumi.ToStringArray(t.Values),
		})
	}
	if len(targets) > 0 {
		args.Targets = targets
	}

	if spec.ScheduleExpression != "" {
		args.ScheduleExpression = pulumi.String(spec.ScheduleExpression)
	}
	if spec.ApplyOnlyAtCronInterval {
		args.ApplyOnlyAtCronInterval = pulumi.Bool(true)
	}
	if spec.ComplianceSeverity != "" {
		args.ComplianceSeverity = pulumi.String(spec.ComplianceSeverity)
	}
	if spec.SyncCompliance != "" {
		args.SyncCompliance = pulumi.String(spec.SyncCompliance)
	}
	if spec.MaxConcurrency != "" {
		args.MaxConcurrency = pulumi.String(spec.MaxConcurrency)
	}
	if spec.MaxErrors != "" {
		args.MaxErrors = pulumi.String(spec.MaxErrors)
	}
	if spec.AutomationTargetParameterName != "" {
		args.AutomationTargetParameterName = pulumi.String(spec.AutomationTargetParameterName)
	}
	if len(spec.CalendarNames) > 0 {
		args.CalendarNames = pulumi.ToStringArray(spec.CalendarNames)
	}

	if spec.OutputLocation != nil {
		outputLocationArgs := &ssm.AssociationOutputLocationArgs{
			S3BucketName: pulumi.String(spec.OutputLocation.S3BucketName.GetValue()),
		}
		if spec.OutputLocation.S3KeyPrefix != "" {
			outputLocationArgs.S3KeyPrefix = pulumi.String(spec.OutputLocation.S3KeyPrefix)
		}
		if spec.OutputLocation.S3Region != "" {
			outputLocationArgs.S3Region = pulumi.String(spec.OutputLocation.S3Region)
		}
		args.OutputLocation = outputLocationArgs
	}

	if spec.WaitForSuccessTimeoutSeconds != 0 {
		args.WaitForSuccessTimeoutSeconds = pulumi.Int(int(spec.WaitForSuccessTimeoutSeconds))
	}

	createdAssociation, err := ssm.NewAssociation(ctx, "association", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create association")
	}

	ctx.Export(OpAssociationId, createdAssociation.AssociationId)
	ctx.Export(OpAssociationArn, createdAssociation.Arn)
	ctx.Export(OpDocumentName, createdAssociation.Name)
	return nil
}
