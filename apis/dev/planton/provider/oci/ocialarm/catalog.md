# Alarm on OCI

Deploys an Oracle Cloud Infrastructure Monitoring alarm -- a rule that evaluates metrics via Monitoring Query Language (MQL) expressions and triggers notifications to ONS topics or Streaming endpoints when thresholds are breached. Supports multi-threshold evaluation through overrides with independent query, severity, body, and pending duration per threshold. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Monitoring Alarm** -- a `monitoring.Alarm` in the specified compartment with MQL query evaluation, severity classification, and notification delivery to one or more ONS topics or Streaming streams
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the alarm

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the alarm in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A compartment containing the metrics to evaluate. This can be the same compartment or a different one (cross-compartment monitoring).
- At least one ONS notification topic or Streaming stream OCID for alarm delivery.
- The metric namespace and MQL query for the target service (e.g., `oci_computeagent` for compute metrics, `oci_blockstore` for block volume metrics).

## Deploy

### Console

Open the deployment store, find **Alarm on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CPU Utilization Critical** preset in the [Presets](#presets) tab to pre-populate a production-ready alarm for compute CPU monitoring.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciAlarm
metadata:
  name: high-cpu-alarm
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  metricCompartmentId:
    value: "ocid1.compartment.oc1..example"
  namespace: oci_computeagent
  query: "CpuUtilization[5m].mean() > 80"
  severity: critical
  destinations:
    - "ocid1.onstopic.oc1..example"
  isEnabled: true
  pendingDuration: PT5M
```

```shell
planton apply -f alarm.yaml
```

This creates an alarm that fires when mean CPU utilization exceeds 80% for 5 consecutive minutes. Notification body, message format, and overrides are not configured -- the alarm delivers raw-format notifications.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the alarm to compartments deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: monitoring
      fieldPath: status.outputs.compartmentId
  metricCompartmentId:
    valueFrom:
      kind: OciCompartment
      name: workloads
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartments first, then provisions the alarm with the resolved OCIDs.

## Key Configuration

These are the most important decisions when configuring an alarm. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**MQL query** -- The `query` field accepts any Monitoring Query Language expression. The query specifies the metric, statistic, interval, and threshold (e.g., `CpuUtilization[5m].mean() > 80`). The `namespace` field identifies the source service (e.g., `oci_computeagent`, `oci_blockstore`, `oci_autonomous_database`).

**Severity and overrides** -- Set `severity` to `critical`, `error`, `warning`, or `info` for the base rule. Use `overrides` for multi-threshold evaluation -- each override can specify its own query, severity, body, and pending duration. Overrides are evaluated in order before the base rule, enabling escalation patterns (warn at 70%, critical at 90%).

**Pending duration** -- Set `pendingDuration` in ISO 8601 format (e.g., `PT5M` for 5 minutes, `PT1H` for 1 hour). The alarm must stay in the breached state for this duration before transitioning to FIRING. Range: PT1M to PT1H. Shorter durations catch transient spikes; longer durations reduce alert noise.

**Cross-compartment monitoring** -- Set `metricCompartmentId` to a different compartment than `compartmentId` to monitor metrics from one compartment while placing the alarm in another. Set `metricCompartmentIdInSubtree: true` to include all sub-compartments of the metric compartment (requires the metric compartment to be a tenancy OCID).

**Notification format** -- Set `messageFormat` to `pretty_json` for human-readable JSON payloads or `ons_optimized` for compact email-optimized delivery. The default `raw` format works with both Notifications and Streaming destinations.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciCompartment** | `metricCompartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `alarm_id` | OCID of the monitoring alarm | Alarm management, suppression, dashboard integration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**CPU utilization critical** -- A single-threshold alarm that fires when mean CPU utilization exceeds 80% for 5 minutes, delivering pretty-JSON notifications. The standard monitoring starting point. Start from the **CPU Utilization Critical** preset.

**Multi-threshold escalation** -- An alarm with a warning base rule at 70% and a critical override at 90%. Overrides evaluate first, enabling tiered alerting with different notification bodies and pending durations per severity level. Start from the **Multi-Threshold Escalation** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment for alarm placement and metric source scoping