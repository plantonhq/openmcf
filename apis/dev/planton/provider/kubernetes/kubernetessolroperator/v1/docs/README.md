# KubernetesSolrOperator: Research and Design

## Introduction

KubernetesSolrOperator installs the Apache Solr Operator from the
official `solr-operator` Helm chart (https://solr.apache.org/charts,
pinned 0.9.1) as a single Helm release named after `metadata.name`.
The operator is the ENGINE of the Solr story in this catalog:
KubernetesSolr declares `SolrCloud` custom resources, and this
operator reconciles them into node StatefulSets with shard-aware
rolling updates, replica-moving scale operations, TLS, basic-auth
bootstrap, and registered backup repositories. `SolrBackup` and
`SolrPrometheusExporter` are the operator's other reconciled types.

## The Deployment Landscape

SolrCloud is a coordinated distributed system: shards, replicas, and
an external ZooKeeper quorum for cluster state. Restarting nodes in
the wrong order takes shards offline; scaling down without moving
replicas loses them. The Apache Solr project ships its own operator to
encode that Day-2 expertise — which is why the catalog splits the
concern in two: this kind installs the engine once, KubernetesSolr
declares each cluster.

### The bundled zookeeper-operator

Every SolrCloud needs a ZooKeeper ensemble. The chart declares the
zookeeper-operator as a dependency (version 0.2.15, from
charts.pravega.io, conditioned on `zookeeper-operator.install` —
default true), which is what makes KubernetesSolr's "provided
ZooKeeper" arm work out of the box: the Solr operator asks the
zookeeper-operator to provision an ensemble per cluster.

The second-install posture matters: the zookeeper-operator carries
fixed-name cluster-scoped RBAC, so a second install in the same
cluster CONFLICTS. When one already runs, set `install: false` and
`use_existing: true` (the chart's `use` knob) so the Solr operator
still knows a zookeeper-operator is present. When every Solr cluster
will use an external ensemble, `install: false` alone suffices.

## The CRD Lifecycle: Module-Owned, Keep-on-Uninstall

Unlike most operator charts, the solr-operator chart ships NO CRDs —
upstream publishes them as separate release artifacts. The modules
therefore own all four:

- the three `solr.apache.org` CRDs — SolrCloud, SolrBackup,
  SolrPrometheusExporter — and
- the `zookeeper.pravega.io` ZookeeperCluster CRD of the bundled
  dependency.

Mechanics, per engine:

- **Terraform**: `kubectl_manifest` with `apply_only = true` — the
  provider's Delete is a NO-OP (verified in the alekc/kubectl provider
  source), so `terraform destroy` removes the operator release but
  leaves the CRDs (and therefore every SolrCloud / SolrBackup /
  ZookeeperCluster resource cluster-wide) untouched.
- **Pulumi**: one classic yaml `ConfigGroup` per CRD with
  `retainOnDelete` delivered through a resource TRANSFORMATION — the
  one mechanism the SDK propagates to the ConfigGroup's in-process
  children.
- **Both engines key each CRD resource by the CRD's OWN
  `metadata.name`**, so state addresses stay stable regardless of file
  naming or ordering. Document splitting is line-anchored ("---" at
  column 0) because the CRD schemas embed "---" inside description
  text — a substring split would corrupt them.
- **`zookeeper-operator.crd.create` is pinned false unconditionally**
  — the one always-rendered chart value: the bundled subchart must
  never install its own ZookeeperCluster CRD and put it under Helm's
  delete-on-uninstall lifecycle.
- **Server-side apply matters twice over**: the SolrCloud CRD's schema
  blows past the client-side last-applied-configuration annotation
  size limit, and SSA field ownership lets a restaged (upgraded) CRD
  file apply as an in-place update.
- **The release depends on the CRDs** — the operator's controllers
  watch these types immediately, and the bundled zookeeper-operator
  refuses to start without the ZookeeperCluster CRD present.

### The chart/CRD version pairing

Chart and operator versions move together: chart 0.9.1 ships operator
v0.9.1 (the chart version carries no `v` prefix; the operator image
and CRD artifacts do). The operator is pre-1.0 — minor versions can
change the CRD API — so a `chart_version` bump means restaging the
matching CRD files alongside it; SSA turns the restage into an
in-place update.

## The mTLS Client Identity

When a KubernetesSolr cluster enforces client certificates
(`tls.client_auth: Need`), every caller must present one — including
the operator itself, whose probes and reconciliation calls otherwise
fail. The `mtls` block is that identity: `client_cert_secret`
(tls.crt/tls.key — REQUIRED whenever the block is declared, because an
mtls block without a client certificate would render nothing and
silently leave the operator without an identity), an optional CA
secret (key `ca-cert.pem` by chart default), `insecure_skip_verify`
(chart default true — the operator calls pods by IP, which rarely
matches certificate SANs), and `watch_for_updates` (restart the
operator when the certificate rotates, chart default true).

## Design Decisions

- **The install is blocking.** The Helm release waits for the operator
  to become Available (atomic, 600s timeout, cleanup on fail): an
  unpullable image fails THIS apply with a readiness timeout instead
  of surfacing later as SolrCloud resources that mysteriously never
  reconcile.
- **Chart-default-matching values render only on divergence** — the
  rendered values stay minimal on both engines, with the one
  `zookeeper-operator.crd.create: false` exception above.
- **`watch_namespaces` renders as the chart's comma-joined string**
  (the chart splits it back apart in its helpers), not a YAML list.
  Unlike some operators, the watched namespaces need only exist by the
  time SolrCloud resources appear in them.
- **`deployment_name` replays the chart's fullname helper** — the
  release name, suffixed with the chart name when absent, truncated to
  63 characters. A `helm_values` override of
  `nameOverride`/`fullnameOverride` would break it.
- **The chart accepts exactly ONE image pull secret**
  (`image.imagePullSecret`, a singular string) — the spec models it as
  a single field rather than pretending a list works.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `solr-operator` at https://solr.apache.org/charts | Pinned 0.9.1 (spec default) |
| Operator image | `apache/solr-operator:v<chart_version>` | The `v` prefix is on the image, not the chart |
| CRDs | SolrCloud, SolrBackup, SolrPrometheusExporter (`solr.apache.org`), ZookeeperCluster (`zookeeper.pravega.io`) | Module-owned, keep-on-uninstall; restage on chart bumps |
| Bundled dependency | zookeeper-operator 0.2.15 (charts.pravega.io) | `install` default true; `crd.create` pinned false |
| Watch scope | all namespaces (chart default) | `watch_namespaces` fences to an explicit set |
| Deployment | the chart's fullname helper | Exported as `deployment_name` |

## IaC Twins

Pulumi (`module/crds.go` + `module/values.go`) and Terraform
(`main.tf` + `locals.tf`) render identical chart values, the same
module-owned CRD set keyed the same way, and the same
keep-on-uninstall posture. Keep the typed-value rendering and the
fullname derivation in lockstep.
