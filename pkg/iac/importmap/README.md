# importmap — import-ID recipes, loaded and resolved

State import needs two pieces of knowledge per resource: **which address**
the engine tracks it under, and **which identifier** the provider imports it
by. Addresses are never authored — the engine enumerates them per spec at
import time (a read-only preview lists them), which handles
module-constructed names, repeated (`for_each`/`count`) resources, and
conditional resources by construction. Only the identifier knowledge is
authored, in two tiers:

| Tier | Kind | Location | Owns |
|------|------|----------|------|
| Provider | `ProviderImportCatalog` | `apis/dev/planton/provider/{provider}/aa_import/catalog.yaml` | Import-ID **format** per resource type (`"{bucket}"`, `"{vpc_id}"`), plus `config_only_attributes` — attributes that exist only in IaC configuration and can never round-trip through import |
| Component | `ComponentImportMap` | `{component}/v1/iac/import-map.yaml` | The **value source** per `{placeholder}`: metadata.name, a spec field, a stack output, a pasted ARN's part, or the enumerated address's instance key — with "where to find this" guidance for anything only the user can supply |

Both are proto-backed KRM documents (`iac.planton.dev/v1`, protos under
`apis/dev/planton/iac/`), parsed through `pkg/protobufyaml` like the E2E
profiles.

## Correctness is machine-proven, never review-trusted

- **Offline** (`conformance_test.go`, runs in `make test`): every resource
  type a mapped component's OpenTofu module declares has an id_format;
  every placeholder those formats reference is declared by the component
  map; `from_spec_field` paths resolve to scalar leaves on the kind's spec
  proto; `from_stack_output` keys are real StackOutputs fields. Enrollment
  is allowlisted (`mappedKinds`), mirroring the `variables.tf` drift guard.
- **Live** (the E2E `IMPORT-RT` phase, opt-in via
  `PLANTON_E2E_IMPORT_ROUNDTRIP=1`): deploy the fixture, set its state
  aside, re-import every resource *blind* through these recipes, and
  require the follow-up plan to propose no real change — in-place updates
  are tolerated only for declared `config_only_attributes`
  (e.g. `aws_s3_bucket.force_destroy`, engine delete behavior with no
  cloud-side existence). The destroy that follows runs through the
  re-imported state, proving it fully owns the resources.

## Adding a kind

1. Add the component's resource types (with import-ID formats) to the
   provider catalog if absent; declare any config-only attributes.
2. Author `{component}/v1/iac/import-map.yaml` naming each placeholder's
   derivations (prefer derivable sources; `where_to_find` is mandatory when
   nothing derives).
3. Enroll the component in `mappedKinds` (conformance guard) and run the
   live round-trip lane for it before shipping.
