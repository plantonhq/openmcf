# providerparity

Measures the catalog against the Terraform provider it deploys through:
for each cloud provider, what the pinned Terraform provider can configure
versus what the catalog's kinds expose and consume.

**This is provider parity, not cross-engine parity.** Cross-engine parity
(`pkg/iac/MODULE_PARITY.md`, `PARITY-EXCEPTION:` comments, the component
audit's parity focus) is one kind's Terraform and Pulumi modules
implementing the same contract identically. Provider parity — this package —
is the catalog covering the full configurable surface of the pinned
provider, with every omission recorded as a decision.

## Design

Detection is fully mechanical; judgment is recorded once, where tooling
reads it. Three independent censuses feed the measurement:

1. **The provider describes itself.** `schemas/` holds one committed,
   distilled `providers schema -json` artifact per provider at the exact
   release its pin resolves to (produced by `distiller/`, refreshed via
   `make generate-provider-schemas`). No provider source parsing. Parity is
   always declared against a **named version** — the pin — never against
   "latest"; that is what makes a freshness promise well-defined. Bumping a
   pin replaces the artifact, and the resulting check failures are the
   migration work list.
2. **The contract side comes from descriptors.** `spec_census.go` walks the
   kind registry via `pkg/crkreflect` and enumerates every spec leaf field —
   the same walk shape as `pkg/secretcoverage` (StringValueOrRef as one
   leaf, map/list handling, recursion guard). Text/regex counting of protos
   undercounts nested specs and is banned for parity numbers.
3. **The consumption side comes from the modules.** `module_census.go` scans
   **every** `*.tf` file of each kind's `iac/tf/` module for resource
   declarations and provider pins — never `main.tf` alone, because modules
   may split resources across sibling files and a partial scan undercounts
   silently.

`report.go` joins the three into the single aggregation every consumer
renders from (CLI output, CI gate, the public parity report generator), so
numbers can never disagree between surfaces.

## Layout notes

- Artifact identity lives in the content (`provider` field), never the
  filename; loaders reject two artifacts for one provider so the pin stays
  unambiguous.
- Artifacts are **loaded from the repo tree, not `go:embed`-ded**: they are
  repo-audit data that grows with every provider, and embedding them would
  carry megabytes into the customer-facing binary for data only repo gates
  read (contrast `pkg/protodocs`, whose small index the CLI serves at
  runtime and therefore embeds).
- The core is provider-agnostic: `google` is configuration in the Makefile
  target, not code. Adding a provider is a new `--provider` flag and its
  catalog's registry entries.

## Running

```bash
go test ./pkg/providerparity/          # gates + hermetic fixtures + live measurement
make generate-provider-schemas        # refresh schemas/ after a pin change
```
