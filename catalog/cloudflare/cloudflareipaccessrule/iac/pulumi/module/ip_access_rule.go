package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ipAccessRule creates the IP Access rule on exactly one scope (the spec's CEL
// guarantees account XOR zone -- the provider would silently prefer the account
// if both arrived, so the module never sends both).
//
// API limitation taught by the provider's own tests: only mode and notes update
// in place; a configuration (target/value) change does not stick -- recreate the
// rule to change what it matches.
func ipAccessRule(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareIpAccessRule.Spec

	args := &cloudflare.AccessRuleArgs{
		Mode: pulumi.String(spec.Mode),
		Configuration: &cloudflare.AccessRuleConfigurationArgs{
			Target: pulumi.StringPtr(spec.Configuration.Target),
			Value:  pulumi.StringPtr(spec.Configuration.Value),
		},
	}

	if spec.AccountId != "" {
		args.AccountId = pulumi.StringPtr(spec.AccountId)
	}
	if spec.ZoneId.GetValue() != "" {
		args.ZoneId = pulumi.StringPtr(spec.ZoneId.GetValue())
	}
	if spec.Notes != "" {
		args.Notes = pulumi.String(spec.Notes)
	}

	createdRule, err := cloudflare.NewAccessRule(
		ctx,
		"ip_access_rule",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create ip access rule")
	}

	ctx.Export(OpRuleId, createdRule.ID())
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))
	ctx.Export(OpAccountId, pulumi.String(spec.AccountId))

	return nil
}
