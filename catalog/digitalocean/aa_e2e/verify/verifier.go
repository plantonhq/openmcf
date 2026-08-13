// Package verify implements per-component resource verification for the
// DigitalOcean E2E harness. Each verifier answers two questions through the
// DigitalOcean REST API (godo): does the resource exist after deploy, and is
// it gone after destroy. A 404 from the API is the ONLY absence signal; every
// other error surfaces as a genuine failure, so flaky credentials or rate
// limits can never masquerade as a passing destroy verification.
//
// The one exception to godo is the Spaces bucket verifier: Spaces is an
// S3-compatible credential plane the API token cannot reach, so that verifier
// speaks the S3 API against the bucket's regional Spaces endpoint (see
// bucket.go).
package verify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/digitalocean/godo"
	pkgerrors "github.com/pkg/errors"
)

// Verifier checks a single component's DigitalOcean resource for
// existence/absence. DigitalOcean lookups are account-scoped by id, so there
// is no region parameter (region is a property of the resource, not of the
// API endpoint).
type Verifier interface {
	// IDOutputKey is the stack-output key carrying the identifier used to
	// verify the resource (e.g. "vpc_id"). The key names come from each
	// kind's outputs.proto -- they are contract, not convention.
	IDOutputKey() string
	// VerifyExists returns an error unless the resource exists.
	VerifyExists(ctx context.Context, client *godo.Client, id string) error
	// VerifyAbsent returns an error unless the resource is gone.
	VerifyAbsent(ctx context.Context, client *godo.Client, id string) error
}

// OutputsVerifier inspects the full stack output map when a single string id
// is insufficient (e.g. a DNS record is addressed by domain + record id, and
// a Spaces bucket by region + name).
type OutputsVerifier interface {
	Verifier
	VerifyExistsFromOutputs(ctx context.Context, client *godo.Client, outputs map[string]interface{}) error
	VerifyAbsentFromOutputs(ctx context.Context, client *godo.Client, outputs map[string]interface{}) error
}

// verifiers maps component slugs (the catalog directory names) to their
// verifiers. Every kind that appears in another kind's registry prerequisites
// MUST have an entry here, or composed scenarios fail at DEPENDENCIES-UP.
var verifiers = map[string]Verifier{
	"digitaloceanappplatformservice": &appVerifier{component: "digitaloceanappplatformservice", idOutputKey: "app_id"},
	"digitaloceanbucket":             &bucketVerifier{},
	"digitaloceancertificate":        &certificateVerifier{},
	"digitaloceancontainerregistry":  &containerRegistryVerifier{},
	"digitaloceandatabasecluster":    &databaseClusterVerifier{},
	"digitaloceandnsrecord":          &dnsRecordVerifier{},
	"digitaloceandnszone":            &dnsZoneVerifier{},
	"digitaloceandroplet":            &dropletVerifier{},
	"digitaloceanfirewall":           &firewallVerifier{},
	"digitaloceanfunction":           &appVerifier{component: "digitaloceanfunction", idOutputKey: "function_id"},
	"digitaloceankubernetescluster":  &kubernetesClusterVerifier{},
	"digitaloceankubernetesnodepool": &kubernetesNodePoolVerifier{},
	"digitaloceanloadbalancer":       &loadBalancerVerifier{},
	"digitaloceanvolume":             &volumeVerifier{},
	"digitaloceanvpc":                &vpcVerifier{},
}

// GetVerifier returns the verifier for a component, or an error if none is registered.
func GetVerifier(component string) (Verifier, error) {
	v, ok := verifiers[component]
	if !ok {
		return nil, pkgerrors.Errorf("no DigitalOcean verifier registered for component %q", component)
	}
	return v, nil
}

// isNotFound reports whether err is the DigitalOcean API's 404. godo wraps
// every non-2xx response in *godo.ErrorResponse carrying the raw
// *http.Response, so the status code is checked typed, never by matching
// error strings.
func isNotFound(err error) bool {
	var errResp *godo.ErrorResponse
	if errors.As(err, &errResp) && errResp.Response != nil {
		return errResp.Response.StatusCode == http.StatusNotFound
	}
	return false
}

// StringOutput reads a string-valued stack output, tolerating non-string
// scalars: DigitalOcean's numeric ids (droplets, DNS records) may decode as
// float64 or json.Number depending on the engine's JSON path, and a float64
// rendered with %v would turn 12345678 into "1.2345678e+07". Exported because
// the harness reads outputs with the same care.
func StringOutput(outputs map[string]interface{}, key string) string {
	if outputs == nil {
		return ""
	}
	v, ok := outputs[key]
	if !ok {
		return ""
	}
	switch n := v.(type) {
	case string:
		return n
	case float64:
		return strconv.FormatFloat(n, 'f', -1, 64)
	case json.Number:
		return n.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
