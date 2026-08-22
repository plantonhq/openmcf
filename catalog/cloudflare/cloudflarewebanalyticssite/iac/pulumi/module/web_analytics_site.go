package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// webAnalyticsSite creates the Web Analytics (RUM) site plus its folded
// measurement rules. The site is identified by host OR zone (spec
// validation enforces exactly one).
//
// Cloudflare stores include/exclude rules as separate objects under the
// site's ruleset; this module manages one rule object per declared row,
// keyed by position. The provider never reads rules back after writing them
// (its refresh is deliberately blind), so each apply re-asserts exactly the
// declared rows.
func webAnalyticsSite(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareWebAnalyticsSite.Spec

	args := &cloudflare.WebAnalyticsSiteArgs{
		AccountId: pulumi.String(spec.AccountId),
	}

	if spec.Host != "" {
		args.Host = pulumi.StringPtr(spec.Host)
	}
	if spec.ZoneTag.GetValue() != "" {
		args.ZoneTag = pulumi.StringPtr(spec.ZoneTag.GetValue())
	}
	if spec.AutoInstall != nil {
		args.AutoInstall = pulumi.BoolPtr(spec.GetAutoInstall())
	}
	if spec.Enabled != nil {
		args.Enabled = pulumi.BoolPtr(spec.GetEnabled())
	}
	if spec.Lite != nil {
		args.Lite = pulumi.BoolPtr(spec.GetLite())
	}

	createdSite, err := cloudflare.NewWebAnalyticsSite(
		ctx,
		"web_analytics_site",
		args,
		pulumi.Provider(cloudflareProvider),
		// site_token and snippet carry the measurement credential; keep
		// them out of plain-text stack state.
		pulumi.AdditionalSecretOutputs([]string{"siteToken", "snippet"}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create web analytics site")
	}

	rulesetId := createdSite.Ruleset.Id().Elem()

	for index, row := range spec.Rules {
		ruleArgs := &cloudflare.WebAnalyticsRuleArgs{
			AccountId: pulumi.String(spec.AccountId),
			RulesetId: rulesetId,
		}
		if row.Host != "" {
			ruleArgs.Host = pulumi.StringPtr(row.Host)
		}
		if len(row.Paths) > 0 {
			ruleArgs.Paths = pulumi.ToStringArray(row.Paths)
		}
		if row.Inclusive != nil {
			ruleArgs.Inclusive = pulumi.BoolPtr(row.GetInclusive())
		}
		if row.IsPaused != nil {
			ruleArgs.IsPaused = pulumi.BoolPtr(row.GetIsPaused())
		}

		if _, err := cloudflare.NewWebAnalyticsRule(
			ctx,
			fmt.Sprintf("web_analytics_rule_%d", index),
			ruleArgs,
			pulumi.Provider(cloudflareProvider),
			pulumi.Parent(createdSite),
		); err != nil {
			return errors.Wrapf(err, "failed to create web analytics rule %d", index)
		}
	}

	ctx.Export(OpSiteTag, createdSite.SiteTag)
	ctx.Export(OpSiteToken, createdSite.SiteToken)
	ctx.Export(OpSnippet, createdSite.Snippet)
	ctx.Export(OpRulesetId, rulesetId)

	return nil
}
