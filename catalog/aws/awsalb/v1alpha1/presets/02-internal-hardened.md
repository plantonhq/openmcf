# Internal Hardened ALB

This preset creates a VPC-internal Application Load Balancer with the HTTP
hardening turned all the way up: invalid request headers are dropped, desync
mitigation runs in `strictest` mode, and every request lands in S3 access
logs. Internal traffic is exactly where the strict settings are affordable —
the clients are your own services, so RFC-compliance strictness breaks
nothing while closing the request-smuggling surface. Attach `AwsLbListener`
resources against its `load_balancer_arn` output for ports and routing.

## When to Use

- Service-to-service HTTP routing inside a VPC (the internal API tier)
- Environments with compliance requirements that expect request logging and
  smuggling protection
- Any internal ALB fronting services that themselves sit behind an edge
  proxy — strictness here catches what the edge let through

## Key Configuration Choices

- **Internal scheme** (`internal: true`) — nodes live in private subnets and
  the ALB is reachable only inside the VPC; the scheme is immutable, so this
  is a create-time decision
- **`desyncMitigationMode: strictest`** — blocks every request that is not
  RFC 7230 compliant, not just the ambiguous ones; safe when clients are
  known-good internal services, disruptive when they are arbitrary browsers
- **`dropInvalidHeaderFields: true`** — malformed header names are dropped
  instead of forwarded, the companion hardening to desync mitigation
- **Access logs to S3** — one entry per request; the bucket must carry the
  regional ELB log-delivery bucket policy or delivery fails silently
- **Deletion protection** — an internal ALB accumulates listeners and rules
  from many services; deleting it orphans all of them

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<alb-name>` | Unique name for the ALB (AWS caps it at 32 characters) | Choose a descriptive name (e.g., `internal-api`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<private-subnet-id-az1>` | Private subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az2>` | Private subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<alb-security-group-id>` | Security group admitting internal callers on listener ports | AWS EC2 console or `AwsSecurityGroup` status outputs |
| `<access-logs-bucket-name>` | S3 bucket with the ELB log-delivery bucket policy | AWS S3 console or `AwsS3Bucket` status outputs |
| `<log-prefix>` | Key prefix inside the bucket (e.g., `alb/internal-api`) | Your logging layout |

## Common Additions

- Add `connectionLogs` (same bucket, different prefix) to capture TLS
  handshake failures that never become requests
- Add `healthCheckLogs` when debugging flapping targets
- Add `dns` with a private hosted zone to give the ALB a stable internal
  hostname
- Set `preserveHostHeader: true` if backends generate absolute URLs from the
  Host header

## Related Presets

- **01-internet-facing** — the public entry-point variant with Route53 alias
  DNS
