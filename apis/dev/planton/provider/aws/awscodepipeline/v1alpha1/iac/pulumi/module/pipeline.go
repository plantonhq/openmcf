package module

import (
	"github.com/pkg/errors"
	awscodepipelinev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscodepipeline/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/codepipeline"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// pipeline provisions the CodePipeline pipeline -- a metadata-only
// control-plane resource: create/update/delete are single API calls that
// complete in seconds. The only operational wait is IAM eventual consistency
// on a freshly created pipeline role, which the provider absorbs with a
// bounded retry on create.
func pipeline(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*codepipeline.Pipeline, error) {
	spec := locals.AwsCodePipeline.Spec

	// --- Artifact stores ---
	// AWS models the store two ways: exactly one store without a region
	// (single-region pipeline) or one store per region (cross-region). The
	// spec's CEL enforces the shape; the mapping is uniform here because the
	// SDK carries region as an optional field.
	var artifactStores codepipeline.PipelineArtifactStoreArray
	for _, store := range spec.ArtifactStores {
		storeArgs := &codepipeline.PipelineArtifactStoreArgs{
			Location: pulumi.String(store.Location.GetValue()),
			Type:     pulumi.String("S3"),
		}
		if store.Region != "" {
			storeArgs.Region = pulumi.StringPtr(store.Region)
		}
		if store.EncryptionKeyId != nil && store.EncryptionKeyId.GetValue() != "" {
			storeArgs.EncryptionKey = &codepipeline.PipelineArtifactStoreEncryptionKeyArgs{
				Id:   pulumi.String(store.EncryptionKeyId.GetValue()),
				Type: pulumi.String("KMS"),
			}
		}
		artifactStores = append(artifactStores, storeArgs)
	}

	// --- Stages ---
	var stages codepipeline.PipelineStageArray
	for _, stage := range spec.Stages {
		var actions codepipeline.PipelineStageActionArray
		for _, action := range stage.Actions {
			actionArgs := &codepipeline.PipelineStageActionArgs{
				Name:     pulumi.String(action.Name),
				Category: pulumi.String(action.Category),
				Owner:    pulumi.String(action.Owner),
				Provider: pulumi.String(action.Provider),
				Version:  pulumi.String(action.Version),
			}
			if len(action.Configuration) > 0 {
				actionArgs.Configuration = pulumi.ToStringMap(action.Configuration)
			}
			if len(action.InputArtifacts) > 0 {
				actionArgs.InputArtifacts = pulumi.ToStringArray(action.InputArtifacts)
			}
			if len(action.OutputArtifacts) > 0 {
				actionArgs.OutputArtifacts = pulumi.ToStringArray(action.OutputArtifacts)
			}
			if action.Namespace != "" {
				actionArgs.Namespace = pulumi.StringPtr(action.Namespace)
			}
			if action.Region != "" {
				actionArgs.Region = pulumi.StringPtr(action.Region)
			}
			if action.RoleArn != nil && action.RoleArn.GetValue() != "" {
				actionArgs.RoleArn = pulumi.StringPtr(action.RoleArn.GetValue())
			}
			if action.RunOrder > 0 {
				actionArgs.RunOrder = pulumi.IntPtr(int(action.RunOrder))
			}
			// Per-action timeout override (omit for the provider default).
			if action.TimeoutInMinutes > 0 {
				actionArgs.TimeoutInMinutes = pulumi.IntPtr(int(action.TimeoutInMinutes))
			}
			actions = append(actions, actionArgs)
		}

		stageArgs := &codepipeline.PipelineStageArgs{
			Name:    pulumi.String(stage.Name),
			Actions: actions,
		}

		// Entry gate: rules that must pass before the stage starts (e.g., a
		// DeploymentWindow rule admitting executions only in business hours).
		if stage.BeforeEntry != nil {
			stageArgs.BeforeEntry = &codepipeline.PipelineStageBeforeEntryArgs{
				Condition: buildBeforeEntryCondition(stage.BeforeEntry),
			}
		}

		// Post-success verification: a failing rule fails the stage despite
		// successful actions (e.g., a post-deploy CloudWatchAlarm check).
		if stage.OnSuccess != nil {
			stageArgs.OnSuccess = &codepipeline.PipelineStageOnSuccessArgs{
				Condition: buildOnSuccessCondition(stage.OnSuccess),
			}
		}

		// Failure handling: automatic rollback to the last successful state,
		// automatic retry, or rule-gated handling.
		if stage.OnFailure != nil {
			failureArgs := &codepipeline.PipelineStageOnFailureArgs{}
			if stage.OnFailure.Result != "" {
				failureArgs.Result = pulumi.StringPtr(stage.OnFailure.Result)
			}
			if stage.OnFailure.RetryConfiguration != nil {
				failureArgs.RetryConfiguration = &codepipeline.PipelineStageOnFailureRetryConfigurationArgs{
					RetryMode: pulumi.StringPtr(stage.OnFailure.RetryConfiguration.RetryMode),
				}
			}
			if stage.OnFailure.Condition != nil {
				failureArgs.Condition = buildOnFailureCondition(stage.OnFailure.Condition)
			}
			stageArgs.OnFailure = failureArgs
		}

		stages = append(stages, stageArgs)
	}

	// --- Pipeline args ---
	// The spec defaults to V2/SUPERSEDED, while the PROVIDER defaults
	// pipeline_type to V1 -- an omitted value must deploy the same pipeline
	// an explicit V2 would, so the spec defaults are applied here, never
	// left to the provider.
	pipelineType := spec.GetPipelineType()
	if pipelineType == "" {
		pipelineType = "V2"
	}
	executionMode := spec.GetExecutionMode()
	if executionMode == "" {
		executionMode = "SUPERSEDED"
	}

	args := &codepipeline.PipelineArgs{
		Name:           pulumi.StringPtr(locals.PipelineName),
		RoleArn:        pulumi.String(spec.RoleArn.GetValue()),
		ArtifactStores: artifactStores,
		Stages:         stages,
		PipelineType:   pulumi.StringPtr(pipelineType),
		ExecutionMode:  pulumi.StringPtr(executionMode),
		Tags:           pulumi.ToStringMap(locals.Labels),
	}

	// --- Triggers (V2 git-event execution) ---
	if len(spec.Triggers) > 0 {
		var triggers codepipeline.PipelineTriggerArray
		for _, trigger := range spec.Triggers {
			triggerArgs := &codepipeline.PipelineTriggerArgs{
				ProviderType: pulumi.String(trigger.ProviderType),
			}
			if trigger.GitConfiguration != nil {
				gitArgs := &codepipeline.PipelineTriggerGitConfigurationArgs{
					SourceActionName: pulumi.String(trigger.GitConfiguration.SourceActionName),
				}
				if len(trigger.GitConfiguration.Push) > 0 {
					gitArgs.Pushes = buildPushFilters(trigger.GitConfiguration.Push)
				}
				if len(trigger.GitConfiguration.PullRequest) > 0 {
					gitArgs.PullRequests = buildPullRequestFilters(trigger.GitConfiguration.PullRequest)
				}
				triggerArgs.GitConfiguration = gitArgs
			}
			triggers = append(triggers, triggerArgs)
		}
		args.Triggers = triggers
	}

	// --- Pipeline-level variables (V2) ---
	if len(spec.Variables) > 0 {
		var variables codepipeline.PipelineVariableArray
		for _, v := range spec.Variables {
			varArgs := &codepipeline.PipelineVariableArgs{
				Name: pulumi.String(v.Name),
			}
			if v.DefaultValue != "" {
				varArgs.DefaultValue = pulumi.StringPtr(v.DefaultValue)
			}
			if v.Description != "" {
				varArgs.Description = pulumi.StringPtr(v.Description)
			}
			variables = append(variables, varArgs)
		}
		args.Variables = variables
	}

	created, err := codepipeline.NewPipeline(ctx, "codepipeline", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create codepipeline")
	}

	return created, nil
}

// The SDK generates a distinct Go type per condition path (before_entry /
// on_success / on_failure), so each has its own builder even though the spec
// shape is one message. ruleTypeId defaults (category "Rule", owner "AWS")
// are applied here because the SDK marks category required.

func buildBeforeEntryCondition(c *awscodepipelinev1alpha1.AwsCodePipelineStageCondition) *codepipeline.PipelineStageBeforeEntryConditionArgs {
	condArgs := &codepipeline.PipelineStageBeforeEntryConditionArgs{}
	if c.Result != "" {
		condArgs.Result = pulumi.StringPtr(c.Result)
	}
	var rules codepipeline.PipelineStageBeforeEntryConditionRuleArray
	for _, r := range c.Rules {
		ruleArgs := &codepipeline.PipelineStageBeforeEntryConditionRuleArgs{
			Name:       pulumi.String(r.Name),
			RuleTypeId: buildBeforeEntryRuleTypeId(r.RuleTypeId),
		}
		applyCommonRuleFields(r,
			func(m pulumi.StringMapInput) { ruleArgs.Configuration = m },
			func(a pulumi.StringArrayInput) { ruleArgs.Commands = a },
			func(a pulumi.StringArrayInput) { ruleArgs.InputArtifacts = a },
			func(s pulumi.StringPtrInput) { ruleArgs.Region = s },
			func(s pulumi.StringPtrInput) { ruleArgs.RoleArn = s },
			func(i pulumi.IntPtrInput) { ruleArgs.TimeoutInMinutes = i },
		)
		rules = append(rules, ruleArgs)
	}
	condArgs.Rules = rules
	return condArgs
}

func buildOnSuccessCondition(c *awscodepipelinev1alpha1.AwsCodePipelineStageCondition) *codepipeline.PipelineStageOnSuccessConditionArgs {
	condArgs := &codepipeline.PipelineStageOnSuccessConditionArgs{}
	if c.Result != "" {
		condArgs.Result = pulumi.StringPtr(c.Result)
	}
	var rules codepipeline.PipelineStageOnSuccessConditionRuleArray
	for _, r := range c.Rules {
		ruleArgs := &codepipeline.PipelineStageOnSuccessConditionRuleArgs{
			Name:       pulumi.String(r.Name),
			RuleTypeId: buildOnSuccessRuleTypeId(r.RuleTypeId),
		}
		applyCommonRuleFields(r,
			func(m pulumi.StringMapInput) { ruleArgs.Configuration = m },
			func(a pulumi.StringArrayInput) { ruleArgs.Commands = a },
			func(a pulumi.StringArrayInput) { ruleArgs.InputArtifacts = a },
			func(s pulumi.StringPtrInput) { ruleArgs.Region = s },
			func(s pulumi.StringPtrInput) { ruleArgs.RoleArn = s },
			func(i pulumi.IntPtrInput) { ruleArgs.TimeoutInMinutes = i },
		)
		rules = append(rules, ruleArgs)
	}
	condArgs.Rules = rules
	return condArgs
}

func buildOnFailureCondition(c *awscodepipelinev1alpha1.AwsCodePipelineStageCondition) *codepipeline.PipelineStageOnFailureConditionArgs {
	condArgs := &codepipeline.PipelineStageOnFailureConditionArgs{}
	if c.Result != "" {
		condArgs.Result = pulumi.StringPtr(c.Result)
	}
	var rules codepipeline.PipelineStageOnFailureConditionRuleArray
	for _, r := range c.Rules {
		ruleArgs := &codepipeline.PipelineStageOnFailureConditionRuleArgs{
			Name:       pulumi.String(r.Name),
			RuleTypeId: buildOnFailureRuleTypeId(r.RuleTypeId),
		}
		applyCommonRuleFields(r,
			func(m pulumi.StringMapInput) { ruleArgs.Configuration = m },
			func(a pulumi.StringArrayInput) { ruleArgs.Commands = a },
			func(a pulumi.StringArrayInput) { ruleArgs.InputArtifacts = a },
			func(s pulumi.StringPtrInput) { ruleArgs.Region = s },
			func(s pulumi.StringPtrInput) { ruleArgs.RoleArn = s },
			func(i pulumi.IntPtrInput) { ruleArgs.TimeoutInMinutes = i },
		)
		rules = append(rules, ruleArgs)
	}
	condArgs.Rules = rules
	return condArgs
}

// applyCommonRuleFields maps the optional rule fields shared by all three
// condition paths through setter callbacks, keeping the per-path builders
// down to their type-specific skeletons.
func applyCommonRuleFields(
	r *awscodepipelinev1alpha1.AwsCodePipelineRule,
	setConfiguration func(pulumi.StringMapInput),
	setCommands func(pulumi.StringArrayInput),
	setInputArtifacts func(pulumi.StringArrayInput),
	setRegion func(pulumi.StringPtrInput),
	setRoleArn func(pulumi.StringPtrInput),
	setTimeout func(pulumi.IntPtrInput),
) {
	if len(r.Configuration) > 0 {
		setConfiguration(pulumi.ToStringMap(r.Configuration))
	}
	if len(r.Commands) > 0 {
		setCommands(pulumi.ToStringArray(r.Commands))
	}
	if len(r.InputArtifacts) > 0 {
		setInputArtifacts(pulumi.ToStringArray(r.InputArtifacts))
	}
	if r.Region != "" {
		setRegion(pulumi.StringPtr(r.Region))
	}
	if r.RoleArn != nil && r.RoleArn.GetValue() != "" {
		setRoleArn(pulumi.StringPtr(r.RoleArn.GetValue()))
	}
	if r.TimeoutInMinutes > 0 {
		setTimeout(pulumi.IntPtr(int(r.TimeoutInMinutes)))
	}
}

func buildBeforeEntryRuleTypeId(id *awscodepipelinev1alpha1.AwsCodePipelineRuleTypeId) *codepipeline.PipelineStageBeforeEntryConditionRuleRuleTypeIdArgs {
	args := &codepipeline.PipelineStageBeforeEntryConditionRuleRuleTypeIdArgs{
		Category: pulumi.String(ruleCategory(id)),
		Provider: pulumi.String(id.Provider),
	}
	if owner := ruleOwner(id); owner != "" {
		args.Owner = pulumi.StringPtr(owner)
	}
	if id.Version != "" {
		args.Version = pulumi.StringPtr(id.Version)
	}
	return args
}

func buildOnSuccessRuleTypeId(id *awscodepipelinev1alpha1.AwsCodePipelineRuleTypeId) *codepipeline.PipelineStageOnSuccessConditionRuleRuleTypeIdArgs {
	args := &codepipeline.PipelineStageOnSuccessConditionRuleRuleTypeIdArgs{
		Category: pulumi.String(ruleCategory(id)),
		Provider: pulumi.String(id.Provider),
	}
	if owner := ruleOwner(id); owner != "" {
		args.Owner = pulumi.StringPtr(owner)
	}
	if id.Version != "" {
		args.Version = pulumi.StringPtr(id.Version)
	}
	return args
}

func buildOnFailureRuleTypeId(id *awscodepipelinev1alpha1.AwsCodePipelineRuleTypeId) *codepipeline.PipelineStageOnFailureConditionRuleRuleTypeIdArgs {
	args := &codepipeline.PipelineStageOnFailureConditionRuleRuleTypeIdArgs{
		Category: pulumi.String(ruleCategory(id)),
		Provider: pulumi.String(id.Provider),
	}
	if owner := ruleOwner(id); owner != "" {
		args.Owner = pulumi.StringPtr(owner)
	}
	if id.Version != "" {
		args.Version = pulumi.StringPtr(id.Version)
	}
	return args
}

// ruleCategory/ruleOwner apply the spec defaults ("Rule"/"AWS") when the
// presence-carrying optionals are unset -- AWS accepts no other values today.
func ruleCategory(id *awscodepipelinev1alpha1.AwsCodePipelineRuleTypeId) string {
	if id.GetCategory() != "" {
		return id.GetCategory()
	}
	return "Rule"
}

func ruleOwner(id *awscodepipelinev1alpha1.AwsCodePipelineRuleTypeId) string {
	if id.GetOwner() != "" {
		return id.GetOwner()
	}
	return "AWS"
}

func buildPushFilters(pushes []*awscodepipelinev1alpha1.AwsCodePipelineGitPush) codepipeline.PipelineTriggerGitConfigurationPushArray {
	var result codepipeline.PipelineTriggerGitConfigurationPushArray
	for _, push := range pushes {
		pushArgs := &codepipeline.PipelineTriggerGitConfigurationPushArgs{}
		if push.Branches != nil {
			pushArgs.Branches = &codepipeline.PipelineTriggerGitConfigurationPushBranchesArgs{
				Includes: pulumi.ToStringArray(push.Branches.Includes),
				Excludes: pulumi.ToStringArray(push.Branches.Excludes),
			}
		}
		if push.FilePaths != nil {
			pushArgs.FilePaths = &codepipeline.PipelineTriggerGitConfigurationPushFilePathsArgs{
				Includes: pulumi.ToStringArray(push.FilePaths.Includes),
				Excludes: pulumi.ToStringArray(push.FilePaths.Excludes),
			}
		}
		if push.Tags != nil {
			pushArgs.Tags = &codepipeline.PipelineTriggerGitConfigurationPushTagsArgs{
				Includes: pulumi.ToStringArray(push.Tags.Includes),
				Excludes: pulumi.ToStringArray(push.Tags.Excludes),
			}
		}
		result = append(result, pushArgs)
	}
	return result
}

func buildPullRequestFilters(prs []*awscodepipelinev1alpha1.AwsCodePipelineGitPullRequest) codepipeline.PipelineTriggerGitConfigurationPullRequestArray {
	var result codepipeline.PipelineTriggerGitConfigurationPullRequestArray
	for _, pr := range prs {
		prArgs := &codepipeline.PipelineTriggerGitConfigurationPullRequestArgs{}
		if pr.Branches != nil {
			prArgs.Branches = &codepipeline.PipelineTriggerGitConfigurationPullRequestBranchesArgs{
				Includes: pulumi.ToStringArray(pr.Branches.Includes),
				Excludes: pulumi.ToStringArray(pr.Branches.Excludes),
			}
		}
		if pr.FilePaths != nil {
			prArgs.FilePaths = &codepipeline.PipelineTriggerGitConfigurationPullRequestFilePathsArgs{
				Includes: pulumi.ToStringArray(pr.FilePaths.Includes),
				Excludes: pulumi.ToStringArray(pr.FilePaths.Excludes),
			}
		}
		if len(pr.Events) > 0 {
			prArgs.Events = pulumi.ToStringArray(pr.Events)
		}
		result = append(result, prArgs)
	}
	return result
}
