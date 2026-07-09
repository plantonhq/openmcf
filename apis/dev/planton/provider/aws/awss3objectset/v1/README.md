# Overview

The AWS S3 Object Set API resource provides a declarative interface for uploading and managing one or more objects in an Amazon S3 bucket. By abstracting the complexity of individual S3 object uploads, this resource allows you to manage configuration files, static assets, seed data, and other content alongside your infrastructure.

## Why We Created This API Resource

Managing S3 objects alongside infrastructure requires coordinating bucket creation with object uploads. This resource solves that by:

- **Declarative Object Management**: Define objects as code alongside your infrastructure, ensuring they are created, updated, or removed consistently.
- **Foreign Key Integration**: Reference an `AwsS3Bucket` component directly, so the bucket name is resolved automatically when both resources are managed together.
- **Batch Uploads**: Upload multiple objects to the same bucket in a single deployment, reducing configuration duplication.
- **Content Flexibility**: Support both inline text content (UTF-8) and base64-encoded binary content for any file type.

## Key Features

### Bucket Reference via Foreign Key

- **Literal Bucket Name**: Provide the bucket name directly as a string value.
- **Component Reference**: Reference an `AwsS3Bucket` component by name; the bucket ID is resolved from `status.outputs.bucket_id` automatically.

### Multi-Object Support

- Upload one or more objects per deployment. Each object's identity is its key: adding, removing, or reordering entries never churns unrelated objects.
- Set-level tags are merged with object-level tags (object tags take precedence).

### Content Sources

- **Inline Text** (`content`): For configuration files, JSON, YAML, HTML, and other text formats.
- **Base64 Binary** (`content_base64`): For images, compiled assets, or any binary data small enough to carry in a manifest.

### Per-Object HTTP and Storage Surface

- **Presentation headers**: content type, cache control, content encoding, content disposition, content language, and website redirect targets.
- **User metadata**: lowercase-keyed `x-amz-meta-*` entries stored with the object.
- **Storage class**: STANDARD through INTELLIGENT_TIERING, infrequent-access, and archive tiers.
- **Upload integrity checksums**: CRC32/CRC32C/CRC64NVME/SHA1/SHA256 stored alongside the object.

### Security Posture: Bucket First, Object Override Second

Uniform encryption and access posture belongs on the bucket — `AwsS3Bucket` models default encryption, versioning, public-access blocking, and ownership controls, and S3 applies them to every uploaded object automatically. The per-object `server_side_encryption` / `kms_key` / `bucket_key_enabled` / `acl` fields here are overrides for individual objects that must diverge. Object Lock retention (`object_lock_mode` + `object_lock_retain_until_date`, plus legal holds) is available for objects in Object Lock-enabled buckets.

### Clean Destroys, Even Under Retention

Destroying an object removes all of its versions on versioned buckets. For objects under GOVERNANCE-mode Object Lock retention or legal holds, `force_destroy` sends the governance-bypass flag with the delete (valid only on Object Lock-enabled buckets) so retained objects can still be torn down deliberately.

## Stack Outputs

- **bucket_id**: The bucket the objects were uploaded to.
- **object_arns**: Map of object key to ARN, for IAM policy Resource lists.
- **object_etags**: Map of object key to ETag (content hash) for cache invalidation.
- **object_version_ids**: Map of object key to version ID (when bucket versioning is enabled).

## Benefits

- **Infrastructure as Code**: Manage S3 objects declaratively alongside buckets and other resources.
- **Consistency**: Ensure objects are always in sync with infrastructure deployments.
- **Simplicity**: Single resource manages multiple objects with shared tag inheritance.
- **Flexibility**: Text and binary content, the full per-object HTTP/storage/encryption surface, and honest composition with the bucket's own posture.
