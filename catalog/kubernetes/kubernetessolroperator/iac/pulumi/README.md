# KubernetesSolrOperator Pulumi Module

Installs the Apache Solr Operator from the official `solr-operator` Helm
chart (`https://solr.apache.org/charts`) as a single Helm release named
after `metadata.name`. The typed spec renders into chart values in
`module/values.go`; the `helm_values` escape hatch merges LAST over them
with Helm `-f` semantics (maps deep-merge, later document wins, lists
replace) — the exact semantic twin of the Terraform module's
`helm_release` with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace must
   already exist
2. **The four operator CRDs** — see CRD ownership below
3. **Helm Release `<metadata.name>`** — the `solr-operator` chart
   (pinned default 0.9.1; chart versions carry NO `v` prefix while the
   operator image/CRD artifacts do — chart 0.9.1 ships operator v0.9.1)

## CRD Ownership and the Keep-on-Uninstall Mechanism

Unlike most operator charts, the `solr-operator` chart ships NO CRDs —
upstream publishes them as separate release artifacts. The module OWNS
them: the four files staged at `../crds` are applied before the release,
one `yaml.ConfigGroup` per CRD keyed by the CRD's own `metadata.name`:

- `solrclouds.solr.apache.org`
- `solrbackups.solr.apache.org`
- `solrprometheusexporters.solr.apache.org`
- `zookeeperclusters.zookeeper.pravega.io` (the bundled
  zookeeper-operator dependency's CRD)

**Keep mechanism (Pulumi):** every applied CRD carries `retainOnDelete`,
so `pulumi destroy` removes the operator but NEVER the CRDs — SolrCloud
/ SolrBackup / ZookeeperCluster resources are never cascade-deleted
cluster-wide. Because the yaml SDK forwards only version/pluginDownloadURL
to a ConfigGroup's children, the option is delivered through a resource
TRANSFORMATION (inherited parent→child) — see `module/crds.go`. The
Terraform twin is `apply_only = true` on `kubectl_manifest` (the
provider's Delete is a no-op).

Because the module owns the ZookeeperCluster CRD, the bundled subchart's
own CRD switch is pinned off — `zookeeper-operator.crd.create: false` is
the ONE value rendered unconditionally (it must never fall under Helm's
delete-on-uninstall lifecycle).

**Version note:** the staged CRD files match the pinned default chart
version (0.9.1). The operator is pre-1.0 — when upgrading
`chart_version` across a minor version, restage the matching CRD files
(the module applies with server-side apply semantics, so a restage is an
in-place update).

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
| `mtls.client_cert_secret` | `mTLS.clientCertSecret` | value-or-ref resolved to the secret name |
| `mtls.ca_cert_secret` | `mTLS.caCertSecret` | value-or-ref resolved to the secret name |
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

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. An operator that never becomes ready
(an unpullable image from a private mirror is the classic case) fails
THIS deploy with a readiness timeout instead of surfacing later as
SolrCloud resources that mysteriously never reconcile.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |
| `deployment_name` | Name of the operator Deployment — the chart's fullname template replayed: the release name itself when it contains "solr-operator", otherwise `<release>-solr-operator`, truncated to 63 chars |

The `deployment_name` derivation assumes the chart's default
`nameOverride`/`fullnameOverride` (the typed spec sets neither); a
`helm_values` override of either changes the real name without being
reflected in the output.

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → CRDs → operator release → output exports
- `module/crds.go`: staged-CRD application keyed by CRD name, with the
  retain-on-delete transformation and the comment-header/document-
  separator handling
- `module/values.go`: typed-spec → chart values rendering and the
  escape-hatch merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version, deployment name (the chart fullname
  replay) — kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository, `solr-operator`), the
  pinned default version (0.9.1), the 600s timeout, the staged CRDs
  directory
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
