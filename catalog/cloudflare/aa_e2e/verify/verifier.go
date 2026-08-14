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
	"time"

	"github.com/pkg/errors"
)

// EnvelopePresence selects which 200-body signals count as ABSENT. The zero
// value treats any 200 as present (same answer as ResourceExists, but
// ResourcePresent still parses the v4 envelope -- never use it on raw-body
// endpoints).
type EnvelopePresence struct {
	// SoftDeleted: a 200 whose result.deleted_at is set counts as absent
	// (cloudflared tunnels, teamnet routes).
	SoftDeleted bool
	// AbsentStatuses: a 200 whose result.status is one of these counts as
	// absent (certificate packs, custom hostnames, fallback origin --
	// pending_deletion / deleted).
	AbsentStatuses []string
}

// API is the surface verifiers need from the harness client: an
// existence probe plus the account scope resolved at Setup.
type API interface {
	// ResourceExists GETs a path relative to /client/v4/ and reports
	// presence (200) vs absence (404, or Cloudflare's 400/7003 unknown-
	// object answer); any other status is an error.
	ResourceExists(ctx context.Context, path string) (bool, error)
	// ResourceActive is the soft-delete-aware sibling: a 200 whose envelope
	// carries a non-null result.deleted_at counts as ABSENT. Used only by
	// verifiers whose family soft-deletes (cloudflared tunnels, teamnet
	// routes) -- it parses the v4 envelope, so it must never replace
	// ResourceExists on raw-body endpoints (e.g. the KV value endpoint).
	ResourceActive(ctx context.Context, path string) (bool, error)
	// ResourcePresent is the envelope-aware probe: 404 / 400-7003 are
	// always absent; a 200 is absent when SoftDeleted sees deleted_at or
	// result.status matches AbsentStatuses. Parses the v4 envelope --
	// never use on raw-body endpoints.
	ResourcePresent(ctx context.Context, path string, opts EnvelopePresence) (bool, error)
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
	// softDeleted opts into the deleted_at-aware probe for families whose
	// destroy never yields a 404: the GET keeps answering 200 with a
	// non-null deleted_at (measured on cloudflared tunnels and teamnet
	// routes at provider v5.23.0 -- the provider's own destroy check asserts
	// deleted_at != nil, and its test sweepers filter on it). Without this,
	// verify-destroyed would report a false "still exists" on every lane.
	softDeleted bool
	// absentStatuses opts into the status-enum absence probe: a 200 whose
	// result.status is one of these counts as gone (certificate packs,
	// custom hostnames, fallback origin -- pending_deletion / deleted).
	absentStatuses []string
	// absentRetries is the number of VerifyAbsent probes, 1s apart. List
	// items delete through an async bulk POST the provider never polls; a
	// single GET immediately after destroy races the 404. Zero/one means
	// a single probe (everyone else's default).
	absentRetries int
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
	exists, err := v.probe(ctx, api, path)
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
	attempts := v.absentRetries
	if attempts < 1 {
		attempts = 1
	}
	var exists bool
	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		var probeErr error
		exists, probeErr = v.probe(ctx, api, path)
		if probeErr != nil {
			return errors.Wrapf(probeErr, "%s verify-absent failed", v.component)
		}
		if !exists {
			return nil
		}
	}
	return errors.Errorf("%s %s still exists after destroy (GET %s)",
		v.component, outputs[v.IDOutputKey()], path)
}

// probe selects the existence check: the plain status-code probe, or the
// envelope-aware one for soft-deleting / status-enum families.
func (v *apiPathVerifier) probe(ctx context.Context, api API, path string) (bool, error) {
	if v.softDeleted || len(v.absentStatuses) > 0 {
		return api.ResourcePresent(ctx, path, EnvelopePresence{
			SoftDeleted:    v.softDeleted,
			AbsentStatuses: v.absentStatuses,
		})
	}
	return api.ResourceExists(ctx, path)
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
	// The Access kinds are dual-scope (account XOR zone) but their live
	// scenarios are account-scoped -- the common, reusable case. A
	// zone-scoped arm would need zones/%s/access/... variants reading a
	// zone_id the outputs do not carry today (the ruleset precedent:
	// recorded as not-wired, a follow-up if a zone-scoped arm is ever
	// queued).
	"cloudflarezerotrustaccessapplication": &apiPathVerifier{
		component:     "cloudflarezerotrustaccessapplication",
		pathFormat:    "accounts/%s/access/apps/%s",
		outputKeys:    []string{"application_id"},
		accountScoped: true,
	},
	"cloudflarezerotrustaccessgroup": &apiPathVerifier{
		component:     "cloudflarezerotrustaccessgroup",
		pathFormat:    "accounts/%s/access/groups/%s",
		outputKeys:    []string{"group_id"},
		accountScoped: true,
	},
	// The policy is also the Access application's self-hosted fixture (a
	// scenario-declared prerequisite), so this entry serves both.
	"cloudflarezerotrustaccesspolicy": &apiPathVerifier{
		component:     "cloudflarezerotrustaccesspolicy",
		pathFormat:    "accounts/%s/access/policies/%s",
		outputKeys:    []string{"policy_id"},
		accountScoped: true,
	},
	// The tunnel is also the tunnel route's registry-prerequisite fixture.
	// Cloudflared tunnels SOFT-DELETE: the GET answers 200 with deleted_at
	// set after destroy, so the probe is the deleted_at-aware one.
	"cloudflarezerotrusttunnel": &apiPathVerifier{
		component:     "cloudflarezerotrusttunnel",
		pathFormat:    "accounts/%s/cfd_tunnel/%s",
		outputKeys:    []string{"tunnel_id"},
		accountScoped: true,
		softDeleted:   true,
	},
	// Routes and virtual networks live under teamnet/, not cfd_tunnel/.
	// Routes soft-delete like tunnels; virtual networks 404 honestly.
	"cloudflarezerotrusttunnelroute": &apiPathVerifier{
		component:     "cloudflarezerotrusttunnelroute",
		pathFormat:    "accounts/%s/teamnet/routes/%s",
		outputKeys:    []string{"route_id"},
		accountScoped: true,
		softDeleted:   true,
	},
	// Also the route's isolated-vnet fixture (a scenario-declared
	// prerequisite), so this entry serves both.
	"cloudflarezerotrusttunnelvirtualnetwork": &apiPathVerifier{
		component:     "cloudflarezerotrusttunnelvirtualnetwork",
		pathFormat:    "accounts/%s/teamnet/virtual_networks/%s",
		outputKeys:    []string{"virtual_network_id"},
		accountScoped: true,
	},
	// Certificate packs, custom hostnames, and the fallback origin answer
	// GET 200 with status pending_deletion/deleted after destroy rather
	// than 404ing (measured in the provider's own destroy checks at
	// v5.23.0). The status-enum probe treats those as gone.
	"cloudflarecertificatepack": &apiPathVerifier{
		component:      "cloudflarecertificatepack",
		pathFormat:     "zones/%s/ssl/certificate_packs/%s",
		outputKeys:     []string{"zone_id", "certificate_pack_id"},
		absentStatuses: []string{"pending_deletion", "deleted"},
	},
	"cloudflarecustomhostname": &apiPathVerifier{
		component:      "cloudflarecustomhostname",
		pathFormat:     "zones/%s/custom_hostnames/%s",
		outputKeys:     []string{"zone_id", "custom_hostname_id"},
		absentStatuses: []string{"pending_deletion", "deleted"},
	},
	// Zone singleton: no resource id -- zone_id is the identity and the
	// only output the verifier keys on (email-routing-settings precedent).
	"cloudflarecustomhostnamefallbackorigin": &apiPathVerifier{
		component:      "cloudflarecustomhostnamefallbackorigin",
		pathFormat:     "zones/%s/custom_hostnames/fallback_origin",
		outputKeys:     []string{"zone_id"},
		absentStatuses: []string{"pending_deletion", "deleted"},
	},
	"cloudflarelist": &apiPathVerifier{
		component:     "cloudflarelist",
		pathFormat:    "accounts/%s/rules/lists/%s",
		outputKeys:    []string{"list_id"},
		accountScoped: true,
	},
	// List-item delete is an async bulk POST the provider never polls;
	// a single GET immediately after destroy races the 404.
	"cloudflarelistitem": &apiPathVerifier{
		component:     "cloudflarelistitem",
		pathFormat:    "accounts/%s/rules/lists/%s/items/%s",
		outputKeys:    []string{"list_id", "item_id"},
		accountScoped: true,
		absentRetries: 5,
	},
	"cloudflareloadbalancermonitor": &apiPathVerifier{
		component:     "cloudflareloadbalancermonitor",
		pathFormat:    "accounts/%s/load_balancers/monitors/%s",
		outputKeys:    []string{"monitor_id"},
		accountScoped: true,
	},
	// Identity is the project NAME, not a UUID -- the provider's Read is
	// GET accounts/{a}/pages/projects/{name}.
	"cloudflarepagesproject": &apiPathVerifier{
		component:     "cloudflarepagesproject",
		pathFormat:    "accounts/%s/pages/projects/%s",
		outputKeys:    []string{"project_name"},
		accountScoped: true,
	},
	// Identity is the sitekey. Path segment is challenges/widgets, not
	// turnstile/.
	"cloudflareturnstilewidget": &apiPathVerifier{
		component:     "cloudflareturnstilewidget",
		pathFormat:    "accounts/%s/challenges/widgets/%s",
		outputKeys:    []string{"sitekey"},
		accountScoped: true,
	},
	// User-scoped (no accounts/ or zones/ prefix). Delete is a revoke;
	// a post-destroy 200 is possible and recorded on the watch-list --
	// the plain 404 probe is the honest starting point.
	"cloudflareorigincacertificate": &apiPathVerifier{
		component:  "cloudflareorigincacertificate",
		pathFormat: "certificates/%s",
		outputKeys: []string{"certificate_id"},
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
