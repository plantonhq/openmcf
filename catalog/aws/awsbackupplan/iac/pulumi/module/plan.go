package module

import (
	"github.com/pkg/errors"
	awsbackupplanv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbackupplan/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/backup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// plan creates the backup plan with its rules and folded selections,
// and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the plan's identity at AWS is a generated UUID, not the name;
//     the name forces replacement;
//   - the provider CANNOT send an explicit zero for the lifecycle day
//     counts (zero is dropped as unset) - the spec presence-types them
//     so the truth is explicit;
//   - opt_in_to_archive_for_supported_resources is transmitted only
//     when true (the provider never sends an explicit false);
//   - selections are fully ForceNew (no update path) and AWS refuses
//     to delete a plan while selections exist - the provider retries
//     the plan delete while the folded selections drain.
func plan(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var rules backup.PlanRuleArray
	for _, r := range spec.Rules {
		ruleArgs := &backup.PlanRuleArgs{
			RuleName:        pulumi.String(r.Name),
			TargetVaultName: pulumi.String(r.TargetVaultName.GetValue()),
		}
		// Rendered only on an explicit choice so the module never
		// fights the provider defaults (timezone Etc/UTC, start window
		// 60, completion window 180).
		if r.Schedule != "" {
			ruleArgs.Schedule = pulumi.String(r.Schedule)
		}
		if r.ScheduleExpressionTimezone != "" {
			ruleArgs.ScheduleExpressionTimezone = pulumi.String(r.ScheduleExpressionTimezone)
		}
		if r.StartWindowMinutes != 0 {
			ruleArgs.StartWindow = pulumi.Int(int(r.StartWindowMinutes))
		}
		if r.CompletionWindowMinutes != 0 {
			ruleArgs.CompletionWindow = pulumi.Int(int(r.CompletionWindowMinutes))
		}
		if r.EnableContinuousBackup {
			ruleArgs.EnableContinuousBackup = pulumi.Bool(true)
		}
		if len(r.RecoveryPointTags) > 0 {
			ruleArgs.RecoveryPointTags = pulumi.ToStringMap(r.RecoveryPointTags)
		}
		if r.TargetLogicallyAirGappedBackupVaultArn.GetValue() != "" {
			ruleArgs.TargetLogicallyAirGappedBackupVaultArn = pulumi.String(r.TargetLogicallyAirGappedBackupVaultArn.GetValue())
		}
		if r.Lifecycle != nil {
			lifecycleArgs := &backup.PlanRuleLifecycleArgs{}
			if r.Lifecycle.ColdStorageAfterDays != nil {
				lifecycleArgs.ColdStorageAfter = pulumi.Int(int(*r.Lifecycle.ColdStorageAfterDays))
			}
			if r.Lifecycle.DeleteAfterDays != nil {
				lifecycleArgs.DeleteAfter = pulumi.Int(int(*r.Lifecycle.DeleteAfterDays))
			}
			if r.Lifecycle.OptInToArchiveForSupportedResources != nil {
				lifecycleArgs.OptInToArchiveForSupportedResources = pulumi.Bool(*r.Lifecycle.OptInToArchiveForSupportedResources)
			}
			ruleArgs.Lifecycle = lifecycleArgs
		}
		var copyActions backup.PlanRuleCopyActionArray
		for _, c := range r.CopyActions {
			copyArgs := &backup.PlanRuleCopyActionArgs{
				DestinationVaultArn: pulumi.String(c.DestinationVaultArn.GetValue()),
			}
			if c.Lifecycle != nil {
				lifecycleArgs := &backup.PlanRuleCopyActionLifecycleArgs{}
				if c.Lifecycle.ColdStorageAfterDays != nil {
					lifecycleArgs.ColdStorageAfter = pulumi.Int(int(*c.Lifecycle.ColdStorageAfterDays))
				}
				if c.Lifecycle.DeleteAfterDays != nil {
					lifecycleArgs.DeleteAfter = pulumi.Int(int(*c.Lifecycle.DeleteAfterDays))
				}
				if c.Lifecycle.OptInToArchiveForSupportedResources != nil {
					lifecycleArgs.OptInToArchiveForSupportedResources = pulumi.Bool(*c.Lifecycle.OptInToArchiveForSupportedResources)
				}
				copyArgs.Lifecycle = lifecycleArgs
			}
			copyActions = append(copyActions, copyArgs)
		}
		if len(copyActions) > 0 {
			ruleArgs.CopyActions = copyActions
		}
		var scanActions backup.PlanRuleScanActionArray
		for _, s := range r.ScanActions {
			scanActions = append(scanActions, &backup.PlanRuleScanActionArgs{
				MalwareScanner: pulumi.String(s.MalwareScanner),
				ScanMode:       pulumi.String(s.ScanMode),
			})
		}
		if len(scanActions) > 0 {
			ruleArgs.ScanActions = scanActions
		}
		rules = append(rules, ruleArgs)
	}

	args := &backup.PlanArgs{
		// metadata.name is the plan name on both engines. Changing it
		// replaces the plan (and re-parents every selection).
		Name:  pulumi.String(locals.Target.Metadata.Name),
		Rules: rules,
		Tags:  pulumi.ToStringMap(locals.AwsTags),
	}

	var advancedSettings backup.PlanAdvancedBackupSettingArray
	for _, a := range spec.AdvancedBackupSettings {
		advancedSettings = append(advancedSettings, &backup.PlanAdvancedBackupSettingArgs{
			ResourceType:  pulumi.String(a.ResourceType),
			BackupOptions: pulumi.ToStringMap(a.BackupOptions),
		})
	}
	if len(advancedSettings) > 0 {
		args.AdvancedBackupSettings = advancedSettings
	}

	if spec.ScanSetting != nil {
		args.ScanSettings = backup.PlanScanSettingArray{
			&backup.PlanScanSettingArgs{
				MalwareScanner: pulumi.String(spec.ScanSetting.MalwareScanner),
				ResourceTypes:  pulumi.ToStringArray(spec.ScanSetting.ResourceTypes),
				ScannerRoleArn: pulumi.String(spec.ScanSetting.ScannerRoleArn.GetValue()),
			},
		}
	}

	createdPlan, err := backup.NewPlan(ctx, "plan", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create plan")
	}

	// Folded selections, keyed by name (fully ForceNew at the provider
	// - any change replaces the selection, never edits it).
	selectionIds := pulumi.StringMap{}
	for _, sel := range spec.Selections {
		createdSelection, err := selection(ctx, createdPlan, sel, provider)
		if err != nil {
			return errors.Wrapf(err, "selection %s", sel.Name)
		}
		selectionIds[sel.Name] = createdSelection.ID().ToStringOutput()
	}

	ctx.Export(OpPlanId, createdPlan.ID())
	ctx.Export(OpPlanArn, createdPlan.Arn)
	ctx.Export(OpPlanVersion, createdPlan.Version)
	ctx.Export(OpSelectionIds, selectionIds)
	return nil
}

// selection creates one folded backup selection under the plan.
func selection(ctx *pulumi.Context, createdPlan *backup.Plan,
	sel *awsbackupplanv1alpha1.AwsBackupPlanSelection, provider *aws.Provider) (*backup.Selection, error) {

	args := &backup.SelectionArgs{
		Name:       pulumi.String(sel.Name),
		PlanId:     createdPlan.ID(),
		IamRoleArn: pulumi.String(sel.IamRoleArn.GetValue()),
	}
	if len(sel.Resources) > 0 {
		args.Resources = pulumi.ToStringArray(sel.Resources)
	}
	if len(sel.NotResources) > 0 {
		args.NotResources = pulumi.ToStringArray(sel.NotResources)
	}
	var selectionTags backup.SelectionSelectionTagArray
	for _, t := range sel.SelectionTags {
		selectionTags = append(selectionTags, &backup.SelectionSelectionTagArgs{
			Type:  pulumi.String(t.Type),
			Key:   pulumi.String(t.Key),
			Value: pulumi.String(t.Value),
		})
	}
	if len(selectionTags) > 0 {
		args.SelectionTags = selectionTags
	}
	if sel.Conditions != nil {
		conditionArgs := &backup.SelectionConditionArgs{}
		var stringEquals backup.SelectionConditionStringEqualArray
		for _, p := range sel.Conditions.StringEquals {
			stringEquals = append(stringEquals, &backup.SelectionConditionStringEqualArgs{
				Key: pulumi.String(p.Key), Value: pulumi.String(p.Value),
			})
		}
		if len(stringEquals) > 0 {
			conditionArgs.StringEquals = stringEquals
		}
		var stringNotEquals backup.SelectionConditionStringNotEqualArray
		for _, p := range sel.Conditions.StringNotEquals {
			stringNotEquals = append(stringNotEquals, &backup.SelectionConditionStringNotEqualArgs{
				Key: pulumi.String(p.Key), Value: pulumi.String(p.Value),
			})
		}
		if len(stringNotEquals) > 0 {
			conditionArgs.StringNotEquals = stringNotEquals
		}
		var stringLikes backup.SelectionConditionStringLikeArray
		for _, p := range sel.Conditions.StringLike {
			stringLikes = append(stringLikes, &backup.SelectionConditionStringLikeArgs{
				Key: pulumi.String(p.Key), Value: pulumi.String(p.Value),
			})
		}
		if len(stringLikes) > 0 {
			conditionArgs.StringLikes = stringLikes
		}
		var stringNotLikes backup.SelectionConditionStringNotLikeArray
		for _, p := range sel.Conditions.StringNotLike {
			stringNotLikes = append(stringNotLikes, &backup.SelectionConditionStringNotLikeArgs{
				Key: pulumi.String(p.Key), Value: pulumi.String(p.Value),
			})
		}
		if len(stringNotLikes) > 0 {
			conditionArgs.StringNotLikes = stringNotLikes
		}
		args.Conditions = backup.SelectionConditionArray{conditionArgs}
	}

	return backup.NewSelection(ctx, "selection-"+sel.Name, args,
		pulumi.Provider(provider), pulumi.Parent(createdPlan))
}
