package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// kmsKeyVerifier probes a Cloud KMS crypto key by its fully qualified
// resource path (projects/{p}/locations/{l}/keyRings/{r}/cryptoKeys/{name})
// — the key_id output is exactly the CMEK handle every downstream consumer
// passes to its kms_key_name-style field, so verifying with it doubles as
// proof the composition handle is honest. Posture assertions include the
// platform attribution label (both engines must stamp the identical
// planton-ai_* set — a permanently guarded live regression for the
// label-parity defect class) and, for keys that mint an initial version,
// an ENABLED primary.
type kmsKeyVerifier struct{}

func (v *kmsKeyVerifier) IDOutputKey() string { return "key_id" }

func (v *kmsKeyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	keyID := outputs["key_id"]
	if keyID == "" {
		return errors.New("key_id output missing after deploy")
	}

	key, err := svc.CloudKms.Projects.Locations.KeyRings.CryptoKeys.Get(keyID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "kms crypto key %s not found after deploy", keyID)
	}

	if key.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("kms crypto key %s missing the planton-ai_resource attribution label after deploy (labels: %v)", keyID, key.Labels)
	}

	// The primary outputs are populated only for ENCRYPT_DECRYPT keys with
	// at least one version; when the deploy exported them, the live key
	// must agree — an ENABLED primary is the proof the key can encrypt.
	if wantPrimary := outputs["primary_version_name"]; wantPrimary != "" {
		if key.Primary == nil {
			return errors.Errorf("kms crypto key %s exported primary version %s but the live key has no primary", keyID, wantPrimary)
		}
		if key.Primary.Name != wantPrimary {
			return errors.Errorf("kms crypto key %s primary is %q, outputs claim %q", keyID, key.Primary.Name, wantPrimary)
		}
		if key.Primary.State != "ENABLED" {
			return errors.Errorf("kms crypto key %s primary version is %s (want ENABLED) after deploy", keyID, key.Primary.State)
		}
	}
	return nil
}

func (v *kmsKeyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	keyID := outputs["key_id"]
	if keyID == "" {
		return nil
	}

	// Crypto keys have no delete API: destroy schedules every version for
	// destruction and disables rotation, and the key object persists in the
	// ring forever. The honest "destroyed" assertion for this resource
	// class is therefore that no version remains usable and the rotation
	// clock is off — not that the key is gone.
	key, err := svc.CloudKms.Projects.Locations.KeyRings.CryptoKeys.Get(keyID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing kms crypto key %s after destroy", keyID)
	}

	if key.RotationPeriod != "" {
		return errors.Errorf("kms crypto key %s still has rotation_period %q after destroy (destroy must disable rotation)", keyID, key.RotationPeriod)
	}

	versions, err := svc.CloudKms.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.List(keyID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "failed to list versions of kms crypto key %s after destroy", keyID)
	}
	for _, version := range versions.CryptoKeyVersions {
		if version.State != "DESTROYED" && version.State != "DESTROY_SCHEDULED" {
			return errors.Errorf("kms crypto key %s version %s is still %s after destroy (want DESTROYED or DESTROY_SCHEDULED)", keyID, version.Name, version.State)
		}
	}
	return nil
}
