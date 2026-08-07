# KubernetesOtelOperator Terraform Module

Installs the OpenTelemetry Operator from the official
`opentelemetry-operator` Helm chart
(`https://open-telemetry.github.io/opentelemetry-helm-charts`) as a
single Helm release named after `metadata.name`. The typed spec renders
into chart values in `locals.tf` (`local.typed_values`); the
`helm_values` escape hatch is passed as a SECOND values document the
provider merges over the first with Helm `-f` semantics, and a THIRD
document re-pins the two design-load-bearing keys LAST — the exact
semantic twin of the Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The module (not the chart) owns the CRDs** — the chart templates
  its opentelemetry.io CRDs as release-owned resources, so a
  Helm-owned install would cascade-delete every collector declaration
  on uninstall. The module pins `crds.create: false` UNCONDITIONALLY
  in the rendered values and applies the four CRD files staged at
  `../crds/` itself, one `kubectl_manifest` per CRD keyed by the CRD's
  own `metadata.name`. `skip_crds` is the bring-your-own-CRDs arm (the
  CRDs are owned elsewhere, e.g. a GitOps-managed bundle).
- **CRD keep-on-uninstall: `apply_only = true`** — the provider's
  Delete becomes a NO-OP (verified in the provider source: "When true,
  Delete is a no-op"), so destroying this module removes the CRDs from
  state WITHOUT deleting them from the cluster; collector declarations
  survive an operator uninstall. Server-side apply is REQUIRED, not
  just preferred: the collector CRD is ~418 KB, far past the
  262144-byte cap on the client-side last-applied-configuration
  annotation. The exact semantic twin of the Pulumi module's
  `retainOnDelete` on each CRD.
- **The staged CRD files are TOKENIZED renders of the pinned chart** —
  this chart TEMPLATES its CRDs: the collector CRD carries the
  `cert-manager.io/inject-ca-from` annotation and a version-conversion
  webhook `clientConfig`, both derived from the release's identity.
  The staged files carry `__PLANTON_RELEASE_NAME__` /
  `__PLANTON_NAMESPACE__` tokens, substituted in `locals.tf` (and
  identically in the Pulumi module), so the kept CRDs always point at
  THIS release's webhook Service and cert-manager Certificate.
- **Two keys are re-pinned AFTER the escape-hatch merge** — the
  deliberate exceptions to `helm_values`' last-word contract:
  `crds.create: false` (handing the CRDs to Helm would arm the
  uninstall cascade-delete this design exists to prevent) and
  `admissionWebhooks.certManager.enabled: true` (the kept CRDs'
  conversion trust rides cert-manager's CA injector; disabling it
  would leave module-owned CRDs pointing at a Certificate that no
  longer exists and silently break collector-CR conversion).
- **The release depends on the CRDs** (and the optional namespace), so
  the operator never starts against an unregistered API group.
- **The pinned default chart version is 0.120.0** — the newest SERVED
  stable chart (= operator appVersion 0.156.0, verified against the
  repository index). A `chart_version` bump must re-stage the
  `../crds/` files from the new pin — the staged files ARE the chart's
  CRDs at this version.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. The manager pod mounts the
  cert-manager-issued webhook Secret, so an install without a working
  cert-manager (or with an unpullable image) fails THIS apply with a
  readiness timeout instead of surfacing later as collectors that
  never reconcile.
- **The name budget is enforced with a precondition** —
  `metadata.name` must be 30 characters or fewer: the module pins
  `fullnameOverride` to it, the chart derives a
  `<name>-controller-manager-service-cert` Secret (33-char suffix),
  and Kubernetes caps names at 63.
- **The module (not Helm) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource
  carrying the standard governance labels;
  `helm_release.create_namespace` is always false.

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
| `helm_values` | merged as the second values document, Helm `-f` semantics |

## Rendering Quirks

- **Chart-default-matching values render only on divergence** — every
  `manager.*` key, the issuer reference, and the scheduling fields
  render only when the spec carries a value, so an empty spec installs
  the chart exactly as upstream ships it (except the two always-pinned
  keys above).
- **The `manager.resources` key renders only when the spec sets it** —
  the chart ships NO default requests/limits for the manager. Helm
  deep-merges per key, so a partial spec block overrides only the
  halves it carries.
- **The collector image splits by the last `:` after the last `/`** —
  registry ports carry `:` too (`reg.example.com:5000/x` has no tag).
  The chart renders `--collector-image` only when BOTH repository and
  tag are present; a repository-only override deep-merges with the
  chart's default tag, so the flag still renders.
- **Pull secrets render as the raw Kubernetes object list at the
  chart's TOP LEVEL** (`imagePullSecrets: [{name: ...}]`) — this chart
  does NOT nest them under `manager`.
- **`nodeSelector` deep-merges over the chart's default
  `{kubernetes.io/os: linux}`**, so the OS pin survives a spec
  selector.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.otel_operator` | `spec.create_namespace` |
| `kubectl_manifest.crds` (one per staged CRD, keyed by CRD name) | unless `spec.skip_crds` |
| `helm_release.otel_operator` | always |

## Usage

```bash
planton tofu apply --manifest otel-operator.yaml
```

## Local Development

```bash
tofu init -backend=false
tofu validate
tofu plan -var-file=terraform.tfvars.json
tofu apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto). The spec arrives from the proto→tfvars converter in
snake_case with the `namespace` foreign key (KubernetesNamespace)
resolved to a literal string before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`; the chart's fullname is pinned to it) |
| `webhook_service` | The operator's webhook Service (`<name>-webhook`, port 443) — where the API server sends admission reviews and CRD conversion calls |
| `webhook_cert_secret_name` | The Secret holding the webhook serving certificate (`<name>-controller-manager-service-cert`, cert-manager-issued) |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
chart identity and pinned default version (0.120.0), same
`metadata.name` release name and 30-character budget, same values
rendering (the divergence-only manager block, the collector-image
split, the two post-merge re-pins), same module-owned tokenized-CRD
posture (`apply_only` here, `retainOnDelete` there), same atomic/wait
posture, same outputs.
