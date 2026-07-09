package verify

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// vertexAiDeployedIndexVerifier probes a Vertex AI deployed index. A
// DeployedIndex is not a first-class GET-able resource — it lives inside its
// index endpoint — so the probe fetches the endpoint (via the index endpoint
// verifier's shared GET) and locates the deployment's entry in
// deployedIndexes[] by its user-chosen ID. The index_endpoint output carries
// the parent path; deployed_index_id carries the entry's key.
type vertexAiDeployedIndexVerifier struct{}

func (v *vertexAiDeployedIndexVerifier) IDOutputKey() string { return "deployed_index_id" }

// find fetches the parent endpoint and returns the matching deployedIndexes[]
// entry (nil when the deployment is absent). The HTTP status is the
// endpoint GET's status — 404 means the whole endpoint is gone, which also
// means the deployment is.
func (v *vertexAiDeployedIndexVerifier) find(ctx context.Context, svc *Services, outputs map[string]string) (*vertexAiIndexEndpointDeploy, int, error) {
	endpointID := outputs["index_endpoint"]
	if endpointID == "" {
		return nil, 0, errors.New("index_endpoint output missing — cannot locate the parent endpoint")
	}
	deployedIndexID := outputs["deployed_index_id"]

	endpoint, status, err := getVertexAiIndexEndpoint(ctx, svc, endpointID)
	if err != nil {
		return nil, status, err
	}
	for i := range endpoint.DeployedIndexes {
		if endpoint.DeployedIndexes[i].ID == deployedIndexID {
			return &endpoint.DeployedIndexes[i], status, nil
		}
	}
	return nil, status, nil
}

func (v *vertexAiDeployedIndexVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	deployedIndexID := outputs["deployed_index_id"]
	if deployedIndexID == "" {
		return errors.New("deployed_index_id output missing after deploy")
	}

	entry, _, err := v.find(ctx, svc, outputs)
	if err != nil {
		return errors.Wrapf(err, "failed to probe vertex ai deployed index %s after deploy", deployedIndexID)
	}
	if entry == nil {
		return errors.Errorf("vertex ai deployed index %s not present in endpoint %s deployedIndexes after deploy", deployedIndexID, outputs["index_endpoint"])
	}
	return nil
}

func (v *vertexAiDeployedIndexVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	deployedIndexID := outputs["deployed_index_id"]
	if deployedIndexID == "" {
		return nil
	}

	entry, status, err := v.find(ctx, svc, outputs)
	if err != nil {
		// The endpoint itself being gone (chain destroy) also proves the
		// deployment is gone.
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing vertex ai deployed index %s after destroy", deployedIndexID)
	}
	if entry != nil {
		return errors.Errorf("vertex ai deployed index %s still present in endpoint %s deployedIndexes after destroy", deployedIndexID, outputs["index_endpoint"])
	}
	return nil
}
