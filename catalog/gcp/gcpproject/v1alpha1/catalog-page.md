# GCP Project

Creates and configures a Google Cloud project within your resource
hierarchy — the Layer-0 container every other GCP resource lives in. The
component handles placement under an organization or folder, billing
linkage, labels and resource-manager tags, the default-network decision,
API enablement, and the deletion policy.

## What Gets Created

When you deploy a GcpProject resource, Planton provisions:

- **GCP Project** — a `google_project` resource placed under the specified
  organization or folder, with billing attached, labels merged (user labels
  beneath platform attribution labels), optional create-time tags, and the
  configured deletion policy
- **Enabled APIs** — a `google_project_service` resource for each entry in
  `enabledApis` (never disabled on destroy)

IAM grants are NOT created here — model each grant as a first-class
`GcpProjectIamMember` resource.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP organization or folder** — referenced via `parentType` and `parentId`
- **IAM permissions** to create projects under the target parent (`roles/resourcemanager.projectCreator`)
- **A billing account** (`roles/billing.user` on it) if the project will consume billable services

## Quick Start

Create a file `project.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpProject
metadata:
  name: my-project
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpProject.my-project
spec:
  projectId: my-dev-project-01
  parentType: organization
  parentId: "123456789012"
  billingAccountId: 0123AB-4567CD-89EFGH
```

Deploy:

```shell
planton apply -f project.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `projectId` | `string` | Globally-unique project ID. IMMUTABLE; deleted IDs stay reserved ~30 days. | Required, 6-30 chars, lowercase letters/digits/hyphens, starts with a letter, no trailing hyphen |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `displayName` | `string` | `metadata.name` | Human-readable console name (mutable). |
| `parentType` / `parentId` | `enum` / `string` | — | `organization` or `folder` plus the numeric parent ID. Changing migrates the project. |
| `billingAccountId` | `string` | — | `XXXXXX-XXXXXX-XXXXXX` billing account to link. |
| `labels` | `map<string,string>` | — | User labels (the primary cost-allocation dimension); platform attribution labels win on key conflicts. |
| `tags` | `map<string,string>` | — | Resource-manager tags (`tagKeys/{id}` → `tagValues/{id}`), CREATE-TIME only — changing later recreates the project. |
| `autoCreateNetwork` | `bool` | `false` | Whether GCP auto-creates the "default" VPC. Kept false as a hardening default; one network slot of quota is still needed momentarily. |
| `enabledApis` | `string[]` | — | Cloud APIs to pre-enable (each ends with `.googleapis.com`). |
| `deletionPolicy` | `string` | `DELETE` | `DELETE` (shutdown with 30-day recovery), `PREVENT` (destroy fails), or `ABANDON` (remove from state, project lives on). |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `project_id` | `string` | The unique project ID — the value every other kind's `projectId` reference resolves |
| `project_number` | `string` | The numeric identifier assigned by Google |
| `name` | `string` | The display name |

## Related Components

- [GcpProjectIamMember](/docs/catalog/gcp/project-iam-member) — additive IAM grants on the project
- [GcpVpcNetwork](/docs/catalog/gcp/vpc-network) — explicit networks instead of the default network
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — workload identities inside the project
