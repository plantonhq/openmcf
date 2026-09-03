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
2. **The operator CRDs** — derived from the pinned chart; see CRD ownership below
3. **Helm Release `<metadata.name>`** — the `solr-operator` chart
   (pinned default 0.9.1; chart versions carry NO `v` prefix while the
   operator image does — chart 0.9.1 ships operator v0.9.1)

## CRD Ownership and the Keep-on-Uninstall Mechanism

The chart carries its CRDs on both of Helm's surfaces: the three
solr.apache.org CRDs in its `crds/` directory (Helm installs those once
and never upgrades them) and the ZookeeperCluster CRD templated by the
bundled zookeeper-operator subchart behind `zookeeper-operator.crd.create`
(Helm would delete it with the release). The module OWNS all four through
the catalog's shared package for charts that carry CRDs, `keptcrds`: the
pinned chart is rendered in-process (Helm SDK, client-only) with the
release's own values plus the subchart's CRD switch turned on, the
CustomResourceDefinition documents are kept, stamped with
`planton.ai/crd-source-chart` and `planton.ai/crd-source-version`, and
applied one `yaml.ConfigGroup` per CRD keyed by the CRD's own
`metadata.name`:

- `solrclouds.solr.apache.org`
- `solrbackups.solr.apache.org`
- `solrprometheusexporters.solr.apache.org`
- `zookeeperclusters.zookeeper.pravega.io` (rendered only when the
  bundled zookeeper-operator is installed, the chart default)

The release installs with `SkipCrds` and `zookeeper-operator.crd.create:
false` re-pinned AFTER the `helm_values` merge, so Helm never touches the
CRDs whichever way a user override points. `crds.install: false` is the
bring-your-own-CRDs arm (the CRDs are owned elsewhere, e.g. a
GitOps-managed bundle).

**Keep mechanism (Pulumi):** every applied CRD carries `retainOnDelete`
(delivered through a resource transformation, the one option the yaml
SDK propagates to a ConfigGroup's children), so `pulumi destroy` removes
the operator but NEVER the CRDs — SolrCloud / SolrBackup /
ZookeeperCluster resources are never cascade-deleted cluster-wide.
`crds.keep_on_uninstall: false` turns it off. The CRDs ride a dedicated
upsert provider (server-side apply with force) so a reinstall re-adopts
the kept CRDs instead of failing on their existence. The Terraform twin
is `apply_only = true` on `kubectl_manifest`.

**A schema downgrade is refused:** before anything registers, the
module lists the CRDs it has stamped on the cluster; a `chart_version`
below the version they carry fails the preview with what was observed,
what it means, and the next step (pin the higher version, or delete the
CRDs deliberately). A version that is not published in the repository
index is refused the same way, before anything is created.

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
- `module/main.go`: namespace → release values → derived CRDs
  (`keptcrds.Apply`) → operator release → output exports
- `module/values.go`: typed-spec → chart values rendering and the
  escape-hatch merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version, deployment name (the chart fullname
  replay) — kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository, `solr-operator`), the
  pinned default version (0.9.1), the CRD render override, the 600s
  timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
