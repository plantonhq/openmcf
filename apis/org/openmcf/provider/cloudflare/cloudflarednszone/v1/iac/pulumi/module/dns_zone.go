package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	cloudflarednszonev1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare/cloudflarednszone/v1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dnsZone provisions the Cloudflare zone and exports outputs.
func dnsZone(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) (*cloudflare.Zone, error) {

	// 1. Create the zone.
	createdZone, err := cloudflare.NewZone(
		ctx,
		// Use metadata.name as the resource label to mimic Terraform naming.
		strings.ToLower(locals.CloudflareDnsZone.Metadata.Name),
		// 2. Build the arguments straight from proto fields.
		// Note: Plan field was removed in v6 and is now set via a separate API call or zone settings
		&cloudflare.ZoneArgs{
			Account: cloudflare.ZoneAccountArgs{
				Id: pulumi.String(locals.CloudflareDnsZone.Spec.AccountId),
			},
			Name:   pulumi.String(locals.CloudflareDnsZone.Spec.ZoneName),
			Paused: pulumi.BoolPtr(locals.CloudflareDnsZone.Spec.Paused),
			// NOTE: default_proxied and plan are not available at zone‑level in the v6 provider.
			// Plan is now managed separately via zone settings or account configuration.
		},
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudflare zone")
	}

	// 3. Create DNS records if specified.
	if len(locals.CloudflareDnsZone.Spec.Records) > 0 {
		if err := createDnsRecords(ctx, createdZone, locals.CloudflareDnsZone.Spec.Records, cloudflareProvider); err != nil {
			return nil, errors.Wrap(err, "failed to create dns records")
		}
	}

	// 4. Export required outputs.
	ctx.Export(OpZoneId, createdZone.ID())
	ctx.Export(OpNameservers, createdZone.NameServers)

	return createdZone, nil
}

// createDnsRecords creates DNS records within the zone.
func createDnsRecords(
	ctx *pulumi.Context,
	zone *cloudflare.Zone,
	records []*cloudflarednszonev1.CloudflareDnsZoneRecord,
	cloudflareProvider *cloudflare.Provider,
) error {
	for _, record := range records {
		resourceName := fmt.Sprintf("%s-%s", record.Name, record.Type.String())

		recordArgs := &cloudflare.RecordArgs{
			ZoneId:  zone.ID(),
			Name:    pulumi.String(record.Name),
			Type:    pulumi.String(record.Type.String()),
			Content: pulumi.String(record.Value),
			Ttl:     pulumi.Float64(float64(record.Ttl)),
		}

		// proxied is only applicable to A, AAAA, and CNAME records
		if record.Type == cloudflarednszonev1.CloudflareDnsZoneRecord_A ||
			record.Type == cloudflarednszonev1.CloudflareDnsZoneRecord_AAAA ||
			record.Type == cloudflarednszonev1.CloudflareDnsZoneRecord_CNAME {
			recordArgs.Proxied = pulumi.Bool(record.Proxied)
		}

		// priority is only used for MX and SRV records
		if record.Type == cloudflarednszonev1.CloudflareDnsZoneRecord_MX ||
			record.Type == cloudflarednszonev1.CloudflareDnsZoneRecord_SRV {
			recordArgs.Priority = pulumi.Float64Ptr(float64(record.Priority))
		}

		// comment for the DNS record
		if record.Comment != "" {
			recordArgs.Comment = pulumi.String(record.Comment)
		}

		_, err := cloudflare.NewRecord(
			ctx,
			resourceName,
			recordArgs,
			pulumi.Provider(cloudflareProvider),
			pulumi.DependsOn([]pulumi.Resource{zone}),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to create dns record %s", resourceName)
		}
	}
	return nil
}
