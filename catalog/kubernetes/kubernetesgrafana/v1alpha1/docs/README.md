# KubernetesGrafana: Research and Design

## Introduction

KubernetesGrafana deploys one standalone Grafana from the official
`grafana` Helm chart as a single Helm release named after
`metadata.name`. The typed spec renders into chart values; every
credential rides a Kubernetes Secret (chart-owned or referenced) and
none ever lands in rendered values; and the deployment grain is one
Grafana instance-or-HA-group, with its state posture — the decision
that actually separates a demo from a production deployment — typed
as three explicit arms.

## Chart Identity: the Repository Move

The canonical Grafana chart no longer lives where most of the
internet says it does. The `grafana/helm-charts` repository's served
index stalls at chart 10.5.x, while the live line is served from
`grafana-community/helm-charts` (pinned here: chart 12.8.0, shipping
Grafana 13.1.1). The authoritative tell: kube-prometheus-stack's own
Chart.yaml declares its Grafana dependency against the
grafana-community repository. Both modules install from the community
repository; the `chart_version` field documents the move so a future
maintainer bumping the pin verifies against the right index.

The chart itself is Apache-2.0. The Grafana application it references
is AGPLv3 — referenced, never distributed: the customer's cluster
pulls the image directly from Grafana's registry, and nothing of the
application is redistributed by this catalog.

## The Ephemeral-State Trap, Typed

Grafana stores UI-authored dashboards, users, preferences, alert
silences — everything hand-made — in an embedded SQLite database on
the pod's local disk. The chart's default deployment is stateless:
every pod restart erases all of it. This is the single most common
Grafana-on-Kubernetes failure, and it is not an edge case — a node
drain during routine cluster maintenance is enough.

The spec types the three honest postures instead of hiding them
behind chart values:

- **Ephemeral (default).** Correct when everything is provisioned as
  code — declared datasources, sidecar-discovered dashboards — so a
  fresh pod reconstructs the whole UI from configuration. Wrong the
  moment a human builds a dashboard by hand.
- **`storage`.** A ReadWriteOnce PVC under the SQLite file: one
  stateful instance. CEL blocks combining it with `replicas > 1` — a
  RWO volume cannot follow two pods, and two Grafanas sharing one
  SQLite file corrupt it.
- **`database`.** External Postgres/MySQL for Grafana's state — the
  durable path and the HA requirement. CEL enforces
  `replicas > 1 → database`: scaled replicas without a shared
  database silently split sessions and dashboards across pods (each
  keeps its own SQLite), which looks like data loss to users. The
  `host` is a foreign key defaulting to KubernetesPostgres (its
  read-write endpoint) so the database composes as a real dependency
  edge; the password rides an existing Secret through environment
  expansion.

## The Credentials Contract

The chart's admin-secret template is lookup-stable: it generates a
random password only when the Secret does not already exist, so the
credential is created ONCE at first install and survives upgrades
(the same class as Qdrant's chart-owned API-key Secret). That makes
the Secret chart-owned — the module materializes nothing and the
import travels with the Helm release. `admin_secret` is the
existing-Secret arm (`admin.existingSecret` + key overrides) for
teams that manage credentials externally; the outputs echo whichever
Secret is live, so downstream consumers read one field either way.

The chart hard-refuses to render secret material into its config
ConfigMap, which the spec turns into a design rule: every credential
surface (admin, database, datasource basic auth, SMTP) is a Secret
reference expanded through environment variables at runtime.

## Provisioning as Code, and the Sidecar Contract

Datasources and dashboards have a declarative path that survives pod
loss without any persistence:

- **`datasources`** render Grafana's datasource provisioning file.
  Each URL is a `StringValueOrRef` foreign key whose default kind is
  KubernetesKubePrometheusStack, resolving to its exported Prometheus
  endpoint — the one-line wiring from dashboards to cluster metrics,
  with a real dependency edge underneath. `json_data` covers
  type-specific settings (rendered in clear — the comment bans
  credentials there); basic-auth passwords ride Secret references.
- **The dashboard sidecar** (k8s-sidecar, on by default) watches the
  whole cluster for ConfigMaps labeled `grafana_dashboard: "1"` and
  loads them live. This is the composition contract: any component or
  team ships dashboards by creating a labeled ConfigMap — no coupling
  to this resource's spec, no redeploy.
- **`community_dashboards`** import grafana.com dashboards by numeric
  ID at install, each bound to a declared datasource by name.
  Revisions are pinnable because a moving latest revision can change
  a dashboard under its users.

## Grafana 13 Plugin Reality

Grafana 13 moved several once-core datasource plugins (elasticsearch,
cloudwatch among them) out of the core image; installs that silently
worked on 11.x need them declared. The `plugins` field carries
`<id>` or `<id> <version>` entries; the modules also enable the
chart's bundled-plugin shadowing knob so plugin installs succeed on
images with a read-only plugin directory. Plugin installs and
community-dashboard imports both download at pod start — an air-gap
consideration the field comments carry.

## What Was Deliberately Left Out

- **Embedded ingress.** Exposure composes from first-class kinds
  (KubernetesIngress, Gateway API kinds) over the exported service
  handle; `server.root_url` is the one Grafana-side field exposure
  needs (OAuth redirects and rendered links embed it).
- **LDAP/OAuth provider blocks.** High-cardinality, provider-shaped
  configuration that belongs in grafana.ini — it rides `helm_values`
  until demand shapes a typed surface.
- **Grafana Alerting configuration.** Alert rules and contact points
  are data-plane objects with their own provisioning file format;
  the escape hatch carries them for now.

## Configuration Split: Typed Fields vs helm_values

Typed: identity (namespace, chart version, replicas, image), state
(storage/database), credentials (admin, database, datasource,
SMTP — all Secret-backed), provisioning (datasources, sidecar,
community dashboards, plugins), server identity and auth posture,
scheduling, and the ServiceMonitor toggle. `helm_values` (merged
LAST, Helm `-f` semantics, identical engines) carries the long tail:
LDAP/OAuth, the image renderer, alerting provisioning, notifiers,
extra mounts and sidecars. Never secrets.

## Outputs Vocabulary

The outputs export the composition handles: namespace, release name,
the Service name (= the resource name — fullname is pinned), the
in-cluster endpoint, the live admin-Secret name (chart-owned or
echoed), and a port-forward command for workstation access before any
exposure is composed.
