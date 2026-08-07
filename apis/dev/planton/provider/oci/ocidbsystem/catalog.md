# Database System on OCI

Deploys an Oracle Cloud Infrastructure Database System -- a managed Oracle Database running on Virtual Machine or Bare Metal infrastructure. A DB System provisions the underlying compute and storage, a DB Home containing the Oracle Database software, and an initial database instance. These three layers form an inseparable unit at creation time. The component supports single-node deployments and two-node Real Application Clusters (RAC) with fault domain separation, rolling maintenance for zero-downtime patching, and dedicated backup subnets. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DB System** -- the compute and storage infrastructure in the specified compartment, subnet, and availability domain with the chosen shape, CPU cores, and storage allocation
- **DB Home** -- an Oracle Database software installation at the specified version (e.g., 19c, 21c, 23ai) or from a custom software image
- **Initial Database** -- an Oracle Database instance within the DB Home with the specified database name, administrator password, character set, and optional pluggable database
- **Automatic Backups** -- created only when `dbHome.database.dbBackupConfig.autoBackupEnabled` is `true`; daily backups with configurable retention window
- **RAC Nodes** -- created only when `nodeCount` is 2; a second node in a different fault domain with shared ASM storage and automatic failover
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the DB System

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the DB System in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet for the client network. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. Changing the subnet forces recreation.
- An availability domain name (e.g., `Uocm:PHX-AD-1`). Changing the availability domain forces recreation.
- At least one SSH public key in OpenSSH format for administrative node access. SSH keys are required -- this is a VM/BM deployment, not a serverless managed service.
- For RAC deployments: a backup subnet for inter-node communication (cache fusion). This must be a separate subnet from the client network.
- An administrator password for the SYS and SYSTEM users (2-30 characters, must include uppercase, lowercase, and numeric).

## Deploy

### Console

Open the deployment store, find **Database System on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Single-Node VM** preset in the [Presets](#presets) tab to pre-populate a single-node Oracle Database with Standard Edition and automatic backups.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciDbSystem
metadata:
  name: oracle-db
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  availabilityDomain: "Uocm:PHX-AD-1"
  shape: VM.Standard.E4.Flex
  cpuCoreCount: 2
  subnetId:
    value: "ocid1.subnet.oc1.phx..example"
  sshPublicKeys:
    - "ssh-rsa AAAA..."
  hostname: oracledb
  databaseEdition: standard_edition
  licenseModel: license_included
  dataStorageSizeInGb: 256
  nodeCount: 1
  dbHome:
    dbVersion: "19.0.0.0"
    displayName: dbhome1
    database:
      adminPassword: "Ex4mpl3!Passw0rd"
      dbName: mydb
      pdbName: mypdb
```

```shell
planton apply -f oracle-db.yaml
```

This creates a single-node VM DB System with Oracle 19c Standard Edition, 2 OCPUs, 256 GB storage, and a pluggable database. No automatic backups or RAC clustering are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DB System to a compartment, subnets, and security groups deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: db-compartment
      fieldPath: status.outputs.compartmentId
  subnetId:
    valueFrom:
      kind: OciSubnet
      name: client-subnet
      fieldPath: status.outputs.subnetId
  nsgIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: oracle-db-nsg
        fieldPath: status.outputs.networkSecurityGroupId
  backupSubnetId:
    valueFrom:
      kind: OciSubnet
      name: backup-subnet
      fieldPath: status.outputs.subnetId
  backupNetworkNsgIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: backup-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnets, and security groups first, then provisions the DB System with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Database System. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Database edition** -- Set `databaseEdition` to `standard_edition` for most workloads. Use `enterprise_edition` for partitioning, advanced compression, and Data Guard. `enterprise_edition_high_performance` adds in-memory database. `enterprise_edition_extreme_performance` is required for RAC (2-node) deployments and includes all features. The edition is immutable after creation.

**Node count and RAC** -- Set `nodeCount` to 1 for a single-node system or 2 for a two-node RAC cluster. RAC provides active-active clustering with automatic failover and zero-downtime rolling maintenance. RAC requires `enterprise_edition_extreme_performance`, a dedicated `backupSubnetId` for inter-node cache fusion, and `faultDomains` to distribute nodes across fault domains.

**Shape and CPU cores** -- The `shape` field (e.g., `VM.Standard.E4.Flex`, `BM.DenseIO2.52`) determines the compute architecture and core count range. `cpuCoreCount` sets the number of OCPUs (VM) or total cores (BM). For flex shapes, CPU and memory scale independently.

**DB Home and database version** -- The `dbHome.dbVersion` (e.g., `"19.0.0.0"`) sets the Oracle Database software version. Alternatively, use `dbHome.databaseSoftwareImageId` for a custom software image. These are mutually exclusive. The initial database is configured within `dbHome.database` with `dbName`, `adminPassword`, optional `pdbName` (pluggable database), and optional `dbBackupConfig` for automatic backups.

**Storage** -- Set `dataStorageSizeInGb` for the initial data volume (256 to 40960 GB depending on shape). `diskRedundancy` controls mirroring (`normal` = 2-way, `high` = 3-way) for BM systems. `storageVolumePerformanceMode` selects between `balanced` and `high_performance` I/O tiers.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `nsgIds` | `status.outputs.networkSecurityGroupId` |
| **OciSubnet** (optional, RAC) | `backupSubnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional, RAC) | `backupNetworkNsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `db_system_id` | OCID of the DB System | Monitoring, IAM policy scoping, resource management |
| `db_home_id` | OCID of the first DB Home | DB Home lifecycle management, patching |
| `database_id` | OCID of the initial database | Application connection configuration, DBA operations |
| `listener_port` | TCP port on which the database listener accepts connections | Application connection strings (default 1521) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single-node VM** -- A single-node Oracle Database on a VM flex shape with Standard Edition, automatic backups, and a pluggable database. The baseline for workloads needing full DBA-level control over an Oracle Database instance. Start from the **Single-Node VM** preset.

**Two-node RAC** -- A two-node Real Application Clusters deployment with Enterprise Edition Extreme Performance, fault domain separation, rolling maintenance for zero-downtime patching, a dedicated backup subnet, and 60-day backup retention. Start from the **Two-Node RAC** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this DB System
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the client network subnet and optionally the backup subnet for RAC
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for both client and backup VNICs