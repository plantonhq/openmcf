package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// vertexAiNotebookVerifier probes a Vertex AI Workbench instance via the
// Notebooks REST API. The pinned google.golang.org/api line has no typed
// Workbench client, so the probe is a plain authenticated GET on the
// instance's documented v2 resource path. Posture assertions confirm the
// platform attribution labels landed (the label-parity proof) and that the
// instance is in an operational state.
type vertexAiNotebookVerifier struct{}

func (v *vertexAiNotebookVerifier) IDOutputKey() string { return "instance_id" }

type vertexAiNotebook struct {
	Name   string            `json:"name"`
	State  string            `json:"state"`
	Labels map[string]string `json:"labels"`
}

func (v *vertexAiNotebookVerifier) get(ctx context.Context, svc *Services, instanceID string) (*vertexAiNotebook, int, error) {
	url := fmt.Sprintf("https://notebooks.googleapis.com/v2/%s", instanceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build workbench instance GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "workbench instance GET request failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to read workbench instance response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Errorf("workbench instance GET %s returned %d: %s", instanceID, resp.StatusCode, string(body))
	}

	instance := &vertexAiNotebook{}
	if err := json.Unmarshal(body, instance); err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to decode workbench instance")
	}
	return instance, resp.StatusCode, nil
}

func (v *vertexAiNotebookVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	instanceID := outputs["instance_id"]
	if instanceID == "" {
		return errors.New("instance_id output missing after deploy")
	}

	instance, _, err := v.get(ctx, svc, instanceID)
	if err != nil {
		return errors.Wrapf(err, "workbench instance %s not found after deploy", instanceID)
	}

	if instance.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("workbench instance %s missing the planton-ai_resource attribution label after deploy", instanceID)
	}

	// A freshly created notebook should be running or still initializing;
	// anything terminal means the deploy left a broken posture.
	switch instance.State {
	case "ACTIVE", "INITIALIZING", "STARTING":
	default:
		return errors.Errorf("workbench instance %s in state %s after deploy (expected ACTIVE, INITIALIZING, or STARTING)", instanceID, instance.State)
	}

	if got := outputs["instance_name"]; got != "" {
		liveName := lastPathSegment(instance.Name)
		if got != liveName {
			return errors.Errorf("workbench instance %s instance_name output %q does not match live name %q", instanceID, got, liveName)
		}
	}
	return nil
}

func (v *vertexAiNotebookVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	instanceID := outputs["instance_id"]
	if instanceID == "" {
		return nil
	}

	_, status, err := v.get(ctx, svc, instanceID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing workbench instance %s after destroy", instanceID)
	}
	return errors.Errorf("workbench instance %s still exists after destroy", instanceID)
}
