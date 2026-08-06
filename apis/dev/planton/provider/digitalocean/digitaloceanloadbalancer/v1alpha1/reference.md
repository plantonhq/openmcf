# DigitalOceanLoadBalancer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanLoadBalancerSpec defines the specification for creating a DigitalOcean Load Balancer.
It focuses on essential parameters following the 80/20 principle, including region, VPC placement,
target Droplets (by IDs or tag), forwarding rules for traffic, and health checks for backend service health.
Note: Either `droplet_ids` or `droplet_tag` may be provided (mutually exclusive). The load balancer must be associated with a VPC.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.loadBalancerName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.vpc` | `string \| valueFrom` | yes |  | DigitalOceanVpc (`status.outputs.vpc_id`) |
| `spec.forwardingRules` | `[]DigitalOceanLoadBalancerForwardingRule` | yes |  |  |
| `spec.forwardingRules[].entryPort` | `uint32` | yes |  |  |
| `spec.forwardingRules[].entryProtocol` | `enum` | yes |  |  |
| `spec.forwardingRules[].targetPort` | `uint32` | yes |  |  |
| `spec.forwardingRules[].targetProtocol` | `enum` | yes |  |  |
| `spec.forwardingRules[].certificateName` | `string` | yes |  |  |
| `spec.healthCheck` | `DigitalOceanLoadBalancerHealthCheck` |  |  |  |
| `spec.healthCheck.port` | `uint32` | yes |  |  |
| `spec.healthCheck.protocol` | `enum` | yes |  |  |
| `spec.healthCheck.path` | `string` |  |  |  |
| `spec.healthCheck.checkIntervalSec` | `uint32` |  | `10` |  |
| `spec.dropletIds` | `[]string \| valueFrom` |  |  | DigitalOceanDroplet (`status.outputs.droplet_id`) |
| `spec.dropletTag` | `string` | yes |  |  |
| `spec.enableStickySessions` | `bool` |  |  |  |

## Field Details

### spec.loadBalancerName

`string` · required

The name of the Load Balancer. Must be unique per account.
Constraints: 1-64 characters, lowercase alphanumeric and hyphens.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"64","pattern":"^[a-z0-9-]+$"}}

### spec.region

`enum` · required

The DigitalOcean region where the Load Balancer will be created.
Determines the geographical location of the load balancer.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

### spec.vpc

`string | valueFrom` · required

Reference to the DigitalOcean VPC in which to create the Load Balancer.
This should be a foreign key reference to an existing DigitalOceanVpc resource.

- references: DigitalOceanVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.forwardingRules

`[]DigitalOceanLoadBalancerForwardingRule` · required

A list of forwarding rules that define how traffic is routed from the load balancer to backend Droplets.
Each forwarding rule specifies an incoming port/protocol and a corresponding target port/protocol.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.forwardingRules[].entryPort

`uint32` · required

Port on the load balancer that will listen for incoming traffic.

- rule: {"required":true,"uint32":{"lte":65535,"gte":1}}

### spec.forwardingRules[].entryProtocol

`enum` · required

Protocol for incoming traffic on the load balancer's entry port (e.g., HTTP, HTTPS, TCP).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digitalocean_load_balancer_protocol_unspecified`
- `http`
- `https`
- `tcp`

### spec.forwardingRules[].targetPort

`uint32` · required

Port on the Droplet that will receive forwarded traffic.

- rule: {"required":true,"uint32":{"lte":65535,"gte":1}}

### spec.forwardingRules[].targetProtocol

`enum` · required

Protocol for traffic between the load balancer and the Droplet (e.g., HTTP, HTTPS, TCP).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digitalocean_load_balancer_protocol_unspecified`
- `http`
- `https`
- `tcp`

### spec.forwardingRules[].certificateName

`string` · required

The name of a TLS certificate resource uploaded to DigitalOcean.
Required when entry_protocol is HTTPS. The certificate is used for SSL termination.
Use certificate name (not ID) to avoid breaking IaC state when Let's Encrypt auto-renews certificates.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"minLen":"1","maxLen":"255"}}

### spec.healthCheck

`DigitalOceanLoadBalancerHealthCheck`

Health check configuration for the load balancer’s backend Droplets.
This defines how the load balancer will probe the Droplets to check their health.

### spec.healthCheck.port

`uint32` · required

The port on the Droplet to which the health check will be performed.

- rule: {"required":true,"uint32":{"lte":65535,"gte":1}}

### spec.healthCheck.protocol

`enum` · required

Protocol to use for health checking (HTTP, HTTPS, or TCP).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digitalocean_load_balancer_protocol_unspecified`
- `http`
- `https`
- `tcp`

### spec.healthCheck.path

`string`

If using HTTP/HTTPS for health checks, the request path to probe (e.g., "/health").
Ignored for TCP health checks.

### spec.healthCheck.checkIntervalSec

`uint32`

Interval (in seconds) between health check probes.

- default: `10`

### spec.dropletIds

`[]string | valueFrom`

A list of specific Droplet IDs to attach to the Load Balancer.
Mutually exclusive with `droplet_tag`. These can be literal IDs or references to Droplet resources.

- references: DigitalOceanDroplet (`status.outputs.droplet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanDroplet, name: <that resource's name>, fieldPath: status.outputs.droplet_id}} -- a bare string does not parse

### spec.dropletTag

`string` · required

A Droplet tag name. All Droplets with this tag in the specified VPC will be attached to the Load Balancer.
Mutually exclusive with `droplet_ids`.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"minLen":"1","maxLen":"255"}}

### spec.enableStickySessions

`bool`

Enables sticky sessions if true (disabled by default).
When enabled, the load balancer will attempt to direct repeated requests from the same client to the same Droplet.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanLoadBalancer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.load_balancer_id` | `string` | The unique identifier (UUID) of the created DigitalOcean Load Balancer. |
| `status.outputs.ip` | `string` | The public IP address assigned to the Load Balancer. |
| `status.outputs.dns_name` | `string` | The DNS name for the Load Balancer (if applicable). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpc` | DigitalOceanVpc | `status.outputs.vpc_id` |
| `spec.dropletIds` | DigitalOceanDroplet | `status.outputs.droplet_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
