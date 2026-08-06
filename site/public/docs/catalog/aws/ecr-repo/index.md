---
title: "ECR Repo"
description: "ECR Repo deployment documentation"
icon: "package"
order: 100
componentName: "awsecrrepo"
---

# AWS ECR Repo

Deploys an AWS Elastic Container Registry repository at the full repository surface: tag mutability (including exclusion-filtered modes that freeze release tags while `latest` stays movable), AES256/KMS/dual-layer encryption, push-time vulnerability scanning, structured lifecycle rules, and a folded repository policy for cross-account or service access.

## What Gets Created

When you deploy an AwsEcrRepo resource, Planton provisions:

- **ECR Repository** — the repository with the specified name (a slash-namespaced registry path), tag mutability mode plus wildcard exclusion filters, image scanning configuration, encryption configuration (create-time), force-delete behavior, and Planton resource tags
- **Lifecycle Policy** (optional) — the repository's image-expiration rules, built from the structured `lifecycleRules` list (evaluated by ascending priority; an image is expired by the first rule that selects it)
- **Repository Policy** (optional) — the repository's resource-based access policy, from the `repositoryPolicy` document

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A KMS key** if using `KMS` or `KMS_DSSE` encryption (optional; the default `AES256` encryption requires no additional setup)

## Quick Start

Create a file `ecr-repo.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcrRepo
metadata:
  name: my-service
spec:
  region: us-east-1
  repositoryName: my-org/my-service
```

Deploy:

```shell
planton apply -f ecr-repo.yaml
```

This creates an ECR repository with mutable tags, AES256 encryption, and scan-on-push enabled (the defaults).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region where the ECR repository will be created. | Valid AWS region |
| `repositoryName` | `string` | Name of the ECR repository — a slash-namespaced registry path unique within the account and region, e.g., `team-blue/my-microservice`. Create-time immutable. | 2–256 characters; lowercase letters, numbers, `._-` separators, `/` namespaces |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `imageTagMutability` | `string` | `"MUTABLE"` | Tag overwrite behavior: `MUTABLE`, `IMMUTABLE`, `IMMUTABLE_WITH_EXCLUSION` (immutable except tags matching the filters), or `MUTABLE_WITH_EXCLUSION` (mutable except tags matching the filters). |
| `imageTagMutabilityExclusionFilters` | `list(string)` | — | Wildcard tag patterns (max 5, up to two `*` each) that invert the base mutability. Required with, and only valid with, the `*_WITH_EXCLUSION` modes. |
| `encryptionType` | `string` | `"AES256"` | Encryption at rest: `AES256` (AWS-managed), `KMS` (customer-managed key), or `KMS_DSSE` (dual-layer). Create-time immutable — changing it replaces the repository. |
| `kmsKeyId` | `string \| ref` | — | KMS key ARN/ID (or AwsKmsKey reference) for `KMS`/`KMS_DSSE`. Omit to use the AWS-managed `aws/ecr` key. Not valid with `AES256`. |
| `scanOnPush` | `bool` | `true` | Automatic vulnerability scanning when images are pushed. |
| `forceDelete` | `bool` | `false` | When `true`, deleting the repository removes all contained images. When `false`, deletion fails if images exist. |
| `lifecycleRules` | `list(object)` | — | Structured image-expiration rules (see below). If omitted, no lifecycle policy is created. |
| `repositoryPolicy` | `object` | — | Resource-based IAM policy document (cross-account pulls, service principals such as Lambda). |

### Lifecycle Rule Fields

| Field | Type | Description |
|-------|------|-------------|
| `rulePriority` | `int32` | Evaluation order (unique; lower first). A `tagStatus: any` rule must carry the highest priority. |
| `description` | `string` | Human-readable note stored in the policy. |
| `tagStatus` | `string` | Which images the rule selects: `tagged`, `untagged`, or `any`. At most one `untagged` and one `any` rule per policy. |
| `tagPrefixes` | `list(string)` | Tag prefixes selecting tagged images. Exactly one of `tagPrefixes`/`tagPatterns` with `tagged`. |
| `tagPatterns` | `list(string)` | Wildcard tag patterns (up to four `*` each) selecting tagged images. |
| `countType` | `string` | `imageCountMoreThan` (keep last N) or `sinceImagePushed` (expire by age in days). |
| `countNumber` | `int32` | The image count or the number of days, per `countType`. |

## Examples

### Production: frozen release tags, movable `latest`

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcrRepo
metadata:
  name: prod-api
spec:
  region: us-east-1
  repositoryName: my-org/api-server
  imageTagMutability: IMMUTABLE_WITH_EXCLUSION
  imageTagMutabilityExclusionFilters:
    - latest
```

### KMS Encryption

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcrRepo
metadata:
  name: compliant-repo
spec:
  region: us-east-1
  repositoryName: my-org/compliant-service
  imageTagMutability: IMMUTABLE
  encryptionType: KMS
  kmsKeyId:
    value: arn:aws:kms:us-east-1:123456789012:key/abcd-1234-efgh-5678
```

### Lifecycle Rules for Cost Control

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcrRepo
metadata:
  name: ci-images
spec:
  region: us-east-1
  repositoryName: my-org/ci-runner
  lifecycleRules:
    - rulePriority: 1
      description: Expire untagged images after 7 days
      tagStatus: untagged
      countType: sinceImagePushed
      countNumber: 7
    - rulePriority: 2
      description: Keep only the last 10 PR preview images
      tagStatus: tagged
      tagPrefixes:
        - pr-
      countType: imageCountMoreThan
      countNumber: 10
    - rulePriority: 3
      description: Keep at most 50 images overall
      tagStatus: any
      countType: imageCountMoreThan
      countNumber: 50
```

### Cross-Service Pull Access

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcrRepo
metadata:
  name: lambda-images
spec:
  region: us-east-1
  repositoryName: my-org/lambda-functions
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
```

### Development Repository with Force Delete

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsEcrRepo
metadata:
  name: dev-scratch
spec:
  region: us-west-2
  repositoryName: my-org/scratch
  forceDelete: true
  lifecycleRules:
    - rulePriority: 1
      tagStatus: untagged
      countType: sinceImagePushed
      countNumber: 1
    - rulePriority: 2
      tagStatus: any
      countType: imageCountMoreThan
      countNumber: 10
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `repository_name` | `string` | The repository name, matching `spec.repositoryName` |
| `repository_url` | `string` | The repository URL for docker push/pull (e.g., `123456789012.dkr.ecr.us-east-1.amazonaws.com/my-repo`) |
| `repository_arn` | `string` | The repository ARN (e.g., `arn:aws:ecr:us-east-1:123456789012:repository/my-repo`) |
| `registry_id` | `string` | The registry ID (AWS account ID) associated with the repository |

## Related Components

- [AwsKmsKey](/docs/catalog/aws/kms-key) — provides a customer-managed KMS key for repository encryption
- [AwsIamRole](/docs/catalog/aws/iam-role) — grants push/pull permissions to CI/CD pipelines or services
- [AwsEksCluster](/docs/catalog/aws/eks-cluster) — Kubernetes cluster that pulls images from ECR repositories
