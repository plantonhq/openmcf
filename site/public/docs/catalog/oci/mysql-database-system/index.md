---
title: "MySQL Database System"
description: "MySQL Database System deployment documentation"
icon: "package"
order: 100
componentName: "ocimysqldbsystem"
---

# MySQL Database System on OCI

Deploys an Oracle Cloud Infrastructure MySQL HeatWave Database System -- a fully managed MySQL database service with optional High Availability across fault domains, automated backups with point-in-time recovery, read-scaling endpoints, and integrated HeatWave in-memory analytics acceleration. The component manages the DB System resource itself; HeatWave cluster and replication channels are separate OCI resources with independent lifecycles. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **MySQL DB System** -- an `oci_mysql_mysql_db_system` in the specified compartment and subnet, placed in a given availability domain on a chosen compute shape, with a primary read/write endpoint
- **High Availability replicas** -- created only when `isHighlyAvailable` is `true`; three instances are provisioned across different fault domains with automatic failover
- **Automatic backups** -- created only when `backupPolicy` is configured; daily backups with configurable retention and optional point-in-time recovery
- **Read Endpoint** -- created only when `readEndpoint.isEnabled` is `true`; a separate DNS endpoint distributes read queries across HA replicas
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the DB System

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the DB System in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet in an existing VCN. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. Changing the subnet forces recreation.
- An availability domain name (e.g., `Uocm:PHX-AD-1`). Changing the availability domain forces recreation.
- A compute shape name (e.g., `MySQL.VM.Standard.E4.1.8GB`, `MySQL.VM.Standard.E4.4.64GB`).

## Deploy

### Console

Open the deployment store, find **MySQL Database System on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **High Availability** preset in the [Presets](#presets) tab to pre-populate a production MySQL DB System with HA, PITR-enabled backups, and delete protection.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciMysqlDbSystem
metadata:
  name: app-mysql
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  availabilityDomain: "Uocm:PHX-AD-1"
  shapeName: "MySQL.VM.Standard.E4.1.8GB"
  subnetId:
    value: "ocid1.subnet.oc1.phx..example"
  adminUsername: admin
  adminPassword: "Ex4mpl3!Passw0rd"
```

```shell
planton apply -f mysql-db.yaml
```

This creates a single-instance MySQL DB System with the smallest compute shape and Oracle-managed encryption. No HA, backups, or delete protection are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DB System to a compartment, subnet, and security groups deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: app-compartment
      fieldPath: status.outputs.compartmentId
  subnetId:
    valueFrom:
      kind: OciSubnet
      name: db-subnet
      fieldPath: status.outputs.subnetId
  nsgIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: mysql-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the MySQL DB System with the resolved values.

## Key Configuration

These are the most important decisions when configuring a MySQL DB System. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**High Availability** -- Set `isHighlyAvailable` to `true` to provision three instances across fault domains with automatic failover. Standby instances are not directly accessible. HA roughly triples the compute cost but provides production-grade resilience. Enable `readEndpoint` alongside HA to distribute read queries across replicas.

**Compute shape and placement** -- The `shapeName` determines CPU, memory, and network bandwidth (e.g., `MySQL.VM.Standard.E4.1.8GB` for 1 OCPU/8 GB). The `availabilityDomain` places the primary instance. Both are ForceNew fields -- changing either forces recreation.

**Storage and auto-expansion** -- Configure `dataStorage.dataStorageSizeInGb` for initial volume size. Enable `dataStorage.isAutoExpandStorageEnabled` with `maxStorageSizeInGbs` to auto-expand storage when usage nears the limit, eliminating storage-related outages.

**Backups and PITR** -- Configure `backupPolicy` with `isEnabled: true`, a `retentionInDays` value, and optionally a `windowStartTime`. Enable `backupPolicy.pitrPolicy.isEnabled` for point-in-time recovery, allowing restoration to any second within the retention window.

**Deletion policy** -- The `deletionPolicy` message provides three independent safety controls: `isDeleteProtected` (prevents deletion entirely), `finalBackup` (`REQUIRE_FINAL_BACKUP` or `SKIP_FINAL_BACKUP`), and `automaticBackupRetention` (`RETAIN` or `DELETE`). For production, enable all three safeguards.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `db_system_id` | OCID of the MySQL DB System | Monitoring, IAM policy scoping, resource management |
| `endpoint_hostname` | Hostname of the primary (read/write) endpoint | Application connection strings |
| `endpoint_ip_address` | Private IP address of the primary (read/write) endpoint | DNS records, firewall allowlists |
| `endpoint_port` | TCP port of the primary (read/write) endpoint | Application connection configuration (default 3306) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**High availability** -- A production DB System with HA across fault domains, PITR-enabled backups with 30-day retention, auto-expanding storage, delete protection, and crash recovery enabled. Start from the **High Availability** preset.

**Standalone development** -- A single-instance DB System with the smallest shape, minimal storage, short backup retention, and no HA or delete protection for cost-optimized development. Start from the **Standalone Development** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this DB System
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet where the DB System is placed
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for the DB System VNIC