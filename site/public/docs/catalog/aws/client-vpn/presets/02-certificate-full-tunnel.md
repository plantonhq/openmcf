---
title: "Certificate Full-Tunnel VPN"
description: "This preset creates a Client VPN endpoint that routes ALL client traffic through AWS — internet included — for postures where every packet must egress through inspected, NAT-ed infrastructure."
type: "preset"
rank: "02"
presetSlug: "02-certificate-full-tunnel"
componentSlug: "client-vpn"
componentTitle: "Client VPN"
provider: "aws"
icon: "package"
order: 2
---

# Certificate Full-Tunnel VPN

This preset creates a Client VPN endpoint that routes ALL client traffic
through AWS — internet included — for postures where every packet must
egress through inspected, NAT-ed infrastructure.

## When to Use

- Compliance postures requiring all remote-worker traffic to egress
  through corporate infrastructure
- Untrusted client networks (public Wi-Fi) where local breakout is
  unacceptable

## Key Configuration Choices

- **Full tunnel** (`splitTunnel: false`) — all client traffic enters the
  VPN
- **The 0.0.0.0/0 pair** — a route through a NAT-routed subnet AND an
  authorization rule for it; without BOTH, connected clients lose internet
  access (the route carries the traffic, the rule authorizes it)
- **NAT-routed target subnet** — the associated subnet must reach the
  internet through a NAT gateway; the VPN only delivers traffic to the
  subnet

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region code (e.g. `us-west-2`) | Your deployment region |
| `<client-ca-certificate-arn>` | ACM ARN of the client CA chain | ACM console, or an AwsCertManagerCert output |
| `<server-certificate-arn>` | ACM ARN of the server certificate | ACM console, or an AwsCertManagerCert output |
| `<vpc-id>` | The VPC whose security groups apply | Your AwsVpc output |
| `<nat-routed-private-subnet-id>` | Private subnet with a NAT route to the internet | Your AwsSubnet output |
| `<vpn-log-group-name>` | CloudWatch log group for connection events | Your AwsCloudwatchLogGroup output |

## Common Additions

- `clientRouteEnforcementEnabled: true` so managed devices cannot bypass
  the tunnel with client-side route edits
- `dnsServers` pointing at the VPC resolver so DNS also resolves privately

## Related Presets

- **01-certificate-split-tunnel** — VPC-only traffic through the VPN
- **03-saml-sso** — SAML single sign-on with the self-service portal
