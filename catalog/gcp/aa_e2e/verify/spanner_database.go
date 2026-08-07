package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// spannerDatabaseVerifier probes a Cloud Spanner database via the spanner
// admin API. Posture assertions confirm the database reached READY and its
// live name matches the database_name output.
type spannerDatabaseVerifier struct{}

func (v *spannerDatabaseVerifier) IDOutputKey() string { return "database_id" }

func (v *spannerDatabaseVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	databaseID := outputs["database_id"]
	if databaseID == "" {
		return errors.New("database_id output missing after deploy")
	}

	database, err := svc.Spanner.Projects.Instances.Databases.Get(databaseID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "spanner database %s not found after deploy", databaseID)
	}

	if database.State != "READY" {
		return errors.Errorf("spanner database %s state is %q, want READY", databaseID, database.State)
	}

	if wantName := outputs["database_name"]; wantName != "" && pathTail(database.Name) != wantName {
		return errors.Errorf("spanner database %s name mismatch: output %q, live %q",
			databaseID, wantName, pathTail(database.Name))
	}
	return nil
}

func (v *spannerDatabaseVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	databaseID := outputs["database_id"]
	if databaseID == "" {
		return nil
	}

	_, err := svc.Spanner.Projects.Instances.Databases.Get(databaseID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("spanner database %s still exists after destroy", databaseID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	// When the whole chain (instance included) was destroyed, the API answers
	// with a not-found for the parent instance instead of the database.
	if apiErr != nil && apiErr.Code == 400 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing spanner database %s after destroy", databaseID)
}
