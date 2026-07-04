package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ExtraPrerequisitesAnnotation is the scenario-manifest annotation through which
// a scenario declares fixtures BEYOND its kind's registry prerequisites. Registry
// `prerequisites` carry a strict meaning -- the parents a resource cannot exist
// without -- and double as deploy-ordering metadata, so optional composition
// seams (a subnet's route-table attachment, a NAT gateway's public-IP
// association, a peering's second network) must never be encoded there: doing so
// would force every downstream kind's fixture chain to deploy them forever.
// This annotation is the honest home for such seams: the scenario that proves an
// optional edge opts into the fixtures it needs, and no other chain pays for it.
//
// The value is a comma-separated list where each entry is either:
//   - a registered kind name (e.g. "AzureRouteTable") -- installed through the
//     kind's standard install profile (prerequisite.yaml, else minimal scenario),
//     skipped if the registry chain already deploys that kind; or
//   - a repo-relative manifest path (recognized by containing a "/") -- deployed
//     as an EXTRA INSTANCE of whatever kind the manifest declares, for scenarios
//     that need more instances than the standard profiles provide (e.g. a
//     virtual-network peering's remote network).
//
// Entries deploy in listed order after the registry prerequisites, each preceded
// by any of its own transitive registry prerequisites not already deployed. All
// fixtures join the same transitive reference resolution, and teardown runs in
// reverse across the merged chain.
const ExtraPrerequisitesAnnotation = "planton.dev/e2e-extra-prerequisites"

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

// ResolveDependencies returns the ordered, deduplicated list of prerequisite
// deployments a component needs before its own scenario is applied. Dependencies
// come from the component's CloudResourceKindMeta.prerequisites graph in the proto
// registry (resolved transitively, deploy-first order): declaring
// `prerequisites: [X]` on a kind is enough for the harness to install X first,
// with no per-component wiring. Teardown runs in reverse, so the most foundational
// dependency is removed last.
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
func ResolveDependencies(repoRoot, componentProvider, component string) ([]Dependency, error) {
	kind := crkreflect.KindFromString(component)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		// Not a registered kind (or an alias mismatch); no prerequisites.
		return nil, nil
	}

	prereqs, err := crkreflect.TransitivePrerequisites(kind)
	if err != nil {
		return nil, errors.Wrapf(err, "resolving prerequisites for %s", component)
	}

	var deps []Dependency
	for _, p := range prereqs {
		slug := strings.ToLower(p.String())
		manifestPath, err := prerequisiteManifestPath(repoRoot, componentProvider, slug)
		if err != nil {
			return nil, err
		}
		deps = append(deps, Dependency{
			KindSlug:     slug,
			ManifestPath: manifestPath,
		})
	}
	return deps, nil
}

// ResolveAllDependencies returns the registry-driven prerequisites followed by
// any extra fixtures the scenario manifest declares through
// ExtraPrerequisitesAnnotation, deduplicated and in deploy-first order.
func ResolveAllDependencies(repoRoot, componentProvider, component, scenarioManifestPath string) ([]Dependency, error) {
	deps, err := ResolveDependencies(repoRoot, componentProvider, component)
	if err != nil {
		return nil, err
	}
	extras, err := resolveExtraDependencies(repoRoot, componentProvider, scenarioManifestPath, deps)
	if err != nil {
		return nil, err
	}
	return append(deps, extras...), nil
}

// resolveExtraDependencies reads ExtraPrerequisitesAnnotation from the scenario
// manifest and expands it into deployable dependencies. Kind-name entries are
// skipped when the registry chain already deploys that kind; path entries always
// deploy (they exist precisely to add another instance). Every entry is preceded
// by its own transitive registry prerequisites that are not yet scheduled, so an
// extra fixture with parents of its own (e.g. a NAT gateway needing a resource
// group) works without the scenario spelling the whole chain out.
func resolveExtraDependencies(repoRoot, componentProvider, scenarioManifestPath string, registry []Dependency) ([]Dependency, error) {
	if scenarioManifestPath == "" {
		return nil, nil
	}
	raw, err := manifestAnnotation(scenarioManifestPath, ExtraPrerequisitesAnnotation)
	if err != nil {
		return nil, errors.Wrapf(err, "reading %s annotation from %s", ExtraPrerequisitesAnnotation, scenarioManifestPath)
	}
	if raw == "" {
		return nil, nil
	}

	scheduled := make(map[string]bool, len(registry))
	for _, d := range registry {
		scheduled[d.KindSlug] = true
	}

	var extras []Dependency
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		if strings.Contains(entry, "/") {
			// A repo-relative manifest path: an extra instance of its declared kind.
			full := filepath.Join(repoRoot, entry)
			if !pathExists(full) {
				return nil, errors.Errorf("%s annotation entry %q: manifest not found at %s", ExtraPrerequisitesAnnotation, entry, full)
			}
			slug, err := manifestKindSlug(full)
			if err != nil {
				return nil, errors.Wrapf(err, "%s annotation entry %q", ExtraPrerequisitesAnnotation, entry)
			}
			pre, err := unscheduledTransitives(repoRoot, componentProvider, slug, scheduled)
			if err != nil {
				return nil, err
			}
			extras = append(extras, pre...)
			// The path instance itself never marks its kind as scheduled: it is an
			// additional instance, not a substitute for the kind's install profile.
			extras = append(extras, Dependency{KindSlug: slug, ManifestPath: full})
			continue
		}

		// A registered kind name, installed through its standard profile.
		kind := crkreflect.KindFromString(entry)
		if kind == cloudresourcekind.CloudResourceKind_unspecified {
			return nil, errors.Errorf("%s annotation entry %q is neither a registered kind nor a manifest path", ExtraPrerequisitesAnnotation, entry)
		}
		slug := strings.ToLower(kind.String())
		if scheduled[slug] {
			continue
		}
		pre, err := unscheduledTransitives(repoRoot, componentProvider, slug, scheduled)
		if err != nil {
			return nil, err
		}
		extras = append(extras, pre...)
		manifestPath, err := prerequisiteManifestPath(repoRoot, componentProvider, slug)
		if err != nil {
			return nil, err
		}
		extras = append(extras, Dependency{KindSlug: slug, ManifestPath: manifestPath})
		scheduled[slug] = true
	}
	return extras, nil
}

// unscheduledTransitives returns the transitive registry prerequisites of a kind
// that are not yet scheduled, in deploy-first order, marking them scheduled.
func unscheduledTransitives(repoRoot, componentProvider, slug string, scheduled map[string]bool) ([]Dependency, error) {
	kind := crkreflect.KindFromString(slug)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return nil, errors.Errorf("kind %q is not registered", slug)
	}
	prereqs, err := crkreflect.TransitivePrerequisites(kind)
	if err != nil {
		return nil, errors.Wrapf(err, "resolving prerequisites for extra fixture %s", slug)
	}
	var deps []Dependency
	for _, p := range prereqs {
		pSlug := strings.ToLower(p.String())
		if scheduled[pSlug] {
			continue
		}
		manifestPath, err := prerequisiteManifestPath(repoRoot, componentProvider, pSlug)
		if err != nil {
			return nil, err
		}
		deps = append(deps, Dependency{KindSlug: pSlug, ManifestPath: manifestPath})
		scheduled[pSlug] = true
	}
	return deps, nil
}

// prerequisiteManifestPath returns the manifest used to install a prerequisite:
// its published prerequisite.yaml if present, else its minimal scenario. Errors if
// neither exists, so a missing install profile fails loudly rather than silently
// skipping a required dependency.
func prerequisiteManifestPath(repoRoot, componentProvider, slug string) (string, error) {
	base := filepath.Join(repoRoot, "apis", "dev", "planton", "provider", componentProvider, slug, "v1", "e2e")
	prereq := filepath.Join(base, "prerequisite.yaml")
	if pathExists(prereq) {
		return prereq, nil
	}
	minimal := filepath.Join(base, "scenarios", "minimal.yaml")
	if pathExists(minimal) {
		return minimal, nil
	}
	return "", errors.Errorf("no install manifest for prerequisite %q: expected %s or %s", slug, prereq, minimal)
}

// DeployDependencies resolves and deploys all prerequisite deployments for a
// component in order, via Pulumi: the registry prerequisites first, then any
// extra fixtures the scenario manifest declares (ExtraPrerequisitesAnnotation).
// Returns the deployed states (needed for teardown) and any error. On the first
// failure it stops and returns whatever was already deployed so the caller can
// tear it down.
func DeployDependencies(ctx context.Context, repoRoot, componentProvider, component, scenarioManifestPath, backendURL, runID string, harness provider.Harness) ([]DependencyState, error) {
	deps, err := ResolveAllDependencies(repoRoot, componentProvider, component, scenarioManifestPath)
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

	// instanceCounts disambiguates stack names when the same kind deploys more
	// than once in a chain -- whether as several documents of one install profile
	// or as a scenario-declared extra instance of an already-deployed kind.
	instanceCounts := make(map[string]int, len(deps))

	var deployed []DependencyState
	for _, dep := range deps {
		docPaths, err := splitManifestDocuments(dep.ManifestPath)
		if err != nil {
			return deployed, errors.Wrapf(err, "failed to split install profile for dependency %q", dep.KindSlug)
		}
		for _, docPath := range docPaths {
			docDep := Dependency{KindSlug: dep.KindSlug, ManifestPath: docPath}

			resolvedManifestPath, err := ResolveManifestRefs(docDep.ManifestPath, accumulated)
			if err != nil {
				return deployed, errors.Wrapf(err, "failed to resolve references for dependency %q", dep.KindSlug)
			}
			docDep.ManifestPath = resolvedManifestPath

			instanceIndex := instanceCounts[dep.KindSlug]
			instanceCounts[dep.KindSlug]++

			state, err := deployDependency(ctx, repoRoot, componentProvider, docDep, backendURL, runID, harness, instanceIndex)
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
// Terraform). instanceIndex disambiguates the stack name when the chain deploys
// several instances of the same kind.
func deployDependency(ctx context.Context, repoRoot, componentProvider string, dep Dependency, backendURL, runID string, harness provider.Harness, instanceIndex int) (DependencyState, error) {
	moduleDir := filepath.Join(repoRoot, "apis", "dev", "planton", "provider", componentProvider, dep.KindSlug, "v1", "iac", "pulumi")
	if !pathExists(moduleDir) {
		return DependencyState{}, errors.Errorf("dependency %q pulumi module not found at %s", dep.KindSlug, moduleDir)
	}

	stackLabel := "dep-" + dep.KindSlug
	if instanceIndex > 0 {
		stackLabel = fmt.Sprintf("%s-%d", stackLabel, instanceIndex)
	}
	stackName := GenerateStackName(stackLabel, runID)
	if len(stackName) > 50 {
		stackName = stackName[:50]
	}

	fmt.Printf("  [deps] Deploying dependency %s...\n", dep.KindSlug)
	start := time.Now()

	manifestName, err := manifestMetadataName(dep.ManifestPath)
	if err != nil {
		return DependencyState{}, errors.Wrapf(err, "failed to read manifest name for dependency %q", dep.KindSlug)
	}

	stackInputPath, err := BuildStackInput(dep.ManifestPath, moduleDir)
	if err != nil {
		return DependencyState{}, errors.Wrapf(err, "failed to build stack input for dependency %q", dep.KindSlug)
	}

	// The state is constructed BEFORE `pulumi up` runs: a failed deploy may
	// still have created any number of resources (the stack exists the moment
	// the update starts), and returning an empty state here would hide them
	// from teardown -- the partially-created fixture then orphans and blocks
	// the whole reverse teardown chain (a leftover load balancer holds its
	// subnet, the subnet holds the network, the network holds the resource
	// group). Destroying a stack whose update failed is safe: destroy removes
	// whatever was actually created.
	state := DependencyState{
		Dependency:     dep,
		ModuleDir:      moduleDir,
		StackName:      stackName,
		BackendURL:     backendURL,
		StackInputPath: stackInputPath,
		ManifestName:   manifestName,
	}

	if _, err := PulumiDeploy(moduleDir, stackName, backendURL, stackInputPath); err != nil {
		return state, errors.Wrapf(err, "failed to deploy dependency %q", dep.KindSlug)
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

// TeardownDependencies destroys deployed dependencies in reverse order.
func TeardownDependencies(deployed []DependencyState) {
	for i := len(deployed) - 1; i >= 0; i-- {
		dep := deployed[i]
		fmt.Printf("  [deps] Destroying dependency %s...\n", dep.Dependency.KindSlug)

		if _, err := PulumiDestroy(dep.ModuleDir, dep.StackName, dep.BackendURL, dep.StackInputPath); err != nil {
			fmt.Printf("  [WARN] dependency %s destroy failed: %v\n", dep.Dependency.KindSlug, err)
			continue
		}
		if err := PulumiRemoveStack(dep.ModuleDir, dep.StackName, dep.BackendURL); err != nil {
			fmt.Printf("  [WARN] dependency %s stack removal failed: %v\n", dep.Dependency.KindSlug, err)
		}
	}
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

// manifestAnnotation reads a single metadata.annotations value from a KRM
// manifest, returning "" when the annotation (or the annotations map) is absent.
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
	meta := top.Get(metaFd).Message()
	annFd := meta.Descriptor().Fields().ByName("annotations")
	if annFd == nil || !annFd.IsMap() {
		return "", nil
	}
	value := meta.Get(annFd).Map().Get(protoreflect.ValueOfString(key).MapKey())
	if !value.IsValid() {
		return "", nil
	}
	return value.String(), nil
}

// manifestKindSlug derives the lowercase kind slug from a manifest's declared
// kind, erroring on unregistered kinds so a typo in a path-declared extra
// fixture fails at resolution rather than at module lookup.
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
