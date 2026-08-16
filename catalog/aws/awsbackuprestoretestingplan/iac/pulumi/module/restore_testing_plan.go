package module

import (
	"github.com/pkg/errors"
	awsbackuprestoretestingplanv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbackuprestoretestingplan/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/backup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// restoreTestingPlan creates the restore testing plan with its folded
// selections and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - AWS restore testing names forbid hyphens and periods, so the
//     names are spec.plan_name / selection.name (explicit fields),
//     never metadata.name;
//   - a selection covers resources by EXACTLY ONE of explicit ARNs or
//     tag conditions (the provider enforces the same rule
//     resource-wide; the spec's CEL matches it);
//   - AWS returns empty condition lists as present-but-empty; the
//     provider collapses that to absent on both create and read - the
//     module renders the block only when the spec carries conditions;
//   - several Optional+Computed knobs (timezone, start window,
//     selection window, validation window, exclude vaults) keep an
//     AWS-side value once set - they cannot be cleared back to unset.
func restoreTestingPlan(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	rps := spec.RecoveryPointSelection
	selectionArgs := &backup.RestoreTestingPlanRecoveryPointSelectionArgs{
		Algorithm:          pulumi.String(rps.Algorithm),
		IncludeVaults:      pulumi.ToStringArray(rps.IncludeVaults),
		RecoveryPointTypes: pulumi.ToStringArray(rps.RecoveryPointTypes),
	}
	if len(rps.ExcludeVaults) > 0 {
		selectionArgs.ExcludeVaults = pulumi.ToStringArray(rps.ExcludeVaults)
	}
	if rps.SelectionWindowDays != 0 {
		selectionArgs.SelectionWindowDays = pulumi.Int(int(rps.SelectionWindowDays))
	}

	args := &backup.RestoreTestingPlanArgs{
		Name:                   pulumi.String(spec.PlanName),
		ScheduleExpression:     pulumi.String(spec.ScheduleExpression),
		RecoveryPointSelection: selectionArgs,
		Tags:                   pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.ScheduleExpressionTimezone != "" {
		args.ScheduleExpressionTimezone = pulumi.String(spec.ScheduleExpressionTimezone)
	}
	if spec.StartWindowHours != 0 {
		args.StartWindowHours = pulumi.Int(int(spec.StartWindowHours))
	}

	createdPlan, err := backup.NewRestoreTestingPlan(ctx, "restore-testing-plan", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create restore testing plan")
	}

	// Folded per-type selections, keyed by name.
	for _, sel := range spec.Selections {
		if err := testingSelection(ctx, createdPlan, sel, provider); err != nil {
			return errors.Wrapf(err, "selection %s", sel.Name)
		}
	}

	ctx.Export(OpRestoreTestingPlanArn, createdPlan.Arn)
	return nil
}

// testingSelection creates one folded restore testing selection under
// the plan.
func testingSelection(ctx *pulumi.Context, createdPlan *backup.RestoreTestingPlan,
	sel *awsbackuprestoretestingplanv1alpha1.AwsBackupRestoreTestingPlanSelection, provider *aws.Provider) error {

	args := &backup.RestoreTestingSelectionArgs{
		Name:                   pulumi.String(sel.Name),
		RestoreTestingPlanName: createdPlan.Name,
		ProtectedResourceType:  pulumi.String(sel.ProtectedResourceType),
		IamRoleArn:             pulumi.String(sel.IamRoleArn.GetValue()),
	}
	if len(sel.ProtectedResourceArns) > 0 {
		args.ProtectedResourceArns = pulumi.ToStringArray(sel.ProtectedResourceArns)
	}
	if sel.ProtectedResourceConditions != nil {
		conditionsArgs := &backup.RestoreTestingSelectionProtectedResourceConditionsArgs{}
		var stringEquals backup.RestoreTestingSelectionProtectedResourceConditionsStringEqualArray
		for _, p := range sel.ProtectedResourceConditions.StringEquals {
			stringEquals = append(stringEquals, &backup.RestoreTestingSelectionProtectedResourceConditionsStringEqualArgs{
				Key: pulumi.String(p.Key), Value: pulumi.String(p.Value),
			})
		}
		if len(stringEquals) > 0 {
			conditionsArgs.StringEquals = stringEquals
		}
		var stringNotEquals backup.RestoreTestingSelectionProtectedResourceConditionsStringNotEqualArray
		for _, p := range sel.ProtectedResourceConditions.StringNotEquals {
			stringNotEquals = append(stringNotEquals, &backup.RestoreTestingSelectionProtectedResourceConditionsStringNotEqualArgs{
				Key: pulumi.String(p.Key), Value: pulumi.String(p.Value),
			})
		}
		if len(stringNotEquals) > 0 {
			conditionsArgs.StringNotEquals = stringNotEquals
		}
		args.ProtectedResourceConditions = conditionsArgs
	}
	if len(sel.RestoreMetadataOverrides) > 0 {
		args.RestoreMetadataOverrides = pulumi.ToStringMap(sel.RestoreMetadataOverrides)
	}
	if sel.ValidationWindowHours != 0 {
		args.ValidationWindowHours = pulumi.Int(int(sel.ValidationWindowHours))
	}

	_, err := backup.NewRestoreTestingSelection(ctx, "selection-"+sel.Name, args,
		pulumi.Provider(provider), pulumi.Parent(createdPlan))
	return err
}
