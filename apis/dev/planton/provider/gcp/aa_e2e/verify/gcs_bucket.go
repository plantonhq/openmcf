package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// gcsBucketVerifier probes a Cloud Storage bucket by its globally-unique
// name via the storage API. Posture assertions confirm the platform
// attribution labels landed (the cross-engine label-parity canary — the
// Terraform module historically stamped a different label set than
// Pulumi, so this is a permanently guarded regression) and that the
// location output matches the live bucket. Also registered so a bucket
// can serve as a verified E2E prerequisite for origin-consuming kinds
// (e.g. backend buckets).
type gcsBucketVerifier struct{}

func (v *gcsBucketVerifier) IDOutputKey() string { return "bucket_id" }

func (v *gcsBucketVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["bucket_id"]
	bucket, err := svc.Storage.Buckets.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "gcs bucket %s not found after deploy", name)
	}

	// The platform attribution labels are the cross-engine parity canary:
	// a missing set means one engine stamped labels and the other did not.
	if bucket.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("gcs bucket %s missing the planton-ai_resource attribution label after deploy (labels: %v)", name, bucket.Labels)
	}

	// GCS reports location upper-cased; the location output must be the
	// live value on both engines.
	if location := outputs["location"]; location != "" && location != bucket.Location {
		return errors.Errorf("gcs bucket %s location output %q does not match live location %q", name, location, bucket.Location)
	}

	return nil
}

func (v *gcsBucketVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["bucket_id"]
	_, err := svc.Storage.Buckets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing gcs bucket %s after destroy", name)
	}
	return errors.Errorf("gcs bucket %s still exists after destroy", name)
}
