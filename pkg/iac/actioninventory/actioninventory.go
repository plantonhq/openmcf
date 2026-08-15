// Package actioninventory holds the committed snapshots of the cloud
// providers' own permission-action inventories and the matching rules the
// permissions conformance gate uses to prove that every action a runner
// permissions manifest (catalog/<provider>/<kind>/iac/permissions.yaml)
// names actually exists. A least-privilege policy naming an action the
// provider never defined is a fabrication that no syntax check can catch --
// the inventory makes that class of wrong structurally impossible.
//
// The snapshot is generated data used only by this gate: like
// pkg/anatomy's baseline.yaml, it lives beside the code that reads it and
// is refreshed by a make target (`make generate-action-inventory`) that
// talks to the provider's published inventory -- network in make, never in
// CI. CI validates manifests against the committed snapshot. The snapshot
// deliberately covers ONLY the services the committed manifests reference:
// a manifest naming a service the snapshot lacks fails the gate with the
// refresh command, and a refresh that cannot find the service in the
// provider's inventory is a hard error -- a genuinely wrong prefix.
//
// AWS is the first inventory arm, read from AWS's machine-readable service
// reference (https://servicereference.us-east-1.amazonaws.com). Unlike the
// price books' immutable versioned offer documents, the reference serves
// latest-only URLs, so each service records the retrieval date and the
// index's own modification stamp instead of claiming an immutability the
// source does not offer. Providers without a machine-readable inventory
// arm are exempt from existence checking (their structural validation
// lives in pkg/iac/permissions) -- exemption is stated here, never silent.
package actioninventory

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// AwsFileName is the AWS inventory snapshot's name inside this package's
// directory.
const AwsFileName = "aws.yaml"

// Inventory is one provider's committed action-inventory snapshot.
type Inventory struct {
	// Provider is the catalog provider directory name (e.g. "aws").
	Provider string `yaml:"provider"`
	// Services are the per-service action lists, sorted by prefix.
	Services []Service `yaml:"services"`
}

// Service is one cloud service's action list as the provider publishes it.
type Service struct {
	// Prefix is the service's action-name prefix in the provider's own
	// vocabulary -- for AWS, the IAM prefix before the colon in
	// "lambda:CreateFunction".
	Prefix string `yaml:"prefix"`
	// SourceURL is the provider inventory document the actions were read
	// from. The URL serves latest-only content; provenance is the URL plus
	// the dates below, never a claim of immutability.
	SourceURL string `yaml:"source_url"`
	// SourceModified is the provider's own last-modified date for this
	// service's inventory document (UTC, date-only), read from the
	// inventory index at fetch time.
	SourceModified string `yaml:"source_modified"`
	// RetrievedOn is the date this service's actions were fetched (UTC).
	RetrievedOn string `yaml:"retrieved_on"`
	// Actions are the action names the provider defines for this service,
	// sorted, without the service prefix.
	Actions []string `yaml:"actions"`
}

// LoadAws reads and strictly parses the committed AWS inventory snapshot
// and enforces its structural invariants (sorted unique services, sorted
// unique non-empty action lists) so gate results can never depend on
// snapshot ordering accidents.
func LoadAws(dir string) (*Inventory, error) {
	raw, err := os.ReadFile(path.Join(dir, AwsFileName))
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(strings.NewReader(string(raw)))
	decoder.KnownFields(true)
	inv := &Inventory{}
	if err := decoder.Decode(inv); err != nil {
		return nil, fmt.Errorf("parse %s: %w", AwsFileName, err)
	}
	if inv.Provider != "aws" {
		return nil, fmt.Errorf("%s: provider is %q, want \"aws\"", AwsFileName, inv.Provider)
	}
	if len(inv.Services) == 0 {
		return nil, fmt.Errorf("%s: no services -- run `make generate-action-inventory`", AwsFileName)
	}
	seen := map[string]bool{}
	for i, svc := range inv.Services {
		if svc.Prefix == "" || svc.SourceURL == "" || svc.SourceModified == "" || svc.RetrievedOn == "" {
			return nil, fmt.Errorf("%s: service %q: prefix, source_url, source_modified, and retrieved_on are all required", AwsFileName, svc.Prefix)
		}
		if seen[svc.Prefix] {
			return nil, fmt.Errorf("%s: duplicate service prefix %q", AwsFileName, svc.Prefix)
		}
		seen[svc.Prefix] = true
		if i > 0 && inv.Services[i-1].Prefix > svc.Prefix {
			return nil, fmt.Errorf("%s: services not sorted at %q", AwsFileName, svc.Prefix)
		}
		if len(svc.Actions) == 0 {
			return nil, fmt.Errorf("%s: service %q has no actions", AwsFileName, svc.Prefix)
		}
		for j, action := range svc.Actions {
			if action == "" {
				return nil, fmt.Errorf("%s: service %q has an empty action", AwsFileName, svc.Prefix)
			}
			if j > 0 {
				if prev := svc.Actions[j-1]; prev == action {
					return nil, fmt.Errorf("%s: service %q duplicates action %q", AwsFileName, svc.Prefix, action)
				} else if prev > action {
					return nil, fmt.Errorf("%s: service %q actions not sorted at %q", AwsFileName, svc.Prefix, action)
				}
			}
		}
	}
	return inv, nil
}

// ServiceActions returns the snapshot's action list for a service prefix,
// or nil when the snapshot does not cover the service.
func (inv *Inventory) ServiceActions(prefix string) []string {
	for i := range inv.Services {
		if inv.Services[i].Prefix == prefix {
			return inv.Services[i].Actions
		}
	}
	return nil
}

// MatchAction reports how many of a service's published actions an IAM
// action name matches. IAM evaluates action names case-insensitively and
// supports "*" (any run) and "?" (any one character) wildcards, so the
// gate does too: an exact name matches its published spelling in any case,
// and a wildcard pattern must match at least one published action -- a
// pattern matching NOTHING is the fabrication class the inventory exists
// to catch, and wildcards are where it hides best.
func MatchAction(published []string, name string) int {
	pattern := strings.ToLower(name)
	matches := 0
	for _, action := range published {
		if globMatch(pattern, strings.ToLower(action)) {
			matches++
		}
	}
	return matches
}

// globMatch implements IAM's action wildcard semantics ("*" any run, "?"
// any single character) over lowercased inputs. Iterative with
// backtracking on the last "*", the classic linear-scan shape.
func globMatch(pattern, s string) bool {
	pIdx, sIdx := 0, 0
	starIdx, starMatch := -1, 0
	for sIdx < len(s) {
		switch {
		case pIdx < len(pattern) && (pattern[pIdx] == '?' || pattern[pIdx] == s[sIdx]):
			pIdx++
			sIdx++
		case pIdx < len(pattern) && pattern[pIdx] == '*':
			starIdx = pIdx
			starMatch = sIdx
			pIdx++
		case starIdx != -1:
			pIdx = starIdx + 1
			starMatch++
			sIdx = starMatch
		default:
			return false
		}
	}
	for pIdx < len(pattern) && pattern[pIdx] == '*' {
		pIdx++
	}
	return pIdx == len(pattern)
}

// header is the snapshot's leading comment, written by Render so the
// fetcher and any future regeneration path can never disagree on it.
const header = `# GENERATED -- DO NOT EDIT. Refresh with ` + "`make generate-action-inventory`" + `.
# Committed snapshot of the provider's own permission-action inventory,
# scoped to exactly the services the committed runner permissions manifests
# (catalog/*/*/iac/permissions.yaml) reference. The conformance gate in
# this package proves every manifest action exists here -- an action name
# the provider never defined cannot ship. Source URLs serve latest-only
# content; provenance is the URL plus the retrieval date and the source's
# own modification stamp.
`

// Render writes an inventory in its canonical byte form: header, sorted
// services, sorted actions, dates quoted. One renderer home keeps refresh
// diffs meaningful.
func Render(inv *Inventory) string {
	var b strings.Builder
	b.WriteString(header)
	b.WriteString("provider: " + inv.Provider + "\n")
	b.WriteString("services:\n")
	services := append([]Service(nil), inv.Services...)
	sort.Slice(services, func(i, j int) bool { return services[i].Prefix < services[j].Prefix })
	for _, svc := range services {
		b.WriteString("  - prefix: " + svc.Prefix + "\n")
		b.WriteString("    source_url: " + svc.SourceURL + "\n")
		b.WriteString("    source_modified: \"" + svc.SourceModified + "\"\n")
		b.WriteString("    retrieved_on: \"" + svc.RetrievedOn + "\"\n")
		b.WriteString("    actions:\n")
		actions := append([]string(nil), svc.Actions...)
		sort.Strings(actions)
		for _, action := range actions {
			b.WriteString("      - " + action + "\n")
		}
	}
	return b.String()
}
