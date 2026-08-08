//go:build !codegen
// +build !codegen

// Module census: the catalog's consumption side of the parity measurement.
// For every registered kind of one cloud provider, which Terraform resources
// its module actually declares and which provider releases it pins.
//
// The scan reads EVERY `*.tf` file in the module directory, never main.tf
// alone: modules may split resources across sibling files, and a main.tf-only
// scan undercounts silently -- the worst failure mode a parity instrument can
// have. Resource declarations are matched by the same top-level pattern the
// import-map conformance guard scans.

package providerparity

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// catalogRoot is where components live, relative to the repo root. A
// package-local constant by convention (pkg/anatomy, pkg/e2e/profile,
// pkg/iac/importmap do the same); the version segment never appears in
// module paths -- iac/ sits at the component root.
const catalogRoot = "catalog"

// tfResourceRe matches top-level terraform resource declarations.
var tfResourceRe = regexp.MustCompile(`(?m)^resource\s+"([a-z0-9_]+)"\s+"`)

// tfProviderPinRe pulls each required_providers entry's local name, source,
// and version constraint. The catalog's modules declare pins in one uniform
// block shape, so a structural regex is sufficient here, as it is for the
// tf-provider-pins CI guard.
var tfProviderPinRe = regexp.MustCompile(`(?s)([A-Za-z0-9_-]+)\s*=\s*\{[^{}]*?source\s*=\s*"([^"]+)"[^{}]*?version\s*=\s*"([^"]+)"[^{}]*?\}`)

// ModuleCensus is one kind's Terraform-module consumption.
type ModuleCensus struct {
	// Kind is the registry name, e.g. "GcpGcsBucket".
	Kind string
	// ModuleDir is the repo-root-relative Terraform module directory.
	ModuleDir string
	// MissingModule marks a registered, implemented kind with no Terraform
	// module directory at all -- an anatomy-baselined gap in some catalogs.
	// Recorded here (and surfaced as an accounting Finding) rather than
	// aborting the census, so one kind's debt never hides every other
	// kind's numbers.
	MissingModule bool
	// Resources is the sorted distinct resource types the module declares.
	Resources []string
	// Pins maps each declared provider's local name to its version
	// constraint, e.g. {"google": "~> 6.0"}.
	Pins map[string]string
}

// ModuleCensusForProvider scans the Terraform module of every registered kind
// of one cloud provider, sorted by kind. A registered kind whose module
// directory is missing is a per-kind FINDING, never a silent skip and never
// a run-abort: the finding gates through the burn-down baseline, so an
// anatomy-clean catalog (zero baseline entries) still fails CI on the first
// missing module, while a catalog carrying recorded anatomy debt can be
// measured without that debt hiding every other kind's numbers.
func ModuleCensusForProvider(repoRoot string, provider cloudresourcekind.CloudResourceProvider) ([]ModuleCensus, error) {
	var out []ModuleCensus
	for _, kind := range crkreflect.KindsList() {
		if crkreflect.GetProvider(kind) != provider {
			continue
		}
		if _, err := crkreflect.NewInstance(kind); err != nil {
			continue // enum value exists but the kind is not implemented yet
		}
		moduleRel := filepath.Join(catalogRoot, provider.String(), strings.ToLower(kind.String()), "iac", "tf")
		moduleAbs := filepath.Join(repoRoot, moduleRel)
		if _, err := os.Stat(moduleAbs); os.IsNotExist(err) {
			out = append(out, ModuleCensus{Kind: kind.String(), ModuleDir: moduleRel, MissingModule: true})
			continue
		}
		census, err := ScanModule(moduleAbs)
		if err != nil {
			return nil, errors.Wrapf(err, "kind %s", kind.String())
		}
		census.Kind = kind.String()
		census.ModuleDir = moduleRel
		out = append(out, census)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Kind < out[j].Kind })
	return out, nil
}

// ScanModule reads every top-level `*.tf` file in one module directory and
// returns its declared resources and provider pins. Exposed so tests can
// drive it against fixture modules in isolation.
func ScanModule(moduleDir string) (ModuleCensus, error) {
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		return ModuleCensus{}, errors.Wrapf(err, "reading module dir %s", moduleDir)
	}
	seen := map[string]bool{}
	census := ModuleCensus{Pins: map[string]string{}}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tf") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(moduleDir, entry.Name()))
		if err != nil {
			return ModuleCensus{}, errors.Wrapf(err, "reading %s", entry.Name())
		}
		for _, m := range tfResourceRe.FindAllStringSubmatch(string(content), -1) {
			if !seen[m[1]] {
				seen[m[1]] = true
				census.Resources = append(census.Resources, m[1])
			}
		}
		for _, m := range tfProviderPinRe.FindAllStringSubmatch(string(content), -1) {
			census.Pins[m[1]] = m[3]
		}
	}
	sort.Strings(census.Resources)
	return census, nil
}
