package runner

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// IsolateHelmEnvironment points Helm's repository configuration, repository
// cache, and registry credentials at a fresh per-run directory, for this
// process and every engine it spawns (terraform, tofu, pulumi and the module
// binaries inherit the environment).
//
// Why: Helm consults the machine's repository list even for a chart addressed
// by URL, and a stale entry there (a listed repository whose index cache is
// missing) fails every chart render and every release install with
// "no cached repo found". A lane must see the runner's clean Helm state, not
// the laptop's; otherwise a developer's `helm repo add` from months ago fails
// a Helm-based kind's lane for a reason that has nothing to do with the kind.
//
// Returns a cleanup that removes the directory. Variables already set by the
// invoker are respected (a CI job may want a shared cache).
func IsolateHelmEnvironment() (func(), error) {
	dir, err := os.MkdirTemp("", "planton-e2e-helm-*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the per-run Helm directory")
	}
	for name, value := range map[string]string{
		"HELM_REPOSITORY_CONFIG": filepath.Join(dir, "repositories.yaml"),
		"HELM_REPOSITORY_CACHE":  filepath.Join(dir, "cache"),
		"HELM_REGISTRY_CONFIG":   filepath.Join(dir, "registry.json"),
	} {
		if os.Getenv(name) != "" {
			continue
		}
		if err := os.Setenv(name, value); err != nil {
			return nil, errors.Wrapf(err, "failed to set %s", name)
		}
	}
	return func() { _ = os.RemoveAll(dir) }, nil
}
