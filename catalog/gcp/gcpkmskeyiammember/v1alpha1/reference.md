# GcpKmsKeyIamMember

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `gcp.planton.dev/v1alpha1`

GcpKmsKeyIamMemberSpec defines a single ADDITIVE IAM grant ON a KMS crypto
key: one role, to one member, on one key — the least-privilege unit of
customer-managed encryption (CMEK) access control.

The canonical use is granting a Google service agent permission to use a
key for CMEK:
  roles/cloudkms.cryptoKeyEncrypterDecrypter — the role every CMEK
    consumer's service agent needs on the key (Cloud Storage's
    service-<project_number>@gs-project-accounts.iam.gserviceaccount.com,
    BigQuery's bq-<project_number>@..., and so on). Granting it here, on
    the key, is preferable to a project-wide grant: the agent can use THIS
    key and nothing else, and the grant is an explicit dependency edge so
    orchestration orders the encrypted resource after the permission it
    needs exists.
  roles/cloudkms.viewer / roles/cloudkms.admin — read or manage the key's
    metadata and policy without touching key material.
IAM on a key also flows DOWN from its key ring and project — a ring-level
grant covers every key in the ring. Use key-scoped grants when different
keys in a ring have different consumers.

Additive semantics: this grant merges into the key's IAM policy without
touching any other member's bindings on the same role, and removal
subtracts only this exact (role, member) pair. Multiple grants — from
different charts, teams, or tools — never fight over the policy.
(Authoritative binding/policy management, which clobbers everything not
listed, is deliberately not modeled.)

Every field is immutable, mirroring the underlying API: an IAM grant has no
update — changing the key, role, member, or condition replaces the grant
atomically.

## Example

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpKmsKeyIamMember
metadata:
  name: my-sample-key-grant
spec:
  # The crypto key whose IAM policy receives this grant, as its
  # fully-qualified resource path. Reference a GcpKmsKey — its `key_id`
  # output is exactly this value.
  cryptoKeyId:
    value: projects/my-gcp-project-123/locations/us-central1/keyRings/app-ring/cryptoKeys/state-key

  # The role to grant ON the key: typically the CMEK consumer role.
  role:
    value: roles/cloudkms.cryptoKeyEncrypterDecrypter

  # The identity receiving the grant, in GCP IAM member format — most
  # commonly the consuming service's agent for CMEK.
  member:
    value: serviceAccount:service-123456789@gs-project-accounts.iam.gserviceaccount.com

  # Optional IAM Condition restricting when the grant applies
  # condition:
  #   title: expires-2027
  #   expression: request.time < timestamp("2027-01-01T00:00:00Z")
  #   description: Temporary access for the migration
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.cryptoKeyId` | `string \| valueFrom` | yes |  | GcpKmsKey (`status.outputs.key_id`) |
| `spec.role` | `string \| valueFrom` | yes |  | GcpIamCustomRole (`status.outputs.name`) |
| `spec.member` | `string \| valueFrom` | yes |  | GcpServiceAccount (`status.outputs.member`) |
| `spec.condition` | `GcpKmsKeyIamMemberCondition` |  |  |  |
| `spec.condition.title` | `string` | yes |  |  |
| `spec.condition.expression` | `string` | yes |  |  |
| `spec.condition.description` | `string` |  |  |  |

## Field Details

### spec.cryptoKeyId

`string | valueFrom` · required

The crypto key whose IAM policy receives this grant, as its
fully-qualified resource path:
  projects/<project>/locations/<location>/keyRings/<ring>/cryptoKeys/<key>
Reference a GcpKmsKey resource — its `key_id` output is exactly this
value. There is no separate project or location field: both are embedded
in the key path.

- references: GcpKmsKey (`status.outputs.key_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.role

`string | valueFrom` · required

The role to grant on the key. Either a predefined role (typically
"roles/cloudkms.cryptoKeyEncrypterDecrypter" for CMEK consumers) or a
custom role's fully-qualified name
("projects/<project>/roles/<role_id>"). Reference a GcpIamCustomRole
resource to grant a custom role — its `name` output is exactly this value.

- references: GcpIamCustomRole (`status.outputs.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpIamCustomRole, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.member

`string | valueFrom` · required

The identity receiving the grant, in GCP IAM member format:
  serviceAccount:<email>       — a service account or Google service
                                 agent (the most common for CMEK — the
                                 consuming service's agent, e.g.
                                 service-<project_number>@gs-project-accounts.iam.gserviceaccount.com;
                                 for your own accounts reference a
                                 GcpServiceAccount resource — its `member`
                                 output is exactly this value)
  user:<email>                 — a Google account
  group:<email>                — a Google group
  domain:<domain>              — everyone in a Workspace/Cloud Identity domain
  principal:// / principalSet:// — workload identity federation principals
Grants to deleted principals ("deleted:...") are not supported.
Format validation happens at deploy time in the modules rather than here,
because the value usually arrives through a reference that is only
resolved at deploy time.

- references: GcpServiceAccount (`status.outputs.member`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpServiceAccount, name: <that resource's name>, fieldPath: status.outputs.member}} -- a bare string does not parse

### spec.condition

`GcpKmsKeyIamMemberCondition`

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

Reference an output from another manifest as `valueFrom: {kind: GcpKmsKeyIamMember, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.crypto_key_id` | `string` | The fully-qualified resource path of the crypto key whose IAM policy received the grant (projects/<project>/locations/<location>/keyRings/<ring>/cryptoKeys/<key>), after reference resolution. |
| `status.outputs.role` | `string` | The role that was granted (predefined name or fully-qualified custom role name), after reference resolution. |
| `status.outputs.member` | `string` | The member the role was granted to, in IAM member format, after reference resolution. |
| `status.outputs.etag` | `string` | The etag of the key IAM policy after this grant was applied — a fingerprint of the policy version, useful for audit correlation. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.cryptoKeyId` | GcpKmsKey | `status.outputs.key_id` |
| `spec.role` | GcpIamCustomRole | `status.outputs.name` |
| `spec.member` | GcpServiceAccount | `status.outputs.member` |

## See Also

- [Overview](../README.md)
