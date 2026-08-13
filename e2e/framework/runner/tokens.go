package runner

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// RunIDToken is the placeholder scenario and prerequisite manifests use for
// values that must be unique per test run. Some cloud resources reserve their
// user-chosen identifier long after deletion (soft-delete retention windows —
// e.g. workload identity pools and KMS key rings), so a manifest that hardcodes
// such an identifier can never be deployed twice: the second run (even the
// second engine of the same run) collides with the soft-deleted ghost of the
// first. Embedding this token in the identifier field keeps the manifest a
// plain, reviewable YAML file while giving every run a fresh identifier.
const RunIDToken = "${E2E_RUN_ID}"

// UnderscoreRunIDToken is RunIDToken for identifier classes that forbid
// hyphens (letters/numbers/underscores only — e.g. a Vertex AI
// deployed_index_id or a BigQuery dataset id). The engine-scoped run id
// carries a hyphenated engine suffix ("-p"/"-t"), so the plain token would
// break such fields; this variant expands with every hyphen replaced by an
// underscore. Use it whenever a run-scoped identifier lives in an
// underscore-only field.
const UnderscoreRunIDToken = "${E2E_RUN_ID_UNDERSCORE}"

// ScenarioToken is the placeholder for the running scenario's slug (the
// scenario file's base name without extension). It is the reservation-window
// sibling of RunIDToken: prerequisite install manifests redeploy per
// scenario, so a run-scoped cloud-side name alone is still DESTROYED AND
// RECREATED at every scenario boundary within one engine run — and resource
// classes that reserve a deleted name for minutes (Firestore database IDs)
// collide the next scenario's chain with the previous scenario's ghost.
// Embedding the scenario slug alongside the run id gives every scenario its
// own cloud-side identifier, so no name is ever recreated within a run
// (first hit: a Firestore-index kind growing a second scenario whose chain
// reinstalls the same database override).
const ScenarioToken = "${E2E_SCENARIO}"

// EnvTokenPrefix restricts which environment variables scenario manifests may
// reference through ${E2E_ENV:...} tokens. Batched real-cluster lanes need
// batch-specific values a committed manifest cannot carry honestly — IRSA
// role ARNs, bucket names, queue URLs differ per test account, and hardcoding
// one account's identifiers would make the scenario deploy a lie everywhere
// else. The batch bootstrap exports them as PLANTON_E2E_-prefixed variables;
// the prefix keeps expansion opt-in and scoped (a manifest can never read
// arbitrary process environment like AWS_SECRET_ACCESS_KEY).
const EnvTokenPrefix = "PLANTON_E2E_"

// envTokenPattern matches ${E2E_ENV:PLANTON_E2E_*} occurrences. The variable
// name is captured; names outside EnvTokenPrefix never match and therefore
// fail expansion loudly via the residual-token check below.
var envTokenPattern = regexp.MustCompile(`\$\{E2E_ENV:(` + EnvTokenPrefix + `[A-Z0-9_]+)\}`)

// ExpandManifestTokens substitutes RunIDToken and ScenarioToken occurrences
// in the manifest, expands ${E2E_ENV:PLANTON_E2E_*} environment tokens, and
// returns the path to the expanded copy (a temp file, so the source manifest
// is never modified). Manifests without tokens pass through untouched — the
// original path is returned and no file is written. scenario is the running
// scenario's slug ("" when there is none — a manifest carrying ScenarioToken
// then fails loudly rather than deploying the literal token).
//
// Expansion is a plain text substitution performed before the manifest is
// parsed, so the token can appear in any field. Only cloud-side identifier
// fields should carry it: metadata names must stay stable because prerequisite
// FK resolution and human debugging both key off them.
func ExpandManifestTokens(manifestPath, runID, scenario string) (string, error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read manifest %s for token expansion", manifestPath)
	}
	hasEnvTokens := strings.Contains(string(raw), "${E2E_ENV:")
	hasScenarioToken := strings.Contains(string(raw), ScenarioToken)
	if !strings.Contains(string(raw), RunIDToken) && !strings.Contains(string(raw), UnderscoreRunIDToken) && !hasEnvTokens && !hasScenarioToken {
		return manifestPath, nil
	}
	if (strings.Contains(string(raw), RunIDToken) || strings.Contains(string(raw), UnderscoreRunIDToken)) && runID == "" {
		return "", errors.Errorf("manifest %s uses %s but no run id was provided", manifestPath, RunIDToken)
	}
	if hasScenarioToken && scenario == "" {
		return "", errors.Errorf("manifest %s uses %s but no scenario is in scope", manifestPath, ScenarioToken)
	}

	// The underscore variant substitutes first: RunIDToken is not a textual
	// prefix of it (the closing brace differs), but explicit ordering keeps
	// the intent obvious.
	expanded := strings.ReplaceAll(string(raw), UnderscoreRunIDToken, strings.ReplaceAll(runID, "-", "_"))
	expanded = strings.ReplaceAll(expanded, RunIDToken, runID)
	expanded = strings.ReplaceAll(expanded, ScenarioToken, scenario)

	if hasEnvTokens {
		var missing []string
		expanded = envTokenPattern.ReplaceAllStringFunc(expanded, func(token string) string {
			name := envTokenPattern.FindStringSubmatch(token)[1]
			value := os.Getenv(name)
			if value == "" {
				missing = append(missing, name)
			}
			return value
		})
		if len(missing) > 0 {
			return "", errors.Errorf(
				"manifest %s uses environment tokens whose variables are unset: %s (exported by the real-cluster batch bootstrap)",
				manifestPath, strings.Join(missing, ", "))
		}
		// A residual ${E2E_ENV: token means the name fell outside the
		// allowed prefix or its syntax is malformed — never deploy it as
		// literal text.
		if idx := strings.Index(expanded, "${E2E_ENV:"); idx >= 0 {
			end := strings.IndexByte(expanded[idx:], '}')
			residual := expanded[idx:]
			if end >= 0 {
				residual = expanded[idx : idx+end+1]
			}
			return "", errors.Errorf(
				"manifest %s carries an invalid environment token %q: only %s-prefixed variables may be referenced",
				manifestPath, residual, EnvTokenPrefix)
		}
	}

	// A temp file (not next to the scenario, so discovery never picks it up)
	// that KEEPS the scenario's basename: verifier dispatch keys behavioral
	// variants off the scenario name in the manifest path, and a
	// random-only temp name would silently demote every token-carrying
	// behavioral scenario to its plain verifier.
	tmpDir, err := os.MkdirTemp("", "planton-e2e-expanded-*")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir for expanded manifest")
	}
	tmpPath := filepath.Join(tmpDir, filepath.Base(manifestPath))
	if err := os.WriteFile(tmpPath, []byte(expanded), 0o600); err != nil {
		return "", errors.Wrap(err, "failed to write expanded manifest")
	}
	return tmpPath, nil
}

// ScenarioSlug derives the ScenarioToken expansion from a scenario manifest
// path: the base filename without extension. Expanded/resolved manifest
// copies keep their basename (the temp-file contract above), so the slug is
// stable across every copy of the scenario.
func ScenarioSlug(scenarioManifestPath string) string {
	if scenarioManifestPath == "" {
		return ""
	}
	base := filepath.Base(scenarioManifestPath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// EngineScopedRunID derives the value substituted for RunIDToken from the test
// run's id and the engine executing the scenario. Both engines run the same
// scenario within one test invocation, so the run id alone is not unique enough:
// a soft-delete-reserving identifier created by the Pulumi run would collide
// with the Terraform run minutes later. The engine's first letter keeps the
// suffix short (identifier fields like workload identity pool ids cap at 32
// characters) while making each engine's expansion distinct.
func EngineScopedRunID(runID, engine string) string {
	if engine == "" {
		return runID
	}
	return runID + "-" + engine[:1]
}
