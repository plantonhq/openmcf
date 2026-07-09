package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// certManagerDnsAuthorizationVerifier probes a Certificate Manager DNS
// authorization via the certificatemanager API. The authorization_id output
// is the fully-qualified resource name, so the probe consumes it directly.
// Posture assertions confirm the authorized domain and that the validation
// record GCP computed matches the exported composition fields.
type certManagerDnsAuthorizationVerifier struct{}

func (v *certManagerDnsAuthorizationVerifier) IDOutputKey() string { return "authorization_id" }

func (v *certManagerDnsAuthorizationVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	authorizationID := outputs["authorization_id"]
	if authorizationID == "" {
		return errors.New("authorization_id output missing after deploy")
	}

	authorization, err := svc.CertificateManager.Projects.Locations.DnsAuthorizations.Get(authorizationID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "dns authorization %s not found after deploy", authorizationID)
	}

	if wantDomain := outputs["domain"]; wantDomain != "" && authorization.Domain != wantDomain {
		return errors.Errorf("dns authorization %s domain mismatch: output %q, live %q", authorizationID, wantDomain, authorization.Domain)
	}

	if authorization.DnsResourceRecord == nil {
		return errors.Errorf("dns authorization %s has no validation record after deploy", authorizationID)
	}
	if wantRecordName := outputs["dns_record_name"]; wantRecordName != "" && authorization.DnsResourceRecord.Name != wantRecordName {
		return errors.Errorf("dns authorization %s validation record name mismatch: output %q, live %q",
			authorizationID, wantRecordName, authorization.DnsResourceRecord.Name)
	}
	return nil
}

func (v *certManagerDnsAuthorizationVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	authorizationID := outputs["authorization_id"]
	if authorizationID == "" {
		return nil
	}

	_, err := svc.CertificateManager.Projects.Locations.DnsAuthorizations.Get(authorizationID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing dns authorization %s after destroy", authorizationID)
	}
	return errors.Errorf("dns authorization %s still exists after destroy", authorizationID)
}
