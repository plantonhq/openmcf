package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dnsZone provisions the DigitalOcean domain plus its managed DNS records and
// exports stack outputs.
func dnsZone(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.Domain, error) {
	spec := locals.DigitalOceanDnsZone.Spec

	domainArgs := &digitalocean.DomainArgs{
		Name: pulumi.String(spec.DomainName),
	}
	// ip_address is a create-only convenience that seeds an initial apex A
	// record DigitalOcean never tracks afterwards — prefer declaring records.
	if spec.IpAddress != "" {
		domainArgs.IpAddress = pulumi.StringPtr(spec.IpAddress)
	}

	createdDomain, err := digitalocean.NewDomain(
		ctx,
		"dns_zone",
		domainArgs,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean domain")
	}

	// One DigitalOcean record per record value (same name and type). The
	// per-type fields carry the spec's presence semantics: unset stays
	// unset, matching the provider's own GetOk-guarded request building.
	for recIdx, rec := range spec.Records {
		for valIdx, val := range rec.Values {
			// The spec's shared enum value names ARE the DigitalOcean record
			// types (A, AAAA, CNAME, ...), so the type wires through directly.
			recordArgs := &digitalocean.DnsRecordArgs{
				Domain: createdDomain.Name,
				Name:   pulumi.String(rec.Name),
				Type:   pulumi.String(rec.Type.String()),
				Value:  pulumi.String(val.GetValue()),
			}

			// 0 means unset: the ttl attribute is then Computed and
			// DigitalOcean applies its default (1800 seconds).
			if rec.TtlSeconds > 0 {
				recordArgs.Ttl = pulumi.IntPtr(int(rec.TtlSeconds))
			}
			if rec.Priority != nil {
				recordArgs.Priority = pulumi.IntPtr(int(*rec.Priority))
			}
			if rec.Weight != nil {
				recordArgs.Weight = pulumi.IntPtr(int(*rec.Weight))
			}
			if rec.Port != nil {
				recordArgs.Port = pulumi.IntPtr(int(*rec.Port))
			}
			if rec.Flags != nil {
				recordArgs.Flags = pulumi.IntPtr(int(*rec.Flags))
			}
			if rec.Tag != "" {
				recordArgs.Tag = pulumi.StringPtr(rec.Tag)
			}

			resourceName := fmt.Sprintf("%s-%d-%d", rec.Name, recIdx, valIdx)
			if _, err := digitalocean.NewDnsRecord(
				ctx,
				resourceName,
				recordArgs,
				pulumi.Provider(digitalOceanProvider),
			); err != nil {
				return nil, errors.Wrapf(err, "failed to create dns record %s", resourceName)
			}
		}
	}

	ctx.Export(OpZoneName, createdDomain.Name)
	ctx.Export(OpZoneId, createdDomain.ID())
	// DigitalOcean's authoritative name servers are a fixed platform-wide set
	// the API does not return per zone.
	ctx.Export(OpNameServers, pulumi.StringArray{
		pulumi.String("ns1.digitalocean.com"),
		pulumi.String("ns2.digitalocean.com"),
		pulumi.String("ns3.digitalocean.com"),
	})
	// The SDK renames the provider's urn attribute to domainUrn.
	ctx.Export(OpUrn, createdDomain.DomainUrn)

	return createdDomain, nil
}
