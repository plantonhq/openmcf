# OpenStack Application Credential

Deploys a Keystone application credential on OpenStack -- a scoped authentication token that allows applications to authenticate without a user's password. The credential supports role restrictions, fine-grained API access rules, and optional expiration, with the secret generated once at creation time and captured in outputs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Application Credential** -- a Keystone identity credential bound to the creating user's project, with configurable roles, API access rules, unrestricted mode, and optional expiration

## Before You Deploy

### OpenStack Account

- **Project scope** -- the credential is created in the project associated with the current authentication scope. The creating user must have valid Keystone credentials on the target project.
- **Roles** (optional) -- if restricting the credential to specific roles, confirm the role names exist in your Keystone deployment. Run `openstack role list` to find available roles. If omitted, the credential inherits all roles of the creating user.
- **Access rules** (optional) -- if restricting the credential to specific API operations, confirm the service types (e.g., `compute`, `identity`, `block-storage`) match the services deployed in your OpenStack environment.
- **Immutability** -- all fields are ForceNew. Any change to the spec destroys and recreates the credential, generating a new secret. Downstream consumers must update their authentication configuration.

## Deploy

### Console

Open the deployment store, find **OpenStack Application Credential**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Restricted Read-Only** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackApplicationCredential
metadata:
  name: ci-readonly
  org: acme-corp
  env: prod
spec:
  roles:
    - reader
```

```shell
planton apply -f application-credential.yaml
```

This creates a restricted credential with the `reader` role, scoped to the current project, with an auto-generated secret and no expiration. No access rules, custom secret, or expiration are configured.

## Key Configuration

These are the most important decisions when configuring an application credential. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Role scoping** -- `roles` restricts the credential to specific Keystone roles. If omitted, the credential inherits all roles of the creating user on the current project. For least-privilege access, explicitly list only the roles the application needs.

**Access rules** -- `accessRules` further restricts the credential to specific API operations by service type, HTTP method, and URL path pattern. Each rule specifies a service (e.g., `compute`), a method (e.g., `GET`), and a path (e.g., `/v2.1/servers/*`). Combine with role scoping for defense-in-depth.

**Unrestricted mode** -- `unrestricted` defaults to `false`. When `true`, the credential can create additional application credentials or trusts. Leave `false` for automated workloads unless the credential explicitly needs to delegate access.

**Secret management** -- `secret` accepts a user-provided secret string. If omitted, OpenStack generates a random secret. The secret appears in `status.outputs.secret` exactly once at creation time. Store it in a secret manager immediately -- it cannot be retrieved again from the OpenStack API.

**Expiration** -- `expiresAt` sets a hard expiration in RFC3339 format (e.g., `2027-01-01T00:00:00Z`). For long-running services, omit this field for a non-expiring credential. For CI/CD or temporary access, set a bounded expiration.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | UUID of the application credential in Keystone | Automation reference, audit trails |
| `name` | Name of the credential | Monitoring labels, configuration references |
| `secret` | Application credential secret (sensitive) | Authentication in downstream services and CI/CD pipelines |
| `project_id` | UUID of the project this credential is scoped to | Verifying scope alignment with other resources |
| `region` | OpenStack region where the credential was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Read-only monitoring credential** -- A restricted credential with the `reader` role and GET-only access rules for compute, network, and block-storage APIs. Suitable for monitoring agents, dashboards, and audit tools that should never modify infrastructure. Start from the **Restricted Read-Only** preset.

**Compute automation credential** -- A credential scoped to compute (Nova) operations with the `member` role. Allows creating, managing, and deleting instances without access to networking, storage, or identity APIs. Start from the **Compute-Scoped** preset.

## Works With

This component operates independently and does not reference other deployment components.