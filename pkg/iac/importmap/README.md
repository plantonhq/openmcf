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
| Component | `ComponentImportMap` | `{component}/v1/iac/import-map.yaml` | The **value source** per `{placeholder}`: metadata.name (optionally with a literal suffix, for convention-named satellites like `<name>-hpa`), a spec field, a stack output, a pasted ARN's part, the enumerated address's instance key, or a module-hardcoded literal (a typed-CR module's apiVersion/kind) — with "where to find this" guidance for anything only the user can supply |

An id_format has two optionality forms, with deliberately different literal
handling: `{name?}` renders as the empty string when unresolved and KEEPS its
surrounding literal delimiters (the DynamoDB `index:{index_name?}/` shape),
while a bracketed segment group like `[//{namespace}]` disappears WHOLESALE —
delimiters included — when its placeholder does not resolve. The group form
exists for composed IDs whose provider rejects a trailing delimiter:
`kubectl_manifest` imports namespaced CRs as
`apiVersion//kind//name//namespace` and cluster-scoped CRs as the 3-part
`apiVersion//kind//name`, never with a trailing `//`.

When a module declares SEVERAL resources of one type whose ID placeholders
carry different values (a control-plane module installing multiple Helm
releases — each release's `{release_name}` is a different constant), a value
declaration can be scoped to one Terraform logical resource name via
`tofu_resource_name`. At resolve time the declaration scoped to the address's
logical name wins; addresses without a scoped declaration fall back to the
unscoped declaration of the same placeholder. The offline conformance guard
rejects a scope that names no real module resource (a typo'd scope would
otherwise silently fall back and import the wrong resource).

Both are proto-backed KRM documents (`iac.planton.dev/v1`, protos under
`apis/dev/planton/iac/`), parsed through `pkg/protobufyaml` like the E2E
profiles.

## Correctness is machine-proven, never review-trusted

- **Offline** (`conformance_test.go`, runs in `make test`): every resource
  type a mapped component's OpenTofu module declares has an id_format;
  every placeholder those formats reference is declared by the component
  map; `from_spec_field` paths resolve to scalar leaves on the kind's spec
  proto; `from_stack_output` keys are real StackOutputs fields.
- **Live** (the E2E `IMPORT-RT` phase, opt-in via
  `PLANTON_E2E_IMPORT_ROUNDTRIP=1`): deploy the fixture, set its state
  aside, re-import every resource *blind* through these recipes, and
  require the follow-up plan to propose no real change — in-place updates
  are tolerated only for declared `config_only_attributes`. One
  framework-owned pruning happens before declarations apply: attributes the
  plan itself marks WHOLLY unknown post-apply (`after_unknown` literal true —
  plugin-framework providers recompute `id`/`metadata`-class computed
  attributes on every in-place update) are never treated as drift; partially
  unknown objects and all known drift stay under the oracle
  (e.g. `aws_s3_bucket.force_destroy`, engine delete behavior with no
  cloud-side existence) and `write_normalized_attributes`. A tolerance may
  be a dotted sub-path (`"spec.update_strategy"`) when a provider's importer
  fails to read back one nested block — the oracle then prunes exactly that
  sub-path from both plan sides and requires the remainder identical, so
  sibling drift still fails. The destroy that follows runs through the
  re-imported state, proving it fully owns the resources.

## Enrollment is the file itself

The import-map file's presence is the single enrollment signal everywhere:
the offline conformance guard auto-discovers every
`{component}/v1/iac/import-map.yaml` (`DiscoverComponentImportMaps`), the
E2E round-trip gate keys off the same file, and the platform's catalog
pipeline bundles the same file into the import wizard's data. There is
deliberately no allowlist: an authored map that ships to users while
dodging validation — or passes validation while never shipping — is a
divergence this single signal makes impossible.

**The merge discipline that keeps the ledger honest: a new or changed
import map does not merge without a green live round-trip lane.** The
offline guard proves recipes are well-formed; only the live lane proves
they are CORRECT (they import the real resource, not a plausible one).

## Live-proven ledger

The round-trip lane is the only honest enforcement of "proven" — this
table is the human-readable record of which kinds have run it, so a reader
never has to guess which maps are correctness-proven versus merely
well-formed. A recipe for a resource no scenario deploys is offline-validated
only (noted per row); the lane proves exactly what the fixtures exercise.

| Component | Live round-trip | Tolerated declared updates / notes |
|-----------|-----------------|------------------------------------|
| `awss3bucket` | 2026-07-09, both scenarios (13 resources incl. `for_each` keys + the `{bucket}:{intelligent_tiering_name}` composite) | `aws_s3_bucket.force_destroy` (config-only) |
| `awsvpc` | 2026-07-10, minimal | `aws_vpc_ipv4_cidr_block_association` is user-supplied by design (no blind derivation), offline-validated only |
| `awssecuritygroup` | 2026-07-10, rules-rich | `aws_security_group.revoke_rules_on_delete` (config-only) |
| `awsecrrepo` | 2026-07-10, full-surface (repository + lifecycle policy) | `aws_ecr_repository.force_delete` (config-only) |
| `awssubnet` | 2026-07-10, minimal + routed (subnet + route table + `{subnet_id}/{route_table_id}` association composite) | — |
| `awsinternetgateway` | 2026-07-10, smoke | — |
| `awsnatgateway` | 2026-07-10, minimal | — |
| `awsiamrole` | 2026-07-10, smoke (role + inline `{role_name}:{inline_policy_name}` + attachment `{role_name}/{managed_policy_arn}`, both via `for_each` keys) | `aws_iam_role.force_detach_policies` declared (provider-documented config-only; scenario does not set it) |
| `awsdynamodb` | 2026-07-10, on-demand-full-surface (table + resource policy + table- and GSI-level contributor insights incl. the optional `{index_name?}` segment and account id) | `aws_dynamodb_resource_policy.policy`/`revision_id` (write-normalized); `aws_dynamodb_kinesis_streaming_destination` not composed by the scenario, offline-validated only |
| `awssqsqueue` | 2026-07-10, fifo-full-surface (queue URL id) | — |
| `awssnstopic` | 2026-07-10, standard-topic (topic ARN id) | `aws_sns_topic_data_protection_policy` not composed by the scenario, offline-validated only |
| `awskmskey` | 2026-07-10, minimal (key UUID + `alias/...` via `for_each` key) | `aws_kms_key.deletion_window_in_days` declared (provider-documented config-only; scenario does not set it) |
| `kubernetesdeployment` | 2026-07-21, all 13 scenarios (workload + service + HPA + PDB + env Secret + created namespace) | `kubernetes_deployment_v1.wait_for_rollout`, `kubernetes_service_v1.wait_for_load_balancer`, `kubernetes_secret_v1.wait_for_service_account_token` (config-only) |
| `kubernetesstatefulset` | 2026-07-21, all 13 scenarios | same config-only knobs; `kubernetes_stateful_set_v1` `spec.update_strategy` (write-normalized dotted sub-path: the provider importer does not read the block back) |
| `kubernetesdaemonset` | 2026-07-21, all 5 scenarios | `kubernetes_daemon_set_v1.wait_for_rollout` (config-only) |
| `kubernetesjob` | 2026-07-21, all 5 scenarios | — |
| `kubernetescronjob` | 2026-07-21, all 4 scenarios | — |
| `kubernetesservice` | 2026-07-21, all 13 scenarios | `kubernetes_service_v1.wait_for_load_balancer` (config-only) |
| `kubernetesingress` | 2026-07-21, all 4 scenarios (incl. the composed real-backend fixture) | `kubernetes_ingress_v1.wait_for_load_balancer` (config-only) |
| `kubernetesnetworkpolicy` | 2026-07-21, all 4 scenarios | — |
| `kubernetespersistentvolumeclaim` | 2026-07-22, all 3 scenarios (incl. the composed StorageClass fixture and the pinned-empty-class static-binding shape) | `kubernetes_persistent_volume_claim_v1.wait_until_bound` (config-only; deliberately false in the module) |
| `kubernetesstorageclass` | 2026-07-22, both scenarios | — |
| `kubernetesresourcequota` | 2026-07-22, all 3 scenarios (governed-namespace proves the quota + companion LimitRange pair) | — |
| `kubernetespriorityclass` | 2026-07-22, both scenarios | — |
| `kubernetespoddisruptionbudget` | 2026-07-22, both scenarios | — |
| `kuberneteshorizontalpodautoscaler` | 2026-07-22, all 3 scenarios (incl. the composed Deployment target fixture) | — |
| `kubernetescertmanager` | 2026-07-22, both scenarios (Helm release + created namespace) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kuberneteshelmrelease` | 2026-07-22, all 3 scenarios (HTTPS repo, OCI registry, values-override) | `helm_release` install-time attributes (repository, values/set/set_sensitive, lifecycle knobs — Helm does not persist how a release was installed; provider-documented, declared config-only). Computed attributes the plan marks wholly after-unknown (id, metadata) are pruned by the oracle itself, not declared. |
| `kubernetesexternaldns` | 2026-07-22, all 3 scenarios (Helm release + created namespace; the credential-Secret arms are exercised offline by the hack manifest — the kind lanes run credential-less by design) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kubernetesexternalsecretsoperator` | 2026-07-22, both scenarios (Helm release + created namespace) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kubernetesingressnginx` | 2026-07-22, all 3 scenarios (Helm release + created namespace; multi-instance scenario re-imports alongside a live sibling controller) | `helm_release` install-time attributes (config-only, see the catalog row); the module's controller-Service data source is a read, not an imported resource |
| `kubernetesmetricsserver` | 2026-07-22, both scenarios (Helm release + created namespace) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kubernetesgatewayclass` | 2026-07-22, both scenarios (the composed 3-part cluster-scoped ID — the `[//{namespace}]` segment group drops) | `kubectl_manifest` provider-side knobs (config-only) + `yaml_body`/`yaml_body_parsed` (write-normalized: the importer stores the stripped live object), see the catalog row |
| `kubernetesgateway` | 2026-07-22, all 3 scenarios (incl. the composed GatewayClass-fixture FK) | same `kubectl_manifest` tolerances |
| `kuberneteslistenerset` | 2026-07-22, both scenarios (incl. the composed Gateway-fixture parent FK) | same `kubectl_manifest` tolerances |
| `kuberneteshttproute` | 2026-07-22, all 3 scenarios (incl. the composed Gateway + Service fixture chain) | same `kubectl_manifest` tolerances |
| `kubernetesgrpcroute` | 2026-07-22, both scenarios | same `kubectl_manifest` tolerances |
| `kubernetestcproute` | 2026-07-22, both scenarios | same `kubectl_manifest` tolerances |
| `kubernetesudproute` | 2026-07-22, both scenarios | same `kubectl_manifest` tolerances |
| `kubernetestlsroute` | 2026-07-22, both scenarios | same `kubectl_manifest` tolerances |
| `kubernetesreferencegrant` | 2026-07-22, both scenarios | same `kubectl_manifest` tolerances |
| `kubernetesclusterissuer` | 2026-07-22, self-signed scenario (cluster-scoped 3-part composed ID) | same `kubectl_manifest` tolerances; credential Secrets derive blind via `from_address_key` (the module keys them by Secret name) — the self-signed scenario materializes none, so those recipes are offline-validated only |
| `kubernetesissuer` | 2026-07-22, self-signed + composed CA-chain scenarios | same `kubectl_manifest` tolerances; credential-Secret recipes offline-validated only (see `kubernetesclusterissuer`) |
| `kubernetescertificate` | 2026-07-22, minimal + full-surface scenarios (real issuance fixture chain) | same `kubectl_manifest` tolerances |
| `kubernetesclustersecretstore` | 2026-07-22, fake-backend scenario (cluster-scoped 3-part composed ID) | same `kubectl_manifest` tolerances; the conditional `<name>-credentials` Secret derives via `from_address_key` / `from_metadata_name_suffix` — the fake backend declares no credentials, so that recipe is offline-validated only |
| `kubernetessecretstore` | 2026-07-22, fake-backend scenario | same as `kubernetesclustersecretstore` |
| `kubernetesexternalsecret` | 2026-07-22, both scenarios (incl. the real sync-loop fixture chain) | same `kubectl_manifest` tolerances |
| `kubernetesistio` | 2026-07-22, minimal + ambient scenarios (20 resources: FOUR `tofu_resource_name`-scoped Helm releases incl. the ambient cni/ztunnel pair, 15 module-owned CRDs blind-derived via `from_address_key` — the module keys each `kubectl_manifest` by the CRD's own name — and the created namespace) | `helm_release` install-time attributes (config-only, see the catalog row) + `kubectl_manifest` tolerances |
| `kubernetesdestinationrule` | 2026-07-22, minimal scenario | same `kubectl_manifest` tolerances |
| `kubernetesserviceentry` | 2026-07-22, minimal scenario | same `kubectl_manifest` tolerances |
| `kubernetespeerauthentication` | 2026-07-22, minimal scenario | same `kubectl_manifest` tolerances |
| `kubernetesrequestauthentication` | 2026-07-22, minimal scenario | same `kubectl_manifest` tolerances |
| `kubernetesauthorizationpolicy` | 2026-07-22, minimal + behavioral-deny scenarios (re-imported alongside a live mesh) | same `kubectl_manifest` tolerances |
| `kubernetestelemetry` | 2026-07-22, minimal scenario | same `kubectl_manifest` tolerances |
| `kubernetesenvoyfilter` | 2026-07-22, minimal scenario | same `kubectl_manifest` tolerances |
| `kubernetescilium` | 2026-07-23, both scenarios on the cilium-cni cluster profile (Helm release; kube-system installs import no namespace) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kuberneteskeda` | 2026-07-23, both scenarios (Helm release + created namespace; behavioral-scaling re-imports with the live scale-target fixture present) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kubernetesbackendtlspolicy` | 2026-07-23, all 3 scenarios (4-part namespaced composed ID; composed-target resolves the Service-fixture FK) | same `kubectl_manifest` tolerances |
| `kubernetesclusterautoscaler` | 2026-07-23, minimal scenario on the kwok arm (Helm release + created namespace) | `helm_release` install-time attributes (config-only, see the catalog row) |
| `kubernetesvelero` | 2026-07-23, both scenarios (Helm release + created namespace; behavioral re-imports alongside the live MinIO fixture) | `helm_release` install-time attributes (config-only, see the catalog row) |

Kinds where an import map is **deliberately not applicable** (recorded so
absence is never mistaken for an oversight):

| Component | Why no import map |
|-----------|-------------------|
| `kubernetesmanifest` | The kind's state is arbitrary user YAML — there is no per-kind resource schema for a blind round-trip oracle to compare against; each deployment's resource set is defined by the manifest itself. Adopting existing raw resources is done by pasting their YAML into `spec.manifest_yaml` and applying (server-side apply takes ownership of unmanaged fields). |
| `kubernetesgatewayapicrds` (and other multi-document CRD-bundle installers) | The module applies a fetched multi-document bundle with one `kubectl_manifest` per document, keyed by the SPLIT INDEX — a positional key, never identity. The composed import ID needs each document's own apiVersion/kind/name triple, and no single derivation can honestly feed three placeholders per positional address. Adopting an existing CRD install is done by applying over it (server-side apply with force_conflicts takes ownership). |
| `kubernetesistiobasecrds` | Same multi-document CRD-bundle class as `kubernetesgatewayapicrds` — the istio/base CRD bundle applies one `kubectl_manifest` per split document, positionally keyed. |

## Adding a kind

1. Add the component's resource types (with import-ID formats) to the
   provider catalog if absent; declare any config-only attributes. When a
   type is listable via Cloud Control (`aws cloudcontrol list-resources
   --type-name ...` succeeds without extra parameters -- verify
   empirically, never assume), also declare its `cloud_control_type_name`:
   it is the scan-side correspondence account scanning and the mapping
   eval harness (`pkg/iac/mappingeval`) translate through.
2. Author `{component}/v1/iac/import-map.yaml` naming each placeholder's
   derivations (prefer derivable sources; `where_to_find` is mandatory when
   nothing derives). The file's presence enrolls the component in every
   check — nothing else to register.
3. Run the live round-trip lane
   (`PLANTON_E2E_IMPORT_ROUNDTRIP=1 go test -tags=e2e -run Test<Kind>_Terraform ./e2e/...`)
   before merging, and record the result in the ledger above.
