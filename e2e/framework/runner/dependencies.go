package runner

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Dependency is a single prerequisite deployment that must exist before a
// component's own scenario is applied.
type Dependency struct {
	// KindSlug is the lowercase component directory name of the dependency
	// (e.g. "kubernetesgatewayapicrds").
	KindSlug string

	// ManifestPath is the absolute path to the KRM manifest deployed for it.
	ManifestPath string
}

// DependencyState tracks a deployed dependency so it can be torn down later.
type DependencyState struct {
	Dependency     Dependency
	ModuleDir      string
	StackName      string
	BackendURL     string
	StackInputPath string

	// ManifestName is the metadata.name of the deployed manifest. It keys this
	// instance's outputs in DependencyOutputs so a reference can address one
	// specific instance when an install profile deploys several of the same kind.
	ManifestName string

	// Outputs are the dependency's captured stack outputs. They are used both to
	// verify the dependency and to resolve the dependent component's value_from
	// references (see ResolveManifestRefs).
	Outputs map[string]interface{}
}

// scenarioPrerequisitesAnnotation is the manifest annotation through which a
// scenario declares prerequisites BEYOND the component kind's registry graph.
// The registry's `kind_meta.prerequisites` stays the honest statement of what
// the kind REQUIRES to deploy at all; this annotation is how one scenario says
// what it additionally COMPOSES -- optional references (a folded capacity
// provider's auto-scaling group, a subnet's route-table attachment) that only
// that scenario exercises. Without it, live-testing an optional composition
// would force a false prerequisite onto every consumer of the kind.
//
// The value is a comma-separated list where each entry is either:
//   - a registered kind name (e.g. "AwsSubnet,AwsAutoScalingGroup") --
//     installed through the kind's standard install profile
//     (prerequisite.yaml, else its minimal scenario), expanded through its
//     own prerequisite edges so naming the deepest kind of a chain is
//     enough, and skipped if the chain already deploys that kind; or
//   - a repo-relative manifest path (recognized by containing a "/") --
//     deployed as an EXTRA INSTANCE of whatever kind the manifest declares,
//     for scenarios that need more instances than the standard profiles
//     provide (e.g. a virtual-network peering's remote network). Path
//     entries deploy after the kind-driven chain, each preceded by any of
//     its own transitive prerequisites not already deployed, and never
//     substitute for the kind's install profile.
//
// All fixtures join the same transitive reference resolution, and teardown
// runs in reverse across the merged chain.
const scenarioPrerequisitesAnnotation = "planton.dev/e2e-prerequisites"

// ResolveDependencies returns the ordered, deduplicated list of prerequisite
// deployments a component needs before its own scenario is applied. Dependencies
// come from three sources, merged and expanded transitively in deploy-first order:
//
//  1. The component's CloudResourceKindMeta.prerequisites graph in the proto
//     registry: declaring `prerequisites: [X]` on a kind is enough for the
//     harness to install X first, with no per-component wiring.
//  2. The scenario manifest's `planton.dev/e2e-prerequisites` annotation, for
//     compositions that are optional on the kind and therefore must not be
//     registry prerequisites (see scenarioPrerequisitesAnnotation). Kind-name
//     entries join the graph; manifest-path entries deploy as extra instances
//     AFTER it, each preceded by its own not-yet-scheduled prerequisites.
//  3. The same annotation on each prerequisite's OWN install manifest: an
//     install profile that composes fixtures of other kinds (e.g. the
//     zip-backed Lambda scenario referencing the S3 object-set fixture)
//     declares them there, and the harness orders those fixtures BEFORE the
//     declaring kind so its value_from references resolve. Without this, a
//     dependency whose install manifest references sibling fixtures would
//     deploy before them and fail.
//
// Teardown runs in reverse, so the most foundational dependency is removed last.
//
// The install manifest for each prerequisite is, in order of preference:
//   - <dep>/v1/e2e/prerequisite.yaml      (the dependency's published install profile)
//   - <dep>/v1/e2e/scenarios/minimal.yaml (fallback to its minimal scenario)
//
// An install profile may carry MULTIPLE `---`-separated documents of the kind
// when real topologies require more than one instance -- e.g. a load balancer
// needs two subnets in different availability zones, so the subnet profile
// publishes a two-AZ pair. Each document deploys as its own stack and its
// outputs are captured under its own manifest name.
func ResolveDependencies(repoRoot, componentProvider, component, scenarioManifestPath string) ([]Dependency, error) {
	kind := crkreflect.KindFromString(component)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		// Not a registered kind (or an alias mismatch); no prerequisites.
		return nil, nil
	}

	// Copy before appending: Prerequisites returns the registry's own slice.
	roots := append([]cloudresourcekind.CloudResourceKind{}, crkreflect.Prerequisites(kind)...)

	declaredKinds, declaredPaths, err := scenarioDeclaredPrerequisites(scenarioManifestPath)
	if err != nil {
		return nil, err
	}
	for _, d := range declaredKinds {
		if d == kind {
			return nil, errors.Errorf("scenario %s declares the component's own kind %s as an E2E prerequisite", scenarioManifestPath, component)
		}
		roots = append(roots, d)
	}

	visited := make(map[cloudresourcekind.CloudResourceKind]bool)
	prereqs, err := expandPrerequisiteGraph(repoRoot, componentProvider, component, roots, visited)
	if err != nil {
		return nil, err
	}

	var deps []Dependency
	for _, p := range prereqs {
		slug := strings.ToLower(p.String())
		manifestPath, err := prerequisiteManifestPath(repoRoot, componentProvider, component, slug)
		if err != nil {
			return nil, err
		}
		deps = append(deps, Dependency{
			KindSlug:     slug,
			ManifestPath: manifestPath,
		})
	}

	// Manifest-path entries deploy after the kind-driven chain, in listed
	// order. Each is an EXTRA INSTANCE of its declared kind: it never marks
	// the kind as scheduled (it is not a substitute for the kind's install
	// profile), but its own transitive prerequisites join the chain first if
	// the graph has not already scheduled them.
	for _, entry := range declaredPaths {
		full := filepath.Join(repoRoot, entry)
		if !pathExists(full) {
			return nil, errors.Errorf("%s annotation entry %q: manifest not found at %s", scenarioPrerequisitesAnnotation, entry, full)
		}
		slug, err := manifestKindSlug(full)
		if err != nil {
			return nil, errors.Wrapf(err, "%s annotation entry %q", scenarioPrerequisitesAnnotation, entry)
		}
		entryKind := crkreflect.KindFromString(slug)
		pre, err := expandPrerequisiteGraph(repoRoot, componentProvider, component, crkreflect.Prerequisites(entryKind), visited)
		if err != nil {
			return nil, err
		}
		for _, p := range pre {
			pSlug := strings.ToLower(p.String())
			manifestPath, err := prerequisiteManifestPath(repoRoot, componentProvider, component, pSlug)
			if err != nil {
				return nil, err
			}
			deps = append(deps, Dependency{KindSlug: pSlug, ManifestPath: manifestPath})
		}
		deps = append(deps, Dependency{KindSlug: slug, ManifestPath: full})
	}
	return deps, nil
}

// expandPrerequisiteGraph topologically orders the root prerequisites plus
// everything they transitively require, deduplicated, deploy-first. Each
// kind's edges come from two sources: its registry prerequisites (the honest
// statement of what it needs to deploy at all) and the e2e-prerequisites
// annotation on its own install manifest (the fixtures that manifest's
// value_from references compose -- see ResolveDependencies point 3). Cycles
// across either edge source indicate a modeling mistake and fail loudly.
//
// visited carries kinds already scheduled by an earlier expansion (they are
// skipped, and newly visited kinds are added to it), so a scenario's
// manifest-path entries can extend one chain without re-deploying fixtures
// the graph already ordered. Only NEWLY visited kinds are returned.
func expandPrerequisiteGraph(repoRoot, componentProvider, component string, roots []cloudresourcekind.CloudResourceKind, visited map[cloudresourcekind.CloudResourceKind]bool) ([]cloudresourcekind.CloudResourceKind, error) {
	var result []cloudresourcekind.CloudResourceKind
	if visited == nil {
		visited = make(map[cloudresourcekind.CloudResourceKind]bool)
	}
	inStack := make(map[cloudresourcekind.CloudResourceKind]bool)

	var visit func(k cloudresourcekind.CloudResourceKind) error
	visit = func(k cloudresourcekind.CloudResourceKind) error {
		if inStack[k] {
			return errors.Errorf("prerequisite cycle detected at %s while resolving dependencies for %s", k.String(), component)
		}
		if visited[k] {
			return nil
		}
		inStack[k] = true

		edges := append([]cloudresourcekind.CloudResourceKind{}, crkreflect.Prerequisites(k)...)
		manifestDeclared, err := installManifestPrerequisites(repoRoot, componentProvider, component, k)
		if err != nil {
			return err
		}
		edges = append(edges, manifestDeclared...)

		for _, e := range edges {
			if err := visit(e); err != nil {
				return err
			}
		}

		inStack[k] = false
		visited[k] = true
		result = append(result, k)
		return nil
	}

	for _, r := range roots {
		if err := visit(r); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// installManifestPrerequisites reads the e2e-prerequisites annotation from a
// prerequisite kind's install manifest (every document of a multi-document
// profile), returning the extra kinds that must deploy before it. A kind with
// no install manifest contributes no edges here -- the missing manifest fails
// loudly later, when the dependency list is materialized, keeping that the
// single authoritative error. A manifest that cannot be parsed also
// contributes no edges: deploying it will fail with the real parse error, so
// swallowing it here never hides a failure. An annotation naming an unknown
// kind, however, errors immediately -- silently skipping it would deploy the
// manifest without a fixture it relies on.
func installManifestPrerequisites(repoRoot, componentProvider, consumer string, k cloudresourcekind.CloudResourceKind) ([]cloudresourcekind.CloudResourceKind, error) {
	slug := strings.ToLower(k.String())
	// The consumer scoping matters here too: annotations are read from the
	// manifest that will actually deploy, which may be a consumer-scoped
	// override rather than the kind's published profile.
	manifestPath, err := prerequisiteManifestPath(repoRoot, componentProvider, consumer, slug)
	if err != nil {
		return nil, nil
	}

	docPaths, err := splitManifestDocuments(manifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "splitting install profile for prerequisite %q", slug)
	}

	var kinds []cloudresourcekind.CloudResourceKind
	seen := make(map[cloudresourcekind.CloudResourceKind]bool)
	for _, docPath := range docPaths {
		raw, err := manifestAnnotation(docPath, scenarioPrerequisitesAnnotation)
		if err != nil {
			// Unparseable document: its deploy will surface the real error.
			continue
		}
		if raw == "" {
			continue
		}
		for _, token := range strings.Split(raw, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			if strings.Contains(token, "/") {
				// Manifest-path entries are scenario-scoped extra instances;
				// an install manifest's edges must be kind names so the graph
				// can dedup them against the rest of the chain.
				return nil, errors.Errorf("install manifest %s declares manifest-path entry %q in the %s annotation; install-manifest edges must be kind names (path entries are only valid on scenario manifests)", manifestPath, token, scenarioPrerequisitesAnnotation)
			}
			declared := crkreflect.KindFromString(token)
			if declared == cloudresourcekind.CloudResourceKind_unspecified {
				return nil, errors.Errorf("install manifest %s declares unknown kind %q in the %s annotation", manifestPath, token, scenarioPrerequisitesAnnotation)
			}
			if declared == k || seen[declared] {
				continue
			}
			seen[declared] = true
			kinds = append(kinds, declared)
		}
	}
	return kinds, nil
}

// prerequisiteManifestPath returns the manifest used to install a prerequisite,
// in order of preference (each component's version segment follows its
// declared version in the kind registry):
//   - <consumer>/<version>/e2e/prerequisites/<dep>.yaml — a consumer-scoped
//     override, for when the same prerequisite kind needs a different install
//     shape for different consumers (e.g. GcpGlobalAddress as an EXTERNAL VIP
//     for a forwarding rule vs an INTERNAL VPC_PEERING range for a service
//     networking connection);
//   - <dep>/<version>/e2e/prerequisite.yaml — the dependency's published install profile;
//   - <dep>/<version>/e2e/scenarios/minimal.yaml — fallback to its minimal scenario.
//
// Errors if none exist, so a missing install profile fails loudly rather than
// silently skipping a required dependency.
func prerequisiteManifestPath(repoRoot, componentProvider, consumer, slug string) (string, error) {
	if consumer != "" {
		consumerVersionDir, err := crkreflect.ComponentVersionDir(consumer)
		if err != nil {
			return "", err
		}
		consumerPrereq := filepath.Join(repoRoot, "catalog", componentProvider, consumer, consumerVersionDir, "e2e", "prerequisites", slug+".yaml")
		if pathExists(consumerPrereq) {
			return consumerPrereq, nil
		}
	}
	slugVersionDir, err := crkreflect.ComponentVersionDir(slug)
	if err != nil {
		return "", err
	}
	base := filepath.Join(repoRoot, "catalog", componentProvider, slug, slugVersionDir, "e2e")
	prereq := filepath.Join(base, "prerequisite.yaml")
	if pathExists(prereq) {
		return prereq, nil
	}
	minimal := filepath.Join(base, "scenarios", "minimal.yaml")
	if pathExists(minimal) {
		return minimal, nil
	}
	return "", errors.Errorf("no install manifest for prerequisite %q (consumer %q): expected %s or %s (or a consumer-scoped prerequisites/%s.yaml)", slug, consumer, prereq, minimal, slug)
}

// DeployDependencies resolves and deploys all prerequisite deployments for a
// component in order, via Pulumi. scenarioManifestPath (optional, "" to skip)
// lets the scenario under test add composed-but-optional prerequisites via its
// annotation. Returns the deployed states (needed for teardown) and any error.
// On the first failure it stops and returns whatever was already deployed so
// the caller can tear it down.
func DeployDependencies(ctx context.Context, repoRoot, componentProvider, component, scenarioManifestPath, backendURL, runID string, harness provider.Harness) ([]DependencyState, error) {
	deps, err := ResolveDependencies(repoRoot, componentProvider, component, scenarioManifestPath)
	if err != nil {
		return nil, err
	}
	if len(deps) == 0 {
		return nil, nil
	}

	fmt.Printf("  [deps] Deploying %d dependencies for %s\n", len(deps), component)

	// accumulated holds each deployed prerequisite's outputs keyed by kind and
	// manifest name. A later prerequisite that references an earlier one (e.g. an
	// AwsSubnet's vpc_id -> the AwsVpc it sits in) has its value_from refs resolved
	// against this map before it deploys -- the same resolution RunComponentTest
	// applies to the component under test, extended transitively across the
	// prerequisite chain so deep compositions (VPC -> Subnet -> NatGateway) can be
	// tested standalone.
	accumulated := make(DependencyOutputs, len(deps))

	var deployed []DependencyState
	for _, dep := range deps {
		docPaths, err := splitManifestDocuments(dep.ManifestPath)
		if err != nil {
			return deployed, errors.Wrapf(err, "failed to split install profile for dependency %q", dep.KindSlug)
		}
		for docIndex, docPath := range docPaths {
			docDep := Dependency{KindSlug: dep.KindSlug, ManifestPath: docPath}

			// Prerequisites redeploy for every engine run, so their manifests need
			// the same per-run unique-id expansion as the scenario under test (same
			// token, same run id — a dependent's ${E2E_RUN_ID}-suffixed reference
			// lines up with the prerequisite it points at).
			expandedManifestPath, err := ExpandManifestTokens(docDep.ManifestPath, runID)
			if err != nil {
				return deployed, errors.Wrapf(err, "failed to expand manifest tokens for dependency %q", dep.KindSlug)
			}
			docDep.ManifestPath = expandedManifestPath

			resolvedManifestPath, err := ResolveManifestRefs(docDep.ManifestPath, accumulated)
			if err != nil {
				return deployed, errors.Wrapf(err, "failed to resolve references for dependency %q", dep.KindSlug)
			}
			docDep.ManifestPath = resolvedManifestPath

			state, err := deployDependency(ctx, repoRoot, componentProvider, docDep, backendURL, runID, harness, docIndex)
			// A non-empty stack name means Pulumi created resources we must track
			// for teardown, even if verification afterwards failed.
			if state.StackName != "" {
				deployed = append(deployed, state)
			}
			if err != nil {
				return deployed, err
			}

			kind := crkreflect.KindFromString(dep.KindSlug)
			if accumulated[kind] == nil {
				accumulated[kind] = make(map[string]map[string]interface{})
			}
			accumulated[kind][state.ManifestName] = state.Outputs
		}
	}
	return deployed, nil
}

// deployDependency builds the stack input, runs `pulumi up`, and verifies the
// dependency is present. The dependency's own pulumi module is always used
// (dependencies deploy via Pulumi even when the component under test uses
// Terraform). docIndex disambiguates the stack name when an install profile
// deploys several instances of the same kind.
func deployDependency(ctx context.Context, repoRoot, componentProvider string, dep Dependency, backendURL, runID string, harness provider.Harness, docIndex int) (DependencyState, error) {
	versionDir, err := crkreflect.ComponentVersionDir(dep.KindSlug)
	if err != nil {
		return DependencyState{}, err
	}
	moduleDir := filepath.Join(repoRoot, "catalog", componentProvider, dep.KindSlug, versionDir, "iac", "pulumi")
	if !pathExists(moduleDir) {
		return DependencyState{}, errors.Errorf("dependency %q pulumi module not found at %s", dep.KindSlug, moduleDir)
	}

	manifestName, err := manifestMetadataName(dep.ManifestPath)
	if err != nil {
		return DependencyState{}, errors.Wrapf(err, "failed to read manifest name for dependency %q", dep.KindSlug)
	}

	// The stack label carries the MANIFEST name, not just the kind slug: a
	// scenario may chain several extra-instance fixtures of one kind (e.g. a
	// meshed client and backend, both KubernetesDeployment), and kind-only
	// labels would make the second fixture's `pulumi up` land on the first
	// fixture's stack — silently REPLACING it instead of adding an instance.
	stackLabel := "dep-" + dep.KindSlug + "-" + manifestName
	if docIndex > 0 {
		stackLabel = fmt.Sprintf("%s-%d", stackLabel, docIndex)
	}
	stackName := GenerateStackName(stackLabel, runID)
	if len(stackName) > 50 {
		stackName = stackName[:50]
	}

	fmt.Printf("  [deps] Deploying dependency %s (%s)...\n", dep.KindSlug, manifestName)
	start := time.Now()

	stackInputPath, err := BuildStackInput(dep.ManifestPath, moduleDir)
	if err != nil {
		return DependencyState{}, errors.Wrapf(err, "failed to build stack input for dependency %q", dep.KindSlug)
	}

	if _, err := PulumiDeploy(moduleDir, stackName, backendURL, stackInputPath); err != nil {
		return DependencyState{}, errors.Wrapf(err, "failed to deploy dependency %q", dep.KindSlug)
	}

	state := DependencyState{
		Dependency:     dep,
		ModuleDir:      moduleDir,
		StackName:      stackName,
		BackendURL:     backendURL,
		StackInputPath: stackInputPath,
		ManifestName:   manifestName,
	}

	// Capture the dependency's outputs so its verifier can confirm it (cloud
	// verifiers need the resource id from the outputs) and so the dependent
	// component's value_from refs can resolve against them.
	outputsJSON, err := PulumiStackOutputs(moduleDir, stackName, backendURL)
	if err != nil {
		return state, errors.Wrapf(err, "failed to read outputs for dependency %q", dep.KindSlug)
	}
	depStackOutputs, err := parsePulumiOutputs(outputsJSON)
	if err != nil {
		return state, errors.Wrapf(err, "failed to parse outputs for dependency %q", dep.KindSlug)
	}
	state.Outputs = depStackOutputs

	verifyCtx := context.WithValue(ctx, provider.ManifestPathKey{}, dep.ManifestPath)
	if err := harness.VerifyDeployed(verifyCtx, dep.KindSlug, state.Outputs); err != nil {
		return state, errors.Wrapf(err, "dependency %q deployed but verification failed", dep.KindSlug)
	}

	fmt.Printf("  [deps] Dependency %s deployed and verified in %s\n", dep.KindSlug, time.Since(start).Round(time.Second))
	return state, nil
}

// pulumiDestroyFn / pulumiRemoveStackFn are seams over the Pulumi CLI wrappers
// so teardown aggregation can be unit-tested without a live Pulumi backend.
// Production code never overrides them.
var (
	pulumiDestroyFn     = PulumiDestroy
	pulumiRemoveStackFn = PulumiRemoveStack
)

// Dependency destroys retry because some producer-side cleanups are
// asynchronous: deleting a Cloud SQL instance completes minutes before the
// service producer releases its hold on the service networking connection,
// and until then the connection destroy fails with "Producer services are
// still using this connection". Skipping past that failure is worse than an
// orphan: the framework would then force-delete the VPC, stranding a
// producer-side connection record that poisons the NEXT scenario's
// same-named prerequisite chain (its connection create silently attaches to
// the stale record and no peering ever materializes).
//
// The budget is deliberately bounded at ~6 minutes: async releases that GCP
// documents in HOURS (e.g. the serverless address reservation a direct-VPC
// Cloud Run service leaves in its subnetwork for 1-2 hours after deletion)
// can never be covered by retrying, and scenarios whose teardown depends on
// such a release must be excluded from live E2E instead (see the E2E
// coverage policy).
const dependencyDestroyAttempts = 6

// dependencyDestroyBackoff is a variable (not a const) purely so teardown
// unit tests can zero it; production code never changes it.
var dependencyDestroyBackoff = 60 * time.Second

// TeardownDependencies destroys deployed dependencies in reverse order and
// returns the aggregated failures. Each destroy retries with backoff (see
// dependencyDestroyAttempts) before being declared failed. Teardown is
// best-effort in the sense that one dependency's failure never stops the
// remaining teardowns (stopping early would leak everything deployed before
// it) -- but every failure is collected and returned so the caller FAILS the
// run. A destroy that cannot run (for example because the ephemeral backend
// state disappeared before teardown) means real cloud resources may still
// exist; reporting success would leave them leaking silently, invisible
// until someone audits the account.
func TeardownDependencies(deployed []DependencyState) error {
	var failures []error
	for i := len(deployed) - 1; i >= 0; i-- {
		dep := deployed[i]
		fmt.Printf("  [deps] Destroying dependency %s...\n", dep.Dependency.KindSlug)

		var destroyErr error
		for attempt := 1; attempt <= dependencyDestroyAttempts; attempt++ {
			if _, destroyErr = pulumiDestroyFn(dep.ModuleDir, dep.StackName, dep.BackendURL, dep.StackInputPath); destroyErr == nil {
				break
			}
			if attempt < dependencyDestroyAttempts {
				fmt.Printf("  [deps] dependency %s destroy attempt %d/%d failed (waiting %s): %v\n",
					dep.Dependency.KindSlug, attempt, dependencyDestroyAttempts, dependencyDestroyBackoff, destroyErr)
				time.Sleep(dependencyDestroyBackoff)
			}
		}
		if destroyErr != nil {
			failures = append(failures, errors.Wrapf(destroyErr, "dependency %s (stack %s) destroy failed after %d attempts", dep.Dependency.KindSlug, dep.StackName, dependencyDestroyAttempts))
			continue
		}
		if err := pulumiRemoveStackFn(dep.ModuleDir, dep.StackName, dep.BackendURL); err != nil {
			failures = append(failures, errors.Wrapf(err, "dependency %s (stack %s) stack removal failed", dep.Dependency.KindSlug, dep.StackName))
		}
	}
	return stderrors.Join(failures...)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// splitManifestDocuments splits a `---`-separated multi-document YAML file into
// one temp file per document, so each instance in a multi-instance install
// profile deploys as its own stack. A single-document file passes through
// untouched (the common case pays no temp-file cost). Empty documents (e.g. a
// leading separator or trailing whitespace) are skipped.
func splitManifestDocuments(manifestPath string) ([]string, error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read install profile %s", manifestPath)
	}

	var docs []string
	var current []string
	flush := func() {
		doc := strings.TrimSpace(strings.Join(current, "\n"))
		if doc != "" {
			docs = append(docs, strings.Join(current, "\n"))
		}
		current = nil
	}
	for _, line := range strings.Split(string(raw), "\n") {
		// A document separator is a line whose content is exactly "---"
		// (trailing whitespace tolerated); "--- something" would be YAML content.
		if strings.TrimRight(line, " \t") == "---" {
			flush()
			continue
		}
		current = append(current, line)
	}
	flush()

	if len(docs) <= 1 {
		return []string{manifestPath}, nil
	}

	paths := make([]string, 0, len(docs))
	for _, doc := range docs {
		tmpFile, err := os.CreateTemp("", "planton-e2e-prereq-doc-*.yaml")
		if err != nil {
			return nil, errors.Wrap(err, "failed to create temp file for install profile document")
		}
		if _, err := tmpFile.WriteString(doc + "\n"); err != nil {
			tmpFile.Close()
			return nil, errors.Wrap(err, "failed to write install profile document")
		}
		if err := tmpFile.Close(); err != nil {
			return nil, errors.Wrap(err, "failed to close install profile document")
		}
		paths = append(paths, tmpFile.Name())
	}
	return paths, nil
}

// scenarioDeclaredPrerequisites reads the scenario manifest's
// `planton.dev/e2e-prerequisites` annotation and returns the declared kind
// entries and manifest-path entries (recognized by containing a "/") in
// their listed order. An empty path or an absent annotation returns nil
// (the registry graph alone drives resolution); an unknown kind name errors
// loudly rather than silently skipping a dependency the scenario relies on.
func scenarioDeclaredPrerequisites(manifestPath string) ([]cloudresourcekind.CloudResourceKind, []string, error) {
	if manifestPath == "" {
		return nil, nil, nil
	}
	raw, err := manifestAnnotation(manifestPath, scenarioPrerequisitesAnnotation)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "reading scenario manifest %s", manifestPath)
	}
	if raw == "" {
		return nil, nil, nil
	}

	var kinds []cloudresourcekind.CloudResourceKind
	var paths []string
	for _, token := range strings.Split(raw, ",") {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		if strings.Contains(token, "/") {
			paths = append(paths, token)
			continue
		}
		kind := crkreflect.KindFromString(token)
		if kind == cloudresourcekind.CloudResourceKind_unspecified {
			return nil, nil, errors.Errorf("scenario %s declares unknown kind %q in the %s annotation", manifestPath, token, scenarioPrerequisitesAnnotation)
		}
		kinds = append(kinds, kind)
	}
	return kinds, paths, nil
}

// manifestKindSlug reads the kind a manifest declares and returns its lowercase
// component directory slug, erroring on unregistered kinds so a mistyped
// manifest-path fixture fails at resolution time rather than mid-deploy.
func manifestKindSlug(manifestPath string) (string, error) {
	obj, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return "", err
	}
	name := string(obj.ProtoReflect().Descriptor().Name())
	kind := crkreflect.KindFromString(name)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return "", errors.Errorf("manifest %s declares kind %q, which is not a registered cloud resource kind", manifestPath, name)
	}
	return strings.ToLower(kind.String()), nil
}

// ManifestAnnotation reads a single metadata annotation from a KRM manifest,
// returning "" when the annotation (or the annotations map) is absent. The
// exported form serves the test entrypoints' scenario-routing needs (e.g.
// per-scenario cluster-profile selection) without duplicating manifest
// parsing outside the framework.
func ManifestAnnotation(manifestPath, key string) (string, error) {
	return manifestAnnotation(manifestPath, key)
}

// ScenarioEnginesAnnotation restricts a scenario to a subset of IaC engines
// (comma-separated: "pulumi", "terraform"). Absent = both engines, which is
// the norm and the parity bar. Use it ONLY when a scenario exercises a
// surface one engine rejects BY DOCUMENTED DESIGN — e.g. a spec arm the
// Terraform provider cannot express, guarded by a PARITY-EXCEPTION
// precondition: the scenario proving that arm live is honest about which
// engine can run it, and the other engine's lane skips with the reason
// instead of failing on its own designed rejection.
const ScenarioEnginesAnnotation = "planton.dev/e2e-engines"

// ScenarioSupportsEngine reports whether a scenario runs on the given engine
// per its ScenarioEnginesAnnotation (absent = every engine).
func ScenarioSupportsEngine(manifestPath, engine string) (bool, error) {
	value, err := manifestAnnotation(manifestPath, ScenarioEnginesAnnotation)
	if err != nil {
		return false, err
	}
	if value == "" {
		return true, nil
	}
	for _, e := range strings.Split(value, ",") {
		if strings.TrimSpace(e) == engine {
			return true, nil
		}
	}
	return false, nil
}

// ScenarioRequiredEnvAnnotation names the PLANTON_E2E_* environment
// variables (comma-separated) a scenario cannot run without — external
// credentials a committed manifest carries only as ${E2E_ENV:...} tokens
// (a sandbox account PAT, an owner-arranged tenant). Unset token
// variables FAIL expansion loudly by design (the right behavior when a
// batch bootstrap was SUPPOSED to export them), so a scenario whose
// tokens are owner-arranged must declare them here: lanes then SKIP it
// with the reason when the environment does not carry the arrangement —
// an honest deferral instead of a false failure — and run it live
// wherever the variables are exported.
const ScenarioRequiredEnvAnnotation = "planton.dev/e2e-required-env"

// ScenarioMissingRequiredEnv returns the declared-but-unset environment
// variable names per ScenarioRequiredEnvAnnotation (absent = none).
func ScenarioMissingRequiredEnv(manifestPath string) ([]string, error) {
	value, err := manifestAnnotation(manifestPath, ScenarioRequiredEnvAnnotation)
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}
	var missing []string
	for _, name := range strings.Split(value, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if os.Getenv(name) == "" {
			missing = append(missing, name)
		}
	}
	return missing, nil
}

// manifestAnnotation reads a single metadata annotation from a KRM manifest,
// returning "" when the annotation (or the annotations map) is absent.
func manifestAnnotation(manifestPath, key string) (string, error) {
	obj, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return "", err
	}
	top := obj.ProtoReflect()
	metaFd := top.Descriptor().Fields().ByName("metadata")
	if metaFd == nil || metaFd.Kind() != protoreflect.MessageKind {
		return "", nil
	}
	annFd := metaFd.Message().Fields().ByName("annotations")
	if annFd == nil || !annFd.IsMap() {
		return "", nil
	}
	value := top.Get(metaFd).Message().Get(annFd).Map().Get(protoreflect.ValueOfString(key).MapKey())
	if !value.IsValid() {
		return "", nil
	}
	return value.String(), nil
}

// manifestMetadataName reads metadata.name from a KRM manifest. The name keys
// the deployed instance's outputs in DependencyOutputs, which is what lets a
// reference address one specific instance of a multi-instance prerequisite.
func manifestMetadataName(manifestPath string) (string, error) {
	obj, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return "", err
	}
	top := obj.ProtoReflect()
	metaFd := top.Descriptor().Fields().ByName("metadata")
	if metaFd == nil || metaFd.Kind() != protoreflect.MessageKind {
		return "", errors.Errorf("manifest %s has no metadata message", manifestPath)
	}
	nameFd := metaFd.Message().Fields().ByName("name")
	if nameFd == nil {
		return "", errors.Errorf("manifest %s metadata has no name field", manifestPath)
	}
	name := top.Get(metaFd).Message().Get(nameFd).String()
	if name == "" {
		return "", errors.Errorf("manifest %s has an empty metadata.name", manifestPath)
	}
	return name, nil
}
