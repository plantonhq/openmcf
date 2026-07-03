package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// gcsBucketVerifier probes a Cloud Storage bucket by its globally-unique name
// via the storage API. Registered so a bucket can serve as a verified E2E
// prerequisite for origin-consuming kinds (e.g. backend buckets).
type gcsBucketVerifier struct{}

func (v *gcsBucketVerifier) IDOutputKey() string { return "bucket_id" }

func (v *gcsBucketVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["bucket_id"]
	if _, err := svc.Storage.Buckets.Get(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "gcs bucket %s not found after deploy", name)
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
