# AwsCostAnomalyMonitor — Terraform/OpenTofu module

Manages one Cost Explorer anomaly monitor (`aws_ce_anomaly_monitor`)
and its folded alert subscriptions (`aws_ce_anomaly_subscription`, one
per `spec.subscriptions` entry).

Module facts worth knowing before editing:

- **The shape arms are create-only**: `monitor_type`,
  `monitor_dimension`, and `monitor_specification` all force
  replacement; only the display name updates in place.
- **The CUSTOM arm's Struct renders as raw Expression JSON**
  (`jsonencode`) — the provider takes the AWS document verbatim, not
  typed blocks.
- **Subscriptions render with `for_each` keyed by
  `spec.subscriptions[].name`** — the state key and the
  `subscription_arns` output-map key. Each binds to THIS monitor's
  ARN via `monitor_arn_list`.
- **The threshold expression is LEVELED (root → leaf)** — exactly the
  one composition level AWS accepts on subscriptions, so the dynamic
  blocks unroll it 1:1 with no depth checks.
- **Frequency pairs with the channel** (spec-enforced): IMMEDIATE →
  all-SNS subscribers; DAILY/WEEKLY → all-email.

Outputs mirror the Pulumi module key-for-key: `monitor_arn`,
`subscription_arns` (map keyed by subscription name).
