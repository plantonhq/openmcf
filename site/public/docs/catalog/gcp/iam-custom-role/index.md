---
title: "IAM Custom Role"
description: "IAM Custom Role deployment documentation"
icon: "package"
order: 100
componentName: "gcpiamcustomrole"
---

# GCP IAM Custom Role

Creates a project-scoped IAM custom role — a named, least-privilege bundle of permissions defined once and granted wherever needed. Custom roles replace over-broad predefined roles (like `roles/storage.admin` for a workload that only writes objects) with exactly the permissions a workload uses.

## What Gets Created

When you deploy a GcpIamCustomRole resource, Planton provisions:

- **IAM Custom Role** — a `google_project_iam_custom_role` in the target project, grantable as `projects/<project>/roles/<roleId>`

The role is pure definition: it grants nothing until an IAM grant (a GcpProjectIamMember) binds it to an identity.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId`
- **IAM permissions** — `roles/iam.roleAdmin` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpIamCustomRole
metadata:
  name: log-bucket-writer
spec:
  projectId:
    value: my-gcp-project-123
  roleId: logBucketWriter
  title: Log Bucket Writer
  description: Grants exactly the permissions needed to write log objects
  permissions:
    - storage.objects.create
    - storage.objects.get
```

```shell
planton apply -f custom-role.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `roleId` | `string` | — | Required. Unique role ID (3-64 chars; letters, digits, underscores, periods; no hyphens). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the role. Can reference a GcpProject. Immutable. |
| `title` | `string` | — | Required. Console-visible title (max 100 chars). Mutable. |
| `description` | `string` | `""` | What the role is for and who should hold it (max 256 chars). Mutable. |
| `permissions` | `list(string)` | — | Required, min 1. Permission edits propagate immediately to every grant of the role. |
| `stage` | `string` | `GA` | Launch stage label: `ALPHA`, `BETA`, `GA`, `DEPRECATED`, `DISABLED`, `EAP`. `DISABLED` is an IAM kill switch. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Fully-qualified role name (`projects/<project>/roles/<roleId>`) — feed directly into a GcpProjectIamMember's `role` field |
| `role_id` | The bare role ID within the project |
| `deleted` | Whether the role is currently soft-deleted (GCP retains deleted roles for up to 14 days) |

## Related Components

- [GcpProjectIamMember](/docs/catalog/gcp/project-iam-member) — grants this role to an identity
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — the identity most commonly granted custom roles
- [GcpProject](/docs/catalog/gcp/project) — provides the project that owns the role
