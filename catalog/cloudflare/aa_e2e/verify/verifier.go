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
	// Existence GET is the provider Read path (accounts/.../workers/services/{name}),
	// not workers/scripts/{name}. Account comes from the harness, never from outputs.
	"cloudflareworker": &apiPathVerifier{
		component:     "cloudflareworker",
		pathFormat:    "accounts/%s/workers/services/%s",
		outputKeys:    []string{"script_name"},
		accountScoped: true,
	},
	"cloudflareloadbalancer": &apiPathVerifier{
		component:  "cloudflareloadbalancer",
		pathFormat: "zones/%s/load_balancers/%s",
		outputKeys: []string{"zone_id", "load_balancer_id"},
	},
	// The pool is account-scoped (it is the load balancer's registry
	// prerequisite, so its fixture is verified right after install).
	"cloudflareloadbalancerpool": &apiPathVerifier{
		component:     "cloudflareloadbalancerpool",
		pathFormat:    "accounts/%s/load_balancers/pools/%s",
		outputKeys:    []string{"pool_id"},
		accountScoped: true,
	},
	// Email Routing is a zone singleton: the GET returns the zone's routing
	// settings object rather than 404ing for a zone that never enabled it, so
	// the absent-check may need a value-aware probe (read `enabled`) once
	// measured live -- the existence probe is the honest starting point.
	"cloudflareemailroutingzone": &apiPathVerifier{
		component:  "cloudflareemailroutingzone",
		pathFormat: "zones/%s/email/routing",
		outputKeys: []string{"zone_id"},
	},
	"cloudflareemailroutingrule": &apiPathVerifier{
		component:  "cloudflareemailroutingrule",
		pathFormat: "zones/%s/email/routing/rules/%s",
		outputKeys: []string{"zone_id", "rule_id"},
	},
	// Destination addresses are account-scoped (shared across zones).
	"cloudflareemailroutingaddress": &apiPathVerifier{
		component:     "cloudflareemailroutingaddress",
		pathFormat:    "accounts/%s/email/routing/addresses/%s",
		outputKeys:    []string{"address_id"},
		accountScoped: true,
	},
	// Rulesets are dual-scope; the live scenarios are zone-scoped, and the
	// outputs carry zone_id only for zone-scoped rulesets (an account-scoped
	// arm would need an accounts/%s/rulesets/%s variant reading the harness
	// account).
	"cloudflareruleset": &apiPathVerifier{
		component:  "cloudflareruleset",
		pathFormat: "zones/%s/rulesets/%s",
		outputKeys: []string{"zone_id", "ruleset_id"},
	},
	// The namespace is also the KV pair's install fixture, so this entry
	// serves both its own lanes and the pair's fixture verification.
	"cloudflarekvnamespace": &apiPathVerifier{
		component:     "cloudflarekvnamespace",
		pathFormat:    "accounts/%s/storage/kv/namespaces/%s",
		outputKeys:    []string{"namespace_id"},
		accountScoped: true,
	},
	// A KV entry's existence GET is the value endpoint (200 with the raw
	// value when present). Keys are path-escaped by the API client; scenario
	// keys stay slash-free so the import ID stays parseable too.
	"cloudflareworkerskvpair": &apiPathVerifier{
		component:     "cloudflareworkerskvpair",
		pathFormat:    "accounts/%s/storage/kv/namespaces/%s/values/%s",
		outputKeys:    []string{"namespace_id", "key_name"},
		accountScoped: true,
	},
	// Note the singular "database" in the D1 path -- the provider's own Read
	// path, not the plural most Cloudflare collections use.
	"cloudflared1database": &apiPathVerifier{
		component:     "cloudflared1database",
		pathFormat:    "accounts/%s/d1/database/%s",
		outputKeys:    []string{"database_id"},
		accountScoped: true,
	},
	// The queue is also the R2 event-notification scenario's install fixture
	// (a scenario-declared prerequisite), so this entry serves both.
	"cloudflarequeue": &apiPathVerifier{
		component:     "cloudflarequeue",
		pathFormat:    "accounts/%s/queues/%s",
		outputKeys:    []string{"queue_id"},
		accountScoped: true,
	},
	// R2's identity is the bucket NAME (no separate id). Jurisdiction rides
	// a cf-r2-jurisdiction HEADER, not the path -- this probe is honest only
	// for default-jurisdiction buckets, which is all the live scenarios
	// create; a non-default arm would need a header-aware client extension.
	"cloudflarer2bucket": &apiPathVerifier{
		component:     "cloudflarer2bucket",
		pathFormat:    "accounts/%s/r2/buckets/%s",
		outputKeys:    []string{"bucket_name"},
		accountScoped: true,
	},
	"cloudflarehyperdriveconfig": &apiPathVerifier{
		component:     "cloudflarehyperdriveconfig",
		pathFormat:    "accounts/%s/hyperdrive/configs/%s",
		outputKeys:    []string{"hyperdrive_id"},
		accountScoped: true,
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
