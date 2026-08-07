# OciObjectStorageBucket

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciObjectStorageBucketSpec defines the specification for an OCI Object
Storage bucket -- a durable, scalable object store with built-in retention
rules, lifecycle management, and cross-region replication.

Retention rules are managed inline on the bucket resource (max 100).
Lifecycle rules control automatic archival, tiering, and deletion of
objects based on age and name patterns. Replication policies enable
cross-region disaster recovery by asynchronously copying objects to
a destination bucket in another OCI region.

Excluded from v1:
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels
  - pre-authenticated requests -- separate operational concern
  - private endpoints -- separate networking concern
  - individual objects -- application-level concern

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.namespace` | `string` | yes |  |  |
| `spec.name` | `string` | yes |  |  |
| `spec.accessType` | `enum` |  |  |  |
| `spec.storageTier` | `enum` |  |  |  |
| `spec.versioning` | `enum` |  |  |  |
| `spec.autoTiering` | `enum` |  |  |  |
| `spec.objectEventsEnabled` | `bool` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.metadata` | `map<string, string>` |  |  |  |
| `spec.retentionRules` | `[]RetentionRule` |  |  |  |
| `spec.retentionRules[].displayName` | `string` | yes |  |  |
| `spec.retentionRules[].duration` | `Duration` |  |  |  |
| `spec.retentionRules[].duration.timeAmount` | `int64` |  |  |  |
| `spec.retentionRules[].duration.timeUnit` | `enum` |  |  |  |
| `spec.retentionRules[].timeRuleLocked` | `string` |  |  |  |
| `spec.lifecycleRules` | `[]LifecycleRule` |  |  |  |
| `spec.lifecycleRules[].name` | `string` | yes |  |  |
| `spec.lifecycleRules[].action` | `enum` |  |  |  |
| `spec.lifecycleRules[].isEnabled` | `bool` |  |  |  |
| `spec.lifecycleRules[].timeAmount` | `int64` |  |  |  |
| `spec.lifecycleRules[].timeUnit` | `enum` |  |  |  |
| `spec.lifecycleRules[].target` | `string` |  |  |  |
| `spec.lifecycleRules[].objectNameFilter` | `ObjectNameFilter` |  |  |  |
| `spec.lifecycleRules[].objectNameFilter.inclusionPatterns` | `[]string` |  |  |  |
| `spec.lifecycleRules[].objectNameFilter.inclusionPrefixes` | `[]string` |  |  |  |
| `spec.lifecycleRules[].objectNameFilter.exclusionPatterns` | `[]string` |  |  |  |
| `spec.replicationPolicies` | `[]ReplicationPolicy` |  |  |  |
| `spec.replicationPolicies[].name` | `string` | yes |  |  |
| `spec.replicationPolicies[].destinationBucketName` | `string` | yes |  |  |
| `spec.replicationPolicies[].destinationRegionName` | `string` | yes |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the bucket will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.namespace

`string` · required

Object Storage namespace for the tenancy. This is a unique identifier
assigned to each tenancy (e.g. "axe1234abc"). Retrieve it via
`oci os ns get` or from the OCI Console.

- rule: {"string":{"minLen":"1"}}

### spec.name

`string` · required

Bucket name. Must be unique within the namespace. Valid characters:
uppercase/lowercase letters, numbers, hyphens, underscores, periods.
Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.accessType

`enum`

Type of public access enabled on the bucket.
When unspecified, defaults to no_public_access.

Allowed values (use exactly as shown):

- `access_type_unspecified`
- `no_public_access`
- `object_read`
- `object_read_without_list`

### spec.storageTier

`enum`

Storage tier for the bucket. Immutable after creation.
When unspecified, defaults to standard.

Allowed values (use exactly as shown):

- `storage_tier_unspecified`
- `standard`
- `archive`

### spec.versioning

`enum`

Versioning status. Objects in a version-enabled bucket are protected
from overwrites and deletions by maintaining version history.
On create: enabled or disabled. On update: enabled or suspended.

Allowed values (use exactly as shown):

- `versioning_unspecified`
- `enabled`
- `disabled`
- `suspended`

### spec.autoTiering

`enum`

Auto-tiering automatically transitions objects between Standard
and InfrequentAccess tiers based on access patterns.

Allowed values (use exactly as shown):

- `auto_tiering_unspecified`
- `auto_tiering_disabled`
- `infrequent_access`

### spec.objectEventsEnabled

`bool`

When true, events are emitted for object state changes in this bucket
via the OCI Events service.

### spec.kmsKeyId

`string | valueFrom`

OCID of a KMS master encryption key for server-side encryption.
When unset, Oracle-managed keys are used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.metadata

`map<string, string>`

User-defined metadata as key-value pairs. Keys must be lowercase.
Total size limit is 4KB.

### spec.retentionRules

`[]RetentionRule`

Retention rules enforce minimum retention periods on objects.
A maximum of 100 rules are supported per bucket. Rules take
effect within approximately 30 seconds.

### spec.retentionRules[].displayName

`string` · required

User-specified name for the retention rule. Must be unique
within the bucket. Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.retentionRules[].duration

`Duration`

Retention duration. When omitted, the rule applies indefinitely.

### spec.retentionRules[].duration.timeAmount

`int64`

Time amount, interpreted in units defined by time_unit.

- rule: {"int64":{"gte":"1"}}

### spec.retentionRules[].duration.timeUnit

`enum`

Unit of time for time_amount.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `time_unit_unspecified`
- `days`
- `years`

### spec.retentionRules[].timeRuleLocked

`string`

RFC 3339 datetime after which this rule becomes locked.
Once locked, the rule can only be deleted by deleting the bucket,
and only duration increases are allowed.

### spec.lifecycleRules

`[]LifecycleRule`

Lifecycle rules automate object transitions and deletions based
on age. Managed as a single lifecycle policy resource on the bucket.

- rule: lifecycle rules with abort action must target 'multipart-uploads'
- rule: object_name_filter is not valid when target is 'multipart-uploads'

### spec.lifecycleRules[].name

`string` · required

Rule name. Must be unique within the lifecycle policy.

- rule: {"string":{"minLen":"1"}}

### spec.lifecycleRules[].action

`enum`

Action to perform on matching objects.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `lifecycle_action_unspecified`
- `lifecycle_archive`
- `lifecycle_infrequent_access`
- `lifecycle_delete`
- `lifecycle_abort`

### spec.lifecycleRules[].isEnabled

`bool`

Whether this rule is currently active.

### spec.lifecycleRules[].timeAmount

`int64`

Age threshold. Objects older than this are acted upon.

- rule: {"int64":{"gte":"1"}}

### spec.lifecycleRules[].timeUnit

`enum`

Unit for time_amount.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `time_unit_unspecified`
- `days`
- `years`

### spec.lifecycleRules[].target

`string`

Target object type. Valid values: "objects" (default),
"multipart-uploads", "previous-object-versions".
Uses plain strings because the values contain hyphens.

- rule: {"string":{"in":["objects","multipart-uploads","previous-object-versions",""]}}

### spec.lifecycleRules[].objectNameFilter

`ObjectNameFilter`

Filter to narrow which objects the rule applies to.
Not valid when target is "multipart-uploads".

### spec.lifecycleRules[].objectNameFilter.inclusionPatterns

`[]string`

Glob patterns to include. An empty list includes all objects.

### spec.lifecycleRules[].objectNameFilter.inclusionPrefixes

`[]string`

Object name prefixes to include. Kept for backward compatibility;
prefer inclusion_patterns.

### spec.lifecycleRules[].objectNameFilter.exclusionPatterns

`[]string`

Glob patterns to exclude. Takes precedence over inclusions.

### spec.replicationPolicies

`[]ReplicationPolicy`

Cross-region replication policies. Each policy replicates objects
to a destination bucket in another OCI region. Destination buckets
must exist before creating the policy. All fields are immutable
after creation.

### spec.replicationPolicies[].name

`string` · required

Policy name.

- rule: {"string":{"minLen":"1"}}

### spec.replicationPolicies[].destinationBucketName

`string` · required

Name of the destination bucket. Must already exist in the
destination region.

- rule: {"string":{"minLen":"1"}}

### spec.replicationPolicies[].destinationRegionName

`string` · required

OCI region identifier for the destination (e.g. "us-ashburn-1").

- rule: {"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciObjectStorageBucket, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.bucket_id` | `string` | OCID of the Object Storage bucket. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](../README.md)
