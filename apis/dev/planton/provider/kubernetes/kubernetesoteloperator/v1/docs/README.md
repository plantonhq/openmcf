# KubernetesOtelOperator: Research and Design

## Introduction

KubernetesOtelOperator installs the OpenTelemetry Operator from the
official `opentelemetry-operator` Helm chart
(https://open-telemetry.github.io/opentelemetry-helm-charts, pinned
0.120.0) as a single Helm release named after `metadata.name`. The
operator is the ENGINE of the OpenTelemetry story in this catalog:
KubernetesOtelCollector declares `OpenTelemetryCollector` custom
resources, and this operator reconciles them into running collector
fleets — injecting a default collector image into CRs that declare
none, deriving each collector's Service ports from the receivers the
CR declares, and defaulting and validating CRs through its admission
webhooks (failurePolicy Fail).

## The Deployment Landscape

One operator serves the whole cluster: it watches every namespace and
reconciles the collector CRs (and, unmodeled in this catalog,
Instrumentation/OpAMPBridge/TargetAllocator CRs authored directly).
The catalog splits the concern in two — this kind installs the engine
once, KubernetesOtelCollector declares each collector.

### The chart/operator pinning truth

The chart's default manager image tag is its appVersion, and the
served repository index is the source of truth for what actually
installs: **chart 0.120.0 pairs with operator v0.156.0** — the newest
SERVED stable chart, verified against the repository index. Bumping
the pin requires re-staging the module's CRD files from the new chart
version — the staged files ARE the chart's CRDs at that version.

## The Templated-CRD Mechanism

Unlike plain-YAML CRD bundles, this chart TEMPLATES its CRDs: they
render inside the chart's webhook template via `tpl`/`.Files.Get`,
which means the CRD manifests carry RELEASE-DERIVED values. The
collector CRD carries the `cert-manager.io/inject-ca-from` annotation
and a version-conversion webhook `clientConfig` pointing at the
release's webhook Service (`<name>-webhook`, port 443, path
`/convert`).

This is why the modules stage RENDERED, TOKENIZED CRD files rather
than copying raw chart files: the files under `iac/crds/` were
rendered from the pinned chart with the release-derived values
replaced by `__PLANTON_RELEASE_NAME__` / `__PLANTON_NAMESPACE__`
tokens, substituted at apply time (identically on both engines) — so
the kept CRDs always point at THIS release's webhook Service and
cert-manager Certificate. The tokens appear only in the collector CRD;
the substitution is a harmless no-op on the other files.

## The CRD Lifecycle: Module-Owned, Keep-on-Uninstall

The chart templates its opentelemetry.io CRDs as release-owned
resources. Left to Helm, uninstalling the operator would delete the
CRDs — and deleting a CRD cascade-deletes every custom resource of
that type cluster-wide, taking every collector declaration with it.

The modules therefore own the CRD lifecycle end to end:

- **`crds.create` pins false unconditionally** in the rendered chart
  values — and is re-pinned AFTER the `helm_values` escape-hatch
  merge, one of the two deliberate exceptions to the escape hatch's
  last-word contract.
- **The four CRDs the default operator serves**
  (`opentelemetrycollectors`, `instrumentations`, `opampbridges`,
  `targetallocators`) **apply as module-owned resources**, keyed by
  each CRD's OWN `metadata.name` (never a positional index), so state
  addresses stay stable across file renames and reorderings.
- **Keep-on-uninstall, per engine:**
  - Terraform: `kubectl_manifest` with `apply_only = true` — the
    provider's Delete is a NO-OP (verified in the provider source), so
    destroy removes the CRDs from state without touching the cluster.
  - Pulumi: the classic yaml `ConfigGroup` with `retainOnDelete`
    delivered through a resource TRANSFORMATION — the one mechanism
    the SDK propagates to the ConfigGroup's in-process children (the
    yaml SDK forwards only parent/version options to children, and
    yaml/v2 children are created provider-side, beyond the reach of
    any SDK-side option — both verified in the pinned
    pulumi-kubernetes source).
- **Server-side apply is REQUIRED, not just preferred**: the collector
  CRD is ~418 KB — far past the 262144-byte cap on the client-side
  last-applied-configuration annotation.
- **The release depends on the CRDs**, so the operator never starts
  against an unregistered API group.
- **`skip_crds` is the bring-your-own-CRDs arm** — the CRDs are owned
  elsewhere (a GitOps-managed bundle) and the modules must not touch
  them. With the CRDs absent the operator cannot start.
- **The fifth, feature-gated `clusterobservabilities` CRD** (the
  `operator.clusterobservability` alpha gate, off by default) is
  deliberately not staged — enabling that gate via `helm_values`
  without its CRD is unsupported here.

## Webhook Certificates: Why Only cert-manager Is Modeled

The chart offers three ways to serve the admission webhooks'
certificate: cert-manager (`admissionWebhooks.certManager`, the chart
default), a Helm-generated self-signed certificate
(`admissionWebhooks.autoGenerateCert`), and operator-provided
certificate files. This component models exactly one — cert-manager —
and that is a deliberate consequence of the CRD lifecycle:

- The collector CRD carries a version-CONVERSION webhook, and the CRDs
  are retained past the release's lifetime. Their conversion trust
  (the CA bundle the API server uses to call the conversion webhook)
  must be kept current by a RUNNING reconciler — cert-manager's CA
  injector, via the CRDs' `cert-manager.io/inject-ca-from` annotation.
- A certificate embedded once at install time (the chart's
  Helm-generated arm) goes stale on rotation and silently breaks
  collector-CR version conversion long after the install succeeded —
  so no such arm exists in this spec.
- The modules re-pin `admissionWebhooks.certManager.enabled: true`
  AFTER the escape-hatch merge (the second of the two deliberate
  exceptions): disabling it would leave module-owned CRDs pointing at
  a Certificate that no longer exists.

The only typed knob is `webhook.issuer_ref` (Issuer or ClusterIssuer):
empty means the chart creates its own self-signed Issuer — the right
choice for almost everyone, since the webhook certificate only needs
to be trusted by the API server, which cert-manager's CA injection
handles.

## Image Forms

The chart takes its two images in different shapes, and the spec
models each accordingly:

- **The manager image is combined-form**: the chart's default is the
  ghcr.io image
  `ghcr.io/open-telemetry/opentelemetry-operator/opentelemetry-operator`.
  `image_registry` replaces ONLY the registry part (the path stays the
  upstream one) — the air-gap mirror seam for the one image THIS
  component's pods pull.
- **The default collector image is a repository+tag pair**:
  `manager.collectorImage` takes `repository` and `tag` separately,
  and the chart renders the `--collector-image` flag only when BOTH
  are present (a repository-only override deep-merges with the chart's
  default tag, so the flag still renders). The spec carries ONE image
  string (`default_collector_image`) and the modules split it — a tag
  exists when the last `:` comes after the last `/` (registry ports
  carry `:` too). `image_registry` does NOT rewrite this image:
  collector pods pull it, not this component.

## Design Decisions

- **The install is blocking.** The Helm release waits for the operator
  to become Available (atomic, 600s timeout, cleanup on fail): the
  manager pod mounts the cert-manager-issued webhook Secret, so an
  install without a working cert-manager — or with an unpullable image
  from a private mirror — fails THIS apply with a readiness timeout
  instead of surfacing later as collectors that mysteriously never
  reconcile.
- **The module owns namespace creation** (`create_namespace`), never
  the Helm release — pre-existing-namespace installs leave the flag
  false.
- **`fullnameOverride` pins the chart's fullname to the resource
  name** — load-bearing, not cosmetic: the staged CRDs' conversion
  webhook and `inject-ca-from` annotation point at names derived from
  it. The 30-character name budget follows: the chart's longest
  derived suffix is `-controller-manager-service-cert` (33 chars)
  against the Kubernetes 63-character name limit, and both modules
  fail loudly over the budget.
- **Chart-default-matching values render only on divergence** (the
  manager block, the issuer reference, scheduling), so the rendered
  values stay minimal on both engines — with the two re-pinned
  exceptions above.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `opentelemetry-operator` at https://open-telemetry.github.io/opentelemetry-helm-charts | Pinned 0.120.0 (spec default) |
| Operator image | `ghcr.io/open-telemetry/opentelemetry-operator/opentelemetry-operator` | Chart 0.120.0 = operator v0.156.0 |
| Default collector image | `ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-k8s` | The operator's compiled-in default; override fleet-wide via `default_collector_image` |
| CRD API group | `opentelemetry.io` | Four staged CRDs; `clusterobservabilities` (alpha gate) deliberately not staged |
| CRD files | staged with the module (tokenized renders of the chart 0.120.0 set) | Upgrade together with the chart pin |
| Webhook Service | `<name>-webhook` (port 443) | Exported as `webhook_service` |
| Webhook cert Secret | `<name>-controller-manager-service-cert` | Exported as `webhook_cert_secret_name`; the 33-char suffix behind the 30-char name budget |

## IaC Twins

Pulumi (`module/crds.go` + `module/values.go`) and Terraform
(`main.tf` + `locals.tf`) render identical chart values, the same
module-owned CRD set with the same token substitution, and the same
keep-on-uninstall posture (`retainOnDelete` transformation /
`apply_only = true`). Keep the typed-value rendering, the image-split
semantics, and the two post-merge re-pins in lockstep.
