# QUIC / HTTP/3 Pass-Through

This preset creates a `TCP_QUIC` target group for a Network Load Balancer
serving HTTP/3 traffic: QUIC connections pass through natively while the
same group serves clients that fall back to TCP on port 443. QUIC routes
established connections by connection ID rather than by 5-tuple, so each
static registration carries the `quicServerId` the target presents.
Reference this group's `target_group_arn` output from an `AwsLbListener`
with a `TCP_QUIC` protocol.

## When to Use

- HTTP/3 backends behind an NLB, where clients negotiate QUIC but must be
  able to fall back to TCP on the same port
- Media, gaming, or low-latency APIs that benefit from QUIC's connection
  migration (clients roam networks without dropping the session)
- Pure-QUIC services: change `protocol` to `QUIC` when no TCP fallback is
  needed

## Key Configuration Choices

- **`protocol: TCP_QUIC`** -- one group serves both wire protocols on one
  port, the HTTP/3-with-fallback pattern; the protocol decides the load
  balancer family and is create-only
- **`stickiness: source_ip_dest_ip`** -- flow-hash affinity across both
  address families, so a client arriving over IPv4 and IPv6 keeps hitting
  the same backend; use `source_ip_dest_ip_proto` for the narrowest
  affinity
- **`quicServerId` on each registration** -- ties the target to the QUIC
  endpoint identity it serves, which QUIC's connection-ID routing requires
- **HTTP health check on a QUIC group** -- the probe hits a readiness
  endpoint on a dedicated admin port; QUIC itself is not probeable

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | Name for the target group (max 32 chars after truncation) | Your service's name (e.g., `h3-edge`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<vpc-resource-name>` | Name of the AwsVpc resource the targets live in | Your AwsVpc manifest's `metadata.name` |
| `<target-ip>` | IP address of the QUIC backend inside the VPC | Your service's address |
| `<quic-server-id>` | The QUIC server ID the backend presents | Your HTTP/3 server configuration |

## Common Additions

- Drop the `targets` list entirely when an orchestrator (ECS, EKS,
  auto-scaling) registers targets dynamically -- then set `quicServerId`
  through the registering controller
- Set `targetHealthState.enableUnhealthyConnectionTermination: false` to
  let established QUIC sessions ride out transient health blips

## Related Presets

- **02-nlb-tcp-passthrough** -- the plain TCP Layer-4 shape
- **01-ecs-service-http** -- the ALB-family HTTP shape for containerized services
