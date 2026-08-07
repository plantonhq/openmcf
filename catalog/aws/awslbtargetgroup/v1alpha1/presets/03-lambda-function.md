# Lambda Function Target

This preset creates a Lambda target group: the ALB invokes the function
directly for each request, so there is no port, protocol, or VPC to
configure. The function is the group's single target, registered statically
in the spec. Reference this group's `target_group_arn` output from an
`AwsLbListener` forward action or an `AwsLbListenerRule` to put a function
behind a path or host route.

## When to Use

- Webhooks, form handlers, or low-traffic API paths served by a function on
  the same ALB (and domain) as container services
- Migrating one route at a time from containers to functions -- a listener
  rule flips a single path to this group
- Any HTTP workload where per-request billing beats always-on capacity

## Key Configuration Choices

- **No `port`, `protocol`, or `vpcId`** -- a Lambda target is invoked, not
  addressed; the spec rejects these fields for `lambda` groups
- **Function by reference** -- the target resolves the `AwsLambda`'s
  `function_arn` output at deploy time, keeping the two components composed
  in the graph
- **`lambdaMultiValueHeadersEnabled: true`** -- repeated headers and query
  parameters arrive as arrays instead of last-value-wins strings; almost
  always what an HTTP handler expects
- **Health checks left at the AWS default (off)** -- an invocation-based
  health check spends function invocations; enable one only if you need the
  ALB to fail over away from a broken function

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<function-name>` | Name for the target group (max 32 chars after truncation) | Usually the function's name |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<lambda-resource-name>` | Name of the AwsLambda resource to invoke | Your AwsLambda manifest's `metadata.name` |

## Common Additions

- The Lambda function needs a resource-based permission allowing
  `elasticloadbalancing.amazonaws.com` to invoke it for this group's ARN
- Enable `healthCheck` (with `enabled: true` and a path) to detect broken
  deployments at the ALB instead of at the client

## Related Presets

- **01-ecs-service-http** -- the container-backed HTTP shape
- **02-nlb-tcp-passthrough** -- Layer-4 pass-through for non-HTTP protocols
