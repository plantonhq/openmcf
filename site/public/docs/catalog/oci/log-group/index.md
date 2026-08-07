---
title: "Log Group"
description: "Log Group deployment documentation"
icon: "package"
order: 100
componentName: "ociloggroup"
---

# Log Group on OCI

Deploys an Oracle Cloud Infrastructure Log Group with bundled individual logs -- the organizational container for the OCI Logging service. A log group holds service logs (auto-collected from OCI services like VCN flow logs, Object Storage, and API Gateway) and custom logs (ingested via the Logging Ingestion API). Logs are bundled as sub-resources because they cannot exist outside a log group. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Log Group** -- a `logging.LogGroup` in the specified compartment with display name and description
- **Logs** -- one `logging.Log` per entry in the `logs` list. Each log is either a service log (with source configuration) or a custom log (accepting entries via the Ingestion API). Logs depend on the log group and are created after it.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the log group and each log

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the log group in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For service logs: the OCID of the source resource (e.g., a VCN subnet for flow logs, a bucket for Object Storage logs, an API gateway deployment). Provide directly or reference via ValueFromRef.
- Knowledge of the target service's log namespace and category (e.g., `flowlogs`/`all`, `objectstorage`/`write`, `apigateway`/`access`).

## Deploy

### Console

Open the deployment store, find **Log Group on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **VCN Flow Logs** preset in the [Presets](#presets) tab to pre-populate a log group with a service log collecting VCN network traffic data.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciLogGroup
metadata:
  name: network-logs
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  description: VCN flow logs for network traffic auditing
  logs:
    - displayName: vcn-flow-log
      logType: service
      isEnabled: true
      retentionDuration: 90
      configuration:
        service: flowlogs
        resource:
          value: "ocid1.subnet.oc1..example"
        category: all
```

```shell
planton apply -f log-group.yaml
```

This creates a log group with a single service log collecting all VCN flow log categories from the specified subnet, retained for 90 days. Custom logs and additional service logs can be added to the `logs` list.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the log group to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: observability
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the log group with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a log group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Log type** -- Each log entry must set `logType` to `service` or `custom`. Service logs auto-collect from OCI services and require a `configuration` block specifying the source service, resource OCID, and category. Custom logs accept entries pushed via the Logging Ingestion API and do not require configuration. The `logType` is a ForceNew field -- changing it forces log recreation.

**Service log configuration** -- For service logs, `configuration.service` identifies the OCI service (e.g., `flowlogs`, `objectstorage`, `apigateway`, `loadbalancer`, `functionsInvoke`), `configuration.resource` provides the source resource OCID, and `configuration.category` selects the log category. The entire configuration block is ForceNew.

**Retention duration** -- Set `retentionDuration` in 30-day increments between 30 and 180 days. When omitted, OCI defaults to 30 days. Longer retention increases storage costs but is required for compliance (e.g., 90 days for SOC 2, 180 days for extended audit trails). Retention is updatable without recreation.

**Display name uniqueness** -- Each log's `displayName` must be unique within the log group. The display name is used as the resource key in IaC modules. Choose descriptive names that identify the log source (e.g., `vcn-flow-log`, `api-gateway-access`, `app-custom-log`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `log_group_id` | OCID of the log group | Log analytics, monitoring dashboards, IAM policy scoping |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**VCN flow logs** -- A log group with a service log collecting all VCN flow log categories from a subnet, retained for 90 days. The standard configuration for network traffic auditing and troubleshooting. Start from the **VCN Flow Logs** preset.

**Custom application logs** -- A log group with a custom log for application-level entries pushed via the Logging Ingestion API. Retained for 30 days. Used for application observability without infrastructure-level log collection. Start from the **Custom Application Logs** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this log group