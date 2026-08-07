# AliCloudNetworkLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudNetworkLoadBalancerSpec defines the configuration for an Alibaba Cloud
Network Load Balancer (NLB) with bundled server groups and listeners.

NLB is a modern Layer 4 load balancer for TCP, UDP, and TCPSSL traffic,
designed for ultra-high performance and low latency. This component bundles
the NLB, its server groups, and listeners into a single deployable unit
because an NLB without at least one server group and listener is
non-functional.

Server groups are created as empty targets -- backend membership (ECS
instances, ENI IPs, etc.) is managed externally by ACK service controllers,
manual attachment, or other orchestration, matching the ALB pattern.

Billing is hardcoded to PayAsYouGo in the IaC modules, matching the ALB
convention.

Provider resources:
  Terraform: alicloud_nlb_load_balancer + alicloud_nlb_server_group + alicloud_nlb_listener
  Pulumi:    nlb.LoadBalancer + nlb.ServerGroup + nlb.Listener

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudNetworkLoadBalancer
metadata:
  name: test-nlb
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-test123
  zoneMappings:
    - zoneId: cn-hangzhou-a
      vswitchId:
        value: vsw-test-a
    - zoneId: cn-hangzhou-b
      vswitchId:
        value: vsw-test-b
  serverGroups:
    - name: test-backend
      healthCheck:
        healthCheckEnabled: true
  listeners:
    - listenerPort: 80
      listenerProtocol: TCP
      serverGroupName: test-backend
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.loadBalancerName` | `string` |  |  |  |
| `spec.addressType` | `string` |  | `Internet` |  |
| `spec.zoneMappings` | `[]AliCloudNetworkLoadBalancerZoneMapping` | yes |  |  |
| `spec.zoneMappings[].zoneId` | `string` | yes |  |  |
| `spec.zoneMappings[].vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.zoneMappings[].allocationId` | `string \| valueFrom` |  |  | AliCloudEipAddress (`status.outputs.eip_id`) |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.crossZoneEnabled` | `bool` |  | `true` |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.serverGroups` | `[]AliCloudNetworkLoadBalancerServerGroup` |  |  |  |
| `spec.serverGroups[].name` | `string` | yes |  |  |
| `spec.serverGroups[].protocol` | `string` |  | `TCP` |  |
| `spec.serverGroups[].scheduler` | `string` |  | `Wrr` |  |
| `spec.serverGroups[].connectionDrainEnabled` | `bool` |  | `false` |  |
| `spec.serverGroups[].connectionDrainTimeout` | `int32` |  | `10` |  |
| `spec.serverGroups[].preserveClientIpEnabled` | `bool` |  | `true` |  |
| `spec.serverGroups[].healthCheck` | `AliCloudNetworkLoadBalancerHealthCheckConfig` | yes |  |  |
| `spec.serverGroups[].healthCheck.healthCheckEnabled` | `bool` |  |  |  |
| `spec.serverGroups[].healthCheck.healthCheckType` | `string` |  | `TCP` |  |
| `spec.serverGroups[].healthCheck.healthCheckConnectPort` | `int32` |  | `0` |  |
| `spec.serverGroups[].healthCheck.healthCheckConnectTimeout` | `int32` |  | `5` |  |
| `spec.serverGroups[].healthCheck.healthCheckInterval` | `int32` |  | `10` |  |
| `spec.serverGroups[].healthCheck.healthyThreshold` | `int32` |  | `2` |  |
| `spec.serverGroups[].healthCheck.unhealthyThreshold` | `int32` |  | `2` |  |
| `spec.serverGroups[].healthCheck.healthCheckUrl` | `string` |  |  |  |
| `spec.serverGroups[].healthCheck.healthCheckDomain` | `string` |  |  |  |
| `spec.serverGroups[].healthCheck.httpCheckMethod` | `string` |  | `GET` |  |
| `spec.serverGroups[].healthCheck.healthCheckHttpCodes` | `[]string` |  |  |  |
| `spec.listeners` | `[]AliCloudNetworkLoadBalancerListener` |  |  |  |
| `spec.listeners[].listenerPort` | `int32` | yes |  |  |
| `spec.listeners[].listenerProtocol` | `string` | yes |  |  |
| `spec.listeners[].serverGroupName` | `string` | yes |  |  |
| `spec.listeners[].listenerDescription` | `string` |  |  |  |
| `spec.listeners[].idleTimeout` | `int32` |  | `900` |  |
| `spec.listeners[].proxyProtocolEnabled` | `bool` |  | `false` |  |
| `spec.listeners[].certificateIds` | `[]string` |  |  |  |
| `spec.listeners[].securityPolicyId` | `string` |  |  |  |
| `spec.listeners[].caCertificateIds` | `[]string` |  |  |  |
| `spec.listeners[].caEnabled` | `bool` |  | `false` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the NLB will be created.
Must match the region of the VPC and VSwitches.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

VPC ID that the NLB belongs to. The NLB, all zone-mapping VSwitches,
and server groups must reside in the same VPC.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.loadBalancerName

`string`

NLB name. 2-128 characters; must start with a letter; can contain
digits, underscores, periods, and hyphens.
If omitted, defaults to the metadata.name.

- rule: load_balancer_name must be between 2 and 128 characters when set

### spec.addressType

`string` · optional (explicit presence)

Network type of the NLB.
"Internet" creates a public-facing NLB with a DNS name resolvable
from the internet. "Intranet" creates a VPC-internal NLB.
Default: "Internet"

- default: `Internet`
- rule: address_type must be one of: Internet, Intranet

### spec.zoneMappings

`[]AliCloudNetworkLoadBalancerZoneMapping` · required

Availability zone mappings. Each mapping places an NLB node in a zone
with a VSwitch for IP allocation. NLB requires at least 2 zones for
high availability.

- rule: {"required":true,"repeated":{"minItems":"2"}}

### spec.zoneMappings[].zoneId

`string` · required

Availability zone ID within the region.
Examples: "cn-hangzhou-a", "cn-hangzhou-b", "us-west-1-a".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.zoneMappings[].vswitchId

`string | valueFrom` · required

VSwitch ID in this availability zone. The NLB allocates an IP from
this VSwitch. Must belong to the same VPC as the NLB's vpc_id.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.zoneMappings[].allocationId

`string | valueFrom`

EIP allocation ID to bind a fixed public IP to this zone's NLB node.
Only meaningful for internet-facing NLBs that need stable public IPs
(e.g., database access, game servers, DNS A-records).
If omitted, Alibaba Cloud auto-assigns IPs.

- references: AliCloudEipAddress (`status.outputs.eip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudEipAddress, name: <that resource's name>, fieldPath: status.outputs.eip_id}} -- a bare string does not parse

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the NLB is placed in the account's default resource group.

### spec.crossZoneEnabled

`bool` · optional (explicit presence)

Whether to enable cross-zone load balancing. When enabled, NLB
distributes traffic across all healthy backends in all zones. When
disabled, traffic stays within the zone where it was received.
Default: true

- default: `true`

### spec.tags

`map<string, string>`

Tags to apply to the NLB resource.

### spec.serverGroups

`[]AliCloudNetworkLoadBalancerServerGroup`

Server groups that define backend targets for the NLB.
Each server group has its own health check, protocol, and scheduling
configuration. Listeners reference server groups by name.

Server groups are created empty -- backend membership is managed
externally (by ACK service controllers, manual attachment, etc.).

### spec.serverGroups[].name

`string` · required

Server group name. 2-128 characters; must start with a letter.
Used by listeners to reference this server group via server_group_name.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.serverGroups[].protocol

`string` · optional (explicit presence)

Backend protocol for communication between the NLB and servers.
Default: "TCP"

- default: `TCP`
- rule: protocol must be one of: TCP, UDP, TCPSSL

### spec.serverGroups[].scheduler

`string` · optional (explicit presence)

Scheduling algorithm for distributing connections across servers.
"Wrr" -- Weighted Round Robin (default, distributes by weight).
"Rr"  -- Round Robin (equal distribution).
"Sch" -- Source IP Consistent Hashing (same client to same server).
"Tch" -- Four-tuple Consistent Hashing (src/dst IP + port).
"Qch" -- QUIC ID Consistent Hashing.
"Wlc" -- Weighted Least Connections.
Default: "Wrr"

- default: `Wrr`
- rule: scheduler must be one of: Wrr, Rr, Sch, Tch, Qch, Wlc

### spec.serverGroups[].connectionDrainEnabled

`bool` · optional (explicit presence)

Whether to enable connection draining. When enabled, existing
connections to a removed backend are allowed to complete within the
drain timeout before being forcibly closed.
Default: false

- default: `false`

### spec.serverGroups[].connectionDrainTimeout

`int32` · optional (explicit presence)

Maximum time in seconds to wait for in-flight connections to complete
when a backend is removed. Only effective when connection_drain_enabled
is true.
Range: 10-900. Default: 10

- default: `10`
- rule: {"int32":{"lte":900,"gte":10}}

### spec.serverGroups[].preserveClientIpEnabled

`bool` · optional (explicit presence)

Whether to preserve the client's original IP address when forwarding
traffic to backends. The backend sees the real client IP instead of
the NLB's IP.
Default: true

- default: `true`

### spec.serverGroups[].healthCheck

`AliCloudNetworkLoadBalancerHealthCheckConfig` · required

Health check configuration. Required for every server group.

- rule: {"required":true}

### spec.serverGroups[].healthCheck.healthCheckEnabled

`bool`

Whether health checks are enabled for this server group.
When disabled, all servers are considered healthy.

### spec.serverGroups[].healthCheck.healthCheckType

`string` · optional (explicit presence)

Protocol used for health check probes.
"TCP" performs a TCP connection check.
"HTTP" sends an HTTP request and checks the response code.
"UDP" sends a UDP packet and expects a response.
Default: "TCP"

- default: `TCP`
- rule: health_check_type must be one of: TCP, HTTP, UDP

### spec.serverGroups[].healthCheck.healthCheckConnectPort

`int32` · optional (explicit presence)

Backend port used for health checks. 0 means use the port of the
backend server (the default).
Range: 0-65535. Default: 0

- default: `0`
- rule: {"int32":{"lte":65535,"gte":0}}

### spec.serverGroups[].healthCheck.healthCheckConnectTimeout

`int32` · optional (explicit presence)

Maximum time to wait for a health check response, in seconds.
Range: 1-300. Default: 5

- default: `5`
- rule: {"int32":{"lte":300,"gte":1}}

### spec.serverGroups[].healthCheck.healthCheckInterval

`int32` · optional (explicit presence)

Interval between health check probes, in seconds.
Range: 5-50. Default: 10

- default: `10`
- rule: {"int32":{"lte":50,"gte":5}}

### spec.serverGroups[].healthCheck.healthyThreshold

`int32` · optional (explicit presence)

Number of consecutive successful probes before a server is marked healthy.
Range: 2-10. Default: 2

- default: `2`
- rule: {"int32":{"lte":10,"gte":2}}

### spec.serverGroups[].healthCheck.unhealthyThreshold

`int32` · optional (explicit presence)

Number of consecutive failed probes before a server is marked unhealthy.
Range: 2-10. Default: 2

- default: `2`
- rule: {"int32":{"lte":10,"gte":2}}

### spec.serverGroups[].healthCheck.healthCheckUrl

`string`

URL path for HTTP health check probes. Must start with "/".
Only applicable when health_check_type is "HTTP".
Example: "/health", "/api/healthz"

### spec.serverGroups[].healthCheck.healthCheckDomain

`string`

Domain name used in the Host header of HTTP health check probes.
Only applicable when health_check_type is "HTTP".
If omitted, the server's IP address is used.

### spec.serverGroups[].healthCheck.httpCheckMethod

`string` · optional (explicit presence)

HTTP method for health check probes.
Only applicable when health_check_type is "HTTP".
Default: "GET"

- default: `GET`
- rule: http_check_method must be one of: GET, HEAD

### spec.serverGroups[].healthCheck.healthCheckHttpCodes

`[]string`

HTTP status codes that indicate a healthy response.
Only applicable when health_check_type is "HTTP".
Examples: ["http_2xx"], ["http_2xx", "http_3xx"]

### spec.listeners

`[]AliCloudNetworkLoadBalancerListener`

Listeners that define how the NLB accepts incoming traffic.
Each listener binds to a port and protocol and forwards traffic
directly to a server group. TCPSSL listeners support TLS termination.

### spec.listeners[].listenerPort

`int32` · required

Port on which the listener accepts traffic.
Range: 1-65535. Common values: 80 (TCP), 443 (TCPSSL), 3306 (MySQL).

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.listeners[].listenerProtocol

`string` · required

Protocol for this listener.
"TCP" for raw TCP traffic.
"UDP" for UDP traffic.
"TCPSSL" for TLS-encrypted TCP traffic (requires certificate_ids).

- rule: listener_protocol must be one of: TCP, UDP, TCPSSL
- rule: {"required":true}

### spec.listeners[].serverGroupName

`string` · required

Name of the server group that receives traffic from this listener.
Must match a name in the server_groups list.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.listeners[].listenerDescription

`string`

Human-readable description of this listener's purpose.

### spec.listeners[].idleTimeout

`int32` · optional (explicit presence)

Connection idle timeout in seconds. Connections idle longer are closed.
Applies to TCP and TCPSSL only; ignored for UDP.
Range: 1-900. Default: 900

- default: `900`
- rule: {"int32":{"lte":900,"gte":1}}

### spec.listeners[].proxyProtocolEnabled

`bool` · optional (explicit presence)

Whether to enable the Proxy Protocol. When enabled, the NLB inserts
a Proxy Protocol header so backends can read the real client IP/port.
Default: false

- default: `false`

### spec.listeners[].certificateIds

`[]string`

Server certificate IDs for TCPSSL listeners. At least one is required
when listener_protocol is "TCPSSL".
Obtain from Alibaba Cloud Certificate Management Service (CAS).

### spec.listeners[].securityPolicyId

`string`

TLS security policy that defines the supported TLS versions and cipher
suites. Only applicable for TCPSSL listeners.
Examples: "tls_cipher_policy_1_0", "tls_cipher_policy_1_2_strict"

### spec.listeners[].caCertificateIds

`[]string`

CA certificate IDs for mutual TLS authentication on TCPSSL listeners.
When set along with ca_enabled, the NLB validates the client certificate
against these CA certificates.

### spec.listeners[].caEnabled

`bool` · optional (explicit presence)

Whether to enable mutual TLS (client certificate verification)
on TCPSSL listeners. Requires ca_certificate_ids to be set.
Default: false

- default: `false`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudNetworkLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | The NLB instance ID assigned by Alibaba Cloud (e.g., "nlb-xxxxx"). |
| `status.outputs.dns_name` | `string` | The DNS name automatically assigned to the NLB. For internet-facing NLBs, this resolves to the NLB's public address. For intranet NLBs, it resolves to the private VPC address. Use this as a CNAME target for custom domain DNS records. |
| `status.outputs.server_group_ids` | `map<string, string>` | Map of server group names to their IDs. Keys are the names specified in spec.server_groups[].name. Values are the Alibaba Cloud server group IDs (e.g., "sgp-xxxxx"). Useful for downstream components that need to attach backends or reference specific server groups. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.zoneMappings[].vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.zoneMappings[].allocationId` | AliCloudEipAddress | `status.outputs.eip_id` |

## See Also

- [Overview](../README.md)
