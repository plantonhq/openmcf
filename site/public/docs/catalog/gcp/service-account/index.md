---
title: "Service Account"
description: "Service Account deployment documentation"
icon: "package"
order: 100
componentName: "gcpserviceaccount"
---

# GCP Service Account

Deploys a Google Cloud service account — the identity workloads authenticate as — with optional JSON key generation and additive IAM role grants at the project and organization levels. The module creates the account, attaches the specified roles using per-role additive `IAMMember` grants, and exports the identity handles downstream resources reference (email, the ready-made IAM member string, unique ID, resource name, and the key when requested).

## What Gets Created

When you deploy a GcpServiceAccount resource, Planton provisions:

- **Service Account** — a GCP service account in the specified project, with `serviceAccountId` as the account ID and `displayName` (falling back to `metadata.name`) as the console-visible name
- **Service Account Key** (conditional) — a JSON private key for the service account, created only when `createKey` is set to `true`
- **Project IAM Grants** — one additive `IAMMember` per entry in `projectIamRoles`, granting each role to the service account in the target project without touching other members' bindings
- **Organization IAM Grants** — one additive `IAMMember` per entry in `orgIamRoles`, granting each role in the specified organization (requires `orgId`)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** where the service account will be created
- **Organization ID** if you need to assign organization-level IAM roles
- **IAM permissions** to create service accounts and manage IAM bindings in the target project (and organization, if applicable)

## Quick Start

Create a file `service-account.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceAccount
metadata:
  name: my-app-sa
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpServiceAccount.my-app-sa
spec:
  serviceAccountId: my-app-sa
  projectId:
    value: my-gcp-project-123
```

Deploy:

```shell
planton apply -f service-account.yaml
```

This creates a service account `my-app-sa@my-gcp-project-123.iam.gserviceaccount.com` with no keys and no additional IAM roles.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `serviceAccountId` | `string` | Short unique ID for the service account, used to form the email `<serviceAccountId>@<project>.iam.gserviceaccount.com`. Immutable. | Required; 6-30 chars; lowercase letters, digits, hyphens; starts with a letter; cannot end with a hyphen |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | Provider default project | The GCP project in which the service account is created. Can be a literal value or a reference to a GcpProject resource. Immutable. |
| `displayName` | `string` | metadata name | Human-readable identity shown in the GCP console. Mutable. |
| `description` | `string` | `""` | What this identity is for (max 256 bytes). Mutable. |
| `disabled` | `bool` | `false` | Disabled accounts keep their IAM bindings but cannot authenticate — a kill switch for incident response or staged decommissioning. Mutable. |
| `createKey` | `bool` | `false` | When `true`, a JSON private key is generated for the service account and its base64-encoded value is exported in stack outputs. |
| `projectIamRoles` | `string[]` | `[]` | IAM roles granted at the project level (additive member grants; e.g., `roles/logging.logWriter`). For custom roles, conditions, or independently-owned grants, use GcpProjectIamMember instead. |
| `orgId` | `string` | `""` | GCP organization ID (numeric string). Required when `orgIamRoles` is non-empty. |
| `orgIamRoles` | `string[]` | `[]` | IAM roles granted at the organization level (additive; e.g., `roles/resourcemanager.organizationViewer`). Requires `orgId`. |

## Examples

### Service Account with Project IAM Roles

A service account with permissions to write logs and read Cloud Storage buckets:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceAccount
metadata:
  name: backend-worker
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpServiceAccount.backend-worker
spec:
  serviceAccountId: backend-worker
  projectId:
    value: my-gcp-project-123
  projectIamRoles:
    - roles/logging.logWriter
    - roles/storage.objectViewer
```

### Service Account with Key Generation

A CI/CD service account with a generated JSON key and deployment permissions:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceAccount
metadata:
  name: ci-deployer
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: staging.GcpServiceAccount.ci-deployer
spec:
  serviceAccountId: ci-deployer
  projectId:
    value: my-gcp-project-123
  createKey: true
  projectIamRoles:
    - roles/container.developer
    - roles/storage.admin
    - roles/artifactregistry.writer
```

### Service Account with Organization IAM Roles

A service account that needs both project-level and organization-level permissions:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceAccount
metadata:
  name: org-auditor
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpServiceAccount.org-auditor
spec:
  serviceAccountId: org-auditor
  projectId:
    value: my-gcp-project-123
  orgId: "123456789012"
  projectIamRoles:
    - roles/logging.logWriter
  orgIamRoles:
    - roles/resourcemanager.organizationViewer
    - roles/iam.securityReviewer
```

### Using a Foreign Key Reference for Project ID

Reference an Planton-managed GcpProject instead of hardcoding the project ID:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceAccount
metadata:
  name: app-runtime
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpServiceAccount.app-runtime
spec:
  serviceAccountId: app-runtime
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  createKey: true
  projectIamRoles:
    - roles/cloudrun.invoker
    - roles/secretmanager.secretAccessor
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `email` | `string` | The full email address of the created service account (`<serviceAccountId>@<project>.iam.gserviceaccount.com`) |
| `member` | `string` | Ready-made IAM member string (`serviceAccount:<email>`) — feed directly into IAM grants |
| `unique_id` | `string` | Stable numeric ID, never reused across delete/recreate |
| `name` | `string` | Fully-qualified resource name (`projects/<project>/serviceAccounts/<email>`) |
| `key_base64` | `string` | Base64-encoded JSON private key. Only populated when `createKey` is `true`. Sensitive. |

## Related Components

- [GcpProjectIamMember](/docs/catalog/gcp/project-iam-member) — first-class additive grant referencing this account's `member` output
- [GcpIamCustomRole](/docs/catalog/gcp/iam-custom-role) — least-privilege role definitions to grant to this account
- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project where the service account is created
- [GcpGkeCluster](/docs/catalog/gcp/gke-cluster) — GKE clusters that may use this service account for workload identity
- [GcpCloudRun](/docs/catalog/gcp/cloud-run) — Cloud Run services that may run under this service account
