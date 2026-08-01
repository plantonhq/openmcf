# Kubernetes policy engines: Kyverno and Gatekeeper forged at full depth

## What changed

- **KubernetesKyverno (880, new)** — deploys Kyverno, the
  Kubernetes-native policy engine, from the official `kyverno` Helm
  chart pinned 3.8.2 (Kyverno v1.18.2, kyverno.github.io index). The
  typed surface covers the four controllers (admission with replicas,
  sizing, scheduling and an HPA; background, cleanup and reports each
  with an enable toggle, replicas, sizing and scheduling), the policy
  CRD lifecycle (install, keep-on-uninstall — the CRDs are
  chart-templated and otherwise cascade-delete every policy on destroy;
  the post-upgrade migration hook with its reg.kyverno.io image
  documented), the engine's runtime configuration as EDITS to the
  chart's defaults (webhook namespace exclusions that re-include
  kube-system by construction, resource-filter include/exclude,
  principal exclusions, registry mutation), typed feature flags
  (fail-open override, background-scan tuning, native
  ValidatingAdmissionPolicy generation, report toggles, logging,
  omitted event types), a three-arm certificate model whose
  cert-manager arm applies to BOTH webhook servers (admission and
  cleanup), ServiceMonitor fan-out across all four controllers, a
  global image-registry override for air-gapped installs, and the Helm
  values escape hatch.

- **The webhook lifecycle is the design spine.** The chart templates no
  webhook configurations: Kyverno's admission controller registers them
  at runtime and the chart's pre-delete hook removes them at uninstall
  — the module renders that hook explicitly (default on) and the spec
  teaches what strands without it (policy rules default fail-closed, so
  an orphaned `kyverno-*` webhook configuration blocks matched
  admissions cluster-wide) along with the label-verified manual
  cleanup command.

- **KubernetesGatekeeper (881, new)** — deploys OPA Gatekeeper from the
  official `gatekeeper` Helm chart pinned 3.23.0 (chart and app move in
  lockstep, open-policy-agent.github.io index). The typed surface
  covers both webhook postures (enable, failure policy, timeout —
  with the namespace-label check webhook's fail-closed default as its
  own field, deliberately separate from the fail-open policy webhook),
  the audit controller (interval, per-constraint violation limits,
  cache-based audit, constraint-scoped kind matching, chunking,
  sizing), the namespace exemption authorization model (names and
  prefixes), engine capabilities (external data, Kubernetes CEL
  validation, generator expansion, disabled Rego builtins, deny
  logging), scheduling fanned to both deployments, all four lifecycle
  hook jobs, an external webhook-certificate arm for cert-manager, and
  the Helm values escape hatch.

- **A chart asymmetry is compensated in both engines:** enabling
  external certificate injection auto-disables the embedded cert
  rotator on the audit deployment but NOT on the controller-manager
  (template-verified at the pin) — the modules set
  `controllerManager.disableCertRotation` explicitly so the rotator
  never overwrites an injected certificate.

- **E2E surfaces with an enforcement bar.** Every scenario's verifier
  proves the ENGINE, not just the rollout: a verifier-owned policy
  (Kyverno ClusterPolicy in Enforce; Gatekeeper CEL ConstraintTemplate
  plus a deny Constraint) must REJECT a violating Pod and ADMIT a
  compliant one, with all proof artifacts swept. Behavioral lanes add
  Kyverno's mutation proof and Gatekeeper's audit-loop proof (a
  violation created before any constraint exists must be found by
  audit). Destroy assertions encode each chart's designed posture:
  Kyverno's webhook configurations must be GONE; Gatekeeper's engine
  CRDs must be KEPT.

- Presets (dev, production, air-gap for Kyverno; audit-first,
  production-enforce, cert-manager TLS for Gatekeeper), reader docs,
  catalog pages, helm-release import maps, outputs-conformance cases,
  and tier-2 CI wiring ship with both kinds. Profiles are
  `pending_proof` — live E2E follows in the proof lane.

## Why

Policy is the control plane of platform engineering: the catalog now
offers both mainstream engines — Kyverno for policy-as-Kubernetes-YAML
with mutation/generation/cleanup, Gatekeeper for the OPA constraint
framework — each modeled to full configuration depth with the failure
modes of admission control (fail-open vs fail-closed, webhook
stranding, CRD cascade) taught on the exact fields that control them.
