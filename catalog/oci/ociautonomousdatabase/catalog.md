# Autonomous Database on OCI

Deploys an Oracle Cloud Infrastructure Autonomous Database -- a fully managed database service supporting OLTP (Autonomous Transaction Processing), DW (Autonomous Data Warehouse), AJD (Autonomous JSON Database), APEX, and Lakehouse workloads. The database runs on shared (serverless) or dedicated Exadata infrastructure with ECPU or OCPU compute models, independent compute and storage auto-scaling, optional private endpoint networking, and configurable mTLS/TLS-only connections. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Autonomous Database** -- a single autonomous database resource in the specified compartment with the chosen workload type, compute model, storage allocation, and networking configuration
- **Private Endpoint** -- created only when `subnetId` is set; places the database on a private IP within the specified subnet, disabling public secure access
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the database

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the database in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For private endpoint access: a subnet and optionally up to 5 network security groups. Without a subnet, the database is accessible via public secure access endpoints.
- An administrator password (12-30 characters, must include uppercase, lowercase, and numeric). For production, use `secretId` to reference an OCI Vault secret instead of embedding the password in the manifest.

## Deploy

### Console

Open the deployment store, find **Autonomous Database on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Serverless OLTP** preset in the [Presets](#presets) tab to pre-populate a production-ready ATP database with private endpoint and auto-scaling.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciAutonomousDatabase
metadata:
  name: app-database
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  dbName: appdb
  dbWorkload: oltp
  computeModel: ecpu
  computeCount: 4
  dataStorageSizeInTbs: 1
  isAutoScalingEnabled: true
  adminPassword: "Ex4mpl3!Passw0rd"
```

```shell
planton apply -f autonomous-db.yaml
```

This creates a serverless ATP database with 4 ECPUs, 1 TB storage, and compute auto-scaling (bursts to 12 ECPUs). No private endpoint is configured -- the database is accessible via public secure access. Storage auto-scaling and Data Guard are not enabled.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the database to a compartment, subnet, and security groups deployed in the same InfraPipeline:

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
        name: db-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the database with the resolved values.

## Key Configuration

These are the most important decisions when configuring an Autonomous Database. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Workload type** -- Set `dbWorkload` to `oltp` for transactional applications (ATP), `dw` for analytic and reporting workloads (ADW), `ajd` for JSON document storage, `apex` for low-code APEX development, or `lh` for data lake analytics. The workload type optimizes the database engine's storage layout, query optimizer, and available features.

**Compute model and auto-scaling** -- Set `computeModel` to `ecpu` (recommended) or `ocpu` (legacy). `computeCount` sets the baseline compute units. When `isAutoScalingEnabled` is `true`, compute can burst to 3x the baseline during demand spikes. `isAutoScalingForStorageEnabled` independently controls storage auto-expansion.

**Private endpoint** -- Set `subnetId` to place the database on a private IP within a subnet. When set, public secure access endpoints are disabled. Add `nsgIds` to restrict traffic to specific security groups. When `subnetId` is omitted, the database is accessible via public secure access with client IP allowlisting via `whitelistedIps`.

**mTLS vs TLS-only** -- Set `isMtlsConnectionRequired` to `false` to accept standard TLS connections (simpler client configuration). Set to `true` to require wallet-based mutual TLS for all connections. TLS-only is sufficient when the database is behind a private endpoint.

**Database edition** -- Set `databaseEdition` to `standard_edition` for most workloads. Use `enterprise_edition` when you need partitioning, advanced compression, or advanced security features.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** (optional) | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `autonomous_database_id` | OCID of the autonomous database | Monitoring, IAM policy scoping, resource management |
| `connection_string_high` | High-priority connection string for latency-sensitive workloads | Application primary connection pool |
| `connection_string_medium` | Medium-priority connection string for typical application workloads | General-purpose application connections |
| `connection_string_low` | Low-priority connection string for batch and background workloads | ETL jobs, reporting, background tasks |
| `service_console_url` | URL of the OCI Service Console for this database | DBA administration and monitoring |
| `private_endpoint` | Private endpoint IP address (empty when no subnet is configured) | DNS records, firewall allowlists |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Serverless OLTP** -- An ATP database with ECPU compute, private endpoint networking, and auto-scaling for both compute and storage. The standard pattern for production application backends. Start from the **Serverless OLTP** preset.

**Free tier development** -- An Always Free database with 2 ECPUs and 20 GB usable storage for prototyping and learning at zero cost. Not suitable for production. Start from the **Free Tier Development** preset.

**Serverless data warehouse** -- An ADW database with Enterprise Edition for partitioning and advanced compression, optimized for analytic queries and reporting dashboards. Start from the **Serverless Data Warehouse** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this database
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet for private endpoint access
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for the private endpoint