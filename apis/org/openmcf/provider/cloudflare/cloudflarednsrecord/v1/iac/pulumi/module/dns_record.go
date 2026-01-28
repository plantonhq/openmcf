package module

import (
	"strings"

	"github.com/pkg/errors"
	cloudflarednsrecordv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare/cloudflarednsrecord/v1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// recordTypeToString converts the proto enum to the Cloudflare API string.
func recordTypeToString(recordType cloudflarednsrecordv1.CloudflareDnsRecordSpec_RecordType) string {
	switch recordType {
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_A:
		return "A"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_AAAA:
		return "AAAA"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_CNAME:
		return "CNAME"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_MX:
		return "MX"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_TXT:
		return "TXT"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_SRV:
		return "SRV"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_NS:
		return "NS"
	case cloudflarednsrecordv1.CloudflareDnsRecordSpec_CAA:
		return "CAA"
	default:
		return "A" // Default fallback
	}
}

// dnsRecord provisions the Cloudflare DNS record and exports outputs.
func dnsRecord(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) (*cloudflare.DnsRecord, error) {
	spec := locals.CloudflareDnsRecord.Spec
	recordTypeStr := recordTypeToString(spec.Type)

	// Determine TTL (1 = automatic, or specified value).
	ttl := float64(1)
	if spec.Ttl > 0 {
		ttl = float64(spec.Ttl)
	}

	// Resolve zone_id from literal value or reference.
	zoneId := ""
	if spec.ZoneId != nil {
		zoneId = spec.ZoneId.GetValue()
	}

	// Build the record arguments.
	recordArgs := &cloudflare.DnsRecordArgs{
		ZoneId:  pulumi.String(zoneId),
		Name:    pulumi.String(spec.Name),
		Type:    pulumi.String(recordTypeStr),
		Content: pulumi.String(spec.Value),
		Proxied: pulumi.Bool(spec.Proxied),
		Ttl:     pulumi.Float64(ttl),
	}

	// Set priority for MX/SRV records.
	if spec.Type == cloudflarednsrecordv1.CloudflareDnsRecordSpec_MX ||
		spec.Type == cloudflarednsrecordv1.CloudflareDnsRecordSpec_SRV {
		recordArgs.Priority = pulumi.Float64(float64(spec.Priority))
	}

	// Set comment if provided.
	if spec.Comment != "" {
		recordArgs.Comment = pulumi.String(spec.Comment)
	}

	// Create the DNS record using DnsRecord (Record is deprecated in v6).
	createdRecord, err := cloudflare.NewDnsRecord(
		ctx,
		strings.ToLower(locals.CloudflareDnsRecord.Metadata.Name),
		recordArgs,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudflare dns record")
	}

	// Export required outputs.
	ctx.Export(OpRecordId, createdRecord.ID())
	ctx.Export(OpHostname, createdRecord.Name)
	ctx.Export(OpRecordType, pulumi.String(recordTypeStr))
	ctx.Export(OpProxied, createdRecord.Proxied)

	return createdRecord, nil
}
