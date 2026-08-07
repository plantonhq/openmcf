package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// serviceAccountIamMemberVerifier confirms an additive (role, member) pair is
// present in (or absent from) the SERVICE ACCOUNT's IAM policy — the policy on
// the account as a resource, where impersonation grants live. The grant has no
// standalone server-side object — the policy itself is the source of truth,
// so the verifier reads the whole policy and looks for the exact pair.
type serviceAccountIamMemberVerifier struct{}

func (v *serviceAccountIamMemberVerifier) IDOutputKey() string { return "member" }

func (v *serviceAccountIamMemberVerifier) grantInPolicy(ctx context.Context, svc *Services, outputs map[string]string) (bool, error) {
	serviceAccountId := outputs["service_account_id"]
	if serviceAccountId == "" {
		return false, errors.New("service_account_id output missing")
	}

	// Version 3 so bindings with IAM conditions are returned as distinct
	// entries rather than collapsed.
	policy, err := svc.Iam.Projects.ServiceAccounts.GetIamPolicy(serviceAccountId).
		OptionsRequestedPolicyVersion(3).Context(ctx).Do()
	if err != nil {
		return false, err
	}

	for _, binding := range policy.Bindings {
		if binding.Role != outputs["role"] {
			continue
		}
		for _, member := range binding.Members {
			// IAM treats most member emails case-insensitively; compare the
			// same way so server-side casing normalization never fails a probe.
			if strings.EqualFold(member, outputs["member"]) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (v *serviceAccountIamMemberVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	found, err := v.grantInPolicy(ctx, svc, outputs)
	if err != nil {
		return errors.Wrapf(err, "failed to read IAM policy of service account %s", outputs["service_account_id"])
	}
	if !found {
		return errors.Errorf("grant of %s to %s not present in service account IAM policy after deploy",
			outputs["role"], outputs["member"])
	}
	return nil
}

func (v *serviceAccountIamMemberVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	found, err := v.grantInPolicy(ctx, svc, outputs)
	if err != nil {
		// A vanished service account carries no grants — that IS the absent
		// state (the account prerequisite may already be torn down).
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "failed to read IAM policy of service account %s", outputs["service_account_id"])
	}
	if found {
		return errors.Errorf("grant of %s to %s still present in service account IAM policy after destroy",
			outputs["role"], outputs["member"])
	}
	return nil
}
