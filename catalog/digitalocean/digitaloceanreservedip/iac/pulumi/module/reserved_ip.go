package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// reservedIpv4 reserves an IPv4 address, assigning it through the
// resource's own droplet_id argument -- the provider updates it in place
// (assign / re-point / unassign without ever replacing the address).
func reservedIpv4(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) error {
	spec := locals.DigitalOceanReservedIp.Spec

	args := &digitalocean.ReservedIpArgs{
		// The provider lowercases the region into state; the spec's region
		// enum value names are already the lowercase slugs.
		Region: pulumi.String(spec.Region.String()),
	}
	if locals.DropletId != nil {
		args.DropletId = pulumi.IntPtr(*locals.DropletId)
	}

	createdIp, err := digitalocean.NewReservedIp(
		ctx,
		"reserved-ip",
		args,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create digitalocean reserved ip")
	}

	ctx.Export(OpReservedIpAddress, createdIp.IpAddress)
	// The provider's urn attribute is renamed by the Pulumi bridge (it
	// collides with Pulumi's own resource URN).
	ctx.Export(OpUrn, createdIp.ReservedIpUrn)

	return nil
}

// reservedIpv6 reserves an IPv6 address. The v6 resource CANNOT assign
// inline (the provider's v6 create silently ignores a droplet id and has
// no update function), so assignment goes through the separate
// digitalocean_reserved_ipv6_assignment resource; re-pointing replaces
// just the assignment, never the address.
func reservedIpv6(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) error {
	spec := locals.DigitalOceanReservedIp.Spec

	createdIpv6, err := digitalocean.NewReservedIpv6(
		ctx,
		"reserved-ipv6",
		&digitalocean.ReservedIpv6Args{
			// Same region concept as v4; the v6 resource names the
			// argument region_slug.
			RegionSlug: pulumi.String(spec.Region.String()),
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create digitalocean reserved ipv6")
	}

	if locals.DropletId != nil {
		if _, err := digitalocean.NewReservedIpv6Assignment(
			ctx,
			"reserved-ipv6-assignment",
			&digitalocean.ReservedIpv6AssignmentArgs{
				Ip:        createdIpv6.Ip,
				DropletId: pulumi.Int(*locals.DropletId),
			},
			pulumi.Provider(digitalOceanProvider),
		); err != nil {
			return errors.Wrap(err, "failed to assign digitalocean reserved ipv6")
		}
	}

	ctx.Export(OpReservedIpAddress, createdIpv6.Ip)
	// Bridge rename, same as the v4 resource's urn.
	ctx.Export(OpUrn, createdIpv6.ReservedIpv6Urn)

	return nil
}
