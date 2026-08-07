package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// certManagerCertVerifier probes a Certificate Manager certificate via the
// certificatemanager API. The certificate_id output is the fully-qualified
// resource name, so the probe consumes it directly.
//
// The posture contract for managed certificates is CREATION posture, never
// ACTIVE: a managed certificate only turns ACTIVE once public DNS validation
// completes on a domain the test run owns, so waiting on ACTIVE would hang
// the run. The verifier asserts the certificate exists, carries the platform
// attribution label (the live label-parity guard), and — for managed
// certificates — reports a non-failed provisioning state.
type certManagerCertVerifier struct{}

func (v *certManagerCertVerifier) IDOutputKey() string { return "certificate_id" }

func (v *certManagerCertVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	certificateID := outputs["certificate_id"]
	if certificateID == "" {
		return errors.New("certificate_id output missing after deploy")
	}

	certificate, err := svc.CertificateManager.Projects.Locations.Certificates.Get(certificateID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "certificate %s not found after deploy", certificateID)
	}

	// Live label-parity guard: both engines must have applied the platform
	// attribution labels identically.
	if certificate.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("certificate %s is missing the planton-ai_resource attribution label after deploy (labels: %v)",
			certificateID, certificate.Labels)
	}

	// Managed certificates report their provisioning state; a FAILED state
	// means the arm deployed misconfigured. PROVISIONING is the expected
	// steady state for a domain the run does not own.
	if wantState := outputs["managed_state"]; wantState != "" {
		if certificate.Managed == nil {
			return errors.Errorf("certificate %s exported managed_state %q but has no managed block live", certificateID, wantState)
		}
		if strings.EqualFold(certificate.Managed.State, "FAILED") {
			return errors.Errorf("certificate %s managed provisioning FAILED after deploy", certificateID)
		}
	}
	return nil
}

func (v *certManagerCertVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	certificateID := outputs["certificate_id"]
	if certificateID == "" {
		return nil
	}

	_, err := svc.CertificateManager.Projects.Locations.Certificates.Get(certificateID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing certificate %s after destroy", certificateID)
	}
	return errors.Errorf("certificate %s still exists after destroy", certificateID)
}
