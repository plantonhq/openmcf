# AwsCostAnomalyMonitor — Component Guide

Authored operational judgment for the anomaly-monitor component: the
design decisions behind the spec's shape, and what to know before
operating anomaly detection in production.

## Design decisions

- **Subscriptions fold in.** A subscription's monitor list is the
  structural edge — it exists to route THIS monitor's anomalies, so
  each entry is a name-keyed satellite (the for_each key and the
  `subscription_arns` output-map key). Cross-monitor fan-in (one
  subscription watching many monitors) is the one shape the fold
  cannot express — deploy per-monitor subscriptions instead, which is
  also AWS's own console posture.
- **The monitor's display name is an explicit spec field** (spaces
  legal) and the ONLY field that updates in place — both shape arms
  force replacement.
- **The CUSTOM arm takes the AWS Expression JSON verbatim** (a
  free-form Struct) because the provider models it the same way — a
  raw JSON string, not typed blocks. The DIMENSIONAL arm is the typed,
  recommended path.
- **Frequency pairs with the channel, spec-enforced.** IMMEDIATE
  delivers individual alerts via SNS; DAILY/WEEKLY summaries deliver
  via email. AWS rejects mismatches at apply — the CEL keeps the
  failure at manifest time.

## Operating anomaly detection in production

- **A fresh monitor trains silently** on roughly ten days of history
  before it flags anything — no alerts in the first days is expected,
  not broken.
- **Impact thresholds are the noise dial**: ANOMALY_TOTAL_IMPACT_ABSOLUTE
  (dollars) and ANOMALY_TOTAL_IMPACT_PERCENTAGE (percent above
  normal), composable with AND ("$100 and 10%"). Without one, every
  flagged anomaly alerts.
- **SNS subscriptions need a topic policy allowing
  costalerts.amazonaws.com to publish** — silent alert loss
  otherwise.
- **One DIMENSIONAL/SERVICE monitor per account is AWS's recommended
  baseline**; add CUSTOM monitors for slices that deserve their own
  stream (a team's tag, a member account).
- **Import IDs are ARNs** for both the monitor and its subscriptions.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
