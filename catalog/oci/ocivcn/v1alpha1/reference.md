# OciVcn

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciVcnSpec defines the specification for an Oracle Cloud Infrastructure
Virtual Cloud Network (VCN) and its optional gateway sub-resources.

A VCN is the foundational networking construct in OCI, analogous to an AWS
VPC. It provides an isolated virtual network within a compartment and
supports multiple CIDR blocks, DNS resolution, and IPv6.

Gateway resources (Internet, NAT, Service) are bundled into this component
because they are tightly coupled to the VCN lifecycle and rarely useful on
their own. Each gateway is controlled by a boolean toggle and only created
when enabled.

## Example

```yaml
apiVersion: oci.planton.dev/v1alpha1
kind: OciVcn
metadata:
  name: ocivcn-demo
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  cidrBlocks:
    - "10.0.0.0/16"
  displayName: "demo-vcn"
  dnsLabel: "demovcn"
  isIpv6Enabled: false
  isInternetGatewayEnabled: true
  isNatGatewayEnabled: true
  isServiceGatewayEnabled: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.cidrBlocks` | `[]string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.dnsLabel` | `string` |  |  |  |
| `spec.isIpv6Enabled` | `bool` |  |  |  |
| `spec.isInternetGatewayEnabled` | `bool` |  |  |  |
| `spec.isNatGatewayEnabled` | `bool` |  |  |  |
| `spec.isServiceGatewayEnabled` | `bool` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the VCN and its gateways will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.cidrBlocks

`[]string` · required

IPv4 CIDR blocks for the VCN. At least one is required.
OCI supports multiple non-overlapping CIDRs per VCN (e.g. ["10.0.0.0/16", "172.16.0.0/16"]).
Each block must be between /16 and /30.

- rule: {"repeated":{"minItems":"1"}}

### spec.displayName

`string`

Human-readable name for the VCN shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.dnsLabel

`string`

DNS label for the VCN. Enables DNS hostnames within the VCN.
Must be alphanumeric, start with a letter, and be at most 15 characters.
When set, the VCN domain becomes: <dns_label>.oraclevcn.com

### spec.isIpv6Enabled

`bool`

When true, allocate an Oracle-assigned /56 IPv6 GUA prefix for the VCN.

### spec.isInternetGatewayEnabled

`bool`

When true, create an Internet Gateway attached to this VCN.
Required for resources that need direct inbound/outbound internet access.

### spec.isNatGatewayEnabled

`bool`

When true, create a NAT Gateway attached to this VCN.
Allows private resources to initiate outbound internet connections
without exposing them to inbound traffic.

### spec.isServiceGatewayEnabled

`bool`

When true, create a Service Gateway attached to this VCN.
Provides private access to OCI services (Object Storage, etc.) without
traffic leaving the Oracle network. The gateway is automatically
configured for "All Services in Oracle Services Network".

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciVcn, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vcn_id` | `string` | OCID of the VCN. |
| `status.outputs.default_route_table_id` | `string` | OCID of the default route table created with the VCN. |
| `status.outputs.default_security_list_id` | `string` | OCID of the default security list created with the VCN. |
| `status.outputs.default_dhcp_options_id` | `string` | OCID of the default DHCP options created with the VCN. |
| `status.outputs.internet_gateway_id` | `string` | OCID of the Internet Gateway. Empty when is_internet_gateway_enabled is false. |
| `status.outputs.nat_gateway_id` | `string` | OCID of the NAT Gateway. Empty when is_nat_gateway_enabled is false. |
| `status.outputs.service_gateway_id` | `string` | OCID of the Service Gateway. Empty when is_service_gateway_enabled is false. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciContainerEngineCluster | `spec.vcnId` | `status.outputs.vcn_id` |
| OciSecurityGroup | `spec.vcnId` | `status.outputs.vcn_id` |
| OciSubnet | `spec.vcnId` | `status.outputs.vcn_id` |

## See Also

- [Overview](../README.md)
