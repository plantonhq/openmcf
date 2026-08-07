# KubernetesIngressNginx Pulumi Module

Installs the ingress-nginx controller from the official Helm chart
(`ingress-nginx` at `https://kubernetes.github.io/ingress-nginx`) as a real
Helm release. The typed spec renders into chart values in
`module/values.go`; the `helm_values` escape hatch merges LAST over them
with Helm `-f` semantics (maps deep-merge, later document wins, lists
replace) — the exact semantic twin of the Terraform module's `helm_release`
with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must exist.
   Labels are stamped on this module-created namespace only — never
   injected into the chart's own resources; Helm owns those
2. **Helm Release** — named after `metadata.name` (NOT a fixed chart name:
   multiple controller instances per cluster — public + internal traffic
   splits, each owning its own IngressClass — are a first-class upstream
   pattern). The chart's `fullnameOverride` is pinned to the release name
   so every chart object carries a deterministic, manifest-derived name
   (`<name>-controller`, `<name>-controller-internal`), which also isolates
   leader election per instance (the chart's election ID defaults to
   `<fullname>-leader`)
3. **Controller Service read** (`Service.Get`) — after the release, the
   controller Service is read back for the load-balancer address outputs.
   Gated on the `load_balancer` service type: for `node_port`/`cluster_ip`
   there is no LB status to read and the address outputs stay empty by
   design. Reads run through the provider's read path (no awaiters), so
   this never blocks

The IngressClass controller identifier derives automatically in
`module/locals.go`: the chart default `k8s.io/ingress-nginx` for class
`nginx`, otherwise `k8s.io/<class-name>` — additional controllers isolate
without the user inventing a vocabulary.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (300s
timeout) for the release's resources to become ready — a controller that
never starts (bad image, unschedulable pod, webhook certgen failure) fails
THIS deploy, not the first Ingress. Helm's readiness check on a
LoadBalancer-type Service also waits for the cloud LB address, so on
clusters WITHOUT a cloud LB controller (kind, bare metal) a `load_balancer`
service type times out loudly here — deliberate: use `node_port`/host
access on such clusters, and the failure names the real problem instead of
leaving a silently Pending entry point.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the controller is installed in |
| `release_name` | Helm release name (equals `metadata.name`; controller resources are named `<release>-controller`) |
| `ingress_class_name` | The IngressClass this controller owns — what KubernetesIngress resources reference |
| `controller_service_name` | The controller's external Service — the traffic entry point |
| `internal_service_name` | The internal Service (empty unless `service.internal.enabled`) |
| `load_balancer_ip` | External IP of the cloud LB (providers that populate an IP; empty otherwise) |
| `load_balancer_hostname` | External hostname of the cloud LB (providers that populate a DNS name; empty otherwise) |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release → controller Service read →
  exports
- `module/values.go`: typed-spec → chart values rendering (ingress-class
  identity, replicas-vs-autoscaling exclusivity, service shaping,
  host-network DNS policy, default-TLS `extraArgs` flag, default-backend
  image split), escape-hatch merge
- `module/namespace.go`: conditional namespace creation
- `module/locals.go`: resolved names (release, ingress class + controller
  value, controller/internal Service names) — kept in lockstep with the
  Terraform module's `locals.tf`
- `module/vars.go`: chart identity and the pinned default chart version
  (kept byte-identical with the Terraform module — cross-engine chart-name
  drift installs different software per engine)
