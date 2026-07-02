package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// projectIamMemberVerifier confirms an additive (role, member) pair is present
// in (or absent from) the project's IAM policy. The grant has no standalone
// server-side object — the policy itself is the source of truth, so the
// verifier reads the whole policy and looks for the exact pair.
type projectIamMemberVerifier struct{}

func (v *projectIamMemberVerifier) IDOutputKey() string { return "member" }

func (v *projectIamMemberVerifier) grantInPolicy(ctx context.Context, svc *Services, outputs map[string]string) (bool, error) {
	project := outputs["project_id"]
	if project == "" {
		project = svc.Project
	}

	// Version 3 so bindings with IAM conditions are returned as distinct
	// entries rather than collapsed.
	policy, err := svc.Crm.Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{
		Options: &cloudresourcemanager.GetPolicyOptions{RequestedPolicyVersion: 3},
	}).Context(ctx).Do()
	if err != nil {
		return false, errors.Wrapf(err, "failed to read IAM policy of project %s", project)
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

func (v *projectIamMemberVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	found, err := v.grantInPolicy(ctx, svc, outputs)
	if err != nil {
		return err
	}
	if !found {
		return errors.Errorf("grant of %s to %s not present in project IAM policy after deploy",
			outputs["role"], outputs["member"])
	}
	return nil
}

func (v *projectIamMemberVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	found, err := v.grantInPolicy(ctx, svc, outputs)
	if err != nil {
		return err
	}
	if found {
		return errors.Errorf("grant of %s to %s still present in project IAM policy after destroy",
			outputs["role"], outputs["member"])
	}
	return nil
}
