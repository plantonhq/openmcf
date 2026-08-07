# ScalewayObjectBucket

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayObjectBucketSpec defines the specification for a Scaleway Object
Storage bucket.

Scaleway Object Storage is an S3-compatible service that allows you to
store and retrieve any amount of data from anywhere on the web. Buckets
are the top-level containers for objects and can be configured with
versioning, lifecycle rules, and CORS policies.

This is a **standalone resource** wrapping a single
`scaleway_object_bucket` Terraform resource. It covers the 80% use case
of creating a bucket with optional versioning, lifecycle rules, and CORS
for web applications.

Object Storage buckets are **regional** resources. Available regions:
  - "fr-par" (Paris, France)
  - "nl-ams" (Amsterdam, Netherlands)
  - "pl-waw" (Warsaw, Poland)

**Naming constraint:** Bucket names must be globally unique across all
Scaleway Object Storage (like AWS S3). The bucket name is derived from
`metadata.name`, so choose a name that is DNS-compatible and unlikely
to collide with other users.

**S3 compatibility:** Scaleway Object Storage implements the S3 API,
so you can use AWS CLI, s3cmd, rclone, and any S3-compatible SDK by
pointing the endpoint to `s3.<region>.scw.cloud`.

**Deferred features (not in v1):**
  - Bucket ACL (deprecated on the main resource; use BucketAcl resource)
  - Bucket Policy (JSON IAM policies)
  - Object Lock Configuration (retention/legal hold rules)
  - Website Configuration (static site hosting)
  These can be managed via separate Terraform resources or added in
  future versions of this kind.

**Composition pattern:** This is a leaf resource with no upstream
`StringValueOrRef` dependencies. Downstream resources (serverless
functions, containers, applications) can reference
`status.outputs.bucket_id` and `status.outputs.endpoint` for
S3-compatible access.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.versioningEnabled` | `bool` |  | `false` |  |
| `spec.objectLockEnabled` | `bool` |  |  |  |
| `spec.lifecycleRules` | `[]ScalewayObjectBucketLifecycleRule` |  |  |  |
| `spec.lifecycleRules[].id` | `string` | yes |  |  |
| `spec.lifecycleRules[].enabled` | `bool` |  |  |  |
| `spec.lifecycleRules[].prefix` | `string` |  |  |  |
| `spec.lifecycleRules[].tags` | `map<string, string>` |  |  |  |
| `spec.lifecycleRules[].expirationDays` | `int32` |  |  |  |
| `spec.lifecycleRules[].transitions` | `[]ScalewayObjectBucketLifecycleTransition` |  |  |  |
| `spec.lifecycleRules[].transitions[].days` | `int32` |  |  |  |
| `spec.lifecycleRules[].transitions[].storageClass` | `string` | yes |  |  |
| `spec.lifecycleRules[].abortIncompleteMultipartUploadDays` | `int32` |  |  |  |
| `spec.corsRules` | `[]ScalewayObjectBucketCorsRule` |  |  |  |
| `spec.corsRules[].allowedMethods` | `[]string` | yes |  |  |
| `spec.corsRules[].allowedOrigins` | `[]string` | yes |  |  |
| `spec.corsRules[].allowedHeaders` | `[]string` |  |  |  |
| `spec.corsRules[].exposeHeaders` | `[]string` |  |  |  |
| `spec.corsRules[].maxAgeSeconds` | `int32` |  |  |  |
| `spec.forceDestroy` | `bool` |  | `false` |  |

## Field Details

### spec.region

`string` · required

The Scaleway region where the bucket will be created.

Available regions: "fr-par", "nl-ams", "pl-waw"

Cannot be changed after creation (would require recreating the bucket
and migrating all objects).

Choose the region closest to your application for lowest latency.
Consider data residency requirements (EU data in EU regions).

- rule: {"required":true}

### spec.versioningEnabled

`bool`

Whether to enable S3-compatible object versioning.

When enabled, every PUT operation creates a new version of the
object. Previous versions are retained and can be retrieved by
version ID. DELETE operations insert a delete marker instead of
removing the object.

IMPORTANT: Once enabled, versioning cannot be fully disabled --
only suspended. Suspended versioning stops creating new versions
but retains all existing versions.

Use versioning for:
  - Protecting against accidental deletions
  - Maintaining audit trails of object changes
  - Enabling point-in-time recovery

Combine with lifecycle rules to expire old versions and control
storage costs.

Default: false

- default: `false`

### spec.objectLockEnabled

`bool`

Whether to enable S3 Object Lock on this bucket.

Object Lock enables WORM (Write Once Read Many) protection. When
enabled, objects can be locked to prevent deletion or modification
for a specified retention period.

IMPORTANT: Object Lock requires versioning to be enabled. This
constraint is enforced via CEL validation below.

IMPORTANT: Object Lock can only be enabled at bucket creation time.
It cannot be added to an existing bucket, and cannot be removed
once enabled.

Default retention rules are configured separately via the
`scaleway_object_bucket_lock_configuration` Terraform resource
(deferred to future version of this kind).

Default: false

### spec.lifecycleRules

`[]ScalewayObjectBucketLifecycleRule`

Lifecycle rules for automated object management.

Lifecycle rules define actions that Scaleway Object Storage applies
to groups of objects based on prefix, tags, or age. Common uses:

  - **Expire old objects**: Delete logs, temp files after N days
  - **Transition to cold storage**: Move infrequently accessed data
    to GLACIER or ONEZONE_IA to save costs
  - **Expire old versions**: Clean up old versions when versioning
    is enabled (requires versioning_enabled = true)
  - **Abort incomplete multipart uploads**: Clean up stale uploads

Rules are evaluated daily. Changes take effect within 24 hours.

Optional. If empty, no lifecycle automation is applied.

### spec.lifecycleRules[].id

`string` · required

Unique identifier for this rule.

Used for identification in Terraform state and in Scaleway's
lifecycle management API. Must be unique within the bucket.

- rule: {"string":{"minLen":"1"}}

### spec.lifecycleRules[].enabled

`bool`

Whether this rule is currently active.

Disabled rules are retained in configuration but not evaluated.
Useful for temporarily suspending a rule without deleting it.

### spec.lifecycleRules[].prefix

`string`

Object key prefix filter.

When set, the rule applies only to objects whose key starts with
this prefix. For example, "logs/" applies to all objects in the
logs directory.

Empty string or omitted = rule applies to all objects in the
bucket (subject to tag filter if present).

### spec.lifecycleRules[].tags

`map<string, string>`

Tag-based filter for objects.

When set, the rule applies only to objects that have ALL of these
tags. Combined with prefix for more specific targeting.

Example: {"environment": "staging"} applies only to objects tagged
with environment=staging.

### spec.lifecycleRules[].expirationDays

`int32`

Number of days after object creation to expire (delete) the object.

When an object reaches this age, it is automatically deleted.
For versioned buckets, a delete marker is inserted.

0 = no expiration (default).

### spec.lifecycleRules[].transitions

`[]ScalewayObjectBucketLifecycleTransition`

Storage class transitions.

Move objects to cheaper storage classes after a specified number
of days. Multiple transitions can be defined to create a tiered
storage lifecycle (e.g., Standard -> ONEZONE_IA -> GLACIER).

Available storage classes:
  - "GLACIER" -- Cold storage for archival data (minutes-to-hours
    retrieval). Significantly cheaper than Standard.
  - "ONEZONE_IA" -- Infrequent Access in a single zone. Cheaper
    than Standard but less redundant.

Transition days must increase across transitions (e.g., 30 days
to ONEZONE_IA, then 90 days to GLACIER).

### spec.lifecycleRules[].transitions[].days

`int32`

Number of days after object creation to transition.

Must be a positive integer. When defining multiple transitions in
a rule, each transition's days must be greater than the previous.

- rule: {"int32":{"gt":0}}

### spec.lifecycleRules[].transitions[].storageClass

`string` · required

Target storage class.

Supported values:
  - "GLACIER" -- Cold storage for archival (cheapest, slow retrieval)
  - "ONEZONE_IA" -- Infrequent Access (cheaper, single-zone redundancy)

- rule: {"required":true,"string":{"in":["GLACIER","ONEZONE_IA"]}}

### spec.lifecycleRules[].abortIncompleteMultipartUploadDays

`int32`

Number of days to abort incomplete multipart uploads.

Multipart uploads that are not completed within this period are
automatically aborted and their parts are deleted.

Recommended: 7 days. Incomplete multipart uploads accumulate
storage costs silently.

0 = incomplete uploads are never aborted (default).

### spec.corsRules

`[]ScalewayObjectBucketCorsRule`

Cross-Origin Resource Sharing (CORS) rules for the bucket.

CORS rules control which web origins can make cross-origin requests
to the bucket. This is essential when:
  - A web application uploads files directly to the bucket
  - Browser-based JavaScript reads objects from the bucket
  - A frontend SPA hosted on one domain accesses bucket content

Without CORS rules, browsers block cross-origin requests to the
bucket's S3 endpoint.

Optional. If empty, no CORS headers are served (browsers will
block cross-origin requests).

### spec.corsRules[].allowedMethods

`[]string` · required

HTTP methods allowed for cross-origin requests.

Common values: "GET", "PUT", "POST", "DELETE", "HEAD"

For read-only access (e.g., serving images): ["GET", "HEAD"]
For upload access: ["GET", "PUT", "POST", "DELETE", "HEAD"]

- rule: {"repeated":{"minItems":"1","unique":true}}

### spec.corsRules[].allowedOrigins

`[]string` · required

Origins allowed to make cross-origin requests.

Format: scheme + domain (e.g., "https://example.com").
Use "*" to allow all origins (not recommended for production).

Examples:
  - "https://app.example.com" -- Allow specific origin
  - "https://*.example.com"   -- Allow subdomains
  - "*"                       -- Allow all (development only)

- rule: {"repeated":{"minItems":"1","unique":true}}

### spec.corsRules[].allowedHeaders

`[]string`

Headers that browsers are allowed to include in cross-origin
requests.

Common values: "Content-Type", "Authorization", "x-amz-*"
Use "*" to allow all headers.

Optional. If empty, only simple headers are allowed.

### spec.corsRules[].exposeHeaders

`[]string`

Headers that browsers are allowed to read from the response.

By default, browsers can only read "simple" response headers.
List headers here that your application needs to access.

Common values: "ETag", "x-amz-request-id"

Optional.

### spec.corsRules[].maxAgeSeconds

`int32`

Maximum time (in seconds) that the browser should cache the
preflight response.

Higher values reduce the number of preflight OPTIONS requests
at the cost of slower propagation of CORS policy changes.

Common value: 3600 (1 hour).
0 = browser uses its own default (usually 5 seconds).

### spec.forceDestroy

`bool`

Whether to force-destroy the bucket even if it contains objects.

When true, all objects (including locked objects and all versions)
are deleted before the bucket is destroyed.

WARNING: This is irreversible. All data in the bucket will be
permanently lost.

Recommended settings:
  - false for production buckets (prevents accidental data loss)
  - true for dev/test buckets (enables clean teardown)

Default: false

- default: `false`

## Validation Rules

- `object_lock_requires_versioning`: object_lock_enabled requires versioning_enabled to be true -- S3 Object Lock depends on versioning

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayObjectBucket, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.bucket_id` | `string` | The unique identifier of the created bucket. Format: "{region}/{bucket-name}" (e.g., "fr-par/my-app-media"). This is the primary output referenced by downstream resources. |
| `status.outputs.endpoint` | `string` | The FQDN endpoint URL of the bucket. Format: "{bucket-name}.s3.{region}.scw.cloud" This is the URL used by S3-compatible clients, CDNs, and applications to access objects in the bucket. |
| `status.outputs.api_endpoint` | `string` | The S3 API endpoint URL for the bucket's region. Format: "https://s3.{region}.scw.cloud" Used when configuring S3 clients that require the API endpoint separately from the bucket name (e.g., AWS CLI --endpoint-url). |
| `status.outputs.bucket_name` | `string` | The name of the bucket as it exists in Scaleway Object Storage. This is the bucket name that S3 clients use in their requests. Typically matches `metadata.name` but exported explicitly for downstream reference convenience. |
| `status.outputs.region` | `string` | The region where the bucket is deployed. Exported for downstream resources that need region-aware configuration (e.g., Lambda@Edge, CDN origin configuration). |

## See Also

- [Overview](../README.md)
