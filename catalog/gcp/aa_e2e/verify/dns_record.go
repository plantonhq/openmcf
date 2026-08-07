package verify

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// dnsRecordVerifier probes a Cloud DNS record set via the dns API. A record
// set is addressed by (zone, fqdn, type) — all three ride the component's
// outputs. Posture assertions confirm the record answers (rrdatas or a
// routing policy present) and the TTL matches the deployed intent.
type dnsRecordVerifier struct{}

func (v *dnsRecordVerifier) IDOutputKey() string { return "fqdn" }

func (v *dnsRecordVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	fqdn := outputs["fqdn"]
	zoneName := outputs["managed_zone"]
	recordType := outputs["record_type"]
	if fqdn == "" || zoneName == "" || recordType == "" {
		return errors.New("fqdn/managed_zone/record_type outputs missing after deploy")
	}

	project := outputs["project_id"]
	if project == "" {
		project = svc.Project
	}

	recordSet, err := svc.DNS.ResourceRecordSets.Get(project, zoneName, fqdn, recordType).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "dns record set %s (%s) not found in zone %s after deploy", fqdn, recordType, zoneName)
	}

	// A record set answers with static rrdatas XOR a routing policy.
	if len(recordSet.Rrdatas) == 0 && recordSet.RoutingPolicy == nil {
		return errors.Errorf("dns record set %s has neither rrdatas nor a routing policy after deploy", fqdn)
	}

	if wantTtl := outputs["ttl_seconds"]; wantTtl != "" {
		if liveTtl := strconv.FormatInt(recordSet.Ttl, 10); liveTtl != wantTtl {
			return errors.Errorf("dns record set %s ttl mismatch: output %s, live %s", fqdn, wantTtl, liveTtl)
		}
	}
	return nil
}

func (v *dnsRecordVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	fqdn := outputs["fqdn"]
	zoneName := outputs["managed_zone"]
	recordType := outputs["record_type"]
	if fqdn == "" || zoneName == "" || recordType == "" {
		return nil
	}

	project := outputs["project_id"]
	if project == "" {
		project = svc.Project
	}

	_, err := svc.DNS.ResourceRecordSets.Get(project, zoneName, fqdn, recordType).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		// The zone itself is usually destroyed alongside the record in the
		// prerequisite teardown; a missing zone equally proves absence.
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing dns record set %s after destroy", fqdn)
	}
	return errors.Errorf("dns record set %s still exists after destroy", fqdn)
}
