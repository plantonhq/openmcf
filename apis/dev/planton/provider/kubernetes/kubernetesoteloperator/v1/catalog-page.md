# OpenTelemetry Operator

The controller that turns `OpenTelemetryCollector` declarations into
running collector fleets, from the official `opentelemetry-operator`
Helm chart (0.120.0 = operator v0.156.0). One operator per cluster: it
watches every namespace and reconciles the collectors declared with
`KubernetesOtelCollector`, injecting a default collector image and
deriving Service ports from each CR's declared receivers.

## Highlights

- **Removing the operator never deletes the fleet** — the module (not
  Helm) owns the four opentelemetry.io CRDs, applied outside the
  release and retained on destroy; a Helm-owned install would
  cascade-delete every collector declaration on uninstall.
- **cert-manager keeps conversion trust alive** — the collector CRD
  carries a version-conversion webhook, and the retained CRDs' trust
  rides cert-manager's CA injector; a one-shot install-time
  certificate would go stale on rotation and silently break CR
  conversion, which is why no self-signed arm exists.
- **The escape hatch has two guarded keys** — `helm_values` merges
  last with Helm `-f` semantics, but `crds.create: false` and
  `admissionWebhooks.certManager.enabled: true` are re-pinned after
  the merge; both are load-bearing design, not defaults.
- **Fleet-wide collector image, one field** — `default_collector_image`
  sets what the operator injects into CRs that declare none (the
  air-gap path for collector pods); `image_registry` mirrors the
  operator's own image separately.
- **Fail-loud, not fail-later** — names over the 30-character budget
  are rejected at apply time, and the atomic install waits for the
  operator (and its cert-manager-issued webhook Secret) to become
  ready instead of surfacing as collectors that never reconcile.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
