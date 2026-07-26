---
title: "OPA Gatekeeper"
description: "OPA Gatekeeper deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesgatekeeper"
---

# OPA Gatekeeper

Constraint-based policy for Kubernetes, from the OPA project. Policies
are ConstraintTemplates (Rego or Kubernetes CEL) instantiated as
typed Constraint resources, enforced at admission and continuously
audited against what already runs — the framework of choice where
policy teams standardize on OPA across more than just Kubernetes.

## Highlights

- **Adopt without breaking anything** — fail-open webhook and the
  audit loop by default: see every violation on a brownfield cluster
  before anything blocks; one typed field flips to fail-closed
  enforcement when audit results are boring.
- **The exemption model, typed** — namespaces and prefixes authorized
  to opt out, with the label-guard webhook (fail-closed, narrow blast
  radius) protecting the exemption mechanism itself.
- **Audit tuned for scale** — interval, constraint-scoped kind
  matching, chunking, and per-constraint violation limits as typed
  fields, not flags buried in a Deployment.
- **Policy library survives the engine** — engine and
  constraint-template CRDs are kept on uninstall by design; destroy
  removes the webhooks and workloads, never your constraints.
- **cert-manager option done right** — external webhook certificates
  with the embedded rotator disabled on BOTH deployments (the chart
  only handles one — the module compensates).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
