# AliCloudStorageBucket

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudStorageBucketSpec defines the configuration for an Alibaba Cloud Object
Storage Service (OSS) bucket.

OSS is a cloud-native, S3-compatible object storage service. A bucket is the
top-level container for objects. Bucket names are globally unique across all
Alibaba Cloud accounts.

This component creates a single OSS bucket with optional versioning,
server-side encryption, lifecycle management, CORS rules, and access logging.
Features like bucket policies and referer configuration are managed as
separate provider resources when needed.

Important: storage_class and redundancy_type are immutable after creation --
changing either value requires destroying and recreating the bucket.

Provider resources:
  Terraform: alicloud_oss_bucket
  Pulumi:    oss.Bucket

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudStorageBucket
metadata:
  name: alicloudstoragebucket-demo
spec:
  region: cn-hangzhou
  bucketName: alicloudstoragebucket-demo
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.bucketName` | `string` | yes |  |  |
| `spec.acl` | `string` |  | `private` |  |
| `spec.storageClass` | `string` |  | `Standard` |  |
| `spec.redundancyType` | `string` |  | `LRS` |  |
| `spec.versioningEnabled` | `bool` |  |  |  |
| `spec.serverSideEncryption` | `AliCloudStorageBucketEncryption` |  |  |  |
| `spec.serverSideEncryption.sseAlgorithm` | `string` | yes |  |  |
| `spec.serverSideEncryption.kmsMasterKeyId` | `string` |  |  |  |
| `spec.lifecycleRules` | `[]AliCloudStorageBucketLifecycleRule` |  |  |  |
| `spec.lifecycleRules[].prefix` | `string` |  |  |  |
| `spec.lifecycleRules[].enabled` | `bool` |  |  |  |
| `spec.lifecycleRules[].expirationDays` | `int32` |  |  |  |
| `spec.lifecycleRules[].transitions` | `[]AliCloudStorageBucketLifecycleTransition` |  |  |  |
| `spec.lifecycleRules[].transitions[].days` | `int32` |  |  |  |
| `spec.lifecycleRules[].transitions[].storageClass` | `string` | yes |  |  |
| `spec.lifecycleRules[].abortMultipartUploadDays` | `int32` |  |  |  |
| `spec.lifecycleRules[].noncurrentVersionExpirationDays` | `int32` |  |  |  |
| `spec.corsRules` | `[]AliCloudStorageBucketCorsRule` |  |  |  |
| `spec.corsRules[].allowedOrigins` | `[]string` | yes |  |  |
| `spec.corsRules[].allowedMethods` | `[]string` | yes |  |  |
| `spec.corsRules[].allowedHeaders` | `[]string` |  |  |  |
| `spec.corsRules[].exposeHeaders` | `[]string` |  |  |  |
| `spec.corsRules[].maxAgeSeconds` | `int32` |  |  |  |
| `spec.logging` | `AliCloudStorageBucketLogging` |  |  |  |
| `spec.logging.targetBucket` | `string` | yes |  |  |
| `spec.logging.targetPrefix` | `string` |  |  |  |
| `spec.forceDestroy` | `bool` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the bucket will be created.
The region determines the physical location of stored data and the
available endpoints.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.bucketName

`string` · required

Globally unique bucket name. Must be 3-63 characters, lowercase letters,
digits, and hyphens only. Cannot start or end with a hyphen.
OSS uses the bucket name as the resource ID.

- rule: {"required":true,"string":{"minLen":"3","maxLen":"63"}}

### spec.acl

`string` · optional (explicit presence)

Bucket access control list.
Controls the default access permissions for all objects in the bucket.

Uses the provider's exact values: "private", "public-read", "public-read-write".
Default: "private"

- default: `private`
- rule: acl must be one of: private, public-read, public-read-write

### spec.storageClass

`string` · optional (explicit presence)

Storage class for the bucket. Determines cost, availability, and access
latency characteristics. This is immutable after creation (ForceNew).

Uses the provider's exact values:
  "Standard"      -- frequent access, highest cost, lowest latency
  "IA"            -- infrequent access, lower cost, retrieval fee
  "Archive"       -- archival, minutes to restore
  "ColdArchive"   -- cold storage, hours to restore
  "DeepColdArchive" -- deep cold, 12-48h to restore
Default: "Standard"

- default: `Standard`
- rule: storage_class must be one of: Standard, IA, Archive, ColdArchive, DeepColdArchive

### spec.redundancyType

`string` · optional (explicit presence)

Data redundancy type. Controls how data is replicated within the region.
This is immutable after creation (ForceNew).

  "LRS" -- Locally Redundant Storage: 3 copies in a single zone
  "ZRS" -- Zone-Redundant Storage: 3 copies across multiple zones
Default: "LRS"

- default: `LRS`
- rule: redundancy_type must be one of: LRS, ZRS

### spec.versioningEnabled

`bool`

Enable object versioning on this bucket.
When true, OSS preserves all versions of every object, allowing recovery
from accidental overwrites or deletes. Combine with
lifecycle_rules.noncurrent_version_expiration_days to control version
retention costs.
Default: false

### spec.serverSideEncryption

`AliCloudStorageBucketEncryption`

Server-side encryption configuration.
When set, all objects stored in this bucket are encrypted at rest using
the specified algorithm. Omit to store objects without default encryption.

### spec.serverSideEncryption.sseAlgorithm

`string` · required

Encryption algorithm.
"AES256" uses OSS-managed keys (no additional configuration needed).
"KMS" uses Alibaba Cloud KMS; optionally specify kms_master_key_id for
a customer-managed key, or omit for the OSS default KMS key.

- rule: sse_algorithm must be one of: AES256, KMS
- rule: {"required":true}

### spec.serverSideEncryption.kmsMasterKeyId

`string`

KMS master key ID for customer-managed encryption.
Only applicable when sse_algorithm is "KMS". When omitted with KMS,
OSS uses the default service key.

### spec.lifecycleRules

`[]AliCloudStorageBucketLifecycleRule`

Lifecycle management rules for automatic object transitions and expiration.
Rules are evaluated independently; multiple rules can match the same object.
Maximum 1000 rules per bucket.

### spec.lifecycleRules[].prefix

`string`

Object name prefix that scopes this rule. An empty prefix applies
the rule to all objects in the bucket.

### spec.lifecycleRules[].enabled

`bool`

Whether this rule is active. Disabled rules are retained but not evaluated.

### spec.lifecycleRules[].expirationDays

`int32`

Number of days after object creation to expire (permanently delete) objects.
Set to 0 to skip expiration. Only one of date-based or days-based expiration
can be active; this component exposes the days-based approach as the standard
pattern.

### spec.lifecycleRules[].transitions

`[]AliCloudStorageBucketLifecycleTransition`

Storage class transitions. Objects are moved to a cheaper storage tier
after the specified number of days.

### spec.lifecycleRules[].transitions[].days

`int32`

Number of days after object creation to transition to the target
storage class.

- rule: {"int32":{"gt":0}}

### spec.lifecycleRules[].transitions[].storageClass

`string` · required

Target storage class for the transition.
Transitions must move objects to a colder tier (cost-decreasing order):
  Standard -> IA -> Archive -> ColdArchive -> DeepColdArchive

- rule: storage_class must be one of: IA, Archive, ColdArchive, DeepColdArchive
- rule: {"required":true}

### spec.lifecycleRules[].abortMultipartUploadDays

`int32`

Number of days after which incomplete multipart uploads are automatically
aborted. This prevents orphaned upload parts from consuming storage.
Recommended: 7 days for most workloads.

### spec.lifecycleRules[].noncurrentVersionExpirationDays

`int32`

Number of days after which noncurrent object versions are expired.
Only meaningful when versioning is enabled on the bucket.
Recommended: 30-90 days depending on recovery requirements.

### spec.corsRules

`[]AliCloudStorageBucketCorsRule`

Cross-Origin Resource Sharing rules.
Required when browser-based clients need to access OSS directly.
Maximum 10 rules per bucket.

### spec.corsRules[].allowedOrigins

`[]string` · required

Origins allowed to make cross-origin requests. Use "*" to allow all origins.
Example: ["https://example.com", "https://app.example.com"]

- rule: {"repeated":{"minItems":"1"}}

### spec.corsRules[].allowedMethods

`[]string` · required

HTTP methods allowed for cross-origin requests.
Valid values: GET, PUT, POST, DELETE, HEAD

- rule: {"repeated":{"minItems":"1"}}

### spec.corsRules[].allowedHeaders

`[]string`

Request headers that are allowed in a preflight OPTIONS request.
Use "*" to allow all headers.

### spec.corsRules[].exposeHeaders

`[]string`

Response headers that the browser is allowed to access.
Example: ["x-oss-request-id", "ETag"]

### spec.corsRules[].maxAgeSeconds

`int32`

Maximum time (in seconds) that the browser should cache the preflight
response. Reduces the number of OPTIONS requests for repeated access.

### spec.logging

`AliCloudStorageBucketLogging`

Server access logging configuration.
When set, OSS writes access logs to the specified target bucket.
Omit to disable access logging.

### spec.logging.targetBucket

`string` · required

Destination bucket for access log objects. Can be the same bucket or a
different bucket in the same region.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.logging.targetPrefix

`string`

Key prefix for log objects. Useful for organizing logs within the
target bucket. Example: "logs/my-bucket/"

### spec.forceDestroy

`bool`

Allow force-destroying a non-empty bucket.
When true, Terraform/Pulumi will delete all objects in the bucket before
destroying it. Use with caution in production.
Default: false

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the bucket is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the bucket.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudStorageBucket, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.bucket_name` | `string` | The bucket name. OSS uses the bucket name as the resource identifier. Consumers reference this to upload objects, configure logging targets, or set up replication sources. |
| `status.outputs.extranet_endpoint` | `string` | The public internet endpoint for accessing this bucket. Format: {bucket_name}.oss-{region}.aliyuncs.com Used by external clients and CDN origins. |
| `status.outputs.intranet_endpoint` | `string` | The VPC-internal endpoint for accessing this bucket. Format: {bucket_name}.oss-{region}-internal.aliyuncs.com Used by ECS instances, functions, and containers within the same region for zero-cost, low-latency access. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
