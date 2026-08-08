//go:build !codegen
// +build !codegen

// Per-kind mapping manifests: the recorded-judgment side of provider parity.
// Detection is mechanical (the censuses); judgment -- "this provider argument
// is this spec field under another name", "this one is deliberately not
// modeled, and here is why" -- is recorded ONCE, in the kind's own manifest,
// where the tooling reads it and the next author finds it.
//
// The manifest lives at catalog/<provider>/<kind>/iac/provider-parity.yaml,
// beside iac/import-map.yaml, and file presence IS enrollment (the import-map
// convention): a kind with a manifest is held to total accounting; a kind
// without one is an accepted gap in the burn-down baseline. There is
// deliberately no separate allowlist to keep in sync.
//
// The matcher the manifest feeds carries ZERO name heuristics: an argument
// either matches its derived spec path EXACTLY or its divergence is recorded
// here. Cleverness (singularization, suffix stripping) would eventually
// false-match and silently overclaim coverage -- the one failure mode a
// parity instrument must make impossible.

package providerparity

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"gopkg.in/yaml.v3"
)

// ManifestFileName is the per-kind mapping manifest, at the component's
// iac/ root beside import-map.yaml.
const ManifestFileName = "provider-parity.yaml"

// Manifest is one kind's recorded provider-parity judgment.
type Manifest struct {
	// Resources carries one entry per consumed Terraform resource that needs
	// recorded judgment. A consumed resource absent from the manifest is
	// held to pure exact matching under the default spec root.
	Resources map[string]*ResourceManifest `yaml:"resources"`
	// SpecExclusions are spec leaf fields with no provider counterpart --
	// platform-level concepts the reverse accounting must not flag. Each
	// carries its reason.
	SpecExclusions []SpecExclusion `yaml:"specExclusions"`
}

// ResourceManifest is the judgment for one consumed resource.
type ResourceManifest struct {
	// Internal marks the whole resource as module plumbing (e.g. API
	// enablement via google_project_service): its arguments are module
	// decisions with fixed or derived values, not spec surface. The value
	// is the mandatory reason. Mutually exclusive with every other field.
	Internal string `yaml:"internal,omitempty"`
	// SpecRoot is the spec subtree this resource's arguments match under.
	// Defaults to "spec" (the kind's primary resource); a secondary
	// resource names the subtree that instantiates it, e.g.
	// "spec.iam_members" for the kind's *_iam_member resource.
	SpecRoot string `yaml:"specRoot,omitempty"`
	// Mappings record name divergences: a provider argument (or argument
	// subtree) that is this spec field (or subtree) under another name.
	// Exact matching resumes below a mapped subtree, so mappings stay
	// O(divergences), never O(fields).
	Mappings []Mapping `yaml:"mappings,omitempty"`
	// Exclusions are configurable arguments deliberately not modeled, each
	// with its mandatory reason -- the recorded decision that makes the
	// omission loud instead of silent. An exclusion may name an argument
	// subtree (a block): one judgment then covers everything under it,
	// mirroring mappings and specExclusions -- the form for a resource
	// that embeds another kind's whole surface inline (the exclusion stays
	// one recorded decision, not one line per leaf).
	Exclusions []ArgExclusion `yaml:"exclusions,omitempty"`
}

// Mapping is one recorded rename: arg (relative to the resource's root
// block, e.g. "lifecycle_rule.condition.age") is spec (an absolute
// spec-rooted path, e.g. "spec.lifecycle_rules.condition.age_days").
// Both sides may name a subtree; matching resumes exactly below it.
type Mapping struct {
	Spec string `yaml:"spec"`
	Arg  string `yaml:"arg"`
}

// ArgExclusion is one deliberately unmodeled argument (or argument
// subtree) and its reason.
type ArgExclusion struct {
	Arg    string `yaml:"arg"`
	Reason string `yaml:"reason"`
}

// SpecExclusion is one spec leaf field with no provider counterpart.
type SpecExclusion struct {
	Field  string `yaml:"field"`
	Reason string `yaml:"reason"`
}

// ManifestPath composes the manifest location for one kind, mirroring the
// module census's path grammar (iac/ sits at the component root; the
// version segment never appears).
func ManifestPath(repoRoot string, provider cloudresourcekind.CloudResourceProvider, kind cloudresourcekind.CloudResourceKind) string {
	return filepath.Join(repoRoot, catalogRoot, provider.String(),
		strings.ToLower(kind.String()), "iac", ManifestFileName)
}

// LoadManifest reads and structurally validates one manifest. Parsing is
// STRICT (unknown fields rejected) and validation failures are hard errors,
// never findings: a manifest is an explicit authoring act, so a malformed
// one means the author's judgment was not recorded the way they intended --
// stop loudly rather than account against a half-read record.
func LoadManifest(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "opening manifest %s", path)
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	var m Manifest
	if err := dec.Decode(&m); err != nil {
		return nil, errors.Wrapf(err, "manifest %s is not valid (strict) YAML", path)
	}
	if err := m.validate(); err != nil {
		return nil, errors.Wrapf(err, "manifest %s", path)
	}
	return &m, nil
}

// LoadKindManifest loads one kind's manifest, returning (nil, nil) when the
// kind ships none -- absence is the enrollment signal, not an error.
func LoadKindManifest(repoRoot string, provider cloudresourcekind.CloudResourceProvider, kind cloudresourcekind.CloudResourceKind) (*Manifest, error) {
	path := ManifestPath(repoRoot, provider, kind)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "checking manifest %s", path)
	}
	return LoadManifest(path)
}

// specPathRoot is the root of every census spec path (CollectSpecPaths
// prefixes with "spec"), and therefore of every manifest spec reference.
const specPathRoot = "spec"

func isSpecPath(p string) bool {
	return p == specPathRoot || strings.HasPrefix(p, specPathRoot+".")
}

func (m *Manifest) validate() error {
	for resource, rm := range m.Resources {
		if resource == "" {
			return errors.New("resources must be keyed by Terraform resource type")
		}
		if rm == nil {
			return errors.Errorf("resource %s carries no judgment -- an empty entry records nothing; remove it", resource)
		}
		if rm.Internal != "" {
			if rm.SpecRoot != "" || len(rm.Mappings) > 0 || len(rm.Exclusions) > 0 {
				return errors.Errorf(
					"resource %s: internal is the whole judgment -- it cannot carry specRoot/mappings/exclusions", resource)
			}
			continue
		}
		if rm.SpecRoot == "" && len(rm.Mappings) == 0 && len(rm.Exclusions) == 0 {
			return errors.Errorf("resource %s carries no judgment -- an empty entry records nothing; remove it", resource)
		}
		if rm.SpecRoot != "" && !isSpecPath(rm.SpecRoot) {
			return errors.Errorf("resource %s: specRoot %q must be spec-rooted (e.g. spec.iam_members)", resource, rm.SpecRoot)
		}
		seenArgs := map[string]bool{}
		for _, mp := range rm.Mappings {
			if mp.Arg == "" {
				return errors.Errorf("resource %s: mapping for spec %q names no arg", resource, mp.Spec)
			}
			if !isSpecPath(mp.Spec) {
				return errors.Errorf("resource %s: mapping for arg %q must be spec-rooted, got %q", resource, mp.Arg, mp.Spec)
			}
			if seenArgs[mp.Arg] {
				return errors.Errorf("resource %s: arg %q is judged twice", resource, mp.Arg)
			}
			seenArgs[mp.Arg] = true
		}
		for _, ex := range rm.Exclusions {
			if ex.Arg == "" {
				return errors.Errorf("resource %s: exclusion names no arg", resource)
			}
			if ex.Reason == "" {
				return errors.Errorf("resource %s: exclusion of %q carries no reason -- an unexplained omission is the failure class this file exists to prevent", resource, ex.Arg)
			}
			if seenArgs[ex.Arg] {
				return errors.Errorf("resource %s: arg %q is judged twice", resource, ex.Arg)
			}
			seenArgs[ex.Arg] = true
		}
	}
	seenFields := map[string]bool{}
	for _, ex := range m.SpecExclusions {
		if !isSpecPath(ex.Field) {
			return errors.Errorf("specExclusions: field %q must be spec-rooted", ex.Field)
		}
		if ex.Reason == "" {
			return errors.Errorf("specExclusions: %q carries no reason", ex.Field)
		}
		if seenFields[ex.Field] {
			return errors.Errorf("specExclusions: %q is listed twice", ex.Field)
		}
		seenFields[ex.Field] = true
	}
	return nil
}
