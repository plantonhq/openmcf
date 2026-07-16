package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// gkeWorkloadIdentityBindingVerifier probes the target service account's IAM
// policy for the roles/iam.workloadIdentityUser grant to the constructed
// workload-identity principal. An IAM grant has no id of its own — existence
// IS membership in the policy — so both phases read the policy and look for
// the (role, member) pair.
type gkeWorkloadIdentityBindingVerifier struct{}

const workloadIdentityUserRole = "roles/iam.workloadIdentityUser"

func (v *gkeWorkloadIdentityBindingVerifier) IDOutputKey() string { return "member" }

func (v *gkeWorkloadIdentityBindingVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	member := outputs["member"]
	email := outputs["service_account_email"]
	if member == "" || email == "" {
		return errors.New("member/service_account_email outputs missing after deploy")
	}

	found, err := workloadIdentityGrantPresent(ctx, svc, email, member)
	if err != nil {
		return errors.Wrapf(err, "failed to read IAM policy of service account %s after deploy", email)
	}
	if !found {
		return errors.Errorf("service account %s has no %s grant for %s after deploy", email, workloadIdentityUserRole, member)
	}
	return nil
}

func (v *gkeWorkloadIdentityBindingVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	member := outputs["member"]
	email := outputs["service_account_email"]
	if member == "" || email == "" {
		return nil
	}

	found, err := workloadIdentityGrantPresent(ctx, svc, email, member)
	if err != nil {
		var apiErr *googleapi.Error
		// The service account itself is destroyed alongside the grant in the
		// prerequisite teardown; a missing SA equally proves absence.
		if errors.As(err, &apiErr) && (apiErr.Code == 404 || apiErr.Code == 403) {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing IAM policy of service account %s after destroy", email)
	}
	if found {
		return errors.Errorf("service account %s still carries the %s grant for %s after destroy", email, workloadIdentityUserRole, member)
	}
	return nil
}

func workloadIdentityGrantPresent(ctx context.Context, svc *Services, email, member string) (bool, error) {
	// The "-" project wildcard infers the SA's project from the email — the
	// same fully-qualified form the modules pass the provider.
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", email)
	policy, err := svc.Iam.Projects.ServiceAccounts.GetIamPolicy(resource).Context(ctx).Do()
	if err != nil {
		return false, err
	}
	for _, binding := range policy.Bindings {
		if binding.Role != workloadIdentityUserRole {
			continue
		}
		for _, m := range binding.Members {
			if m == member {
				return true, nil
			}
		}
	}
	return false, nil
}
