# KubernetesSolrOperator Terraform Module

Installs the Apache Solr Operator from the official `solr-operator` Helm
chart (`https://solr.apache.org/charts`) as a single Helm release named
after `metadata.name`. The typed spec renders into chart values in
`locals.tf` (`local.typed_values`); the `helm_values` escape hatch is
passed as a SECOND values document the provider merges over the first
with Helm `-f` semantics — the exact semantic twin of the Pulumi
module's `buildHelmValues` + `mergeMaps`.

## CRD Ownership and the Keep-on-Uninstall Mechanism

Unlike most operator charts, the `solr-operator` chart ships NO CRDs —
upstream publishes them as separate release artifacts. The module OWNS
them: the four files staged at `../crds` are applied before the release,
one `kubectl_manifest` per CRD keyed by the CRD's own `metadata.name`:

- `solrclouds.solr.apache.org`
- `solrbackups.solr.apache.org`
- `solrprometheusexporters.solr.apache.org`
- `zookeeperclusters.zookeeper.pravega.io` (the bundled
  zookeeper-operator dependency's CRD)

**Keep mechanism (Terraform):** `apply_only = true` makes the
alekc/kubectl provider's Delete a NO-OP (verified in the provider
source), so `terraform destroy` removes the operator release but leaves
the CRDs — and therefore every SolrCloud / SolrBackup /
ZookeeperCluster resource cluster-wide — untouched. The Pulumi twin is
`retainOnDelete` on each CRD. `server_side_apply = true` matters twice
over: the SolrCloud CRD's schema exceeds the client-side
last-applied-configuration annotation size limit, and SSA lets a
restaged (upgraded) CRD file apply as an in-place update.

The staged files carry upstream license-comment headers and the
ZookeeperCluster file a doubled leading `---`; `locals.crd_documents`
splits on separator LINES (`"\n---"` — never a bare `"---"` substring
split, which would corrupt the `rw-rw----` / `--- TODO` text the CRD
schemas embed in descriptions) and drops the comment-only fragments
with a `can(yamldecode)` filter.

Because the module owns the ZookeeperCluster CRD, the bundled
subchart's own CRD switch is pinned off —
`zookeeper-operator.crd.create: false` is the ONE value rendered
unconditionally.

**Version note:** the staged CRD files match the pinned default chart
version (0.9.1). The operator is pre-1.0 — when upgrading
`chart_version` across a minor version, restage the matching CRD files.

## Module Behavior

- **The release name is `metadata.name`.**
- **The chart version has NO `v` prefix** — chart 0.9.1 ships operator
  image v0.9.1; the pinned default (0.9.1) must exist as a SERVED chart
  in the repository index at https://solr.apache.org/charts.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An operator that never becomes
  ready (an unpullable image from a private mirror is the classic case)
  fails THIS apply with a readiness timeout instead of surfacing later
  as SolrCloud resources that never reconcile.
- **The module (not Helm) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource
  carrying the standard governance labels;
  `helm_release.create_namespace` is always false.

## Values Mapping

| Spec field | Chart value | Notes |
|---|---|---|
| `replicas` | `replicaCount` | on presence |
| `watch_namespaces` | `watchNamespaces` | COMMA-JOINED string (the chart's format); empty = watch ALL namespaces |
| `zookeeper_operator.install` | `zookeeper-operator.install` | on presence (chart default true) |
| `zookeeper_operator.use_existing` | `zookeeper-operator.use` | only when true (chart default false; ignored when install=true) |
| — | `zookeeper-operator.crd.create` | ALWAYS `false` — the module owns the ZookeeperCluster CRD |
| `leader_election_enabled` | `leaderElection.enable` | on presence |
| `metrics_enabled` | `metrics.enable` | on presence |
| `mtls.client_cert_secret` | `mTLS.clientCertSecret` | value-or-ref resolved to the secret name before Terraform runs |
| `mtls.ca_cert_secret` | `mTLS.caCertSecret` | value-or-ref resolved to the secret name before Terraform runs |
| `mtls.ca_cert_secret_key` | `mTLS.caCertSecretKey` | on presence (chart default `ca-cert.pem`) |
| `mtls.insecure_skip_verify` | `mTLS.insecureSkipVerify` | on presence (chart default true) |
| `mtls.watch_for_updates` | `mTLS.watchForUpdates` | on presence (chart default true) |
| `resources` | `resources` | on presence (chart ships none — the operator is lightweight) |
| `node_selector` | `nodeSelector` | when non-empty |
| `tolerations` | `tolerations` | when non-empty |
| `image_pull_secret` | `image.imagePullSecret` | SINGULAR string — the chart accepts exactly one |
| `image.repository` | `image.repository` | when non-empty |
| `image.tag` | `image.tag` | when non-empty |
| `helm_values` | (merged LAST) | Helm `-f` semantics, identical on both engines |

Note the dash in `zookeeper-operator`: it is the SUBCHART name,
addressed as a quoted map key in `locals.tf` — never dot-path syntax.

## Rendering Quirks

- **Chart-default-matching values render only on divergence** — an
  empty spec installs the chart exactly as upstream ships it, except
  for the pinned `zookeeper-operator.crd.create: false`.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.solr_operator` | `spec.create_namespace` |
| `kubectl_manifest.solr_operator_crds` (×4, keyed by CRD name) | always |
| `helm_release.solr_operator` | always |

## Usage

```bash
planton tofu apply --manifest solr-operator.yaml
```

## Local Development

```bash
tofu init -backend=false
tofu validate
```

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto). The spec arrives from the proto→tfvars converter in
snake_case with the `namespace` foreign key (KubernetesNamespace) and
the mTLS secret foreign keys (KubernetesSecret) resolved to literal
strings before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |
| `deployment_name` | Name of the operator Deployment — the chart's fullname template replayed: the release name itself when it contains "solr-operator", otherwise `<release>-solr-operator`, truncated to 63 chars |

The `deployment_name` derivation assumes the chart's default
`nameOverride`/`fullnameOverride` (the typed spec sets neither); a
`helm_values` override of either changes the real name without being
reflected in the output.

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
chart identity and pinned default version (0.9.1), same `metadata.name`
release name, same CRD ownership keyed by CRD name (apply_only here,
retainOnDelete there), same values rendering (the comma-joined watch
scope, the always-rendered `zookeeper-operator.crd.create: false`, the
singular image pull secret), same atomic/wait posture, same outputs.
