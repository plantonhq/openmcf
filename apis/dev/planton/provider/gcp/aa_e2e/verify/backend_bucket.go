package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// backendBucketVerifier probes a Compute Engine backend bucket by name via
// the compute API, additionally confirming it points at the expected origin
// bucket (the one FK-resolved from the GcsBucket prerequisite).
type backendBucketVerifier struct{}

func (v *backendBucketVerifier) IDOutputKey() string { return "self_link" }

func (v *backendBucketVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["backend_bucket_name"]
	backendBucket, err := svc.Compute.BackendBuckets.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "backend bucket %s not found after deploy", name)
	}
	// The origin wiring is the point of the composition — confirm the deployed
	// backend bucket references the origin the outputs claim.
	if want := outputs["bucket_name"]; want != "" && backendBucket.BucketName != want {
		return errors.Errorf("backend bucket %s points at origin %q, expected %q",
			name, backendBucket.BucketName, want)
	}
	return nil
}

func (v *backendBucketVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["backend_bucket_name"]
	_, err := svc.Compute.BackendBuckets.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing backend bucket %s after destroy", name)
	}
	return errors.Errorf("backend bucket %s still exists after destroy", name)
}
