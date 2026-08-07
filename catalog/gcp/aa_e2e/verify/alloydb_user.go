package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

// alloydbUserVerifier probes an AlloyDB database user via the alloydb API.
// Posture assertions confirm the live user resource name matches the user_id
// output (the database role identifier).
type alloydbUserVerifier struct{}

func (v *alloydbUserVerifier) IDOutputKey() string { return "name" }

func (v *alloydbUserVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return errors.New("name output missing after deploy")
	}

	user, err := svc.AlloyDB.Projects.Locations.Clusters.Users.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "alloydb user %s not found after deploy", name)
	}

	wantUserID := outputs["user_id"]
	if wantUserID == "" {
		return nil
	}
	gotUserID := alloyDBUserPathTail(user.Name)
	if gotUserID != wantUserID {
		return errors.Errorf("alloydb user %s user_id mismatch: output %q, live %q", name, wantUserID, gotUserID)
	}
	return nil
}

func (v *alloydbUserVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return nil
	}

	_, err := svc.AlloyDB.Projects.Locations.Clusters.Users.Get(name).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("alloydb user %s still exists after destroy", name)
	}
	if isAlloyDBNotFound(err) {
		return nil
	}
	if clusterID := alloyDBClusterFromChildPath(name, "/users/"); clusterID != "" {
		if _, clusterErr := svc.AlloyDB.Projects.Locations.Clusters.Get(clusterID).Context(ctx).Do(); clusterErr != nil && isAlloyDBNotFound(clusterErr) {
			return nil
		}
	}
	return errors.Wrapf(err, "unexpected error probing alloydb user %s after destroy", name)
}

func alloyDBUserPathTail(userName string) string {
	const marker = "/users/"
	if idx := strings.LastIndex(userName, marker); idx >= 0 {
		return userName[idx+len(marker):]
	}
	return userName
}
