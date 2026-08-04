# HetznerCloudPrimaryIp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1`

HetznerCloudPrimaryIpSpec defines the specification for a Hetzner Cloud Primary IP.

A Primary IP is a managed public IP address (IPv4 or IPv6) that persists
independently of any server. Primary IPs are created at a specific location
and can be assigned to servers at creation time or later. Unlike auto-assigned
public IPs, Primary IPs survive server deletion and can be reassigned, making
them suitable for stable endpoints that must outlive individual servers.

The IP address is allocated by Hetzner Cloud when the resource is created and
cannot be chosen. For IPv4, a single address is allocated. For IPv6, a /64
network block is allocated.

An optional reverse DNS (rDNS) record can be set for the allocated IP address.
rDNS maps an IP to a hostname and is commonly required for mail servers and
services that perform reverse lookups for identity verification.

The Primary IP is referenced by HetznerCloudServer via its primary_ip_id
output through StringValueOrRef.

Fields not exposed in this spec (hardcoded in IaC modules):
  - auto_delete:   Always false. In Planton's component model, resources are
                   managed independently. Auto-deletion on server removal would
                   silently destroy a separately-managed resource.
  - assignee_type: Always "server" (the only type Hetzner Cloud supports).
  - assignee_id:   Not set. Server assignment is handled by the
                   HetznerCloudServer component, not by the IP component.
  - datacenter:    Deprecated in the provider (removal after 2026-07-01).
                   Use location instead.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudPrimaryIp
metadata:
  name: hetznercloudprimaryip-demo
spec:
  type: ipv4
  location: fsn1
  dnsPtr: mail.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.type` | `enum` | yes |  |  |
| `spec.location` | `string` | yes |  |  |
| `spec.dnsPtr` | `string` |  |  |  |
| `spec.deleteProtection` | `bool` |  |  |  |

## Field Details

### spec.type

`enum` · required

IP address type. Determines whether an IPv4 address or an IPv6 /64
network block is allocated.

Changing this value forces replacement of the Primary IP resource.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `ip_type_unspecified`
- `ipv4`
- `ipv6`

### spec.location

`string` · required

Hetzner Cloud location where the IP is allocated (e.g., "fsn1", "nbg1",
"hel1", "ash", "hil", "sin"). The Primary IP can only be assigned to
servers in the same location.

Changing this value forces replacement of the Primary IP resource.

- rule: {"string":{"minLen":"1"}}

### spec.dnsPtr

`string`

Reverse DNS pointer record for the allocated IP address. If set, an
hcloud_rdns resource is created that maps the IP back to this hostname.

Commonly used for mail servers (SPF/DKIM verification relies on matching
forward and reverse DNS) and any service where clients verify the server's
identity via reverse lookup.

Example: "mail.example.com"

### spec.deleteProtection

`bool`

Prevent accidental deletion of the Primary IP via the Hetzner Cloud API.
When enabled, the IP cannot be deleted until protection is removed.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudPrimaryIp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.primary_ip_id` | `string` | The Hetzner Cloud numeric ID of the created Primary IP (as a string). Referenced by HetznerCloudServer via StringValueOrRef to assign the IP to a server at creation time. |
| `status.outputs.ip_address` | `string` | The allocated IP address. For IPv4, this is a single address (e.g., "203.0.113.42"). For IPv6, this is the first address in the allocated /64 block (e.g., "2001:db8::1"). |
| `status.outputs.ip_network` | `string` | The allocated IPv6 network in CIDR notation (e.g., "2001:db8::/64"). Empty for IPv4 Primary IPs. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudServer | `spec.publicNet.ipv4` | `status.outputs.primary_ip_id` |
| HetznerCloudServer | `spec.publicNet.ipv6` | `status.outputs.primary_ip_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
