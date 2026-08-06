# Kubernetes Argo CD — design notes

## Grain

One resource = one Argo CD control plane from the official `argo-cd`
chart (argoproj index). The release is named after `metadata.name` and
`fullnameOverride` pins every child name (`<name>-server`,
`<name>-application-controller`, ...), so the exported outputs are
deterministic. Names are capped at 37 characters — the chart appends
component suffixes up to 26 characters and truncates silently at 63,
which would break the naming contract; both engines fail loudly
instead. The initial-admin Secret's name is fixed by the APPLICATION
(not the chart), so one generated-password instance runs per namespace.

## The composition seam

- **In:** Applications/AppProjects/ApplicationSets as plain CRs
  (`KubernetesManifest`, charts, the UI/CLI) once the control plane
  runs; repo and SSO credentials as labeled Secrets composed next to it.
- **Out:** `server_service` / `server_kube_endpoint` for composed
  exposure, `initial_admin_secret_name` for first login, and the
  port-forward command for workstations.

## Application-owned vs chart-owned vs module-owned

The admin password is APPLICATION-owned (generated at first start into
`argocd-initial-admin-secret`); the Redis auth Secret is CHART-owned
(the redis-secret-init hook Job generates it, lookup-stable); the
namespace is MODULE-owned. Nothing credential-bearing is module-owned
by design — this component transports no secret material at all.

## The config surfaces

`configs.cm` (OIDC/dex/exec/reconciliation), `configs.params`
(`server.insecure`), `configs.rbac` and `configs.repositories` render
only DECLARED entries, so the chart's own defaults stay authoritative.
The cm helper coerces values through `toString` — booleans are safe in
ConfigMap data.

## Cross-engine parity

Both engines render byte-identical chart values. Argo CD's image is
the COMBINED `global.image.repository` form (registry included — no
split mapping); dex and redis images stay chart-default under an
override (reachable via `helm_values` when mirroring those too). The
single `service_monitors_enabled` toggle fans out to every component
with a serviceMonitor at the pin (controller, server, repo server,
applicationSet, notifications, dex); the commit server has none.

## Deliberate exclusions

Per-component env/volumes/ingress, notification notifiers and
templates, server extensions, network policies, and the
`dynamicClusterDistribution` Deployment mode — reachable through
`helm_values`, never the primary interface. The chart's own
ingress/httproute blocks are never modeled: exposure composes from
first-class kinds.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
