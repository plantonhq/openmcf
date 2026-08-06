# AliCloudLogProject

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudLogProjectSpec defines the configuration for an Alibaba Cloud
Simple Log Service (SLS) project with optional bundled log stores.

An SLS project is the top-level container for log data in Alibaba Cloud.
It groups related log stores and provides a namespace for log collection,
querying, and analysis. Each project is region-scoped and must have a
globally unique name.

This component bundles the project with its log stores and full-text indexes
(per DD07) because a project without stores is an empty shell, and stores
without indexes are unsearchable.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudLogProject
metadata:
  name: alicloudlogproject-demo
spec:
  region: cn-hangzhou
  projectName: planton-demo-logs
  logStores:
    - name: app-logs
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.projectName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.logStores` | `[]AliCloudLogStore` |  |  |  |
| `spec.logStores[].name` | `string` | yes |  |  |
| `spec.logStores[].retentionDays` | `int32` |  | `30` |  |
| `spec.logStores[].shardCount` | `int32` |  | `2` |  |
| `spec.logStores[].autoSplit` | `bool` |  | `true` |  |
| `spec.logStores[].maxSplitShardCount` | `int32` |  | `64` |  |
| `spec.logStores[].enableIndex` | `bool` |  | `true` |  |
| `spec.logStores[].appendMeta` | `bool` |  | `true` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the SLS project will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.projectName

`string` · required

SLS project name. Must be globally unique within Alibaba Cloud.
3-63 characters, lowercase letters, digits, and hyphens only.
Must start and end with a letter or digit.

- rule: {"required":true,"string":{"minLen":"3","maxLen":"63"}}

### spec.description

`string`

Human-readable description of the project.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the project is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the SLS project.

### spec.logStores

`[]AliCloudLogStore`

Log stores to create within this project.
Each log store is an independent unit of log collection and storage.

### spec.logStores[].name

`string` · required

Log store name. Must be unique within the project.
3-63 characters, lowercase letters, digits, hyphens, and underscores.

- rule: {"required":true,"string":{"minLen":"3","maxLen":"63"}}

### spec.logStores[].retentionDays

`int32` · optional (explicit presence)

Data retention period in days. Range: 1-3650.
Set to 3650 for permanent retention.
Default: 30

- default: `30`
- rule: {"int32":{"lte":3650,"gte":1}}

### spec.logStores[].shardCount

`int32` · optional (explicit presence)

Number of shards for the log store. More shards = higher write throughput.
Default: 2

- default: `2`
- rule: {"int32":{"lte":256,"gte":1}}

### spec.logStores[].autoSplit

`bool` · optional (explicit presence)

Enable automatic shard splitting when write throughput exceeds shard capacity.
Default: true

- default: `true`

### spec.logStores[].maxSplitShardCount

`int32` · optional (explicit presence)

Maximum number of shards after auto-splitting. Only effective when auto_split is true.
Default: 64

- default: `64`
- rule: {"int32":{"lte":256,"gte":1}}

### spec.logStores[].enableIndex

`bool` · optional (explicit presence)

Create a full-text search index for this log store.
When true, a default full-text index is created with case-insensitive matching
and standard tokenization, making logs immediately searchable.
Default: true

- default: `true`

### spec.logStores[].appendMeta

`bool` · optional (explicit presence)

Automatically append log receive time and client IP as metadata fields.
Default: true

- default: `true`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudLogProject, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.project_name` | `string` | The SLS project name (also serves as the project identifier in SLS APIs). |
| `status.outputs.project_id` | `string` | The SLS project ID. |
| `status.outputs.log_store_names` | `map<string, string>` | Map of log store names created within the project. Key: the log store name as specified in spec. Value: the log store name (identical). This output exists so downstream components can reference specific log store names via StringValueOrRef (e.g., AckManagedCluster referencing a log store for audit logs). |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudFunction | `spec.logConfig.project` | `status.outputs.project_name` |
| AliCloudKubernetesCluster | `spec.logging.controlPlaneLogProject` | `status.outputs.project_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
