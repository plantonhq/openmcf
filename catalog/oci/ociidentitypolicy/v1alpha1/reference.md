# OciIdentityPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciIdentityPolicySpec defines the specification for an Oracle Cloud
Infrastructure IAM policy.

Policies are the mechanism for granting access in OCI. Each policy consists
of one or more human-readable statements written in OCI's policy language
(e.g. "Allow group Admins to manage all-resources in compartment Production").
Policies are attached to a compartment and grant permissions within that
compartment and all of its children.

Policy names must be unique across the tenancy and cannot be changed after
creation. Statements and description are updatable.

## Example

```yaml
apiVersion: oci.planton.dev/v1alpha1
kind: OciIdentityPolicy
metadata:
  name: ociidentitypolicy-demo
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  name: "demo-admin-policy"
  description: "Demo policy granting admin access for testing OciIdentityPolicy deployment"
  statements:
    - "Allow group Admins to manage all-resources in compartment DemoCompartment"
    - "Allow group Developers to use instances in compartment DemoCompartment"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.name` | `string` |  |  |  |
| `spec.description` | `string` | yes |  |  |
| `spec.statements` | `[]string` | yes |  |  |
| `spec.versionDate` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where this policy will be created.
For tenancy-level policies, use the tenancy OCID.
For compartment-scoped policies, use the target compartment's OCID.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.name

`string`

Name assigned to the policy. Must be unique across all policies in the
tenancy and cannot be changed after creation.
Falls back to metadata.name if not provided.

### spec.description

`string` · required

Description of the policy's purpose.
Required by the OCI API. Updatable after creation.

- rule: {"string":{"minLen":"1"}}

### spec.statements

`[]string` · required

Policy statements written in OCI's policy language.
At least one statement is required. Each statement follows the syntax:
  "Allow <subject> to <verb> <resource-type> in <location> [where <conditions>]"
See https://docs.oracle.com/iaas/Content/Identity/Concepts/policies.htm

- rule: {"repeated":{"minItems":"1"}}

### spec.versionDate

`string`

Version date for policy evaluation (YYYY-MM-DD format).
When set, the policy is evaluated according to the behavior of OCI
services on that date, providing a stable policy interpretation.
When empty, the policy uses the current service behavior at evaluation time.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciIdentityPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.policy_id` | `string` | OCID of the created policy. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](../README.md)
