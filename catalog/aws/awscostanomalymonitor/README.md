<p align="center">
  <img src="logo.svg" alt="AWS Cost Anomaly Monitor" width="80"/>
</p>

# AWS Cost Anomaly Monitor

Manage a [Cost Explorer anomaly monitor](https://docs.aws.amazon.com/cost-management/latest/userguide/manage-ad.html)
— the ML-driven watcher that flags unusual spend — with its folded
alert subscriptions.

## What Gets Managed

- **The monitor** (`spec.monitorName` is the display name — an
  explicit field because monitor names legally carry spaces
  `metadata.name` cannot): its shape is **DIMENSIONAL** (segment spend
  by ONE built-in dimension — SERVICE is AWS's recommended posture) or
  **CUSTOM** (watch the slice a Cost Explorer expression selects,
  supplied as the AWS Expression JSON verbatim). Both shape arms are
  create-only.
- **Subscriptions** — the folded satellites: who hears about
  anomalies, how often (IMMEDIATE individual alerts via SNS;
  DAILY/WEEKLY summaries via email — AWS pairs channel with
  frequency, spec-enforced), and above what impact threshold (a
  leveled expression over the ANOMALY_TOTAL_IMPACT dimensions).

Anomaly detection is free — AWS bills nothing for monitors or
subscriptions.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
