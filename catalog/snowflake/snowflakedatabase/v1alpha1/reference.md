# SnowflakeDatabase

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `snowflake.planton.dev/v1alpha1`

snowflake-database spec

## Example

```yaml
apiVersion: snowflake.planton.dev/v1alpha1
kind: SnowflakeDatabase
metadata:
  name: test-database
  id: snowdb-test-001
  org: engineering
  env: development
  labels:
    team: data-platform
    component: snowflake-database
    managed-by: planton
spec:
  # Database name (required)
  name: TEST_DATABASE_DEV
  
  # Database description
  comment: "Test database for development and CI/CD - transient for cost savings"
  
  # Cost optimization: transient database eliminates 7-day Fail-safe period
  is_transient: true
  
  # Minimal Time Travel retention for dev environment (1 day)
  data_retention_time_in_days: 1
  
  # Iceberg configuration (optional)
  catalog: ""
  external_volume: ""
  
  # Collation settings (optional)
  default_ddl_collation: ""
  
  # Schema management
  drop_public_schema_on_creation: false
  
  # Logging and debugging
  enable_console_output: true
  log_level: "INFO"
  trace_level: "ON_EVENT"
  
  # Data management settings
  max_data_extension_time_in_days: 0
  quoted_identifiers_ignore_case: false
  replace_invalid_characters: false
  
  # Storage configuration
  storage_serialization_policy: ""
  
  # Task management settings
  suspend_task_after_num_failures: 5
  task_auto_retry_attempts: 2
  
  # User task configuration
  user_task:
    managed_initial_warehouse_size: "XSMALL"
    minimum_trigger_interval_in_seconds: 60
    timeout_ms: 3600000
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.catalog` | `string` |  |  |  |
| `spec.comment` | `string` |  |  |  |
| `spec.dataRetentionTimeInDays` | `int32` |  |  |  |
| `spec.defaultDdlCollation` | `string` |  |  |  |
| `spec.dropPublicSchemaOnCreation` | `bool` |  |  |  |
| `spec.enableConsoleOutput` | `bool` |  |  |  |
| `spec.externalVolume` | `string` |  |  |  |
| `spec.isTransient` | `bool` |  |  |  |
| `spec.logLevel` | `string` |  |  |  |
| `spec.maxDataExtensionTimeInDays` | `int32` |  |  |  |
| `spec.name` | `string` |  |  |  |
| `spec.quotedIdentifiersIgnoreCase` | `bool` |  |  |  |
| `spec.replaceInvalidCharacters` | `bool` |  |  |  |
| `spec.storageSerializationPolicy` | `string` |  |  |  |
| `spec.suspendTaskAfterNumFailures` | `int32` |  |  |  |
| `spec.taskAutoRetryAttempts` | `int32` |  |  |  |
| `spec.traceLevel` | `string` |  |  |  |
| `spec.userTask` | `SnowflakeDatabaseUserTask` |  |  |  |
| `spec.userTask.managedInitialWarehouseSize` | `string` |  |  |  |
| `spec.userTask.minimumTriggerIntervalInSeconds` | `int32` |  |  |  |
| `spec.userTask.timeoutMs` | `int32` |  |  |  |

## Field Details

### spec.catalog

`string`

The database parameter that specifies the default catalog to use for Iceberg tables
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#catalog_yaml
https://docs.snowflake.com/en/sql-reference/parameters#catalog

### spec.comment

`string`

Specifies a comment for the database
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#comment_yaml

### spec.dataRetentionTimeInDays

`int32`

Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database,
as well as specifying the default Time Travel retention time for all schemas created in the database.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#dataretentiontimeindays_yaml
https://docs.snowflake.com/en/user-guide/data-time-travel

### spec.defaultDdlCollation

`string`

Specifies a default collation specification for all schemas and tables added to the database.
It can be overridden on schema or table level.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#defaultddlcollation_yaml
https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification

### spec.dropPublicSchemaOnCreation

`bool`

Specifies whether to drop public schema on creation or not. Modifying the parameter after database is
already created won't have any effect.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#droppublicschemaoncreation_yaml

### spec.enableConsoleOutput

`bool`

If true, enables stdout/stderr fast path logging for anonymous stored procedures.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#enableconsoleoutput_yaml

### spec.externalVolume

`string`

The database parameter that specifies the default external volume to use for Iceberg tables
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#externalvolume_yaml
https://docs.snowflake.com/en/sql-reference/parameters#external-volume

### spec.isTransient

`bool`

Specifies the database as transient. Transient databases do not have a Fail-safe period so they do not incur
additional storage costs once they leave Time Travel; however, this means they are also not protected by
Fail-safe in the event of a data loss.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#istransient_yaml

### spec.logLevel

`string`

Specifies the severity level of messages that should be ingested and made available in the active event table.
Valid options are: [TRACE DEBUG INFO WARN ERROR FATAL OFF]. Messages at the specified level (and at more severe levels) are ingested.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#loglevel_yaml
https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level

### spec.maxDataExtensionTimeInDays

`int32`

Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period
for tables in the database to prevent streams on the tables from becoming stale.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#maxdataextensiontimeindays_yaml
https://docs.snowflake.com/en/sql-reference/parameters.html#label-max-data-extension-time-in-days

### spec.name

`string`

Specifies the identifier for the database; must be unique for your account. As a best practice for Database
Replication and Failover, it is recommended to give each secondary database the same name as its primary database.
This practice supports referencing fully-qualified objects (i.e. '\n\n.\n\n.\n\n') by other objects in the
same database, such as querying a fully-qualified table name in a view. If a secondary database has a
different name from the primary database, then these object references would break in the secondary database.
Due to technical limitations (read more here), avoid using the following characters: |, ., (, ), "

### spec.quotedIdentifiersIgnoreCase

`bool`

If true, the case of quoted identifiers is ignored
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#quotedidentifiersignorecase_yaml
https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case

### spec.replaceInvalidCharacters

`bool`

Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�) in query results
for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#replaceinvalidcharacters_yaml
https://docs.snowflake.com/en/sql-reference/parameters#replace-invalid-characters

### spec.storageSerializationPolicy

`string`

The storage serialization policy for Iceberg tables that use Snowflake as the catalog.
Valid options are: [COMPATIBLE OPTIMIZED]. COMPATIBLE: Snowflake performs encoding and compression of data
files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and
compression of data files that ensures the best table performance within Snowflake.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#storageserializationpolicy_yaml
https://docs.snowflake.com/en/sql-reference/parameters#storage-serialization-policy

### spec.suspendTaskAfterNumFailures

`int32`

How many times a task must fail in a row before it is automatically suspended. 0 disables auto-suspending.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#suspendtaskafternumfailures_yaml
https://docs.snowflake.com/en/sql-reference/parameters#suspend-task-after-num-failures

### spec.taskAutoRetryAttempts

`int32`

Maximum automatic retries allowed for a user task
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#taskautoretryattempts_yaml
https://docs.snowflake.com/en/sql-reference/parameters#task-auto-retry-attempts

### spec.traceLevel

`string`

Controls how trace events are ingested into the event table. Valid options are: [ALWAYS ON*EVENT OFF]
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#tracelevel_yaml
https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level

### spec.userTask

`SnowflakeDatabaseUserTask`

snowflake database user task

### spec.userTask.managedInitialWarehouseSize

`string`

The initial size of warehouse to use for managed warehouses in the absence of history.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#usertaskmanagedinitialwarehousesize_yaml
https://docs.snowflake.com/en/sql-reference/parameters#user-task-managed-initial-warehouse-size

### spec.userTask.minimumTriggerIntervalInSeconds

`int32`

Minimum amount of time between Triggered Task executions in seconds.
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#usertaskminimumtriggerintervalinseconds_yaml

### spec.userTask.timeoutMs

`int32`

User task execution timeout in milliseconds
https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#usertasktimeoutms_yaml
https://docs.snowflake.com/en/sql-reference/parameters#user-task-timeout-ms

## Outputs

Reference an output from another manifest as `valueFrom: {kind: SnowflakeDatabase, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | The provider-assigned unique ID for this managed resource. https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#id_yaml |
| `status.outputs.name` | `string` | The fully-qualified name of the created database https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#name_yaml |
| `status.outputs.owner` | `string` | The owner role of the database https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#owner_yaml |
| `status.outputs.created_on` | `string` | Timestamp when the database was created https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#createdon_yaml |
| `status.outputs.is_transient` | `string` | Indicates if the database is transient ("true" or "false") https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#istransient_yaml |
| `status.outputs.data_retention_time_in_days` | `string` | Configured data retention time in days (as string) https://www.pulumi.com/registry/packages/snowflake/api-docs/database/#dataretentiontimeindays_yaml |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
