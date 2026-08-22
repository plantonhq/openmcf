package module

import (
	"github.com/pkg/errors"
	cloudflareemailroutingzonev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareemailroutingzone/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// emailRoutingZone enables Email Routing on the zone (which provisions the
// required DNS records), optionally folds in the single catch-all rule, and
// optionally locks the DNS records.
func emailRoutingZone(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareEmailRoutingZone.Spec

	zoneId := ""
	if spec.ZoneId != nil {
		zoneId = spec.ZoneId.GetValue()
	}

	settings, err := cloudflare.NewEmailRoutingSettings(
		ctx,
		"email-routing-settings",
		&cloudflare.EmailRoutingSettingsArgs{
			ZoneId: pulumi.String(zoneId),
		},
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to enable email routing settings")
	}

	if ca := spec.CatchAll; ca != nil {
		// Map each typed catch-all action onto the provider's generic
		// {type, value[]}: forward -> destination addresses; worker -> the
		// single script name; drop -> no values.
		actions := cloudflare.EmailRoutingCatchAllActionArray{}
		for _, a := range ca.Actions {
			values := pulumi.StringArray{}
			switch a.Type {
			case cloudflareemailroutingzonev1alpha1.CloudflareEmailRoutingCatchAllActionType_forward:
				for _, f := range a.ForwardTo {
					values = append(values, pulumi.String(f.GetValue()))
				}
			case cloudflareemailroutingzonev1alpha1.CloudflareEmailRoutingCatchAllActionType_worker:
				if a.Worker != nil {
					values = append(values, pulumi.String(a.Worker.GetValue()))
				}
			}
			actions = append(actions, cloudflare.EmailRoutingCatchAllActionArgs{
				Type:   pulumi.String(a.Type.String()),
				Values: values,
			})
		}

		catchAllArgs := &cloudflare.EmailRoutingCatchAllArgs{
			ZoneId:  pulumi.String(zoneId),
			Enabled: pulumi.Bool(ca.Enabled),
			// "all" is the only matcher type Cloudflare permits on the
			// catch-all, so the module supplies it.
			Matchers: cloudflare.EmailRoutingCatchAllMatcherArray{
				cloudflare.EmailRoutingCatchAllMatcherArgs{Type: pulumi.String("all")},
			},
			Actions: actions,
		}
		if ca.Name != "" {
			catchAllArgs.Name = pulumi.String(ca.Name)
		}

		// NOTE: the provider's Delete for this resource is a genuine no-op (no
		// API call) -- destroying it drops it from state and the zone keeps its
		// last catch-all configuration. The zone-level destroy (disabling Email
		// Routing) is what actually retires the behavior.
		_, err := cloudflare.NewEmailRoutingCatchAll(
			ctx,
			"email-routing-catch-all",
			catchAllArgs,
			pulumi.Provider(cloudflareProvider),
			pulumi.DependsOn([]pulumi.Resource{settings}),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create email routing catch-all")
		}
	}

	if spec.LockDnsRecords {
		// dns_name targets a subdomain of the zone; empty routes the apex.
		dnsArgs := &cloudflare.EmailRoutingDnsArgs{
			ZoneId: pulumi.String(zoneId),
		}
		if spec.DnsName != "" {
			dnsArgs.Name = pulumi.String(spec.DnsName)
		}
		_, err := cloudflare.NewEmailRoutingDns(
			ctx,
			"email-routing-dns",
			dnsArgs,
			pulumi.Provider(cloudflareProvider),
			pulumi.DependsOn([]pulumi.Resource{settings}),
		)
		if err != nil {
			return errors.Wrap(err, "failed to lock email routing dns records")
		}
	}

	ctx.Export(OpZoneId, pulumi.String(zoneId))
	ctx.Export(OpEnabled, settings.Enabled)
	ctx.Export(OpStatus, settings.Status)
	ctx.Export(OpName, settings.Name)

	return nil
}
