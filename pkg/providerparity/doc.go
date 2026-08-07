// Package providerparity measures the catalog against the Terraform
// provider it deploys through: for each cloud provider, what the pinned
// Terraform provider can configure versus what the catalog's kinds expose.
//
// Naming note -- this is PROVIDER parity, a different axis from the
// cross-engine parity the audit machinery already enforces:
//
//   - Cross-engine parity (pkg/iac/MODULE_PARITY.md, PARITY-EXCEPTION
//     comments, the audit rule's --parity focus): the Terraform and Pulumi
//     modules of ONE kind implement the same contract identically.
//   - Provider parity (this package): the catalog as a whole, kind by kind,
//     covers the full configurable surface of the pinned Terraform provider,
//     with every omission recorded as a decision.
//
// The measurement has three independent censuses, each mechanical:
//
//   - The provider side is self-describing: committed, distilled
//     `providers schema -json` artifacts per pinned provider version
//     (schemas/, produced by the distiller sub-package). No source parsing.
//   - The catalog's contract side is a proto field census from descriptors
//     via pkg/crkreflect (spec_census.go) -- the same registry walk as
//     pkg/secretcoverage; a reader who knows one walk knows both.
//   - The catalog's module side is a consumed-resource and provider-pin
//     census over every `*.tf` file of every kind's Terraform module
//     (module_census.go). Every file, never main.tf alone: modules may
//     split resources across sibling files, and a main.tf-only scan
//     undercounts silently.
//
// Parity is always measured against a NAMED provider version -- the pin --
// never against "latest"; that is what makes a freshness promise
// well-defined.
//
// On top of the measurement sits the total-accounting layer (accounting.go):
// every configurable, non-deprecated argument of every consumed resource is
// exact-matched to a spec field, mapped by the kind's recorded judgment
// (manifest.go -- catalog/<provider>/<kind>/iac/provider-parity.yaml), or
// excluded there with a reason; in reverse, every spec leaf reaches provider
// surface or carries an exclusion (the reverse-drift check); and every GA
// resource carries exactly one breadth disposition (dispositions.go). The
// matcher carries zero name heuristics -- divergence is recorded, never
// guessed. Gaps gate through the burn-down baseline (baseline.go) in the
// pkg/anatomy / pkg/secretcoverage grain, surfaced by the `planton
// provider-parity` developer command and the lint.provider-parity CI lane.
package providerparity
