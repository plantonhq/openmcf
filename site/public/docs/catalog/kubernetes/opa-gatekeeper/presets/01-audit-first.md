---
title: "Audit-first preset"
description: "Gatekeeper as it ships: three webhook replicas, the policy webhook fail-OPEN (`failurePolicy: Ignore`), and the audit controller re-checking existing resources every 60 seconds. This is the right..."
type: "preset"
rank: "01"
presetSlug: "01-audit-first"
componentSlug: "opa-gatekeeper"
componentTitle: "OPA Gatekeeper"
provider: "kubernetes"
icon: "package"
order: 1
---

# Audit-first preset

Gatekeeper as it ships: three webhook replicas, the policy webhook
fail-OPEN (`failurePolicy: Ignore`), and the audit controller
re-checking existing resources every 60 seconds. This is the right
first posture for adopting policy on a running cluster: apply your
ConstraintTemplates and Constraints with `enforcementAction: warn` (or
`dryrun`), read violations out of each constraint's status and the
audit results, and only then decide what deserves to block admissions.

Fail-open is the deliberate trade here. An engine outage never stops
the cluster admitting resources — but a request that slips through
during one is not evaluated. The one fail-CLOSED piece is the
namespace-label check webhook, which guards Gatekeeper's own exemption
label (`admission.gatekeeper.sh/ignore`) so nobody can exempt a
namespace while the engine is down; its blast radius is namespace
label edits only.

Know the CRD posture: the engine CRDs install once from the chart's
crds/ directory and SURVIVE uninstall, as do the per-template
constraint CRDs Gatekeeper creates at runtime — destroying the engine
does not delete your policy library.

Change first: nothing, until audit results are boring. Then move to
the production-enforce preset.

See [01-audit-first.yaml](./01-audit-first.yaml) for the manifest.
