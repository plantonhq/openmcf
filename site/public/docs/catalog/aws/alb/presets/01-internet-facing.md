---
title: "Internet-Facing ALB"
description: "This preset creates the standard public entry point: an internet-facing Application Load Balancer across two public subnets, with an explicit security group, deletion protection, and Route53 alias..."
type: "preset"
rank: "01"
presetSlug: "01-internet-facing"
componentSlug: "alb"
componentTitle: "ALB"
provider: "aws"
icon: "package"
order: 1
---

# Internet-Facing ALB

This preset creates the standard public entry point: an internet-facing
Application Load Balancer across two public subnets, with an explicit
security group, deletion protection, and Route53 alias records pointing your
hostnames at it. It carries no listeners by design — attach `AwsLbListener`
resources (typically the HTTPS-forward and HTTP-redirect pair) against its
`load_balancer_arn` output to start accepting traffic.

## When to Use

- The shared front door for public web applications, APIs, or microservices
- The first piece of the composable stack: ALB now, listeners and per-service
  rules as services deploy
- Any architecture where the load balancer should outlive the services
  behind it

## Key Configuration Choices

- **Two public subnets across AZs** — the AWS minimum, and what buys zonal
  redundancy; the spec rejects fewer
- **Explicit security group** — omit it and AWS attaches the VPC default,
  which admits far more than your listener ports; open exactly 80/443 (or
  whatever your listeners will use)
- **Deletion protection** (`deleteProtectionEnabled: true`) — deleting an ALB
  silently orphans every listener and rule attached to it, so production load
  balancers should refuse casual deletion
- **Alias DNS** (`dns.enabled: true`) — alias A records work at the zone
  apex and cost nothing per query, unlike CNAMEs
- **AWS defaults everywhere else** — timeouts, HTTP/2, and desync mitigation
  stay at AWS defaults until you have a reason to move them

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<alb-name>` | Unique name for the ALB (AWS caps it at 32 characters) | Choose a descriptive name (e.g., `main-web`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<public-subnet-id-az1>` | Public subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<public-subnet-id-az2>` | Public subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<alb-security-group-id>` | Security group opening the listener ports (80/443) | AWS EC2 console or `AwsSecurityGroup` status outputs |
| `<route53-hosted-zone-id>` | Route53 hosted zone for your domain | AWS Route53 console or `AwsRoute53Zone` status outputs |
| `<your-domain.com>` | Hostname that should resolve to this ALB | Your DNS plan |

## Common Additions

- Use `valueFrom` references to `AwsSubnet` and `AwsSecurityGroup` resources
  instead of literal IDs so the dependency graph stays explicit
- Add `accessLogs` with an S3 bucket (carrying the ELB log-delivery policy)
  before you need them in an incident
- Set `ipAddressType: dualstack` to serve IPv6 clients, or
  `dualstack-without-public-ipv4` to also avoid public-IPv4 charges
- Raise `idleTimeoutSeconds` above your slowest response time if you serve
  long-running requests

## Related Presets

- **02-internal-hardened** — the VPC-internal variant with header-smuggling
  hardening and access logs
