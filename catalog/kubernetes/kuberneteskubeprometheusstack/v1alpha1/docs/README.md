# KubernetesKubePrometheusStack: Research and Design

## Introduction

KubernetesKubePrometheusStack deploys one cluster monitoring stack
from the official `kube-prometheus-stack` Helm chart
(https://prometheus-community.github.io/helm-charts, pinned 87.19.1,
pairing Prometheus Operator v0.92.1) as a single Helm release named
after `metadata.name`. The typed spec renders into chart values; the
deployment grain is the whole bundle — operator, Prometheus,
Alertmanager, exporters, rules and the bundled Grafana — because that
is the unit the chart composes, tests and versions together.

## Why the Operator-Based Stack Is the Baseline

Running Prometheus as a raw container is easy; running it correctly
is not. Scrape configuration is a living document (every new workload
changes it), rule files need validation before they reach a running
server (a syntax error stalls config reload for everything), and HA,
retention and storage all have server-lifecycle implications. The
Prometheus Operator answers this with CRDs: ServiceMonitor and
PodMonitor declare scraping next to the workload that owns it,
PrometheusRule declares alerting next to the service it watches, and
the operator compiles all of them into a validated server
configuration continuously.

The kube-prometheus-stack chart is the canonical packaging of that
model: the operator, a Prometheus and an Alertmanager declared
through it, kube-state-metrics (object-state metrics), node-exporter
(host metrics), a curated rule set encoding years of upstream
operational knowledge, and a Grafana pre-loaded with the matching
dashboards. It is the single most-installed monitoring artifact in
the Kubernetes ecosystem, and every catalog component's
`service_monitor_enabled` toggle assumes its CRDs.

One stack per cluster is the grain. The CRDs are cluster-scoped
singletons; a second stack must skip CRD installation and fence its
discovery — supported (`skip_crds`, `discovery`), but an advanced
posture, not the default.

## Chart Identity and the CRD Subchart

The chart's CRD posture is unusual and drives real design decisions.
The CRDs do NOT sit in a top-level `crds/` directory: they ship in a
local `crds` SUBCHART (`charts/crds/`, condition `crds.enabled`,
default true) whose manifests sit in the subchart's own `crds/`
directory. Helm semantics for that layout: installed ONCE at first
install, NEVER touched by chart upgrades, KEPT on uninstall.

Consequences designed in:

- **Keep-on-uninstall is the contract, not an accident.**
  ServiceMonitors and rules across the cluster survive removal of the
  stack. The E2E destroy phase asserts the CRDs remain — tolerating
  their disappearance would hide a regression in the designed
  posture.
- **Chart upgrades do not upgrade CRDs.** Crossing operator versions
  without new CRDs is how operator upgrades break silently. The chart
  ships an optional `upgradeJob` (a pre-install/pre-upgrade hook Job
  that server-side-applies the CRD bundle — SSA matters, these CRDs
  are among the largest in the ecosystem and exceed the client-side
  annotation limit). The spec exposes it as `crd_upgrade_job` with
  the teaching on the `chart_version` field.
- **`skip_crds` is a bring-your-own-CRDs arm**, for the second-stack
  and GitOps-managed-CRD postures — with the CRDs absent the install
  fails; it is never a "lighter" install.

## The Naming Contract and the 26-Character Budget

The modules pin `fullnameOverride` to `metadata.name`, which makes
every child name deterministic: `<name>-operator`,
`prometheus-<name>-prometheus` (the operator names the StatefulSet
`prometheus-<CR name>` and the chart names the CR
`<fullname>-prometheus`), `alertmanager-<name>-alertmanager`,
`<name>-grafana`. Deterministic names are what the stack outputs, the
import maps and the verifiers stand on.

The chart truncates its fullname at 26 characters — its own headroom
for the longest child name it derives. A silently truncated fullname
would break the naming contract everywhere downstream, so both
modules fail loudly at plan/preview time when `metadata.name` exceeds
26 characters.

## Discovery: Cluster-Wide by Default, Deliberately

The chart's own default discovers only monitor/rule objects labeled
by its release — upstream's most-tripped-over behavior ("my
ServiceMonitor is ignored"). This component inverts that: the default
is cluster-wide discovery (the selector-nil-uses-helm-values family
of values set to discover everything), because the catalog's
composition model depends on it — any component's
`service_monitor_enabled` toggle, and any user-authored monitor, must
light up without touching the stack's own manifest.
`discovery: release_managed_only` restores the chart's fence for
multi-tenant clusters with deliberate ownership boundaries.

## Persistence: PVCs by Default, Against the Chart's Grain

The chart's default storage for both Prometheus and Alertmanager is
an emptyDir: every metric, silence and notification-log entry
vanishes on pod restart. That default is wrong for anything beyond a
demo, so this component provisions PersistentVolumeClaims by default
(`disk_size` per replica, `storage_class` by literal or
KubernetesStorageClass reference) and makes the chart's posture an
explicit opt-in (`ephemeral: true`), with CEL keeping the two
mutually exclusive. The CEL tolerates the platform's defaulting
middleware stamping default sizes onto every manifest — an ephemeral
manifest that never declared a disk size must stay expressible.

Prometheus retention is two-dimensional: `retention` (time) and
`retention_size` (bytes) — samples trim when either is reached.
`retention_size` must sit below the volume size: when the volume
itself fills, Prometheus crash-loops instead of trimming.

## Remote Write: the Cross-Cloud Surface

Remote write is how this stack composes with managed backends, and
each cloud authenticates differently — all four arms are typed, from
the Prometheus CRD's remoteWrite schema:

- **Basic auth** (Grafana Cloud, Mimir): the CRD wants BOTH halves
  from Secrets. Passwords ride existing Secrets you reference; the
  non-secret usernames are materialized by the module into a
  deterministic `<name>-remote-write-auth` Secret (keyed
  `username-<i>` per destination) so the CR has a Secret to point at
  without the manifest carrying one.
- **Bearer token** (API-token backends): an existing Secret
  reference.
- **SigV4** (Amazon Managed Prometheus): keyless by default — the
  pod's ambient identity signs (IRSA on EKS), with `role_arn` for
  assume-role and a static-keys arm (both-or-neither CEL) when
  ambient identity is unavailable.
- **Azure AD managed identity** (Azure Monitor managed Prometheus):
  the workload-identity client ID and cloud selector.

`enable_remote_write_receiver` covers the inverse composition: this
Prometheus as the aggregation point other agents push to.

## The Managed-Cloud Scraper Truth

The chart scrapes the full self-managed control plane by default. On
EKS/GKE/AKS the controller-manager, scheduler and etcd run on
provider-internal machines the cluster network can never reach —
leaving those scrapers on produces permanently-down targets and
alerts that fire forever. kube-proxy is its own trap: its metrics
port binds to localhost on several managed platforms, and clusters
running Cilium's kube-proxy replacement have nothing to scrape at
all.

The spec types every toggle (`control_plane_scrapers`) with the
per-cloud truth in the field comments, and pairs it with
`default_rules.disabled_groups` — a disabled scraper without its rule
group disabled leaves alerts that can never fire truthfully. The
presets carry the exact managed-cloud set.

## The Bundled Grafana vs the Standalone Kind

The chart's Grafana subchart arrives pre-wired: a datasource pointing
at this Prometheus and the full curated dashboard set. That is the
product for the one-stack-one-cluster case, so it is a first-class
typed block, on by default. The boundary with `KubernetesGrafana` is
role, not preference: the bundled Grafana serves THIS stack's
dashboards; the standalone kind is the composition hub — many
datasources, external database state, HA replicas. The standalone
kind's datasource URL is a foreign key resolving to this stack's
`prometheus_endpoint`, so the two compose with a real dependency
edge.

Bundled-Grafana credentials follow the chart's own mechanism: a
random admin password generated ONCE at first install
(lookup-stable across upgrades) in the chart-owned `<name>-grafana`
Secret, or an existing Secret referenced through `admin_secret`.
Either way the Secret name lands in the stack outputs and no
credential ever appears in rendered values.

## Secret Discipline

Every credential surface is a Secret reference: remote-write
passwords/tokens/keys, the bundled-Grafana admin Secret, image-pull
Secret names. The single module-materialized Secret
(`<name>-remote-write-auth`) carries only usernames — non-secret
values that the CRD nevertheless requires from a Secret.
`alertmanager.config_yaml` is the one seam where a user could inline
a credential (webhook URLs, API keys); the field comment teaches the
`_file`-reference pattern instead, and `helm_values` carries the
Secret mount.

## Configuration Split: Typed Fields vs helm_values

Typed: everything above — the surfaces with operational consequences
that deserve teaching and validation. `helm_values` (merged LAST,
Helm `-f` semantics, identical engines) carries the long tail:
Thanos sidecar/ruler scale-out, windows monitoring, scrape classes,
additional Alertmanager template files, per-component
securityContexts, ingress-per-replica. Never secrets.

## Outputs Vocabulary

The outputs export the full downstream composition surface: the
namespace and release name, Service names and in-cluster endpoints
for Prometheus, Alertmanager and the bundled Grafana (empty when the
half is disabled), the bundled Grafana's admin-Secret name, and
port-forward commands for workstation access. These are the handles
the standalone Grafana's datasource FK, the catalog's
`service_monitor_enabled` arms, and composed exposure all stand on.
