# Kubernetes Gatekeeper — design notes

## Grain

One resource = one Gatekeeper Helm release (official `gatekeeper`
chart, open-policy-agent.github.io index; chart version = app
version). The chart HARDCODES its resource names
(`gatekeeper-webhook-service`, `gatekeeper-webhook-server-cert`, the
webhook configuration names) — no fullname derivation exists, so the
engine is a per-cluster singleton by construction and no name budget
applies.

## The composition seam

- **In:** ConstraintTemplates / Constraints (and Config, mutators,
  ExpansionTemplates) applied as `KubernetesManifest` resources or
  GitOps after the engine installs.
- **Out:** `webhook_service_name` and `webhook_cert_secret_name`
  (chart-fixed), namespace and release name.
- **Certificates:** the `external_cert` arm composes
  `KubernetesCertificate` (secret-name reference); DNS names must
  cover the webhook Service.

## Webhook posture model

Both webhooks are typed as posture blocks: enable/disable, failure
policy, timeout — plus the delete-operations expansion on the
validating side and mutation annotations on the mutating side.
`validatingWebhookCheckIgnoreFailurePolicy` (the label-guard webhook)
is typed separately because its default (`Fail`) intentionally
differs from the policy webhook's (`Ignore`) — collapsing them into
one field would force a false choice.

## The external-cert asymmetry (chart-truth at the pin)

With `externalCertInjection.enabled`, the audit deployment's
`--disable-cert-rotation` flag is forced true by the template
(`or .Values.audit.disableCertRotation .Values.externalCertInjection.enabled`)
but the controller-manager's reads only its own value — without
`controllerManager.disableCertRotation: true` the embedded rotator
keeps overwriting the injected Secret. Both modules set it explicitly
whenever the arm is declared.

## Scheduling fan-out

The typed scheduling block applies to BOTH deployments
(controller-manager and audit): placing the engine on dedicated nodes
must move the audit loop with it, or audit findings and enforcement
observe different clusters.

## Cross-engine parity

The Terraform and Pulumi modules render byte-identical chart values.
The chart's image keys are `repository` + `release` (the tag key is
literally named "release") with a separate `crdRepository` for hook
containers — the typed image override maps repo/tag onto
repository/release; crdRepository and the curl probe image ride
`helm_values` for full air-gap installs.

## Deliberate exclusions

Violation export (exportBackend + the sidecar reader), pod-count
limits, PDB/network-policy/resource-quota toggles, PSS label sets,
per-hook images, custom webhook rules/match conditions, and
metrics-backend selection — reachable through `helm_values`, never
the primary interface. Deliberately NOT modeled at all: `emitAdmission
Events`-style involved-namespace event routing (operational nuance
without composition value).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
