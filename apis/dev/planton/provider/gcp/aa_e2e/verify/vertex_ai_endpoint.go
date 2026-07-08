package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// vertexAiEndpointVerifier probes a Vertex AI endpoint via the AI Platform
// REST API. The pinned google.golang.org/api line has no typed Vertex client
// for endpoints at the harness's granularity, so the probe is a plain
// authenticated GET on the endpoint's documented resource path (the
// memorystore precedent). Posture assertions confirm the platform attribution
// labels landed (the label-parity proof) and that the endpoint_name output
// matches the live numeric ID — the cross-engine determinism contract both
// modules share.
type vertexAiEndpointVerifier struct{}

func (v *vertexAiEndpointVerifier) IDOutputKey() string { return "endpoint_id" }

type vertexAiEndpoint struct {
	Name                 string            `json:"name"`
	DisplayName          string            `json:"displayName"`
	Labels               map[string]string `json:"labels"`
	DedicatedEndpointDns string            `json:"dedicatedEndpointDns"`
}

func (v *vertexAiEndpointVerifier) get(ctx context.Context, svc *Services, endpointID string) (*vertexAiEndpoint, int, error) {
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/%s",
		regionFromVertexResource(endpointID), endpointID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build vertex ai endpoint GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "vertex ai endpoint GET request failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to read vertex ai endpoint response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Errorf("vertex ai endpoint GET %s returned %d: %s", endpointID, resp.StatusCode, string(body))
	}

	endpoint := &vertexAiEndpoint{}
	if err := json.Unmarshal(body, endpoint); err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to decode vertex ai endpoint")
	}
	return endpoint, resp.StatusCode, nil
}

func (v *vertexAiEndpointVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	endpointID := outputs["endpoint_id"]
	if endpointID == "" {
		return errors.New("endpoint_id output missing after deploy")
	}

	endpoint, _, err := v.get(ctx, svc, endpointID)
	if err != nil {
		return errors.Wrapf(err, "vertex ai endpoint %s not found after deploy", endpointID)
	}

	if endpoint.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("vertex ai endpoint %s missing the planton-ai_resource attribution label after deploy", endpointID)
	}

	// The numeric endpoint_name output must match the live resource ID —
	// proof that both engines derived (or accepted) the same cloud-side name.
	if got := outputs["endpoint_name"]; got != "" {
		liveName := lastPathSegment(endpoint.Name)
		if got != liveName {
			return errors.Errorf("vertex ai endpoint %s endpoint_name output %q does not match live name %q", endpointID, got, liveName)
		}
	}

	if dns := outputs["dedicated_endpoint_dns"]; dns != "" && endpoint.DedicatedEndpointDns != dns {
		return errors.Errorf("vertex ai endpoint %s dedicated_endpoint_dns output %q does not match live %q", endpointID, dns, endpoint.DedicatedEndpointDns)
	}
	return nil
}

func (v *vertexAiEndpointVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	endpointID := outputs["endpoint_id"]
	if endpointID == "" {
		return nil
	}

	_, status, err := v.get(ctx, svc, endpointID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing vertex ai endpoint %s after destroy", endpointID)
	}
	return errors.Errorf("vertex ai endpoint %s still exists after destroy", endpointID)
}

// regionFromVertexResource extracts the region from a fully qualified Vertex
// resource path (projects/{p}/locations/{region}/...) — Vertex AI endpoints
// are served from regional API hosts.
func regionFromVertexResource(name string) string {
	const marker = "/locations/"
	i := strings.Index(name, marker)
	if i < 0 {
		return "us-central1"
	}
	rest := name[i+len(marker):]
	if j := strings.IndexByte(rest, '/'); j >= 0 {
		return rest[:j]
	}
	return rest
}

// lastPathSegment returns the substring after the final '/' — the short
// resource name at the end of a fully qualified resource path.
func lastPathSegment(name string) string {
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		return name[i+1:]
	}
	return name
}
