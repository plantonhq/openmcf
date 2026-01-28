# GCP Cloud SQL API - Examples

Here are examples of how to create and configure a **GcpCloudSql** API resource using the OpenMCF CLI. The examples cover MySQL and PostgreSQL instances with various configuration options.

## Create using CLI

First, create a YAML file using the examples provided below. After the YAML file is created, you can apply the configuration using the following command:

```shell
openmcf pulumi up --manifest <yaml-path> --stack <org>/<stack-name>/<environment>
```

Or using Terraform:

```shell
openmcf tofu apply --manifest <yaml-path> --auto-approve
```

## Basic MySQL Example (Literal Value)

This example demonstrates how to create a basic MySQL 8.0 Cloud SQL instance with a hardcoded project ID.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-db
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-n1-standard-1
  storageGb: 10
  rootPassword: MySecurePassword123!
```

## Basic MySQL Example (Value From Reference)

This example demonstrates referencing a GcpProject resource for the project ID.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-db
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-n1-standard-1
  storageGb: 10
  rootPassword: MySecurePassword123!
```

## Basic PostgreSQL Example

This example demonstrates how to create a basic PostgreSQL 15 Cloud SQL instance.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-db
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-n1-standard-2
  storageGb: 20
  rootPassword: MySecurePassword123!
```

## MySQL with Private IP and VPC (Literal Values)

This example creates a MySQL instance with private IP connectivity through a VPC network using literal values.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-private
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-n1-standard-1
  storageGb: 10
  network:
    vpcId:
      value: projects/my-gcp-project/global/networks/my-vpc
    privateIpEnabled: true
  rootPassword: MySecurePassword123!
```

## MySQL with Private IP and VPC (Value From References)

This example creates a MySQL instance with private IP connectivity by referencing existing GcpProject and GcpVpc resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-private
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-n1-standard-1
  storageGb: 10
  network:
    vpcId:
      valueFrom:
        kind: GcpVpc
        name: main-vpc
        fieldPath: status.outputs.network_id
    privateIpEnabled: true
  rootPassword: MySecurePassword123!
```

## PostgreSQL with Public IP and Authorized Networks

This example creates a PostgreSQL instance with public IP access restricted to specific CIDR ranges.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-public
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-n1-standard-1
  storageGb: 10
  network:
    authorizedNetworks:
      - 203.0.113.0/24
      - 198.51.100.0/24
  rootPassword: MySecurePassword123!
```

## MySQL with High Availability

This example creates a highly available MySQL instance with automatic failover to a secondary zone.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-ha
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-n1-standard-2
  storageGb: 50
  highAvailability:
    enabled: true
    zone: us-central1-b
  rootPassword: MySecurePassword123!
```

## PostgreSQL with Automated Backups

This example creates a PostgreSQL instance with automated daily backups configured.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-backup
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-n1-standard-1
  storageGb: 20
  backup:
    enabled: true
    startTime: "03:00"
    retentionDays: 7
  rootPassword: MySecurePassword123!
```

## Production-Grade PostgreSQL Instance (Literal Values)

This comprehensive example creates a production-ready PostgreSQL instance with high availability, private networking, automated backups, and custom database flags using literal values.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-production
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-n1-highmem-4
  storageGb: 100
  network:
    vpcId:
      value: projects/my-gcp-project/global/networks/production-vpc
    privateIpEnabled: true
  highAvailability:
    enabled: true
    zone: us-central1-c
  backup:
    enabled: true
    startTime: "02:00"
    retentionDays: 30
  databaseFlags:
    maxConnections: "200"
    sharedBuffers: "262144"
    effectiveCacheSize: "786432"
  rootPassword: ProductionSecurePassword123!
```

## Production-Grade PostgreSQL Instance (Value From References)

This comprehensive example creates a production-ready PostgreSQL instance by referencing existing GcpProject and GcpVpc resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-production
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: prod-project
      fieldPath: status.outputs.project_id
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-n1-highmem-4
  storageGb: 100
  network:
    vpcId:
      valueFrom:
        kind: GcpVpc
        name: production-vpc
        fieldPath: status.outputs.network_id
    privateIpEnabled: true
  highAvailability:
    enabled: true
    zone: us-central1-c
  backup:
    enabled: true
    startTime: "02:00"
    retentionDays: 30
  databaseFlags:
    maxConnections: "200"
    sharedBuffers: "262144"
    effectiveCacheSize: "786432"
  rootPassword: ProductionSecurePassword123!
```

## MySQL with Database Flags

This example creates a MySQL instance with custom database configuration flags.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-custom
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-n1-standard-2
  storageGb: 30
  databaseFlags:
    maxConnections: "150"
    innodbBufferPoolSize: "268435456"
    slowQueryLog: "on"
  rootPassword: MySecurePassword123!
```

## Large Storage PostgreSQL Instance

This example creates a PostgreSQL instance optimized for large datasets.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-large-storage
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_15
  tier: db-n1-highmem-8
  storageGb: 1000
  network:
    vpcId:
      value: projects/my-gcp-project/global/networks/data-vpc
    privateIpEnabled: true
  backup:
    enabled: true
    startTime: "01:00"
    retentionDays: 14
  rootPassword: MySecurePassword123!
```

## Development MySQL Instance

This example creates a small, cost-effective MySQL instance suitable for development environments.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-dev
spec:
  projectId:
    value: my-dev-project
  region: us-central1
  databaseEngine: MYSQL
  databaseVersion: MYSQL_8_0
  tier: db-f1-micro
  storageGb: 10
  rootPassword: DevPassword123!
```

## Notes

- **Root Password**: Always use strong passwords with at least 8 characters including letters, numbers, and special characters.
- **Private IP**: Requires VPC peering to be configured between your VPC and the Cloud SQL service.
- **High Availability**: Increases cost but provides better uptime and automatic failover.
- **Backup Retention**: Longer retention periods increase storage costs but provide more recovery options.
- **Database Flags**: Refer to MySQL or PostgreSQL documentation for available flags and their values.
- **Instance Tiers**: Choose appropriate tier based on CPU and memory requirements of your workload.
- **Storage**: Start with minimum required and enable automatic storage increase in production.

## StringValueOrRef Fields

The following fields support both literal values and references to other resources:

### projectId
Use a literal value or reference a GcpProject resource:
```yaml
# Literal value
projectId:
  value: my-gcp-project

# Reference to GcpProject
projectId:
  valueFrom:
    kind: GcpProject
    name: my-project
    fieldPath: status.outputs.project_id
```

### network.vpcId
Use a literal value or reference a GcpVpc resource:
```yaml
# Literal value
network:
  vpcId:
    value: projects/my-project/global/networks/my-vpc
  privateIpEnabled: true

# Reference to GcpVpc
network:
  vpcId:
    valueFrom:
      kind: GcpVpc
      name: my-vpc
      fieldPath: status.outputs.network_id
  privateIpEnabled: true
```

