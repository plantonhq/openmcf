package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/plantonhq/planton/e2e/framework/provider"
)

// SetupScriptAnnotation names a repo-committed bash script the runner
// executes as the SETUP phase -- after the dependency chain is deployed and
// the scenario's references are resolved, before the component under test
// deploys. Its value is a repo-root-relative path.
//
// The phase exists for DATA-PLANE seeding: assets a scenario needs that no
// catalog kind can create because they are content inside a fixture, not an
// ARM/cloud control-plane object (the first user: a trivial MLflow model
// registered into the fixture ML workspace, which Azure requires before it
// will provision a managed online deployment). Everything a catalog kind CAN
// create must keep entering through registry prerequisites or the
// e2e-prerequisites annotation -- this hook is not an alternative fixture
// mechanism, and a script that creates control-plane resources would dodge
// the teardown and orphan-sweep guarantees those paths carry.
//
// Contract:
//   - The script runs once per engine lane, via bash, from the repo root.
//   - It inherits the process environment (cloud CLI logins, ARM_* exports)
//     plus E2E_RUN_ID (engine-scoped) and E2E_SCENARIO.
//   - A non-zero exit fails the lane BEFORE the component deploys; the
//     dependency chain still tears down.
//   - It must be IDEMPOTENT per lane (an in-lane retry may re-run it) and
//     must seed ONLY into fixture-owned resources, so DEPENDENCIES-DOWN
//     destroys everything it created and the zero-orphan sweep stays honest.
//     No teardown pair exists by design.
const SetupScriptAnnotation = "planton.dev/e2e-setup-script"

// runSetupScript executes the scenario's setup script. The caller resolves
// the annotation; an empty scriptRel is the caller's bug, not a skip.
func runSetupScript(tc *provider.ComponentTestContext, scriptRel, engineScopedRunID string) error {
	scriptPath := filepath.Join(tc.RepoRoot, scriptRel)
	if _, err := os.Stat(scriptPath); err != nil {
		return errors.Wrapf(err, "setup script %s (from the %s annotation) is not readable", scriptRel, SetupScriptAnnotation)
	}

	fmt.Printf("  [setup] running %s\n", scriptRel)
	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = tc.RepoRoot
	cmd.Env = append(os.Environ(),
		"E2E_RUN_ID="+engineScopedRunID,
		"E2E_SCENARIO="+ScenarioSlug(tc.ManifestPath),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "setup script %s failed", scriptRel)
	}
	fmt.Printf("  [setup] %s completed\n", scriptRel)
	return nil
}
