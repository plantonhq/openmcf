---
title: "PostgreSQL Database System"
description: "PostgreSQL Database System deployment documentation"
icon: "package"
order: 100
componentName: "ocipostgresqldbsystem"
---

# PostgreSQL Database System on OCI

Deploys an Oracle Cloud Infrastructure PostgreSQL Database System -- a fully managed PostgreSQL service running on dedicated compute shapes with configurable storage durability, flexible OCPU/memory sizing, read replicas, and built-in backup policies. The DB System supports regional storage durability (replicated across availability domains) or AD-local storage, with read scaling via configurable instance counts and an optional reader endpoint. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **PostgreSQL DB System** -- a managed PostgreSQL database in the specified compartment and subnet with the chosen shape, storage configuration, and credential setup
- **Read Replicas** -- created when `instanceCount` is 2 or higher; additional instances receive replicated data for read scaling
- **Reader Endpoint** -- created only when `networkDetails.isReaderEndpointEnabled` is `true`; a separate DNS endpoint distributes read queries across replica instances
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the DB System

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the DB System in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet in an existing VCN. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. The subnet determines the network placement of all DB System instances.
- For AD-local storage: an availability domain name (e.g., `Uocm:PHX-AD-1`). Not needed when `storageDetails.isRegionallyDurable` is `true`.
- Database credentials: either a plain-text password or an OCI Vault secret OCID for production environments.

## Deploy

### Console

Open the deployment store, find **PostgreSQL Database System on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Regionally Durable** preset in the [Presets](#presets) tab to pre-populate a production PostgreSQL DB System with regional storage durability and a read replica.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciPostgresqlDbSystem
metadata:
  name: app-postgres
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  dbVersion: "16"
  shape: VM.Standard.E4.Flex
  instanceOcpuCount: 4
  instanceMemorySizeInGbs: 32
  instanceCount: 1
  networkDetails:
    subnetId:
      value: "ocid1.subnet.oc1.phx..example"
  storageDetails:
    isRegionallyDurable: true
  credentials:
    username: postgres
    passwordDetails:
      passwordType: plain_text
      password: "Ex4mpl3!Passw0rd"
```

```shell
planton apply -f postgresql-db.yaml
```

This creates a single-instance PostgreSQL 16 DB System with 4 OCPUs, 32 GB memory, and regionally durable storage. No read replicas or reader endpoint are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DB System to a compartment, subnet, and security groups deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: app-compartment
      fieldPath: status.outputs.compartmentId
  networkDetails:
    subnetId:
      valueFrom:
        kind: OciSubnet
        name: db-subnet
        fieldPath: status.outputs.subnetId
    nsgIds:
      - valueFrom:
          kind: OciSecurityGroup
          name: postgres-nsg
          fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the DB System with the resolved values.

## Key Configuration

These are the most important decisions when configuring a PostgreSQL DB System. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Storage durability** -- Set `storageDetails.isRegionallyDurable` to `true` for production workloads -- data is replicated across availability domains so an entire AD failure does not cause data loss. Set to `false` for development environments where AD-local storage is cheaper; requires specifying `storageDetails.availabilityDomain`. This is a ForceNew field -- changing it forces recreation.

**Instance count and read scaling** -- Set `instanceCount` to 1 for standalone, 2+ for primary plus read replicas. Enable `networkDetails.isReaderEndpointEnabled` to create a separate DNS endpoint that load-balances read queries across replicas. Each instance gets the OCPU and memory allocation specified by `instanceOcpuCount` and `instanceMemorySizeInGbs`.

**Flex shape sizing** -- The `shape` field (e.g., `VM.Standard.E4.Flex`) determines the compute architecture. `instanceOcpuCount` and `instanceMemorySizeInGbs` independently set the CPU and memory for each instance. Both are updatable without recreation.

**Backup schedule** -- Configure `managementPolicy.backupPolicy` with `kind` set to `daily`, `weekly`, `monthly`, or `none`. Daily requires `backupStart` (UTC hour). Weekly adds `daysOfTheWeek`. Monthly adds `daysOfTheMonth`. Set `retentionDays` for how long backups are retained after DB System deletion.

**Credentials** -- The `credentials` block is immutable after creation. Set `passwordDetails.passwordType` to `plain_text` for development or `vault_secret` for production (references an OCI Vault secret by OCID, keeping passwords out of the manifest).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `networkDetails.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `networkDetails.nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `db_system_id` | OCID of the PostgreSQL DB System | Monitoring, IAM policy scoping, resource management |
| `primary_db_endpoint_private_ip` | Private IP address of the primary (read-write) endpoint | Application connection strings, DNS records |
| `admin_username` | Administrator username (computed after creation) | Application connection configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Regionally durable** -- A production DB System with regional storage durability, a read replica, reader endpoint, daily backups with 30-day retention, and Vault-managed credentials. Start from the **Regionally Durable** preset.

**Standalone development** -- A single-instance DB System with AD-local storage, minimal compute, short backup retention, and plain-text credentials for development and testing. Start from the **Standalone Development** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this DB System
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet for the DB System network placement
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for the DB System instances