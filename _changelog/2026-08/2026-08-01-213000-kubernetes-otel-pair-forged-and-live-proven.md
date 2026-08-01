# Kubernetes OpenTelemetry pair forged and live-proven: the operator and the collector complete the observability tier

## What changed

- **Two new kinds built to full depth and proven against a live
  cluster, both engines, in one pass** — KubernetesOtelOperator (the
  controller that turns collector declarations into running fleets,
  from the official `opentelemetry-operator` chart 0.120.0 = operator
  v0.156.0) and KubernetesOtelCollector (one OpenTelemetryCollector
  declaration per collector: the pipeline document is the product, with
  typed deployment/daemonset/statefulset/sidecar modes, the
  operator-managed autoscaler, secret-safe credential wiring through
  `${env:...}` references, and hostPath volume plumbing for node-log
  collection). Every behavioral promise ran with verifier-output
  evidence: a real OTLP push observed in the declared pipeline's output
  on every collector lane; node log files actually ingested on the
  daemonset lane (the run-as-root pattern the new
  `pod_security_context` field exists for); a verifier-patched pipeline
  rolled onto live behavior by the operator (declared config change →
  fresh telemetry carrying the new attribute); the operator's
  fail-closed admission webhook REJECTING an invalid collector; and the
  version-CONVERSION webhook proven through cert-manager's CA-injected
  trust on every lane — a v1beta1-written collector read back through
  the v1alpha1 API. Both kinds entered the green E2E CI matrix with
  blind import round-trips proven.

- **The operator owns its CRDs; uninstalling it never deletes the
  fleet.** The chart templates its CRDs release-owned (a Helm uninstall
  would cascade-delete every collector declaration); the modules stage
  the four CRDs rendered from the pinned chart and apply them outside
  the release, retained on destroy — asserted live. Because the
  collector CRD carries a conversion webhook whose trust must outlive
  any single install, cert-manager is a REQUIRED prerequisite (its CA
  injector keeps the kept CRDs' conversion trust current); a one-shot
  self-signed certificate arm is deliberately not offered — it would
  break CR conversion silently on rotation.

- **A framework-level defect in the Terraform E2E harness is fixed for
  every module that stages sibling files.** The harness copied only a
  module directory's top-level files into its isolated workspace, so a
  module reading staged files through `../crds` silently planned ZERO
  of them — an install could "pass" riding CRDs a previous run left
  behind. The workspace copy now reproduces the module's parent
  directory tree, and the OTel modules add fail-loud staged-file count
  guards in both engines so the class can never pass silently again.

- **Retained module-owned resources now re-adopt cleanly on Pulumi.**
  A new shared provider seam (`upsertExistingObjects` + forced field
  ownership, scoped to retained resources only) lets a fresh install
  adopt CRDs an earlier install retained — previously a plain create
  failed AlreadyExists. Terraform's server-side apply always had this.

- **Sidecar-mode misconfigurations are validation errors.** The CRD
  rejects tolerations and priorityClassName in sidecar mode at
  admission; the spec now mirrors both rules so they fail at validate
  time with a message explaining that sidecar collectors live inside
  the target pods.

- **Catalog hygiene rides along**: eight secret-coverage findings
  (pointless exemptions and one missing annotation) fixed across
  Keycloak, Harbor, OpenBao and OpenFGA specs; a dangling
  Gatekeeper→Certificate foreign-key field path corrected; the
  containment ledger regenerated with three missing namespace verdicts.
