# NLB TCP Pass-Through

This preset creates a TCP target group for a Network Load Balancer fronting a
non-HTTP protocol (shown here on 5432, the PostgreSQL port). Traffic passes
through at Layer 4 untouched; the application terminates its own TLS if it
needs encryption. Reference this group's `target_group_arn` output from an
`AwsLbListener` with a `TCP` protocol.

## When to Use

- Databases, message brokers, or custom TCP protocols behind an NLB
- gRPC or WebSocket backends that need Layer-4 pass-through rather than ALB
  termination
- Any long-lived-connection workload where the load balancer must not touch
  the byte stream

## Key Configuration Choices

- **HTTP health check on a TCP group** -- the probe hits a readiness endpoint
  on a dedicated admin port (`8081`) instead of the traffic port; a TCP check
  would only prove the port is open, not that the service can serve
- **`preserveClientIp: true`** -- backends see the real client address in the
  IP header. Note the constraint: a target cannot reliably connect to its own
  group through the NLB with preservation on (hairpinning), and security
  groups on targets must allow the client ranges, not just the NLB
- **`connectionTermination: true`** -- long-lived connections would otherwise
  pin draining targets forever; this closes them when the 120-second drain
  expires
- **Proxy Protocol v2 left off** -- enable `proxyProtocolV2` only when the
  backend is configured to parse the header; enabling it against an unaware
  backend breaks every connection

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<service-name>` | Name for the target group (max 32 chars after truncation) | Your service's name (e.g., `postgres-pool`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<vpc-resource-name>` | Name of the AwsVpc resource the targets live in | Your AwsVpc manifest's `metadata.name` |

## Common Additions

- Set `stickiness` (`type: source_ip`) if clients must keep hitting the same
  backend across connections
- Set `targetHealthState.enableUnhealthyConnectionTermination: false` to let
  established sessions ride out transient health blips
- Change `port`/`protocol` for other wire protocols (`UDP` for DNS/QUIC,
  `TCP_UDP` for both, `TLS` for NLB-side termination)

## Related Presets

- **01-ecs-service-http** -- the ALB-family HTTP shape for containerized services
- **03-lambda-function** -- invoke a Lambda function instead of addressing targets
