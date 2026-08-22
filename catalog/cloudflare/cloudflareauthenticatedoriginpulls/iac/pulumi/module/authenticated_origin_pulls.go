package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authenticatedOriginPulls manages the zone-wide toggle and the per-hostname
// certificate associations. Destroy semantics differ per surface and neither
// is a delete: the zone toggle has NO delete at Cloudflare (destroy abandons
// the live value), and an association is removed by a revert write
// (enabled: null) the provider issues from state.
func authenticatedOriginPulls(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareAuthenticatedOriginPulls.Spec

	// The zone-wide toggle. Managed only when the spec sets zone_enabled --
	// unset means "leave the zone's toggle alone".
	if spec.ZoneEnabled != nil {
		_, err := cloudflare.NewAuthenticatedOriginPullsSettings(
			ctx,
			"zone_toggle",
			&cloudflare.AuthenticatedOriginPullsSettingsArgs{
				ZoneId:  pulumi.String(spec.ZoneId.GetValue()),
				Enabled: pulumi.Bool(spec.GetZoneEnabled()),
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create zone toggle")
		}
	}

	// One association resource per hostname row. The provider requires the
	// config list to hold exactly one element per resource (it hard-fails
	// otherwise), so each row fans out to its own resource. An omitted
	// enabled is sent as true: Cloudflare treats null as "void the
	// association", and a declared row is meant to exist.
	for _, association := range spec.HostnameAssociations {
		configArgs := cloudflare.AuthenticatedOriginPullsConfigArgs{
			Hostname: pulumi.StringPtr(association.Hostname),
			Enabled:  pulumi.BoolPtr(association.Enabled == nil || association.GetEnabled()),
		}
		if association.CertificateId.GetValue() != "" {
			configArgs.CertId = pulumi.StringPtr(association.CertificateId.GetValue())
		}

		_, err := cloudflare.NewAuthenticatedOriginPulls(
			ctx,
			"association_"+association.Hostname,
			&cloudflare.AuthenticatedOriginPullsArgs{
				ZoneId:  pulumi.String(spec.ZoneId.GetValue()),
				Configs: cloudflare.AuthenticatedOriginPullsConfigArray{configArgs},
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to create association for hostname %s", association.Hostname)
		}
	}

	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	return nil
}
