# NoSQL Table on OCI

Deploys an Oracle Cloud Infrastructure NoSQL table with DDL-defined schema, configurable throughput capacity (provisioned or on-demand), and optional secondary indexes. The table schema is defined entirely through a DDL statement -- OCI NoSQL's native schema mechanism -- supporting columns of type STRING, INTEGER, JSON, TIMESTAMP, and more. Indexes are bundled as immutable sub-resources. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **NoSQL Table** -- an `oci_nosql_table` in the specified compartment with the schema defined by the DDL statement and throughput limits set by the capacity mode
- **Secondary Indexes** -- one `oci_nosql_index` per entry in the `indexes` list. Each index is immutable; any change to an existing index forces its recreation. Indexes support plain columns and JSON field paths within JSON-typed columns.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the table

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the table in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A valid DDL statement (`CREATE TABLE` for new tables). The table name in the DDL must match the `name` field exactly.

## Deploy

### Console

Open the deployment store, find **NoSQL Table on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Provisioned Throughput** preset in the [Presets](#presets) tab to pre-populate a table with explicit read/write capacity limits.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciNosqlTable
metadata:
  name: sessions
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  name: sessions
  ddlStatement: >-
    CREATE TABLE sessions (
      session_id STRING,
      user_id STRING,
      expires_at TIMESTAMP(3),
      PRIMARY KEY(session_id)
    )
  tableLimits:
    maxReadUnits: 100
    maxWriteUnits: 100
    maxStorageInGbs: 25
```

```shell
planton apply -f nosql-table.yaml
```

This creates a NoSQL table with provisioned throughput (100 read units, 100 write units) and 25 GB of storage. No secondary indexes are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the table to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: app-compartment
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the NoSQL table with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a NoSQL table. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Capacity mode** -- Set `tableLimits.capacityMode` to `provisioned` for predictable workloads with well-understood access patterns (cheaper at sustained throughput). Set to `on_demand` for bursty or unpredictable traffic (scales automatically, pay-per-request). When provisioned, `maxReadUnits` and `maxWriteUnits` are required and set explicit throughput ceilings.

**DDL statement** -- The `ddlStatement` defines the table schema using OCI NoSQL's DDL syntax. Use `CREATE TABLE` for new tables and `ALTER TABLE` for schema evolution. Column order should not change; new columns can only be appended. The table name in the DDL must match the `name` field.

**Secondary indexes** -- Each entry in `indexes` creates an immutable `oci_nosql_index`. Changes to an existing index force its recreation. Indexes support plain columns and JSON field paths -- for JSON columns, specify `jsonFieldType` and `jsonPath` to index specific nested fields.

**Auto-reclaimable** -- Set `isAutoReclaimable: true` for temporary tables that OCI may reclaim after extended inactivity. This is a ForceNew field -- changing it forces recreation. Omit for persistent tables.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `table_id` | OCID of the NoSQL table | Application configuration, monitoring, IAM policy scoping |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Provisioned throughput** -- A table with explicit read/write capacity limits and a timestamp-indexed schema. Predictable costs for workloads with well-understood access patterns. Start from the **Provisioned Throughput** preset.

**On-demand** -- A table with automatic throughput scaling and pay-per-request pricing. Ideal for bursty traffic or new workloads where capacity cannot be estimated. Start from the **On-Demand** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this NoSQL table