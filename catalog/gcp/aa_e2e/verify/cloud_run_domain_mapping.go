package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// cloudRunDomainMappingVerifier probes a Cloud Run domain mapping through
// the Knative-style v1 API at its REGIONAL endpoint
// ({region}-run.googleapis.com — the same URL shape the provider itself
// addresses; domain mappings are outside the run/v2 typed client's
// surface, so the probe rides the shared ADC-authenticated RestClient,
// the regional-secret grain).
//
// Posture: the mapping's Ready condition is asserted to be the provider's
// own create-success contract — True, or Unknown with reason
// CertificatePending (the documented resting state of a mapping whose DNS
// records are not yet published; the E2E lane deliberately never
// publishes them).
type cloudRunDomainMappingVerifier struct{}

func (v *cloudRunDomainMappingVerifier) IDOutputKey() string { return "domain" }

// mappingURL builds the regional Knative-API URL for the mapping.
func (v *cloudRunDomainMappingVerifier) mappingURL(svc *Services, outputs map[string]string) (string, error) {
	domain := outputs["domain"]
	region := outputs["region"]
	if domain == "" || region == "" {
		return "", errors.New("domain or region output missing")
	}
	return fmt.Sprintf("https://%s-run.googleapis.com/apis/domains.cloudrun.com/v1/namespaces/%s/domainmappings/%s",
		region, svc.Project, domain), nil
}

// domainMappingObject is the slice of the Knative response the assertions
// read.
type domainMappingObject struct {
	Spec struct {
		RouteName string `json:"routeName"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
			Reason string `json:"reason"`
		} `json:"conditions"`
	} `json:"status"`
}

func (v *cloudRunDomainMappingVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	url, err := v.mappingURL(svc, outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	status, body, err := v.get(ctx, svc, url)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.Errorf("domain mapping %s probe returned %d after deploy", url, status)
	}

	mapping := domainMappingObject{}
	if err := json.Unmarshal(body, &mapping); err != nil {
		return errors.Wrap(err, "failed to decode domain mapping response")
	}

	// The mapping must point at the service the outputs claim.
	if mappedRoute := outputs["mapped_route_name"]; mappedRoute != "" && mapping.Spec.RouteName != mappedRoute {
		return errors.Errorf("domain mapping route mismatch: output %q, live %q",
			mappedRoute, mapping.Spec.RouteName)
	}

	// Ready=True (DNS published, certificate issued) or the documented
	// CertificatePending resting state both prove a healthy mapping; any
	// other condition is a real failure (unverified domain, missing
	// route, ...).
	for _, condition := range mapping.Status.Conditions {
		if condition.Type != "Ready" {
			continue
		}
		if condition.Status == "True" {
			return nil
		}
		if condition.Status == "Unknown" && condition.Reason == "CertificatePending" {
			return nil
		}
		return errors.Errorf("domain mapping Ready condition is %s (reason %s), want True or Unknown/CertificatePending",
			condition.Status, condition.Reason)
	}
	return errors.New("domain mapping has no Ready condition")
}

func (v *cloudRunDomainMappingVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	url, err := v.mappingURL(svc, outputs)
	if err != nil {
		// Without a URL there is nothing left to probe; treat as gone.
		return nil
	}

	// Domain-mapping deletion is asynchronous behind the Knative API: a GET
	// issued seconds after a successful destroy can still answer 200 while
	// the delete settles (live-caught: the provider's destroy returned
	// clean and the probe 1.6s later still read the mapping; a later probe
	// read 404). Poll briefly before declaring the mapping still present --
	// the read-after-delete posture the service-account and peering
	// verifiers established.
	deadline := time.Now().Add(90 * time.Second)
	for {
		status, _, err := v.get(ctx, svc, url)
		if err != nil {
			return err
		}
		if status == http.StatusNotFound {
			return nil
		}
		if status != http.StatusOK {
			return errors.Errorf("unexpected status %d probing domain mapping %s after destroy", status, url)
		}
		if time.Now().After(deadline) {
			return errors.Errorf("domain mapping %s still exists after destroy", url)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
		}
	}
}

// get issues an authenticated GET and returns the status code and body.
func (v *cloudRunDomainMappingVerifier) get(ctx context.Context, svc *Services, url string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to build domain mapping GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "domain mapping GET request failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to read domain mapping response")
	}
	return resp.StatusCode, body, nil
}
