package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tt "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/pkg/crkreflect"
)

// Phase represents a stage in the E2E test lifecycle.
type Phase string

const (
	PhaseDepsUp      Phase = "DEPENDENCIES-UP"
	PhaseSetup       Phase = "SETUP"
	PhaseValidate    Phase = "VALIDATE"
	PhaseDeploy      Phase = "DEPLOY"
	PhaseIdempotency Phase = "IDEMPOTENCY"
	PhaseVerifyOut   Phase = "VERIFY-OUT"
	PhaseVerifyRes   Phase = "VERIFY-RES"
	PhaseImportRT    Phase = "IMPORT-RT"
	PhaseDestroy     Phase = "DESTROY"
	PhaseVerifyCln   Phase = "VERIFY-CLN"
	PhaseDepsDn      Phase = "DEPENDENCIES-DOWN"
)

// PhaseResult captures the outcome of a single lifecycle phase.
type PhaseResult struct {
	Phase    Phase
	Passed   bool
	Duration time.Duration
	Error    error
}

// TestResult captures the full 6-phase lifecycle outcome for a component.
type TestResult struct {
	Component string
	Engine    string
	Phases    []PhaseResult
	Passed    bool
	Duration  time.Duration
}

// RunComponentTest executes the E2E lifecycle for a single component.
// If the component has dependencies (registry prerequisites), they are deployed
// first and torn down last, wrapping the standard 6-phase lifecycle.
func RunComponentTest(ctx context.Context, tc *provider.ComponentTestContext, harness provider.Harness) *TestResult {
	start := time.Now()
	result := &TestResult{
		Component: tc.Component,
		Engine:    tc.Engine,
		Passed:    true,
	}

	// Expand per-run unique-id tokens before anything parses the manifest, so
	// identifiers that cloud providers reserve across soft-delete windows get a
	// fresh value on every run (and on each engine within a run).
	expandRunID := EngineScopedRunID(tc.RunID, tc.Engine)
	expandedPath, err := ExpandManifestTokens(tc.ManifestPath, expandRunID, ScenarioSlug(tc.ManifestPath))
	if err != nil {
		result.Passed = false
		result.Phases = append(result.Phases, PhaseResult{
			Phase:  PhaseValidate,
			Passed: false,
			Error:  errors.Wrap(err, "failed to expand manifest tokens"),
		})
		result.Duration = time.Since(start)
		return result
	}
	tc.ManifestPath = expandedPath

	verifyCtx := context.WithValue(ctx, provider.ManifestPathKey{}, tc.ManifestPath)

	// Phase 0: deploy dependencies (registry prerequisites merged with any the
	// scenario manifest declares via its e2e-prerequisites annotation)
	var dependencyStates []DependencyState
	if tc.RepoRoot != "" {
		depStart := time.Now()
		var err error
		// The engine-scoped id is passed down so prerequisite manifests expand to
		// the same values as the scenario under test (their tokens must line up),
		// and so each engine's prerequisite deploys get distinct identifiers.
		dependencyStates, err = DeployDependencies(ctx, tc.RepoRoot, tc.Provider, tc.Component, tc.ManifestPath, tc.BackendURL, expandRunID, harness)
		pr := PhaseResult{
			Phase:    PhaseDepsUp,
			Duration: time.Since(depStart),
			Passed:   err == nil,
			Error:    err,
		}
		if len(dependencyStates) > 0 || err != nil {
			result.Phases = append(result.Phases, pr)
		}
		if err != nil {
			result.Passed = false
			// The run already failed, but a teardown failure still gets its own
			// phase entry: it means prerequisite resources may be leaking.
			if tdErr := TeardownDependencies(dependencyStates); tdErr != nil {
				result.Phases = append(result.Phases, PhaseResult{
					Phase:  PhaseDepsDn,
					Passed: false,
					Error:  tdErr,
				})
			}
			result.Duration = time.Since(start)
			return result
		}

		// Resolve the component manifest's value_from refs against the deployed
		// prerequisites' outputs -- the orchestrator's resolution step, performed
		// here so a composed topology (e.g. subnet -> vpc) can be tested standalone.
		if len(dependencyStates) > 0 {
			depOutputs := make(DependencyOutputs, len(dependencyStates))
			for _, depState := range dependencyStates {
				kind := crkreflect.KindFromString(depState.Dependency.KindSlug)
				if depOutputs[kind] == nil {
					depOutputs[kind] = make(map[string]map[string]interface{})
				}
				depOutputs[kind][depState.ManifestName] = depState.Outputs
			}
			resolvedPath, resolveErr := ResolveManifestRefs(tc.ManifestPath, depOutputs)
			if resolveErr != nil {
				result.Passed = false
				result.Phases = append(result.Phases, PhaseResult{
					Phase:    PhaseValidate,
					Duration: 0,
					Passed:   false,
					Error:    errors.Wrap(resolveErr, "failed to resolve manifest references from dependency outputs"),
				})
				if tdErr := TeardownDependencies(dependencyStates); tdErr != nil {
					result.Phases = append(result.Phases, PhaseResult{
						Phase:  PhaseDepsDn,
						Passed: false,
						Error:  tdErr,
					})
				}
				result.Duration = time.Since(start)
				return result
			}
			tc.ManifestPath = resolvedPath
			verifyCtx = context.WithValue(ctx, provider.ManifestPathKey{}, tc.ManifestPath)
		}
	}

	// SETUP: the scenario's data-plane seeding hook (see SetupScriptAnnotation).
	// It runs after the dependency chain and reference resolution because the
	// assets it seeds live INSIDE fixtures, and before VALIDATE because a
	// seeding failure must stop the lane before any component deploy.
	if setupScript, annErr := ManifestAnnotation(tc.ManifestPath, SetupScriptAnnotation); annErr == nil && setupScript != "" {
		setupStart := time.Now()
		setupErr := runSetupScript(tc, setupScript, expandRunID)
		result.Phases = append(result.Phases, PhaseResult{
			Phase:    PhaseSetup,
			Duration: time.Since(setupStart),
			Passed:   setupErr == nil,
			Error:    setupErr,
		})
		if setupErr != nil {
			result.Passed = false
			if tdErr := TeardownDependencies(dependencyStates); tdErr != nil {
				result.Phases = append(result.Phases, PhaseResult{
					Phase:  PhaseDepsDn,
					Passed: false,
					Error:  tdErr,
				})
			}
			result.Duration = time.Since(start)
			return result
		}
	}

	// Phases 1-6 (7 with the opt-in import round-trip): standard lifecycle.
	type lifecyclePhase struct {
		phase Phase
		fn    func() error
	}
	phases := []lifecyclePhase{
		{PhaseValidate, func() error { return runValidate(tc) }},
		{PhaseDeploy, func() error { return runDeploy(tc) }},
	}
	// The idempotency gate re-plans the configuration DEPLOY just applied and
	// fails on any pending change: a dirty second plan means the module and
	// the provider disagree about the applied state (the send-omitted-value /
	// Optional+Computed echo defect classes), which users would meet as a
	// perpetual diff on every re-apply. Provider profiles arm it via
	// assert_apply_idempotency; prerequisite fixtures are deliberately outside
	// its scope (they belong to other kinds' contracts).
	if tc.AssertApplyIdempotency {
		phases = append(phases, lifecyclePhase{PhaseIdempotency, func() error { return runIdempotency(tc) }})
	}
	phases = append(phases,
		lifecyclePhase{PhaseVerifyOut, func() error { return runVerifyOutputs(tc) }},
		lifecyclePhase{PhaseVerifyRes, func() error { return runVerifyResources(verifyCtx, tc, harness) }},
	)
	// The import round-trip slots between VERIFY-RES and DESTROY: it needs the
	// deployed fixture live (imports read the real cloud) and the destroy that
	// follows tears down through the re-imported state -- itself part of the
	// proof that the blind import fully owns the resources.
	if importRoundTripEnabled(tc) {
		phases = append(phases, lifecyclePhase{PhaseImportRT, func() error { return runImportRoundTrip(tc) }})
	}
	phases = append(phases,
		lifecyclePhase{PhaseDestroy, func() error { return runDestroy(tc) }},
		lifecyclePhase{PhaseVerifyCln, func() error { return runVerifyCleanup(verifyCtx, tc, harness) }},
	)

	for _, p := range phases {
		phaseStart := time.Now()
		err := p.fn()
		pr := PhaseResult{
			Phase:    p.phase,
			Duration: time.Since(phaseStart),
			Passed:   err == nil,
			Error:    err,
		}
		result.Phases = append(result.Phases, pr)

		if err != nil {
			result.Passed = false
			if p.phase == PhaseDeploy || p.phase == PhaseIdempotency || p.phase == PhaseVerifyOut || p.phase == PhaseVerifyRes || p.phase == PhaseImportRT {
				cleanupErr := runDestroy(tc)
				if cleanupErr != nil {
					fmt.Printf("  [WARN] cleanup destroy also failed: %v\n", cleanupErr)
				}
			}
			break
		}
	}

	// Phase 7: teardown dependencies in reverse order. A teardown failure
	// FAILS the run even when every lifecycle phase passed: it means
	// prerequisite cloud resources may still exist, and a green result would
	// hide that leak until someone audits the account.
	if len(dependencyStates) > 0 {
		depStart := time.Now()
		tdErr := TeardownDependencies(dependencyStates)
		result.Phases = append(result.Phases, PhaseResult{
			Phase:    PhaseDepsDn,
			Duration: time.Since(depStart),
			Passed:   tdErr == nil,
			Error:    tdErr,
		})
		if tdErr != nil {
			result.Passed = false
		}
	}

	// Clean up Terraform working directory if one was created
	if tc.TerraformCleanup != nil {
		tc.TerraformCleanup()
	}

	result.Duration = time.Since(start)
	return result
}

func runValidate(tc *provider.ComponentTestContext) error {
	if tc.ManifestPath == "" {
		return errors.New("manifest path is empty")
	}

	// The opt-in per-component provider-config fixture (nil for the default
	// ambient-credential posture) -- one resolution serving both engines.
	providerConfig, err := LoadProviderConfigFixture(tc.ModuleDir, tc.ManifestPath)
	if err != nil {
		return errors.Wrap(err, "validation failed: provider-config fixture")
	}

	switch tc.Engine {
	case "pulumi":
		stackInputPath, err := BuildStackInput(tc.ManifestPath, providerConfig)
		if err != nil {
			return errors.Wrap(err, "validation failed: cannot build stack input from manifest")
		}
		tc.StackInputFilePath = stackInputPath

	case "terraform":
		workDir, cleanup, err := PrepareWorkDir(tc.ModuleDir)
		if err != nil {
			return errors.Wrap(err, "validation failed: cannot prepare terraform working directory")
		}
		tc.TerraformWorkDir = workDir
		tc.TerraformCleanup = cleanup

		input, err := BuildTerraformInput(tc.ManifestPath, workDir, providerConfig)
		if err != nil {
			cleanup()
			return errors.Wrap(err, "validation failed: cannot build terraform input from manifest")
		}

		tc.TerraformOpts = BuildTerratestOptions(tc.T, workDir, input.TfvarsPath, input.EnvVars)

	default:
		return errors.Errorf("unsupported engine for validation: %s", tc.Engine)
	}

	return nil
}

func runDeploy(tc *provider.ComponentTestContext) error {
	switch tc.Engine {
	case "pulumi":
		_, err := PulumiDeploy(tc.ModuleDir, tc.StackName, tc.BackendURL, tc.StackInputFilePath)
		return err
	case "terraform":
		opts, ok := tc.TerraformOpts.(*tt.Options)
		if !ok || opts == nil {
			return errors.New("terraform options not initialized (runValidate must run first)")
		}
		_, err := TerraformDeploy(tc.T, opts)
		return err
	default:
		return errors.Errorf("unsupported engine: %s", tc.Engine)
	}
}

// runIdempotency re-plans the configuration runDeploy just applied and fails
// on any pending change. See the phase-insertion comment in RunComponentTest
// for why this gate exists and what a failure means.
func runIdempotency(tc *provider.ComponentTestContext) error {
	switch tc.Engine {
	case "pulumi":
		_, err := PulumiPreviewExpectNoChanges(tc.ModuleDir, tc.StackName, tc.BackendURL, tc.StackInputFilePath)
		return err
	case "terraform":
		opts, ok := tc.TerraformOpts.(*tt.Options)
		if !ok || opts == nil {
			return errors.New("terraform options not initialized (runValidate must run first)")
		}
		_, err := TerraformPlanNoChanges(tc.T, opts)
		return err
	default:
		return errors.Errorf("unsupported engine: %s", tc.Engine)
	}
}

func runVerifyOutputs(tc *provider.ComponentTestContext) error {
	switch tc.Engine {
	case "pulumi":
		outputJSON, err := PulumiStackOutputs(tc.ModuleDir, tc.StackName, tc.BackendURL)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve pulumi stack outputs")
		}
		rawOutputs, parseErr := parsePulumiOutputs(outputJSON)
		if parseErr != nil {
			return errors.Wrap(parseErr, "failed to parse pulumi stack outputs JSON")
		}
		tc.Outputs = rawOutputs

	case "terraform":
		opts, ok := tc.TerraformOpts.(*tt.Options)
		if !ok || opts == nil {
			return errors.New("terraform options not initialized (runValidate must run first)")
		}
		rawOutputs, err := TerraformOutputs(tc.T, opts)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve terraform outputs")
		}
		tc.Outputs = rawOutputs

	default:
		return nil
	}

	if len(tc.Outputs) == 0 {
		fmt.Printf("  [outputs] %s: no outputs captured, skipping transformation validation\n", tc.Component)
		return nil
	}

	msg, flatOutputs, err := VerifyOutputTransformation(tc.Component, tc.Outputs, tc.ModuleDir)
	if err != nil {
		return err
	}
	tc.FlatOutputs = flatOutputs
	tc.TransformedOutputs = msg
	return nil
}

func runVerifyResources(ctx context.Context, tc *provider.ComponentTestContext, harness provider.Harness) error {
	return harness.VerifyDeployed(ctx, tc.Component, tc.Outputs)
}

func runDestroy(tc *provider.ComponentTestContext) error {
	switch tc.Engine {
	case "pulumi":
		_, err := PulumiDestroy(tc.ModuleDir, tc.StackName, tc.BackendURL, tc.StackInputFilePath)
		if err != nil {
			return err
		}
		return PulumiRemoveStack(tc.ModuleDir, tc.StackName, tc.BackendURL)
	case "terraform":
		opts, ok := tc.TerraformOpts.(*tt.Options)
		if !ok || opts == nil {
			return errors.New("terraform options not initialized")
		}
		_, err := TerraformDestroy(tc.T, opts)
		return err
	default:
		return errors.Errorf("unsupported engine: %s", tc.Engine)
	}
}

func runVerifyCleanup(ctx context.Context, tc *provider.ComponentTestContext, harness provider.Harness) error {
	return harness.VerifyDestroyed(ctx, tc.Component)
}

// parsePulumiOutputs converts the JSON string from `pulumi stack output --json`
// into a map[string]interface{} compatible with tc.Outputs.
func parsePulumiOutputs(outputJSON string) (map[string]interface{}, error) {
	if outputJSON == "" {
		return nil, nil
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(outputJSON), &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
