package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// workflowVerifier probes a Cloud Workflows workflow and confirms it
// exists with deployed source. The workflow_id output is the full
// resource name (projects/{p}/locations/{region}/workflows/{name}) —
// exactly what the Workflows GET takes.
type workflowVerifier struct{}

func (v *workflowVerifier) IDOutputKey() string { return "workflow_id" }

func (v *workflowVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["workflow_id"]
	workflow, err := svc.Workflows.Projects.Locations.Workflows.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "workflow %s not found after deploy", name)
	}
	if workflow.RevisionId == "" {
		return errors.Errorf("workflow %s reports no deployed revision", name)
	}
	return nil
}

func (v *workflowVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["workflow_id"]
	_, err := svc.Workflows.Projects.Locations.Workflows.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing workflow %s after destroy", name)
	}
	return errors.Errorf("workflow %s still exists after destroy", name)
}
