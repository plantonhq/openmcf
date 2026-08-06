---
title: "Artifact Registry Repository"
description: "Artifact Registry Repository deployment documentation"
icon: "package"
order: 100
componentName: "gcpartifactregistryrepo"
---

# GCP Artifact Registry Repository

Creates a Google Cloud Artifact Registry repository — the package store for container images (Docker/OCI), language packages (Maven, npm, Python, Go), and OS packages (Apt, Yum) — in any of its three serving modes: standard (push target), remote (pull-through cache of an upstream registry), or virtual (priority-ordered aggregation of other repositories). Access is granted additively per repository.

## What Gets Created

- The Artifact Registry API is enabled on the project (never disabled on destroy)
- A `google_artifact_registry_repository` carrying your labels merged beneath Planton's attribution labels (`planton-ai_resource`, `planton-ai_name`, `planton-ai_kind`, plus org/env/id when set)
- One `google_artifact_registry_repository_iam_member` per `iamMembers` entry — additive grants that compose safely with grants made elsewhere

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **IAM permissions** — `roles/artifactregistry.admin` on the target project
- For CMEK: a `GcpKmsKey` the Artifact Registry service agent can use

## Quick Start

Create a file `artifact-repo.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpArtifactRegistryRepo
metadata:
  name: app-images
spec:
  location: us-central1
  format: DOCKER
  dockerConfig:
    immutableTags: true
  cleanupPolicies:
    - id: delete-old-untagged
      action: DELETE
      condition:
        olderThan: "2592000s"
        tagState: UNTAGGED
    - id: keep-last-10
      action: KEEP
      mostRecentVersions:
        keepCount: 10
```

Deploy:

```shell
planton apply -f artifact-repo.yaml
```

This creates a private Docker repository with immutable tags whose storage stays bounded: superseded digests are deleted after 30 days while the 10 most recent versions of every image are always kept.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Target project; empty uses the provider default. Immutable |
| `repositoryId` | string | No | URL path segment; defaults to `metadata.name`. Immutable |
| `location` | string | Yes | Region or multi-region (`us`, `europe`, `asia`). Immutable |
| `format` | string | Yes | `DOCKER`, `MAVEN`, `NPM`, `PYTHON`, `GO`, `APT`, `YUM`, `GENERIC`, `KFP`, ... Immutable |
| `mode` | string | No | `STANDARD_REPOSITORY` (default), `REMOTE_REPOSITORY`, `VIRTUAL_REPOSITORY`. Immutable |
| `description` | string | No | Human-readable description |
| `labels` | map | No | User labels (merged beneath platform labels) |
| `kmsKeyName` | StringValueOrRef | No | CMEK key (reference a `GcpKmsKey`). Immutable |
| `dockerConfig.immutableTags` | bool | No | Pushed tags permanently identify one digest |
| `mavenConfig.versionPolicy` | string | No | `RELEASE` or `SNAPSHOT` |
| `mavenConfig.allowSnapshotOverwrites` | bool | No | Allow re-publishing an existing version |
| `cleanupPolicies[]` | list | No | `id` + `action` (`DELETE`/`KEEP`) + `condition`/`mostRecentVersions` |
| `cleanupPolicyDryRun` | bool | No | Log cleanup matches without deleting |
| `remoteRepositoryConfig` | object | No | Exactly one upstream arm (public registries, Apt/Yum mirrors, or `commonRepository.uri`); optional Secret-Manager-referenced credentials |
| `virtualRepositoryConfig.upstreamPolicies[]` | list | No | `id` + `repository` (reference another repo's `repository_path`) + `priority` |
| `vulnerabilityScanningEnablement` | string | No | `INHERITED` (default) or `DISABLED` |
| `iamMembers[]` | list | No | Additive grants: `role` + `member` (+ optional IAM `condition`) |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Short repository name |
| `repository_path` | `projects/{p}/locations/{l}/repositories/{r}` — what composing resources reference |
| `registry_uri` | Push/pull endpoint (`us-central1-docker.pkg.dev/{project}/{repo}`) |
| `location` | Repository location |

## Related Resources

- **GcpCloudFunction** — `dockerRepository` references this repository's `repository_path`
- **GcpKmsKey** — CMEK protection via `kmsKeyName`
- **GcpServiceAccount** — pull/push identities granted via `iamMembers`
- **GcpGcsBucket** — the object-storage sibling for non-package artifacts
