package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	pkgerrors "github.com/pkg/errors"
)

// ecrRegistrySettingsVerifier verifies an AwsEcrRegistrySettings - a
// SETTINGS SINGLETON whose identity is the region's registry (the
// account id). The registry itself always exists, so existence rides
// the outputs: each keyed collection entry (cache rules, creation
// templates, exclusions) is asserted present. Absence is the arms'
// RESET semantics, not object deletion: after destroy the registry
// policy must be gone, replication must have zero rules, and the
// keyed collections must be empty. Account settings PERSIST after
// destroy by AWS design - deliberately not asserted absent.
type ecrRegistrySettingsVerifier struct{}

func (*ecrRegistrySettingsVerifier) IDOutputKey() string { return "registry_id" }

func (*ecrRegistrySettingsVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := ecr.NewFromConfig(cfg).DescribeRegistry(ctx, &ecr.DescribeRegistryInput{})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecrregistrysettings verify-exists failed for %q", id)
	}
	if out.RegistryId == nil || *out.RegistryId != id {
		return pkgerrors.Errorf("awsecrregistrysettings registry id %q does not match deployed %q", aws.ToString(out.RegistryId), id)
	}
	return nil
}

// VerifyAbsent asserts the RESET posture: policy gone, replication
// emptied, and no cache rules / templates / exclusions left. The
// registry object itself (and the persisted account settings) remain
// by AWS design.
func (*ecrRegistrySettingsVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := ecr.NewFromConfig(cfg)

	if _, err := client.GetRegistryPolicy(ctx, &ecr.GetRegistryPolicyInput{}); err == nil {
		return pkgerrors.Errorf("awsecrregistrysettings %q registry policy still exists after destroy", id)
	} else {
		var notFound *ecrtypes.RegistryPolicyNotFoundException
		if !pkgerrors.As(err, &notFound) {
			return pkgerrors.Wrapf(err, "awsecrregistrysettings %q policy verify-absent failed", id)
		}
	}

	registry, err := client.DescribeRegistry(ctx, &ecr.DescribeRegistryInput{})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecrregistrysettings %q replication verify-absent failed", id)
	}
	if registry.ReplicationConfiguration != nil && len(registry.ReplicationConfiguration.Rules) > 0 {
		return pkgerrors.Errorf("awsecrregistrysettings %q still carries %d replication rules after destroy", id, len(registry.ReplicationConfiguration.Rules))
	}

	cacheRules, err := client.DescribePullThroughCacheRules(ctx, &ecr.DescribePullThroughCacheRulesInput{})
	if err == nil && len(cacheRules.PullThroughCacheRules) > 0 {
		return pkgerrors.Errorf("awsecrregistrysettings %q still carries %d pull-through cache rules after destroy", id, len(cacheRules.PullThroughCacheRules))
	}

	templates, err := client.DescribeRepositoryCreationTemplates(ctx, &ecr.DescribeRepositoryCreationTemplatesInput{})
	if err == nil && len(templates.RepositoryCreationTemplates) > 0 {
		return pkgerrors.Errorf("awsecrregistrysettings %q still carries %d repository creation templates after destroy", id, len(templates.RepositoryCreationTemplates))
	}

	return nil
}

// VerifyExistsFromOutputs walks each keyed collection: every declared
// cache rule, creation template, and pull-time exclusion must exist
// under its recorded key.
func (v *ecrRegistrySettingsVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	registryId, _ := outputs["registry_id"].(string)
	if registryId == "" {
		return pkgerrors.New("awsecrregistrysettings outputs carry no registry_id")
	}
	if err := v.VerifyExists(ctx, cfg, registryId, region); err != nil {
		return err
	}
	client := ecr.NewFromConfig(cfg)

	if cacheRuleIds, _ := outputs["pull_through_cache_rule_registry_ids"].(map[string]interface{}); len(cacheRuleIds) > 0 {
		prefixes := make([]string, 0, len(cacheRuleIds))
		for prefix := range cacheRuleIds {
			prefixes = append(prefixes, prefix)
		}
		out, err := client.DescribePullThroughCacheRules(ctx, &ecr.DescribePullThroughCacheRulesInput{
			EcrRepositoryPrefixes: prefixes,
		})
		if err != nil {
			return pkgerrors.Wrap(err, "awsecrregistrysettings cache-rule read failed")
		}
		if len(out.PullThroughCacheRules) != len(prefixes) {
			return pkgerrors.Errorf("awsecrregistrysettings declares %d cache rules, AWS reports %d", len(prefixes), len(out.PullThroughCacheRules))
		}
	}

	if templateIds, _ := outputs["repository_creation_template_registry_ids"].(map[string]interface{}); len(templateIds) > 0 {
		prefixes := make([]string, 0, len(templateIds))
		for prefix := range templateIds {
			prefixes = append(prefixes, prefix)
		}
		out, err := client.DescribeRepositoryCreationTemplates(ctx, &ecr.DescribeRepositoryCreationTemplatesInput{
			Prefixes: prefixes,
		})
		if err != nil {
			return pkgerrors.Wrap(err, "awsecrregistrysettings creation-template read failed")
		}
		if len(out.RepositoryCreationTemplates) != len(prefixes) {
			return pkgerrors.Errorf("awsecrregistrysettings declares %d creation templates, AWS reports %d", len(prefixes), len(out.RepositoryCreationTemplates))
		}
	}

	if exclusionArns, _ := outputs["pull_time_update_exclusion_arns"].(map[string]interface{}); len(exclusionArns) > 0 {
		out, err := client.ListPullTimeUpdateExclusions(ctx, &ecr.ListPullTimeUpdateExclusionsInput{})
		if err != nil {
			return pkgerrors.Wrap(err, "awsecrregistrysettings exclusion read failed")
		}
		listed := map[string]bool{}
		for _, exclusionArn := range out.PullTimeUpdateExclusions {
			listed[exclusionArn] = true
		}
		for arn := range exclusionArns {
			if !listed[arn] {
				return pkgerrors.Errorf("awsecrregistrysettings pull-time exclusion %q not found after deploy", arn)
			}
		}
	}

	return nil
}
