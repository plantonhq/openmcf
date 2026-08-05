# OpenStackApplicationCredential

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackApplicationCredentialSpec defines the configuration for an OpenStack
Identity (Keystone) application credential.

An application credential is a scoped authentication token that allows
applications to authenticate without using a user's password. Application
credentials are bound to the project that was in scope when they were created,
and they can optionally be restricted to specific roles and API access rules.

IMPORTANT: This is an immutable resource. All fields are ForceNew in the
Terraform provider -- any change to the spec requires destroying and
recreating the credential. The secret is generated once at creation time
and cannot be retrieved again after the initial API response.

The credential name is derived from metadata.name.

Terraform resource: openstack_identity_application_credential_v3
Pulumi resource: openstack.identity.ApplicationCredential

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackApplicationCredential
metadata:
  name: test-app-cred
spec:
  description: Test application credential for local development
  roles:
    - member
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.unrestricted` | `bool` |  | `false` |  |
| `spec.secret` | `string` (sensitive) |  |  |  |
| `spec.roles` | `[]string` |  |  |  |
| `spec.accessRules` | `[]AccessRule` |  |  |  |
| `spec.accessRules[].path` | `string` | yes |  |  |
| `spec.accessRules[].method` | `string` | yes |  |  |
| `spec.accessRules[].service` | `string` | yes |  |  |
| `spec.expiresAt` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.description

`string`

description is a human-readable description of the application credential.
ForceNew: cannot be changed after creation.

### spec.unrestricted

`bool` · optional (explicit presence)

unrestricted controls whether the application credential can create
additional application credentials or trusts.
Default: false (restricted). Setting to true is a security risk --
an unrestricted credential can be used to create sub-credentials.
ForceNew: cannot be changed after creation.

- default: `false`

### spec.secret

`string` · sensitive

secret is the application credential secret.
If omitted, OpenStack generates a random secret.
If provided, the user-specified secret is used.
ForceNew: cannot be changed after creation.
The secret is sensitive and should be stored securely.

### spec.roles

`[]string`

roles is a list of role names to scope the application credential.
The credential can only perform actions allowed by these roles.
If omitted, the credential inherits all roles of the creating user
on the current project.
ForceNew: cannot be changed after creation.

### spec.accessRules

`[]AccessRule`

access_rules restricts the credential to specific API operations.
Each rule specifies a service, HTTP method, and URL path pattern.
When set, the credential can ONLY call the specified APIs.
ForceNew: cannot be changed after creation.

### spec.accessRules[].path

`string` · required

path is the URL path pattern for the API endpoint.
Supports wildcards: "/v2.1/servers/*" matches all server sub-paths.
Required.

- rule: {"string":{"minLen":"1"}}

### spec.accessRules[].method

`string` · required

method is the HTTP method allowed for this rule.
Must be one of the standard HTTP methods.
Required.

- rule: {"string":{"minLen":"1","in":["POST","GET","HEAD","PATCH","PUT","DELETE"]}}

### spec.accessRules[].service

`string` · required

service is the OpenStack service type for this rule.
Examples: "identity", "compute", "block-storage", "image", "network"
Required.

- rule: {"string":{"minLen":"1"}}

### spec.expiresAt

`string`

expires_at is the expiration timestamp in RFC3339 format.
After this time, the credential becomes invalid.
If omitted, the credential does not expire.
ForceNew: cannot be changed after creation.
Example: "2027-01-01T00:00:00Z"

### spec.region

`string`

region overrides the region from the provider config.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackApplicationCredential, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | id is the unique identifier of the application credential in Keystone. |
| `status.outputs.name` | `string` | name is the name of the application credential (from metadata.name). |
| `status.outputs.secret` | `string` | secret is the application credential secret. SENSITIVE: This value is generated once at creation time. If a user-provided secret was specified in the spec, this echoes that value. Otherwise, it contains the auto-generated secret from OpenStack. This value is essential for downstream automation that uses this credential. |
| `status.outputs.project_id` | `string` | project_id is the UUID of the project this credential is scoped to. Computed from the authentication scope used during creation. |
| `status.outputs.region` | `string` | region is the OpenStack region where the credential was created. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
