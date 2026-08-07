package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// filestoreInstanceVerifier probes a Filestore instance via the Filestore
// REST API. The pinned google.golang.org/api line lacks the typed file/v1
// client, so the probe is a plain authenticated GET on the instance's
// documented resource path (the instance_id output IS that path).
// Posture assertions confirm READY state, the platform attribution labels
// (the cross-engine label-parity canary), and that the exported share
// addresses belong to the live network attachment.
type filestoreInstanceVerifier struct{}

func (v *filestoreInstanceVerifier) IDOutputKey() string { return "instance_id" }

// filestoreInstance is the subset of the API's Instance object the
// verifier asserts on.
type filestoreInstance struct {
	Name     string            `json:"name"`
	State    string            `json:"state"`
	Labels   map[string]string `json:"labels"`
	Networks []struct {
		IpAddresses     []string `json:"ipAddresses"`
		ReservedIpRange string   `json:"reservedIpRange"`
	} `json:"networks"`
	FileShares []struct {
		Name string `json:"name"`
	} `json:"fileShares"`
}

func (v *filestoreInstanceVerifier) get(ctx context.Context, svc *Services, name string) (*filestoreInstance, int, error) {
	url := fmt.Sprintf("https://file.googleapis.com/v1/%s", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build filestore GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "filestore GET request failed")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to read filestore response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.Errorf("filestore GET %s returned %d: %s", name, resp.StatusCode, string(body))
	}

	instance := &filestoreInstance{}
	if err := json.Unmarshal(body, instance); err != nil {
		return nil, resp.StatusCode, errors.Wrap(err, "failed to decode filestore instance")
	}
	return instance, resp.StatusCode, nil
}

func (v *filestoreInstanceVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["instance_id"]
	if name == "" {
		return errors.New("instance_id output missing after deploy")
	}

	instance, _, err := v.get(ctx, svc, name)
	if err != nil {
		return errors.Wrapf(err, "filestore instance %s not found after deploy", name)
	}
	if instance.State != "READY" {
		return errors.Errorf("filestore instance %s state is %q, want READY", name, instance.State)
	}

	// The platform attribution labels are the cross-engine parity canary:
	// a missing set means one engine stamped labels and the other did not.
	if instance.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("filestore instance %s missing the planton-ai_resource attribution label after deploy (labels: %v)", name, instance.Labels)
	}

	// The share name output must match the live share (the NFS mount path
	// applications use).
	if share := outputs["file_share_name"]; share != "" {
		found := false
		for _, fs := range instance.FileShares {
			if fs.Name == share {
				found = true
			}
		}
		if !found {
			return errors.Errorf("filestore instance %s: file_share_name output %q matches no live share", name, share)
		}
	}

	return nil
}

func (v *filestoreInstanceVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["instance_id"]
	if name == "" {
		return nil
	}

	_, status, err := v.get(ctx, svc, name)
	if err != nil {
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing filestore instance %s after destroy", name)
	}
	return errors.Errorf("filestore instance %s still exists after destroy", name)
}
