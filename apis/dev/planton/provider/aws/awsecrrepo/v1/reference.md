# AwsEcrRepo

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsEcrRepoSpec defines the configuration for an Amazon ECR (Elastic Container
Registry) repository — the private Docker/OCI image registry that container
workloads (ECS, EKS, Lambda container images, App Runner) pull from.

The spec covers the full repository surface: tag mutability (including the
exclusion-filter modes that let a locked-down repository still float mutable
tags like "latest"), encryption at rest (AES256, KMS, or dual-layer
KMS_DSSE), push-time vulnerability scanning, structured lifecycle rules for
storage cost control, and a folded repository policy for cross-account or
cross-service pull access.

Notes:
- The repository name and the entire encryption configuration are create-time
  choices (ForceNew): changing either replaces the repository. Images are NOT
  migrated on replacement.
- Lifecycle rules and the repository policy are separate AWS API resources
  keyed 1:1 by the repository; they are folded here because they have no
  identity of their own and follow the repository's lifecycle.
- Registry-level ECR settings (registry policy, registry-wide scanning
  configuration, replication configuration, pull-through cache rules) are
  account-scoped, not repository-scoped, and are deliberately not part of
  this spec.

## Example

```yaml
# AWS ECR Repository — full-surface example
#
# Demonstrates the complete repository shape: exclusion-filtered tag
# immutability, KMS encryption, push-time scanning, structured lifecycle
# rules, and a folded repository policy.
#
# Usage:
#   planton apply -f manifest.yaml

apiVersion: aws.planton.dev/v1
kind: AwsEcrRepo
metadata:
  name: production-ecr-repo
spec:
  region: us-west-2

  # The immutable registry path (ForceNew). Slash namespaces are first-class.
  repositoryName: my-company/production-api

  # Release tags are frozen; floating tags matching the filters stay movable.
  imageTagMutability: IMMUTABLE_WITH_EXCLUSION
  imageTagMutabilityExclusionFilters:
    - latest
    - dev-*

  # Vulnerability scanning on every push (the spec default).
  scanOnPush: true

  # Customer-managed encryption (create-time; changing it replaces the repo).
  encryptionType: KMS
  kmsKeyId:
    value: arn:aws:kms:us-west-2:123456789012:key/00000000-0000-0000-0000-000000000000

  # Deletion fails while images exist (production safety default).
  forceDelete: false

  # Structured lifecycle rules — evaluated by ascending priority; the "any"
  # rule carries the highest priority (an AWS requirement).
  lifecycleRules:
    - rulePriority: 1
      description: Expire untagged images after 14 days
      tagStatus: untagged
      countType: sinceImagePushed
      countNumber: 14
    - rulePriority: 2
      description: Keep only the last 10 PR preview images
      tagStatus: tagged
      tagPrefixes:
        - pr-
      countType: imageCountMoreThan
      countNumber: 10
    - rulePriority: 3
      description: Keep at most 200 images overall
      tagStatus: any
      countType: imageCountMoreThan
      countNumber: 200

  # Folded repository policy: let AWS Lambda pull container images.
  repositoryPolicy:
    Version: "2012-10-17"
    Statement:
      - Sid: AllowLambdaPull
        Effect: Allow
        Principal:
          Service: lambda.amazonaws.com
        Action:
          - ecr:BatchGetImage
          - ecr:GetDownloadUrlForLayer

---
# Minimal example — secure defaults (MUTABLE tags, AES256 encryption,
# scanning on, deletion protected).

apiVersion: aws.planton.dev/v1
kind: AwsEcrRepo
metadata:
  name: minimal-ecr-repo
spec:
  region: us-west-2
  repositoryName: my-company/simple-service
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.repositoryName` | `string` | yes |  |  |
| `spec.imageTagMutability` | `string` |  | `MUTABLE` |  |
| `spec.imageTagMutabilityExclusionFilters` | `[]string` |  |  |  |
| `spec.encryptionType` | `string` |  | `AES256` |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.scanOnPush` | `bool` |  | `true` |  |
| `spec.forceDelete` | `bool` |  |  |  |
| `spec.lifecycleRules` | `[]AwsEcrRepoLifecycleRule` |  |  |  |
| `spec.lifecycleRules[].rulePriority` | `int32` | yes |  |  |
| `spec.lifecycleRules[].description` | `string` |  |  |  |
| `spec.lifecycleRules[].tagStatus` | `string` | yes |  |  |
| `spec.lifecycleRules[].tagPrefixes` | `[]string` |  |  |  |
| `spec.lifecycleRules[].tagPatterns` | `[]string` |  |  |  |
| `spec.lifecycleRules[].countType` | `string` | yes |  |  |
| `spec.lifecycleRules[].countNumber` | `int32` | yes |  |  |
| `spec.repositoryPolicy` | `object` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.repositoryName

`string` · required

Name of the ECR repository. Must be unique within the AWS account and
region, and cannot be changed after creation (ForceNew).

ECR repository names support slash-separated namespaces (e.g.
"team-blue/checkout-service"), which is why the name is an explicit spec
field rather than being derived from metadata.name: the namespace path is
registry structure, a genuinely different concept from the graph node's
display name.

Allowed characters: lowercase letters, numbers, and separators
(hyphen, underscore, period) between alphanumeric runs; path segments
separated by "/".

- rule: {"required":true,"string":{"minLen":"2","maxLen":"256","pattern":"^(?:[a-z0-9]+(?:[._-][a-z0-9]+)*/)*[a-z0-9]+(?:[._-][a-z0-9]+)*$"}}

### spec.imageTagMutability

`string` · optional (explicit presence)

Controls whether image tags can be overwritten after they are pushed.

Valid values:
- "MUTABLE" (AWS default): any tag can be overwritten by a later push.
  Convenient for development, but a moving "v1.2.3" undermines deploy
  reproducibility.
- "IMMUTABLE": no tag can ever be overwritten. The production
  recommendation — a tag permanently identifies one image digest.
- "IMMUTABLE_WITH_EXCLUSION": immutable EXCEPT for tags matching
  image_tag_mutability_exclusion_filters. The best of both: release tags
  are frozen while floating tags like "latest" or "dev-*" stay movable.
- "MUTABLE_WITH_EXCLUSION": mutable EXCEPT for tags matching the
  exclusion filters (the filters select the IMMUTABLE tags).

- default: `MUTABLE`
- rule: {"string":{"in":["MUTABLE","IMMUTABLE","IMMUTABLE_WITH_EXCLUSION","MUTABLE_WITH_EXCLUSION"]}}

### spec.imageTagMutabilityExclusionFilters

`[]string`

Wildcard tag filters that invert the repository's base mutability for
matching tags. Only allowed (and then required) when image_tag_mutability
is IMMUTABLE_WITH_EXCLUSION or MUTABLE_WITH_EXCLUSION.

Each filter is a tag pattern of letters, numbers, ".", "_", "-" and up to
two "*" wildcards. Examples: "latest", "dev-*", "*-snapshot".
Maximum 5 filters. AWS currently supports only wildcard-type filters, so
the filter type is implied and not modeled.

- rule: {"repeated":{"maxItems":"5","items":{"string":{"minLen":"1","maxLen":"128","pattern":"^[a-zA-Z0-9._-]*(\\*[a-zA-Z0-9._-]*){0,2}$"}}}}

### spec.encryptionType

`string` · optional (explicit presence)

How ECR encrypts stored images. Changing this after creation REPLACES the
repository (images are not migrated).

Valid values:
- "AES256" (default): server-side encryption with an Amazon S3-managed key.
  Zero configuration, no KMS cost.
- "KMS": server-side encryption with AWS KMS — either the AWS-managed
  "aws/ecr" key (when kms_key_id is omitted) or a customer-managed key.
  Required when compliance demands key rotation control or CloudTrail
  visibility of decrypt calls.
- "KMS_DSSE": dual-layer server-side encryption (two independent KMS
  envelope layers). For workloads subject to DoD CC SRG / top-secret
  data-at-rest requirements.

- default: `AES256`
- rule: {"string":{"in":["AES256","KMS","KMS_DSSE"]}}

### spec.kmsKeyId

`string | valueFrom`

Customer-managed KMS key used when encryption_type is "KMS" or "KMS_DSSE".
Accepts a key ARN or key ID, or a reference to an AwsKmsKey resource.
When omitted with KMS/KMS_DSSE, AWS uses the AWS-managed "aws/ecr" key.
Must not be set when encryption_type is "AES256".
Create-time only (ForceNew, together with encryption_type).

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.scanOnPush

`bool` · optional (explicit presence)

Enables automatic vulnerability scanning on every image push (the
repository-level "basic scanning" toggle). A production security
essential — shift-left detection of known CVEs before deploy.
Note: if the account has registry-level enhanced scanning configured
(Amazon Inspector), that account-wide setting supersedes this toggle.

- default: `true`

### spec.forceDelete

`bool`

When true, deleting the repository also deletes all contained images.
When false (default), deletion fails if any image is present — protection
against losing the only copy of a production image.

### spec.lifecycleRules

`[]AwsEcrRepoLifecycleRule`

Automated image-expiration rules for storage cost control. Active CI/CD
pipelines generate images far faster than anyone deletes them; without
lifecycle rules, storage grows unbounded.

Rules are evaluated by ascending rule_priority; an image is expired by the
first rule that selects it. AWS constraints (enforced here so they fail at
authoring time, not at apply): rule priorities must be unique, and at most
one rule may use tag_status "any" (AWS additionally requires that rule to
carry the highest priority — keep it last).

- rule: tag_status 'tagged' requires exactly one of tag_prefixes or tag_patterns; 'untagged' and 'any' take neither

### spec.lifecycleRules[].rulePriority

`int32` · required

Evaluation order: lower priorities are evaluated first, and an image is
expired by the first rule that selects it. Must be unique across rules.
A tag_status "any" rule must carry the highest priority in the policy.

- rule: {"required":true,"int32":{"gte":1}}

### spec.lifecycleRules[].description

`string`

Human-readable description of what the rule does, stored in the policy.
Example: "Expire untagged images after 14 days".

### spec.lifecycleRules[].tagStatus

`string` · required

Which images the rule selects by tag state:
- "tagged": only images with tags matching tag_prefixes or tag_patterns
  (exactly one of the two lists is required).
- "untagged": only images with no tags (typically superseded layers and
  failed builds). No prefix/pattern lists allowed.
- "any": every image in the repository. No prefix/pattern lists allowed;
  at most one such rule per policy and it must have the highest priority.

- rule: {"required":true,"string":{"in":["tagged","untagged","any"]}}

### spec.lifecycleRules[].tagPrefixes

`[]string`

Tag prefixes selecting tagged images (e.g. "release-" matches
"release-1.2.3"). Only with tag_status "tagged", and mutually exclusive
with tag_patterns. Maximum 100 prefixes.

- rule: {"repeated":{"maxItems":"100","items":{"string":{"minLen":"1","maxLen":"128"}}}}

### spec.lifecycleRules[].tagPatterns

`[]string`

Wildcard tag patterns selecting tagged images (e.g. "*-snapshot",
"v1.*"). Each pattern may contain up to four "*" wildcards. Only with
tag_status "tagged", and mutually exclusive with tag_prefixes.
Maximum 100 patterns.

- rule: {"repeated":{"maxItems":"100","items":{"string":{"minLen":"1","maxLen":"128","pattern":"^[a-zA-Z0-9._-]*(\\*[a-zA-Z0-9._-]*){0,4}$"}}}}

### spec.lifecycleRules[].countType

`string` · required

How the rule counts images for expiry:
- "imageCountMoreThan": keep the newest count_number selected images and
  expire the rest ("keep last N").
- "sinceImagePushed": expire selected images older than count_number days
  ("expire by age").

- rule: {"required":true,"string":{"in":["imageCountMoreThan","sinceImagePushed"]}}

### spec.lifecycleRules[].countNumber

`int32` · required

The count for count_type: an image count for "imageCountMoreThan", or a
number of days for "sinceImagePushed" (the "days" unit is the only unit
ECR supports and is implied).

- rule: {"required":true,"int32":{"gte":1}}

### spec.repositoryPolicy

`object`

Resource-based access policy for the repository, as a standard IAM policy
document. The primary uses are cross-account pulls (granting another
account's principals ecr:BatchGetImage / ecr:GetDownloadUrlForLayer) and
service access such as AWS Lambda pulling container images. AWS models
this as a separate API resource keyed by the repository; it is folded here
because the policy has no identity of its own and follows the repository's
lifecycle.

## Validation Rules

- `exclusion_filters_require_exclusion_mode`: image_tag_mutability_exclusion_filters can only be set when image_tag_mutability is 'IMMUTABLE_WITH_EXCLUSION' or 'MUTABLE_WITH_EXCLUSION'
- `exclusion_mode_requires_filters`: at least one image_tag_mutability_exclusion_filters entry is required when image_tag_mutability is 'IMMUTABLE_WITH_EXCLUSION' or 'MUTABLE_WITH_EXCLUSION'
- `kms_key_requires_kms_encryption`: kms_key_id can only be set when encryption_type is 'KMS' or 'KMS_DSSE' (AES256 uses an Amazon-managed key)
- `lifecycle_rule_priorities_unique`: each lifecycle_rules entry must have a unique rule_priority
- `lifecycle_single_any_rule`: at most one lifecycle_rules entry may use tag_status 'any' (AWS requires it to be the highest-priority rule)
- `lifecycle_single_untagged_rule`: at most one lifecycle_rules entry may use tag_status 'untagged'

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsEcrRepo, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.repository_name` | `string` | The repository name, matching spec.repository_name |
| `status.outputs.repository_url` | `string` | The repository URL, e.g. "123456789012.dkr.ecr.us-east-1.amazonaws.com/my-repo" |
| `status.outputs.repository_arn` | `string` | The repository ARN, e.g. "arn:aws:ecr:us-east-1:123456789012:repository/my-repo" |
| `status.outputs.registry_id` | `string` | The registry ID associated with this repository (i.e., the AWS Account ID) |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
