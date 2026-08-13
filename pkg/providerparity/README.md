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

## Total accounting

On top of the measurement sits the judgment layer: the accounting
(`accounting.go`) that makes silence impossible in both directions.

**Depth (per kind).** Every configurable, non-deprecated argument of every
consumed resource must be exact-matched to a spec field, mapped by the
kind's manifest, or excluded there with a reason; and — the reverse
direction — every spec leaf field must reach a provider argument or carry a
recorded exclusion. The reverse check is what catches reverse drift: a spec
field the pinned provider no longer serves. The matcher carries **zero name
heuristics** — an argument either matches its derived spec path exactly or
its divergence is recorded — because any clever normalization would
eventually false-match and silently overclaim, the one failure mode a
parity instrument must make impossible. Two argument classes sit outside
accounting by standing rule: Terraform's own machinery (`id`, `timeouts.*`
— the tool's surface, not the resource's) and schema-deprecated arguments.

**The per-kind manifest** (`catalog/<provider>/<kind>/iac/provider-parity.yaml`,
beside `import-map.yaml`; file presence IS enrollment, no allowlist):

```yaml
resources:
  google_storage_bucket:
    mappings:                       # renames and subtree mappings; exact
      - spec: spec.bucket_name      # matching resumes below a mapped
        arg: name                   # subtree, so entries stay
      - spec: spec.lifecycle_rules  # O(divergences), never O(fields)
        arg: lifecycle_rule
    exclusions:                     # deliberately unmodeled, reason mandatory;
      - arg: lifecycle_rule.condition.send_age_if_zero
        reason: derived by the module from the optional field's presence
      - arg: legacy_inline_block    # an exclusion may name a block: one
        reason: fixed by the module # judgment covers its whole subtree
  google_storage_bucket_iam_member:
    specRoot: spec.iam_members      # secondary resource: the spec subtree
    exclusions:                     # that instantiates it
      - arg: bucket
        reason: wired by the module to the bucket it creates
  google_project_service:
    internal: API enablement plumbing -- arguments are module decisions
  azapi_resource:                   # externally specified: the kind's primary
    external: >-                    # surface lives outside every loaded schema
      raw-ARM surface at a pinned type@api-version, admitted by recorded
      decision (name the pin and the tracking issue here)
specExclusions:                     # spec fields with no provider counterpart
  - field: spec.platform_only
    reason: platform-level concept
```

Spec references use the census's `spec.`-rooted proto-name dot paths (the
same path grammar as `pkg/secretcoverage`). Mappings are many-to-many in
both directions: one spec map field may be realized by several arguments
(the name/value idiom — record one mapping per argument), and one
map-typed argument may realize several spec fields (the FAN-IN idiom —
record one mapping per spec field on the same arg; the argument counts
once, every mapped field is covered in reverse). Manifest entries
referencing surface that no longer exists (after a pin bump or spec
change) are findings, not warnings — after a pin bump those failures ARE
the migration work list. Malformed manifests fail hard at load: judgment
half-read is worse than judgment absent.

**Externally specified resources.** `external` is `internal`'s mirror: where
internal says "module plumbing below the spec", external says "the kind's
primary surface, whose contract is pinned OUTSIDE every loaded provider
schema" — a raw-API resource (e.g. azapi's `azapi_resource`) admitted per
kind by a recorded decision naming an exact `type@api-version`. No argument
walk is possible (there is no schema to walk), and a kind consuming only
external/internal resources runs no reverse spec walk either: its spec's
depth is accounted against the external contract, where the admission
records it. The judgment ratchets toward native support: the moment a
loaded schema serves the resource at the pin, the external entry becomes a
staleness finding — the exit-to-native migration's mechanical reminder. A
kind mixing schema-served and external resources keeps the reverse walk;
its external-fed spec fields carry `specExclusions` naming the external
resource.

**Externally specified resources.** `external` is `internal`'s mirror: where
internal says "module plumbing below the spec", external says "the kind's
primary surface, whose contract is pinned OUTSIDE every loaded provider
schema" — a raw-API resource (e.g. azapi's `azapi_resource`) admitted per
kind by a recorded decision naming an exact `type@api-version`. No argument
walk is possible (there is no schema to walk), and a kind consuming only
external/internal resources runs no reverse spec walk either: its spec's
depth is accounted against the external contract, where the admission
records it. The judgment ratchets toward native support: the moment a
loaded schema serves the resource at the pin, the external entry becomes a
staleness finding — the exit-to-native migration's mechanical reminder. A
kind mixing schema-served and external resources keeps the reverse walk;
its external-fed spec fields carry `specExclusions` naming the external
resource.

**Breadth (per GA resource).** Every GA resource carries exactly one
disposition. Two classes are computed — `modeled` (the module census proves
consumption) and `iam-covered` (the `*_iam_member/binding/policy` pattern,
covered by the owning kinds' additive `iam_members` fields) — plus
schema-flagged deprecations. The rest is recorded judgment in the
dispositions ledger (`dispositions/<schema>.yaml`): `composed`,
`model-planned`, `deferred`, and doc-level `excluded-deprecated`, reason
mandatory. Computed classes always win: a ledger entry shadowed by one
(the resource became consumed, matches the IAM pattern, or is
schema-flagged deprecated) is a staleness finding — recorded judgment that
duplicates derived judgment is judgment nobody re-evaluates.

**The gate.** Findings group by kind/resource into the burn-down baseline
(`baseline.yaml`, the `pkg/anatomy` / `pkg/secretcoverage` shape): new
findings fail, fixed-but-still-listed entries fail, the list only ever
shrinks truthfully. The baseline holds only what has not been judged YET —
the permanent record lives in the manifests and the ledger. CI lane:
`.github/workflows/lint.provider-parity.yaml`.

## The public parity report

The same accounting renders the published parity page — the artifact that
makes the coverage claim verifiable instead of trusted:

- **Home:** `catalog/<provider>/terraform-parity.md`, a COMMITTED generated
  file (the `reference-index.md` grain). The site build's copy script
  carries it to the docs site and links it from the provider index.
- **Content:** the measurement baseline (pinned schema versions, pin
  distribution), the per-kind depth accounting with live-proof status, the
  breadth disposition totals, and the full enumerated per-resource record.
  Deterministic markdown — no timestamps; freshness is the named pin.
- **PROVEN** is joined mechanically from the component E2E profiles
  (`e2eproof.go`): a kind is proven when its profile is green with BOTH
  provisioners validated. Claim language stays "built for 100% Terraform
  parity" — the page never asserts achieved parity; it shows the measured
  state.
- **Drift gate:** each page embeds its generation parameters, and the drift
  test regenerates every committed page and requires byte-identical output —
  a hand edit or a stale page fails CI. Regenerate with
  `make generate-provider-parity-report`. A provider enrolls by generating
  its page once with `--write-report`; file presence is enrollment.

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
make generate-provider-parity-report  # regenerate the committed parity page(s)

# The developer CLI (from the repo root; registered in the standalone
# binary beside e2e -- never in the embedded engine set, because it reads
# repo-tree artifacts that do not exist outside a checkout):
planton provider-parity --provider gcp --ga-schema google              # summary
planton provider-parity --provider gcp --ga-schema google --kind GcpGcsBucket
planton provider-parity --provider gcp --ga-schema google --check     # the CI gate
planton provider-parity --provider gcp --ga-schema google --output json
planton provider-parity --provider gcp --ga-schema google --write-report
```
