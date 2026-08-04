# OciSubnet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciSubnetSpec defines the specification for an Oracle Cloud Infrastructure
subnet within a Virtual Cloud Network (VCN).

A subnet is a subdivision of a VCN that provides a contiguous range of IP
addresses. Subnets can be regional (span all availability domains) or
AD-specific. Private subnets prevent public IP assignment on VNICs; public
subnets allow it.

An optional custom route table can be created inline via route_rules. When
route_rules are provided, a dedicated route table is created and attached to
this subnet. Alternatively, an existing route table can be referenced via
route_table_id. If neither is provided, the VCN's default route table is
used. The two fields are mutually exclusive.

## Example

```yaml
apiVersion: oci.planton.dev/v1
kind: OciSubnet
metadata:
  name: ocisubnet-demo
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  vcnId:
    value: "ocid1.vcn.oc1.iad.example"
  cidrBlock: "10.0.1.0/24"
  displayName: "demo-private-subnet"
  dnsLabel: "priv1"
  prohibitPublicIpOnVnic: true
  prohibitInternetIngress: true
  routeRules:
    - destination: "0.0.0.0/0"
      destinationType: "cidr_block"
      networkEntityId:
        value: "ocid1.natgateway.oc1.iad.example"
      description: "Route internet traffic via NAT gateway"
    - destination: "all-iad-services-in-oracle-services-network"
      destinationType: "service_cidr_block"
      networkEntityId:
        value: "ocid1.servicegateway.oc1.iad.example"
      description: "Route OCI services via service gateway"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.vcnId` | `string \| valueFrom` | yes |  | OciVcn (`status.outputs.vcn_id`) |
| `spec.cidrBlock` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.dnsLabel` | `string` |  |  |  |
| `spec.availabilityDomain` | `string` |  |  |  |
| `spec.prohibitPublicIpOnVnic` | `bool` |  |  |  |
| `spec.prohibitInternetIngress` | `bool` |  |  |  |
| `spec.dhcpOptionsId` | `string \| valueFrom` |  |  |  |
| `spec.routeTableId` | `string \| valueFrom` |  |  |  |
| `spec.securityListIds` | `[]string \| valueFrom` |  |  |  |
| `spec.ipv6CidrBlock` | `string` |  |  |  |
| `spec.routeRules` | `[]RouteRule` |  |  |  |
| `spec.routeRules[].destination` | `string` | yes |  |  |
| `spec.routeRules[].destinationType` | `enum` |  |  |  |
| `spec.routeRules[].networkEntityId` | `string \| valueFrom` | yes |  |  |
| `spec.routeRules[].description` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the subnet will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.vcnId

`string | valueFrom` · required

OCID of the VCN that this subnet belongs to.

- references: OciVcn (`status.outputs.vcn_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciVcn, name: <that resource's name>, fieldPath: status.outputs.vcn_id}} -- a bare string does not parse

### spec.cidrBlock

`string` · required

IPv4 CIDR block for the subnet (e.g. "10.0.1.0/24").
Must be within one of the VCN's CIDR blocks and not overlap with other subnets.

- rule: {"required":true}

### spec.displayName

`string`

Human-readable name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.dnsLabel

`string`

DNS label for the subnet. Combined with the VCN and subnet domain, forms
the FQDN: <dns_label>.<vcn_dns_label>.oraclevcn.com
Must be alphanumeric, start with a letter, and be at most 15 characters.

### spec.availabilityDomain

`string`

Availability domain name (e.g. "Iocq:US-ASHBURN-AD-1").
When omitted, the subnet is regional and spans all ADs in the region.
When set, the subnet is scoped to a single AD.

### spec.prohibitPublicIpOnVnic

`bool`

When true, VNICs in this subnet cannot have public IP addresses.
This is the primary control for making a subnet private.

### spec.prohibitInternetIngress

`bool`

When true, the subnet blocks all inbound internet traffic to VNICs,
even if a security rule or NSG would otherwise allow it.

### spec.dhcpOptionsId

`string | valueFrom`

OCID of custom DHCP options to use instead of the VCN's default.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.routeTableId

`string | valueFrom`

OCID of an existing route table to associate with this subnet.
Mutually exclusive with route_rules. If neither is provided, the VCN's
default route table is used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.securityListIds

`[]string | valueFrom`

Security list OCIDs to associate with this subnet. OCI allows a maximum
of 5 security lists per subnet.

- rule: {"repeated":{"maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.ipv6CidrBlock

`string`

IPv6 CIDR block for dual-stack subnets (e.g. "2001:0db8:0123:1111::/64").
Only valid when the parent VCN has IPv6 enabled.

### spec.routeRules

`[]RouteRule`

Route rules for a custom route table that will be created and owned by
this subnet. Mutually exclusive with route_table_id.

### spec.routeRules[].destination

`string` · required

Target IP range in CIDR notation (e.g. "0.0.0.0/0") or a service CIDR
label (e.g. "all-iad-services-in-oracle-services-network").

- rule: {"required":true}

### spec.routeRules[].destinationType

`enum`

Whether destination is a CIDR block or a service CIDR block.

Allowed values (use exactly as shown):

- `unspecified`
- `cidr_block`
- `service_cidr_block`

### spec.routeRules[].networkEntityId

`string | valueFrom` · required

OCID of the network entity to route matching traffic to
(Internet Gateway, NAT Gateway, DRG, Service Gateway, Local Peering Gateway, etc.).

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.routeRules[].description

`string`

Optional human-readable description for this rule.

## Validation Rules

- `route_table_mutual_exclusivity`: route_table_id and route_rules are mutually exclusive; provide one or neither

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciSubnet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.subnet_id` | `string` | OCID of the subnet. |
| `status.outputs.subnet_domain_name` | `string` | Fully qualified domain name of the subnet (e.g. "subnet1.vcn1.oraclevcn.com"). |
| `status.outputs.virtual_router_ip` | `string` | IP address of the virtual router in this subnet. |
| `status.outputs.virtual_router_mac` | `string` | MAC address of the virtual router in this subnet. |
| `status.outputs.route_table_id` | `string` | OCID of the route table associated with this subnet. This is either the custom route table created from route_rules, the externally referenced route_table_id, or the VCN's default. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.vcnId` | OciVcn | `status.outputs.vcn_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciApiGateway | `spec.subnetId` | `status.outputs.subnet_id` |
| OciApplicationLoadBalancer | `spec.subnetIds` | `status.outputs.subnet_id` |
| OciAutonomousDatabase | `spec.subnetId` | `status.outputs.subnet_id` |
| OciBastion | `spec.targetSubnetId` | `status.outputs.subnet_id` |
| OciComputeInstance | `spec.createVnicDetails.subnetId` | `status.outputs.subnet_id` |
| OciContainerEngineCluster | `spec.endpointConfig.subnetId` | `status.outputs.subnet_id` |
| OciContainerEngineCluster | `spec.options.serviceLbSubnetIds` | `status.outputs.subnet_id` |
| OciContainerEngineNodePool | `spec.nodeConfigDetails.placementConfigs[].subnetId` | `status.outputs.subnet_id` |
| OciContainerEngineNodePool | `spec.nodeConfigDetails.podNetworkOptionDetails.podSubnetIds` | `status.outputs.subnet_id` |
| OciContainerInstance | `spec.vnics[].subnetId` | `status.outputs.subnet_id` |
| OciDbSystem | `spec.subnetId` | `status.outputs.subnet_id` |
| OciDbSystem | `spec.backupSubnetId` | `status.outputs.subnet_id` |
| OciFileSystem | `spec.mountTarget.subnetId` | `status.outputs.subnet_id` |
| OciFunctionsApplication | `spec.subnetIds` | `status.outputs.subnet_id` |
| OciMysqlDbSystem | `spec.subnetId` | `status.outputs.subnet_id` |
| OciNetworkFirewall | `spec.subnetId` | `status.outputs.subnet_id` |
| OciNetworkLoadBalancer | `spec.subnetId` | `status.outputs.subnet_id` |
| OciPostgresqlDbSystem | `spec.networkDetails.subnetId` | `status.outputs.subnet_id` |
| OciRedisCluster | `spec.subnetId` | `status.outputs.subnet_id` |
| OciStreamPool | `spec.privateEndpointSettings.subnetId` | `status.outputs.subnet_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
