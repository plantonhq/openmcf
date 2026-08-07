# GcpServiceAccountIamMember

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `gcp.planton.dev/v1alpha1`

GcpServiceAccountIamMemberSpec defines a single ADDITIVE IAM grant ON a
service account: one role, to one member, on one service account resource.

A GCP service account is both an identity (it can be granted roles on other
resources) and a resource (other identities can be granted roles ON it).
This kind covers the resource side — the grants that control who may USE or
MANAGE the service account itself:
  roles/iam.workloadIdentityUser      — lets a federated principal (GitHub
                                        Actions, GKE workload, external OIDC
                                        identity) impersonate the account;
                                        the keyless-authentication hop.
  roles/iam.serviceAccountTokenCreator — lets the member mint short-lived
                                        access/ID tokens as the account.
  roles/iam.serviceAccountUser        — lets the member deploy or run
                                        workloads AS the account (actAs);
                                        required by Cloud Run, GCE, Cloud
                                        Functions deployments that attach it.
Prefer these account-scoped grants over their project-level equivalents:
a project-level serviceAccountUser grant allows acting as EVERY service
account in the project.

Additive semantics: this grant merges into the service account's IAM policy
without touching any other member's bindings on the same role, and removal
subtracts only this exact (role, member) pair. Multiple grants — from
different charts, teams, or tools — never fight over the policy.
(Authoritative binding/policy management, which clobbers everything not
listed, is deliberately not modeled.)

Every field is immutable, mirroring the underlying API: an IAM grant has no
update — changing the account, role, member, or condition replaces the
grant atomically.

## Example

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpServiceAccountIamMember
metadata:
  name: my-sample-sa-grant
spec:
  # The service account whose IAM policy receives this grant, as its
  # fully-qualified resource name. Reference a GcpServiceAccount — its `name`
  # output is exactly this value.
  serviceAccountId:
    value: projects/my-gcp-project-123/serviceAccounts/deployer@my-gcp-project-123.iam.gserviceaccount.com

  # The role to grant ON the account: typically an impersonation/usage role.
  role:
    value: roles/iam.serviceAccountTokenCreator

  # The identity receiving the grant, in GCP IAM member format.
  # Reference a GcpServiceAccount — its `member` output is exactly this value.
  member:
    value: serviceAccount:caller@my-gcp-project-123.iam.gserviceaccount.com

  # Optional IAM Condition restricting when the grant applies
  # condition:
  #   title: expires-2027
  #   expression: request.time < timestamp("2027-01-01T00:00:00Z")
  #   description: Temporary impersonation for the migration
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.serviceAccountId` | `string \| valueFrom` | yes |  | GcpServiceAccount (`status.outputs.name`) |
| `spec.role` | `string \| valueFrom` | yes |  | GcpIamCustomRole (`status.outputs.name`) |
| `spec.member` | `string \| valueFrom` | yes |  | GcpServiceAccount (`status.outputs.member`) |
| `spec.condition` | `GcpServiceAccountIamMemberCondition` |  |  |  |
| `spec.condition.title` | `string` | yes |  |  |
| `spec.condition.expression` | `string` | yes |  |  |
| `spec.condition.description` | `string` |  |  |  |

## Field Details

### spec.serviceAccountId

`string | valueFrom` · required

The service account whose IAM policy receives this grant, as its
fully-qualified resource name:
  projects/<project>/serviceAccounts/<email>
Reference a GcpServiceAccount resource — its `name` output is exactly
this value. There is no separate project field: the account's project is
embedded in the resource name.

- references: GcpServiceAccount (`status.outputs.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpServiceAccount, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.role

`string | valueFrom` · required

The role to grant on the service account. Either a predefined role
(typically one of the service-account usage roles listed above) or a
custom role's fully-qualified name
("projects/<project>/roles/<role_id>"). Reference a GcpIamCustomRole
resource to grant a custom role — its `name` output is exactly this value.

- references: GcpIamCustomRole (`status.outputs.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpIamCustomRole, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.member

`string | valueFrom` · required

The identity receiving the grant, in GCP IAM member format:
  serviceAccount:<email>       — another service account (reference a
                                 GcpServiceAccount resource — its `member`
                                 output is exactly this value)
  principal:// / principalSet:// — workload identity federation principals;
                                 the workloadIdentityUser + principalSet
                                 pair is how an external identity (e.g. a
                                 GitHub repository) gains impersonation
  user:<email>                 — a Google account
  group:<email>                — a Google group
  domain:<domain>              — everyone in a Workspace/Cloud Identity domain
Grants to deleted principals ("deleted:...") are not supported.
Format validation happens at deploy time in the modules rather than here,
because the value usually arrives through a reference that is only
resolved at deploy time.

- references: GcpServiceAccount (`status.outputs.member`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpServiceAccount, name: <that resource's name>, fieldPath: status.outputs.member}} -- a bare string does not parse

### spec.condition

`GcpServiceAccountIamMemberCondition`

Optional IAM Condition restricting when this grant applies (e.g. only
before an expiry date). The condition is part of the grant's identity:
the same role with and without a condition are two independent grants
that do not interfere.

### spec.condition.title

`string` · required

Short human-readable title identifying the condition's intent,
e.g. "expires-2026-12-31".

- rule: {"required":true,"string":{"maxLen":"100"}}

### spec.condition.expression

`string` · required

The CEL condition expression, e.g.
request.time < timestamp("2027-01-01T00:00:00Z").

- rule: {"required":true}

### spec.condition.description

`string`

Optional longer explanation of what the condition does and why it exists.

- rule: {"string":{"maxLen":"256"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: GcpServiceAccountIamMember, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.service_account_id` | `string` | The fully-qualified resource name of the service account whose IAM policy received the grant (projects/<project>/serviceAccounts/<email>), after reference resolution. |
| `status.outputs.role` | `string` | The role that was granted (predefined name or fully-qualified custom role name), after reference resolution. |
| `status.outputs.member` | `string` | The member the role was granted to, in IAM member format, after reference resolution. |
| `status.outputs.etag` | `string` | The etag of the service account IAM policy after this grant was applied — a fingerprint of the policy version, useful for audit correlation. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serviceAccountId` | GcpServiceAccount | `status.outputs.name` |
| `spec.role` | GcpIamCustomRole | `status.outputs.name` |
| `spec.member` | GcpServiceAccount | `status.outputs.member` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
