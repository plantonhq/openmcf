# OciSecurityGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciSecurityGroupSpec defines the specification for an Oracle Cloud
Infrastructure Network Security Group (NSG) and its inline security rules.

An NSG acts as a virtual firewall for compute instances and other resources.
Unlike security lists (which apply at the subnet level), NSGs provide
per-VNIC traffic control and are OCI's recommended approach for fine-grained
network security. Each NSG belongs to a single VCN and contains up to 120
security rules total (ingress + egress combined).

Rules are split into ingress_rules and egress_rules so that direction is
implicit from the field name, eliminating the error-prone "direction +
conditional source/destination" pattern from the raw provider API.

## Example

```yaml
apiVersion: oci.planton.dev/v1alpha1
kind: OciSecurityGroup
metadata:
  name: web-tier-nsg
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  vcnId:
    value: "ocid1.vcn.oc1.iad.example"
  displayName: "Web Tier NSG"
  ingressRules:
    - source: "0.0.0.0/0"
      sourceType: cidr_block
      protocol: tcp
      description: "Allow HTTPS from internet"
      tcpOptions:
        destinationPortRange:
          min: 443
          max: 443
    - source: "0.0.0.0/0"
      sourceType: cidr_block
      protocol: tcp
      description: "Allow HTTP from internet"
      tcpOptions:
        destinationPortRange:
          min: 80
          max: 80
    - source: "10.0.0.0/16"
      sourceType: cidr_block
      protocol: icmp
      description: "Allow ICMP type 3 code 4 from VCN (Path MTU Discovery)"
      icmpOptions:
        type: 3
        code: 4
    - source: "10.0.0.0/16"
      sourceType: cidr_block
      protocol: icmp
      description: "Allow ICMP type 3 from VCN (Destination Unreachable)"
      icmpOptions:
        type: 3
  egressRules:
    - destination: "0.0.0.0/0"
      destinationType: cidr_block
      protocol: all
      description: "Allow all outbound traffic"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.vcnId` | `string \| valueFrom` | yes |  | OciVcn (`status.outputs.vcn_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.ingressRules` | `[]IngressRule` |  |  |  |
| `spec.ingressRules[].source` | `string` | yes |  |  |
| `spec.ingressRules[].sourceType` | `enum` |  |  |  |
| `spec.ingressRules[].protocol` | `enum` | yes |  |  |
| `spec.ingressRules[].description` | `string` |  |  |  |
| `spec.ingressRules[].stateless` | `bool` |  |  |  |
| `spec.ingressRules[].tcpOptions` | `TcpOptions` |  |  |  |
| `spec.ingressRules[].tcpOptions.destinationPortRange` | `PortRange` |  |  |  |
| `spec.ingressRules[].tcpOptions.destinationPortRange.min` | `int32` |  |  |  |
| `spec.ingressRules[].tcpOptions.destinationPortRange.max` | `int32` |  |  |  |
| `spec.ingressRules[].tcpOptions.sourcePortRange` | `PortRange` |  |  |  |
| `spec.ingressRules[].tcpOptions.sourcePortRange.min` | `int32` |  |  |  |
| `spec.ingressRules[].tcpOptions.sourcePortRange.max` | `int32` |  |  |  |
| `spec.ingressRules[].udpOptions` | `UdpOptions` |  |  |  |
| `spec.ingressRules[].udpOptions.destinationPortRange` | `PortRange` |  |  |  |
| `spec.ingressRules[].udpOptions.destinationPortRange.min` | `int32` |  |  |  |
| `spec.ingressRules[].udpOptions.destinationPortRange.max` | `int32` |  |  |  |
| `spec.ingressRules[].udpOptions.sourcePortRange` | `PortRange` |  |  |  |
| `spec.ingressRules[].udpOptions.sourcePortRange.min` | `int32` |  |  |  |
| `spec.ingressRules[].udpOptions.sourcePortRange.max` | `int32` |  |  |  |
| `spec.ingressRules[].icmpOptions` | `IcmpOptions` |  |  |  |
| `spec.ingressRules[].icmpOptions.type` | `int32` |  |  |  |
| `spec.ingressRules[].icmpOptions.code` | `int32` |  |  |  |
| `spec.egressRules` | `[]EgressRule` |  |  |  |
| `spec.egressRules[].destination` | `string` | yes |  |  |
| `spec.egressRules[].destinationType` | `enum` |  |  |  |
| `spec.egressRules[].protocol` | `enum` | yes |  |  |
| `spec.egressRules[].description` | `string` |  |  |  |
| `spec.egressRules[].stateless` | `bool` |  |  |  |
| `spec.egressRules[].tcpOptions` | `TcpOptions` |  |  |  |
| `spec.egressRules[].tcpOptions.destinationPortRange` | `PortRange` |  |  |  |
| `spec.egressRules[].tcpOptions.destinationPortRange.min` | `int32` |  |  |  |
| `spec.egressRules[].tcpOptions.destinationPortRange.max` | `int32` |  |  |  |
| `spec.egressRules[].tcpOptions.sourcePortRange` | `PortRange` |  |  |  |
| `spec.egressRules[].tcpOptions.sourcePortRange.min` | `int32` |  |  |  |
| `spec.egressRules[].tcpOptions.sourcePortRange.max` | `int32` |  |  |  |
| `spec.egressRules[].udpOptions` | `UdpOptions` |  |  |  |
| `spec.egressRules[].udpOptions.destinationPortRange` | `PortRange` |  |  |  |
| `spec.egressRules[].udpOptions.destinationPortRange.min` | `int32` |  |  |  |
| `spec.egressRules[].udpOptions.destinationPortRange.max` | `int32` |  |  |  |
| `spec.egressRules[].udpOptions.sourcePortRange` | `PortRange` |  |  |  |
| `spec.egressRules[].udpOptions.sourcePortRange.min` | `int32` |  |  |  |
| `spec.egressRules[].udpOptions.sourcePortRange.max` | `int32` |  |  |  |
| `spec.egressRules[].icmpOptions` | `IcmpOptions` |  |  |  |
| `spec.egressRules[].icmpOptions.type` | `int32` |  |  |  |
| `spec.egressRules[].icmpOptions.code` | `int32` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the NSG will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.vcnId

`string | valueFrom` · required

OCID of the VCN that this NSG belongs to. Changing this forces recreation.

- references: OciVcn (`status.outputs.vcn_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciVcn, name: <that resource's name>, fieldPath: status.outputs.vcn_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.ingressRules

`[]IngressRule`

Inbound security rules. Each rule defines traffic that is allowed TO
resources associated with this NSG.

### spec.ingressRules[].source

`string` · required

Traffic source: a CIDR block (e.g. "10.0.0.0/16"), a service CIDR label
(e.g. "all-iad-services-in-oracle-services-network"), or the OCID of
another NSG.

- rule: {"string":{"minLen":"1"}}

### spec.ingressRules[].sourceType

`enum`

Interpretation of the source field. Defaults to cidr_block when unset.

Allowed values (use exactly as shown):

- `target_type_unspecified`
- `cidr_block`
- `service_cidr_block`
- `network_security_group`

### spec.ingressRules[].protocol

`enum` · required

Transport protocol for this rule.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `all`
- `icmp`
- `tcp`
- `udp`
- `icmpv6`

### spec.ingressRules[].description

`string`

Optional human-readable description.

### spec.ingressRules[].stateless

`bool`

When true the rule is stateless (return traffic must be explicitly
allowed). When false (the default) the rule is stateful.

### spec.ingressRules[].tcpOptions

`TcpOptions`

TCP port constraints. Only valid when protocol is tcp.

### spec.ingressRules[].tcpOptions.destinationPortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.ingressRules[].tcpOptions.destinationPortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].tcpOptions.destinationPortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].tcpOptions.sourcePortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.ingressRules[].tcpOptions.sourcePortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].tcpOptions.sourcePortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].udpOptions

`UdpOptions`

UDP port constraints. Only valid when protocol is udp.

### spec.ingressRules[].udpOptions.destinationPortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.ingressRules[].udpOptions.destinationPortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].udpOptions.destinationPortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].udpOptions.sourcePortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.ingressRules[].udpOptions.sourcePortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].udpOptions.sourcePortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.ingressRules[].icmpOptions

`IcmpOptions`

ICMP type/code constraints. Only valid when protocol is icmp or icmpv6.

### spec.ingressRules[].icmpOptions.type

`int32`

ICMP message type (e.g. 3 for "Destination Unreachable", 8 for "Echo").

### spec.ingressRules[].icmpOptions.code

`int32` · optional (explicit presence)

ICMP message code. When omitted, all codes for the given type are matched.
Use optional so that code=0 (a valid ICMP code) is distinguishable from "not set".

### spec.egressRules

`[]EgressRule`

Outbound security rules. Each rule defines traffic that is allowed FROM
resources associated with this NSG.

### spec.egressRules[].destination

`string` · required

Traffic destination: a CIDR block, a service CIDR label, or the OCID of
another NSG.

- rule: {"string":{"minLen":"1"}}

### spec.egressRules[].destinationType

`enum`

Interpretation of the destination field. Defaults to cidr_block when unset.

Allowed values (use exactly as shown):

- `target_type_unspecified`
- `cidr_block`
- `service_cidr_block`
- `network_security_group`

### spec.egressRules[].protocol

`enum` · required

Transport protocol for this rule.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `all`
- `icmp`
- `tcp`
- `udp`
- `icmpv6`

### spec.egressRules[].description

`string`

Optional human-readable description.

### spec.egressRules[].stateless

`bool`

When true the rule is stateless. When false (the default) the rule is stateful.

### spec.egressRules[].tcpOptions

`TcpOptions`

TCP port constraints. Only valid when protocol is tcp.

### spec.egressRules[].tcpOptions.destinationPortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.egressRules[].tcpOptions.destinationPortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].tcpOptions.destinationPortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].tcpOptions.sourcePortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.egressRules[].tcpOptions.sourcePortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].tcpOptions.sourcePortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].udpOptions

`UdpOptions`

UDP port constraints. Only valid when protocol is udp.

### spec.egressRules[].udpOptions.destinationPortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.egressRules[].udpOptions.destinationPortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].udpOptions.destinationPortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].udpOptions.sourcePortRange

`PortRange`

- rule: min must be less than or equal to max

### spec.egressRules[].udpOptions.sourcePortRange.min

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].udpOptions.sourcePortRange.max

`int32`

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.egressRules[].icmpOptions

`IcmpOptions`

ICMP type/code constraints. Only valid when protocol is icmp or icmpv6.

### spec.egressRules[].icmpOptions.type

`int32`

ICMP message type (e.g. 3 for "Destination Unreachable", 8 for "Echo").

### spec.egressRules[].icmpOptions.code

`int32` · optional (explicit presence)

ICMP message code. When omitted, all codes for the given type are matched.
Use optional so that code=0 (a valid ICMP code) is distinguishable from "not set".

## Validation Rules

- `max_120_rules`: an NSG supports at most 120 rules (ingress + egress combined)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciSecurityGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_security_group_id` | `string` | OCID of the network security group. |

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
| OciApiGateway | `spec.networkSecurityGroupIds` | `status.outputs.network_security_group_id` |
| OciApplicationLoadBalancer | `spec.networkSecurityGroupIds` | `status.outputs.network_security_group_id` |
| OciAutonomousDatabase | `spec.nsgIds` | `status.outputs.network_security_group_id` |
| OciComputeInstance | `spec.createVnicDetails.nsgIds` | `status.outputs.network_security_group_id` |
| OciContainerEngineCluster | `spec.endpointConfig.nsgIds` | `status.outputs.network_security_group_id` |
| OciContainerEngineCluster | `spec.options.serviceLbConfig.backendNsgIds` | `status.outputs.network_security_group_id` |
| OciContainerEngineNodePool | `spec.nodeConfigDetails.nsgIds` | `status.outputs.network_security_group_id` |
| OciContainerEngineNodePool | `spec.nodeConfigDetails.podNetworkOptionDetails.podNsgIds` | `status.outputs.network_security_group_id` |
| OciContainerInstance | `spec.vnics[].nsgIds` | `status.outputs.network_security_group_id` |
| OciDbSystem | `spec.nsgIds` | `status.outputs.network_security_group_id` |
| OciDbSystem | `spec.backupNetworkNsgIds` | `status.outputs.network_security_group_id` |
| OciFileSystem | `spec.mountTarget.nsgIds` | `status.outputs.network_security_group_id` |
| OciFunctionsApplication | `spec.networkSecurityGroupIds` | `status.outputs.network_security_group_id` |
| OciMysqlDbSystem | `spec.nsgIds` | `status.outputs.network_security_group_id` |
| OciNetworkFirewall | `spec.networkSecurityGroupIds` | `status.outputs.network_security_group_id` |
| OciNetworkLoadBalancer | `spec.networkSecurityGroupIds` | `status.outputs.network_security_group_id` |
| OciPostgresqlDbSystem | `spec.networkDetails.nsgIds` | `status.outputs.network_security_group_id` |
| OciRedisCluster | `spec.nsgIds` | `status.outputs.network_security_group_id` |
| OciStreamPool | `spec.privateEndpointSettings.nsgIds` | `status.outputs.network_security_group_id` |

## See Also

- [Overview](../README.md)
