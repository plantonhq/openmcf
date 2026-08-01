---
title: "Apache Flink Operator"
description: "Apache Flink Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesflinkoperator"
---

# Apache Flink Operator

The official ASF controller that turns `FlinkDeployment` declarations
into running Flink clusters, from the official chart served per
version by Apache (1.15.0 pinned; chart version = operator version =
image tag). One operator per namespace by construction — the chart
fixes its webhook artifact names — and one cluster-wide-watching
operator is the normal posture; Flink pipelines are declared with
`KubernetesFlinkDeployment`.

## Highlights

- **The webhook lifecycle, told honestly** — with the webhook on (the
  upstream default), cert-manager is a hard prerequisite: the chart
  renders Issuer/Certificate unconditionally, there is no self-signed
  fallback, and both webhook configurations fail CLOSED — bad Flink
  declarations are rejected at admission with a real message, and an
  unreachable webhook rejects every flink.apache.org admission in
  scope. `webhookEnabled: false` trades admission-time validation for
  reconcile-loop validation and drops the dependency.
- **The hardcoded password never ships** — the chart's default
  webhook keystore credential is a hardcoded public password; the
  modules generate a random per-install Secret instead and re-pin
  `useDefaultPassword: false` after the escape-hatch merge so it
  cannot resurface.
- **Removing the operator never deletes Flink declarations** — the
  four flink.apache.org CRDs ride the chart's `crds/` directory:
  installed once, never upgraded by chart bumps, kept on uninstall by
  upstream design; the `flink` job service account survives uninstall
  too, so running jobs keep their identity.
- **The fence scopes RBAC and the webhook together** — a
  `watchNamespaces` list confines the operator's reconcile RBAC AND
  the admission webhook's namespaceSelector to exactly that list.
- **Fail-loud, not fail-later** — names over the 45-character budget
  are rejected at apply time, and the atomic install waits for the
  operator (and, webhook arm, its cert-manager certificate) to become
  ready instead of surfacing as deployments that never reconcile.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
