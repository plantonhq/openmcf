package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// serviceAccountVerifier probes a service account by email via the IAM API.
type serviceAccountVerifier struct{}

func (v *serviceAccountVerifier) IDOutputKey() string { return "email" }

// resourceName addresses the account by email under the "-" wildcard project,
// which the IAM API resolves to the account's home project.
func (v *serviceAccountVerifier) resourceName(outputs map[string]string) string {
	return fmt.Sprintf("projects/-/serviceAccounts/%s", outputs["email"])
}

func (v *serviceAccountVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	account, err := svc.Iam.Projects.ServiceAccounts.Get(v.resourceName(outputs)).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "service account %s not found after deploy", outputs["email"])
	}
	if account.Disabled && outputs["disabled"] != "true" {
		return errors.Errorf("service account %s exists but is unexpectedly disabled", outputs["email"])
	}
	return nil
}

func (v *serviceAccountVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	_, err := svc.Iam.Projects.ServiceAccounts.Get(v.resourceName(outputs)).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("service account %s still exists after destroy", outputs["email"])
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && (apiErr.Code == 404 || apiErr.Code == 403) {
		// 404 is the definitive gone signal; the IAM API returns 403 for
		// recently-deleted accounts in some propagation windows.
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing service account %s after destroy", outputs["email"])
}
