---
title: "Production Immutable ECR Repository"
description: "This preset creates an ECR repository where release tags are frozen forever while the floating `latest` tag stays movable, with automatic vulnerability scanning and lifecycle rules that balance cost..."
type: "preset"
rank: "01"
presetSlug: "01-production-immutable"
componentSlug: "ecr-repo"
componentTitle: "ECR Repo"
provider: "aws"
icon: "package"
order: 1
---

# Production Immutable ECR Repository

This preset creates an ECR repository where release tags are frozen forever while the floating `latest` tag stays movable, with automatic vulnerability scanning and lifecycle rules that balance cost control with rollback capability. Exclusion-filtered immutability gives production the best of both worlds: `v1.2.3` can never be overwritten by a CI/CD pipeline, and `latest` keeps working as a convenience pointer.

## When to Use

- Production container registries where image tag integrity is critical
- Regulated environments requiring immutable artifacts for audit trails
- Any CI/CD pipeline pushing Docker images to ECR for production workloads

## Key Configuration Choices

- **Exclusion-filtered immutability** (`imageTagMutability: IMMUTABLE_WITH_EXCLUSION` + `latest` filter) -- Every tag is frozen once pushed, except tags matching the filters; guarantees `v1.2.3` always refers to the same image while `latest` stays movable
- **Scan on push** (`scanOnPush: true`) -- Automatic vulnerability scanning when images are pushed
- **AES256 encryption** (`encryptionType: AES256`) -- AWS-managed server-side encryption at rest (default, no additional cost)
- **Untagged cleanup rule** (priority 1, `sinceImagePushed` 7 days) -- Removes orphaned layers and failed builds quickly to control costs
- **Keep-last-100 rule** (priority 2, `imageCountMoreThan` 100 on `any`) -- Keeps recent images available for rollback; the `any` rule carries the highest priority, as AWS requires
- **Force delete disabled** (`forceDelete: false`) -- Repository cannot be deleted while it contains images

## Placeholders to Replace

The example `repositoryName` (`team-blue/checkout-service`) is a realistic
registry path — replace it with your own naming convention.

| Placeholder | Description | Where to Find |
| --- | --- | --- |
