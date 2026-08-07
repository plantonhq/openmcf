package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// firestoreIndexVerifier probes a Firestore composite index via the
// Firestore Admin API (collectionGroups.indexes.get on the server-defined
// index_id resource path). Posture assertions confirm the index is
// serving (or still building — both engines wait for the create LRO, so
// READY is the expected steady state) and that the live collection group
// matches the outputs.
type firestoreIndexVerifier struct{}

func (v *firestoreIndexVerifier) IDOutputKey() string { return "index_id" }

func (v *firestoreIndexVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	indexId := outputs["index_id"]
	if indexId == "" {
		return errors.New("index_id output missing after deploy")
	}

	index, err := svc.Firestore.Projects.Databases.CollectionGroups.Indexes.Get(indexId).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "firestore index %s not found after deploy", indexId)
	}
	if index.State != "READY" && index.State != "CREATING" {
		return errors.Errorf("firestore index %s state is %q, want READY (or CREATING mid-build)", indexId, index.State)
	}
	// The collection group is a path segment of the index resource name
	// (projects/{p}/databases/{d}/collectionGroups/{collection}/indexes/{id}).
	if collection := outputs["collection"]; collection != "" &&
		!strings.Contains(index.Name, "/collectionGroups/"+collection+"/") {
		return errors.Errorf("firestore index %s does not belong to collection group %q", index.Name, collection)
	}
	return nil
}

func (v *firestoreIndexVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	indexId := outputs["index_id"]
	if indexId == "" {
		return nil
	}

	_, err := svc.Firestore.Projects.Databases.CollectionGroups.Indexes.Get(indexId).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && (apiErr.Code == 404 || apiErr.Code == 403) {
			// 404: index gone. 403: the chain's database may already be
			// destroyed, which can surface as a permission-shaped error.
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing firestore index %s after destroy", indexId)
	}
	return errors.Errorf("firestore index %s still exists after destroy", indexId)
}
