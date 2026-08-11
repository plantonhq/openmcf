# GCP Monitoring Alert Policy

Creates a Cloud Monitoring alerting policy — the rule that watches metrics or logs (threshold, absence, log-match, MQL, PromQL, or scheduled-SQL conditions) and opens an incident, notifying the referenced channels, when its conditions are met. Channels wire in through ValueFromRef, so the whole paging path — check, policy, channel — is one declarative graph.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Alert Policy** -- a `monitoring.AlertPolicy` with the configured conditions, combiner, severity, notification channels, alert strategy, and runbook documentation
- **Monitoring API enablement** -- `monitoring.googleapis.com` enabled in the target project (never disabled on destroy)
- **GCP Labels** -- resource metadata labels applied automatically as `user_labels`

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A GCP project** where the policy is created (directly or via a GcpProject reference).
- **IAM**: the deploying identity needs `roles/monitoring.alertPolicyEditor` or broader.
- **Notification channels** to page — deploy GcpMonitoringNotificationChannel resources first (or in the same chart).

## Deploy

### Console

Open the deployment store, find **GCP Monitoring Alert Policy**, and click **Deploy**. Start from the **CPU Threshold** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpMonitoringAlertPolicy
metadata:
  name: cpu-saturation
  org: acme-corp
  env: prod
spec:
  combiner: OR
  severity: WARNING
  conditions:
    - displayName: cpu above 80%
      conditionThreshold:
        filter: metric.type="compute.googleapis.com/instance/cpu/utilization" AND resource.type="gce_instance"
        comparison: COMPARISON_GT
        thresholdValue: 0.8
        duration: 300s
        aggregations:
          - alignmentPeriod: 60s
            perSeriesAligner: ALIGN_MEAN
```

```shell
planton apply -f alert-policy.yaml
```

This opens a WARNING incident when any instance's mean CPU stays above 80% for five sustained minutes.

### InfraChart

Wire the paging path in one chart:

```yaml
spec:
  notificationChannels:
    - valueFrom:
        kind: GcpMonitoringNotificationChannel
        name: oncall-email
        fieldPath: status.outputs.channel_name
```

The InfraPipeline deploys the channel first, then the policy with the resolved channel name.

## Key Configuration

**Conditions** -- one to six, each with exactly one condition arm. `conditionThreshold` is the workhorse (metric crosses a value for a duration); `conditionAbsent` alerts on silence; `conditionMatchedLog` alerts on log entries (and requires the rate limit); `conditionPrometheusQueryLanguage` carries ported Prometheus rules; `conditionSql` runs scheduled SQL over log analytics.

**Combiner** -- how multiple conditions merge into one incident decision (`OR` is the default posture; required by GCP even with one condition).

**Alert strategy** -- `autoClose` for incidents whose condition stopped violating, `notificationRateLimit` for log-based policies, `notificationChannelStrategy` for re-notification cadence per channel subset.

**Documentation** -- the runbook the on-call engineer sees, with `${resource.label.*}` substitution and up to three links.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpMonitoringNotificationChannel** | `notificationChannels[]` | `status.outputs.channel_name` |
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `policy_name` | `projects/{project}/alertPolicies/{id}` | Snooze configs, dashboards, Monitoring API cross-references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**CPU threshold** -- the canonical infrastructure alert: sustained CPU above a bound. Start from the **CPU Threshold** preset.

**Uptime-check failure** -- pages when a GcpMonitoringUptimeCheck stops passing — the other half of every availability monitor. Start from the **Uptime Check Failure** preset.

**Error-log match** -- pages on matching log entries (with the required rate limit), extracting labels into the incident. Start from the **Error Log Match** preset.

## Works With

- [**GCP Monitoring Notification Channel**](/cloud-catalog/gcp-monitoring-notification-channel) -- where incidents are delivered
- [**GCP Monitoring Uptime Check**](/cloud-catalog/gcp-monitoring-uptime-check) -- the availability probe policies typically watch
- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the policy is created
