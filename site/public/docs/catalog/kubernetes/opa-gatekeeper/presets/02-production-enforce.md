---
title: "Production enforce preset"
description: "Gatekeeper with teeth. The policy webhook goes fail-CLOSED (`failurePolicy: Fail`): a request the engine cannot evaluate is REJECTED, closing the window an attacker could time against an engine..."
type: "preset"
rank: "02"
presetSlug: "02-production-enforce"
componentSlug: "opa-gatekeeper"
componentTitle: "OPA Gatekeeper"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production enforce preset

Gatekeeper with teeth. The policy webhook goes fail-CLOSED
(`failurePolicy: Fail`): a request the engine cannot evaluate is
REJECTED, closing the window an attacker could time against an engine
outage. That posture is only responsible with the webhook highly
available — three replicas here, spread across nodes by the chart's
default anti-affinity — and with a timeout (5s) short enough that a
sick engine degrades admissions instead of hanging them.

The audit loop is tuned for a real cluster: `matchKindOnly` stops the
controller listing every resource kind in existence (only kinds some
constraint actually matches), 120s spacing halves steady-state API
load, and 50 recorded violations per constraint gives dashboards
enough to page on without bloating etcd objects. `logDenies` writes
every blocked admission to the controller log — your audit trail of
what enforcement actually did.

`exemptNamespacePrefixes: [kube-]` AUTHORIZES kube-* namespaces to
carry the exemption label; the chart's post-install hook labels only
Gatekeeper's own namespace. Exempting the control plane is the safe
default — a constraint that blocks kube-system pods can take the
cluster down with it.

Change first: run the audit-first preset until violations are
understood — flipping straight to Fail on a brownfield cluster is how
platform teams page themselves at 3am.

See [02-production-enforce.yaml](./02-production-enforce.yaml) for
the manifest.
