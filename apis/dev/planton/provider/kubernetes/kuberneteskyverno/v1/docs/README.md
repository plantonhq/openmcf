# Kubernetes Kyverno — design notes

## Grain

One resource = one Kyverno Helm release (official `kyverno` chart,
kyverno.github.io index; chart 3.x = Kyverno 1.18+). One install per
cluster — webhook registration and the policy CRDs are cluster-global.
`fullnameOverride` pins the fullname to `metadata.name`. Chart-truth
on naming: the controller Deployments derive from the CHART name
(constant `kyverno-admission-controller`, `kyverno-background-controller`,
...), while the webhook Service (`<fullname>-svc`), the runtime
ConfigMap (the fullname) and the pre-delete hook Job
(`<fullname>-hook-pre-delete` — the longest fullname suffix, 16
characters) derive from the fullname; the name budget is therefore 47
characters and both engines fail loudly past it.

## The composition seam

- **In:** policies (ClusterPolicy / Policy and the policies.kyverno.io
  v1 families) applied as `KubernetesManifest` resources or GitOps
  after the engine installs.
- **Out:** `admission_service_name` (the webhook backend Service),
  `config_map_name` (the runtime skip-list to inspect when a resource
  is unexpectedly skipped or policed), namespace and release name.
- **Certificates:** the `certificates.cert_manager` arm composes
  `KubernetesCertManager` (issuer reference → `KubernetesClusterIssuer`),
  applied to BOTH webhook servers — admission and cleanup.
- **Metrics:** `metrics.service_monitor` fans a ServiceMonitor across
  all four controllers (requires the kube-prometheus-stack CRDs).

## Runtime-registered webhooks (the load-bearing lifecycle fact)

The chart templates no webhook configurations; the admission
controller registers and maintains them at runtime, scoped to the
installed policies (autoUpdateWebhooks). Uninstall runs the chart's
pre-delete hook (`webhooks_cleanup_enabled`, default on) AND a
module-owned cleanup: a destroy-ordered sentinel ConfigMap the release
depends on, whose delete (after the release — and the admission pods —
are gone) removes `kyverno-*` webhook configurations by label. The
module cleanup is load-bearing, not belt-and-suspenders: at the pinned
release the chart hook's delete-webhooks helper targets the wrong API
group and leaves every ValidatingWebhookConfiguration behind, and pods
still dying can re-register webhooks the hook did delete. The default
per-rule failure policy is Fail — a stranded webhook blocks matched
admissions cluster-wide, which is why the hook toggle carries the
warning and the unstick command lives on the spec's top comment. The
Terraform cleanup is a destroy-time provisioner (kubectl on the runner,
kubeconfig from the process environment); the Pulumi twin is a
BeforeDelete resource hook, which fires only when destroy runs with
`--run-program` (both Planton runners pass it).

## Config semantics

The chart's default resourceFilters list (control-plane namespaces,
Kyverno's own resources) is EDITED, never replaced:
`resource_filters_include` appends, `resource_filters_exclude` removes
matching default entries. `webhook_exclude_namespaces` rebuilds the
webhook namespaceSelector and re-includes kube-system by construction.
`exclude_groups`, when declared, REPLACES the chart default
(system:nodes) — the field comment says to re-include it.

## Cross-engine parity

The Terraform and Pulumi modules render byte-identical chart values.
Admission-controller container resources sit under
`container.resources` (init + main container split); the other three
controllers take `resources` directly. Every feature flag renders the
chart's nested `{enabled: ...}` shape. The Terraform locals use the
null-prune idiom throughout — this module's own offline plan gate
caught three ternary type-unification sites during authoring.

## Deliberate exclusions

The reports-server subchart and openreports API (default-off,
separate storage posture), the grafana dashboard subchart, per-hook
image pinning, PodDisruptionBudgets, tracing/metering blocks, and
apiPriorityAndFairness — reachable through `helm_values`, never the
primary interface. The chart's credential-creating `imagePullSecrets`
map is deliberately not modeled (credentials belong in
`KubernetesSecret`, not chart values); `image_pull_secrets` maps to
`existingImagePullSecrets` (names only).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
