package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// onlineConfigs creates the bundle's continuous-evaluation configs and
// exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - the spec's single `enabled` field fans out to the provider's TWO
//     lifecycle fields: enable_on_create (create-time intent the API
//     never echoes back) and execution_status (the post-create
//     declarative state) -- one declarative knob, both engines wire it
//     identically;
//   - evaluator entries naming in-bundle evaluators resolve to the
//     created evaluator's AWS-generated ID, which also orders the
//     config's create behind the evaluator's;
//   - creation waits CREATING -> ACTIVE, retrying briefly upstream
//     while the execution role becomes assumable.
func onlineConfigs(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdEvaluators map[string]*createdEvaluator) error {
	spec := locals.Spec

	configIds := pulumi.StringMap{}
	configArns := pulumi.StringMap{}
	outputLogGroups := pulumi.StringMap{}

	// Iteration is name-sorted for deterministic previews.
	for _, c := range sortedOnlineConfigs(spec.OnlineEvaluationConfigs) {
		enabled := true
		if c.Enabled != nil {
			enabled = *c.Enabled
		}

		dataSource := &bedrock.AgentcoreOnlineEvaluationConfigDataSourceConfigArgs{
			CloudwatchLogs: &bedrock.AgentcoreOnlineEvaluationConfigDataSourceConfigCloudwatchLogsArgs{
				LogGroupNames: svrsToStringArray(c.DataSource.LogGroupNames),
				ServiceNames:  pulumi.ToStringArray(c.DataSource.ServiceNames),
			},
		}

		var evaluatorEntries bedrock.AgentcoreOnlineEvaluationConfigEvaluatorArray
		for _, e := range c.EvaluatorIds {
			// In-bundle evaluator names resolve to the created resource's
			// ID output (and create the dependency edge); builtins and
			// full custom IDs pass through as literals.
			if created, ok := createdEvaluators[e]; ok {
				evaluatorEntries = append(evaluatorEntries, &bedrock.AgentcoreOnlineEvaluationConfigEvaluatorArgs{
					EvaluatorId: created.resource.EvaluatorId,
				})
				continue
			}
			evaluatorEntries = append(evaluatorEntries, &bedrock.AgentcoreOnlineEvaluationConfigEvaluatorArgs{
				EvaluatorId: pulumi.String(e),
			})
		}

		rule := &bedrock.AgentcoreOnlineEvaluationConfigRuleArgs{
			SamplingConfig: &bedrock.AgentcoreOnlineEvaluationConfigRuleSamplingConfigArgs{
				SamplingPercentage: pulumi.Float64(c.Rule.SamplingPercentage),
			},
		}
		if len(c.Rule.Filters) > 0 {
			var filters bedrock.AgentcoreOnlineEvaluationConfigRuleFilterArray
			for _, f := range c.Rule.Filters {
				value := &bedrock.AgentcoreOnlineEvaluationConfigRuleFilterValueArgs{}
				// Exactly one typed value (spec-validated).
				if f.StringValue != "" {
					value.StringValue = pulumi.String(f.StringValue)
				}
				if f.BooleanValue != nil {
					value.BooleanValue = pulumi.Bool(*f.BooleanValue)
				}
				if f.DoubleValue != nil {
					value.DoubleValue = pulumi.Float64(*f.DoubleValue)
				}
				filters = append(filters, &bedrock.AgentcoreOnlineEvaluationConfigRuleFilterArgs{
					Key:      pulumi.String(f.Key),
					Operator: pulumi.String(f.Operator),
					Value:    value,
				})
			}
			rule.Filters = filters
		}
		if c.Rule.SessionTimeoutMinutes > 0 {
			rule.SessionConfig = &bedrock.AgentcoreOnlineEvaluationConfigRuleSessionConfigArgs{
				SessionTimeoutMinutes: pulumi.Int(int(c.Rule.SessionTimeoutMinutes)),
			}
		}

		args := &bedrock.AgentcoreOnlineEvaluationConfigArgs{
			OnlineEvaluationConfigName: pulumi.String(c.Name),
			EvaluationExecutionRoleArn: pulumi.String(c.ExecutionRoleArn.GetValue()),
			EnableOnCreate:             pulumi.Bool(enabled),
			ExecutionStatus:            pulumi.String(map[bool]string{true: "ENABLED", false: "DISABLED"}[enabled]),
			DataSourceConfig:           dataSource,
			Evaluators:                 evaluatorEntries,
			Rule:                       rule,
			Tags:                       pulumi.ToStringMap(locals.AwsTags),
		}
		if c.Description != "" {
			args.Description = pulumi.String(c.Description)
		}

		resource, err := bedrock.NewAgentcoreOnlineEvaluationConfig(ctx, "online-config-"+c.Name, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create online evaluation config %q", c.Name)
		}
		configIds[c.Name] = resource.OnlineEvaluationConfigId
		configArns[c.Name] = resource.OnlineEvaluationConfigArn
		// The server-assigned results log group (first output config
		// entry when present).
		outputLogGroups[c.Name] = resource.OutputConfigs.ApplyT(func(configs []bedrock.AgentcoreOnlineEvaluationConfigOutputConfig) string {
			if len(configs) > 0 && len(configs[0].CloudwatchConfigs) > 0 {
				return configs[0].CloudwatchConfigs[0].LogGroupName
			}
			return ""
		}).(pulumi.StringOutput)
	}

	ctx.Export(OpOnlineEvaluationConfigIds, configIds)
	ctx.Export(OpOnlineEvaluationConfigArns, configArns)
	ctx.Export(OpOnlineEvaluationOutputLogGroups, outputLogGroups)
	return nil
}
