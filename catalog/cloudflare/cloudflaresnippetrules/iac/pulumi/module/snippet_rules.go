package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// snippetRules applies the zone's snippet routing table. A zone SINGLETON:
// Cloudflare's API replaces the WHOLE list on every apply, so the spec's list
// is the entire table, and destroy deletes ALL snippet rules in the zone --
// including rules created outside this resource (the snippets survive).
//
// enabled defaults to TRUE here even though Cloudflare's own default is FALSE:
// a declared rule should run. The explicit else-branch (rather than omitting
// the field) is what makes the spec's promised default real on this engine.
func snippetRules(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareSnippetRules.Spec

	rules := make(cloudflare.SnippetRulesRuleArray, 0, len(spec.Rules))
	for _, rule := range spec.Rules {
		ruleArgs := &cloudflare.SnippetRulesRuleArgs{
			Expression:  pulumi.String(rule.Expression),
			SnippetName: pulumi.String(rule.SnippetName.GetValue()),
		}
		if rule.Description != "" {
			ruleArgs.Description = pulumi.StringPtr(rule.Description)
		}
		if rule.Enabled != nil {
			ruleArgs.Enabled = pulumi.BoolPtr(rule.GetEnabled())
		} else {
			ruleArgs.Enabled = pulumi.BoolPtr(true)
		}
		rules = append(rules, ruleArgs)
	}

	args := &cloudflare.SnippetRulesArgs{
		ZoneId: pulumi.String(spec.ZoneId.GetValue()),
		Rules:  rules,
	}

	_, err := cloudflare.NewSnippetRules(
		ctx,
		"snippet_rules",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to apply snippet rules")
	}

	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
