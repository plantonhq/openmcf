# Preset: PromQL Alarm on a Prometheus Workspace

**Use case:** Alarm on metrics that live in Amazon Managed Service for
Prometheus, with Prometheus-native semantics.

In PromQL mode the query itself expresses the alarm condition — the
`> 0.8` inside the query is the threshold, `pendingPeriod` is the Prometheus
`for:` clause, and `recoveryPeriod` dampens flapping. The threshold-based
fields (`comparisonOperator`, `evaluationPeriods`, `threshold`) do not apply
and must be omitted.

## What You Get

- A CloudWatch alarm evaluating a PromQL query every 60 seconds
- Prometheus-style pending/recovery windows instead of M-of-N periods
- Standard CloudWatch actions (SNS, OpsItems) on state transitions —
  Prometheus data, CloudWatch alerting plumbing
- Outputs: `alarm_arn`, `alarm_name`

## When to Use

- Workloads instrumented with Prometheus exporters (Kubernetes node metrics,
  application histograms) whose data lands in a Prometheus workspace
- Teams that already express alerting rules in PromQL and want CloudWatch's
  action routing without re-modeling rules as metric math

## Key Configuration Choices

- **Query as contract** — the boolean-style PromQL expression carries the
  condition; adjust the threshold inside the query, not on the alarm
- **`pendingPeriod: 300`** — five continuous minutes of breach before ALARM,
  suppressing transient spikes
- **`evaluationInterval: 60`** — valid values are 10, 20, 30, or multiples
  of 60 seconds

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | Region of the Prometheus workspace and the alarm | Your deployment region |
