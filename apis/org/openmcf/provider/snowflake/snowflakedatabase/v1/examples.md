## Minimal Production Database

A basic production database with standard configuration and 30-day Time Travel retention.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: prod-analytics
spec:
  name: prod_analytics
  comment: "Production analytics database"
  data_retention_time_in_days: 30
  is_transient: false
```

## Transient Staging Database

A cost-optimized transient database suitable for staging environments.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: staging-analytics
spec:
  name: staging_analytics
  comment: "Staging analytics database for testing"
  data_retention_time_in_days: 7
  is_transient: true
  drop_public_schema_on_creation: false
```

## Database with Iceberg Configuration

Database configured for Apache Iceberg tables with external catalog and storage.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: iceberg-datalake
spec:
  name: iceberg_datalake
  comment: "Data lake with Iceberg tables"
  catalog: iceberg_catalog
  external_volume: s3_external_volume
  storage_serialization_policy: "COMPATIBLE"
  replace_invalid_characters: true
  data_retention_time_in_days: 30
  is_transient: false
```

## High-Compliance Database

Database with maximum data retention and comprehensive logging for compliance requirements.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: compliance-finance
  labels:
    environment: production
    compliance: sox-pci
    criticality: tier-1
spec:
  name: finance_compliance
  comment: "SOX and PCI-DSS compliant finance database"
  data_retention_time_in_days: 90
  max_data_extension_time_in_days: 30
  is_transient: false
  drop_public_schema_on_creation: true
  log_level: "INFO"
  trace_level: "ON_EVENT"
  enable_console_output: false
  quoted_identifiers_ignore_case: false
  default_ddl_collation: "en_US"
```

## Task-Intensive ETL Database

Database optimized for task-based ETL workflows with retry and auto-suspend configuration.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: etl-pipeline
spec:
  name: etl_pipeline
  comment: "ETL pipeline database with task management"
  data_retention_time_in_days: 14
  is_transient: false
  suspend_task_after_num_failures: 5
  task_auto_retry_attempts: 3
  user_task:
    managed_initial_warehouse_size: "MEDIUM"
    minimum_trigger_interval_in_seconds: 60
    timeout_ms: 3600000
```

## Development Database with Debug Logging

Development database with verbose logging for troubleshooting and development.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: dev-analytics
spec:
  name: dev_analytics
  comment: "Development analytics database"
  data_retention_time_in_days: 1
  is_transient: true
  drop_public_schema_on_creation: false
  log_level: "DEBUG"
  trace_level: "ALWAYS"
  enable_console_output: true
  quoted_identifiers_ignore_case: true
```

## Multi-Region Iceberg Database

Database for multi-region data lake with optimized serialization for Snowflake-only access.

```yaml
apiVersion: snowflake.openmcf.org/v1
kind: SnowflakeDatabase
metadata:
  name: global-datalake
  labels:
    scope: global
    architecture: multi-region
spec:
  name: global_datalake
  comment: "Global multi-region data lake"
  catalog: global_iceberg_catalog
  external_volume: s3_multi_region_volume
  storage_serialization_policy: "OPTIMIZED"
  data_retention_time_in_days: 60
  is_transient: false
  quoted_identifiers_ignore_case: false
  default_ddl_collation: "en_US"
```

## CLI Workflows

### Validate Manifest

```bash
openmcf validate --manifest snowflake-database.yaml
```

### Deploy with Pulumi

```bash
openmcf pulumi up --manifest snowflake-database.yaml --stack org/project/stack
```

### Deploy with Terraform

```bash
openmcf tofu apply --manifest snowflake-database.yaml --auto-approve
```

### Check Database Status

```bash
openmcf get --manifest snowflake-database.yaml
```

### Update Database Configuration

```bash
# Edit your manifest file with desired changes
openmcf pulumi up --manifest snowflake-database.yaml --stack org/project/stack
```

### Destroy Database

```bash
openmcf pulumi destroy --manifest snowflake-database.yaml --stack org/project/stack
```
