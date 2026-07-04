---
title: "Project IAM Member"
description: "Project IAM Member deployment documentation"
icon: "package"
order: 100
componentName: "gcpprojectiammember"
---

# GCP Project IAM Member

Grants one role to one identity on one GCP project — the safe, additive unit of GCP access control. The grant merges into the project's IAM policy without touching any other member's bindings, and removal subtracts only this exact pair, so grants from different charts, teams, and tools never fight.

## What Gets Created

When you deploy a GcpProjectIamMember resource, Planton provisions:

- **Project IAM Member** — one (role, member[, condition]) entry merged into the target project's IAM policy

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId`
- **IAM permissions** — `roles/resourcemanager.projectIamAdmin` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpProjectIamMember
metadata:
  name: app-logs-writer-grant
spec:
  projectId:
    value: my-gcp-project-123
  role:
    value: roles/logging.logWriter
  member:
    value: serviceAccount:my-app@my-gcp-project-123.iam.gserviceaccount.com
```

```shell
planton apply -f iam-grant.yaml
```

Compose with first-class nodes instead of literals: reference a GcpIamCustomRole's `name` output for the role and a GcpServiceAccount's `member` output for the member.

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | Project whose IAM policy receives the grant. Can reference a GcpProject. Immutable. |
| `role` | `StringValueOrRef` | — | Required. Predefined role (`roles/...`) or custom role name (`projects/.../roles/...`). References GcpIamCustomRole by default. Immutable. |
| `member` | `StringValueOrRef` | — | Required. Identity in IAM member format (`serviceAccount:`, `user:`, `group:`, `domain:`, `principal://`, `allUsers`, ...). References GcpServiceAccount by default. Immutable. |
| `condition` | `object` | — | Optional IAM Condition (`title`, `expression`, optional `description`). Part of the grant's identity. Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `project_id` | The project whose policy received the grant (after reference resolution) |
| `role` | The granted role (after reference resolution) |
| `member` | The granted member (after reference resolution) |
| `etag` | The project IAM policy etag after the grant — useful for audit correlation |

## Related Components

- [GcpIamCustomRole](/docs/catalog/gcp/iam-custom-role) — defines the custom role this grant can reference
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — the identity most commonly granted
- [GcpProject](/docs/catalog/gcp/project) — provides the project whose policy receives the grant
