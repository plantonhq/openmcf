# HetznerCloudFloatingIp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudFloatingIpSpec defines the specification for a Hetzner Cloud Floating IP.

A Floating IP is a reassignable public IP address (IPv4 or IPv6) that can be
moved between servers in the same location. Unlike auto-assigned public IPs,
Floating IPs persist independently of any server and can be reassigned at any
time, making them ideal for failover scenarios where a stable endpoint must
survive server replacement.

The IP address is allocated by Hetzner Cloud when the resource is created and
cannot be chosen. For IPv4, a single address is allocated. For IPv6, a /64
network block is allocated.

An optional reverse DNS (rDNS) record can be set for the allocated IP address.
rDNS maps an IP to a hostname and is commonly required for mail servers and
services that perform reverse lookups for identity verification.

An optional server assignment attaches the Floating IP to a server at creation
time. The assignment can be changed later without replacing the Floating IP.
The server must be in the same location as the Floating IP's home_location.

Fields not exposed in this spec (hardcoded or derived in IaC modules):
  - name:   Derived from metadata.name.
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudFloatingIp
metadata:
  name: hetznercloudfloatingip-demo
spec:
  type: ipv4
  homeLocation: fsn1
  description: Demo floating IP for failover testing
  dnsPtr: demo.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.type` | `enum` | yes |  |  |
| `spec.homeLocation` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.serverId` | `string \| valueFrom` |  |  | HetznerCloudServer (`status.outputs.server_id`) |
| `spec.dnsPtr` | `string` |  |  |  |
| `spec.deleteProtection` | `bool` |  |  |  |

## Field Details

### spec.type

`enum` · required

IP address type. Determines whether an IPv4 address or an IPv6 /64
network block is allocated.

Changing this value forces replacement of the Floating IP resource.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `ip_type_unspecified`
- `ipv4`
- `ipv6`

### spec.homeLocation

`string` · required

Hetzner Cloud location where the Floating IP is homed (e.g., "fsn1",
"nbg1", "hel1", "ash", "hil", "sin"). The Floating IP can only be
assigned to servers in the same location.

Changing this value forces replacement of the Floating IP resource.

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

Human-readable description for the Floating IP. Visible in the
Hetzner Cloud console and API. Useful for documenting the IP's
purpose (e.g., "Production web frontend failover IP").

### spec.serverId

`string | valueFrom`

Server to assign this Floating IP to. Optional.

Accepts a literal Hetzner Cloud server ID (as a string) or a reference
to a HetznerCloudServer resource's output via valueFrom. The referenced
server must be in the same location as home_location.

If omitted, the Floating IP is created unassigned (reserved). Assignment
can be added later by updating the spec.

Example (literal):
  serverId:
    value: "12345678"

Example (reference, once HetznerCloudServer is available):
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: my-server
      fieldPath: status.outputs.server_id

- references: HetznerCloudServer (`status.outputs.server_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudServer, name: <that resource's name>, fieldPath: status.outputs.server_id}} -- a bare string does not parse

### spec.dnsPtr

`string`

Reverse DNS pointer record for the allocated IP address. If set, an
hcloud_rdns resource is created that maps the IP to this hostname.

Commonly used for mail servers (SPF/DKIM verification relies on matching
forward and reverse DNS) and any service where clients verify the server's
identity via reverse lookup.

Example: "mail.example.com"

### spec.deleteProtection

`bool`

Prevent accidental deletion of the Floating IP via the Hetzner Cloud API.
When enabled, the IP cannot be deleted until protection is removed.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudFloatingIp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.floating_ip_id` | `string` | The Hetzner Cloud numeric ID of the created Floating IP (as a string). Can be referenced by other components via StringValueOrRef for cross-component wiring (e.g., monitoring or DNS record creation). |
| `status.outputs.ip_address` | `string` | The allocated IP address. For IPv4, this is a single address (e.g., "203.0.113.42"). For IPv6, this is the first address in the allocated /64 block (e.g., "2001:db8::1"). |
| `status.outputs.ip_network` | `string` | The allocated IPv6 network in CIDR notation (e.g., "2001:db8::/64"). Empty for IPv4 Floating IPs. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serverId` | HetznerCloudServer | `status.outputs.server_id` |

## See Also

- [Overview](../README.md)
