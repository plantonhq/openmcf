package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/amp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// workspace creates the AMP workspace and its folded satellites, and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the workspace ALIAS can never be unset once set - AWS offers no
//     un-alias, so clearing spec.alias replaces the workspace (the
//     provider's ForceNewIfChange contract);
//   - kms_key_arn replaces the workspace on change;
//   - the workspace CONFIGURATION is created via update and has NO
//     delete API - removing the block is a no-op at AWS and the
//     last-applied retention/limits persist (the settings-retention
//     class);
//   - the alert manager definition is strictly one per workspace (its
//     provider ID is the workspace ID);
//   - the resource policy is revision-guarded server-side; the
//     provider's optional revision_id input is a state-managed
//     concurrency token, deliberately not modeled;
//   - every satellite waits on the workspace through its WorkspaceId
//     reference - no explicit DependsOn needed.
func workspace(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	workspaceArgs := &amp.WorkspaceArgs{
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Alias != "" {
		workspaceArgs.Alias = pulumi.String(spec.Alias)
	}
	if spec.KmsKeyArn.GetValue() != "" {
		workspaceArgs.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}
	if spec.Logging != nil {
		workspaceArgs.LoggingConfiguration = &amp.WorkspaceLoggingConfigurationArgs{
			LogGroupArn: pulumi.String(wildcardLogGroupArn(spec.Logging.LogGroupArn.GetValue())),
		}
	}

	createdWorkspace, err := amp.NewWorkspace(ctx, "workspace", workspaceArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create workspace")
	}

	if configuration := spec.Configuration; configuration != nil {
		configurationArgs := &amp.WorkspaceConfigurationArgs{
			WorkspaceId: createdWorkspace.ID(),
		}
		if configuration.RetentionPeriodInDays != nil {
			configurationArgs.RetentionPeriodInDays = pulumi.Int(int(*configuration.RetentionPeriodInDays))
		}
		if configuration.OutOfOrderTimeWindowInSeconds != nil {
			configurationArgs.OutOfOrderTimeWindowInSeconds = pulumi.Int(int(*configuration.OutOfOrderTimeWindowInSeconds))
		}
		if configuration.RuleQueryOffsetInSeconds != nil {
			configurationArgs.RuleQueryOffsetInSeconds = pulumi.Int(int(*configuration.RuleQueryOffsetInSeconds))
		}
		if len(configuration.LimitsPerLabelSet) > 0 {
			limits := amp.WorkspaceConfigurationLimitsPerLabelSetArray{}
			for _, limit := range configuration.LimitsPerLabelSet {
				limits = append(limits, &amp.WorkspaceConfigurationLimitsPerLabelSetArgs{
					LabelSet: pulumi.ToStringMap(limit.LabelSet),
					Limits: &amp.WorkspaceConfigurationLimitsPerLabelSetLimitsArgs{
						MaxSeries: pulumi.Int(int(limit.MaxSeries)),
					},
				})
			}
			configurationArgs.LimitsPerLabelSets = limits
		}
		if _, err := amp.NewWorkspaceConfiguration(ctx, "workspace_configuration", configurationArgs, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "create workspace configuration")
		}
	}

	if spec.AlertManagerDefinition != "" {
		if _, err := amp.NewAlertManagerDefinition(ctx, "alert_manager_definition", &amp.AlertManagerDefinitionArgs{
			WorkspaceId: createdWorkspace.ID(),
			Definition:  pulumi.String(spec.AlertManagerDefinition),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "create alert manager definition")
		}
	}

	namespaceArns := pulumi.StringMap{}
	for _, namespace := range spec.RuleGroupNamespaces {
		createdNamespace, err := amp.NewRuleGroupNamespace(ctx, "rule_group_namespace-"+namespace.Name, &amp.RuleGroupNamespaceArgs{
			WorkspaceId: createdWorkspace.ID(),
			Name:        pulumi.String(namespace.Name),
			Data:        pulumi.String(namespace.Data),
			Tags:        pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create rule group namespace %s", namespace.Name)
		}
		namespaceArns[namespace.Name] = createdNamespace.Arn
	}

	if queryLogging := spec.QueryLogging; queryLogging != nil {
		destinations := amp.QueryLoggingConfigurationDestinationArray{}
		for _, destination := range queryLogging.Destinations {
			destinations = append(destinations, &amp.QueryLoggingConfigurationDestinationArgs{
				CloudwatchLogs: &amp.QueryLoggingConfigurationDestinationCloudwatchLogsArgs{
					LogGroupArn: pulumi.String(wildcardLogGroupArn(destination.LogGroupArn.GetValue())),
				},
				Filters: &amp.QueryLoggingConfigurationDestinationFiltersArgs{
					QspThreshold: pulumi.Int(int(destination.QspThreshold)),
				},
			})
		}
		if _, err := amp.NewQueryLoggingConfiguration(ctx, "query_logging_configuration", &amp.QueryLoggingConfigurationArgs{
			WorkspaceId:  createdWorkspace.ID(),
			Destinations: destinations,
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "create query logging configuration")
		}
	}

	if spec.ResourcePolicy != nil {
		policyJson, err := json.Marshal(spec.ResourcePolicy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal resource policy")
		}
		if _, err := amp.NewResourcePolicy(ctx, "resource_policy", &amp.ResourcePolicyArgs{
			WorkspaceId:    createdWorkspace.ID(),
			PolicyDocument: pulumi.String(string(policyJson)),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "create resource policy")
		}
	}

	detectorIds := pulumi.StringMap{}
	detectorArns := pulumi.StringMap{}
	for _, detector := range spec.AnomalyDetectors {
		rcfArgs := &amp.AnomalyDetectorConfigurationRandomCutForestArgs{
			Query: pulumi.String(detector.Query),
		}
		if detector.SampleSize != nil {
			rcfArgs.SampleSize = pulumi.Int(int(*detector.SampleSize))
		}
		if detector.ShingleSize != nil {
			rcfArgs.ShingleSize = pulumi.Int(int(*detector.ShingleSize))
		}
		if band := detector.IgnoreNearExpectedFromAbove; band != nil {
			bandArgs := &amp.AnomalyDetectorConfigurationRandomCutForestIgnoreNearExpectedFromAboveArgs{}
			if band.Amount != nil {
				bandArgs.Amount = pulumi.Float64(*band.Amount)
			}
			if band.Ratio != nil {
				bandArgs.Ratio = pulumi.Float64(*band.Ratio)
			}
			rcfArgs.IgnoreNearExpectedFromAbove = bandArgs
		}
		if band := detector.IgnoreNearExpectedFromBelow; band != nil {
			bandArgs := &amp.AnomalyDetectorConfigurationRandomCutForestIgnoreNearExpectedFromBelowArgs{}
			if band.Amount != nil {
				bandArgs.Amount = pulumi.Float64(*band.Amount)
			}
			if band.Ratio != nil {
				bandArgs.Ratio = pulumi.Float64(*band.Ratio)
			}
			rcfArgs.IgnoreNearExpectedFromBelow = bandArgs
		}

		// The provider models the action as an exactly-one pair of
		// must-be-true bools; the spec models it honestly as an enum.
		missingDataActionArgs := &amp.AnomalyDetectorMissingDataActionArgs{}
		if detector.MissingDataAction == "MARK_AS_ANOMALY" {
			missingDataActionArgs.MarkAsAnomaly = pulumi.Bool(true)
		} else {
			missingDataActionArgs.Skip = pulumi.Bool(true)
		}

		detectorArgs := &amp.AnomalyDetectorArgs{
			WorkspaceId: createdWorkspace.ID(),
			Alias:       pulumi.String(detector.Alias),
			Configuration: &amp.AnomalyDetectorConfigurationArgs{
				RandomCutForest: rcfArgs,
			},
			MissingDataAction: missingDataActionArgs,
			Tags:              pulumi.ToStringMap(locals.AwsTags),
		}
		if detector.EvaluationIntervalInSeconds != nil {
			detectorArgs.EvaluationIntervalInSeconds = pulumi.Int(int(*detector.EvaluationIntervalInSeconds))
		}
		if len(detector.Labels) > 0 {
			detectorArgs.Labels = pulumi.ToStringMap(detector.Labels)
		}

		createdDetector, err := amp.NewAnomalyDetector(ctx, "anomaly_detector-"+detector.Alias, detectorArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create anomaly detector %s", detector.Alias)
		}
		detectorIds[detector.Alias] = createdDetector.ID().ToStringOutput()
		detectorArns[detector.Alias] = createdDetector.Arn
	}

	ctx.Export(OpWorkspaceId, createdWorkspace.ID())
	ctx.Export(OpWorkspaceArn, createdWorkspace.Arn)
	ctx.Export(OpPrometheusEndpoint, createdWorkspace.PrometheusEndpoint)
	ctx.Export(OpRuleGroupNamespaceArns, namespaceArns)
	ctx.Export(OpAnomalyDetectorIds, detectorIds)
	ctx.Export(OpAnomalyDetectorArns, detectorArns)
	return nil
}
