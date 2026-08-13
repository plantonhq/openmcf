package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// This file is the OFFLINE half of the reference-resolution contract that
// refresolve.go and dependencies.go execute at deploy time. A value_from
// reference that cannot resolve against the prerequisite chain fails only
// DEEP into a live run -- after minutes of fixture deploys (an ambiguous
// name), or worse, silently: an unresolvable reference is left untouched by
// design, the module receives no value, and the failure surfaces as a
// provider validation error at DEPLOY that reads like a module defect (see
// e2e/README.md, "Bare polymorphic references"). None of the other offline
// gates see this class: the manifest validates and the plan renders without
// the reference ever being followed.
//
// CheckScenarioFixtureIntegrity therefore replays, statically, exactly what
// DeployDependencies + ResolveManifestRefs will do: resolve the chain, walk
// each install manifest's references against the instances deployed BEFORE
// it, then walk the scenario's references against the full chain. The rules
// mirror lookupRefValue one for one -- including the sole-instance fallback
// (a name that matches no deployed instance still resolves when exactly one
// instance of the kind exists; scenario references name real-world topology,
// not install-profile files, so that is legal by design and NOT a finding).

// FixtureIntegrityFinding describes one reference (or chain-resolution
// failure) that would break a live run. ManifestPath is the manifest carrying
// the problem -- the scenario itself or one of its chain's install manifests.
type FixtureIntegrityFinding struct {
	ScenarioPath string
	ManifestPath string
	Field        string
	RefKind      string
	RefName      string
	Reason       string
}

func (f FixtureIntegrityFinding) String() string {
	ref := ""
	if f.RefKind != "" || f.RefName != "" {
		ref = fmt.Sprintf(" ref %s/%q on field %q:", f.RefKind, f.RefName, f.Field)
	}
	manifest := f.ManifestPath
	if manifest == f.ScenarioPath {
		manifest = "scenario"
	}
	return fmt.Sprintf("%s [%s]%s %s", f.ScenarioPath, manifest, ref, f.Reason)
}

// CheckScenarioFixtureIntegrity statically verifies that every value_from
// reference in the scenario -- and in every install manifest of its resolved
// prerequisite chain -- will resolve when the runner deploys that chain.
// Findings are defects to fix in the manifests (or the kind's registry
// prerequisites); the returned error is reserved for I/O-level failures of
// the checker itself.
func CheckScenarioFixtureIntegrity(repoRoot, componentProvider, component, scenarioPath string) ([]FixtureIntegrityFinding, error) {
	deps, err := ResolveDependencies(repoRoot, componentProvider, component, scenarioPath)
	if err != nil {
		// A chain that cannot even resolve (missing install manifest, unknown
		// annotation kind, cycle) is the first thing a live run would die on.
		return []FixtureIntegrityFinding{{
			ScenarioPath: scenarioPath,
			ManifestPath: scenarioPath,
			Reason:       "prerequisite chain resolution failed: " + err.Error(),
		}}, nil
	}

	// deployed accumulates (kind -> metadata.name set) in deploy order,
	// mirroring the `accumulated` outputs map in DeployDependencies: an
	// install manifest's references resolve only against instances deployed
	// before it; the scenario's resolve against the whole chain.
	deployed := make(map[cloudresourcekind.CloudResourceKind]map[string]bool)
	var findings []FixtureIntegrityFinding

	for _, dep := range deps {
		docPaths, err := splitManifestDocuments(dep.ManifestPath)
		if err != nil {
			findings = append(findings, FixtureIntegrityFinding{
				ScenarioPath: scenarioPath,
				ManifestPath: dep.ManifestPath,
				Reason:       "install profile cannot be split into documents: " + err.Error(),
			})
			continue
		}
		kind := crkreflect.KindFromString(dep.KindSlug)
		for _, docPath := range docPaths {
			findings = append(findings, manifestRefFindings(scenarioPath, dep.ManifestPath, docPath, deployed)...)

			name, err := manifestMetadataName(docPath)
			if err != nil {
				findings = append(findings, FixtureIntegrityFinding{
					ScenarioPath: scenarioPath,
					ManifestPath: dep.ManifestPath,
					Reason:       "install manifest has no readable metadata.name: " + err.Error(),
				})
				continue
			}
			if deployed[kind] == nil {
				deployed[kind] = make(map[string]bool)
			}
			deployed[kind][name] = true
		}
	}

	findings = append(findings, manifestRefFindings(scenarioPath, scenarioPath, scenarioPath, deployed)...)
	return findings, nil
}

// manifestRefFindings loads one manifest document and checks every value_from
// reference in its spec against the instances deployed so far. reportPath is
// the on-disk file to name in findings (docPath may be a temp per-document
// split of a multi-document profile).
func manifestRefFindings(scenarioPath, reportPath, docPath string, deployed map[cloudresourcekind.CloudResourceKind]map[string]bool) []FixtureIntegrityFinding {
	manifestObject, err := manifest.LoadManifest(docPath)
	if err != nil {
		return []FixtureIntegrityFinding{{
			ScenarioPath: scenarioPath,
			ManifestPath: reportPath,
			Reason:       "manifest failed to load: " + err.Error(),
		}}
	}

	top := manifestObject.ProtoReflect()
	specFd := top.Descriptor().Fields().ByName("spec")
	if specFd == nil || specFd.Kind() != protoreflect.MessageKind {
		return nil
	}

	var findings []FixtureIntegrityFinding
	// The visitor never replaces anything -- this is a read-only replay of the
	// deploy-time resolution through the same traversal it uses.
	_, walkErr := forEachRefField(top.Mutable(specFd).Message(), func(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef) (*foreignkeyv1.StringValueOrRef, error) {
		if ref.GetValueFrom() == nil {
			return nil, nil
		}
		if finding := checkRefResolvable(fd, ref, deployed); finding != nil {
			finding.ScenarioPath = scenarioPath
			finding.ManifestPath = reportPath
			findings = append(findings, *finding)
		}
		return nil, nil
	})
	if walkErr != nil {
		findings = append(findings, FixtureIntegrityFinding{
			ScenarioPath: scenarioPath,
			ManifestPath: reportPath,
			Reason:       "reference walk failed: " + walkErr.Error(),
		})
	}
	return findings
}

// checkRefResolvable mirrors lookupRefValue's resolution rules against the
// statically known instance names instead of live outputs. Returns nil when
// the reference will resolve.
func checkRefResolvable(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef, deployed map[cloudresourcekind.CloudResourceKind]map[string]bool) *FixtureIntegrityFinding {
	valueFrom := ref.GetValueFrom()
	kind := referencedKind(fd, ref)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return &FixtureIntegrityFinding{
			Field:   string(fd.Name()),
			RefName: valueFrom.GetName(),
			Reason: "reference carries no kind and the field declares no default_kind -- the runner leaves it " +
				"unresolved and the module deploys without the value; add an explicit `kind:` to the valueFrom",
		}
	}

	instances := deployed[kind]
	if len(instances) == 0 {
		return &FixtureIntegrityFinding{
			Field:   string(fd.Name()),
			RefKind: kind.String(),
			RefName: valueFrom.GetName(),
			Reason: "the prerequisite chain deploys no instance of this kind before this manifest -- the reference " +
				"can never resolve; add the kind to the component's registry prerequisites or the scenario's " +
				"planton.dev/e2e-prerequisites annotation",
		}
	}
	if instances[valueFrom.GetName()] || len(instances) == 1 {
		// Exact name match, or the runner's sole-instance fallback.
		return nil
	}
	names := make([]string, 0, len(instances))
	for n := range instances {
		names = append(names, n)
	}
	sort.Strings(names)
	return &FixtureIntegrityFinding{
		Field:   string(fd.Name()),
		RefKind: kind.String(),
		RefName: valueFrom.GetName(),
		Reason: fmt.Sprintf("name matches none of the %d deployed instances (%s) -- the runner fails the run on "+
			"this ambiguity after the fixtures are already deployed", len(instances), strings.Join(names, ", ")),
	}
}

// CheckCatalogFixtureIntegrity runs the scenario check across every component
// scenario in the repository's catalog and returns all findings, keeping one
// component's defect from hiding another's.
func CheckCatalogFixtureIntegrity(repoRoot string) ([]FixtureIntegrityFinding, error) {
	catalogDir := filepath.Join(repoRoot, "catalog")
	providers, err := os.ReadDir(catalogDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read catalog dir %s", catalogDir)
	}

	var findings []FixtureIntegrityFinding
	for _, provider := range providers {
		if !provider.IsDir() {
			continue
		}
		providerDir := filepath.Join(catalogDir, provider.Name())
		components, err := os.ReadDir(providerDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read provider dir %s", providerDir)
		}
		for _, component := range components {
			if !component.IsDir() {
				continue
			}
			scenariosDir := filepath.Join(providerDir, component.Name(), "e2e", "scenarios")
			scenarios, err := os.ReadDir(scenariosDir)
			if err != nil {
				continue // no scenarios -- nothing to check
			}
			for _, scenario := range scenarios {
				if scenario.IsDir() || !strings.HasSuffix(scenario.Name(), ".yaml") {
					continue
				}
				scenarioPath := filepath.Join(scenariosDir, scenario.Name())
				scenarioFindings, err := CheckScenarioFixtureIntegrity(repoRoot, provider.Name(), component.Name(), scenarioPath)
				if err != nil {
					return nil, errors.Wrapf(err, "fixture integrity check failed for %s", scenarioPath)
				}
				findings = append(findings, scenarioFindings...)
			}
		}
	}
	return findings, nil
}
