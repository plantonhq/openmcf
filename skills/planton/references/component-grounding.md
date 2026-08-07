# Component Grounding

How to discover cloud resource kinds and learn the exact field names, types,
validation rules, and stack outputs before writing a single line of template
YAML. Grounding is offline-first: you never need a control plane or login to
author a chart correctly.

## Two instruments, one division of work

- **The multi-cloud-catalog skill is the research layer.** For discovery
  ("what exists for this provider?"), component depth (required fields,
  outputs, a validated example), wiring questions in BOTH directions ("what
  can reference this?"), and per-component judgment, read the catalog
  skill's reference pack — its own SKILL.md names which file answers which
  question and in what order. Whole-catalog questions that would take a
  dozen serial lookups resolve there in one or two file reads.
- **`planton explain` is the drill-down and recovery oracle.** Reach for it
  to drill one field's exact contract, to match a BUILD ERROR's field path
  to the schema (never guess a correction), and as the complete grounding
  path whenever no catalog pack is reachable — it is compiled into the CLI
  and works fully offline.

## The explain toolchain

| Step | Command | When |
|------|---------|------|
| Discover kinds | `planton explain --list` | You need a PascalCase kind name (e.g. `AwsVpc`) and the pack is not reachable |
| Read a kind's schema | `planton explain AwsVpc` | Before writing or fixing any manifest of that kind |
| Drill into one field | `planton explain aws-vpc.spec.instanceTenancy` | You need a field's allowed values or exact contract |
| JSON for parsing | `planton explain AwsVpc -o json` | When you need to programmatically scan fields or outputs |
| Validate a draft manifest | `planton validate draft.yaml` | After writing one resource in isolation (offline) |
| Compile the whole chart | `planton chart build <dir> -o json` | After every edit batch (needs control plane) |

`planton explain` is fully offline — the definitions AND their documentation
are compiled into the CLI binary. It is the authoritative grounding surface
for chart authors, and it covers platform API resources too (`planton explain
infra-chart` is the chart manifest's own schema, with fields the tooling
assembles labeled so you never hand-author them).
(`planton cloud-resource schema <Kind>` still works as an alias that renders
the same reference for cloud resource kinds.)

## Reading an explain report

Example (abbreviated):

```text
KIND:    AwsVpc
VERSION: aws.planton.dev/v1alpha1

DESCRIPTION:
  AwsVpcSpec defines an Amazon VPC ...

MANIFEST:
  apiVersion: aws.planton.dev/v1alpha1
  kind: AwsVpc
  metadata:
    name: <resource-name>
  spec:
    <the fields below>

RULES:
  - exactly one of cidrBlock or an IPAM-based source must be set

SPEC:
  cidrBlock <string>
    The primary IPv4 CIDR block for the VPC ...

  instanceTenancy <enum>
    values: default, dedicated

OUTPUTS (reference as valueFrom fieldPath: status.outputs.<name>):
  vpcId <string>
    The VPC's identifier ...
```

What each section means:

- **MANIFEST** — the envelope skeleton: exactly what the top level of your
  YAML must look like.
- **SPEC** — every field you can write under `spec:` in YAML. The name is
  the exact key (protojson camelCase — `cidrBlock`, not `cidr_block`).
  `-required-` means the build fails if you omit it. `values:` lists the only
  legal enum spellings — use them exactly as shown. Fields typed
  `string | valueFrom` accept either a `{value: <literal>}` or a
  cross-resource reference block (see `dependencies.md`); a bare string does
  not parse. Fields marked `(assembled)` or `(computed)` are filled by
  tooling or the platform — never write them by hand.
- **RULES** — cross-field constraints. Read these before combining fields in
  non-obvious ways.
- **OUTPUTS** — what this kind publishes after deployment. These are the
  leaves you use in `fieldPath: status.outputs.<name>` when wiring another
  resource to this one.

When the build rejects a field, re-run `planton explain <Kind>` (or drill
straight to the reported path: `planton explain <Kind>.spec.<field>`) and
match the error message's field path to the report — never guess a
correction. Manifest load errors themselves point back at `planton explain`
with the exact path to look up.

## Deeper grounding (when the explain report is not enough)

All of these live in the open-source repo `github.com/plantonhq/planton`.
Read them from a local clone when one exists; otherwise fetch raw files
directly, e.g.
`https://raw.githubusercontent.com/plantonhq/planton/main/apis/dev/planton/provider/aws/awsvpc/v1alpha1/spec.proto`
(same path shape for every artifact below). No clone and no network? The
explain report alone is sufficient for correct composition.

1. **Official presets** — worked, valid manifests for common shapes:
   `apis/dev/planton/provider/<provider>/<component>/v1/presets/`.
   Read a preset when you need a realistic starting point for a kind you
   have not authored before.
2. **Catalog pages** — human-oriented docs beside each component:
   `…/v1/catalog-page.md`. Good for topology context ("this resource usually
   sits behind a VPC") but the explain report is the field-name authority.
3. **Chart fleet** — production charts that wire many kinds together:
   `charts/<provider>/<chart-name>/` in the same repo.
   Copy wiring patterns from a chart that already solves a similar topology.
4. **Proto source** — when you need edge cases beyond the rendered report:
   `…/v1/spec.proto` and `…/v1/stack_outputs.proto` in the same component
   folder. Same field names and the same comments — the explain report is
   generated from these files, so the report normally suffices.

## Composition doctrine (per-component knowledge)

You do **not** install hundreds of component skills — one workflow skill
(this one) plus the multi-cloud-catalog research skill cover the whole
catalog. Component facts and per-component wisdom live in the catalog
skill's pack and are READ AT ANSWER TIME, never recalled from memory —
whatever you remember about a component's schema is stale by construction.
When no pack is reachable, the explain report plus a preset and a fleet
example for the kind is the complete degradation path.

## What not to use

- **Search or catalog APIs** for field names — they return metadata and links,
  not inline schemas. Use `planton explain` instead.
- **Guessing from AWS/Azure/GCP docs** — Planton kinds wrap provider APIs with
  their own field names and validation. Always ground from the explain report.
- **MCP or web UIs** — optional conveniences; the CLI commands above are the
  contract this skill is written against.
