# AliCloudApplicationLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudApplicationLoadBalancerSpec defines the configuration for an Alibaba Cloud
Application Load Balancer (ALB) with bundled server groups and listeners.

ALB is a modern Layer 7 load balancer for HTTP, HTTPS, and QUIC traffic.
This component bundles the ALB, its server groups, and listeners into a
single deployable unit because an ALB without at least one server group
and listener is non-functional. Forwarding rules (alicloud_alb_rule) are
intentionally excluded -- they are extremely complex (9 action types x
9 condition types) and can be managed separately. The listener's
default_actions handle the 80% use case of forwarding traffic to a
server group.

Server groups are created as empty targets -- backend membership (ECS
instances, ENI IPs, etc.) is managed externally by ACK ingress controllers,
SAE bindings, or manual attachment, matching the Azure LoadBalancer pattern.

Billing is hardcoded to PayAsYouGo in the IaC modules because ALB does not
support subscription billing in the current provider.

Provider resources:
  Terraform: alicloud_alb_load_balancer + alicloud_alb_server_group + alicloud_alb_listener
  Pulumi:    alb.LoadBalancer + alb.ServerGroup + alb.Listener

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudApplicationLoadBalancer
metadata:
  name: test-alb
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
      healthCheckConfig:
        healthCheckEnabled: true
        healthCheckPath: /health
  listeners:
    - listenerPort: 80
      listenerProtocol: HTTP
      defaultActionServerGroupName: test-backend
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.loadBalancerName` | `string` |  |  |  |
| `spec.addressType` | `string` |  | `Internet` |  |
| `spec.loadBalancerEdition` | `string` |  | `Standard` |  |
| `spec.zoneMappings` | `[]AliCloudApplicationLoadBalancerZoneMapping` | yes |  |  |
| `spec.zoneMappings[].zoneId` | `string` | yes |  |  |
| `spec.zoneMappings[].vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.accessLogConfig` | `AliCloudApplicationLoadBalancerAccessLogConfig` |  |  |  |
| `spec.accessLogConfig.logProject` | `string` | yes |  |  |
| `spec.accessLogConfig.logStore` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.serverGroups` | `[]AliCloudApplicationLoadBalancerServerGroup` |  |  |  |
| `spec.serverGroups[].name` | `string` | yes |  |  |
| `spec.serverGroups[].protocol` | `string` |  | `HTTP` |  |
| `spec.serverGroups[].scheduler` | `string` |  | `Wrr` |  |
| `spec.serverGroups[].healthCheckConfig` | `AliCloudApplicationLoadBalancerHealthCheckConfig` | yes |  |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckEnabled` | `bool` |  |  |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckProtocol` | `string` |  | `HTTP` |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckPath` | `string` |  |  |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckHost` | `string` |  |  |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckMethod` | `string` |  | `HEAD` |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckConnectPort` | `int32` |  | `0` |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckInterval` | `int32` |  | `2` |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckTimeout` | `int32` |  | `5` |  |
| `spec.serverGroups[].healthCheckConfig.healthyThreshold` | `int32` |  | `3` |  |
| `spec.serverGroups[].healthCheckConfig.unhealthyThreshold` | `int32` |  | `3` |  |
| `spec.serverGroups[].healthCheckConfig.healthCheckCodes` | `[]string` |  |  |  |
| `spec.serverGroups[].stickySessionConfig` | `AliCloudApplicationLoadBalancerStickySessionConfig` |  |  |  |
| `spec.serverGroups[].stickySessionConfig.stickySessionEnabled` | `bool` |  |  |  |
| `spec.serverGroups[].stickySessionConfig.stickySessionType` | `string` |  |  |  |
| `spec.serverGroups[].stickySessionConfig.cookie` | `string` |  |  |  |
| `spec.serverGroups[].stickySessionConfig.cookieTimeout` | `int32` |  | `1000` |  |
| `spec.listeners` | `[]AliCloudApplicationLoadBalancerListener` |  |  |  |
| `spec.listeners[].listenerPort` | `int32` | yes |  |  |
| `spec.listeners[].listenerProtocol` | `string` | yes |  |  |
| `spec.listeners[].defaultActionServerGroupName` | `string` | yes |  |  |
| `spec.listeners[].listenerDescription` | `string` |  |  |  |
| `spec.listeners[].certificateId` | `string` |  |  |  |
| `spec.listeners[].securityPolicyId` | `string` |  |  |  |
| `spec.listeners[].gzipEnabled` | `bool` |  | `true` |  |
| `spec.listeners[].http2Enabled` | `bool` |  | `true` |  |
| `spec.listeners[].idleTimeout` | `int32` |  | `60` |  |
| `spec.listeners[].requestTimeout` | `int32` |  | `60` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the ALB will be created.
Must match the region of the VPC and VSwitches.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

VPC ID that the ALB belongs to. The ALB, all zone-mapping VSwitches,
and server groups must reside in the same VPC.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.loadBalancerName

`string`

ALB name. 2-128 characters; must start with a letter; can contain
digits, underscores, periods, and hyphens.
If omitted, defaults to the metadata.name.

- rule: load_balancer_name must be between 2 and 128 characters when set

### spec.addressType

`string` · optional (explicit presence)

Network type of the ALB.
"Internet" creates a public-facing ALB with a DNS name resolvable
from the internet. "Intranet" creates a VPC-internal ALB.
Default: "Internet"

- default: `Internet`
- rule: address_type must be one of: Internet, Intranet

### spec.loadBalancerEdition

`string` · optional (explicit presence)

ALB edition that determines feature availability and performance.
"Basic" supports basic L7 load balancing.
"Standard" adds advanced features (WAF integration, custom routing).
"StandardWithWaf" includes integrated WAF protection.
Default: "Standard"

- default: `Standard`
- rule: load_balancer_edition must be one of: Basic, Standard, StandardWithWaf

### spec.zoneMappings

`[]AliCloudApplicationLoadBalancerZoneMapping` · required

Availability zone mappings. Each mapping places an ALB node in a zone
with a VSwitch for IP allocation. ALB requires at least 2 zones for
high availability.

- rule: {"required":true,"repeated":{"minItems":"2"}}

### spec.zoneMappings[].zoneId

`string` · required

Availability zone ID within the region.
Examples: "cn-hangzhou-a", "cn-hangzhou-b", "us-west-1-a".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.zoneMappings[].vswitchId

`string | valueFrom` · required

VSwitch ID in this availability zone. The ALB allocates an IP from
this VSwitch. Must belong to the same VPC as the ALB's vpc_id.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the ALB is placed in the account's default resource group.

### spec.accessLogConfig

`AliCloudApplicationLoadBalancerAccessLogConfig`

Access log configuration for shipping ALB access logs to SLS (Log Service).
When configured, all listener access logs are sent to the specified
SLS log project and log store.

### spec.accessLogConfig.logProject

`string` · required

SLS log project name. Must already exist in the same region as the ALB.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.accessLogConfig.logStore

`string` · required

SLS log store name within the log project.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.tags

`map<string, string>`

Tags to apply to the ALB resource.

### spec.serverGroups

`[]AliCloudApplicationLoadBalancerServerGroup`

Server groups that define backend targets for the ALB.
Each server group has its own health check, protocol, and session
stickiness configuration. Listeners reference server groups by name.

Server groups are created empty -- backend membership is managed
externally (by ACK ingress, SAE, or manual attachment).

### spec.serverGroups[].name

`string` · required

Server group name. 2-128 characters; must start with a letter.
Used by listeners to reference this server group via
default_action_server_group_name.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.serverGroups[].protocol

`string` · optional (explicit presence)

Backend protocol for communication between the ALB and servers.
Default: "HTTP"

- default: `HTTP`
- rule: protocol must be one of: HTTP, HTTPS, GRPC

### spec.serverGroups[].scheduler

`string` · optional (explicit presence)

Scheduling algorithm for distributing requests across servers.
"Wrr" -- Weighted Round Robin (default, distributes by weight).
"Wlc" -- Weighted Least Connections (routes to least-busy server).
"Sch" -- Source IP Consistent Hashing (same client IP to same server).
Default: "Wrr"

- default: `Wrr`
- rule: scheduler must be one of: Wrr, Wlc, Sch

### spec.serverGroups[].healthCheckConfig

`AliCloudApplicationLoadBalancerHealthCheckConfig` · required

Health check configuration. Required for every server group.

- rule: {"required":true}

### spec.serverGroups[].healthCheckConfig.healthCheckEnabled

`bool`

Whether health checks are enabled for this server group.
When disabled, all servers are considered healthy.

### spec.serverGroups[].healthCheckConfig.healthCheckProtocol

`string` · optional (explicit presence)

Protocol used for health check probes.
Default: "HTTP"

- default: `HTTP`
- rule: health_check_protocol must be one of: HTTP, HTTPS, TCP, GRPC

### spec.serverGroups[].healthCheckConfig.healthCheckPath

`string`

URL path for HTTP/HTTPS health check probes. Must start with "/".
Ignored for TCP probes.
Example: "/health", "/api/healthz"

### spec.serverGroups[].healthCheckConfig.healthCheckHost

`string`

Domain name used in the Host header of HTTP/HTTPS health check probes.
If omitted, the server's IP address is used.

### spec.serverGroups[].healthCheckConfig.healthCheckMethod

`string` · optional (explicit presence)

HTTP method for health check probes.
Default: "HEAD"

- default: `HEAD`
- rule: health_check_method must be one of: GET, POST, HEAD

### spec.serverGroups[].healthCheckConfig.healthCheckConnectPort

`int32` · optional (explicit presence)

Backend port used for health checks. 0 means use the port of the
backend server (the default).
Range: 0-65535. Default: 0

- default: `0`
- rule: {"int32":{"lte":65535,"gte":0}}

### spec.serverGroups[].healthCheckConfig.healthCheckInterval

`int32` · optional (explicit presence)

Interval between health check probes, in seconds.
Range: 1-50. Default: 2

- default: `2`
- rule: {"int32":{"lte":50,"gte":1}}

### spec.serverGroups[].healthCheckConfig.healthCheckTimeout

`int32` · optional (explicit presence)

Maximum time to wait for a health check response, in seconds.
Range: 1-300. Default: 5

- default: `5`
- rule: {"int32":{"lte":300,"gte":1}}

### spec.serverGroups[].healthCheckConfig.healthyThreshold

`int32` · optional (explicit presence)

Number of consecutive successful probes before a server is marked healthy.
Range: 2-10. Default: 3

- default: `3`
- rule: {"int32":{"lte":10,"gte":2}}

### spec.serverGroups[].healthCheckConfig.unhealthyThreshold

`int32` · optional (explicit presence)

Number of consecutive failed probes before a server is marked unhealthy.
Range: 2-10. Default: 3

- default: `3`
- rule: {"int32":{"lte":10,"gte":2}}

### spec.serverGroups[].healthCheckConfig.healthCheckCodes

`[]string`

HTTP status codes that indicate a healthy response.
Only applicable when health_check_protocol is HTTP or HTTPS.
Examples: ["http_2xx"], ["http_2xx", "http_3xx"]

### spec.serverGroups[].stickySessionConfig

`AliCloudApplicationLoadBalancerStickySessionConfig`

Session stickiness configuration. When enabled, requests from the same
client are routed to the same backend server.

### spec.serverGroups[].stickySessionConfig.stickySessionEnabled

`bool`

Whether session stickiness is enabled.

### spec.serverGroups[].stickySessionConfig.stickySessionType

`string` · optional (explicit presence)

Session stickiness method.
"Insert" -- ALB inserts a cookie (SERVERID) into responses.
"Server" -- ALB uses a cookie set by the backend server.

- rule: sticky_session_type must be one of: Insert, Server

### spec.serverGroups[].stickySessionConfig.cookie

`string`

Cookie name when sticky_session_type is "Server".
The ALB reads this cookie from backend responses to identify the server.

### spec.serverGroups[].stickySessionConfig.cookieTimeout

`int32` · optional (explicit presence)

Cookie timeout in seconds when sticky_session_type is "Insert".
Range: 1-86400. Default: 1000

- default: `1000`
- rule: {"int32":{"lte":86400,"gte":1}}

### spec.listeners

`[]AliCloudApplicationLoadBalancerListener`

Listeners that define how the ALB accepts incoming traffic.
Each listener binds to a port and protocol and forwards traffic to a
server group via default_actions. HTTPS listeners require a certificate.

### spec.listeners[].listenerPort

`int32` · required

Port on which the listener accepts traffic.
Range: 1-65535. Common values: 80 (HTTP), 443 (HTTPS).

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.listeners[].listenerProtocol

`string` · required

Protocol for this listener.
"HTTP" for unencrypted HTTP traffic.
"HTTPS" for TLS-encrypted traffic (requires certificate_id).
"QUIC" for HTTP/3 over QUIC.

- rule: listener_protocol must be one of: HTTP, HTTPS, QUIC
- rule: {"required":true}

### spec.listeners[].defaultActionServerGroupName

`string` · required

Name of the server group that receives traffic from this listener.
Must match a name in the server_groups list. The listener creates a
default ForwardGroup action routing all traffic to this server group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.listeners[].listenerDescription

`string`

Human-readable description of this listener's purpose.
2-256 characters.

### spec.listeners[].certificateId

`string`

Server certificate ID for HTTPS and QUIC listeners.
Required when listener_protocol is "HTTPS" or "QUIC".
Obtain from Alibaba Cloud Certificate Management Service (CAS).

### spec.listeners[].securityPolicyId

`string`

TLS security policy that defines the supported TLS versions and cipher
suites. Only applicable for HTTPS and QUIC listeners.
Examples: "tls_cipher_policy_1_0", "tls_cipher_policy_1_2_strict"

### spec.listeners[].gzipEnabled

`bool` · optional (explicit presence)

Whether to enable gzip compression for HTTP responses.
Default: true

- default: `true`

### spec.listeners[].http2Enabled

`bool` · optional (explicit presence)

Whether to enable HTTP/2 for this listener.
Only applicable for HTTPS listeners.
Default: true

- default: `true`

### spec.listeners[].idleTimeout

`int32` · optional (explicit presence)

Connection idle timeout in seconds. Connections idle longer are closed.
Range: 1-60. Default: 60

- default: `60`
- rule: {"int32":{"lte":60,"gte":1}}

### spec.listeners[].requestTimeout

`int32` · optional (explicit presence)

Request timeout in seconds. The ALB returns 504 if the backend does
not respond within this period.
Range: 1-180. Default: 60

- default: `60`
- rule: {"int32":{"lte":180,"gte":1}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudApplicationLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | The ALB instance ID assigned by Alibaba Cloud (e.g., "alb-xxxxx"). |
| `status.outputs.dns_name` | `string` | The DNS name automatically assigned to the ALB. For internet-facing ALBs, this resolves to the ALB's public address. For intranet ALBs, it resolves to the private VPC address. Use this as a CNAME target for custom domain DNS records. |
| `status.outputs.server_group_ids` | `map<string, string>` | Map of server group names to their IDs. Keys are the names specified in spec.server_groups[].name. Values are the Alibaba Cloud server group IDs (e.g., "sgp-xxxxx"). Useful for downstream components that need to attach backends or create forwarding rules referencing specific server groups. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.zoneMappings[].vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
