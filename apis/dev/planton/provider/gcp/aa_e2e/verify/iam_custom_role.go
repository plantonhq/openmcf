package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// iamCustomRoleVerifier probes a project custom role by its fully-qualified
// name (projects/<project>/roles/<role_id>) via the IAM API.
type iamCustomRoleVerifier struct{}

func (v *iamCustomRoleVerifier) IDOutputKey() string { return "name" }

func (v *iamCustomRoleVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	role, err := svc.Iam.Projects.Roles.Get(outputs["name"]).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "custom role %s not found after deploy", outputs["name"])
	}
	if role.Deleted {
		return errors.Errorf("custom role %s exists but is soft-deleted after deploy", outputs["name"])
	}
	return nil
}

func (v *iamCustomRoleVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	role, err := svc.Iam.Projects.Roles.Get(outputs["name"]).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing custom role %s after destroy", outputs["name"])
	}
	// GCP soft-deletes custom roles: for up to 14 days the role remains
	// readable with deleted=true while rejecting grants. That IS the destroyed
	// state — only a live (deleted=false) role means the destroy failed.
	if !role.Deleted {
		return errors.Errorf("custom role %s still live (not even soft-deleted) after destroy", outputs["name"])
	}
	return nil
}
