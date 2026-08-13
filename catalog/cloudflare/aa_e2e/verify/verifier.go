// Package verify checks that Cloudflare resources created by an E2E scenario
// exist after DEPLOY and are gone after DESTROY, by GETting each resource on
// Cloudflare's own v4 REST API (the SaaS-provider pattern: one HTTP client,
// per-kind GET paths -- no cloud CLI).
//
// Standing convention: verification reads everything it needs from the
// component's stack outputs. A zone-scoped kind's outputs must therefore
// carry the zone_id alongside the resource's own identifier (a Cloudflare
// resource's API identity is compound -- zones/{zone_id}/<collection>/{id}).
// When enrolling a kind whose outputs lack its scope, add the output to the
// contract and both engine modules first; never fish the scope out of the
// manifest here.
package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// API is the surface verifiers need from the harness client: an
// existence probe plus the account scope resolved at Setup.
type API interface {
	// ResourceExists GETs a path relative to /client/v4/ and reports
	// presence (200) vs absence (404, or Cloudflare's 400/7003 unknown-
	// object answer); any other status is an error.
	ResourceExists(ctx context.Context, path string) (bool, error)
	// AccountID returns the Cloudflare account the harness is scoped to.
	AccountID() string
}

// Verifier checks a single component's Cloudflare resource for
// existence/absence from its stack outputs.
type Verifier interface {
	// IDOutputKey is the stack-output key carrying the resource's primary
	// identifier -- used to confirm the deploy produced a verifiable handle.
	IDOutputKey() string
	// VerifyExists returns an error unless the resource exists.
	VerifyExists(ctx context.Context, api API, outputs map[string]string) error
	// VerifyAbsent returns an error unless the resource is gone.
	VerifyAbsent(ctx context.Context, api API, outputs map[string]string) error
}

// apiPathVerifier is the common implementation: one GET path template whose
// placeholders are filled, in order, from the named stack-output keys.
// Account-scoped resources set accountScoped, which prepends the harness
// account to the path ("accounts/%s/..." with the first placeholder filled
// from API.AccountID(), never from outputs -- account IDs are harness scope,
// not deploy results).
type apiPathVerifier struct {
	component     string
	pathFormat    string
	outputKeys    []string
	accountScoped bool
}

// IDOutputKey returns the last output key -- the resource's own identifier
// (earlier keys are parent scope, e.g. zone_id).
func (v *apiPathVerifier) IDOutputKey() string {
	return v.outputKeys[len(v.outputKeys)-1]
}

func (v *apiPathVerifier) VerifyExists(ctx context.Context, api API, outputs map[string]string) error {
	path, err := v.buildPath(api, outputs)
	if err != nil {
		return err
	}
	exists, err := api.ResourceExists(ctx, path)
	if err != nil {
		return errors.Wrapf(err, "%s verify-exists failed", v.component)
	}
	if !exists {
		return errors.Errorf("%s %s not found after deploy (GET %s)",
			v.component, outputs[v.IDOutputKey()], path)
	}
	return nil
}

func (v *apiPathVerifier) VerifyAbsent(ctx context.Context, api API, outputs map[string]string) error {
	path, err := v.buildPath(api, outputs)
	if err != nil {
		return err
	}
	exists, err := api.ResourceExists(ctx, path)
	if err != nil {
		return errors.Wrapf(err, "%s verify-absent failed", v.component)
	}
	if exists {
		return errors.Errorf("%s %s still exists after destroy (GET %s)",
			v.component, outputs[v.IDOutputKey()], path)
	}
	return nil
}

// buildPath fills the path template from the account scope (when account-
// scoped) and the declared output keys, failing loudly on any missing output
// -- a blank segment would probe the wrong endpoint and report a false
// absence.
func (v *apiPathVerifier) buildPath(api API, outputs map[string]string) (string, error) {
	values := make([]interface{}, 0, len(v.outputKeys)+1)
	if v.accountScoped {
		if api.AccountID() == "" {
			return "", errors.Errorf("%s is account-scoped but the harness has no account ID", v.component)
		}
		values = append(values, api.AccountID())
	}
	for _, key := range v.outputKeys {
		value := outputs[key]
		if value == "" {
			return "", errors.Errorf("%s outputs carry no %q -- cannot verify", v.component, key)
		}
		values = append(values, value)
	}
	return fmt.Sprintf(v.pathFormat, values...), nil
}

// verifiers maps a component directory name to its verifier. A kind
// registers here in the wave that enrolls it for E2E (its profile +
// scenarios), and every kind that appears in another kind's registry
// prerequisites needs its entry before that consumer's lane can run --
// fixtures are verified right after install.
var verifiers = map[string]Verifier{
	// The zone is the central fixture: most zone-scoped lanes install it as
	// their prerequisite before the component under test deploys.
	"cloudflarednszone": &apiPathVerifier{
		component:  "cloudflarednszone",
		pathFormat: "zones/%s",
		outputKeys: []string{"zone_id"},
	},
	"cloudflarednsrecord": &apiPathVerifier{
		component:  "cloudflarednsrecord",
		pathFormat: "zones/%s/dns_records/%s",
		outputKeys: []string{"zone_id", "record_id"},
	},
}

// GetVerifier returns the verifier for a component, or an error if none is
// registered.
func GetVerifier(component string) (Verifier, error) {
	v, ok := verifiers[component]
	if !ok {
		return nil, errors.Errorf("no Cloudflare verifier registered for component %q", component)
	}
	return v, nil
}
