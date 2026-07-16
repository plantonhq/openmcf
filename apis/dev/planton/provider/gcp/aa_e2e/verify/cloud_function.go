package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudFunctionVerifier probes a Cloud Functions (Gen 2) function via the
// Cloud Functions API v2. Posture assertions confirm the function reconciled
// to ACTIVE in the Gen 2 environment and that the underlying Cloud Run
// service produced a serving URI.
type cloudFunctionVerifier struct{}

func (v *cloudFunctionVerifier) IDOutputKey() string { return "function_id" }

// functionPath returns the fully qualified resource name. The function_id
// output IS the projects/*/locations/*/functions/* path the API expects.
func (v *cloudFunctionVerifier) functionPath(outputs map[string]string) (string, error) {
	path := outputs["function_id"]
	if path == "" {
		return "", errors.New("function_id output missing")
	}
	return path, nil
}

func (v *cloudFunctionVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.functionPath(outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	function, err := svc.Functions.Projects.Locations.Functions.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "cloud function %s not found after deploy", path)
	}

	if function.State != "ACTIVE" {
		return errors.Errorf("cloud function %s state is %q, want ACTIVE", path, function.State)
	}
	if function.Environment != "GEN_2" {
		return errors.Errorf("cloud function %s environment is %q, want GEN_2", path, function.Environment)
	}
	if function.ServiceConfig == nil || function.ServiceConfig.Uri == "" {
		return errors.Errorf("cloud function %s has no serving uri after deploy", path)
	}
	return nil
}

func (v *cloudFunctionVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.functionPath(outputs)
	if err != nil {
		return nil
	}

	_, err = svc.Functions.Projects.Locations.Functions.Get(path).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("cloud function %s still exists after destroy", path)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing cloud function %s after destroy", path)
}
