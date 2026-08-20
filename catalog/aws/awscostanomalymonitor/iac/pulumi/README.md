# AwsCostAnomalyMonitor — Pulumi module (Go)

Manages one Cost Explorer anomaly monitor
(`costexplorer.AnomalyMonitor`) and its folded alert subscriptions
(`costexplorer.AnomalySubscription`, one per `spec.subscriptions`
entry).

Module facts worth knowing before editing:

- **The shape arms are create-only**: `MonitorType`,
  `MonitorDimension`, and `MonitorSpecification` all force
  replacement; only the display name updates in place.
- **The CUSTOM arm's Struct renders as raw Expression JSON** — the
  provider takes the AWS document verbatim, not typed blocks.
- **Subscriptions key by `spec.subscriptions[].name`** — the logical
  name and the `subscription_arns` output-map key. Each binds to THIS
  monitor's ARN via `MonitorArnLists`.
- **The threshold expression is LEVELED (root → leaf)** — exactly the
  one composition level AWS accepts on subscriptions, so the builder
  is a 1:1 typed walk with no depth checks.
- **Frequency pairs with the channel** (spec-enforced): IMMEDIATE →
  all-SNS subscribers; DAILY/WEEKLY → all-email.

Outputs mirror the Terraform module key-for-key: `monitor_arn`,
`subscription_arns` (map keyed by subscription name).
