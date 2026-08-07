package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// firestoreDatabaseVerifier probes a Firestore database via the Firestore
// Admin API (projects.databases.get on the fully qualified database_id).
// Posture assertions confirm the server-generated UID matches the outputs
// and the database reports a concrete type — proof the deploy created a
// real database, not just a name.
type firestoreDatabaseVerifier struct{}

func (v *firestoreDatabaseVerifier) IDOutputKey() string { return "database_id" }

func (v *firestoreDatabaseVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	databaseId := outputs["database_id"]
	if databaseId == "" {
		return errors.New("database_id output missing after deploy")
	}

	db, err := svc.Firestore.Projects.Databases.Get(databaseId).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "firestore database %s not found after deploy", databaseId)
	}
	if db.Type == "" {
		return errors.Errorf("firestore database %s reports no type — unexpected API state", databaseId)
	}
	if uid := outputs["uid"]; uid != "" && db.Uid != uid {
		return errors.Errorf("firestore database %s uid mismatch: output %q, live %q", databaseId, uid, db.Uid)
	}
	return nil
}

func (v *firestoreDatabaseVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	databaseId := outputs["database_id"]
	if databaseId == "" {
		return nil
	}

	_, err := svc.Firestore.Projects.Databases.Get(databaseId).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing firestore database %s after destroy", databaseId)
	}
	return errors.Errorf("firestore database %s still exists after destroy", databaseId)
}
