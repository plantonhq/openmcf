//go:build !codegen
// +build !codegen

// Enrollment: which catalog providers are under provider-parity accounting,
// and against which GA schema each one is measured. The pairing is explicit
// by design -- the same choice as the CLI's always-explicit --provider /
// --ga-schema flags -- because the mapping is not derivable: a catalog
// provider name ("gcp") and its parity-baseline schema name ("google") need
// not match, and guessing the pairing is exactly the kind of cleverness a
// parity instrument must refuse. Enrolling a provider is one reviewed line,
// added when its schema artifact is committed and its catalog is ready to
// be gated.
//
// The enrollment table is also what makes the shared baseline SAFE to
// write: baseline.yaml is one file for every provider, so any write must be
// computed from ALL enrolled providers' findings. EnrolledAccountings +
// MergeFindings are the single write-input path for both the CI gate's
// regeneration mode and the CLI's --write-baseline -- one source of truth,
// mirroring how Gate() is shared by the CI test and the CLI --check, so no
// caller can truncate the baseline to a single provider's findings.

package providerparity

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// Enrollment binds one catalog provider to the GA schema its parity is
// declared against.
type Enrollment struct {
	Provider cloudresourcekind.CloudResourceProvider
	// GASchema names the committed schema artifact (schemas/<name>-*.json.gz)
	// that is this provider's parity baseline.
	GASchema string
}

// Enrollments is every catalog provider under provider-parity accounting.
// Order is stable (report/log order); append new providers at the end.
var Enrollments = []Enrollment{
	{Provider: cloudresourcekind.CloudResourceProvider_gcp, GASchema: "google"},
	{Provider: cloudresourcekind.CloudResourceProvider_aws, GASchema: "aws"},
}

// EnrolledAccountings runs the full accounting for every enrolled provider.
// It is the input side of every baseline write and of the CI gate: gating or
// writing from any subset would report the other providers' baseline entries
// as stale (the gate) or silently drop them (the write).
func EnrolledAccountings(repoRoot string, schemas map[string]*Schema) ([]Accounting, error) {
	accountings := make([]Accounting, 0, len(Enrollments))
	for _, e := range Enrollments {
		acc, err := BuildAccounting(repoRoot, e.Provider, schemas, e.GASchema, "")
		if err != nil {
			return nil, errors.Wrapf(err, "accounting %s catalog against GA schema %q", e.Provider, e.GASchema)
		}
		accountings = append(accountings, acc)
	}
	return accountings, nil
}

// MergeFindings flattens the enrolled accountings' findings into the one
// slice WriteBaseline and Gate consume. Trivial on purpose: the value is
// that every baseline write and gate check goes through the same union.
func MergeFindings(accountings []Accounting) []Finding {
	var all []Finding
	for _, acc := range accountings {
		all = append(all, acc.Findings...)
	}
	return all
}
