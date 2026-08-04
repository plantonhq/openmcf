# HetznerCloudFirewall

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1`

HetznerCloudFirewallSpec defines the specification for a Hetzner Cloud firewall.

A firewall is a set of rules that control inbound and outbound network traffic
for servers. When applied to a server (via the server's firewall_ids field),
all inbound traffic not matching any rule is dropped, while outbound traffic
is allowed unless explicitly restricted by outbound rules.

Rules are defined inline. Each rule specifies a direction, protocol, optional
port (required for TCP/UDP), and source or destination CIDR blocks. Hetzner
Cloud supports up to 50 rules per firewall.

The firewall itself does not specify which servers it applies to. Instead,
servers reference firewalls via their firewall_ids field through
StringValueOrRef, keeping the dependency graph unidirectional.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudFirewall
metadata:
  name: hetznercloudfirewall-demo
spec:
  rules:
    - direction: in
      protocol: tcp
      port: "22"
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "allow SSH from anywhere"
    - direction: in
      protocol: tcp
      port: "80"
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "allow HTTP"
    - direction: in
      protocol: tcp
      port: "443"
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "allow HTTPS"
    - direction: in
      protocol: icmp
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "allow ping"
    - direction: out
      protocol: tcp
      port: "any"
      destinationIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "allow all outbound TCP"
    - direction: out
      protocol: udp
      port: "any"
      destinationIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "allow all outbound UDP"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.rules` | `[]Rule` |  |  |  |
| `spec.rules[].direction` | `enum` | yes |  |  |
| `spec.rules[].protocol` | `enum` | yes |  |  |
| `spec.rules[].port` | `string` |  |  |  |
| `spec.rules[].sourceIps` | `[]string` |  |  |  |
| `spec.rules[].destinationIps` | `[]string` |  |  |  |
| `spec.rules[].description` | `string` |  |  |  |

## Field Details

### spec.rules

`[]Rule`

Firewall rules that define allowed traffic.

Hetzner Cloud firewalls are deny-by-default for inbound traffic: when a
firewall is applied to a server, all inbound packets not matching a rule
are dropped. Outbound traffic is allowed by default unless explicitly
restricted by outbound rules.

An empty rules list creates a firewall that blocks all inbound traffic
and allows all outbound traffic when applied to a server.

- rule: port is required when protocol is tcp or udp
- rule: source_ips is required when direction is in
- rule: destination_ips is required when direction is out

### spec.rules[].direction

`enum` · required

Traffic direction this rule applies to.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `direction_unspecified`
- `in`
- `out`

### spec.rules[].protocol

`enum` · required

IP protocol this rule matches.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `protocol_unspecified`
- `icmp`
- `tcp`
- `udp`
- `esp`
- `gre`

### spec.rules[].port

`string`

Port or port range for TCP and UDP rules.

Accepts a single port ("80"), a range ("80-443"), or "any" for all ports.
Required when protocol is tcp or udp. Must not be set for icmp, esp, or gre.

### spec.rules[].sourceIps

`[]string`

CIDR blocks allowed as traffic sources for inbound rules (direction = in).

Use ["0.0.0.0/0", "::/0"] to allow all IPv4 and IPv6 traffic.
Required when direction is in.

### spec.rules[].destinationIps

`[]string`

CIDR blocks allowed as traffic destinations for outbound rules (direction = out).

Use ["0.0.0.0/0", "::/0"] to allow all IPv4 and IPv6 traffic.
Required when direction is out.

### spec.rules[].description

`string`

Human-readable description of what this rule allows.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudFirewall, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.firewall_id` | `string` | The Hetzner Cloud numeric ID of the created firewall (as a string). Referenced by HetznerCloudServer.firewall_ids via StringValueOrRef. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudServer | `spec.firewallIds` | `status.outputs.firewall_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
