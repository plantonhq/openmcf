# OciDynamicGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciDynamicGroupSpec defines the specification for an Oracle Cloud
Infrastructure dynamic group.

Dynamic groups allow OCI resources (compute instances, functions, etc.) to
be grouped and granted permissions via IAM policies -- enabling patterns
like instance principal authentication where workloads authenticate to
OCI services without stored credentials.

The matching rule uses OCI's rule syntax to select which resources belong
to the group. Dynamic groups must be created in the tenancy (root)
compartment.

## Example

```yaml
apiVersion: oci.planton.dev/v1alpha1
kind: OciDynamicGroup
metadata:
  name: ocidynamicgroup-demo
spec:
  compartmentId:
    value: "ocid1.tenancy.oc1..example"
  name: "demo-compute-dynamic-group"
  description: "Dynamic group for compute instances in the demo compartment"
  matchingRule: "Any {instance.compartment.id = 'ocid1.compartment.oc1..example'}"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.name` | `string` |  |  |  |
| `spec.description` | `string` | yes |  |  |
| `spec.matchingRule` | `string` | yes |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the tenancy (root compartment).
Dynamic groups are tenancy-level IAM resources and must be created
in the tenancy compartment, not a child compartment.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.name

`string`

Name assigned to the dynamic group.
Must be unique across all groups (including user groups) in the tenancy
and cannot be changed after creation.
Falls back to metadata.name if not provided.

### spec.description

`string` · required

Description of the dynamic group's purpose.
Required by the OCI API. Updatable after creation.

- rule: {"string":{"minLen":"1"}}

### spec.matchingRule

`string` · required

Matching rule that defines which resources belong to this dynamic group.
Uses OCI's rule syntax, e.g.:
  "Any {instance.compartment.id = 'ocid1.compartment.oc1..xxx'}"
  "All {resource.type = 'fnfunc', resource.compartment.id = 'ocid1...'}"
See https://docs.oracle.com/iaas/Content/Identity/Tasks/managingdynamicgroups.htm

- rule: {"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciDynamicGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.dynamic_group_id` | `string` | OCID of the created dynamic group. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
