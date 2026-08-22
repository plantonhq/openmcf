package module

import (
	"github.com/pkg/errors"
	cloudflarednszonev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarednszone/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// hold places a hold on the zone's hostname: while active, Cloudflare rejects
// attempts to create this hostname (and optionally any subdomain) as a zone in
// any other account. An unset hold_after is left nil: the provider treats an
// empty string as "hold from the current time", which would produce a perpetual
// diff against an unset spec field.
func hold(
	ctx *pulumi.Context,
	resourceName string,
	zone *cloudflare.Zone,
	spec *cloudflarednszonev1alpha1.CloudflareDnsZoneHold,
	cloudflareProvider *cloudflare.Provider,
) error {
	args := &cloudflare.ZoneHoldArgs{
		ZoneId:            zone.ID(),
		IncludeSubdomains: pulumi.Bool(spec.IncludeSubdomains),
	}
	if spec.HoldAfter != "" {
		args.HoldAfter = pulumi.String(spec.HoldAfter)
	}

	if _, err := cloudflare.NewZoneHold(
		ctx,
		resourceName+"-hold",
		args,
		pulumi.Provider(cloudflareProvider),
	); err != nil {
		return errors.Wrap(err, "failed to place zone hold")
	}
	return nil
}
