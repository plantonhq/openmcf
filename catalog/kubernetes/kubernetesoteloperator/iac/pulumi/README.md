# KubernetesOtelOperator Pulumi Module

Installs the OpenTelemetry Operator from the official
`opentelemetry-operator` Helm chart
(`https://open-telemetry.github.io/opentelemetry-helm-charts`) as a
single Helm release named after `metadata.name`. The typed spec renders
into chart values in `module/values.go`; the `helm_values` escape hatch
merges LAST over them with Helm `-f` semantics (maps deep-merge, later
document wins, lists replace), and the two design-load-bearing keys are
re-pinned AFTER that merge — the exact semantic twin of the Terraform
module's `helm_release` with `values = [typed, helm_values, re-pins]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace must
   already exist
2. **The four opentelemetry.io CRDs** (module-owned, unless
   `skip_crds`) — `opentelemetrycollectors`, `instrumentations`,
   `opampbridges`, `targetallocators`, applied from the tokenized files
   staged at `../crds/`, one resource per CRD keyed by the CRD's own
   `metadata.name`, with `retainOnDelete`
3. **Helm Release `<metadata.name>`** — the `opentelemetry-operator`
   chart (pinned default 0.120.0 — the newest SERVED stable chart,
   = operator appVersion 0.156.0, verified against the repository
   index), depending on both of the above

## CRD Keep-on-Uninstall (the load-bearing design)

The chart templates its opentelemetry.io CRDs as release-owned
resources — a Helm-owned install would cascade-delete every collector
declaration on uninstall. The module therefore owns the CRD lifecycle
itself:

- `crds.create: false` is pinned UNCONDITIONALLY in the rendered
  values — and re-pinned AFTER the escape-hatch merge, so no override
  can hand the CRDs back to Helm.
- `module/crds.go` applies each staged CRD file with `retainOnDelete`:
  on `pulumi destroy` the CRDs are dropped from state WITHOUT being
  deleted from the cluster, so destroying the operator never
  cascade-deletes `OpenTelemetryCollector` resources. This is the
  exact semantic twin of the Terraform module's `kubectl_manifest`
  with `apply_only = true`.
- The option must reach the ConfigGroup's CHILDREN (the actual CRD
  resources), and neither yaml package forwards ordinary resource
  options to them — the classic yaml SDK passes only parent/version
  options, and yaml/v2 children are created provider-side, beyond the
  reach of any SDK-side option (both verified in the pinned
  pulumi-kubernetes source). Hence the CLASSIC yaml `ConfigGroup`
  here, with `retainOnDelete` delivered through a resource
  TRANSFORMATION — the one mechanism the SDK propagates down the
  parent chain to in-process children.
- **The staged files are TOKENIZED renders of the pinned chart**: this
  chart TEMPLATES its CRDs — the collector CRD carries the
  `cert-manager.io/inject-ca-from` annotation and a version-conversion
  webhook `clientConfig`, both derived from the release's identity.
  The staged files carry `__PLANTON_RELEASE_NAME__` /
  `__PLANTON_NAMESPACE__` tokens, substituted in `crds.go` (and
  identically in the Terraform module), so the kept CRDs always point
  at THIS release's webhook Service and cert-manager Certificate.
- A `chart_version` upgrade re-applies the staged CRD files — re-stage
  the set together with the pinned default version.

## The Two Post-Merge Re-Pins

The deliberate exceptions to the escape hatch's last-word contract,
applied at the bottom of `buildHelmValues`:

- **`crds.create: false`** — the module owns the CRD lifecycle;
  handing them to Helm would arm the uninstall cascade-delete this
  design exists to prevent.
- **`admissionWebhooks.certManager.enabled: true`** — the kept CRDs'
  conversion trust rides cert-manager's CA injector; disabling it
  would leave module-owned CRDs pointing at a Certificate that no
  longer exists and silently break collector-CR conversion.

## Rendering Notes

- **Chart-default-matching values render only on divergence** — every
  `manager.*` key, the issuer reference, and the scheduling fields
  render only when the spec carries a value, so an empty spec installs
  the chart exactly as upstream ships it (except the two re-pins
  above and the pinned `fullnameOverride`).
- **The `manager.resources` key renders only when the spec sets it** —
  the chart ships NO default requests/limits for the manager. Helm
  deep-merges per key, so a partial spec block overrides only the
  halves it carries.
- **The collector image splits by the last `:` after the last `/`**
  (`splitImageRef` — registry ports carry `:` too). The chart renders
  `--collector-image` only when BOTH repository and tag are present; a
  repository-only override deep-merges with the chart's default tag,
  so the flag still renders.
- **`image_registry` replaces ONLY the registry part of the manager
  image** (the path stays
  `open-telemetry/opentelemetry-operator/opentelemetry-operator`); the
  default collector image the operator injects into CRs is mirrored
  via `default_collector_image` instead — collector pods pull that
  one, not this component.
- **Pull secrets render as the raw Kubernetes object list at the
  chart's TOP LEVEL** (`imagePullSecrets: [{name: ...}]`) — this chart
  does NOT nest them under `manager`.
- **Scheduling keys are TOP-LEVEL in this chart**; `nodeSelector`
  deep-merges over the chart's default `{kubernetes.io/os: linux}`, so
  the OS pin survives a spec selector.

## Values Mapping

| Spec field | Chart value |
|---|---|
| `replicas` | `replicaCount` |
| `resources` | `manager.resources` |
| `image_registry` | `manager.image.repository` (registry + upstream path) |
| `default_collector_image` | `manager.collectorImage.repository` / `.tag` (split) |
| `service_monitor_enabled` | `manager.serviceMonitor.enabled` (only true renders) |
| `webhook.issuer_ref` | `admissionWebhooks.certManager.issuerRef` |
| `image_pull_secrets` | `imagePullSecrets` (TOP-LEVEL, list of `{name}`) |
| `scheduling.node_selector` | `nodeSelector` (TOP-LEVEL) |
| `scheduling.tolerations` | `tolerations` (TOP-LEVEL) |
| `scheduling.priority_class_name` | `priorityClassName` (TOP-LEVEL) |
| (always) | `fullnameOverride: <metadata.name>` |
| (always, re-pinned last) | `crds.create: false` |
| (always, re-pinned last) | `admissionWebhooks.certManager.enabled: true` |
| `helm_values` | merged LAST (before the re-pins), Helm `-f` semantics |

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. The manager pod mounts the
cert-manager-issued webhook Secret, so an install without a working
cert-manager (or with an unpullable image from a private mirror) fails
THIS deploy with a readiness timeout instead of surfacing later as
collectors that mysteriously never reconcile. The module also fails
loudly when `metadata.name` exceeds 30 characters — the chart derives
a `<name>-controller-manager-service-cert` Secret (33-char suffix)
against the Kubernetes 63-character name cap.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`; the chart's fullname is pinned to it) |
| `webhook_service` | The operator's webhook Service (`<name>-webhook`, port 443) — where the API server sends admission reviews and CRD conversion calls |
| `webhook_cert_secret_name` | The Secret holding the webhook serving certificate (`<name>-controller-manager-service-cert`, cert-manager-issued) |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: name-budget guard → namespace → CRDs → operator
  release → output exports
- `module/crds.go`: the module-owned CRD apply (per-CRD ConfigGroups
  keyed by CRD name, the token substitution, `retainOnDelete` via
  transformation)
- `module/values.go`: typed-spec → chart values rendering (the manager
  block, the collector-image split, the issuer reference, scheduling),
  the escape-hatch merge, and the two post-merge re-pins
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version, the chart-derived webhook Service
  and cert-Secret names — kept in lockstep with the Terraform module's
  `locals.tf`
- `module/namespace.go`: the optional governance-labeled namespace
- `module/vars.go`: chart identity (repository,
  `opentelemetry-operator`), the pinned default version (0.120.0), the
  staged-CRD directory, the manager image path, the 600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
