# OciCompartment

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciCompartmentSpec defines the specification for an Oracle Cloud Infrastructure
compartment.

Compartments are the fundamental building block of OCI's resource isolation
model, analogous to GCP folders or AWS Organizational Units. Every OCI
resource exists within exactly one compartment, forming a hierarchy rooted
at the tenancy. Compartments enable fine-grained access control via IAM
policies, cost tracking, and logical grouping of resources.

This component creates a single compartment within a parent compartment.
Nested hierarchies are built by chaining OciCompartment resources -- each
child references its parent via compartment_id.

## Example

```yaml
apiVersion: oci.planton.dev/v1alpha1
kind: OciCompartment
metadata:
  name: ocicompartment-demo
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  name: "demo-compartment"
  description: "Demo compartment for testing OciCompartment deployment"
  enableDelete: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.name` | `string` |  |  |  |
| `spec.description` | `string` | yes |  |  |
| `spec.enableDelete` | `bool` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the parent compartment where this compartment will be created.
For top-level compartments, this should be the tenancy OCID.
For nested compartments, this is the OCID of the parent compartment.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.name

`string`

Name assigned to the compartment. Must be unique across all compartments
within the parent. This is the identifier shown in the OCI Console and
used in IAM policy statements.
Falls back to metadata.name if not provided.

### spec.description

`string` · required

Description of the compartment's purpose.
Required by the OCI API. Use this to document what resources or teams
the compartment is intended for.

- rule: {"string":{"minLen":"1"}}

### spec.enableDelete

`bool`

When true, the compartment will be deleted on resource destruction.
When false (the default), the compartment is retained even after the
IaC resource is destroyed -- this is OCI's safety mechanism to prevent
accidental deletion of compartments containing active resources.
Set to true only for ephemeral or development compartments.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciCompartment, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.compartment_id` | `string` | OCID of the created compartment. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciAlarm | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciAlarm | `spec.metricCompartmentId` | `status.outputs.compartment_id` |
| OciApiGateway | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciApplicationLoadBalancer | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciAutonomousDatabase | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciBastion | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciBlockVolume | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciCompartment | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciComputeInstance | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciContainerEngineCluster | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciContainerEngineNodePool | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciContainerInstance | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciDbSystem | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciDevopsProject | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciDnsZone | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciDynamicGroup | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciDynamicRoutingGateway | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciFileSystem | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciFunctionsApplication | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciIdentityPolicy | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciKmsKey | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciKmsVault | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciLogGroup | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciLogGroup | `spec.logs[].configuration.compartmentId` | `status.outputs.compartment_id` |
| OciMysqlDbSystem | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciNetworkFirewall | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciNetworkLoadBalancer | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciNosqlTable | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciObjectStorageBucket | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciPostgresqlDbSystem | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciPublicIp | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciQueue | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciRedisCluster | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciSecurityGroup | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciStreamPool | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciSubnet | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciVaultSecret | `spec.compartmentId` | `status.outputs.compartment_id` |
| OciVcn | `spec.compartmentId` | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
