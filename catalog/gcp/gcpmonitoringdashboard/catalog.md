# GCP Monitoring Dashboard

Creates a Cloud Monitoring dashboard — the console page of charts, scorecards, tables, and text widgets a team watches to understand a service's health at a glance. The dashboard body is one JSON document in the Monitoring API's own format: build visually in the GCP console, export the JSON, and manage it declaratively from then on.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Dashboard** -- a `monitoring.Dashboard` from the spec's `dashboardJson` document
- **Monitoring API enablement** -- `monitoring.googleapis.com` enabled in the target project (never disabled on destroy)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A GCP project** where the dashboard is created (directly or via a GcpProject reference).
- **IAM**: the deploying identity needs `roles/monitoring.dashboardEditor` or broader.

## Deploy

### Console

Open the deployment store, find **GCP Monitoring Dashboard**, and click **Deploy**. Start from the **Golden Signals** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpMonitoringDashboard
metadata:
  name: api-health
  org: acme-corp
  env: prod
spec:
  dashboardJson: |
    {
      "displayName": "API health",
      "gridLayout": {
        "columns": "2",
        "widgets": [
          {
            "title": "CPU utilization",
            "xyChart": {
              "dataSets": [{
                "timeSeriesQuery": {
                  "timeSeriesFilter": {
                    "filter": "metric.type=\"compute.googleapis.com/instance/cpu/utilization\" resource.type=\"gce_instance\"",
                    "aggregation": {"perSeriesAligner": "ALIGN_MEAN"}
                  }
                }
              }]
            }
          }
        ]
      }
    }
```

```shell
planton apply -f dashboard.yaml
```

## Key Configuration

**dashboardJson** -- the whole dashboard as one JSON document (the Monitoring API's [Dashboard format](https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards)): `displayName` plus exactly one layout (`gridLayout`, `mosaicLayout`, `rowLayout`, or `columnLayout`) whose widgets carry the charts. The practical workflow: build the dashboard visually in the GCP console, open its **JSON editor**, and paste the export here — server-assigned keys round-trip cleanly.

**deletionPolicy** -- `DELETE` (default), `PREVENT` (protects a team's primary operational view from accidental teardown), or `ABANDON` (remove from management, keep in GCP).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `dashboard_name` | `projects/{project}/dashboards/{id}` | Alert-documentation links, Monitoring API cross-references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Golden signals** -- the SRE starter dashboard: traffic, errors, latency, saturation on one page. Start from the **Golden Signals** preset.

**Infrastructure overview** -- fleet-level CPU, memory, disk, and network for Compute Engine workloads. Start from the **Infrastructure Overview** preset.

## Works With

- [**GCP Monitoring Alert Policy**](/cloud-catalog/gcp-monitoring-alert-policy) -- link the dashboard from alert runbook documentation
- [**GCP Log Metric**](/cloud-catalog/gcp-log-metric) -- chart log-based metrics on dashboard widgets
- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the dashboard is created
