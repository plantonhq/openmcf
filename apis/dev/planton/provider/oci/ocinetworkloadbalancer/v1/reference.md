# OciNetworkLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciNetworkLoadBalancerSpec defines the specification for an Oracle Cloud
Infrastructure Network Load Balancer (Layer 4).

An OCI network load balancer provides high-performance, low-latency traffic
distribution for TCP, UDP, and mixed-protocol workloads. Unlike the
application load balancer (OciApplicationLoadBalancer) which operates at Layer 7, the
network load balancer operates at Layer 4 and preserves the original source
IP address by default -- critical for firewalls, logging, and security
appliances that need to see the true client IP.

Key characteristics:
  - Fully elastic bandwidth (no shape configuration required)
  - Source IP preservation via is_preserve_source_destination
  - Tuple-based load balancing policies (FIVE_TUPLE, THREE_TUPLE, TWO_TUPLE)
  - Health checking via HTTP, HTTPS, TCP, UDP, and DNS protocols
  - Instant failover and fail-open capabilities
  - Single subnet deployment (unlike the L7 LB which supports multiple)

This component bundles the network load balancer with its backend sets,
backends, and listeners into a single deployment unit.

For Layer 7 (HTTP/HTTPS) load balancing with SSL termination, virtual
hostname routing, and rule-based traffic manipulation, use OciApplicationLoadBalancer.

## Example

```yaml
apiVersion: oci.planton.dev/v1
kind: OciNetworkLoadBalancer
metadata:
  name: ocinetworkloadbalancer-demo
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  subnetId:
    value: "ocid1.subnet.oc1.iad.example"
  isPrivate: true
  isPreserveSourceDestination: true
  backendSets:
    - name: tcp-backend
      policy: five_tuple
      healthChecker:
        protocol: tcp
        port: 8080
      backends:
        - ipAddress: "10.0.1.10"
          port: 8080
          weight: 1
        - ipAddress: "10.0.1.11"
          port: 8080
          weight: 1
  listeners:
    - name: tcp-listener
      port: 80
      protocol: tcp
      defaultBackendSetName: tcp-backend
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.isPrivate` | `bool` |  |  |  |
| `spec.isPreserveSourceDestination` | `bool` |  |  |  |
| `spec.isSymmetricHashEnabled` | `bool` |  |  |  |
| `spec.networkSecurityGroupIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.nlbIpVersion` | `string` |  |  |  |
| `spec.reservedIps` | `[]ReservedIp` |  |  |  |
| `spec.reservedIps[].id` | `string` | yes |  |  |
| `spec.backendSets` | `[]BackendSet` | yes |  |  |
| `spec.backendSets[].name` | `string` | yes |  |  |
| `spec.backendSets[].policy` | `enum` |  |  |  |
| `spec.backendSets[].healthChecker` | `HealthChecker` | yes |  |  |
| `spec.backendSets[].healthChecker.protocol` | `enum` |  |  |  |
| `spec.backendSets[].healthChecker.port` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.urlPath` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.returnCode` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.responseBodyRegex` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.intervalInMillis` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.timeoutInMillis` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.retries` | `int32` |  |  |  |
| `spec.backendSets[].healthChecker.requestData` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.responseData` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.dnsHealthCheck` | `DnsHealthCheck` |  |  |  |
| `spec.backendSets[].healthChecker.dnsHealthCheck.domainName` | `string` | yes |  |  |
| `spec.backendSets[].healthChecker.dnsHealthCheck.queryClass` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.dnsHealthCheck.queryType` | `string` |  |  |  |
| `spec.backendSets[].healthChecker.dnsHealthCheck.rcodes` | `[]string` |  |  |  |
| `spec.backendSets[].healthChecker.dnsHealthCheck.transportProtocol` | `string` |  |  |  |
| `spec.backendSets[].backends` | `[]Backend` |  |  |  |
| `spec.backendSets[].backends[].port` | `int32` |  |  |  |
| `spec.backendSets[].backends[].ipAddress` | `string` |  |  |  |
| `spec.backendSets[].backends[].targetId` | `string` |  |  |  |
| `spec.backendSets[].backends[].weight` | `int32` |  |  |  |
| `spec.backendSets[].backends[].isBackup` | `bool` |  |  |  |
| `spec.backendSets[].backends[].isDrain` | `bool` |  |  |  |
| `spec.backendSets[].backends[].isOffline` | `bool` |  |  |  |
| `spec.backendSets[].backends[].name` | `string` |  |  |  |
| `spec.backendSets[].isPreserveSource` | `bool` |  |  |  |
| `spec.backendSets[].isFailOpen` | `bool` |  |  |  |
| `spec.backendSets[].isInstantFailoverEnabled` | `bool` |  |  |  |
| `spec.backendSets[].isInstantFailoverTcpResetEnabled` | `bool` |  |  |  |
| `spec.backendSets[].areOperationallyActiveBackendsPreferred` | `bool` |  |  |  |
| `spec.backendSets[].ipVersion` | `string` |  |  |  |
| `spec.listeners` | `[]Listener` | yes |  |  |
| `spec.listeners[].name` | `string` | yes |  |  |
| `spec.listeners[].port` | `int32` |  |  |  |
| `spec.listeners[].protocol` | `enum` |  |  |  |
| `spec.listeners[].defaultBackendSetName` | `string` | yes |  |  |
| `spec.listeners[].ipVersion` | `string` |  |  |  |
| `spec.listeners[].isPpv2Enabled` | `bool` |  |  |  |
| `spec.listeners[].tcpIdleTimeout` | `int32` |  |  |  |
| `spec.listeners[].udpIdleTimeout` | `int32` |  |  |  |
| `spec.listeners[].l3ipIdleTimeout` | `int32` |  |  |  |
| `spec.assignedIpv6` | `string` |  |  |  |
| `spec.assignedPrivateIpv4` | `string` |  |  |  |
| `spec.subnetIpv6cidr` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the network load balancer will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name for the network load balancer shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.subnetId

`string | valueFrom` · required

OCID of the subnet where the network load balancer will be spawned.
Unlike the L7 load balancer which accepts multiple subnets, the NLB
is deployed into a single subnet. Changing this after creation forces
recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.isPrivate

`bool`

When true, creates a private network load balancer that is not accessible
from the public internet. Private NLBs receive only private IP addresses.
Note: OCI defaults NLBs to private (true). The proto3 default is false,
so explicitly set this to true for private NLBs.
Changing this after creation forces recreation.

### spec.isPreserveSourceDestination

`bool`

When true, the NLB preserves the source IP address and destination IP
address in the IP header. This automatically enables skipSourceDestinationCheck
on the NLB's VNIC, allowing packets to reach backends with the original
client IP intact. Essential for firewalls, intrusion detection systems,
and applications that need the true client IP.

### spec.isSymmetricHashEnabled

`bool`

When true, enables symmetric hashing for the NLB. This can only be enabled
when the NLB is working in transparent mode with source/destination
preservation enabled. Removes the dependency on backends (like firewalls)
to perform SNAT.

### spec.networkSecurityGroupIds

`[]string | valueFrom`

OCIDs of network security groups applied to the network load balancer.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.nlbIpVersion

`string`

IP version for the network load balancer. Accepted values:
"IPV4", "IPV6", "IPV4_AND_IPV6".
When omitted, defaults to "IPV4".

### spec.reservedIps

`[]ReservedIp`

Pre-created reserved public IPs to assign to the network load balancer.
When omitted, OCI assigns ephemeral public IPs (for public NLBs).

### spec.reservedIps[].id

`string` · required

OCID of the reserved public IP to assign to the network load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.backendSets

`[]BackendSet` · required

Backend sets define groups of backend servers with load balancing
policies and health checking. At least one backend set is required.
Each listener routes traffic to exactly one default backend set.

- rule: {"repeated":{"minItems":"1"}}

### spec.backendSets[].name

`string` · required

Unique name for this backend set within the network load balancer.
Listeners reference backend sets by this name.

- rule: {"string":{"minLen":"1"}}

### spec.backendSets[].policy

`enum`

Load balancing policy that determines how traffic is distributed
across backends. NLB uses tuple-based hashing rather than the
round-robin/least-connections policies of the L7 load balancer.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `policy_unspecified`
- `five_tuple`
- `three_tuple`
- `two_tuple`

### spec.backendSets[].healthChecker

`HealthChecker` · required

Health checker configuration that monitors backend availability.
The NLB removes unhealthy backends from the rotation until they
pass health checks again.

- rule: {"required":true}

### spec.backendSets[].healthChecker.protocol

`enum`

Protocol used for health checks.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `http`
- `https`
- `tcp`
- `udp`
- `dns`

### spec.backendSets[].healthChecker.port

`int32`

Port on the backend server to probe. When omitted or set to 0,
the health checker uses the backend's traffic port.

### spec.backendSets[].healthChecker.urlPath

`string`

URL path for HTTP/HTTPS health checks (e.g., "/health", "/ready").
Required when protocol is http or https. Ignored for other protocols.

### spec.backendSets[].healthChecker.returnCode

`int32`

Expected HTTP status code from healthy backends (e.g., 200).
Used with http and https protocols.

### spec.backendSets[].healthChecker.responseBodyRegex

`string`

Regex pattern to match against the response body. The backend is
considered healthy only if the response body matches this pattern.

### spec.backendSets[].healthChecker.intervalInMillis

`int32`

Interval between consecutive health checks in milliseconds.
When omitted, defaults to 10000 (10 seconds).

### spec.backendSets[].healthChecker.timeoutInMillis

`int32`

Maximum time to wait for a health check response in milliseconds.
When omitted, defaults to 3000 (3 seconds).

### spec.backendSets[].healthChecker.retries

`int32`

Number of consecutive failed health checks before marking a backend
as unhealthy. When omitted, defaults to 3.

### spec.backendSets[].healthChecker.requestData

`string`

Base64-encoded request data for UDP and TCP health checks.
The NLB sends this data as the health check probe payload.

### spec.backendSets[].healthChecker.responseData

`string`

Base64-encoded expected response data for UDP and TCP health checks.
The backend is healthy only if the response matches this data.

### spec.backendSets[].healthChecker.dnsHealthCheck

`DnsHealthCheck`

DNS health check configuration. Required when protocol is dns.
Ignored for other protocols.

### spec.backendSets[].healthChecker.dnsHealthCheck.domainName

`string` · required

Fully qualified domain name to query.

- rule: {"string":{"minLen":"1"}}

### spec.backendSets[].healthChecker.dnsHealthCheck.queryClass

`string`

DNS query class. Accepted values: "IN" (Internet), "CH" (Chaos).
When omitted, defaults to "IN".

### spec.backendSets[].healthChecker.dnsHealthCheck.queryType

`string`

DNS query type. Accepted values: "A", "AAAA", "TXT".
When omitted, defaults to "A".

### spec.backendSets[].healthChecker.dnsHealthCheck.rcodes

`[]string`

Acceptable DNS response codes. The backend is healthy if the
response code matches any value in this list.
Example values: "NOERROR", "NXDOMAIN".
When omitted, defaults to ["NOERROR"].

### spec.backendSets[].healthChecker.dnsHealthCheck.transportProtocol

`string`

Transport protocol for DNS queries. Accepted values: "UDP", "TCP".
When omitted, defaults to "UDP".

### spec.backendSets[].backends

`[]Backend`

Backend servers in this set. When omitted, the backend set is
created without backends (useful when backends are added dynamically
or via target_id references to compute instances).

### spec.backendSets[].backends[].port

`int32`

Port on which the backend server listens for traffic.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.backendSets[].backends[].ipAddress

`string`

IP address of the backend server. When omitted, target_id must be
provided instead.

### spec.backendSets[].backends[].targetId

`string`

OCID of a compute instance or private IP to use as the backend target.
When provided, OCI resolves the IP address automatically. When omitted,
ip_address must be provided instead.

### spec.backendSets[].backends[].weight

`int32`

Relative weight for traffic distribution. Higher weights receive
proportionally more traffic. When omitted, defaults to 1.

### spec.backendSets[].backends[].isBackup

`bool`

When true, this backend only receives traffic when all non-backup
backends are unhealthy.

### spec.backendSets[].backends[].isDrain

`bool`

When true, the backend is in drain mode -- existing connections
complete but no new connections are sent to this backend.

### spec.backendSets[].backends[].isOffline

`bool`

When true, the backend is temporarily taken offline. No traffic
is sent to offline backends.

### spec.backendSets[].backends[].name

`string`

Optional unique name identifying the backend within the backend set.
When omitted, OCI auto-generates a name in the format "IP:port".

### spec.backendSets[].isPreserveSource

`bool`

When true, the NLB preserves the source IP of the packet when
forwarding to backends in this set. Defaults to true in OCI.

### spec.backendSets[].isFailOpen

`bool`

When true, the NLB continues distributing traffic to all backends
even when all backends are marked unhealthy. Prevents total service
outage at the cost of sending traffic to degraded backends.

### spec.backendSets[].isInstantFailoverEnabled

`bool`

When true, existing connections are immediately forwarded to an
alternative healthy backend as soon as the current backend becomes
unhealthy, rather than waiting for the connection to time out.

### spec.backendSets[].isInstantFailoverTcpResetEnabled

`bool`

When true (and instant failover is enabled), the NLB sends a TCP
RST to clients for existing connections instead of silently failing
over. Gives clients an immediate signal to reconnect.

### spec.backendSets[].areOperationallyActiveBackendsPreferred

`bool`

When true, enables active-standby backend support. The NLB
preferentially routes traffic to operationally active backends.

### spec.backendSets[].ipVersion

`string`

IP version for this backend set. When omitted, inherits from the
NLB's nlb_ip_version setting.

### spec.listeners

`[]Listener` · required

Listeners define the ports and protocols on which the network load
balancer accepts connections. At least one listener is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.listeners[].name

`string` · required

Unique name for this listener within the network load balancer.

- rule: {"string":{"minLen":"1"}}

### spec.listeners[].port

`int32`

Port on which the listener accepts connections.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.listeners[].protocol

`enum`

Protocol for the listener. NLB supports Layer 4 protocols only.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `tcp`
- `udp`
- `tcp_and_udp`
- `any`

### spec.listeners[].defaultBackendSetName

`string` · required

Name of the backend set that receives traffic from this listener.
Must match a backend set defined in backend_sets.

- rule: {"string":{"minLen":"1"}}

### spec.listeners[].ipVersion

`string`

IP version for this listener. When omitted, inherits from the
NLB's nlb_ip_version setting.

### spec.listeners[].isPpv2Enabled

`bool`

When true, enables Proxy Protocol v2 (PPv2) on this listener.
PPv2 prepends connection metadata (source IP, port, protocol)
to the TCP stream, allowing backends to see the original client
information even when source IP preservation is not possible.

### spec.listeners[].tcpIdleTimeout

`int32`

TCP idle timeout in seconds. Connections idle longer than this
are closed by the NLB. When omitted, OCI uses its default.

### spec.listeners[].udpIdleTimeout

`int32`

UDP idle timeout in seconds. When omitted, OCI uses its default.

### spec.listeners[].l3ipIdleTimeout

`int32`

L3IP idle timeout in seconds. Applies when the protocol is ANY.
When omitted, OCI uses its default.

### spec.assignedIpv6

`string`

IPv6 address to assign to the network load balancer. Must be part of
one of the prefixes supported by the subnet.
Example: "2607:9b80:9a0a:9a7e:abcd:ef01:2345:6789"

### spec.assignedPrivateIpv4

`string`

Private IPv4 address to assign to the network load balancer. Must be
in the CIDR range of the subnet. Changing this after creation forces
recreation. Example: "10.0.0.1"

### spec.subnetIpv6cidr

`string`

IPv6 subnet prefix selection. When provided, the NLB IPv6 address is
assigned within this CIDR block.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciNetworkLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.network_load_balancer_id` | `string` | OCID of the created network load balancer. |
| `status.outputs.ip_addresses` | `string` | Comma-separated IP addresses assigned to the network load balancer. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.networkSecurityGroupIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
