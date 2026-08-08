//go:build !codegen
// +build !codegen

// The total-accounting check: every configurable, non-deprecated argument of
// every consumed provider resource is exact-matched to a spec field, mapped
// by recorded judgment, or excluded with a reason -- and, in reverse, every
// spec leaf field lands on a provider argument or carries an exclusion. The
// reverse direction is not symmetry for its own sake: specs are authored
// against a moving provider, and a field the provider no longer serves is
// drift the forward walk alone can never see.
//
// Anything unaccounted in either direction is a Finding, and Findings gate
// through the burn-down baseline (baseline.go) in the pkg/anatomy /
// pkg/secretcoverage grain: new findings fail, fixed-but-still-listed
// baseline entries fail, and the ledger only ever shrinks truthfully.

package providerparity

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// machineryArgPaths are Terraform's own operational surface, present on
// essentially every resource and deliberately outside parity accounting: a
// spec models the RESOURCE, not the tool driving it. One standing exclusion
// here beats the same boilerplate copied into every kind's manifest -- and
// keeps the parity denominator principled (resource configuration only).
//
//   - "id": the Terraform meta-argument legacy SDK resources expose as an
//     optional attribute.
//   - "timeouts.*": per-operation apply/destroy deadlines, an engine
//     concern the platform owns globally.
func isMachineryArg(path string) bool {
	return path == "id" || path == "timeouts" || strings.HasPrefix(path, "timeouts.")
}

// iamResourceRe classifies the provider's per-resource IAM triplets
// (*_iam_member / *_iam_binding / *_iam_policy). The catalog's standing
// design covers this class with additive iam_members fields on the owning
// kinds (authoritative binding/policy forms are deliberately not modeled),
// so the class is dispositioned by pattern, never by 300+ ledger entries.
var iamResourceRe = regexp.MustCompile(`_iam_(member|binding|policy)$`)

// Disposition names for the breadth accounting. Every GA resource carries
// exactly one.
const (
	// DispositionModeled: consumed by at least one kind's Terraform module.
	DispositionModeled = "modeled"
	// DispositionIamCovered: an *_iam_* triplet resource, covered by the
	// owning kind's additive iam_members field.
	DispositionIamCovered = "iam-covered"
	// DispositionExcludedDeprecated: deprecated surface, excluded from
	// parity (schema-flagged automatically; doc-level deprecations enter
	// via the ledger).
	DispositionExcludedDeprecated = "excluded-deprecated"
	// DispositionComposed: covered by fields of an existing kind rather
	// than a kind of its own (ledger-recorded, with the covering kind
	// named in the reason).
	DispositionComposed = "composed"
	// DispositionModelPlanned: judged to be covered by a planned kind that
	// is not built yet -- a kind of its own, or a planned kind it composes
	// into (ledger-recorded, with the planned kind named in the reason).
	// When the kind ships, the entry flips to modeled (computed) or
	// composed, and the staleness findings force the cleanup.
	DispositionModelPlanned = "model-planned"
	// DispositionDeferred: deliberately not offered, with the reason
	// recorded (ledger-recorded).
	DispositionDeferred = "deferred"
)

// ledgerDispositions are the judgments the ledger may record. The other two
// classes (modeled, iam-covered) are computed, never hand-written.
var ledgerDispositions = map[string]bool{
	DispositionComposed:           true,
	DispositionModelPlanned:       true,
	DispositionDeferred:           true,
	DispositionExcludedDeprecated: true,
}

// Finding is one unaccounted or stale item. Findings group into the
// burn-down baseline by BaselineKey, so the baseline reads as a work list
// (one line per kind or resource), not a field dump.
type Finding struct {
	// BaselineKey is "kind:<Kind>" (depth accounting) or
	// "resource:<name>" (breadth disposition).
	BaselineKey string `json:"baselineKey"`
	// Detail names the exact gap and, where possible, the fix.
	Detail string `json:"detail"`
}

// KindAccounting is one kind's depth accounting against the pinned provider.
type KindAccounting struct {
	Kind        string `json:"kind"`
	HasManifest bool   `json:"hasManifest"`
	// MissingModule marks a kind with no Terraform module directory: its
	// depth cannot be measured at all, which is itself the finding -- the
	// per-argument walk is skipped so the one real gap is not buried under
	// a spec-field noise storm.
	MissingModule bool `json:"missingModule,omitempty"`
	// TotalArgs is the surface under accounting: non-deprecated,
	// non-machinery configurable arguments across the kind's consumed
	// resources (internal-dispositioned resources excluded).
	TotalArgs int `json:"totalArgs"`
	// MatchedArgs matched their derived spec path exactly.
	MatchedArgs int `json:"matchedArgs"`
	// MappedArgs matched through a recorded mapping.
	MappedArgs int `json:"mappedArgs"`
	// ExcludedArgs carry a recorded exclusion.
	ExcludedArgs int `json:"excludedArgs"`
	// InternalResources are consumed resources dispositioned as module
	// plumbing by the manifest.
	InternalResources []string `json:"internalResources,omitempty"`
	// UnaccountedArgs ("resource: arg") have no match, mapping, or
	// exclusion -- each is a Finding.
	UnaccountedArgs []string `json:"unaccountedArgs,omitempty"`
	// UncoveredSpecFields are spec leaves no argument matched and no
	// exclusion covers -- reverse drift candidates. Each is a Finding.
	UncoveredSpecFields []string `json:"uncoveredSpecFields,omitempty"`
	// ManifestStale are manifest entries referencing surface that no
	// longer exists (after a pin bump or spec change). Each is a Finding.
	ManifestStale []string `json:"manifestStale,omitempty"`
}

// Accounted reports whether the kind is at total accounting.
func (k KindAccounting) Accounted() bool {
	return !k.MissingModule &&
		len(k.UnaccountedArgs) == 0 && len(k.UncoveredSpecFields) == 0 && len(k.ManifestStale) == 0
}

// ResourceDisposition is one GA resource's recorded breadth judgment.
type ResourceDisposition struct {
	Resource    string `json:"resource"`
	Disposition string `json:"disposition,omitempty"`
	// Detail carries the consuming kinds (modeled), the recorded reason
	// (ledger), or the classifying rule (pattern/schema classes).
	Detail string `json:"detail,omitempty"`
}

// Accounting is the total-accounting result for one cloud provider's
// catalog: per-kind depth accounting plus the GA breadth disposition, with
// every gap surfaced as a Finding.
type Accounting struct {
	CloudProvider string `json:"cloudProvider"`
	// GASchema/GASchemaVersion name the parity baseline the accounting ran
	// against (google-beta capability enters per kind through the
	// enumerated admission list, never through this accounting).
	GASchema        string `json:"gaSchema"`
	GASchemaVersion string `json:"gaSchemaVersion"`
	// Kinds is the per-kind depth accounting, sorted by kind.
	Kinds []KindAccounting `json:"kinds"`
	// Dispositions covers every GA resource, sorted by name. Resources
	// with no recorded disposition have an empty Disposition and a
	// Finding.
	Dispositions []ResourceDisposition `json:"dispositions"`
	// DispositionTotals counts resources per disposition ("" counts the
	// undispositioned).
	DispositionTotals map[string]int `json:"dispositionTotals"`
	// Findings is every gap, sorted by baseline key then detail.
	Findings []Finding `json:"findings,omitempty"`
}

// BuildAccounting runs the total-accounting check for one cloud provider's
// catalog: censuses and manifests from the tree, schemas from the committed
// artifacts, the dispositions ledger from dispositionsPath (empty string for
// the default). gaSchema names the parity-baseline schema (e.g. "google").
func BuildAccounting(repoRoot string, provider cloudresourcekind.CloudResourceProvider, schemas map[string]*Schema, gaSchema, dispositionsPath string) (Accounting, error) {
	if _, ok := schemas[gaSchema]; !ok {
		return Accounting{}, errors.Errorf("GA schema %q is not among the loaded schemas", gaSchema)
	}
	spec := SpecCensus(provider)
	modules, err := ModuleCensusForProvider(repoRoot, provider)
	if err != nil {
		return Accounting{}, err
	}
	manifests := map[string]*Manifest{}
	for _, m := range modules {
		kind := crkreflect.KindFromString(m.Kind)
		manifest, err := LoadKindManifest(repoRoot, provider, kind)
		if err != nil {
			return Accounting{}, err
		}
		if manifest != nil {
			manifests[m.Kind] = manifest
		}
	}
	if dispositionsPath == "" {
		dispositionsPath = filepath.Join(repoRoot, DefaultDispositionsDir, gaSchema+".yaml")
	}
	ledger, err := LoadLedger(dispositionsPath, gaSchema)
	if err != nil {
		return Accounting{}, err
	}
	return buildAccounting(provider.String(), spec, modules, schemas, gaSchema, manifests, ledger), nil
}

// buildAccounting is the pure join -- everything I/O-free so the hermetic
// tests can drive every accounting shape without a repo tree.
func buildAccounting(cloudProvider string, spec []KindCensus, modules []ModuleCensus, schemas map[string]*Schema, gaSchema string, manifests map[string]*Manifest, ledger []LedgerEntry) Accounting {
	schemaNames := make([]string, 0, len(schemas))
	for name := range schemas {
		schemaNames = append(schemaNames, name)
	}
	sort.Strings(schemaNames)

	specPaths := map[string][]string{}
	for _, k := range spec {
		specPaths[k.Kind] = k.SpecFieldPaths
	}

	acc := Accounting{
		CloudProvider:     cloudProvider,
		GASchema:          gaSchema,
		GASchemaVersion:   schemas[gaSchema].Version,
		DispositionTotals: map[string]int{},
	}

	consumedBy := map[string][]string{}
	for _, m := range modules {
		for _, res := range m.Resources {
			consumedBy[res] = append(consumedBy[res], m.Kind)
		}
	}

	for _, m := range modules {
		if m.MissingModule {
			acc.Kinds = append(acc.Kinds, KindAccounting{Kind: m.Kind, MissingModule: true})
			acc.Findings = append(acc.Findings, Finding{"kind:" + m.Kind,
				fmt.Sprintf("no Terraform module directory at %s -- forge the module, or record the debt in the anatomy baseline and here", m.ModuleDir)})
			continue
		}
		manifest := manifests[m.Kind]
		if manifest == nil {
			manifest = &Manifest{}
		}
		ka := accountKind(m, specPaths[m.Kind], schemas, schemaNames, manifests[m.Kind] != nil, manifest)
		acc.Kinds = append(acc.Kinds, ka)
		key := "kind:" + m.Kind
		for _, arg := range ka.UnaccountedArgs {
			acc.Findings = append(acc.Findings, Finding{key,
				fmt.Sprintf("unaccounted provider argument %s -- match it, map it, or exclude it with a reason in the kind's %s", arg, ManifestFileName)})
		}
		for _, field := range ka.UncoveredSpecFields {
			acc.Findings = append(acc.Findings, Finding{key,
				fmt.Sprintf("spec field %s reaches no provider argument -- reverse drift, a missing mapping, or a platform field needing a specExclusions entry", field)})
		}
		for _, stale := range ka.ManifestStale {
			acc.Findings = append(acc.Findings, Finding{key, stale})
		}
	}

	// Breadth: every GA resource carries exactly one disposition.
	ledgerByResource := map[string]LedgerEntry{}
	for _, e := range ledger {
		ledgerByResource[e.Resource] = e
	}
	ga := schemas[gaSchema]
	gaNames := make([]string, 0, len(ga.Resources))
	for name := range ga.Resources {
		gaNames = append(gaNames, name)
	}
	sort.Strings(gaNames)
	// Computed classes always win over the ledger, and a ledger entry
	// shadowed by a computed class is a staleness finding: recorded judgment
	// that duplicates what the instrument derives on its own is judgment
	// nobody re-evaluates -- exactly the rot class this package exists to
	// eliminate.
	for _, name := range gaNames {
		d := ResourceDisposition{Resource: name}
		entry, inLedger := ledgerByResource[name]
		switch {
		case len(consumedBy[name]) > 0:
			d.Disposition = DispositionModeled
			kinds := append([]string(nil), consumedBy[name]...)
			sort.Strings(kinds)
			d.Detail = "consumed by " + strings.Join(kinds, ", ")
			if inLedger {
				acc.Findings = append(acc.Findings, Finding{"resource:" + name,
					"stale ledger entry: the resource is consumed by a module (modeled) -- remove it from the ledger"})
			}
		case iamResourceRe.MatchString(name):
			d.Disposition = DispositionIamCovered
			d.Detail = "per-resource IAM triplet, covered by the owning kind's additive iam_members field"
			if inLedger {
				acc.Findings = append(acc.Findings, Finding{"resource:" + name,
					"stale ledger entry: the IAM triplet is covered by the computed iam-covered class -- remove it from the ledger"})
			}
		case ga.Resources[name].Deprecated:
			d.Disposition = DispositionExcludedDeprecated
			d.Detail = "deprecated in the provider schema"
			if inLedger {
				acc.Findings = append(acc.Findings, Finding{"resource:" + name,
					"stale ledger entry: the schema already flags the resource deprecated (computed) -- the ledger's excluded-deprecated is only for doc-level deprecations; remove it"})
			}
		case inLedger:
			d.Disposition = entry.Disposition
			d.Detail = entry.Reason
		default:
			acc.Findings = append(acc.Findings, Finding{"resource:" + name,
				"GA resource has no recorded disposition -- consume it, or record composed/model-planned/deferred/excluded-deprecated in the dispositions ledger"})
		}
		acc.DispositionTotals[d.Disposition]++
		acc.Dispositions = append(acc.Dispositions, d)
	}
	for _, e := range ledger {
		if _, ok := ga.Resources[e.Resource]; !ok {
			acc.Findings = append(acc.Findings, Finding{"resource:" + e.Resource,
				fmt.Sprintf("stale ledger entry: no resource %s in %s@%s -- it was removed or renamed at the pin", e.Resource, gaSchema, ga.Version)})
		}
	}

	sort.Slice(acc.Findings, func(i, j int) bool {
		if acc.Findings[i].BaselineKey != acc.Findings[j].BaselineKey {
			return acc.Findings[i].BaselineKey < acc.Findings[j].BaselineKey
		}
		return acc.Findings[i].Detail < acc.Findings[j].Detail
	})
	return acc
}

// argMatcher applies one resource's recorded judgment: exclusions by exact
// path, then the longest recorded mapping prefix, then the default spec
// root. No name heuristics, by design.
type argMatcher struct {
	specRoot   string
	mappings   []Mapping // sorted by arg length, longest first
	exclusions map[string]string
}

func newArgMatcher(rm *ResourceManifest) argMatcher {
	m := argMatcher{specRoot: specPathRoot, exclusions: map[string]string{}}
	if rm == nil {
		return m
	}
	if rm.SpecRoot != "" {
		m.specRoot = rm.SpecRoot
	}
	m.mappings = append(m.mappings, rm.Mappings...)
	sort.Slice(m.mappings, func(i, j int) bool { return len(m.mappings[i].Arg) > len(m.mappings[j].Arg) })
	for _, ex := range rm.Exclusions {
		m.exclusions[ex.Arg] = ex.Reason
	}
	return m
}

// derive returns the spec path an argument must match, and whether a
// recorded mapping produced it.
func (m argMatcher) derive(argPath string) (string, bool) {
	for _, mp := range m.mappings {
		if argPath == mp.Arg {
			return mp.Spec, true
		}
		if strings.HasPrefix(argPath, mp.Arg+".") {
			return mp.Spec + argPath[len(mp.Arg):], true
		}
	}
	return m.specRoot + "." + argPath, false
}

// accountKind runs both accounting directions for one kind.
func accountKind(m ModuleCensus, kindSpecPaths []string, schemas map[string]*Schema, schemaNames []string, hasManifest bool, manifest *Manifest) KindAccounting {
	ka := KindAccounting{Kind: m.Kind, HasManifest: hasManifest}
	specSet := map[string]bool{}
	for _, p := range kindSpecPaths {
		specSet[p] = true
	}
	underPath := func(paths []string, prefix string) bool {
		for _, p := range paths {
			if p == prefix || strings.HasPrefix(p, prefix+".") {
				return true
			}
		}
		return false
	}

	coveredSpec := map[string]bool{}
	consumed := map[string]bool{}
	for _, res := range m.Resources {
		consumed[res] = true
		var block *Block
		for _, name := range schemaNames {
			if b, ok := schemas[name].Resources[res]; ok {
				block = b
				break
			}
		}
		if block == nil {
			ka.ManifestStale = append(ka.ManifestStale,
				fmt.Sprintf("consumed resource %s is unknown to every loaded schema -- the module outruns its pin or the artifacts are stale", res))
			continue
		}
		rm := manifest.Resources[res]
		if rm != nil && rm.Internal != "" {
			ka.InternalResources = append(ka.InternalResources, res)
			continue
		}
		matcher := newArgMatcher(rm)

		argPaths := map[string]bool{}
		for _, arg := range block.ConfigurableArgs("") {
			if arg.Deprecated || isMachineryArg(arg.Path) {
				continue
			}
			argPaths[arg.Path] = true
			ka.TotalArgs++
			if _, excluded := matcher.exclusions[arg.Path]; excluded {
				ka.ExcludedArgs++
				continue
			}
			derived, mapped := matcher.derive(arg.Path)
			if specSet[derived] {
				coveredSpec[derived] = true
				if mapped {
					ka.MappedArgs++
				} else {
					ka.MatchedArgs++
				}
				continue
			}
			ka.UnaccountedArgs = append(ka.UnaccountedArgs, res+": "+arg.Path)
		}

		if rm == nil {
			continue
		}
		// Manifest hygiene: judgment referencing surface that no longer
		// exists is a finding, not a warning -- after a pin bump these ARE
		// the migration work list.
		sortedArgs := make([]string, 0, len(argPaths))
		for p := range argPaths {
			sortedArgs = append(sortedArgs, p)
		}
		sort.Strings(sortedArgs)
		if rm.SpecRoot != "" && !underPath(kindSpecPaths, rm.SpecRoot) {
			ka.ManifestStale = append(ka.ManifestStale,
				fmt.Sprintf("%s: specRoot %s matches no spec field", res, rm.SpecRoot))
		}
		for _, mp := range rm.Mappings {
			if !underPath(sortedArgs, mp.Arg) {
				ka.ManifestStale = append(ka.ManifestStale,
					fmt.Sprintf("%s: mapping arg %s matches no configurable argument at the pin", res, mp.Arg))
			}
			if !underPath(kindSpecPaths, mp.Spec) {
				ka.ManifestStale = append(ka.ManifestStale,
					fmt.Sprintf("%s: mapping spec %s matches no spec field", res, mp.Spec))
			}
		}
		for _, ex := range rm.Exclusions {
			if !argPaths[ex.Arg] {
				ka.ManifestStale = append(ka.ManifestStale,
					fmt.Sprintf("%s: exclusion %s matches no configurable argument at the pin -- remove it", res, ex.Arg))
			}
		}
	}

	for res := range manifest.Resources {
		if !consumed[res] {
			ka.ManifestStale = append(ka.ManifestStale,
				fmt.Sprintf("manifest judges resource %s, which the module no longer consumes", res))
		}
	}

	// Reverse direction: every spec leaf reaches provider surface or
	// carries a recorded exclusion.
	specExcluded := func(path string) bool {
		for _, ex := range manifest.SpecExclusions {
			if path == ex.Field || strings.HasPrefix(path, ex.Field+".") {
				return true
			}
		}
		return false
	}
	for _, p := range kindSpecPaths {
		if !coveredSpec[p] && !specExcluded(p) {
			ka.UncoveredSpecFields = append(ka.UncoveredSpecFields, p)
		}
	}
	for _, ex := range manifest.SpecExclusions {
		if !underPath(kindSpecPaths, ex.Field) {
			ka.ManifestStale = append(ka.ManifestStale,
				fmt.Sprintf("specExclusions: %s matches no spec field", ex.Field))
		}
	}

	sort.Strings(ka.InternalResources)
	sort.Strings(ka.UnaccountedArgs)
	sort.Strings(ka.UncoveredSpecFields)
	sort.Strings(ka.ManifestStale)
	return ka
}
