package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// botManagement applies the zone's Bot Management configuration. The surface is
// a zone SINGLETON: create adopts whatever configuration the zone already
// carries (the API's create is a PUT), unset spec fields are never sent (the
// zone keeps its current values), and destroy is a NO-OP at Cloudflare -- the
// last-applied configuration stays live. Retire a setting by applying its off
// value BEFORE destroying.
//
// Plan gates are Cloudflare's, enforced at the API: setting a field the zone's
// plan does not include fails at apply, and on non-entitled zones the API omits
// those fields from responses (the provider's own issue tracker records refresh
// drift for that case). Manage only what the plan includes.
func botManagement(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareBotManagement.Spec

	args := &cloudflare.BotManagementArgs{
		ZoneId: pulumi.String(spec.ZoneId.GetValue()),
	}

	if spec.FightMode != nil {
		args.FightMode = pulumi.BoolPtr(spec.GetFightMode())
	}
	if spec.SbfmDefinitelyAutomated != nil {
		args.SbfmDefinitelyAutomated = pulumi.StringPtr(spec.GetSbfmDefinitelyAutomated())
	}
	if spec.SbfmLikelyAutomated != nil {
		args.SbfmLikelyAutomated = pulumi.StringPtr(spec.GetSbfmLikelyAutomated())
	}
	if spec.SbfmVerifiedBots != nil {
		args.SbfmVerifiedBots = pulumi.StringPtr(spec.GetSbfmVerifiedBots())
	}
	if spec.SbfmStaticResourceProtection != nil {
		args.SbfmStaticResourceProtection = pulumi.BoolPtr(spec.GetSbfmStaticResourceProtection())
	}
	if spec.OptimizeWordpress != nil {
		args.OptimizeWordpress = pulumi.BoolPtr(spec.GetOptimizeWordpress())
	}
	if spec.AutoUpdateModel != nil {
		args.AutoUpdateModel = pulumi.BoolPtr(spec.GetAutoUpdateModel())
	}
	if spec.SuppressSessionScore != nil {
		args.SuppressSessionScore = pulumi.BoolPtr(spec.GetSuppressSessionScore())
	}
	if spec.EnableJs != nil {
		args.EnableJs = pulumi.BoolPtr(spec.GetEnableJs())
	}
	if spec.BmCookieEnabled != nil {
		args.BmCookieEnabled = pulumi.BoolPtr(spec.GetBmCookieEnabled())
	}
	if spec.AiBotsProtection != nil {
		args.AiBotsProtection = pulumi.StringPtr(spec.GetAiBotsProtection())
	}
	if spec.CrawlerProtection != nil {
		args.CrawlerProtection = pulumi.StringPtr(spec.GetCrawlerProtection())
	}
	if spec.ContentBotsProtection != nil {
		args.ContentBotsProtection = pulumi.StringPtr(spec.GetContentBotsProtection())
	}
	if spec.CfRobotsVariant != nil {
		args.CfRobotsVariant = pulumi.StringPtr(spec.GetCfRobotsVariant())
	}
	if spec.IsRobotsTxtManaged != nil {
		args.IsRobotsTxtManaged = pulumi.BoolPtr(spec.GetIsRobotsTxtManaged())
	}

	_, err := cloudflare.NewBotManagement(
		ctx,
		"bot_management",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to apply bot management configuration")
	}

	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
