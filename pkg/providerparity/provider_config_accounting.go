//go:build !codegen
// +build !codegen

// Provider-config accounting: total accounting for the PROVIDER BLOCK, the
// depth accounting's twin one level up. Every configurable, non-deprecated
// argument of the provider's own configuration block is exact-matched to a
// provider-config proto field, mapped by recorded judgment, module-owned by
// recorded judgment, or excluded with a reason -- and, in reverse, every
// config proto leaf lands on a provider-block argument or carries an
// exclusion. Provider-block arguments set inside catalog modules' own
// provider blocks (found mechanically by the HCL module census) must carry
// recorded judgment too, so a behavior flag hand-placed in a module can
// never be an invisible leak again.
//
// Enrollment is manifest presence (provider_config_manifest.go); findings
// gate through the same burn-down baseline as depth and breadth, under the
// third baseline-key class "provider:<cloud>" -- one work-list line per
// provider, the "kind:"/"resource:" grain applied to the provider block.

package providerparity

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ProviderConfigAccounting is one provider's provider-block accounting.
type ProviderConfigAccounting struct {
	// TotalArgs is the provider block's non-deprecated configurable
	// argument surface at the pin.
	TotalArgs int `json:"totalArgs"`
	// MatchedArgs matched their derived config path exactly.
	MatchedArgs int `json:"matchedArgs"`
	// MappedArgs matched through a recorded mapping.
	MappedArgs int `json:"mappedArgs"`
	// ModuleOwnedArgs carry a recorded module-owned judgment.
	ModuleOwnedArgs int `json:"moduleOwnedArgs"`
	// ExcludedArgs carry a recorded exclusion.
	ExcludedArgs int `json:"excludedArgs"`
	// UnaccountedArgs have no match, mapping, module-owned judgment, or
	// exclusion -- each is a Finding.
	UnaccountedArgs []string `json:"unaccountedArgs,omitempty"`
	// UncoveredConfigFields are config proto leaves no argument matched and
	// no exclusion covers -- each is a Finding.
	UncoveredConfigFields []string `json:"uncoveredConfigFields,omitempty"`
	// UnjudgedModuleArgs are provider-block arguments set inside catalog
	// modules with no recorded judgment ("arg (set by Kind1, Kind2)") --
	// each is a Finding.
	UnjudgedModuleArgs []string `json:"unjudgedModuleArgs,omitempty"`
	// ManifestStale are manifest entries referencing surface that no longer
	// exists. Each is a Finding.
	ManifestStale []string `json:"manifestStale,omitempty"`
}

// Accounted reports whether the provider block is at total accounting.
func (p ProviderConfigAccounting) Accounted() bool {
	return len(p.UnaccountedArgs) == 0 && len(p.UncoveredConfigFields) == 0 &&
		len(p.UnjudgedModuleArgs) == 0 && len(p.ManifestStale) == 0
}

// buildProviderConfigAccounting is the pure join for one provider's
// provider-block accounting. configPaths is the config proto census; block
// is the distilled schema's provider block; moduleSetArgs maps each
// module-set provider-block argument to the kinds whose modules set it.
func buildProviderConfigAccounting(cloudProvider string, configPaths []string, block *Block,
	manifest *ProviderConfigManifest, moduleSetArgs map[string][]string) (ProviderConfigAccounting, []Finding) {

	pc := ProviderConfigAccounting{}
	key := "provider:" + cloudProvider

	// Reuse the per-kind argMatcher verbatim: same longest-mapping-wins,
	// same collapse semantics, same subtree exclusions -- one matching law
	// for both accounting levels.
	rm := &ResourceManifest{SpecRoot: configPathRoot}
	for _, mp := range manifest.Mappings {
		rm.Mappings = append(rm.Mappings, Mapping{Spec: mp.Config, Arg: mp.Arg, Collapse: mp.Collapse})
	}
	for _, ex := range manifest.Exclusions {
		rm.Exclusions = append(rm.Exclusions, ex)
	}
	matcher := newArgMatcher(rm)

	moduleOwned := func(argPath string) bool {
		for _, mo := range manifest.ModuleOwned {
			if argPath == mo.Arg || strings.HasPrefix(argPath, mo.Arg+".") {
				return true
			}
		}
		return false
	}

	// Patterns are pre-validated compilable by the manifest loader; compile
	// here for the walk and track per-pattern hits for staleness.
	patterns := make([]*regexp.Regexp, len(manifest.ExclusionPatterns))
	patternHits := make([]int, len(manifest.ExclusionPatterns))
	for i, pe := range manifest.ExclusionPatterns {
		patterns[i] = regexp.MustCompile(pe.Pattern)
	}
	patternExcluded := func(argPath string) bool {
		hit := false
		for i, re := range patterns {
			if re.MatchString(argPath) {
				patternHits[i]++
				hit = true
			}
		}
		return hit
	}

	configSet := map[string]bool{}
	for _, p := range configPaths {
		configSet[p] = true
	}
	underPath := func(paths []string, prefix string) bool {
		for _, p := range paths {
			if p == prefix || strings.HasPrefix(p, prefix+".") {
				return true
			}
		}
		return false
	}

	coveredConfig := map[string]bool{}
	argPaths := map[string]bool{}
	for _, arg := range block.ConfigurableArgs("") {
		if arg.Deprecated || isMachineryArg(arg.Path) {
			continue
		}
		argPaths[arg.Path] = true
		pc.TotalArgs++
		if matcher.excluded(arg.Path) || patternExcluded(arg.Path) {
			pc.ExcludedArgs++
			continue
		}
		if moduleOwned(arg.Path) {
			pc.ModuleOwnedArgs++
			continue
		}
		derived, mapped := matcher.derive(arg.Path)
		matchedAny := false
		for _, d := range derived {
			if configSet[d] {
				coveredConfig[d] = true
				matchedAny = true
			}
		}
		if matchedAny {
			if mapped {
				pc.MappedArgs++
			} else {
				pc.MatchedArgs++
			}
			continue
		}
		pc.UnaccountedArgs = append(pc.UnaccountedArgs, arg.Path)
	}

	// Module-set arguments must carry judgment: module-owned is the natural
	// class, but an argument that reaches a config field (the module wires
	// the config value through, e.g. a TF_VAR-fed token -- whether the
	// correspondence is exact-named or mapped) and an excluded one are
	// judgments too. Only a judgment-free module-set argument is a leak.
	moduleArgs := make([]string, 0, len(moduleSetArgs))
	for arg := range moduleSetArgs {
		moduleArgs = append(moduleArgs, arg)
	}
	sort.Strings(moduleArgs)
	for _, arg := range moduleArgs {
		if moduleOwned(arg) || matcher.excluded(arg) || patternExcluded(arg) {
			continue
		}
		derived, _ := matcher.derive(arg)
		reachesConfig := false
		for _, d := range derived {
			if configSet[d] {
				reachesConfig = true
			}
		}
		if reachesConfig {
			continue
		}
		kinds := append([]string(nil), moduleSetArgs[arg]...)
		sort.Strings(kinds)
		pc.UnjudgedModuleArgs = append(pc.UnjudgedModuleArgs,
			fmt.Sprintf("%s (set by %s)", arg, strings.Join(kinds, ", ")))
	}

	// Manifest hygiene: judgment referencing surface that no longer exists.
	sortedArgs := make([]string, 0, len(argPaths))
	for p := range argPaths {
		sortedArgs = append(sortedArgs, p)
	}
	sort.Strings(sortedArgs)
	for _, mp := range manifest.Mappings {
		if !underPath(sortedArgs, mp.Arg) {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("mapping arg %s matches no provider-block argument at the pin", mp.Arg))
		}
		if !underPath(configPaths, mp.Config) {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("mapping config %s matches no provider-config field", mp.Config))
		}
		if mp.Collapse && !configSet[mp.Config] {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("collapse mapping config %s must name a config census leaf, not a subtree", mp.Config))
		}
	}
	for _, ex := range manifest.Exclusions {
		if !underPath(sortedArgs, ex.Arg) {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("exclusion %s matches no provider-block argument at the pin -- remove it", ex.Arg))
		}
	}
	for _, mo := range manifest.ModuleOwned {
		if !underPath(sortedArgs, mo.Arg) {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("moduleOwned %s matches no provider-block argument at the pin -- remove it", mo.Arg))
			continue
		}
		if !underPath(moduleArgs, mo.Arg) {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("moduleOwned %s is set by no module's provider block -- remove it or record it as an exclusion", mo.Arg))
		}
	}
	for i, pe := range manifest.ExclusionPatterns {
		if patternHits[i] == 0 {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("exclusionPatterns %s matches no provider-block argument at the pin -- remove it", pe.Pattern))
		}
	}

	// Reverse direction: every config leaf reaches the provider block or
	// carries a recorded exclusion.
	configExcluded := func(path string) bool {
		for _, ex := range manifest.ConfigExclusions {
			if path == ex.Field || strings.HasPrefix(path, ex.Field+".") {
				return true
			}
		}
		return false
	}
	for _, p := range configPaths {
		if !coveredConfig[p] && !configExcluded(p) {
			pc.UncoveredConfigFields = append(pc.UncoveredConfigFields, p)
		}
	}
	for _, ex := range manifest.ConfigExclusions {
		if !underPath(configPaths, ex.Field) {
			pc.ManifestStale = append(pc.ManifestStale,
				fmt.Sprintf("configExclusions: %s matches no provider-config field", ex.Field))
		}
	}

	sort.Strings(pc.UnaccountedArgs)
	sort.Strings(pc.UncoveredConfigFields)
	sort.Strings(pc.UnjudgedModuleArgs)
	sort.Strings(pc.ManifestStale)

	var findings []Finding
	for _, arg := range pc.UnaccountedArgs {
		findings = append(findings, Finding{key,
			fmt.Sprintf("unaccounted provider-block argument %s -- match it, map it, or exclude it with a reason in %s", arg, ProviderConfigManifestFileName)})
	}
	for _, field := range pc.UncoveredConfigFields {
		findings = append(findings, Finding{key,
			fmt.Sprintf("provider-config field %s reaches no provider-block argument -- reverse drift, a missing mapping, or a platform field needing a configExclusions entry", field)})
	}
	for _, arg := range pc.UnjudgedModuleArgs {
		findings = append(findings, Finding{key,
			fmt.Sprintf("provider-block argument %s is set inside module provider blocks with no recorded judgment -- record it as moduleOwned (or map/exclude it) in %s", arg, ProviderConfigManifestFileName)})
	}
	for _, stale := range pc.ManifestStale {
		findings = append(findings, Finding{key, stale})
	}
	return pc, findings
}
