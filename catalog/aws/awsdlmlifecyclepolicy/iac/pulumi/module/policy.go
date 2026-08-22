package module

import (
	"github.com/pkg/errors"
	awsdlmlifecyclepolicyv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsdlmlifecyclepolicy/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dlm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// lifecyclePolicy creates the DLM policy (default XOR custom mode)
// and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the provider's DefaultPolicy top-level argument and the
//     PolicyLanguage are DERIVED here from the configured arm
//     (SIMPLIFIED for default mode, STANDARD for custom) - a spec
//     field for either could contradict the arm;
//   - default mode expresses "AWS-default cadence" by OMITTING
//     arguments (the provider diff-suppresses create/retain interval
//     defaults);
//   - Schedule.CopyTags is ForceNew (replaces the whole schedule);
//     everything else in a schedule updates in place;
//   - the bridge flattens two provider MaxItems-1 lists to scalars:
//     ResourceLocations and CreateRule.Times are single strings here.
func lifecyclePolicy(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	description := spec.Description
	if description == "" {
		description = locals.Target.Metadata.Name
	}
	state := "ENABLED"
	if spec.Disabled {
		state = "DISABLED"
	}

	args := &dlm.LifecyclePolicyArgs{
		Description:      pulumi.String(description),
		ExecutionRoleArn: pulumi.String(spec.ExecutionRoleArn.GetValue()),
		State:            pulumi.String(state),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.DefaultPolicy != nil {
		args.DefaultPolicy = pulumi.String(spec.DefaultPolicy.ResourceType)
		args.PolicyDetails = buildDefaultPolicyDetails(spec.DefaultPolicy)
	} else {
		args.PolicyDetails = buildCustomPolicyDetails(spec.CustomPolicy)
	}

	createdPolicy, err := dlm.NewLifecyclePolicy(ctx, "policy", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create lifecycle policy")
	}

	ctx.Export(OpPolicyId, createdPolicy.ID())
	ctx.Export(OpPolicyArn, createdPolicy.Arn)
	return nil
}

// buildDefaultPolicyDetails renders DEFAULT mode: the resource type
// plus the simplified dials, omitting anything AWS defaults.
func buildDefaultPolicyDetails(defaultPolicy *awsdlmlifecyclepolicyv1alpha1.AwsDlmDefaultPolicy) dlm.LifecyclePolicyPolicyDetailsArgs {
	details := dlm.LifecyclePolicyPolicyDetailsArgs{
		ResourceType: pulumi.String(defaultPolicy.ResourceType),
	}
	if defaultPolicy.CreateIntervalDays > 0 {
		details.CreateInterval = pulumi.Int(int(defaultPolicy.CreateIntervalDays))
	}
	if defaultPolicy.RetainIntervalDays > 0 {
		details.RetainInterval = pulumi.Int(int(defaultPolicy.RetainIntervalDays))
	}
	if defaultPolicy.CopyTags {
		details.CopyTags = pulumi.Bool(true)
	}
	if defaultPolicy.ExtendDeletion {
		details.ExtendDeletion = pulumi.Bool(true)
	}
	if exclusions := defaultPolicy.Exclusions; exclusions != nil {
		exclusionsArgs := &dlm.LifecyclePolicyPolicyDetailsExclusionsArgs{}
		if exclusions.ExcludeBootVolumes {
			exclusionsArgs.ExcludeBootVolumes = pulumi.Bool(true)
		}
		if len(exclusions.ExcludeTags) > 0 {
			exclusionsArgs.ExcludeTags = pulumi.ToStringMap(exclusions.ExcludeTags)
		}
		if len(exclusions.ExcludeVolumeTypes) > 0 {
			exclusionsArgs.ExcludeVolumeTypes = pulumi.ToStringArray(exclusions.ExcludeVolumeTypes)
		}
		details.Exclusions = exclusionsArgs
	}
	return details
}

// buildCustomPolicyDetails renders CUSTOM mode: tag-targeted
// schedules or an event-based policy.
func buildCustomPolicyDetails(customPolicy *awsdlmlifecyclepolicyv1alpha1.AwsDlmCustomPolicy) dlm.LifecyclePolicyPolicyDetailsArgs {
	policyType := customPolicy.PolicyType
	if policyType == "" {
		policyType = "EBS_SNAPSHOT_MANAGEMENT"
	}

	details := dlm.LifecyclePolicyPolicyDetailsArgs{
		PolicyType: pulumi.String(policyType),
	}
	if len(customPolicy.ResourceTypes) > 0 {
		details.ResourceTypes = pulumi.ToStringArray(customPolicy.ResourceTypes)
	}
	// The provider's MaxItems-1 list arrives flattened as one string
	// at the bridge.
	if len(customPolicy.ResourceLocations) > 0 {
		details.ResourceLocations = pulumi.String(customPolicy.ResourceLocations[0])
	}
	if len(customPolicy.TargetTags) > 0 {
		details.TargetTags = pulumi.ToStringMap(customPolicy.TargetTags)
	}

	if parameters := customPolicy.Parameters; parameters != nil {
		parametersArgs := &dlm.LifecyclePolicyPolicyDetailsParametersArgs{}
		if parameters.ExcludeBootVolume {
			parametersArgs.ExcludeBootVolume = pulumi.Bool(true)
		}
		if parameters.NoReboot {
			parametersArgs.NoReboot = pulumi.Bool(true)
		}
		details.Parameters = parametersArgs
	}

	if eventSource := customPolicy.EventSource; eventSource != nil {
		details.EventSource = &dlm.LifecyclePolicyPolicyDetailsEventSourceArgs{
			Type: pulumi.String("MANAGED_CWE"),
			Parameters: &dlm.LifecyclePolicyPolicyDetailsEventSourceParametersArgs{
				DescriptionRegex: pulumi.String(eventSource.DescriptionRegex),
				EventType:        pulumi.String(eventSource.EventType),
				SnapshotOwners:   pulumi.ToStringArray(eventSource.SnapshotOwners),
			},
		}
	}

	if action := customPolicy.Action; action != nil {
		copies := dlm.LifecyclePolicyPolicyDetailsActionCrossRegionCopyArray{}
		for _, copy := range action.CrossRegionCopies {
			encryptionConfiguration := &dlm.LifecyclePolicyPolicyDetailsActionCrossRegionCopyEncryptionConfigurationArgs{
				Encrypted: pulumi.Bool(copy.Encrypted),
			}
			if copy.CmkArn.GetValue() != "" {
				encryptionConfiguration.CmkArn = pulumi.String(copy.CmkArn.GetValue())
			}
			copyArgs := &dlm.LifecyclePolicyPolicyDetailsActionCrossRegionCopyArgs{
				Target:                  pulumi.String(copy.Target),
				EncryptionConfiguration: encryptionConfiguration,
			}
			if copy.RetainRule != nil {
				copyArgs.RetainRule = &dlm.LifecyclePolicyPolicyDetailsActionCrossRegionCopyRetainRuleArgs{
					Interval:     pulumi.Int(int(copy.RetainRule.Interval)),
					IntervalUnit: pulumi.String(copy.RetainRule.IntervalUnit),
				}
			}
			copies = append(copies, copyArgs)
		}
		details.Action = &dlm.LifecyclePolicyPolicyDetailsActionArgs{
			Name:              pulumi.String(action.Name),
			CrossRegionCopies: copies,
		}
	}

	schedules := dlm.LifecyclePolicyPolicyDetailsScheduleArray{}
	for _, schedule := range customPolicy.Schedules {
		schedules = append(schedules, buildSchedule(schedule))
	}
	if len(schedules) > 0 {
		details.Schedules = schedules
	}

	return details
}

// buildSchedule renders one named schedule with all its rules.
func buildSchedule(schedule *awsdlmlifecyclepolicyv1alpha1.AwsDlmSchedule) *dlm.LifecyclePolicyPolicyDetailsScheduleArgs {
	createRule := &dlm.LifecyclePolicyPolicyDetailsScheduleCreateRuleArgs{}
	if schedule.CreateRule.IntervalHours > 0 {
		createRule.Interval = pulumi.Int(int(schedule.CreateRule.IntervalHours))
		createRule.IntervalUnit = pulumi.String("HOURS")
	}
	// The provider's MaxItems-1 list arrives flattened as one string
	// at the bridge.
	if len(schedule.CreateRule.Times) > 0 {
		createRule.Times = pulumi.String(schedule.CreateRule.Times[0])
	}
	if schedule.CreateRule.CronExpression != "" {
		createRule.CronExpression = pulumi.String(schedule.CreateRule.CronExpression)
	}
	if schedule.CreateRule.Location != "" {
		createRule.Location = pulumi.String(schedule.CreateRule.Location)
	}
	if scripts := schedule.CreateRule.Scripts; scripts != nil {
		scriptsArgs := &dlm.LifecyclePolicyPolicyDetailsScheduleCreateRuleScriptsArgs{
			ExecutionHandler: pulumi.String(scripts.ExecutionHandler),
		}
		if len(scripts.Stages) > 0 {
			scriptsArgs.Stages = pulumi.ToStringArray(scripts.Stages)
		}
		if scripts.ExecutionHandlerService != "" {
			scriptsArgs.ExecutionHandlerService = pulumi.String(scripts.ExecutionHandlerService)
		}
		if scripts.ExecuteOperationOnScriptFailure {
			scriptsArgs.ExecuteOperationOnScriptFailure = pulumi.Bool(true)
		}
		if scripts.ExecutionTimeoutSeconds > 0 {
			scriptsArgs.ExecutionTimeout = pulumi.Int(int(scripts.ExecutionTimeoutSeconds))
		}
		if scripts.MaximumRetryCount > 0 {
			scriptsArgs.MaximumRetryCount = pulumi.Int(int(scripts.MaximumRetryCount))
		}
		createRule.Scripts = scriptsArgs
	}

	retainRule := &dlm.LifecyclePolicyPolicyDetailsScheduleRetainRuleArgs{}
	if schedule.RetainRule.Count > 0 {
		retainRule.Count = pulumi.Int(int(schedule.RetainRule.Count))
	}
	if schedule.RetainRule.Interval > 0 {
		retainRule.Interval = pulumi.Int(int(schedule.RetainRule.Interval))
		retainRule.IntervalUnit = pulumi.String(schedule.RetainRule.IntervalUnit)
	}

	scheduleArgs := &dlm.LifecyclePolicyPolicyDetailsScheduleArgs{
		Name:       pulumi.String(schedule.Name),
		CreateRule: createRule,
		RetainRule: retainRule,
	}
	if schedule.CopyTags {
		scheduleArgs.CopyTags = pulumi.Bool(true)
	}
	if len(schedule.TagsToAdd) > 0 {
		scheduleArgs.TagsToAdd = pulumi.ToStringMap(schedule.TagsToAdd)
	}
	if len(schedule.VariableTags) > 0 {
		scheduleArgs.VariableTags = pulumi.ToStringMap(schedule.VariableTags)
	}

	if archive := schedule.ArchiveRule; archive != nil {
		tier := &dlm.LifecyclePolicyPolicyDetailsScheduleArchiveRuleArchiveRetainRuleRetentionArchiveTierArgs{}
		if archive.Count > 0 {
			tier.Count = pulumi.Int(int(archive.Count))
		}
		if archive.Interval > 0 {
			tier.Interval = pulumi.Int(int(archive.Interval))
			tier.IntervalUnit = pulumi.String(archive.IntervalUnit)
		}
		scheduleArgs.ArchiveRule = &dlm.LifecyclePolicyPolicyDetailsScheduleArchiveRuleArgs{
			ArchiveRetainRule: &dlm.LifecyclePolicyPolicyDetailsScheduleArchiveRuleArchiveRetainRuleArgs{
				RetentionArchiveTier: tier,
			},
		}
	}

	crossRegionRules := dlm.LifecyclePolicyPolicyDetailsScheduleCrossRegionCopyRuleArray{}
	for _, rule := range schedule.CrossRegionCopyRules {
		ruleArgs := &dlm.LifecyclePolicyPolicyDetailsScheduleCrossRegionCopyRuleArgs{
			TargetRegion: pulumi.String(rule.TargetRegion),
			Encrypted:    pulumi.Bool(rule.Encrypted),
		}
		if rule.CmkArn.GetValue() != "" {
			ruleArgs.CmkArn = pulumi.String(rule.CmkArn.GetValue())
		}
		if rule.CopyTags {
			ruleArgs.CopyTags = pulumi.Bool(true)
		}
		if rule.RetainRule != nil {
			ruleArgs.RetainRule = &dlm.LifecyclePolicyPolicyDetailsScheduleCrossRegionCopyRuleRetainRuleArgs{
				Interval:     pulumi.Int(int(rule.RetainRule.Interval)),
				IntervalUnit: pulumi.String(rule.RetainRule.IntervalUnit),
			}
		}
		if rule.DeprecateRule != nil {
			ruleArgs.DeprecateRule = &dlm.LifecyclePolicyPolicyDetailsScheduleCrossRegionCopyRuleDeprecateRuleArgs{
				Interval:     pulumi.Int(int(rule.DeprecateRule.Interval)),
				IntervalUnit: pulumi.String(rule.DeprecateRule.IntervalUnit),
			}
		}
		crossRegionRules = append(crossRegionRules, ruleArgs)
	}
	if len(crossRegionRules) > 0 {
		scheduleArgs.CrossRegionCopyRules = crossRegionRules
	}

	if deprecate := schedule.DeprecateRule; deprecate != nil {
		deprecateArgs := &dlm.LifecyclePolicyPolicyDetailsScheduleDeprecateRuleArgs{}
		if deprecate.Count > 0 {
			deprecateArgs.Count = pulumi.Int(int(deprecate.Count))
		}
		if deprecate.Interval > 0 {
			deprecateArgs.Interval = pulumi.Int(int(deprecate.Interval))
			deprecateArgs.IntervalUnit = pulumi.String(deprecate.IntervalUnit)
		}
		scheduleArgs.DeprecateRule = deprecateArgs
	}

	if fastRestore := schedule.FastRestoreRule; fastRestore != nil {
		fastRestoreArgs := &dlm.LifecyclePolicyPolicyDetailsScheduleFastRestoreRuleArgs{
			AvailabilityZones: pulumi.ToStringArray(fastRestore.AvailabilityZones),
		}
		if fastRestore.Count > 0 {
			fastRestoreArgs.Count = pulumi.Int(int(fastRestore.Count))
		}
		if fastRestore.Interval > 0 {
			fastRestoreArgs.Interval = pulumi.Int(int(fastRestore.Interval))
			fastRestoreArgs.IntervalUnit = pulumi.String(fastRestore.IntervalUnit)
		}
		scheduleArgs.FastRestoreRule = fastRestoreArgs
	}

	if share := schedule.ShareRule; share != nil {
		shareArgs := &dlm.LifecyclePolicyPolicyDetailsScheduleShareRuleArgs{
			TargetAccounts: pulumi.ToStringArray(share.TargetAccounts),
		}
		if share.UnshareInterval > 0 {
			shareArgs.UnshareInterval = pulumi.Int(int(share.UnshareInterval))
			shareArgs.UnshareIntervalUnit = pulumi.String(share.UnshareIntervalUnit)
		}
		scheduleArgs.ShareRule = shareArgs
	}

	return scheduleArgs
}
