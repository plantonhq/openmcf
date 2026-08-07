package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// projectVerifier probes a Google Cloud project via the cloudresourcemanager
// API. Destroyed projects enter a 30-day pending-deletion window rather than
// vanishing, so the absence contract accepts DELETE_REQUESTED as destroyed —
// the honest posture for GCP's soft-deleted resource class.
type projectVerifier struct{}

func (v *projectVerifier) IDOutputKey() string { return "project_id" }

func (v *projectVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	projectID := outputs["project_id"]
	if projectID == "" {
		return errors.New("project_id output missing after deploy")
	}

	project, err := svc.Crm.Projects.Get(projectID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "project %s not found after deploy", projectID)
	}
	if project.LifecycleState != "ACTIVE" {
		return errors.Errorf("project %s lifecycle state is %s after deploy, want ACTIVE", projectID, project.LifecycleState)
	}

	// Live label-parity guard: both engines must have applied the platform
	// attribution labels identically.
	if project.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("project %s is missing the planton-ai_resource attribution label after deploy (labels: %v)",
			projectID, project.Labels)
	}
	return nil
}

func (v *projectVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	projectID := outputs["project_id"]
	if projectID == "" {
		return nil
	}

	project, err := svc.Crm.Projects.Get(projectID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && (apiErr.Code == 404 || apiErr.Code == 403) {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing project %s after destroy", projectID)
	}
	// Deleted projects linger in DELETE_REQUESTED for the 30-day recovery
	// window — that IS the destroyed posture for this resource class.
	if project.LifecycleState == "DELETE_REQUESTED" {
		return nil
	}
	return errors.Errorf("project %s still exists after destroy (lifecycle state %s)", projectID, project.LifecycleState)
}
