package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// vertexAiIndexEndpointVerifier probes a Vertex AI Vector Search index
// endpoint via the AI Platform REST API — a plain authenticated GET on the
// endpoint's documented resource path (no typed Vertex client on the pinned
// google.golang.org/api line). Posture assertions confirm the platform
// attribution labels landed and that the index_endpoint_name output matches
// the live numeric ID — the cross-engine determinism contract both modules
// share. The endpoint GET also carries deployedIndexes[], which the deployed
// index verifier reuses.
type vertexAiIndexEndpointVerifier struct{}

func (v *vertexAiIndexEndpointVerifier) IDOutputKey() string { return "index_endpoint_id" }

type vertexAiIndexEndpoint struct {
	Name                     string                        `json:"name"`
	DisplayName              string                        `json:"displayName"`
	Labels                   map[string]string             `json:"labels"`
	PublicEndpointDomainName string                        `json:"publicEndpointDomainName"`
	DeployedIndexes          []vertexAiIndexEndpointDeploy `json:"deployedIndexes"`
}

// vertexAiIndexEndpointDeploy is one deployedIndexes[] entry — the live
// representation of a GcpVertexAiDeployedIndex.
type vertexAiIndexEndpointDeploy struct {
	ID    string `json:"id"`
	Index string `json:"index"`
}

func getVertexAiIndexEndpoint(ctx context.Context, svc *Services, endpointID string) (*vertexAiIndexEndpoint, int, error) {
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/%s",
		regionFromVertexResource(endpointID), endpointID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build vertex ai index endpoint GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "vertex ai index endpoint GET request failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to read vertex ai index endpoint response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Errorf("vertex ai index endpoint GET %s returned %d: %s", endpointID, resp.StatusCode, string(body))
	}

	endpoint := &vertexAiIndexEndpoint{}
	if err := json.Unmarshal(body, endpoint); err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to decode vertex ai index endpoint")
	}
	return endpoint, resp.StatusCode, nil
}

func (v *vertexAiIndexEndpointVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	endpointID := outputs["index_endpoint_id"]
	if endpointID == "" {
		return errors.New("index_endpoint_id output missing after deploy")
	}

	endpoint, _, err := getVertexAiIndexEndpoint(ctx, svc, endpointID)
	if err != nil {
		return errors.Wrapf(err, "vertex ai index endpoint %s not found after deploy", endpointID)
	}

	if endpoint.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("vertex ai index endpoint %s missing the planton-ai_resource attribution label after deploy", endpointID)
	}

	// The numeric index_endpoint_name output must match the live resource
	// ID — proof that both engines derived the same cloud-side name.
	if got := outputs["index_endpoint_name"]; got != "" {
		liveName := lastPathSegment(endpoint.Name)
		if got != liveName {
			return errors.Errorf("vertex ai index endpoint %s index_endpoint_name output %q does not match live name %q", endpointID, got, liveName)
		}
	}

	if domain := outputs["public_endpoint_domain_name"]; domain != "" && endpoint.PublicEndpointDomainName != domain {
		return errors.Errorf("vertex ai index endpoint %s public_endpoint_domain_name output %q does not match live %q", endpointID, domain, endpoint.PublicEndpointDomainName)
	}
	return nil
}

func (v *vertexAiIndexEndpointVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	endpointID := outputs["index_endpoint_id"]
	if endpointID == "" {
		return nil
	}

	_, status, err := getVertexAiIndexEndpoint(ctx, svc, endpointID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing vertex ai index endpoint %s after destroy", endpointID)
	}
	return errors.Errorf("vertex ai index endpoint %s still exists after destroy", endpointID)
}
