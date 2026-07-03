package runner

import (
	"os"
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

// ExpandManifestTokens substitutes RunIDToken occurrences in the manifest with
// the run's unique id and returns the path to the expanded copy (a temp file,
// so the source manifest is never modified). Manifests without the token pass
// through untouched — the original path is returned and no file is written.
//
// Expansion is a plain text substitution performed before the manifest is
// parsed, so the token can appear in any field. Only cloud-side identifier
// fields should carry it: metadata names must stay stable because prerequisite
// FK resolution and human debugging both key off them.
func ExpandManifestTokens(manifestPath, runID string) (string, error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read manifest %s for token expansion", manifestPath)
	}
	if !strings.Contains(string(raw), RunIDToken) {
		return manifestPath, nil
	}
	if runID == "" {
		return "", errors.Errorf("manifest %s uses %s but no run id was provided", manifestPath, RunIDToken)
	}

	expanded := strings.ReplaceAll(string(raw), RunIDToken, runID)

	// A temp file (not next to the scenario) so scenario discovery never picks it up.
	tmpFile, err := os.CreateTemp("", "planton-e2e-expanded-*.yaml")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp file for expanded manifest")
	}
	if _, err := tmpFile.WriteString(expanded); err != nil {
		tmpFile.Close()
		return "", errors.Wrap(err, "failed to write expanded manifest")
	}
	if err := tmpFile.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close expanded manifest")
	}
	return tmpFile.Name(), nil
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
