package module

import (
	"strings"

	"github.com/pkg/errors"
	cloudflarednsrecordv1 "github.com/plantonhq/project-planton/apis/org/project_planton/provider/cloudflare/cloudflarednsrecord/v1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// recordTypeToString converts the proto enum to the Cloudflare API string.
func recordTypeToString(recordType cloudflarednsrecordv1.CloudflareDnsRecordType) string {
	switch recordType {
	case cloudflarednsrecordv1.CloudflareDnsRecordType_A:
		return "A"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_AAAA:
		return "AAAA"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_CNAME:
		return "CNAME"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_MX:
		return "MX"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_TXT:
		return "TXT"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_SRV:
		return "SRV"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_NS:
		return "NS"
	case cloudflarednsrecordv1.CloudflareDnsRecordType_CAA:
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

	// Build the record arguments.
	recordArgs := &cloudflare.DnsRecordArgs{
		ZoneId:  pulumi.String(spec.ZoneId),
		Name:    pulumi.String(spec.Name),
		Type:    pulumi.String(recordTypeStr),
		Content: pulumi.String(spec.Value),
		Proxied: pulumi.Bool(spec.Proxied),
		Ttl:     pulumi.Float64(ttl),
	}

	// Set priority for MX/SRV records.
	if spec.Type == cloudflarednsrecordv1.CloudflareDnsRecordType_MX ||
		spec.Type == cloudflarednsrecordv1.CloudflareDnsRecordType_SRV {
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
