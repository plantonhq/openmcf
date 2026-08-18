package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// registry manages the region's registry-level ECR configuration and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the policy/scanning/replication arms are one-per-registry
//     singletons keyed by the account id; scanning and replication
//     have RESET-not-delete semantics at the provider (destroy puts
//     the empty default back);
//   - account settings are PutAccountSetting upserts with a NO-OP
//     delete - the last-applied values persist after destroy;
//   - pull-through cache rules and creation templates are prefix-keyed
//     (ForceNew prefixes - the for_each keys below); clearing a cache
//     rule's credential/custom-role ARN back to empty is NOT
//     propagated by the provider (replace the rule to drop
//     credentials);
//   - pull-time update exclusions are immutable per principal ARN.
func registry(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// The registry's identity is the account itself - resolved once
	// and exported for imports and the pull-URL join key.
	callerIdentity, err := aws.GetCallerIdentity(ctx, nil, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "resolve caller identity")
	}
	registryId := callerIdentity.AccountId

	if spec.RegistryPolicy != "" {
		if _, err := ecr.NewRegistryPolicy(ctx, "registry-policy", &ecr.RegistryPolicyArgs{
			Policy: pulumi.String(spec.RegistryPolicy),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "put registry policy")
		}
	}

	if scanning := spec.Scanning; scanning != nil {
		rules := ecr.RegistryScanningConfigurationRuleArray{}
		for _, rule := range scanning.Rules {
			filters := ecr.RegistryScanningConfigurationRuleRepositoryFilterArray{}
			for _, filter := range rule.Filters {
				filters = append(filters, &ecr.RegistryScanningConfigurationRuleRepositoryFilterArgs{
					Filter: pulumi.String(filter),
					// WILDCARD is the only filter type AWS supports -
					// pinned here, never spec surface.
					FilterType: pulumi.String("WILDCARD"),
				})
			}
			rules = append(rules, &ecr.RegistryScanningConfigurationRuleArgs{
				ScanFrequency:     pulumi.String(rule.ScanFrequency),
				RepositoryFilters: filters,
			})
		}
		if _, err := ecr.NewRegistryScanningConfiguration(ctx, "scanning", &ecr.RegistryScanningConfigurationArgs{
			ScanType: pulumi.String(scanning.ScanType),
			Rules:    rules,
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "put scanning configuration")
		}
	}

	if len(spec.ReplicationRules) > 0 {
		rules := ecr.ReplicationConfigurationReplicationConfigurationRuleArray{}
		for _, rule := range spec.ReplicationRules {
			destinations := ecr.ReplicationConfigurationReplicationConfigurationRuleDestinationArray{}
			for _, destination := range rule.Destinations {
				destinations = append(destinations, &ecr.ReplicationConfigurationReplicationConfigurationRuleDestinationArgs{
					Region:     pulumi.String(destination.Region),
					RegistryId: pulumi.String(destination.RegistryId),
				})
			}
			filters := ecr.ReplicationConfigurationReplicationConfigurationRuleRepositoryFilterArray{}
			for _, filter := range rule.RepositoryFilters {
				filters = append(filters, &ecr.ReplicationConfigurationReplicationConfigurationRuleRepositoryFilterArgs{
					Filter: pulumi.String(filter),
					// PREFIX_MATCH is the only filter type AWS
					// supports - pinned here, never spec surface.
					FilterType: pulumi.String("PREFIX_MATCH"),
				})
			}
			rules = append(rules, &ecr.ReplicationConfigurationReplicationConfigurationRuleArgs{
				Destinations:      destinations,
				RepositoryFilters: filters,
			})
		}
		if _, err := ecr.NewReplicationConfiguration(ctx, "replication", &ecr.ReplicationConfigurationArgs{
			ReplicationConfiguration: &ecr.ReplicationConfigurationReplicationConfigurationArgs{
				Rules: rules,
			},
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "put replication configuration")
		}
	}

	cacheRuleRegistryIds := pulumi.StringMap{}
	for _, cacheRule := range spec.PullThroughCacheRules {
		cacheRuleArgs := &ecr.PullThroughCacheRuleArgs{
			EcrRepositoryPrefix: pulumi.String(cacheRule.EcrRepositoryPrefix),
			UpstreamRegistryUrl: pulumi.String(cacheRule.UpstreamRegistryUrl),
		}
		if cacheRule.UpstreamRepositoryPrefix != "" {
			cacheRuleArgs.UpstreamRepositoryPrefix = pulumi.String(cacheRule.UpstreamRepositoryPrefix)
		}
		if cacheRule.CredentialArn != nil && cacheRule.CredentialArn.GetValue() != "" {
			cacheRuleArgs.CredentialArn = pulumi.String(cacheRule.CredentialArn.GetValue())
		}
		if cacheRule.CustomRoleArn != nil && cacheRule.CustomRoleArn.GetValue() != "" {
			cacheRuleArgs.CustomRoleArn = pulumi.String(cacheRule.CustomRoleArn.GetValue())
		}
		createdCacheRule, err := ecr.NewPullThroughCacheRule(ctx,
			fmt.Sprintf("cache-rule-%s", cacheRule.EcrRepositoryPrefix),
			cacheRuleArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create cache rule %s", cacheRule.EcrRepositoryPrefix)
		}
		cacheRuleRegistryIds[cacheRule.EcrRepositoryPrefix] = createdCacheRule.RegistryId
	}

	templateRegistryIds := pulumi.StringMap{}
	for _, template := range spec.RepositoryCreationTemplates {
		templateArgs := &ecr.RepositoryCreationTemplateArgs{
			Prefix:      pulumi.String(template.Prefix),
			AppliedFors: pulumi.ToStringArray(template.AppliedFor),
		}
		if template.Description != "" {
			templateArgs.Description = pulumi.String(template.Description)
		}
		if template.CustomRoleArn != nil && template.CustomRoleArn.GetValue() != "" {
			templateArgs.CustomRoleArn = pulumi.String(template.CustomRoleArn.GetValue())
		}
		if template.ImageTagMutability != "" {
			templateArgs.ImageTagMutability = pulumi.String(template.ImageTagMutability)
		}
		if len(template.ImageTagMutabilityExclusionFilters) > 0 {
			filters := ecr.RepositoryCreationTemplateImageTagMutabilityExclusionFilterArray{}
			for _, filter := range template.ImageTagMutabilityExclusionFilters {
				filters = append(filters, &ecr.RepositoryCreationTemplateImageTagMutabilityExclusionFilterArgs{
					Filter: pulumi.String(filter),
					// WILDCARD is the only filter type AWS supports.
					FilterType: pulumi.String("WILDCARD"),
				})
			}
			templateArgs.ImageTagMutabilityExclusionFilters = filters
		}
		if encryption := template.Encryption; encryption != nil {
			encryptionArgs := &ecr.RepositoryCreationTemplateEncryptionConfigurationArgs{}
			if encryption.Type != "" {
				encryptionArgs.EncryptionType = pulumi.String(encryption.Type)
			}
			if encryption.KmsKey != nil && encryption.KmsKey.GetValue() != "" {
				encryptionArgs.KmsKey = pulumi.String(encryption.KmsKey.GetValue())
			}
			templateArgs.EncryptionConfigurations = ecr.RepositoryCreationTemplateEncryptionConfigurationArray{encryptionArgs}
		}
		if template.LifecyclePolicy != "" {
			templateArgs.LifecyclePolicy = pulumi.String(template.LifecyclePolicy)
		}
		if template.RepositoryPolicy != "" {
			templateArgs.RepositoryPolicy = pulumi.String(template.RepositoryPolicy)
		}
		if len(template.ResourceTags) > 0 {
			templateArgs.ResourceTags = pulumi.ToStringMap(template.ResourceTags)
		}
		createdTemplate, err := ecr.NewRepositoryCreationTemplate(ctx,
			fmt.Sprintf("creation-template-%s", template.Prefix),
			templateArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create creation template %s", template.Prefix)
		}
		templateRegistryIds[template.Prefix] = createdTemplate.RegistryId
	}

	if settings := spec.AccountSettings; settings != nil {
		// Each toggle is its own PutAccountSetting upsert; only the
		// configured ones are managed (unset toggles keep the
		// account's current values - and all of them PERSIST after
		// destroy).
		settingValues := map[string]string{}
		if settings.BasicScanTypeVersion != "" {
			settingValues["BASIC_SCAN_TYPE_VERSION"] = settings.BasicScanTypeVersion
		}
		if settings.BlobMounting != nil {
			if *settings.BlobMounting {
				settingValues["BLOB_MOUNTING"] = "ENABLED"
			} else {
				settingValues["BLOB_MOUNTING"] = "DISABLED"
			}
		}
		if settings.RegistryPolicyScope != "" {
			settingValues["REGISTRY_POLICY_SCOPE"] = settings.RegistryPolicyScope
		}
		for _, settingName := range []string{"BASIC_SCAN_TYPE_VERSION", "BLOB_MOUNTING", "REGISTRY_POLICY_SCOPE"} {
			settingValue, configured := settingValues[settingName]
			if !configured {
				continue
			}
			if _, err := ecr.NewAccountSetting(ctx,
				fmt.Sprintf("account-setting-%s", settingName),
				&ecr.AccountSettingArgs{
					Name:  pulumi.String(settingName),
					Value: pulumi.String(settingValue),
				}, pulumi.Provider(provider)); err != nil {
				return errors.Wrapf(err, "put account setting %s", settingName)
			}
		}
	}

	exclusionArns := pulumi.StringMap{}
	for _, principal := range spec.PullTimeUpdateExclusions {
		resolvedArn := principal.GetValue()
		createdExclusion, err := ecr.NewPullTimeUpdateExclusion(ctx,
			fmt.Sprintf("pull-time-exclusion-%s", resolvedArn),
			&ecr.PullTimeUpdateExclusionArgs{
				PrincipalArn: pulumi.String(resolvedArn),
			}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "register pull-time exclusion %s", resolvedArn)
		}
		exclusionArns[resolvedArn] = createdExclusion.PrincipalArn
	}

	ctx.Export(OpRegistryId, pulumi.String(registryId))
	ctx.Export(OpRegistryUrl, pulumi.Sprintf("%s.dkr.ecr.%s.amazonaws.com", registryId, spec.Region))
	ctx.Export(OpPullThroughCacheRuleRegistryIds, cacheRuleRegistryIds)
	ctx.Export(OpRepositoryCreationTemplateRegistryIds, templateRegistryIds)
	ctx.Export(OpPullTimeUpdateExclusionArns, exclusionArns)
	return nil
}
