package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	pkgerrors "github.com/pkg/errors"
)

// route53DnsRecordVerifier verifies an AwsRoute53DnsRecord by listing the
// record sets of its zone at the record's name/type coordinates. A record has
// no standalone describe API — its identity is (zone, name, type[, set
// identifier]) — so verification needs three outputs (zone_id, fqdn,
// record_type), which is the OutputsVerifier path. Absence of the (name,
// type) pair in the listing is absence of the record (the zone's own NS/SOA
// pair never collides with scenario records). A destroyed ZONE also means the
// record is gone (NoSuchHostedZone = absent).
type route53DnsRecordVerifier struct{}

func (*route53DnsRecordVerifier) IDOutputKey() string { return "fqdn" }

func (*route53DnsRecordVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	return pkgerrors.New("awsroute53dnsrecord verify-exists requires full outputs (zone_id + record_type); use OutputsVerifier path")
}

func (*route53DnsRecordVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	return pkgerrors.New("awsroute53dnsrecord verify-absent requires full outputs (zone_id + record_type); use OutputsVerifier path")
}

func (*route53DnsRecordVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	exists, err := route53RecordExists(ctx, cfg, outputs)
	if err != nil {
		return pkgerrors.Wrap(err, "awsroute53dnsrecord verify-exists")
	}
	if !exists {
		return pkgerrors.Errorf("awsroute53dnsrecord %q not found after deploy", stringOutputMap(outputs, "fqdn"))
	}
	return nil
}

func (*route53DnsRecordVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	exists, err := route53RecordExists(ctx, cfg, outputs)
	if err != nil {
		return pkgerrors.Wrap(err, "awsroute53dnsrecord verify-absent")
	}
	if exists {
		return pkgerrors.Errorf("awsroute53dnsrecord %q still exists after destroy", stringOutputMap(outputs, "fqdn"))
	}
	return nil
}

func route53RecordExists(ctx context.Context, cfg aws.Config, outputs map[string]interface{}) (bool, error) {
	zoneId := stringOutputMap(outputs, "zone_id")
	fqdn := stringOutputMap(outputs, "fqdn")
	recordType := stringOutputMap(outputs, "record_type")
	if zoneId == "" || fqdn == "" || recordType == "" {
		return false, pkgerrors.New("outputs must carry zone_id, fqdn, and record_type -- cannot verify")
	}

	client := route53.NewFromConfig(cfg)
	out, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneId),
		StartRecordName: aws.String(fqdn),
		StartRecordType: route53types.RRType(recordType),
		MaxItems:        aws.Int32(5),
	})
	if err != nil {
		// The record cannot outlive its zone.
		var noZone *route53types.NoSuchHostedZone
		if errors.As(err, &noZone) {
			return false, nil
		}
		return false, err
	}

	// The listing starts AT the requested coordinates when the record exists;
	// compare with Route 53's trailing-dot normalization.
	want := strings.TrimSuffix(strings.ToLower(fqdn), ".")
	for _, rrset := range out.ResourceRecordSets {
		got := strings.TrimSuffix(strings.ToLower(aws.ToString(rrset.Name)), ".")
		if got == want && string(rrset.Type) == recordType {
			return true, nil
		}
	}
	return false, nil
}
