package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// logBucketVerifier probes a Cloud Logging bucket by its full resource
// name (the kind's composition output) and confirms it is ACTIVE. Live
// lanes exercise the project scope only — folder/org buckets need org
// credentials (the recorded deferral) and the test tenant has no billing
// fixture — but the v2 Buckets service addresses any scope's full name,
// so the probe needs no scope switch. The Logging API holds deleted
// buckets in a DELETE_REQUESTED lifecycle state for 7 days, so
// VerifyAbsent accepts pending deletion as gone-for-management.
type logBucketVerifier struct{}

func (v *logBucketVerifier) IDOutputKey() string { return "bucket_name" }

func (v *logBucketVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["bucket_name"]
	bucket, err := svc.Logging.Projects.Locations.Buckets.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "log bucket %s not found after deploy", name)
	}
	if bucket.LifecycleState != "ACTIVE" {
		return errors.Errorf("log bucket %s reports lifecycle state %s, want ACTIVE", name, bucket.LifecycleState)
	}
	return nil
}

func (v *logBucketVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["bucket_name"]
	bucket, err := svc.Logging.Projects.Locations.Buckets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing log bucket %s after destroy", name)
	}
	// Pending deletion counts as gone: the API keeps the bucket visible in
	// DELETE_REQUESTED for 7 days by design.
	if bucket.LifecycleState == "DELETE_REQUESTED" {
		return nil
	}
	return errors.Errorf("log bucket %s still exists after destroy (lifecycle state %s)", name, bucket.LifecycleState)
}
