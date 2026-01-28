package module

import (
	"strings"

	"github.com/pkg/errors"
	civodnsrecordv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/civo/civodnsrecord/v1"
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// recordTypeToString converts the proto enum to the Civo API string.
func recordTypeToString(recordType civodnsrecordv1.CivoDnsRecordSpec_RecordType) string {
	switch recordType {
	case civodnsrecordv1.CivoDnsRecordSpec_A:
		return "A"
	case civodnsrecordv1.CivoDnsRecordSpec_AAAA:
		return "AAAA"
	case civodnsrecordv1.CivoDnsRecordSpec_CNAME:
		return "CNAME"
	case civodnsrecordv1.CivoDnsRecordSpec_MX:
		return "MX"
	case civodnsrecordv1.CivoDnsRecordSpec_TXT:
		return "TXT"
	case civodnsrecordv1.CivoDnsRecordSpec_SRV:
		return "SRV"
	case civodnsrecordv1.CivoDnsRecordSpec_NS:
		return "NS"
	default:
		return "A" // Default fallback
	}
}

// dnsRecord provisions the Civo DNS record and exports outputs.
func dnsRecord(
	ctx *pulumi.Context,
	locals *Locals,
	civoProvider *civo.Provider,
) (*civo.DnsDomainRecord, error) {
	spec := locals.CivoDnsRecord.Spec
	recordTypeStr := recordTypeToString(spec.Type)

	// Determine TTL (default to 3600 if not specified).
	ttl := 3600
	if spec.Ttl > 0 {
		ttl = int(spec.Ttl)
	}

	// Build the record arguments.
	recordArgs := &civo.DnsDomainRecordArgs{
		DomainId: pulumi.String(locals.ZoneId),
		Name:     pulumi.String(spec.Name),
		Type:     pulumi.String(recordTypeStr),
		Value:    pulumi.String(spec.Value),
		Ttl:      pulumi.Int(ttl),
	}

	// Set priority for MX/SRV records.
	if spec.Type == civodnsrecordv1.CivoDnsRecordSpec_MX ||
		spec.Type == civodnsrecordv1.CivoDnsRecordSpec_SRV {
		recordArgs.Priority = pulumi.Int(int(spec.Priority))
	}

	// Create the DNS record.
	createdRecord, err := civo.NewDnsDomainRecord(
		ctx,
		strings.ToLower(locals.CivoDnsRecord.Metadata.Name),
		recordArgs,
		pulumi.Provider(civoProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create civo dns record")
	}

	// Export required outputs.
	ctx.Export(OpRecordId, createdRecord.ID())
	ctx.Export(OpHostname, createdRecord.Name)
	ctx.Export(OpRecordType, pulumi.String(recordTypeStr))
	ctx.Export(OpAccountId, createdRecord.AccountId)

	return createdRecord, nil
}
