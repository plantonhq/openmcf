# GCP Cloud SQL Pulumi Module - Examples

Here are examples demonstrating how to use the `GcpCloudSql` API resource with the Pulumi module for deploying MySQL and PostgreSQL databases on Google Cloud SQL.

## Example 1: Basic MySQL Instance (Literal Value)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-basic
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  database_engine: MYSQL
  database_version: MYSQL_8_0
  tier: db-n1-standard-1
  storage_gb: 10
  root_password: SecurePassword123!
```

Deploy with:
```shell
openmcf pulumi up --manifest mysql-basic.yaml --stack myorg/platform/dev
```

## Example 2: Basic MySQL Instance (Value From Reference)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-basic
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  region: us-central1
  database_engine: MYSQL
  database_version: MYSQL_8_0
  tier: db-n1-standard-1
  storage_gb: 10
  root_password: SecurePassword123!
```

## Example 3: PostgreSQL with Private IP (Literal Values)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-private
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  database_engine: POSTGRESQL
  database_version: POSTGRES_15
  tier: db-n1-standard-2
  storage_gb: 20
  network:
    vpc_id:
      value: projects/my-gcp-project/global/networks/my-vpc
    private_ip_enabled: true
  root_password: SecurePassword123!
```

## Example 4: PostgreSQL with Private IP (Value From References)

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-private
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  region: us-central1
  database_engine: POSTGRESQL
  database_version: POSTGRES_15
  tier: db-n1-standard-2
  storage_gb: 20
  network:
    vpc_id:
      valueFrom:
        kind: GcpVpc
        name: main-vpc
        fieldPath: status.outputs.network_id
    private_ip_enabled: true
  root_password: SecurePassword123!
```

## Example 5: MySQL with High Availability

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-ha
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  database_engine: MYSQL
  database_version: MYSQL_8_0
  tier: db-n1-standard-2
  storage_gb: 50
  high_availability:
    enabled: true
    zone: us-central1-b
  backup:
    enabled: true
    start_time: "03:00"
    retention_days: 7
  root_password: SecurePassword123!
```

## Example 6: PostgreSQL Production Setup (Value From References)

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
  database_engine: POSTGRESQL
  database_version: POSTGRES_15
  tier: db-n1-highmem-4
  storage_gb: 100
  network:
    vpc_id:
      valueFrom:
        kind: GcpVpc
        name: production-vpc
        fieldPath: status.outputs.network_id
    private_ip_enabled: true
  high_availability:
    enabled: true
    zone: us-central1-c
  backup:
    enabled: true
    start_time: "02:00"
    retention_days: 30
  database_flags:
    max_connections: "200"
    shared_buffers: "262144"
  root_password: ProductionSecurePassword123!
```

## Example 7: MySQL with Public IP and Authorized Networks

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: mysql-public
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  database_engine: MYSQL
  database_version: MYSQL_8_0
  tier: db-n1-standard-1
  storage_gb: 10
  network:
    authorized_networks:
      - 203.0.113.0/24
      - 198.51.100.0/24
  root_password: SecurePassword123!
```

## Example 8: PostgreSQL with Custom Database Flags

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudSql
metadata:
  name: postgres-custom
spec:
  projectId:
    value: my-gcp-project
  region: us-central1
  database_engine: POSTGRESQL
  database_version: POSTGRES_15
  tier: db-n1-standard-2
  storage_gb: 30
  database_flags:
    max_connections: "150"
    shared_preload_libraries: "pg_stat_statements"
    log_statement: "all"
  backup:
    enabled: true
    start_time: "04:00"
    retention_days: 14
  root_password: SecurePassword123!
```

## Deployment Commands

### Preview Changes
```shell
openmcf pulumi preview --manifest gcpcloudsql.yaml --stack myorg/platform/dev
```

### Deploy Instance
```shell
openmcf pulumi up --manifest gcpcloudsql.yaml --stack myorg/platform/dev
```

### Refresh State
```shell
openmcf pulumi refresh --manifest gcpcloudsql.yaml --stack myorg/platform/dev
```

### Destroy Instance
```shell
openmcf pulumi destroy --manifest gcpcloudsql.yaml --stack myorg/platform/dev
```

## Notes

- **Root Password**: Use strong passwords with at least 8 characters including letters, numbers, and special characters.
- **Private IP**: Requires VPC peering to be configured between your VPC and the Cloud SQL service networking.
- **High Availability**: Provides 99.95% uptime SLA with automatic failover capabilities.
- **Database Flags**: Refer to MySQL or PostgreSQL documentation for available configuration flags.
- **Storage**: Minimum 10GB, maximum 65,536GB (65TB). SSD storage is used by default.
- **Backup Retention**: Longer retention periods increase storage costs but provide more recovery options.

## StringValueOrRef Fields

The `projectId` and `network.vpc_id` fields support both literal values and references:

### Using Literal Values
```yaml
projectId:
  value: my-gcp-project
```

### Using Value From References
```yaml
projectId:
  valueFrom:
    kind: GcpProject
    name: my-project
    fieldPath: status.outputs.project_id
```

This enables dynamic infrastructure composition where Cloud SQL instances can automatically reference project IDs and VPC networks from other deployed resources.

