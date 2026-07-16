package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// vertexAiIndexVerifier probes a Vertex AI Vector Search index via the AI
// Platform REST API. The pinned google.golang.org/api line has no typed
// Vertex client at the harness's granularity, so the probe is a plain
// authenticated GET on the index's documented resource path (the sibling
// endpoint verifier's precedent). Posture assertions confirm the platform
// attribution labels landed (the label-parity proof) and that the
// index_name output matches the live numeric ID — the cross-engine
// determinism contract both modules share.
type vertexAiIndexVerifier struct{}

func (v *vertexAiIndexVerifier) IDOutputKey() string { return "index_id" }

type vertexAiIndex struct {
	Name              string            `json:"name"`
	DisplayName       string            `json:"displayName"`
	Labels            map[string]string `json:"labels"`
	MetadataSchemaUri string            `json:"metadataSchemaUri"`
}

func (v *vertexAiIndexVerifier) get(ctx context.Context, svc *Services, indexID string) (*vertexAiIndex, int, error) {
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/%s",
		regionFromVertexResource(indexID), indexID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build vertex ai index GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "vertex ai index GET request failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to read vertex ai index response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Errorf("vertex ai index GET %s returned %d: %s", indexID, resp.StatusCode, string(body))
	}

	index := &vertexAiIndex{}
	if err := json.Unmarshal(body, index); err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to decode vertex ai index")
	}
	return index, resp.StatusCode, nil
}

func (v *vertexAiIndexVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	indexID := outputs["index_id"]
	if indexID == "" {
		return errors.New("index_id output missing after deploy")
	}

	index, _, err := v.get(ctx, svc, indexID)
	if err != nil {
		return errors.Wrapf(err, "vertex ai index %s not found after deploy", indexID)
	}

	if index.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("vertex ai index %s missing the planton-ai_resource attribution label after deploy", indexID)
	}

	// The numeric index_name output must match the live resource ID —
	// proof that both engines derived the same cloud-side name.
	if got := outputs["index_name"]; got != "" {
		liveName := lastPathSegment(index.Name)
		if got != liveName {
			return errors.Errorf("vertex ai index %s index_name output %q does not match live name %q", indexID, got, liveName)
		}
	}

	if uri := outputs["metadata_schema_uri"]; uri != "" && index.MetadataSchemaUri != uri {
		return errors.Errorf("vertex ai index %s metadata_schema_uri output %q does not match live %q", indexID, uri, index.MetadataSchemaUri)
	}
	return nil
}

func (v *vertexAiIndexVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	indexID := outputs["index_id"]
	if indexID == "" {
		return nil
	}

	_, status, err := v.get(ctx, svc, indexID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing vertex ai index %s after destroy", indexID)
	}
	return errors.Errorf("vertex ai index %s still exists after destroy", indexID)
}
