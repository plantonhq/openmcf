package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// kmsKeyRingVerifier probes a Cloud KMS key ring by its fully qualified
// resource path (projects/{p}/locations/{l}/keyRings/{name}) — the key_ring_id
// output is exactly the handle a GcpKmsKey's key_ring_id reference consumes,
// so verifying with it doubles as proof the composition handle is honest.
type kmsKeyRingVerifier struct{}

func (v *kmsKeyRingVerifier) IDOutputKey() string { return "key_ring_id" }

func (v *kmsKeyRingVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	keyRingID := outputs["key_ring_id"]
	if keyRingID == "" {
		return errors.New("key_ring_id output missing after deploy")
	}

	ring, err := svc.CloudKms.Projects.Locations.KeyRings.Get(keyRingID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "kms key ring %s not found after deploy", keyRingID)
	}
	if ring.Name != keyRingID {
		return errors.Errorf("kms key ring resolved to %q, want %q", ring.Name, keyRingID)
	}
	return nil
}

func (v *kmsKeyRingVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	keyRingID := outputs["key_ring_id"]
	if keyRingID == "" {
		return nil
	}

	// Key rings have no delete API: GCP keeps the (free, inert) ring forever
	// and the engines remove it from state only. The ring persisting after
	// destroy IS the destroyed state for this resource class — the probe
	// exists to catch the unexpected (a permission regression or an API
	// error class change), not to demand absence.
	_, err := svc.CloudKms.Projects.Locations.KeyRings.Get(keyRingID).Context(ctx).Do()
	if err == nil {
		return nil
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		// Absent would only happen if the parent project vanished — still a
		// valid "gone" answer.
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing kms key ring %s after destroy", keyRingID)
}
